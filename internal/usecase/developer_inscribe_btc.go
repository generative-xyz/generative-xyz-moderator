package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

func (u Usecase) DeveloperCreateInscribe(ctx context.Context, input structure.InscribeBtcReceiveAddrRespReq) (*entity.DeveloperInscribe, error) {

	u.Logger.Info("input", input)

	walletAddress := &entity.DeveloperInscribe{}
	err := copier.Copy(walletAddress, input)
	if err != nil {
		u.Logger.Error("u.CreateDeveloperInscribeBTC.Copy", err.Error(), err)
		return nil, err
	}

	// create wallet name
	userWallet := helpers.CreateBTCOrdWallet(input.WalletAddress)

	// create master wallet:
	resp, err := u.OrdServiceDeveloper.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userWallet,
			"wallet",
			"create",
		},
	})

	if err != nil {
		u.Logger.Error("u.OrdServiceDeveloper.Exec.create.Wallet", err.Error(), err)
		return nil, err
	}
	walletAddress.Mnemonic = resp.Stdout

	u.Logger.Info("DeveloperCreateInscribe.createdWallet", resp)
	resp, err = u.OrdServiceDeveloper.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userWallet,
			"wallet",
			"receive",
		},
	})

	if err != nil {
		u.Logger.Error("u.OrdServiceDeveloper.Exec.create.receive", err.Error(), err)
		return nil, err
	}

	// parse json to get address:
	// ex: {"mnemonic": "chaos dawn between remember raw credit pluck acquire satoshi rain one valley","passphrase": ""}

	jsonStr := strings.ReplaceAll(resp.Stdout, "\n", "")
	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")

	var receiveResp ord_service.ReceiveCmdStdoputRespose

	err = json.Unmarshal([]byte(jsonStr), &receiveResp)
	if err != nil {
		u.Logger.Error("CreateDeveloperInscribeBTC.Unmarshal", err.Error(), err)
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

	u.Logger.Info("CreateDeveloperInscribeBTC.calculateMintPrice", resp)
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
	walletAddress.OrdAddress = receiveResp.Address
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = input.File
	walletAddress.InscriptionID = ""
	walletAddress.FeeRate = input.FeeRate
	walletAddress.ExpiredAt = time.Now().Add(time.Hour * time.Duration(expiredTime))
	walletAddress.FileName = input.FileName
	walletAddress.UserUuid = input.UserUuid
	if input.NeedVerifyAuthentic() {
		pags, err := u.ListDeveloperInscribeBTC(&entity.FilterDeveloperInscribeBT{
			BaseFilters: entity.BaseFilters{
				Page:  1,
				Limit: 1,
			},
			TokenAddress: &input.TokenAddress,
			TokenId:      &input.TokenId,
			NeStatuses:   []entity.StatusDeveloperInscribe{entity.StatusDeveloperInscribe_TxMintFailed},
		})
		if err != nil {
			return nil, err
		}
		inscribers := pags.Result.([]entity.DeveloperInscribeBTCResp)
		if len(inscribers) > 0 {
			return nil, errors.New("Inscribe was minted")
		}
	}

	err = u.Repo.InsertDeveloperInscribeBTC(walletAddress)
	if err != nil {
		u.Logger.Error("u.CreateDeveloperInscribeBTC.InsertDeveloperInscribeBTC", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) ListDeveloperInscribeBTC(req *entity.FilterDeveloperInscribeBT) (*entity.Pagination, error) {
	return u.Repo.ListDeveloperInscribeBTC(req)
}

func (u Usecase) DetailDeveloperInscribeBTC(uuid string) (*entity.DeveloperInscribeBTCResp, error) {
	return u.Repo.FindDeveloperInscribeBTCByNftID(uuid)
}

func (u Usecase) RetryDeveloperInscribeBTC(id string) error {
	item, _ := u.Repo.FindDeveloperInscribeBTC(id)
	u.Logger.Info("item: ", item, id)
	if item != nil {
		if item.Status == entity.StatusDeveloperInscribe_NotEnoughBalance {
			item.Status = entity.StatusDeveloperInscribe_Pending
			_, err := u.Repo.UpdateDeveloperInscribe(item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// JOBs:
// step 1: job check balance for list inscribe
func (u Usecase) JobDeveloperInscribe_WaitingBalance() error {

	_, bs, err := u.buildBTCClient()

	if err != nil {
		go u.trackDeveloperInscribeHistory("", "JobDeveloperInscribe_WaitingBalance", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error())
		return err
	}
	listPending, _ := u.Repo.DeveloperListBTCInscribePending()
	if len(listPending) == 0 {
		// go u.trackDeveloperInscribeHistory("", "JobDeveloperInscribe_WaitingBalance", "", "", "ListBTCInscribePending", "[]")
		return nil
	}

	for _, item := range listPending {

		// check balance:
		balance, confirm, err := bs.GetBalance(item.SegwitAddress)

		fmt.Println("GetBalance response: ", balance, confirm, err)

		if err != nil {
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_WaitingBalance", item.TableName(), item.Status, "GetBalance - with err", err.Error())
			continue
		}
		if balance == nil {
			err = errors.New("balance is nil")
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_WaitingBalance", item.TableName(), item.Status, "GetBalance", err.Error())
			continue
		}

		if balance.Uint64() == 0 {
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_WaitingBalance", item.TableName(), item.Status, "GetBalance", "0")
			continue
		}

		// get required amount to check vs temp wallet balance:
		amount, ok := big.NewInt(0).SetString(item.Amount, 10)
		if !ok {
			err := errors.New("cannot parse amount")
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_WaitingBalance", item.TableName(), item.Status, "SetString(amount) err", err.Error())
			continue
		}

		if amount.Uint64() == 0 {
			err := errors.New("balance is zero")
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_WaitingBalance", item.TableName(), item.Status, "amount.Uint64() err", err.Error())
			continue
		}

		if balance.Uint64() < amount.Uint64() {
			err := fmt.Errorf("Not enough amount %d < %d ", balance.Uint64(), amount.Uint64())
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_WaitingBalance", item.TableName(), item.Status, "compare balance err", err.Error())

			item.Status = entity.StatusDeveloperInscribe_NotEnoughBalance
			u.Repo.UpdateDeveloperInscribe(&item)
			continue
		}

		// received fund:
		item.Status = entity.StatusDeveloperInscribe_ReceivedFund
		item.IsConfirm = true

		_, err = u.Repo.UpdateDeveloperInscribe(&item)
		if err != nil {
			fmt.Printf("Could not UpdateDeveloperInscribe id %s - with err: %v", item.ID, err)
			continue
		}

		go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_WaitingBalance", item.TableName(), item.Status, "Updated StatusDeveloperInscribe_ReceivedFund", "ok")
		u.Logger.Info(fmt.Sprintf("JobDeveloperInscribe_WaitingBalance.CheckReceiveBTC.%s", item.SegwitAddress), item)
		u.Notify("JobDeveloperInscribe_WaitingBalance", item.SegwitAddress, fmt.Sprintf("%s received BTC %d from [InscriptionID] %s", item.SegwitAddress, item.Status, item.InscriptionID))

	}

	return nil
}

// step 2: job send all fund from segwit address to ord wallet:
func (u Usecase) JobDeveloperInscribe_SendBTCToOrdWallet() error {

	_, bs, err := u.buildBTCClient()

	if err != nil {
		go u.trackDeveloperInscribeHistory("", "JobDeveloperInscribe_SendBTCToOrdWallet", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error())
		return err
	}

	listTosendBtc, _ := u.Repo.DeveloperListBTCInscribeByStatus([]entity.StatusDeveloperInscribe{entity.StatusDeveloperInscribe_ReceivedFund})
	if len(listTosendBtc) == 0 {
		// go u.trackDeveloperInscribeHistory("", "JobDeveloperInscribe_SendBTCToOrdWallet", "", "", "DeveloperListBTCInscribeByStatus", "[]")
		return nil
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusDeveloperInscribe_ReceivedFund {

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
				go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_SendBTCToOrdWallet", item.TableName(), item.Status, "SendTransactionWithPreferenceFromSegwitAddress err", err.Error())
				continue
			}

			item.TxSendBTC = txID
			item.Status = entity.StatusDeveloperInscribe_SendingBTCFromSegwitAddrToOrdAddr
			// item.ErrCount = 0 // reset error count!
			// TODO: update item
			_, err = u.Repo.UpdateDeveloperInscribe(&item)
			if err != nil {
				fmt.Printf("Could not UpdateDeveloperInscribe id %s - with err: %v", item.ID, err)
			}

		}
	}

	return nil
}

// job check 3 tx send: tx user send to temp wallet, tx mint, tx send nft to user
func (u Usecase) JobDeveloperInscribe_CheckTxSend() error {
	_, bs, err := u.buildBTCClient()
	if err != nil {
		logger.AtLog.Logger.Error("Could not initialize Bitcoin RPCClient failed", zap.Error(err))
		return err
	}

	// get list sending tx:
	listTosendBtc, _ := u.Repo.DeveloperListBTCInscribeByStatus([]entity.StatusDeveloperInscribe{entity.StatusDeveloperInscribe_Minting, entity.StatusDeveloperInscribe_SendingBTCFromSegwitAddrToOrdAddr})
	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {
		fields := []zapcore.Field{
			zap.String("id", item.ID.Hex()),
			zap.String("file_name", item.FileName),
		}

		statusSuccess := entity.StatusDeveloperInscribe_Minted
		txHashDb := item.TxMintNft

		if item.Status == entity.StatusDeveloperInscribe_SendingBTCFromSegwitAddrToOrdAddr {
			statusSuccess = entity.StatusDeveloperInscribe_SentBTCFromSegwitAddrToOrdAdd
			txHashDb = item.TxSendBTC
		}
		if item.Status == entity.StatusDeveloperInscribe_Minting {
			item.IsMinted = true
		}

		logger.AtLog.Logger.With(fields...).Error("Could not GetTransaction Bitcoin RPCClient")
		go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_CheckTxSend", item.TableName(), item.Status, "btcClient.GetTransaction: "+txHashDb, err.Error())

		go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_CheckTxSend", item.TableName(), item.Status, "bs.CheckTx: "+txHashDb, "Begin check tx via api.")

		// check with api:
		txInfo, err := bs.CheckTx(txHashDb)
		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.AtLog.Logger.With(fields...).Error("Could not CheckTx")
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_CheckTxSend", item.TableName(), item.Status, "bs.CheckTx: "+txHashDb, err.Error())
		}

		if txInfo.Confirmations >= 1 {
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_CheckTxSend", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txHashDb, txInfo.Confirmations)
			// send nft ok now:
			item.Status = statusSuccess
			item.IsSuccess = statusSuccess == entity.StatusDeveloperInscribe_Minted
			_, err = u.Repo.UpdateDeveloperInscribe(&item)
			if err != nil {
				fields = append(fields, zap.Error(err))
				logger.AtLog.Logger.With(fields...).Error("Could not UpdateDeveloperInscribe")
			}
		}
	}

	return nil
}

// job 4: mint nft:
func (u Usecase) JobDeveloperInscribe_MintNft() error {
	listTosendBtc, _ := u.Repo.DeveloperListBTCInscribeByStatus([]entity.StatusDeveloperInscribe{entity.StatusDeveloperInscribe_SentBTCFromSegwitAddrToOrdAdd})
	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {
		fields := []zapcore.Field{
			zap.String("id", item.ID.Hex()),
			zap.String("file_name", item.FileName),
		}

		logger.AtLog.Logger.With(fields...).Info("Mint nft now...")

		// - Upload the Animation URL to GCS
		typeFile := ""

		if len(item.FileName) == 0 {
			err := errors.New("File name invalid")
			u.Logger.Error("JobDeveloperInscribe_MintNft.len(Filename)", err.Error(), err)
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_MintNft", item.TableName(), item.Status, "CheckFileName", err.Error())
			continue
		}

		typeFiles := strings.Split(item.FileName, ".")
		if len(typeFiles) < 2 {
			err := errors.New("File name invalid")
			u.Logger.Error("JobDeveloperInscribe_MintNft.len(Filename)", err.Error(), err)
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_MintNft", item.TableName(), item.Status, "CheckFileName", err.Error())
			continue
		}

		typeFile = typeFiles[len(typeFiles)-1]
		fields = append(fields, zap.String("type_file", typeFile))
		logger.AtLog.Logger.Info("TypeFile", fields...)

		// update google clound: TODO need to move into api to avoid create file many time.
		_, base64Str, err := decodeFileBase64(item.FileURI)
		if err != nil {
			u.Logger.Error("JobDeveloperInscribe_MintNft.decodeFileBase64", err.Error(), err)
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_MintNft", item.TableName(), item.Status, "helpers.decodeFileBase64", err.Error())
			continue
		}

		now := time.Now().UTC().Unix()
		uploaded, err := u.GCS.UploadBaseToBucket(base64Str, fmt.Sprintf("btc-projects/%s/%d.%s", item.OrdAddress, now, typeFile))
		if err != nil {
			u.Logger.Error("JobDeveloperInscribe_MintNft.helpers.UploadBaseToBucket.Base64DecodeRaw", err.Error(), err)
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_MintNft", item.TableName(), item.Status, "helpers.BUploadBaseToBucket.ase64DecodeRaw", err.Error())
			continue
		}
		item.LocalLink = uploaded.FullPath

		fileURI := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
		item.FileURI = fileURI

		mintData := ord_service.MintRequest{
			WalletName:        item.UserAddress,
			FileUrl:           fileURI,
			FeeRate:           int(item.FeeRate),
			DryRun:            false,
			AutoFeeRateSelect: false,
			RequestId:         item.UUID,
			// new key for ord v5.1, support mint + send in 1 tx:
			DestinationAddress: item.OriginUserAddress,
		}

		resp, err := u.OrdServiceDeveloper.Mint(mintData)

		if err != nil {
			u.Logger.Error("OrdService.Mint", err.Error(), err)
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_MintNft", item.TableName(), item.Status, mintData, err.Error())
			continue
		}
		// if not err => update status ok now:
		//TODO: handle log err: Database already open. Cannot acquire lock

		item.Status = entity.StatusDeveloperInscribe_Minting
		// item.ErrCount = 0 // reset error count!

		item.OutputMintNFT = resp

		_, err = u.Repo.UpdateDeveloperInscribe(&item)
		if err != nil {
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_MintNft", item.TableName(), item.Status, "JobDeveloperInscribe_MintNft.UpdateDeveloperInscribe", err.Error())
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
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_MintNft", item.TableName(), item.Status, "JobDeveloperInscribe_MintNft.Unmarshal(btcMintResp)", err.Error())
			continue
		}

		item.TxMintNft = btcMintResp.Reveal
		item.InscriptionID = btcMintResp.Inscription
		_, err = u.Repo.UpdateDeveloperInscribe(&item)
		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.AtLog.Logger.With(fields...).Error("Could not UpdateDeveloperInscribe")
			go u.trackDeveloperInscribeHistory(item.ID.String(), "JobDeveloperInscribe_MintNft", item.TableName(), item.Status, "JobDeveloperInscribe_MintNft.UpdateDeveloperInscribe", err.Error())
		}
	}

	return nil
}

func (u Usecase) Developer_SendTokenByWallet(receiveAddr, inscriptionID, walletAddressName string, rate int) (*ord_service.ExecRespose, error) {

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

	resp, err := u.OrdServiceDeveloper.Exec(sendTokenReq)
	return resp, err
}

func (u Usecase) DeveloperDeveloperGetNftsOwnerOf(walletName string) (*ord_service.ExecRespose, error) {

	listNFTsReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			walletName,
			"wallet",
			"inscriptions",
		}}

	u.Logger.Info("listNFTsReq", listNFTsReq)
	resp, err := u.OrdServiceDeveloper.Exec(listNFTsReq)
	defer u.Notify("DeveloperGetNftsOwnerOf", "ord_marketplace_master", "inscriptions")
	if err != nil {
		u.Logger.Info("u.OrdServiceDeveloper.Exec.Error", err.Error())
		u.Logger.Error("u.OrdServiceDeveloper.Exec", err.Error(), err)
		return nil, err
	}
	u.Logger.Info("listNFTsRep", resp)
	return resp, err
}

func (u *Usecase) trackDeveloperInscribeHistory(id, name, table string, status interface{}, requestMsg interface{}, responseMsg interface{}) {
	trackData := &entity.DeveloperInscribeBTCLogs{
		RecordID:    id,
		Name:        name,
		Table:       table,
		Status:      status,
		RequestMsg:  requestMsg,
		ResponseMsg: responseMsg,
	}
	err := u.Repo.DeveloperCreateDeveloperInscribeBTCLog(trackData)
	if err != nil {
		fmt.Printf("trackDeveloperInscribeHistory.%s.Error:%s", name, err.Error())
	}

}
