package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/helpers"
)

// for api create a new mint:
func (u Usecase) CreateMintReceiveAddress(input structure.MintNftBtcData) (*entity.MintNftBtc, error) {

	if len(input.ProjectID) == 0 || len(input.WalletAddress) == 0 || len(input.RefundUserAddress) == 0 {
		return nil, errors.New("data invalid")
	}

	walletAddress := &entity.MintNftBtc{}

	receiveAddress := ""
	privateKey := ""
	var err error

	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
	if err != nil {
		u.Logger.Error("u.CreateMintReceiveAddress.FindProjectByTokenID", err.Error(), err)
		return nil, errors.New("project not found")
	}

	// find Project and make sure index < max supply
	if p.MintingInfo.Index >= p.MaxSupply {
		err = fmt.Errorf("project %s is minted out", input.ProjectID)
		u.Logger.Error("projectIsMintedOut", err.Error(), err)
		return nil, err
	}

	// verify paytype:
	if input.PayType != utils.NETWORK_BTC && input.PayType != utils.NETWORK_ETH {
		err = errors.New("only support payType is eth or btc")
		u.Logger.Error("u.CreateMintReceiveAddress.Check(payType)", err.Error(), err)
		return nil, err
	}

	// cal min price:
	mintPriceStr := "0"
	mintPriceInt, err := strconv.Atoi(p.MintPrice)
	if err != nil {
		u.Logger.Error("u.CreateMintReceiveAddress.FindProjectByTokenID", err.Error(), err)
		return nil, err
	}
	networkFee, err := strconv.Atoi(p.NetworkFee)
	if err == nil {
		mintPriceInt += networkFee
	}

	var btcRate, ethRate float64

	// check type:
	if input.PayType == utils.NETWORK_BTC {
		privateKey, _, receiveAddress, err = btc.GenerateAddressSegwit()
		if err != nil {
			u.Logger.Error("u.CreateMintReceiveAddress.GenerateAddressSegwit", err.Error(), err)
			return nil, err
		}
		mintPriceStr = strconv.Itoa(mintPriceInt)

		_, btcRate, ethRate, err = u.convertBTCToETH("1")

	} else if input.PayType == utils.NETWORK_ETH {
		ethClient := eth.NewClient(nil)

		privateKey, _, receiveAddress, err = ethClient.GenerateAddress()
		if err != nil {
			u.Logger.Error("CreateMintReceiveAddress.ethClient.GenerateAddress", err.Error(), err)
			return nil, err
		}
		mintPriceStr, btcRate, ethRate, err = u.convertBTCToETH(fmt.Sprintf("%f", float64(mintPriceInt)/1e8))
		if err != nil {
			u.Logger.Error("convertBTCToETH", err.Error(), err)
			return nil, err
		}
		fmt.Println("mintPriceStr: ", mintPriceStr)
	}

	if len(receiveAddress) == 0 || len(privateKey) == 0 {
		err = errors.New("can not create the wallet")
		u.Logger.Error("u.CreateMintReceiveAddress.GenerateAddress", err.Error(), err)
		return nil, err
	}

	// set temp wallet info:
	walletAddress.PayType = input.PayType

	if len(os.Getenv("SECRET_KEY")) == 0 {
		err = errors.New("please config SECRET_KEY")
		u.Logger.Error("u.CreateMintReceiveAddress.GenerateAddress", err.Error(), err)
		return nil, err
	}

	privateKeyEnCrypt, err := encrypt.EncryptToString(privateKey, os.Getenv("SECRET_KEY"))
	if err != nil {
		u.Logger.Error("u.CreateMintReceiveAddress.Encrypt", err.Error(), err)
		return nil, err
	}

	walletAddress.PrivateKey = privateKeyEnCrypt
	walletAddress.ReceiveAddress = receiveAddress
	walletAddress.RefundUserAdress = input.RefundUserAddress

	// for analytics:
	walletAddress.ProjectNetworkFee = networkFee
	walletAddress.ProjectMintPrice = mintPriceInt
	walletAddress.BtcRate = btcRate
	walletAddress.EthRate = ethRate

	u.Logger.Info("CreateMintReceiveAddress.receive", receiveAddress)

	expiredTime := utils.INSCRIBE_TIMEOUT
	if u.Config.ENV == "develop" {
		expiredTime = 1
	}

	walletAddress.Amount = mintPriceStr
	walletAddress.OriginUserAddress = input.WalletAddress
	walletAddress.Status = entity.StatusMint_Pending
	walletAddress.ProjectID = input.ProjectID
	walletAddress.Balance = "0"
	walletAddress.ExpiredAt = time.Now().Add(time.Hour * time.Duration(expiredTime))

	// insert now:
	err = u.Repo.InsertMintNftBtc(walletAddress)
	if err != nil {
		u.Logger.Error("u.CreateMintReceiveAddress.InsertMintNftBtc", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) CancelMintNftBtc(wallet, uuid string) error {
	mintItem, _ := u.Repo.FindMintNftBtcByNftID(uuid)
	if mintItem == nil {
		return errors.New("item not found")
	}
	if !strings.EqualFold(wallet, mintItem.OriginUserAddress) {
		return errors.New("perminsion denied")
	}
	if mintItem.Status != entity.StatusMint_Pending {
		return errors.New("Can not cancel this, the item is in progress.")
	}
	return u.Repo.UpdateCancelMintNftBtc(mintItem.UUID)
}

func (u Usecase) GetDetalMintNftBtc(uuid string) (*structure.MintingInscription, error) {
	mintItem, _ := u.Repo.FindMintNftBtcByNftID(uuid)
	if mintItem == nil {
		return nil, errors.New("item not found")
	}

	projectInfo, _ := u.Repo.FindProjectByTokenID(mintItem.OriginUserAddress)
	if mintItem == nil {
		return nil, errors.New("item not found")
	}

	isCancel := int(mintItem.Status) == 0

	minting := &structure.MintingInscription{
		ID:            mintItem.UUID,
		CreatedAt:     mintItem.CreatedAt,
		Status:        entity.StatusMintToText[mintItem.Status],
		StatusIndex:   int(mintItem.Status),
		FileURI:       mintItem.FileURI,
		ProjectID:     mintItem.ProjectID,
		ProjectImage:  projectInfo.Thumbnail,
		ProjectName:   projectInfo.Name,
		InscriptionID: mintItem.InscriptionID,

		ReceiveAddress: mintItem.ReceiveAddress,
		IsCancel:       isCancel,
		TxMint:         mintItem.TxMintNft,
		TxSendNft:      mintItem.TxSendNft,
	}
	return minting, nil
}

// JOBs mint begin:
// step 1: job check balance for list mint_nft_btc
func (u Usecase) JobMint_CheckBalance() error {

	_, bs, err := u.buildBTCClient()

	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMint_CheckBalance", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error(), true)
		return err
	}

	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMint_CheckBalance", "", "", "Could not initialize Ether RPCClient - with err", err.Error(), true)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	// get list mint pending to check balance:
	listPending, _ := u.Repo.ListMintNftBtcPending()
	if len(listPending) == 0 {
		// go u.trackMintNftBtcHistory("", "JobMint_CheckBalance", "", "", "ListMintNftBtcPending", "[]", false)
		return nil
	}

	for _, item := range listPending {

		time.Sleep(1 * time.Second)

		// check balance:
		balance := big.NewInt(0)
		confirm := -1

		if item.PayType == utils.NETWORK_BTC {

			balance, confirm, err = bs.GetBalance(item.ReceiveAddress)
			fmt.Println("GetBalance btc response: ", balance, confirm, err)

		} else if item.PayType == utils.NETWORK_ETH {
			// check eth balance:

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			balance, err = ethClient.GetBalance(ctx, item.ReceiveAddress)
			fmt.Println("GetBalance eth response: ", balance, err)

			confirm = 1
		}

		if err != nil {
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "GetBalance - with err", err.Error(), true)
			continue
		}
		if balance == nil {
			err = errors.New("balance is nil")
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "GetBalance", err.Error(), false)
			continue
		}

		if balance.Uint64() == 0 {
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "GetBalance", "0", false)
			continue
		}

		// get required amount to check vs temp wallet balance:
		amount, ok := big.NewInt(0).SetString(item.Amount, 10)
		if !ok {
			err := errors.New("cannot parse amount")
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "SetString(amount) err", err.Error(), true)
			continue
		}

		if amount.Uint64() == 0 {
			err := errors.New("balance is zero")
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "amount.Uint64() err", err.Error(), true)
			continue
		}

		// set receive balance:
		item.Balance = amount.String()

		if balance.Uint64() < amount.Uint64() {
			err := fmt.Errorf("Not enough amount %d < %d ", balance.Uint64(), amount.Uint64())
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "compare balance err", err.Error(), true)

			item.Status = entity.StatusMint_NeedToRefund
			item.ReasonRefund = "Not enough balance"
			u.Repo.UpdateMintNftBtc(&item)
			continue
		}

		if confirm == 0 {
			item.Status = entity.StatusMint_WaitingForConfirms
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "Updated StatusMint_WaitingForConfirms", "0", true)
		}
		if confirm >= 1 {
			// received fund:
			item.Status = entity.StatusMint_ReceivedFund
			item.IsConfirm = true

			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "Updated StatusMint_ReceivedFund", "ok", true)
			u.Logger.Info(fmt.Sprintf("JobMint_CheckBalance.CheckReceiveFund.%s", item.ReceiveAddress), item)
			go u.Notify("JobMint_CheckBalance", item.ReceiveAddress, fmt.Sprintf("%s received %s %d from [UUID] %s", item.ReceiveAddress, item.PayType, item.Status, item.UUID))
		}

		_, err = u.Repo.UpdateMintNftBtc(&item)
		if err != nil {
			fmt.Printf("Could not UpdateMintNftBtc uuid %s - with err: %v", item.UUID, err)
			continue
		}

	}

	return nil
}

