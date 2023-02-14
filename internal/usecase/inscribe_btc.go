package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/helpers"
)

type BitcoinTokenMintFee struct {
	Amount       string
	MintFee      string
	SentTokenFee string
}

func calculateMintPrice(input structure.InscribeBtcReceiveAddrRespReq) (*BitcoinTokenMintFee, error) {
	// base64String := input.File
	// base64String = strings.ReplaceAll(base64String, "data:text/html;base64,", "")
	// base64String = strings.ReplaceAll(base64String, "data:image/png;base64,", "")
	// dec, err := base64.StdEncoding.DecodeString(base64String)
	// if err != nil {
	// 	return nil, err
	// }
	fileSize := len([]byte(input.File))
	mintFee := int32(fileSize) / 4 * input.FeeRate

	sentTokenFee := utils.FEE_BTC_SEND_AGV * 2
	totalFee := int(mintFee) + sentTokenFee

	return &BitcoinTokenMintFee{
		Amount:       strconv.FormatInt(int64(totalFee), 10),
		MintFee:      strconv.FormatInt(int64(mintFee), 10),
		SentTokenFee: strconv.FormatInt(int64(sentTokenFee), 10),
	}, nil
}

func (u Usecase) CreateInscribeBTC(rootSpan opentracing.Span, input structure.InscribeBtcReceiveAddrRespReq) (*entity.InscribeBTC, error) {

	span, log := u.StartSpan("CreateInscribeBTC", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("input", input)

	walletAddress := &entity.InscribeBTC{}
	err := copier.Copy(walletAddress, input)
	if err != nil {
		log.Error("u.CreateInscribeBTC.Copy", err.Error(), err)
		return nil, err
	}

	// create wallet name
	userWallet := helpers.CreateBTCOrdWallet(input.WalletAddress)

	log.SetTag(utils.WALLET_ADDRESS_TAG, userWallet)

	// create master wallet:
	resp, err := u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userWallet,
			"wallet",
			"create",
		},
	})

	if err != nil {
		log.Error("u.OrdService.Exec.create.Wallet", err.Error(), err)
		return nil, err
	} else {
		walletAddress.Mnemonic = resp.Stdout
	}

	log.SetData("CreateOrdBTCWalletAddress.createdWallet", resp)
	resp, err = u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userWallet,
			"wallet",
			"receive",
		},
	})

	if err != nil {
		log.Error("u.OrdService.Exec.create.receive", err.Error(), err)
		return nil, err
	}

	// create segwit address
	privKey, _, addressSegwit, err := btc.GenerateAddressSegwit()
	if err != nil {
		log.Error("u.CreateSegwitBTCWalletAddress.GenerateAddressSegwit", err.Error(), err)
		return nil, err
	}
	walletAddress.SegwitKey = privKey
	walletAddress.SegwitAddress = addressSegwit

	log.SetData("CreateInscribeBTC.calculateMintPrice", resp)
	mintFee, err := calculateMintPrice(input)

	if err != nil {
		log.Error("u.CreateSegwitBTCWalletAddress.calculateMintPrice", err.Error(), err)
		return nil, err
	}

	walletAddress.Amount = mintFee.Amount
	walletAddress.MintFee = mintFee.MintFee
	walletAddress.SentTokenFee = mintFee.SentTokenFee
	walletAddress.UserAddress = userWallet // name
	walletAddress.OriginUserAddress = input.WalletAddress
	walletAddress.OrdAddress = strings.ReplaceAll(resp.Stdout, "\n", "")
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = input.File
	walletAddress.InscriptionID = ""
	walletAddress.FeeRate = input.FeeRate
	walletAddress.ExpiredAt = time.Now().Add(time.Hour * 2)

	log.SetTag(userWallet, walletAddress.OrdAddress)

	err = u.Repo.InsertInscribeBTC(walletAddress)
	if err != nil {
		log.Error("u.CreateInscribeBTC.InsertInscribeBTC", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) ListInscribeBTC(rootSpan opentracing.Span, limit, page int64) (*entity.Pagination, error) {
	return u.Repo.ListInscribeBTC(entity.FilterInscribeBT{
		BaseFilters: entity.BaseFilters{Limit: limit, Page: page},
	})
}
func (u Usecase) DetailInscribeBTC(inscriptionID string) (*entity.InscribeBTCResp, error) {
	return u.Repo.FindInscribeBTCByNftID(inscriptionID)
}

// JOBs:
// step 1: job check balance for list inscribe
func (u Usecase) JobInscribeWaitingBalance(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("JobInscribeWaitingBalance", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	_, bs, err := u.buildBTCClient()

	if err != nil {
		go u.trackInscribeHistory("", "JobInscribeWaitingBalance", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error())
		return err
	}
	listPending, _ := u.Repo.ListBTCInscribePending()
	if len(listPending) == 0 {
		// go u.trackInscribeHistory("", "JobInscribeWaitingBalance", "", "", "ListBTCInscribePending", "[]")
		return nil
	}

	for _, item := range listPending {

		// check balance:
		balance, confirm, err := bs.GetBalance(item.SegwitAddress)

		fmt.Println("GetBalance response: ", balance, confirm, err)

		if err != nil {
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "GetBalance - with err", err.Error())
			continue
		}
		if balance == nil {
			err = errors.New("balance is nil")
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "GetBalance", err.Error())
			continue
		}

		if balance.Uint64() == 0 {
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "GetBalance", "0")
			continue
		}

		// get required amount to check vs temp wallet balance:
		amount, ok := big.NewInt(0).SetString(item.Amount, 10)
		if !ok {
			err := errors.New("cannot parse amount")
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "SetString(amount) err", err.Error())
			continue
		}

		if amount.Uint64() == 0 {
			err := errors.New("balance is zero")
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "amount.Uint64() err", err.Error())
			continue
		}

		if balance.Uint64() < amount.Uint64() {
			err := fmt.Errorf("Not enough amount %d < %d ", balance.Uint64(), amount.Uint64())
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "compare balance err", err.Error())

			item.Status = entity.StatusInscribe_NotEnoughBalance
			u.Repo.UpdateBtcInscribe(&item)
			continue
		}

		// received fund:
		item.Status = entity.StatusInscribe_ReceivedFund
		item.IsConfirm = true

		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			fmt.Printf("Could not UpdateBtcInscribe id %s - with err: %v", item.ID, err)
			continue
		}

		go u.trackInscribeHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "Updated StatusInscribe_ReceivedFund", "ok")
		log.SetData(fmt.Sprintf("JobInscribeWaitingBalance.CheckReceiveBTC.%s", item.SegwitAddress), item)
		u.Notify(rootSpan, "JobInscribeWaitingBalance", item.SegwitAddress, fmt.Sprintf("%s received BTC %s from [InscriptionID] %s", item.SegwitAddress, item.Status, item.InscriptionID))

	}

	return nil
}

