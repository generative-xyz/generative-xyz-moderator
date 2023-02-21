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

	// check type:
	if input.PayType == "btc" { // TODO: move to const config
		privateKey, _, receiveAddress, err = btc.GenerateAddressSegwit()
		if err != nil {
			u.Logger.Error("u.CreateMintReceiveAddress.GenerateAddressSegwit", err.Error(), err)
			return nil, err
		}
	} else if input.PayType == "eth" {
		ethClient := eth.NewClient(nil)

		privateKey, _, receiveAddress, err = ethClient.GenerateAddress()
		if err != nil {
			u.Logger.Error("CreateMintReceiveAddress.ethClient.GenerateAddress", err.Error(), err)
			return nil, err
		}
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

	u.Logger.Info("CreateMintReceiveAddress.receive", receiveAddress)

	mintPrice, err := strconv.Atoi(p.MintPrice)
	if err != nil {
		u.Logger.Error("u.CreateMintReceiveAddress.FindProjectByTokenID", err.Error(), err)
		return nil, err
	}
	networkFee, err := strconv.Atoi(p.NetworkFee)
	if err == nil {
		mintPrice += networkFee
	}

	expiredTime := utils.INSCRIBE_TIMEOUT
	if u.Config.ENV == "develop" {
		expiredTime = 1
	}

	walletAddress.Amount = strconv.Itoa(mintPrice)
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

		// received fund:
		item.Status = entity.StatusMint_ReceivedFund
		item.IsConfirm = true

		_, err = u.Repo.UpdateMintNftBtc(&item)
		if err != nil {
			fmt.Printf("Could not UpdateMintNftBtc uuid %s - with err: %v", item.UUID, err)
			continue
		}

		go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "Updated StatusMint_ReceivedFund", "ok", true)
		u.Logger.Info(fmt.Sprintf("JobMint_CheckBalance.CheckReceiveFund.%s", item.ReceiveAddress), item)
		go u.Notify("JobMint_CheckBalance", item.ReceiveAddress, fmt.Sprintf("%s received %s %d from [UUID] %s", item.ReceiveAddress, item.PayType, item.Status, item.UUID))

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

		// start call rpc mint nft now:
		mintData := ord_service.MintRequest{
			WalletName: os.Getenv("ORD_MASTER_ADDRESS"),
			FileUrl:    baseUrl.String(),
			// FeeRate:    entity.DEFAULT_FEE_RATE, //auto
			DryRun:    false,
			RequestId: item.UUID,      // to track log
			ProjectID: item.ProjectID, // to track log
		}

		u.Logger.Info("mintData", mintData)
		// execute mint:
		resp, err := u.OrdService.Mint(mintData)
		if err != nil {
			u.Logger.Error("JobMint_MintNftBtc.OrdService", err.Error(), err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "OrdService", err.Error(), true)
			continue
		}
		u.Logger.Info("mint.resp", resp)
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

		// update project:
		updated, err := u.Repo.UpdateProject(p.UUID, p)
		if err != nil {
			u.Logger.Error("JobMint_MintNftBtc.UpdateProject", err.Error(), err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "JobMint_MintNftBtc.UpdateProject", err.Error(), true)
		}
		u.Logger.Info("project.Updated", updated)

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

		go u.Notify(fmt.Sprintf("[MintFor][%s][projectID %s]", item.PayType, item.ProjectID), item.ReceiveAddress, fmt.Sprintf("Made mining transaction for %s, waiting network confirm %s", item.UserAddress, resp.Stdout))

	}

	return nil
}