// job 2: mint nft now:
func (u Usecase) JobMint_MintNftBtc() error {

	listToMint, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint(entity.StatusMint_ReceivedFund)})
	if len(listToMint) == 0 {
		// go u.trackMintNftBtcHistory("", "JobMint_MintNftBtc", "", "", "ListMintNftBtcByStatus", "[]")
		return nil
	}

	for _, item := range listToMint {

		// get data from project
		p, err := u.Repo.FindProjectByTokenID(item.ProjectID)
		if err != nil {
			u.Logger.Error("JobMint_MintNftBtc.FindProjectByTokenID", err.Error(), err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "FindProjectByTokenID", err.Error(), true)
			continue
		}

		// check Project and make sure index < max supply
		if p.MintingInfo.Index >= p.MaxSupply {

			// update need to return:
			item.ReasonRefund = "project is minted out"
			item.Status = entity.StatusMint_NeedToRefund

			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "Update need to refund for minted out", err.Error(), true)
			}
			err = fmt.Errorf("project %s is minted out", item.ProjectID)
			u.Logger.Error("projectIsMintedOut", err.Error(), err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "Updated to minted out", err.Error(), true)
			continue
		}

		// - Get project.AnimationURL
		projectNftTokenUri := &structure.ProjectAnimationUrl{}
		err = helpers.Base64DecodeRaw(p.NftTokenUri, projectNftTokenUri)
		if err != nil {
			u.Logger.Error("JobMint_MintNftBtc.Base64DecodeRaw", err.Error(), err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "Base64DecodeRaw", err.Error(), true)
			continue
		}

		// - Upload the Animation URL to GCS
		animation := projectNftTokenUri.AnimationUrl
		u.Logger.Info("animation", animation)

		// for html type:
		if animation != "" {
			animation = strings.ReplaceAll(animation, "data:text/html;base64,", "")
			now := time.Now().UTC().Unix()
			uploaded, err := u.GCS.UploadBaseToBucket(animation, fmt.Sprintf("btc-projects/%s/%d.html", p.TokenID, now))
			if err != nil {
				u.Logger.Error("JobMint_MintNftBtc.UploadBaseToBucket", err.Error(), err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "UploadBaseToBucket", err.Error(), true)
				continue
			}
			item.FileURI = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)

		} else {
			// for image type:
			images := p.Images
			u.Logger.Info("images", len(images))
			if len(images) > 0 {
				item.FileURI = images[0]
				newImages := []string{}
				processingImages := p.ProcessingImages

				//remove the project's image out of the current projects
				for i := 1; i < len(images); i++ {
					newImages = append(newImages, images[i])
				}
				processingImages = append(p.ProcessingImages, item.FileURI)
				p.Images = newImages
				p.ProcessingImages = processingImages
			}
		}
		//end Animation URL
		if len(item.FileURI) == 0 {
			err = errors.New("There is no file uri to mint")
			u.Logger.Error("JobMint_MintNftBtc.UploadBaseToBucket", err.Error(), err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "UploadBaseToBucket", err.Error(), true)
			continue
		}

		baseUrl, err := url.Parse(item.FileURI)
		if err != nil {
			u.Logger.Error("JobMint_MintNftBtc.UploadBaseToBucket", err.Error(), err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "Parse(FileURI)", err.Error(), true)
			continue
		}
		_ = baseUrl
		// start call rpc mint nft now:
		mintData := ord_service.MintRequest{
			WalletName: os.Getenv("ORD_MASTER_ADDRESS"),
			FileUrl:    baseUrl.String(),
			FeeRate:    entity.DEFAULT_FEE_RATE, //auto
			DryRun:     false,
			RequestId:  item.UUID,      // to track log
			ProjectID:  item.ProjectID, // to track log
		}

		u.Logger.Info("mintData", mintData)
		// execute mint:
		resp, err := u.OrdService.Mint(mintData)
		if err != nil {
			u.Logger.Error("JobMint_MintNftBtc.OrdService", err.Error(), err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc.Mint", item.TableName(), item.Status, mintData, err.Error(), true)
			continue
		}
		u.Logger.Info("mint.resp", resp)

		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "MintOK", resp, false)

		//update
		// if not err => update status ok now:
		//TODO: handle log err: Database already open. Cannot acquire lock

		item.Status = entity.StatusMint_Minting
		// item.ErrCount = 0 // reset error count!

		item.OutputMintNFT = resp

		_, err = u.Repo.UpdateMintNftBtc(&item)
		if err != nil {
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "JobMint_MintNftBtc.UpdateMintNftBtc", err.Error(), true)
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
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "JobMint_MintNftBtc.Unmarshal(btcMintResp)", err.Error(), true)
			continue
		}

		item.TxMintNft = btcMintResp.Reveal
		item.InscriptionID = btcMintResp.Inscription
		// TODO: update item
		_, err = u.Repo.UpdateMintNftBtc(&item)
		if err != nil {
			fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err.Error())
		}

		// update project:
		updated, err := u.Repo.UpdateProject(p.UUID, p)
		if err != nil {
			u.Logger.Error("JobMint_MintNftBtc.UpdateProject", err.Error(), err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "JobMint_MintNftBtc.UpdateProject", err.Error(), true)
		}
		u.Logger.Info("project.Updated", updated)

		fmt.Println("update project, token info when minting ...")

		// create entity.TokenURI
		_, err = u.CreateBTCTokenURI(item.ProjectID, item.InscriptionID, item.FileURI, entity.TokenPaidType(item.PayType))
		if err != nil {
			fmt.Printf("Could CreateBTCTokenURI - with err: %v", err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "u.CreateBTCTokenURI()", err.Error(), true)
			continue
		}
		_, err = u.Repo.UpdateMintNftBtcByFilter(item.UUID, bson.M{"$set": bson.M{"isUpdatedNftInfo": true}})
		if err != nil {
			fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "UpdateMintNftBtc", err.Error(), true)
		}

		go u.Notify(fmt.Sprintf("[MintFor][%s][projectID %s]", item.PayType, item.ProjectID), item.ReceiveAddress, fmt.Sprintf("Made mining transaction for %s, waiting network confirm %s", item.UserAddress, resp.Stdout))

		// try to update inscription_index
		go u.getInscribeInfoForMintSuccessToUpdate(item.InscriptionID)

	}

	return nil
}