// step 2: job send all fund from segwit address to ord wallet:
func (u Usecase) JobInscribeSendBTCToOrdWallet(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("JobInscribeSendBTCToOrdWallet", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	_, bs, err := u.buildBTCClient()

	if err != nil {
		go u.trackInscribeHistory("", "JobInscribeSendBTCToOrdWallet", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error())
		return err
	}

	listTosendBtc, _ := u.Repo.ListBTCInscribeByStatus([]entity.StatusInscribe{entity.StatusInscribe_ReceivedFund})
	if len(listTosendBtc) == 0 {
		// go u.trackInscribeHistory("", "JobInscribeSendBTCToOrdWallet", "", "", "ListBTCInscribeByStatus", "[]")
		return nil
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusInscribe_ReceivedFund {

			// send all amount:
			fmt.Println("send all btc from", item.SegwitAddress, "to: ", item.OrdAddress)

			// transfer now:
			txID, err := bs.SendTransactionWithPreferenceFromSegwitAddress(
				item.SegwitKey,
				item.SegwitAddress,
				item.OrdAddress,
				-1,
				btc.PreferenceMedium,
			)
			if err != nil {
				go u.trackInscribeHistory(item.ID.String(), "JobInscribeSendBTCToOrdWallet", item.TableName(), item.Status, "SendTransactionWithPreferenceFromSegwitAddress err", err.Error())
				continue
			}

			item.TxSendBTC = txID
			item.Status = entity.StatusInscribe_SendingBTCFromSegwitAddrToOrdAddr
			// item.ErrCount = 0 // reset error count!
			// TODO: update item
			_, err = u.Repo.UpdateBtcInscribe(&item)
			if err != nil {
				fmt.Printf("Could not UpdateBtcInscribe id %s - with err: %v", item.ID, err)
			}

		}
	}

	return nil
}

// job check 3 tx send: tx user send to temp wallet, tx mint, tx send nft to user
func (u Usecase) JobInscribeCheckTxSend(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("JobInscribeCheckTxSend", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	btcClient, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	// get list sending tx:
	listTosendBtc, _ := u.Repo.ListBTCInscribeByStatus([]entity.StatusInscribe{entity.StatusInscribe_Minting, entity.StatusInscribe_SendingBTCFromSegwitAddrToOrdAddr, entity.StatusInscribe_SendingNFTToUser})
	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {

		statusSuccess := entity.StatusInscribe_Minted
		txHashDb := item.TxMintNft

		if item.Status == entity.StatusInscribe_SendingBTCFromSegwitAddrToOrdAddr {
			statusSuccess = entity.StatusInscribe_SentBTCFromSegwitAddrToOrdAdd
			txHashDb = item.TxSendBTC
		}
		if item.Status == entity.StatusInscribe_SendingNFTToUser {
			statusSuccess = entity.StatusInscribe_SentNFTToUser
			txHashDb = item.TxSendNft
		}
		if item.Status == entity.StatusInscribe_Minting {
			item.IsMinted = true
		}

		txHash, err := chainhash.NewHashFromStr(txHashDb)
		if err != nil {
			fmt.Printf("Could not NewHashFromStr Bitcoin RPCClient - with err: %v", err)
			continue
		}

		txResponse, err := btcClient.GetTransaction(txHash)

		if err == nil {
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeCheckTxSend", item.TableName(), item.Status, "btcClient.txResponse.Confirmations: "+txHashDb, txResponse.Confirmations)
			if txResponse.Confirmations >= 1 {
				// send btc ok now:
				item.Status = statusSuccess
				_, err = u.Repo.UpdateBtcInscribe(&item)
				if err != nil {
					fmt.Printf("Could not JobInscribeCheckTxSend id %s - with err: %v", item.ID, err)
				}
			}
		} else {
			fmt.Printf("Could not GetTransaction Bitcoin RPCClient - with err: %v", err)
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeCheckTxSend", item.TableName(), item.Status, "btcClient.GetTransaction: "+txHashDb, err.Error())

			go u.trackInscribeHistory(item.ID.String(), "JobInscribeCheckTxSend", item.TableName(), item.Status, "bs.CheckTx: "+txHashDb, "Begin check tx via api.")

			// check with api:
			txInfo, err := bs.CheckTx(txHashDb)
			if err != nil {
				fmt.Printf("Could not bs - with err: %v", err)
				go u.trackInscribeHistory(item.ID.String(), "JobInscribeCheckTxSend", item.TableName(), item.Status, "bs.CheckTx: "+txHashDb, err.Error())
			}
			if txInfo.Confirmations >= 1 {
				go u.trackInscribeHistory(item.ID.String(), "JobInscribeCheckTxSend", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txHashDb, txInfo.Confirmations)
				// send nft ok now:
				item.Status = statusSuccess
				item.Success = true
				_, err = u.Repo.UpdateBtcInscribe(&item)
				if err != nil {
					fmt.Printf("Could not UpdateBtcInscribe id %s - with err: %v", item.ID, err)
				}
			}
		}
	}

	return nil
}

// job 4: mint nft:
func (u Usecase) JobInscribeMintNft(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("JobInscribeMintNft", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	listTosendBtc, _ := u.Repo.ListBTCInscribeByStatus([]entity.StatusInscribe{entity.StatusInscribe_SentBTCFromSegwitAddrToOrdAdd})
	if len(listTosendBtc) == 0 {
		// go u.trackInscribeHistory("", "ListBTCInscribeByStatus", "", "", "ListBTCInscribeByStatus", "[]")
		return nil
	}

	for _, item := range listTosendBtc {

		// send all amount:
		fmt.Println("mint nft now ...")

		// - Upload the Animation URL to GCS
		animation := item.FileURI
		animation = strings.ReplaceAll(animation, "data:text/html;base64,", "")
		animation = strings.ReplaceAll(animation, "data:image/png;base64,", "")

		now := time.Now().UTC().Unix()
		uploaded, err := u.GCS.UploadBaseToBucket(animation, fmt.Sprintf("btc-projects/%s/%d.html", item.OrdAddress, now))
		if err != nil {
			log.Error("JobInscribeMintNft.helpers.UploadBaseToBucket.Base64DecodeRaw", err.Error(), err)
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "helpers.BUploadBaseToBucket.ase64DecodeRaw", err.Error())
			continue
		}

		fileURI := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
		item.FileURI = fileURI

		//TODO - enable this
		resp, err := u.OrdService.Mint(ord_service.MintRequest{
			WalletName: item.UserAddress,
			FileUrl:    fileURI,
			FeeRate:    int(item.FeeRate),
			DryRun:     false,
		})

		if err != nil {
			log.Error("OrdService.Mint", err.Error(), err)
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "OrdService.Mint", err.Error())
			continue
		}
		// if not err => update status ok now:
		//TODO: handle log err: Database already open. Cannot acquire lock

		item.Status = entity.StatusInscribe_Minting
		// item.ErrCount = 0 // reset error count!

		item.OutputMintNFT = resp

		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "JobInscribeMintNft.UpdateBtcInscribe", err.Error())
			continue
		}

		tmpText := resp.Stdout
		//tmpText := `{\n  \"commit\": \"7a47732d269d5c005c4df99f2e5cf1e268e217d331d175e445297b1d2991932f\",\n  \"inscription\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2afi0\",\n  \"reveal\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2af\",\n  \"fees\": 185514\n}\n`
		jsonStr := strings.ReplaceAll(tmpText, "\n", "")
		jsonStr = strings.ReplaceAll(jsonStr, "\\", "")

		var btcMintResp ord_service.MintStdoputRespose

		err = json.Unmarshal([]byte(jsonStr), &btcMintResp)
		if err != nil {
			log.Error("BTCMint.helpers.JsonTransform", err.Error(), err)
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "JobInscribeMintNft.Unmarshal(btcMintResp)", err.Error())
			continue
		}

		item.TxMintNft = btcMintResp.Reveal
		item.InscriptionID = btcMintResp.Inscription
		// TODO: update item
		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			fmt.Printf("Could not UpdateBtcInscribe id %s - with err: %v", item.ID, err)
		}

	}

	return nil
}

