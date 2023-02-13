package usecase

import (
	"encoding/base64"
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

const (
	SENT_TOKEN_FEE = 0
)

func calculateMintPrice(input structure.BctWalletAddressDataV2) BitcoinTokenMintFee {
	base64String := input.File
	base64String = strings.ReplaceAll(base64String, "data:text/html;base64,", "")
	base64String = strings.ReplaceAll(base64String, "data:image/png;base64,", "")
	dec, _ := base64.StdEncoding.DecodeString(base64String)
	fileSize := len([]byte(dec))
	mintFee := int32(fileSize) / 4 * input.FeeRate
	return BitcoinTokenMintFee{
		Amount:       strconv.FormatInt(int64(mintFee+SENT_TOKEN_FEE), 10),
		MintFee:      strconv.FormatInt(int64(mintFee), 10),
		SentTokenFee: strconv.FormatInt(int64(SENT_TOKEN_FEE), 10),
	}
}

func (u Usecase) CreateBTCWalletAddressV2(rootSpan opentracing.Span, input structure.BctWalletAddressDataV2) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("CreateBTCWalletAddressV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("input", input)

	walletAddress := &entity.BTCWalletAddressV2{}
	err := copier.Copy(walletAddress, input)
	if err != nil {
		log.Error("u.CreateBTCWalletAddressV2.Copy", err.Error(), err)
		return nil, err
	}

	userWallet := helpers.CreateBTCOrdWallet(input.WalletAddress)
	log.SetTag(utils.WALLET_ADDRESS_TAG, userWallet)
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
		//return nil, err
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
	walletAddress.Mnemonic = privKey
	walletAddress.SegwitAddress = addressSegwit

	log.SetData("CreateBTCWalletAddressV2.receive", resp)
	mintFee := calculateMintPrice(input)
	walletAddress.Amount = mintFee.Amount
	walletAddress.MintFee = mintFee.MintFee
	walletAddress.SentTokenFee = mintFee.SentTokenFee
	walletAddress.UserAddress = userWallet
	walletAddress.OriginUserAddress = input.WalletAddress
	walletAddress.OrdAddress = strings.ReplaceAll(resp.Stdout, "\n", "")
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = input.File
	walletAddress.InscriptionID = ""
	walletAddress.FeeRate = input.FeeRate

	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, walletAddress.OrdAddress)
	err = u.Repo.InsertBtcWalletAddressV2(walletAddress)
	if err != nil {
		log.Error("u.CreateBTCWalletAddressV2.InsertBtcWalletAddressV2", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

// step 1: job check balance for list inscribe
func (u Usecase) JobInscribeWaitingBalance(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("JobInscribeWaitingBalance", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	_, bs, err := u.buildBTCClient()

	if err != nil {
		go u.trackHistory("", "JobInscribeWaitingBalance", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error())
		return err
	}
	listPending, _ := u.Repo.ListBTCInscribePending()
	if len(listPending) == 0 {
		// go u.trackHistory("", "JobInscribeWaitingBalance", "", "", "ListBTCInscribePending", "[]")
		return nil
	}

	for _, item := range listPending {

		// check balance:
		balance, confirm, err := bs.GetBalance(item.SegwitAddress)

		fmt.Println("GetBalance response: ", balance, confirm, err)

		if err != nil {
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			go u.trackHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "GetBalance - with err", err.Error())
			continue
		}
		if balance == nil {
			err = errors.New("balance is nil")
			go u.trackHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "GetBalance", err.Error())
			continue
		}

		if balance.Uint64() == 0 {
			go u.trackHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "GetBalance", "0")
			continue
		}

		// get required amount to check vs temp wallet balance:
		amount, _ := big.NewInt(0).SetString(item.Amount, 10)

		if amount.Uint64() == 0 {
			err := errors.New("balance is zero")
			go u.trackHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "amount.Uint64() err", err.Error())
			continue
		}

		if r := balance.Cmp(amount); r == -1 {
			err := errors.New("Not enough amount")
			go u.trackHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "Receive balance err", err.Error())

			item.Status = entity.StatusInscribe_NotEnoughBalance
			u.Repo.UpdateBtcInscribe(&item)
			continue
		}

		// received fund:
		item.Status = entity.StatusInscribe_ReceivedFund

		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			fmt.Printf("Could not UpdateBtcInscribe id %s - with err: %v", item.ID, err)
			continue
		}

		go u.trackHistory(item.ID.String(), "JobInscribeWaitingBalance", item.TableName(), item.Status, "Updated StatusInscribe_ReceivedFund", "ok")
		log.SetData(fmt.Sprintf("JobInscribeWaitingBalance.CheckReceiveBTC.%s", item.SegwitAddress), item)
		u.Notify(rootSpan, "JobInscribeWaitingBalance", item.SegwitAddress, fmt.Sprintf("%s received BTC %s from [InscriptionID] %s", item.SegwitAddress, item.ReceivedBalance, item.InscriptionID))

	}

	return nil
}