// job check 3 tx mint/send nft
func (u Usecase) JobMint_CheckTxMintSend() error {

	btcClient, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		go u.trackMintNftBtcHistory("", "JobMint_CheckTxMintSend", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error(), true)
		return err
	}

	// get list pending tx:
	listTxToCheck, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint_Minting, entity.StatusMint_SendingNFTToUser})
	if len(listTxToCheck) == 0 {
		return nil
	}

	for _, item := range listTxToCheck {

		var txToCheck string
		var confirm int64 = -1

		if item.Status == entity.StatusMint_Minting {
			txToCheck = item.TxMintNft
		} else if item.Status == entity.StatusMint_SendingNFTToUser {
			txToCheck = item.TxSendNft
		}

		txHash, err := chainhash.NewHashFromStr(txToCheck)
		if err != nil {
			fmt.Printf("Could not NewHashFromStr Bitcoin RPCClient - with err: %v", err)
			continue
		}

		txResponse, err := btcClient.GetTransaction(txHash)

		if err == nil {
			confirm = txResponse.Confirmations

			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "btcClient.txResponse.Confirmations: "+txToCheck, confirm, false)

			if confirm <= 0 {
				continue
			}

		} else {
			fmt.Printf("Could not GetTransaction Bitcoin RPCClient - with err: %v", err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "btcClient.GetTransaction: "+txToCheck, err.Error(), false)

			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "bs.CheckTx: "+txToCheck, "Begin check tx via api.", false)

			// check with api:
			txInfo, err := bs.CheckTx(txToCheck)
			if err != nil {
				fmt.Printf("Could not bs - with err: %v", err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "bs.CheckTx: "+txToCheck, err.Error(), true)
				continue
			}

			confirm = int64(txInfo.Confirmations)

			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txToCheck, txInfo.Confirmations, false)

			if confirm <= 0 {
				continue
			}
		}

		// just check 1 confirm:
		if confirm >= 1 {
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txToCheck, confirm, true)
			// tx ok now:

			if item.Status == entity.StatusMint_Minting {
				item.Status = entity.StatusMint_Minted
				item.IsMinted = true
			} else if item.Status == entity.StatusMint_SendingNFTToUser {
				item.Status = entity.StatusMint_SentNFTToUser
				item.IsSentUser = true
			}

			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err)
				continue
			}
			if item.Status == entity.StatusMint_Minted {
				err = u.Repo.UpdateTokenOnchainStatusByTokenId(item.InscriptionID)
				if err != nil {
					u.Logger.Error(fmt.Sprintf("JobMint_CheckTxMintSend.%s.UpdateTokenOnchainStatusByTokenId.Error", item.InscriptionID), err.Error(), err)
					go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "UpdateTokenOnchainStatusByTokenId()", err.Error(), true)
				}
				// update inscription_index for token uri
				go u.getInscribeInfoForMintSuccessToUpdate(item.InscriptionID)
			}

		}
	}

	return nil
}