// job 5: send nft:
// send nft for buy order records:
func (u Usecase) JobInscribeSendNft(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("JobInscribeSendNft", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	// get list buy order status = StatusInscribe_Minted:
	listTosendNft, _ := u.Repo.ListBTCInscribeByStatus([]entity.StatusInscribe{entity.StatusInscribe_Minted})
	if len(listTosendNft) == 0 {
		return nil
	}

	for _, item := range listTosendNft {

		// check nft in master wallet or not:
		listNFTsRep, err := u.GetNftsOwnerOf(span, item.UserAddress)
		if err != nil {
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeSendNft", item.TableName(), item.Status, "GetNftsOwnerOf.Error", err.Error())
			continue
		}

		go u.trackInscribeHistory(item.ID.String(), "JobInscribeSendNft", item.TableName(), item.Status, "GetNftsOwnerOf.listNFTsRep", listNFTsRep)

		// parse nft data:
		var resp []struct {
			Inscription string `json:"inscription"`
			Location    string `json:"location"`
			Explorer    string `json:"explorer"`
		}

		err = json.Unmarshal([]byte(listNFTsRep.Stdout), &resp)
		if err != nil {
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeSendNft", item.TableName(), item.Status, "GetNftsOwnerOf.Unmarshal(listNFTsRep)", err.Error())
			continue
		}
		owner := false
		for _, nft := range resp {
			if strings.EqualFold(nft.Inscription, item.InscriptionID) {
				owner = true
				break
			}

		}
		go u.trackInscribeHistory(item.ID.String(), "JobInscribeSendNft", item.TableName(), item.Status, "GetNftsOwnerOf.CheckNFTOwner", owner)
		if !owner {
			continue
		}

		// transfer now:
		sentTokenResp, err := u.SendTokenByWallet(item.OriginUserAddress, item.InscriptionID, item.UserAddress, int(item.FeeRate))

		go u.trackInscribeHistory(item.ID.String(), "JobInscribeSendNft", item.TableName(), item.Status, "SendTokenByWallet.sentTokenResp", sentTokenResp)

		if err != nil {
			log.Error(fmt.Sprintf("JobInscribeSendNft.SendTokenMKP.%s.Error", item.OrdAddress), err.Error(), err)
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeSendNft", item.TableName(), item.Status, "SendTokenByWallet.err", err.Error())
			continue
		}

		//TODO: handle log err: Database already open. Cannot acquire lock

		// Update status first if none err:
		item.Status = entity.StatusInscribe_SendingNFTToUser
		// item.ErrCount = 0 // reset error count!

		item.OutputSendNFT = sentTokenResp

		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			errPack := fmt.Errorf("Could not UpdateBtcInscribe id %s - with err: %v", item.ID, err.Error())
			log.Error("BtcSendNFTForBuyOrder.helpers.JsonTransform", errPack.Error(), errPack)
			go u.trackInscribeHistory(item.ID.String(), "UpdateBtcInscribe", item.TableName(), item.Status, "SendTokenMKP.UpdateBtcInscribe", err.Error())
			continue
		}

		txResp := sentTokenResp.Stdout
		//txResp := `fd31946b855cbaaf91df4b2c432f9b173e053e65a9879ac909bad028e21b950e\n`
		txResp = strings.ReplaceAll(txResp, "\n", "")

		// update tx:
		item.TxSendNft = txResp
		// item.ErrCount = 0 // reset error count!
		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			errPack := fmt.Errorf("Could not UpdateBtcInscribe id %s - with err: %v", item.ID, err)
			log.Error("UpdateBtcInscribe.UpdateBtcInscribe", errPack.Error(), errPack)
			go u.trackInscribeHistory(item.ID.String(), "UpdateBtcInscribe", item.TableName(), item.Status, "u.UpdateBtcInscribe.UpdateBTCNFTBuyOrder", err.Error())
		}
		// save log:
		log.SetData(fmt.Sprintf("UpdateBtcInscribe.execResp.%s", item.OrdAddress), sentTokenResp)

	}
	return nil
}