// step 2: job send all fund from segwit address to ord wallet:
func (u Usecase) JobInscribeSendBTCToOrdWallet(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("JobInscribeSendBTCToOrdWallet", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	_, bs, err := u.buildBTCClient()

	if err != nil {
		go u.trackHistory("", "JobInscribeSendBTCToOrdWallet", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error())
		return err
	}

	listTosendBtc, _ := u.Repo.ListBTCInscribeByStatus([]entity.StatusInscribe{entity.StatusInscribe_ReceivedFund})
	if len(listTosendBtc) == 0 {
		// go u.trackHistory("", "JobInscribeSendBTCToOrdWallet", "", "", "ListBTCInscribeByStatus", "[]")
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
				go u.trackHistory(item.ID.String(), "JobInscribeSendBTCToOrdWallet", item.TableName(), item.Status, "SendTransactionWithPreferenceFromSegwitAddress err", err.Error())
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

		txHash, err := chainhash.NewHashFromStr(item.TxSendBTC)
		if err != nil {
			fmt.Printf("Could not NewHashFromStr Bitcoin RPCClient - with err: %v", err)
			continue
		}

		txResponse, err := btcClient.GetTransaction(txHash)

		statusSuccess := entity.StatusInscribe_Minted
		if item.Status == entity.StatusInscribe_SendingBTCFromSegwitAddrToOrdAddr {
			statusSuccess = entity.StatusInscribe_SentBTCFromSegwitAddrToOrdAdd
		}
		if item.Status == entity.StatusInscribe_SendingNFTToUser {
			statusSuccess = entity.StatusInscribe_SentNFTToUser
		}

		if err == nil {
			go u.trackHistory(item.ID.String(), "JobInscribeCheckTxSend", item.TableName(), item.Status, "btcClient.txResponse.Confirmations: "+item.TxSendBTC, txResponse.Confirmations)
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
			go u.trackHistory(item.ID.String(), "JobInscribeCheckTxSend", item.TableName(), item.Status, "btcClient.GetTransaction: "+item.TxSendBTC, err.Error())

			go u.trackHistory(item.ID.String(), "JobInscribeCheckTxSend", item.TableName(), item.Status, "bs.CheckTx: "+item.TxSendBTC, "Begin check tx via api.")

			// check with api:
			txInfo, err := bs.CheckTx(item.TxSendBTC)
			if err != nil {
				fmt.Printf("Could not bs - with err: %v", err)
				go u.trackHistory(item.ID.String(), "JobInscribeCheckTxSend", item.TableName(), item.Status, "bs.CheckTx: "+item.TxSendBTC, err.Error())
			}
			if txInfo.Confirmations >= 1 {
				go u.trackHistory(item.ID.String(), "JobInscribeCheckTxSend", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+item.TxSendBTC, txInfo.Confirmations)
				// send nft ok now:
				item.Status = statusSuccess
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
		// go u.trackHistory("", "ListBTCInscribeByStatus", "", "", "ListBTCInscribeByStatus", "[]")
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
			go u.trackHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "helpers.BUploadBaseToBucket.ase64DecodeRaw", err.Error())
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
			go u.trackHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "OrdService.Mint", err.Error())
			continue
		}
		// if not err => update status ok now:
		//TODO: handle log err: Database already open. Cannot acquire lock

		item.Status = entity.StatusInscribe_Minting
		// item.ErrCount = 0 // reset error count!

		item.OutputMintNFT = resp

		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			go u.trackHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "JobInscribeMintNft.UpdateBtcInscribe", err.Error())
			continue
		}

		tmpText := resp.Stdout
		//tmpText := `{\n  \"commit\": \"7a47732d269d5c005c4df99f2e5cf1e268e217d331d175e445297b1d2991932f\",\n  \"inscription\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2afi0\",\n  \"reveal\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2af\",\n  \"fees\": 185514\n}\n`
		jsonStr := strings.ReplaceAll(tmpText, "\n", "")
		jsonStr = strings.ReplaceAll(jsonStr, "\\", "")

		var btcMintResp ord_service.MintStdoputRespose

		err = json.Unmarshal([]byte(jsonStr), btcMintResp)
		if err != nil {
			log.Error("BTCMint.helpers.JsonTransform", err.Error(), err)
			go u.trackHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "JobInscribeMintNft.Unmarshal(btcMintResp)", err.Error())
			continue
		}

		item.TxMintNft = btcMintResp.Reveal
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
			go u.trackHistory(item.ID.String(), "JobInscribeSendNft", item.TableName(), item.Status, "GetNftsOwnerOf.Error", err.Error())
			continue
		}

		go u.trackHistory(item.ID.String(), "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "GetNftsOwnerOf.listNFTsRep", listNFTsRep)

		// parse nft data:
		var resp []struct {
			Inscription string `json:"inscription"`
			Location    string `json:"location"`
			Explorer    string `json:"explorer"`
		}

		err = json.Unmarshal([]byte(listNFTsRep.Stdout), &resp)
		if err != nil {
			go u.trackHistory(item.ID.String(), "JobInscribeSendNft", item.TableName(), item.Status, "GetNftsOwnerOf.Unmarshal(listNFTsRep)", err.Error())
			continue
		}
		owner := false
		for _, nft := range resp {
			if strings.EqualFold(nft.Inscription, item.InscriptionID) {
				owner = true
				break
			}

		}
		go u.trackHistory(item.ID.String(), "JobInscribeSendNft", item.TableName(), item.Status, "GetNftsOwnerOf.CheckNFTOwner", owner)
		if !owner {
			continue
		}

		// transfer now:
		sentTokenResp, err := u.SendTokenByWallet(item.OrdAddress, item.InscriptionID, item.UserAddress)

		go u.trackHistory(item.ID.String(), "JobInscribeSendNft", item.TableName(), item.Status, "SendTokenMKP.sentTokenResp", sentTokenResp)

		if err != nil {
			log.Error(fmt.Sprintf("JobInscribeSendNft.SendTokenMKP.%s.Error", item.OrdAddress), err.Error(), err)
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
			go u.trackHistory(item.ID.String(), "UpdateBtcInscribe", item.TableName(), item.Status, "SendTokenMKP.UpdateBtcInscribe", err.Error())
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
			go u.trackHistory(item.ID.String(), "UpdateBtcInscribe", item.TableName(), item.Status, "u.UpdateBtcInscribe.UpdateBTCNFTBuyOrder", err.Error())
		}
		// save log:
		log.SetData(fmt.Sprintf("UpdateBtcInscribe.execResp.%s", item.OrdAddress), sentTokenResp)

	}
	return nil
}

