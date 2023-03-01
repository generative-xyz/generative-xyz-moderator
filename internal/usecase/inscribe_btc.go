package usecase

import (
	"context"
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
	"go.uber.org/zap"
	"rederinghub.io/external/nfts"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

type BitcoinTokenMintFee struct {
	Amount       string
	MintFee      string
	SentTokenFee string
	Size         int
}

func decodeFileBase64(file string) (string, string, error) {
	i := strings.Index(file, ",")
	if i < 0 {
		return "", "", errors.New("no comma")
	}

	dec, err := base64.StdEncoding.DecodeString(file[i+1:])
	if err != nil {
		return "", "", err
	}
	return string(dec), file[i+1:], nil
}

func calculateMintPrice(input structure.InscribeBtcReceiveAddrRespReq) (*BitcoinTokenMintFee, error) {
	// base64String := input.File
	// base64String = strings.ReplaceAll(base64String, "data:text/html;base64,", "")
	// base64String = strings.ReplaceAll(base64String, "data:image/png;base64,", "")
	// dec, err := base64.StdEncoding.DecodeString(base64String)
	// if err != nil {
	// 	return nil, err
	// }

	// need to encode file: phuong viet lai:
	fileDecode, _, err := decodeFileBase64(input.File)
	if err != nil {
		return nil, err
	}

	fileSize := len([]byte(fileDecode))

	fmt.Println("fileSize===>", fileSize)

	if fileSize < utils.MIN_FILE_SIZE {
		fileSize = utils.MIN_FILE_SIZE
	}
	fmt.Println("new fileSize===>", fileSize)

	mintFee := int32(fileSize) / 4 * input.FeeRate

	fmt.Println("mintFee===>", mintFee)

	sentTokenFee := utils.FEE_BTC_SEND_AGV * 2
	totalFee := int(mintFee) + sentTokenFee

	fmt.Println("total fee ===>", totalFee)

	return &BitcoinTokenMintFee{
		Amount:       strconv.FormatInt(int64(totalFee), 10),
		MintFee:      strconv.FormatInt(int64(mintFee), 10),
		SentTokenFee: strconv.FormatInt(int64(sentTokenFee), 10),
		Size:         fileSize,
	}, nil
}

func (u Usecase) CreateInscribeBTC(ctx context.Context, input structure.InscribeBtcReceiveAddrRespReq) (*entity.InscribeBTC, error) {

	u.Logger.Info("input", input)

	// todo remove:
	// _, base64Str, err := decodeFileBase64(input.File)
	// if err != nil {
	// 	u.Logger.Error("JobInscribeMintNft.decodeFileBase64", err.Error(), err)
	// 	return nil, err
	// }

	// now := time.Now().UTC().Unix()
	// uploaded, err := u.GCS.UploadBaseToBucket(base64Str, fmt.Sprintf("btc-projects/%s/%d.%s", "bc1p3lh2xp8a63rlwpk8zkxrwhhzwqgskfr9el3lmhceu3atyam4rvmshf24vt", now, "txt"))
	// if err != nil {
	// 	return nil, err
	// }

	// fileURI := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	// fmt.Println("fileURI===> ", fileURI)

	// end remove

	walletAddress := &entity.InscribeBTC{}
	err := copier.Copy(walletAddress, input)
	if err != nil {
		u.Logger.Error("u.CreateInscribeBTC.Copy", err.Error(), err)
		return nil, err
	}

	// create wallet name
	userWallet := helpers.CreateBTCOrdWallet(input.WalletAddress)

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
		u.Logger.Error("u.OrdService.Exec.create.Wallet", err.Error(), err)
		return nil, err
	}
	walletAddress.Mnemonic = resp.Stdout

	u.Logger.Info("CreateOrdBTCWalletAddress.createdWallet", resp)
	resp, err = u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userWallet,
			"wallet",
			"receive",
		},
	})

	if err != nil {
		u.Logger.Error("u.OrdService.Exec.create.receive", err.Error(), err)
		return nil, err
	}

	// create segwit address
	privKey, _, addressSegwit, err := btc.GenerateAddressSegwit()
	if err != nil {
		u.Logger.Error("u.CreateSegwitBTCWalletAddress.GenerateAddressSegwit", err.Error(), err)
		return nil, err
	}
	walletAddress.SegwitKey = privKey
	walletAddress.SegwitAddress = addressSegwit

	u.Logger.Info("CreateInscribeBTC.calculateMintPrice", resp)
	mintFee, err := calculateMintPrice(input)

	if err != nil {
		u.Logger.Error("u.CreateSegwitBTCWalletAddress.calculateMintPrice", err.Error(), err)
		return nil, err
	}

	expiredTime := utils.INSCRIBE_TIMEOUT
	if u.Config.ENV == "develop" {
		expiredTime = 1
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
	walletAddress.ExpiredAt = time.Now().Add(time.Hour * time.Duration(expiredTime))
	walletAddress.FileName = input.FileName
	walletAddress.UserUuid = input.UserUuid

	if input.NeedVerifyAuthentic() {
		if nft, err := u.MoralisNft.GetNftByContractAndTokenID(input.TokenAddress, input.TokenId); err == nil {
			logger.AtLog.Logger.Info("MoralisNft.GetNftByContractAndTokenID",
				zap.Any("raw_data", nft))
			walletAddress.IsAuthentic = true
			walletAddress.TokenAddress = nft.TokenAddress
			walletAddress.TokenId = nft.TokenID
		}
	}

	err = u.Repo.InsertInscribeBTC(walletAddress)
	if err != nil {
		u.Logger.Error("u.CreateInscribeBTC.InsertInscribeBTC", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) ListInscribeBTC(req *entity.FilterInscribeBT) (*entity.Pagination, error) {
	return u.Repo.ListInscribeBTC(req)
}

func (u Usecase) DetailInscribeBTC(inscriptionID string) (*entity.InscribeBTCResp, error) {
	return u.Repo.FindInscribeBTCByNftID(inscriptionID)
}

func (u Usecase) RetryInscribeBTC(id string) error {
	item, _ := u.Repo.FindInscribeBTC(id)
	u.Logger.Info("item: ", item, id)
	if item != nil {
		if item.Status == entity.StatusInscribe_NotEnoughBalance {
			item.Status = entity.StatusInscribe_Pending
			_, err := u.Repo.UpdateBtcInscribe(item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// JOBs:
// step 1: job check balance for list inscribe
func (u Usecase) JobInscribeWaitingBalance() error {

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
		u.Logger.Info(fmt.Sprintf("JobInscribeWaitingBalance.CheckReceiveBTC.%s", item.SegwitAddress), item)
		u.Notify("JobInscribeWaitingBalance", item.SegwitAddress, fmt.Sprintf("%s received BTC %d from [InscriptionID] %s", item.SegwitAddress, item.Status, item.InscriptionID))

	}

	return nil
}

// step 2: job send all fund from segwit address to ord wallet:
func (u Usecase) JobInscribeSendBTCToOrdWallet() error {

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
				btc.PreferenceHigh,
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
func (u Usecase) JobInscribeCheckTxSend() error {

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
				item.IsSuccess = statusSuccess == entity.StatusInscribe_SentNFTToUser
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
func (u Usecase) JobInscribeMintNft() error {

	listTosendBtc, _ := u.Repo.ListBTCInscribeByStatus([]entity.StatusInscribe{entity.StatusInscribe_SentBTCFromSegwitAddrToOrdAdd})
	if len(listTosendBtc) == 0 {
		// go u.trackInscribeHistory("", "ListBTCInscribeByStatus", "", "", "ListBTCInscribeByStatus", "[]")
		return nil
	}

	for _, item := range listTosendBtc {

		// send all amount:
		fmt.Println("mint nft now ...", item.FileName)

		// - Upload the Animation URL to GCS
		typeFile := ""

		if len(item.FileName) == 0 {
			err := errors.New("File name invalid")
			u.Logger.Error("JobInscribeMintNft.len(Filename)", err.Error(), err)
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "CheckFileName", err.Error())
			continue
		}

		typeFiles := strings.Split(item.FileName, ".")
		if len(typeFiles) < 2 {
			err := errors.New("File name invalid")
			u.Logger.Error("JobInscribeMintNft.len(Filename)", err.Error(), err)
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "CheckFileName", err.Error())
			continue
		}

		typeFile = typeFiles[len(typeFiles)-1]
		fmt.Println("typeFile: ", typeFile)

		// update google clound: TODO need to move into api to avoid create file many time.
		_, base64Str, err := decodeFileBase64(item.FileURI)
		if err != nil {
			u.Logger.Error("JobInscribeMintNft.decodeFileBase64", err.Error(), err)
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "helpers.decodeFileBase64", err.Error())
			continue
		}

		now := time.Now().UTC().Unix()
		uploaded, err := u.GCS.UploadBaseToBucket(base64Str, fmt.Sprintf("btc-projects/%s/%d.%s", item.OrdAddress, now, typeFile))
		if err != nil {
			u.Logger.Error("JobInscribeMintNft.helpers.UploadBaseToBucket.Base64DecodeRaw", err.Error(), err)
			go u.trackInscribeHistory(item.ID.String(), "JobInscribeMintNft", item.TableName(), item.Status, "helpers.BUploadBaseToBucket.ase64DecodeRaw", err.Error())
			continue
		}
		item.LocalLink = uploaded.FullPath

		fileURI := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
		item.FileURI = fileURI

		//TODO - enable this
		resp, err := u.OrdService.Mint(ord_service.MintRequest{
			WalletName:        item.UserAddress,
			FileUrl:           fileURI,
			FeeRate:           int(item.FeeRate),
			DryRun:            false,
			AutoFeeRateSelect: false,
		})

		if err != nil {
			u.Logger.Error("OrdService.Mint", err.Error(), err)
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
			u.Logger.Error("BTCMint.helpers.JsonTransform", err.Error(), err)
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
func (u Usecase) JobInscribeSendNft() error {

	// get list buy order status = StatusInscribe_Minted:
	listTosendNft, _ := u.Repo.ListBTCInscribeByStatus([]entity.StatusInscribe{entity.StatusInscribe_Minted})
	if len(listTosendNft) == 0 {
		return nil
	}

	for _, item := range listTosendNft {

		// check nft in master wallet or not:
		listNFTsRep, err := u.GetNftsOwnerOf(item.UserAddress)
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
			u.Logger.Error(fmt.Sprintf("JobInscribeSendNft.SendTokenMKP.%s.Error", item.OrdAddress), err.Error(), err)
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
			u.Logger.Error("BtcSendNFTForBuyOrder.helpers.JsonTransform", errPack.Error(), errPack)
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
			u.Logger.Error("UpdateBtcInscribe.UpdateBtcInscribe", errPack.Error(), errPack)
			go u.trackInscribeHistory(item.ID.String(), "UpdateBtcInscribe", item.TableName(), item.Status, "u.UpdateBtcInscribe.UpdateBTCNFTBuyOrder", err.Error())
		}
		// save log:
		u.Logger.Info(fmt.Sprintf("UpdateBtcInscribe.execResp.%s", item.OrdAddress), sentTokenResp)

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

func (u Usecase) GetNftsOwnerOf(walletName string) (*ord_service.ExecRespose, error) {

	listNFTsReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			walletName,
			"wallet",
			"inscriptions",
		}}

	u.Logger.Info("listNFTsReq", listNFTsReq)
	resp, err := u.OrdService.Exec(listNFTsReq)
	defer u.Notify("GetNftsOwnerOf", "ord_marketplace_master", "inscriptions")
	if err != nil {
		u.Logger.Info("u.OrdService.Exec.Error", err.Error())
		u.Logger.Error("u.OrdService.Exec", err.Error(), err)
		return nil, err
	}
	u.Logger.Info("listNFTsRep", resp)
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

func (u Usecase) ApiCheckListTempAddress() error {
	var autoGenerated []struct {
		SegwitAddress string `json:"segwit_address"`
	}
	listBtc := `[{}]`

	err := json.Unmarshal([]byte(listBtc), &autoGenerated)
	if err != nil {
		fmt.Println("err")
		return nil
	}

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}
	fmt.Println("len(autoGenerated)", len((autoGenerated)))

	for _, btc := range autoGenerated {

		fmt.Println("check address: ", btc.SegwitAddress)

		balance, confirm, err := bs.GetBalance(btc.SegwitAddress)

		fmt.Println("GetBalance response: ", balance, confirm, err)

		if err != nil {
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			continue
		}
		if balance == nil {
			err = errors.New("balance is nil")
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			continue
		}
		if balance.Uint64() > 0 {
			fmt.Println("Balance OK now====>", btc.SegwitAddress)
		}
		time.Sleep(time.Second * 1)

	}

	return nil
}

func (u Usecase) ListNftFromMoralis(ctx context.Context, userWallet, delegateWallet string, pag *entity.Pagination) (map[string]*entity.Pagination, error) {
	var (
		pageSize              = int(pag.PageSize)
		cursor        *string = nil
		resp                  = make(map[string]*entity.Pagination)
		walletAddress string
	)
	if len(pag.Cursor) > 0 {
		cursor = &pag.Cursor
	}
	reqMoralisFilter := nfts.MoralisFilter{
		Limit:  &pageSize,
		Cursor: cursor,
	}

	if delegateWallet == "" {
		delegations, err := u.DelegateService.GetDelegationsByDelegate(ctx, userWallet)
		if err != nil {
			return nil, err
		}
		if len(delegations) > 0 {
			for i := range delegations {
				delegateWalletAddress := delegations[i].Contract.String()
				resp[delegateWalletAddress] = &entity.Pagination{
					Page:     pag.Page,
					PageSize: pag.PageSize,
				}
				nfts, err := u.MoralisNft.GetNftByWalletAddress(delegateWalletAddress, reqMoralisFilter)
				if err != nil {
					return nil, err
				}
				resp[delegateWalletAddress].Result = nfts.Result
				resp[delegateWalletAddress].Total = int64(nfts.Total)
				resp[delegateWalletAddress].SetTotalPage()
			}
		} else {
			walletAddress = userWallet
		}
	} else {
		walletAddress = delegateWallet
	}

	if walletAddress != "" {
		resp[walletAddress] = pag
		nfts, err := u.MoralisNft.GetNftByWalletAddress(walletAddress, reqMoralisFilter)
		if err != nil {
			return nil, err
		}
		pag.Result = nfts.Result
		pag.Total = int64(nfts.Total)
		pag.SetTotalPage()
	}

	return resp, nil
}