func (u Usecase) SendTokenByWallet(receiveAddr, inscriptionID, walletAddressName string, rate int) (*ord_service.ExecRespose, error) {

	sendTokenReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			walletAddressName,
			"wallet",
			"send",
			receiveAddr,
			inscriptionID,
			"--fee-rate",
			fmt.Sprintf("%d", rate),
		}}

	resp, err := u.OrdService.Exec(sendTokenReq)
	return resp, err
}

func (u Usecase) GetNftsOwnerOf(rootSpan opentracing.Span, walletName string) (*ord_service.ExecRespose, error) {
	span, log := u.StartSpan(fmt.Sprintf("GetNftsOwnerOf.%s", "inscriptions"), rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetTag(utils.TOKEN_ID_TAG, "inscriptions")
	log.SetTag(utils.WALLET_ADDRESS_TAG, "ord_marketplace_master")
	listNFTsReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			walletName,
			"wallet",
			"inscriptions",
		}}

	log.SetData("listNFTsReq", listNFTsReq)
	resp, err := u.OrdService.Exec(listNFTsReq)
	defer u.Notify(rootSpan, "GetNftsOwnerOf", "ord_marketplace_master", "inscriptions")
	if err != nil {
		log.SetData("u.OrdService.Exec.Error", err.Error())
		log.Error("u.OrdService.Exec", err.Error(), err)
		return nil, err
	}
	log.SetData("listNFTsRep", resp)
	return resp, err
}

func (u *Usecase) trackInscribeHistory(id, name, table string, status interface{}, requestMsg interface{}, responseMsg interface{}) {
	trackData := &entity.InscribeBTCLogs{
		RecordID:    id,
		Name:        name,
		Table:       table,
		Status:      status,
		RequestMsg:  requestMsg,
		ResponseMsg: responseMsg,
	}
	err := u.Repo.CreateInscribeBTCLog(trackData)
	if err != nil {
		fmt.Printf("trackInscribeHistory.%s.Error:%s", name, err.Error())
	}

}