func (u Usecase) SendTokenByWallet(receiveAddr, inscriptionID, walletAddressName string) (*ord_service.ExecRespose, error) {

	sendTokenReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			walletAddressName,
			"wallet",
			"send",
			receiveAddr,
			inscriptionID,
			"--fee-rate",
			"15",
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

func (u Usecase) CheckbalanceWalletAddressV2(rootSpan opentracing.Span, input structure.CheckBalance) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("CheckbalanceWalletAddressV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, input.Address)
	btc, err := u.Repo.FindBtcWalletAddressByOrdV2(input.Address)
	if err != nil {
		log.Error("u.Repo.FindBtcWalletAddressByOrd", err.Error(), err)
		return nil, err
	}

	balance, err := u.CheckBalanceV2(span, *btc)
	if err != nil {
		log.Error("u.BalanceLogic", err.Error(), err)
		return nil, err
	}

	return balance, nil
}

// deprecated
// change this to mint from ord address
func (u Usecase) BTCMintV2(rootSpan opentracing.Span, input structure.BctMintData) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("BTCMintV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("input", input)
	log.SetTag(utils.WALLET_ADDRESS_TAG, input.Address)

	btc, err := u.Repo.FindBtcWalletAddressByOrdV2(input.Address)
	if err != nil {
		log.Error("BTCMint.FindBtcWalletAddressByOrd", err.Error(), err)
		return nil, err
	}

	//mint logic
	btc, err = u.MintLogicV2(span, btc)
	if err != nil {
		log.Error("BTCMint.MintLogic", err.Error(), err)
		return nil, err
	}

	// - Upload the Animation URL to GCS
	animation := btc.FileURI
	animation = strings.ReplaceAll(animation, "data:text/html;base64,", "")
	animation = strings.ReplaceAll(animation, "data:image/png;base64,", "")

	now := time.Now().UTC().Unix()
	uploaded, err := u.GCS.UploadBaseToBucket(animation, fmt.Sprintf("btc-projects/%s/%d.html", btc.OrdAddress, now))
	if err != nil {
		log.Error("BTCMint.helpers.Base64DecodeRaw", err.Error(), err)
		return nil, err
	}

	fileURI := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	btc.FileURI = fileURI

	//TODO - enable this
	resp, err := u.OrdService.Mint(ord_service.MintRequest{
		WalletName: "ord_master",
		FileUrl:    fileURI,
		FeeRate:    int(btc.FeeRate), //temp
		DryRun:     false,
	})

	if err != nil {
		log.Error("BTCMint.Mint", err.Error(), err)
		return nil, err
	}

	tmpText := resp.Stdout
	//tmpText := `{\n  \"commit\": \"7a47732d269d5c005c4df99f2e5cf1e268e217d331d175e445297b1d2991932f\",\n  \"inscription\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2afi0\",\n  \"reveal\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2af\",\n  \"fees\": 185514\n}\n`
	jsonStr := strings.ReplaceAll(tmpText, "\n", "")
	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")
	btcMintResp := &ord_service.MintStdoputRespose{}

	bytes := []byte(jsonStr)
	err = json.Unmarshal(bytes, btcMintResp)
	if err != nil {
		log.Error("BTCMint.helpers.JsonTransform", err.Error(), err)
		return nil, err
	}

	btc.MintResponse = entity.MintStdoputResponse(*btcMintResp)
	btc, err = u.UpdateBtcMintedStatusV2(span, btc)
	if err != nil {
		log.Error("BTCMint.UpdateBtcMintedStatus", err.Error(), err)
		return nil, err
	}

	return btc, nil
}

func (u Usecase) ReadGCSFolderV2(rootSpan opentracing.Span, input structure.BctWalletAddressData) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("ReadGCSFolderV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("input", input)
	u.GCS.ReadFolder("btc-projects/project-1/")
	return nil, nil
}

// deprecated
func (u Usecase) UpdateBtcMintedStatusV2(rootSpan opentracing.Span, btcWallet *entity.BTCWalletAddressV2) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("UpdateBtcMintedStatusV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("input", btcWallet)

	btcWallet.IsMinted = true
	log.SetTag(utils.WALLET_ADDRESS_TAG, btcWallet.UserAddress)
	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, btcWallet.OrdAddress)
	log.SetTag(utils.TOKEN_ID_TAG, btcWallet.InscriptionID)

	updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddrV2(btcWallet.OrdAddress, btcWallet)
	if err != nil {
		log.Error("BTCMint.helpers.UpdateBtcWalletAddressByOrdAddr", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return btcWallet, nil
}

// deprecated
func (u Usecase) CheckBalanceV2(rootSpan opentracing.Span, btc entity.BTCWalletAddressV2) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("CheckBlanceV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetTag(utils.WALLET_ADDRESS_TAG, btc.UserAddress)
	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, btc.OrdAddress)

	balanceRequest := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			btc.UserAddress,
			"wallet",
			"balance",
		},
	}

	log.SetData("balanceRequest", balanceRequest)
	//userWallet := helpers.CreateBTCOrdWallet(btc.UserAddress)
	resp, err := u.OrdService.Exec(balanceRequest)
	if err != nil {
		log.Error("BTCMint.Exec.balance", err.Error(), err)
		return nil, err
	}

	log.SetData("balanceResponse", resp)
	balance := strings.ReplaceAll(resp.Stdout, "\n", "")
	log.SetData("balance", balance)

	btc.Balance = balance

	go func(rootSpan opentracing.Span, balance *entity.BTCWalletAddressV2) {
		span, log := u.StartSpan("CheckBalance.RoutineUpdate", rootSpan)
		defer u.Tracer.FinishSpan(span, log)

		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddrV2(balance.OrdAddress, balance)
		if err != nil {
			log.Error("u.Repo.UpdateBtcWalletAddressByOrdAddr", err.Error(), err)
			return
		}
		log.SetData("updated", updated)

	}(span, &btc)

	return &btc, nil
}