// job 4: send nft:
func (u Usecase) JobMint_SendNftToUser() error {

	// get list buy order status = StatusInscribe_Minted:
	listTosendNft, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint_Minted})
	if len(listTosendNft) == 0 {
		return nil
	}

	for _, item := range listTosendNft {

		// check nft in master wallet or not:
		if len(item.InscriptionID) == 0 {
			go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "checkEmpty(nftID)", "Nft id empty", true)
			continue
		}
		listNFTsRep, err := u.GetNftsOwnerOf(os.Getenv("ORD_MASTER_ADDRESS"))
		if err != nil {
			go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "GetNftsOwnerOf.Error", err.Error(), true)
			continue
		}

		go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "GetNftsOwnerOf.listNFTsRep", listNFTsRep, false)

		// parse nft data:
		var resp []struct {
			Inscription string `json:"inscription"`
			Location    string `json:"location"`
			Explorer    string `json:"explorer"`
		}

		err = json.Unmarshal([]byte(listNFTsRep.Stdout), &resp)
		if err != nil {
			go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "GetNftsOwnerOf.Unmarshal(listNFTsRep)", err.Error(), true)
			continue
		}
		owner := false
		for _, nft := range resp {
			if strings.EqualFold(nft.Inscription, item.InscriptionID) {
				owner = true
				break
			}

		}

		if !owner {
			go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "GetNftsOwnerOf.CheckNFTOwner", owner, true)
			continue
		}

		// transfer now:
		sendTokenReq := ord_service.ExecRequest{
			Args: []string{
				"--wallet",
				os.Getenv("ORD_MASTER_ADDRESS"),
				"wallet",
				"send",
				item.OriginUserAddress,
				item.InscriptionID,
				"--fee-rate",
				fmt.Sprintf("%d", entity.DEFAULT_FEE_RATE),
			}}

		u.Logger.Info("sendTokenReq", sendTokenReq)
		mintResp, err := u.OrdService.Exec(sendTokenReq)

		go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "SendTokenByWallet.ExecRequest.SendNft()", mintResp, true)

		if err != nil {
			u.Logger.Error(fmt.Sprintf("JobMin_SendNftToUser.SendTokenMKP.%s.Error", item.OriginUserAddress), err.Error(), err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "SendTokenByWallet.err", err.Error(), true)
			continue
		}

		//TODO: handle log err: Database already open. Cannot acquire lock

		// Update status first if none err:
		item.Status = entity.StatusMint_SendingNFTToUser
		// item.ErrCount = 0 // reset error count!

		item.OutputSendNFT = mintResp

		_, err = u.Repo.UpdateMintNftBtc(&item)
		if err != nil {
			err := fmt.Errorf("Could not UpdateMintNftBtc id %s - with err: %v", item.UUID, err.Error())
			u.Logger.Error("JobMin_SendNftToUser.UpdateMintNftBtc", err.Error(), err)
			go u.trackMintNftBtcHistory(item.UUID, "UpdateMintNftBtc", item.TableName(), item.Status, "SendTokenMKP.UpdateMintNftBtc", err.Error(), true)
			continue
		}

		txResp := mintResp.Stdout
		//txResp := `fd31946b855cbaaf91df4b2c432f9b173e053e65a9879ac909bad028e21b950e\n`
		txResp = strings.ReplaceAll(txResp, "\n", "")

		// update tx:
		item.TxSendNft = txResp
		// item.ErrCount = 0 // reset error count!
		_, err = u.Repo.UpdateMintNftBtc(&item)
		if err != nil {
			errPack := fmt.Errorf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err)
			u.Logger.Error("JobMin_SendNftToUser.UpdateMintNftBtc", errPack.Error(), errPack)
			go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "u.UpdateMintNftBtc.JobMin_SendNftToUser", err.Error(), true)
		}

		u.Logger.Info(fmt.Sprintf("JobMin_SendNftToUser.SendNft.execResp.%s", item.OriginUserAddress), mintResp)

	}
	return nil
}