// job check 3 tx mint/send nft
func (u Usecase) JobMint_CheckTxMintSend() error {

	btcClient, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	// get list pending tx:
	listTxToCheck, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint_Minting, entity.StatusMint_SendingNFTToUser})
	if len(listTxToCheck) == 0 {
		return nil
	}

	for _, item := range listTxToCheck {

		txHashDb := item.TxMintNft
		item.Status = entity.StatusMint_Minted

		if item.Status == entity.StatusMint_Minting {
			item.IsMinted = true
		} else if item.Status == entity.StatusMint_SendingNFTToUser {
			txHashDb = item.TxSendNft
			item.IsSentUser = true
			item.Status = entity.StatusMint_SentNFTToUser
		}

		txHash, err := chainhash.NewHashFromStr(txHashDb)
		if err != nil {
			fmt.Printf("Could not NewHashFromStr Bitcoin RPCClient - with err: %v", err)
			continue
		}

		txResponse, err := btcClient.GetTransaction(txHash)

		if err == nil {
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "btcClient.txResponse.Confirmations: "+txHashDb, txResponse.Confirmations, false)
			if txResponse.Confirmations >= 1 {
				// send btc ok now:
				_, err = u.Repo.UpdateMintNftBtc(&item)
				if err != nil {
					fmt.Printf("Could not JobMint_CheckTxMintSend id %s - with err: %v", item.ID, err)
					continue
				}
			}
		} else {
			fmt.Printf("Could not GetTransaction Bitcoin RPCClient - with err: %v", err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "btcClient.GetTransaction: "+txHashDb, err.Error(), false)

			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "bs.CheckTx: "+txHashDb, "Begin check tx via api.", false)

			// check with api:
			txInfo, err := bs.CheckTx(txHashDb)
			if err != nil {
				fmt.Printf("Could not bs - with err: %v", err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "bs.CheckTx: "+txHashDb, err.Error(), true)
				continue
			}

			// just check 1 confirm:
			if txInfo.Confirmations >= 1 {
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txHashDb, txInfo.Confirmations, true)
				// tx ok now:
				_, err = u.Repo.UpdateMintNftBtc(&item)
				if err != nil {
					fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err)
				}
				// update project, token info when mint success:
				if item.Status == entity.StatusMint_Minted {
					// create entity.TokenURI
					_, err = u.CreateBTCTokenURI(item.ProjectID, item.InscriptionID, item.FileURI, entity.TokenPaidType(item.PayType))
					if err != nil {
						fmt.Printf("Could CreateBTCTokenURI - with err: %v", err)
						go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "u.CreateBTCTokenURI()", err.Error(), true)
						continue
					}

					err = u.Repo.UpdateTokenOnchainStatusByTokenId(item.InscriptionID)
					if err != nil {
						u.Logger.Error(fmt.Sprintf("JobMint_CheckTxMintSend.%s.UpdateTokenOnchainStatusByTokenId.Error", item.InscriptionID), err.Error(), err)
						go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "UpdateTokenOnchainStatusByTokenId()", err.Error(), true)
						continue
					}
					item.IsUpdatedNftInfo = true
					_, err = u.Repo.UpdateMintNftBtc(&item)
					if err != nil {
						fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err)
					}
				}

			}
		}
	}

	return nil
}

// job 4: send nft:
func (u Usecase) JobMin_SendNftToUser() error {

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

// job 5:
// send btc from segwit address to master address - it does not call our ORD server
func (u Usecase) JobMint_SendFundToMaster() error {

	listToSentMaster, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint(entity.StatusMint_Minting)})

	if len(listToSentMaster) == 0 {
		return nil
	}

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	for _, item := range listToSentMaster {

		if item.PayType == utils.NETWORK_BTC {

			if len(os.Getenv("SECRET_KEY")) == 0 {
				err = errors.New("please config SECRET_KEY")
				u.Logger.Error("u.JobMint_SendFundToMaster.GenerateAddress", err.Error(), err)
				u.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.SECRET_KEY.%s.Error", utils.MASTER_ADDRESS), err.Error(), err)
				continue
			}

			privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
			if err != nil {
				u.Logger.Error("u.JobMint_SendFundToMaster.Decrypt", err.Error(), err)
				u.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.Decrypt.%s.Error", utils.MASTER_ADDRESS), err.Error(), err)
				continue
			}

			// send master now:
			tx, err := bs.SendTransactionWithPreferenceFromSegwitAddress(privateKeyDeCrypt, item.ReceiveAddress, utils.MASTER_ADDRESS, -1, btc.PreferenceMedium)
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.SendTransactionWithPreferenceFromSegwitAddress.%s.Error", utils.MASTER_ADDRESS), err.Error(), err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_SendFundToMaster", item.TableName(), item.Status, "JobMint_SendFundToMaster.SendTransactionWithPreferenceFromSegwitAddress", err.Error(), true)
				continue
			}
			// save tx:
			item.TxSendMaster = tx
			item.Status = entity.StatusMint_SendingFundToMaster // TODO: need to a job to check tx.
			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				u.Logger.Error(fmt.Sprintf("JobBtcSendBtcToMaster.UpdateBtcWalletAddress.%s.Error", utils.MASTER_ADDRESS), err.Error(), err)
				continue
			}
		} else if item.PayType == utils.NETWORK_ETH {
			// TODO: add code to send ETH master address
		}
		time.Sleep(3 * time.Second)
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