// deprecated
func (u Usecase) BalanceLogicV2(rootSpan opentracing.Span, btc entity.BTCWalletAddressV2) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("BalanceLogicV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	balance, err := u.CheckBalanceV2(span, btc)
	if err != nil {
		log.Error("u.CheckBalance", err.Error(), err)
		return nil, err
	}

	//TODO logic of the checked balance here
	if balance.Balance < btc.Amount {
		err := errors.New("Not enough amount")
		return nil, err
	}
	btc.IsConfirm = true
	return &btc, nil
}

// deprecated
func (u Usecase) MintLogicV2(rootSpan opentracing.Span, btc *entity.BTCWalletAddressV2) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("MintLogicV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	var err error

	log.SetTag(utils.WALLET_ADDRESS_TAG, btc.UserAddress)
	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, btc.OrdAddress)

	//if this was minted, skip it
	if btc.IsMinted {
		err = errors.New("This btc was minted")
		log.Error("BTCMint.Minted", err.Error(), err)
		return nil, err
	}

	if !btc.IsConfirm {
		err = errors.New("This btc must be IsConfirmed")
		log.Error("BTCMint.IsConfirmed", err.Error(), err)
		return nil, err
	}

	log.SetData("btc", btc)
	return btc, nil
}

// deprecated
func (u Usecase) WaitingForBalancingV2(rootSpan opentracing.Span) ([]entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("WaitingForBalancingV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	addreses, err := u.Repo.ListProcessingWalletAddressV2()
	if err != nil {
		log.Error("WillBeProcessWTC.ListProcessingWalletAddress", err.Error(), err)
		return nil, err
	}

	for _, item := range addreses {
		log.SetTag(utils.WALLET_ADDRESS_TAG, item.UserAddress)
		log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, item.OrdAddress)
		newItem, err := u.BalanceLogicV2(span, item)
		if err != nil {
			//log.Error(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s.Error", item.OrdAddress), err.Error(), err)
			continue
		}
		log.SetData(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s", item.OrdAddress), newItem)

		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddrV2(item.OrdAddress, newItem)
		if err != nil {
			log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
			continue
		}
		log.SetData("updated", updated)

		btc, err := u.BTCMintV2(span, structure.BctMintData{Address: newItem.OrdAddress})
		if err != nil {
			log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", newItem.OrdAddress), err.Error(), err)
			continue
		}

		log.SetData("btc.Minted", btc)
	}

	return nil, nil
}

// func (u Usecase) WaitingForMintedV2(rootSpan opentracing.Span) ([]entity.BTCWalletAddressV2, error) {
// 	span, log := u.StartSpan("WaitingForMintedV2", rootSpan)
// 	defer u.Tracer.FinishSpan(span, log)

// 	addreses, err := u.Repo.ListBTCAddressV2()
// 	if err != nil {
// 		log.Error("WillBeProcessWTC.ListBTCAddress", err.Error(), err)
// 		return nil, err
// 	}

// 	for _, item := range addreses {
// 		log.SetTag(utils.WALLET_ADDRESS_TAG, item.UserAddress)
// 		log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, item.OrdAddress)

// 		addr := item.OriginUserAddress
// 		if addr == "" {
// 			addr = item.UserAddress
// 		}

// 		sentTokenResp, err := u.SendToken(rootSpan, addr, item.MintResponse.Inscription)
// 		if err != nil {
// 			log.Error(fmt.Sprintf("ListenTheMintedBTC.sentToken.%s.Error", item.OrdAddress), err.Error(), err)
// 			continue
// 		}

// 		log.SetData(fmt.Sprintf("ListenTheMintedBTC.execResp.%s", item.OrdAddress), sentTokenResp)
// 		// amout, err := strconv.ParseFloat(item.Amount, 10)
// 		// if err != nil {
// 		// 	log.Error("ListenTheMintedBTC.%s. strconv.ParseFloa.Error", err.Error(), err)
// 		// 	continue
// 		// }
// 		// amout = amout * 0.9

// 		// fundData := ord_service.ExecRequest{
// 		// 	Args: []string{
// 		// 		"--wallet",
// 		// 		item.OrdAddress,
// 		// 		"send",
// 		// 		"ord_master",
// 		// 		fmt.Sprintf("%f", amout),
// 		// 		"--fee-rate",
// 		// 		"15",
// 		// 	},
// 		// }

// 		// log.SetData("fundData", fundData)
// 		// fundResp, err := u.OrdService.Exec(fundData)

// 		// if err != nil {
// 		// 	log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.ReFund.Error", item.OrdAddress), err.Error(), err)
// 		// 	continue
// 		// }

// 		// log.SetData("fundResp", fundResp)

// 		item.MintResponse.IsSent = true
// 		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddrV2(item.OrdAddress, &item)
// 		if err != nil {
// 			log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateBtcWalletAddressByOrdAddr.Error", item.OrdAddress), err.Error(), err)
// 			continue
// 		}
// 		log.SetData("updated", updated)

// 		//TODO: - create bitcoin token here
// 		// _, err = u.CreateBTCTokenURI(span, item.ProjectID, item.MintResponse.Inscription)
// 		// if err != nil {
// 		// 	log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.CreateBTCTokenURI.Error", item.OrdAddress), err.Error(), err)
// 		// 	continue
// 		// }
// 	}

// 	return nil, nil
// }

// func (u Usecase) SendTokenV2(rootSpan opentracing.Span, receiveAddr string, inscriptionID string) (*ord_service.ExecRespose, error) {
// 	span, log := u.StartSpan("SendTokenV2", rootSpan)
// 	defer u.Tracer.FinishSpan(span, log)

// 	log.SetData(utils.TOKEN_ID_TAG, inscriptionID)
// 	sendTokenReq := ord_service.ExecRequest{
// 		Args: []string{
// 			"--wallet",
// 			"ord_master",
// 			"wallet",
// 			"send",
// 			receiveAddr,
// 			inscriptionID,
// 			"--fee-rate",
// 			"15",
// 		}}

// 	log.SetData("sendTokenReq", sendTokenReq)
// 	resp, err := u.OrdService.Exec(sendTokenReq)

// 	if err != nil {
// 		log.Error("u.OrdService.Exec", err.Error(), err)
// 		return nil, err
// 	}

// 	log.SetData("sendTokenRes", resp)
// 	return resp, err

// }

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