// job 6:
// refund btc to users:
func (u Usecase) JobMint_RefundBtc() error {

	listToRefund, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint(entity.StatusMint_NeedToRefund)})

	if len(listToRefund) == 0 {
		return nil
	}

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		go u.trackMintNftBtcHistory("", "JobMint_RefundBtc", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error(), true)
		return err
	}

	// eth:
	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMint_RefundBtc", "", "", "Could not initialize Ether RPCClient - with err", err.Error(), true)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	for _, item := range listToRefund {

		if len(item.RefundUserAdress) == 0 {
			continue
		}

		if item.PayType == utils.NETWORK_BTC {

			if len(os.Getenv("SECRET_KEY")) == 0 {
				err = errors.New("please config SECRET_KEY")
				u.Logger.Error("u.JobMint_RefundBtc.GenerateAddress", err.Error(), err)
				continue
			}

			// the user address to refund:
			btcRefundAddress := item.RefundUserAdress

			privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.Decrypt.%s.Error", btcRefundAddress), err.Error(), err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_RefundBtc", item.TableName(), item.Status, "JobMint_RefundBtc.DecryptToString", err.Error(), true)
				continue
			}

			// send user now:
			tx, err := bs.SendTransactionWithPreferenceFromSegwitAddress(privateKeyDeCrypt, item.ReceiveAddress, btcRefundAddress, -1, btc.PreferenceMedium)
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.SendTransactionWithPreferenceFromSegwitAddress.%s.Error", btcRefundAddress), err.Error(), err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_RefundBtc", item.TableName(), item.Status, "JobMint_RefundBtc.SendTransactionWithPreferenceFromSegwitAddress", err.Error(), true)
				continue
			}
			// save tx:
			item.TxRefund = tx
			item.Status = entity.StatusMint_Refunding
			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.UpdateMintNftBtc.%s.Error", btcRefundAddress), err.Error(), err)
				continue
			}
		} else if item.PayType == utils.NETWORK_ETH {

			ethAdressRefund := item.RefundUserAdress

			privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.Decrypt.%s.Error", ethAdressRefund), err.Error(), err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_RefundBtc", item.TableName(), item.Status, "JobMint_RefundBtc.DecryptToString", err.Error(), true)
				continue
			}
			tx, value, err := ethClient.TransferMax(privateKeyDeCrypt, u.Config.MASTER_ADDRESS_CLAIM_ETH)
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.ethClient.TransferMax.%s.Error", ethAdressRefund), err.Error(), err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_RefundBtc", item.TableName(), item.Status, "JobMint_RefundBtc.ethClient.TransferMax", err.Error(), true)
				continue
			}
			// save tx:
			item.TxRefund = tx
			item.AmountRefundUser = value
			item.Status = entity.StatusMint_Refunding
			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.UpdateMintNftBtc.%s.Error", ethAdressRefund), err.Error(), err)
				continue
			}

		}
		time.Sleep(3 * time.Second)
	}

	return nil
}

