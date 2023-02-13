package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) CreateOrdBTCWalletAddress(rootSpan opentracing.Span, input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("CreateOrdBTCWalletAddress", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("input", input)

	// find Project and make sure index < max supply
	projectID := input.ProjectID
	project, err := u.Repo.FindProjectByProjectIdWithoutCache(projectID)
	if err != nil {
		log.Error("u.Repo.FindProjectByProjectIdWithoutCache", err.Error(), err)
		return nil, err
	}
	if project.MintingInfo.Index >= project.MaxSupply {
		err = fmt.Errorf("project %s is minted out", projectID)
		log.Error("projectIsMintedOut", err.Error(), err)
		return nil, err
	}

	walletAddress := &entity.BTCWalletAddress{}
	err = copier.Copy(walletAddress, input)
	if err != nil {
		log.Error("u.CreateOrdBTCWalletAddress.Copy", err.Error(), err)
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

	log.SetData("CreateOrdBTCWalletAddress.receive", resp)
	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
	if err != nil {
		log.Error("u.CreateOrdBTCWalletAddress.FindProjectByTokenID", err.Error(), err)
		return nil, err
	}

	log.SetData("found.Project", p.ID)
	walletAddress.Amount = p.MintPrice
	walletAddress.UserAddress = userWallet
	walletAddress.OriginUserAddress = input.WalletAddress
	walletAddress.OrdAddress = strings.ReplaceAll(resp.Stdout, "\n", "")
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = ""       //find files from google store
	walletAddress.InscriptionID = "" //find files from google store
	walletAddress.ProjectID = input.ProjectID
	walletAddress.Balance = "0"
	walletAddress.BalanceCheckTime = 0

	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, walletAddress.OrdAddress)
	err = u.Repo.InsertBtcWalletAddress(walletAddress)
	if err != nil {
		log.Error("u.CreateOrdBTCWalletAddress.InsertBtcWalletAddress", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) CreateSegwitBTCWalletAddress(rootSpan opentracing.Span, input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
	return nil, nil
}

func (u Usecase) CheckbalanceWalletAddress(rootSpan opentracing.Span, input structure.CheckBalance) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("CheckbalanceWalletAddress", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, input.Address)
	btc, err := u.Repo.FindBtcWalletAddressByOrd(input.Address)
	if err != nil {
		log.Error("u.Repo.FindBtcWalletAddressByOrd", err.Error(), err)
		return nil, err
	}

	blance, err := u.CheckBalance(span, *btc)
	if err != nil {
		log.Error("u.BalanceLogic", err.Error(), err)
		return nil, err
	}

	return blance, nil
}

func (u Usecase) BTCMint(rootSpan opentracing.Span, input structure.BctMintData) (*ord_service.MintStdoputRespose, *string, error) {
	span, log := u.StartSpan("BTCMint", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("input", input)
	log.SetTag(utils.WALLET_ADDRESS_TAG, input.Address)
	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, input.Address)

	btc, err := u.Repo.FindBtcWalletAddressByOrd(input.Address)
	if err != nil {

		btc = &entity.BTCWalletAddress{}
		eth, err := u.Repo.FindEthWalletAddressByOrd(input.Address)
		if err != nil {
			log.Error("BTCMint.FindEthWalletAddressByOrd", err.Error(), err)
			return nil, nil, err
		}

		err = copier.Copy(btc, eth)
		if err != nil {
			log.Error("BTCMint.copier.Copy", err.Error(), err)
			return nil, nil, err
		}
	}

	//mint logic
	btc, err = u.MintLogic(span, btc)
	if err != nil {
		log.Error("BTCMint.MintLogic", err.Error(), err)
		return nil, nil, err
	}

	// get data from project
	p, err := u.Repo.FindProjectByTokenID(btc.ProjectID)
	if err != nil {
		log.Error("BTCMint.FindProjectByTokenID", err.Error(), err)
		return nil, nil, err
	}
	//log.SetData("found.Project", p)

	//prepare data for mint
	// - Get project.AnimationURL
	projectNftTokenUri := &structure.ProjectAnimationUrl{}
	err = helpers.Base64DecodeRaw(p.NftTokenUri, projectNftTokenUri)
	if err != nil {
		log.Error("BTCMint.helpers.Base64DecodeRaw", err.Error(), err)
		return nil, nil, err
	}

	// - Upload the Animation URL to GCS
	animation := projectNftTokenUri.AnimationUrl
	log.SetData("animation", animation)
	if animation != "" {
		animation = strings.ReplaceAll(animation, "data:text/html;base64,", "")
		now := time.Now().UTC().Unix()
		uploaded, err := u.GCS.UploadBaseToBucket(animation, fmt.Sprintf("btc-projects/%s/%d.html", p.TokenID, now))
		if err != nil {
			log.Error("BTCMint.helpers.Base64DecodeRaw", err.Error(), err)
			return nil, nil, err
		}
		btc.FileURI = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)

	} else {
		images := p.Images
		log.SetData("images", len(images))
		if len(images) > 0 {
			btc.FileURI = images[0]
			newImages := []string{}
			processingImages := p.ProcessingImages
			//remove the project's image out of the current projects
			for i := 1; i < len(images); i++ {
				newImages = append(newImages, images[i])
			}
			processingImages = append(p.ProcessingImages, btc.FileURI)
			p.Images = newImages
			p.ProcessingImages = processingImages
			updated, err := u.Repo.UpdateProject(p.UUID, p)
			if err != nil {
				log.Error("BTCMint.UpdateProject", err.Error(), err)
				return nil, nil, err
			}
			log.SetData("updated", updated)
		}
	}
	//end Animation URL
	if btc.FileURI == "" {
		err = errors.New("There is no file uri to mint")
		log.Error("fileURI.empty", err.Error(), err)
		return nil, nil, err
	}

	mintData := ord_service.MintRequest{
		WalletName: "ord_master",
		FileUrl:    btc.FileURI,
		FeeRate:    15, //temp
		DryRun:     false,
	}

	log.SetData("mintData", mintData)
	resp, err := u.OrdService.Mint(mintData)
	if err != nil {
		log.Error("BTCMint.Mint", err.Error(), err)
		return nil, nil, err
	}

	u.Notify(rootSpan, fmt.Sprintf("[MintFor][projectID %s]", btc.ProjectID), btc.UserAddress, fmt.Sprintf("Made mining transaction for %s, waiting network confirm %s", btc.UserAddress, resp.Stdout))
	tmpText := resp.Stdout
	//tmpText := `{\n  \"commit\": \"7a47732d269d5c005c4df99f2e5cf1e268e217d331d175e445297b1d2991932f\",\n  \"inscription\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2afi0\",\n  \"reveal\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2af\",\n  \"fees\": 185514\n}\n`
	jsonStr := strings.ReplaceAll(tmpText, `\n`, "")
	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")
	btcMintResp := &ord_service.MintStdoputRespose{}

	bytes := []byte(jsonStr)
	err = json.Unmarshal(bytes, btcMintResp)
	if err != nil {
		log.Error("BTCMint.helpers.JsonTransform", err.Error(), err)
		return nil, nil, err
	}

	log.SetTag(utils.TOKEN_ID_TAG, btcMintResp.Inscription)
	return btcMintResp, &btc.FileURI, nil
}

func (u Usecase) ReadGCSFolder(rootSpan opentracing.Span, input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("ReadGCSFolder", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("input", input)
	u.GCS.ReadFolder("btc-projects/project-1/")
	return nil, nil
}

func (u Usecase) UpdateBtcMintedStatus(rootSpan opentracing.Span, btcWallet *entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("UpdateBtcMintedStatus", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("input", btcWallet)

	btcWallet.IsMinted = true
	log.SetTag(utils.WALLET_ADDRESS_TAG, btcWallet.UserAddress)
	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, btcWallet.OrdAddress)
	log.SetTag(utils.TOKEN_ID_TAG, btcWallet.InscriptionID)

	updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(btcWallet.OrdAddress, btcWallet)
	if err != nil {
		log.Error("BTCMint.helpers.UpdateBtcWalletAddressByOrdAddr", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return btcWallet, nil
}

func (u Usecase) GetBalanceSegwitBTCWallet(rootSpan opentracing.Span, userAddress string) (string, error) {
	span, log := u.StartSpan("CheckBalance", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	return "", nil
}

func (u Usecase) GetBalanceOrdBTCWallet(rootSpan opentracing.Span, userAddress string) (string, error) {
	span, log := u.StartSpan("CheckBalance", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	balanceRequest := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userAddress,
			"wallet",
			"balance",
		},
	}

	log.SetData("balanceRequest", balanceRequest)
	//userWallet := helpers.CreateBTCOrdWallet(btc.UserAddress)
	resp, err := u.OrdService.Exec(balanceRequest)
	if err != nil {
		log.Error("BTCMint.Exec.balance", err.Error(), err)
		return "", err
	}

	log.SetData("balanceResponse", resp)
	balance := strings.ReplaceAll(resp.Stdout, "\n", "")
	return balance, nil
}

func (u Usecase) CheckBalance(rootSpan opentracing.Span, btc entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("CheckBalance", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetTag(utils.WALLET_ADDRESS_TAG, btc.UserAddress)
	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, btc.OrdAddress)

	// check ord first
	balance, err := u.GetBalanceOrdBTCWallet(rootSpan, btc.UserAddress)
	if err != nil || balance == "" {
		balance, err = u.GetBalanceSegwitBTCWallet(rootSpan, btc.UserAddress)
		if err != nil || balance == "" {
			return nil, err
		}
	}
	log.SetData("balance", balance)

	btc.Balance = balance
	btc.BalanceCheckTime = btc.BalanceCheckTime + 1
	updated, _ := u.Repo.UpdateBtcWalletAddressByOrdAddr(btc.OrdAddress, &btc)
	log.SetData("updated", updated)

	return &btc, nil
}

func (u Usecase) BalanceLogic(rootSpan opentracing.Span, btc entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("BalanceLogic", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	balance, err := u.CheckBalance(span, btc)
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

	updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(btc.OrdAddress, &btc)
	if err != nil {
		log.Error("u.CheckBalance.updatedStatus", err.Error(), err)
		return nil, err
	}
	log.SetData("updated", updated)
	return &btc, nil
}

func (u Usecase) MintLogic(rootSpan opentracing.Span, btc *entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("MintLogic", rootSpan)
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

func (u Usecase) WaitingForBalancing(rootSpan opentracing.Span) ([]entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("WaitingForBalancing", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	addreses, err := u.Repo.ListProcessingWalletAddress()
	if err != nil {
		log.Error("WillBeProcessWTC.ListProcessingWalletAddress", err.Error(), err)
		return nil, err
	}

	log.SetData("addreses", addreses)
	for _, item := range addreses {
		func(rootSpan opentracing.Span, item entity.BTCWalletAddress) {
			span, log := u.StartSpan(fmt.Sprintf("WaitingForMinted.%s", item.UserAddress), rootSpan)
			defer u.Tracer.FinishSpan(span, log)

			log.SetTag(utils.WALLET_ADDRESS_TAG, item.UserAddress)
			log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, item.OrdAddress)

			newItem, err := u.BalanceLogic(span, item)
			if err != nil {
				log.Error(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}
			log.SetData(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s", item.OrdAddress), newItem)
			u.Notify(rootSpan, fmt.Sprintf("[WaitingForBalancing][projectID %s]", item.ProjectID), item.UserAddress, fmt.Sprintf("%s received BTC %s from [user_address] %s", item.OrdAddress, item.Balance, item.UserAddress))
			updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(item.OrdAddress, newItem)
			if err != nil {
				log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}
			log.SetData("updated", updated)

			if item.MintResponse.Inscription != "" {
				err = errors.New("Token is being minted")
				log.Error("Token.minted", err.Error(), err)
				return
			}

			minResp, fileURI, err := u.BTCMint(span, structure.BctMintData{Address: newItem.OrdAddress})
			if err != nil {
				log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", newItem.OrdAddress), err.Error(), err)
				return
			}

			u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
				TokenID:       minResp.Inscription,
				Commit:        minResp.Commit,
				Reveal:        minResp.Reveal,
				Fees:          minResp.Fees,
				MinterAddress: "ord_master",
				Owner:         "",
				ProjectID:     item.ProjectID,
				Action:        entity.MINT,
				Type:          entity.BIT,
				TraceID:       u.Tracer.TraceID(span),
			})

			newItem.MintResponse = entity.MintStdoputResponse(*minResp)
			newItem.IsMinted = true
			newItem.FileURI = *fileURI
			updated, err = u.Repo.UpdateBtcWalletAddressByOrdAddr(item.OrdAddress, newItem)
			if err != nil {
				log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}
			log.SetData("btc.Minted", minResp)

		}(span, item)

		time.Sleep(2 * time.Second)
	}

	return nil, nil
}

func (u Usecase) WaitingForMinted(rootSpan opentracing.Span) ([]entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("WaitingForMinted", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	addreses, err := u.Repo.ListBTCAddress()
	if err != nil {
		log.Error("WillBeProcessWTC.ListBTCAddress", err.Error(), err)
		return nil, err
	}

	log.SetData("addreses", addreses)
	for _, item := range addreses {
		func(rootSpan opentracing.Span, item entity.BTCWalletAddress) {
			span, log := u.StartSpan(fmt.Sprintf("WaitingForMinted.%s", item.UserAddress), rootSpan)
			defer u.Tracer.FinishSpan(span, log)

			log.SetTag(utils.WALLET_ADDRESS_TAG, item.UserAddress)
			log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, item.OrdAddress)

			addr := item.OriginUserAddress
			if addr == "" {
				addr = item.UserAddress
			}

			sentTokenResp, err := u.SendToken(rootSpan, addr, item.MintResponse.Inscription)
			if err != nil {
				log.Error(fmt.Sprintf("ListenTheMintedBTC.sentToken.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}

			log.SetData(fmt.Sprintf("ListenTheMintedBTC.execResp.%s", item.OrdAddress), sentTokenResp)
			// amout, err := strconv.ParseFloat(item.Amount, 10)
			// if err != nil {
			// 	log.Error("ListenTheMintedBTC.%s. strconv.ParseFloa.Error", err.Error(), err)
			// 	continue
			// }
			// amout = amout * 0.9

			// fundData := ord_service.ExecRequest{
			// 	Args: []string{
			// 		"--wallet",
			// 		item.OrdAddress,
			// 		"send",
			// 		"ord_master",
			// 		fmt.Sprintf("%f", amout),
			// 		"--fee-rate",
			// 		"15",
			// 	},
			// }

			// log.SetData("fundData", fundData)
			// fundResp, err := u.OrdService.Exec(fundData)

			// if err != nil {
			// 	log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.ReFund.Error", item.OrdAddress), err.Error(), err)
			// 	continue
			// }

			u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
				TokenID:       item.MintResponse.Inscription,
				Commit:        item.MintResponse.Commit,
				Reveal:        item.MintResponse.Reveal,
				Fees:          item.MintResponse.Fees,
				MinterAddress: "ord_master",
				Owner:         item.UserAddress,
				Action:        entity.SENT,
				ProjectID:     item.ProjectID,
				Type:          entity.BIT,
				TraceID:       u.Tracer.TraceID(span),
			})

			u.Notify(rootSpan, fmt.Sprintf("[SendToken][ProjectID: %s]", item.ProjectID), addr, item.MintResponse.Inscription)

			// log.SetData("fundResp", fundResp
			item.MintResponse.IsSent = true
			updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(item.OrdAddress, &item)
			if err != nil {
				log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateBtcWalletAddressByOrdAddr.Error", item.OrdAddress), err.Error(), err)
				return
			}

			//TODO: - create entity.TokenURI
			_, err = u.CreateBTCTokenURI(span, item.ProjectID, item.MintResponse.Inscription, item.FileURI, entity.BIT)
			if err != nil {
				log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.CreateBTCTokenURI.Error", item.OrdAddress), err.Error(), err)
				return
			}
			log.SetData("updated", updated)

			err = u.Repo.UpdateTokenOnchainStatusByTokenId(item.MintResponse.Inscription)
			if err != nil {
				log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateTokenOnchainStatusByTokenId.Error", item.OrdAddress), err.Error(), err)
				return
			}

		}(span, item)

		time.Sleep(5 * time.Second)
	}

	return nil, nil
}

func (u Usecase) SendToken(rootSpan opentracing.Span, receiveAddr string, inscriptionID string) (*ord_service.ExecRespose, error) {
	span, log := u.StartSpan(fmt.Sprintf("SendToken.%s", inscriptionID), rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetTag(utils.TOKEN_ID_TAG, inscriptionID)
	log.SetTag(utils.WALLET_ADDRESS_TAG, receiveAddr)
	sendTokenReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			"ord_master",
			"wallet",
			"send",
			receiveAddr,
			inscriptionID,
			"--fee-rate",
			"15",
		}}

	log.SetData("sendTokenReq", sendTokenReq)
	resp, err := u.OrdService.Exec(sendTokenReq)

	if err != nil {
		log.Error("u.OrdService.Exec", err.Error(), err)
		return nil, err
	}

	log.SetData("sendTokenRes", resp)
	return resp, err

}

func (u Usecase) Notify(rootSpan opentracing.Span, title string, userAddress string, content string) {
	span, log := u.StartSpan("SendMessageMint", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	//slack
	preText := fmt.Sprintf("[App: %s][traceID %s] - User address: %s, ", os.Getenv("JAEGER_SERVICE_NAME"), u.Tracer.TraceID(span), userAddress)
	c := fmt.Sprintf("%s", content)

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, c); err != nil {
		log.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
	}
}