// job 7:
// send send max fund to master address:
func (u Usecase) JobMint_SendFundToMaster() error {

	listToSentMaster, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint(entity.StatusMint_SentNFTToUser)})

	if len(listToSentMaster) == 0 {
		return nil
	}

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		go u.trackMintNftBtcHistory("", "JobMint_SendFundToMaster", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error(), true)
		return err
	}

	// eth:
	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMint_SendFundToMaster", "", "", "Could not initialize Ether RPCClient - with err", err.Error(), true)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	for _, item := range listToSentMaster {

		if item.PayType == utils.NETWORK_BTC {

			if len(os.Getenv("SECRET_KEY")) == 0 {
				err = errors.New("please config SECRET_KEY")
				u.Logger.Error("u.JobMint_SendFundToMaster.GenerateAddress", err.Error(), err)
				u.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.SECRET_KEY.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_BTC), err.Error(), err)
				continue
			}

			privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.Decrypt.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_BTC), err.Error(), err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_SendFundToMaster", item.TableName(), item.Status, "JobMint_RefundBtc.DecryptToString", err.Error(), true)
				continue
			}

			// send master now:
			tx, err := bs.SendTransactionWithPreferenceFromSegwitAddress(privateKeyDeCrypt, item.ReceiveAddress, u.Config.MASTER_ADDRESS_CLAIM_BTC, -1, btc.PreferenceMedium)
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.SendTransactionWithPreferenceFromSegwitAddress.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_BTC), err.Error(), err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_SendFundToMaster", item.TableName(), item.Status, "JobMint_SendFundToMaster.SendTransactionWithPreferenceFromSegwitAddress", err.Error(), true)
				continue
			}
			// save tx:
			item.TxSendMaster = tx
			item.Status = entity.StatusMint_SendingFundToMaster // TODO: need to a job to check tx.
			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.UpdateBtcWalletAddress.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_BTC), err.Error(), err)
				continue
			}
		} else if item.PayType == utils.NETWORK_ETH {
			privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.Decrypt.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_ETH), err.Error(), err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_SendFundToMaster", item.TableName(), item.Status, "JobMint_SendFundToMaster.DecryptToString", err.Error(), true)
				continue
			}
			tx, amount, err := ethClient.TransferMax(privateKeyDeCrypt, u.Config.MASTER_ADDRESS_CLAIM_ETH)
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.ethClient.TransferMax.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_ETH), err.Error(), err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_SendFundToMaster", item.TableName(), item.Status, "JobMint_SendFundToMaster.ethClient.TransferMax", err.Error(), true)
				continue
			}
			// save tx:
			item.TxSendMaster = tx
			item.AmountSentMaster = amount
			item.Status = entity.StatusMint_SendingFundToMaster
			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.UpdateBtcWalletAddress.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_ETH), err.Error(), err)
				continue
			}

		}
		time.Sleep(3 * time.Second)
	}

	return nil
}

// job 8: check tx send master, refund user:
func (u Usecase) JobMint_CheckTxMasterAndRefund() error {

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		go u.trackMintNftBtcHistory("", "JobMint_CheckTxMasterAndRefund", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error(), true)
		return err
	}

	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMint_CheckBalance", "", "", "Could not initialize Ether RPCClient - with err", err.Error(), true)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	// get list pending tx:
	listTxToCheck, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint_Refunding, entity.StatusMint_SendingFundToMaster})
	if len(listTxToCheck) == 0 {
		return nil
	}

	for _, item := range listTxToCheck {

		var txToCheck string
		var confirm int64 = -1

		if item.Status == entity.StatusMint_Refunding {
			txToCheck = item.TxRefund
		} else if item.Status == entity.StatusMint_SendingFundToMaster {
			txToCheck = item.TxSendMaster
		}

		amountSent := ""

		if item.PayType == utils.NETWORK_ETH {
			fmt.Println("PayType", "eth")
			context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			status, err := ethClient.GetTransaction(context, txToCheck)
			fmt.Println("status, err", status, err)
			if err == nil {
				if status > 0 {
					confirm = 1

				} else {
					continue
				}
			} else {
				// if error maybe tx is pending or rejected
				// TODO check timeout to detect tx is rejected or not.
			}
		} else {
			// check with api btc:
			txInfo, err := bs.CheckTx(txToCheck)
			if err != nil {
				fmt.Printf("Could not bs - with err: %v", err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMasterAndRefund", item.TableName(), item.Status, "bs.CheckTx: "+txToCheck, err.Error(), true)
				continue
			}

			confirm = int64(txInfo.Confirmations)
			amountSent = txInfo.Total.String()

		}

		go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMasterAndRefund", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txToCheck, confirm, false)

		// just check 1 confirm:
		if confirm >= 1 {
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMasterAndRefund", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txToCheck, confirm, true)
			// tx ok now:

			if item.Status == entity.StatusMint_Refunding {
				item.Status = entity.StatusMint_Refunded
				item.IsRefund = true
				if item.PayType == utils.NETWORK_BTC {
					item.AmountRefundUser = amountSent
				}
			} else if item.Status == entity.StatusMint_SendingFundToMaster {
				item.Status = entity.StatusMint_SentFundToMaster
				item.IsSentMaster = true
				if item.PayType == utils.NETWORK_BTC {
					item.AmountSentMaster = amountSent
				}
			}

			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				fmt.Printf("Could not JobMint_CheckTxMasterAndRefund id %s - with err: %v", item.ID, err)
				continue
			}
		}

	}

	return nil
}

func (u *Usecase) trackMintNftBtcHistory(id, name, table string, status interface{}, requestMsg interface{}, responseMsg interface{}, notify bool) {

	trackData := &entity.MintNftBtcLogs{
		RecordID:    id,
		Name:        name,
		Table:       table,
		Status:      status,
		RequestMsg:  requestMsg,
		ResponseMsg: responseMsg,
	}
	err := u.Repo.CreateMintNftBtcLog(trackData)
	if err != nil {
		fmt.Printf("trackMintNftBtcHistory.%s.Error:%s", name, err.Error())
	}

	if notify && requestMsg != nil && responseMsg != nil {
		requestMsgStr := fmt.Sprintf("%v", requestMsg)
		responseMsgStr := fmt.Sprintf("%v", responseMsg)

		preText := fmt.Sprintf("[App: %s][recordID %s] - %s", os.Getenv("JAEGER_SERVICE_NAME"), id, requestMsgStr)

		if _, _, err := u.Slack.SendMessageToSlackWithChannel(os.Getenv("SLACK_MINT_NFT_CHANNEL_ID"), preText, name, responseMsgStr); err != nil {
			fmt.Println("s.Slack.SendMessageToSlack err", err)
		}
	}

}
func (u Usecase) getInscribeInfoForMintSuccessToUpdate(inscriptionID string) error {
	inscribeInfo, err := u.GetInscribeInfo(inscriptionID)
	if err != nil {
		return err
	}
	u.Repo.UpdateTokenInscriptionIndexForMint(inscriptionID, inscribeInfo.Index)

	return nil
}
