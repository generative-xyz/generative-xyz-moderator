package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	discordclient "rederinghub.io/utils/discord"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) CreateOrdBTCWalletAddress(input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
	u.Logger.Info("input", input)

	// find Project and make sure index < max supply
	projectID := input.ProjectID
	project, err := u.Repo.FindProjectByProjectIdWithoutCache(projectID)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}
	if project.MintingInfo.Index >= project.MaxSupply {
		err = fmt.Errorf("project %s is minted out", projectID)
		u.Logger.Error(err)
		return nil, err
	}

	walletAddress := &entity.BTCWalletAddress{}
	err = copier.Copy(walletAddress, input)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	userWallet := helpers.CreateBTCOrdWallet(input.WalletAddress)

	resp, err := u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userWallet,
			"wallet",
			"create",
		},
	})

	if err != nil {
		u.Logger.Error(err)
		//return nil, err
	} else {
		walletAddress.Mnemonic = resp.Stdout
	}

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
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("CreateOrdBTCWalletAddress.receive", resp)
	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("found.Project", p.ID)
	mintPrice, err := strconv.Atoi(p.MintPrice)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}
	networkFee, err := strconv.Atoi(p.NetworkFee)
	if err == nil {
		mintPrice += networkFee
	}
	walletAddress.Amount = strconv.Itoa(mintPrice)
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

	err = u.Repo.InsertBtcWalletAddress(walletAddress)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) CreateSegwitBTCWalletAddress(input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
	walletAddress := &entity.BTCWalletAddress{}
	privKey, _, addressSegwit, err := btc.GenerateAddressSegwit()
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}
	walletAddress.OrdAddress = addressSegwit //TODO: @thaibao/@tri check this field
	walletAddress.Mnemonic = privKey
	walletAddress.UserAddress = helpers.CreateBTCOrdWallet(input.WalletAddress)
	u.Logger.Info("CreateSegwitBTCWalletAddress.receive", addressSegwit)
	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("found.Project", p.ID)
	mintPrice, err := strconv.Atoi(p.MintPrice)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}
	networkFee, err := strconv.Atoi(p.NetworkFee)
	if err == nil {
		mintPrice += networkFee
	}

	walletAddress.Amount = strconv.Itoa(mintPrice)
	walletAddress.OriginUserAddress = input.WalletAddress
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = ""       //find files from google store
	walletAddress.InscriptionID = "" //find files from google store
	walletAddress.ProjectID = input.ProjectID
	walletAddress.Balance = "0"
	walletAddress.BalanceCheckTime = 0

	err = u.Repo.InsertBtcWalletAddress(walletAddress)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) CheckBalanceWalletAddress(input structure.CheckBalance) (*entity.BTCWalletAddress, error) {

	btc, err := u.Repo.FindBtcWalletAddressByOrd(input.Address)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	balance, err := u.CheckBalance(*btc)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	return balance, nil
}

func (u Usecase) BTCMint(input structure.BctMintData) (*ord_service.MintStdoputRespose, *string, error) {
	eth := &entity.ETHWalletAddress{}
	mintype := entity.BIT
	u.Logger.Info("input", input)

	btc, err := u.Repo.FindBtcWalletAddressByOrd(input.Address)
	if err != nil {
		btc = &entity.BTCWalletAddress{}
		eth, err = u.Repo.FindEthWalletAddressByOrd(input.Address)
		if err != nil {
			u.Logger.Error(err)
			return nil, nil, err
		}

		err = copier.Copy(btc, eth)
		if err != nil {
			u.Logger.Error(err)
			return nil, nil, err
		}

		mintype = entity.ETH
	}

	btc, err = u.MintLogic(btc)
	if err != nil {
		u.Logger.Error(err)
		return nil, nil, err
	}

	// get data from project
	p, err := u.Repo.FindProjectByTokenID(btc.ProjectID)
	if err != nil {
		u.Logger.Error(err)
		return nil, nil, err
	}
	//u.Logger.Info("found.Project", p)

	//prepare data for mint
	// - Get project.AnimationURL
	projectNftTokenUri := &structure.ProjectAnimationUrl{}
	err = helpers.Base64DecodeRaw(p.NftTokenUri, projectNftTokenUri)
	if err != nil {
		u.Logger.Error(err)
		return nil, nil, err
	}

	// - Upload the Animation URL to GCS
	animation := projectNftTokenUri.AnimationUrl
	u.Logger.Info("animation", animation)
	if animation != "" {
		animation = strings.ReplaceAll(animation, "data:text/html;base64,", "")
		now := time.Now().UTC().Unix()
		uploaded, err := u.GCS.UploadBaseToBucket(animation, fmt.Sprintf("btc-projects/%s/%d.html", p.TokenID, now))
		if err != nil {
			u.Logger.Error(err)
			return nil, nil, err
		}
		btc.FileURI = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)

	} else {
		images := p.Images
		u.Logger.Info("images", len(images))
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
		}
	}
	//end Animation URL
	if btc.FileURI == "" {
		err = errors.New("There is no file uri to mint")
		u.Logger.Error(err)
		return nil, nil, err
	}

	baseUrl, err := url.Parse(btc.FileURI)
	if err != nil {
		u.Logger.Error(err)
		return nil, nil, err
	}

	mintData := ord_service.MintRequest{
		WalletName: os.Getenv("ORD_MASTER_ADDRESS"),
		FileUrl:    baseUrl.String(),
		FeeRate:    entity.DEFAULT_FEE_RATE, //temp
		DryRun:     false,
	}

	u.Logger.Info("mintData", mintData)
	resp, err := u.OrdService.Mint(mintData)
	if err != nil {
		u.Logger.Error(err)
		return nil, nil, err
	}
	u.Logger.Info("mint.resp", resp)
	//update btc or eth here
	if mintype == entity.BIT {
		btc.IsMinted = true
		btc.FileURI = baseUrl.String()
		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(btc.OrdAddress, btc)
		if err != nil {
			u.Logger.Error(err)
			return nil, nil, err
		}
		u.Logger.Info("updated", updated)

	} else {
		eth.IsMinted = true
		eth.FileURI = baseUrl.String()
		updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(eth.OrdAddress, eth)
		if err != nil {
			u.Logger.Error(err)
			return nil, nil, err
		}
		u.Logger.Info("updated", updated)
	}

	updated, err := u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		u.Logger.Error(err)
		return nil, nil, err
	}
	u.Logger.Info("project.Updated", updated)

	u.Notify(fmt.Sprintf("[MintFor][%s][projectID %s]", mintype, btc.ProjectID), btc.OrdAddress, fmt.Sprintf("Made mining transaction for %s, waiting network confirm %s", btc.UserAddress, resp.Stdout))
	tmpText := resp.Stdout
	//tmpText := `{\n  \"commit\": \"7a47732d269d5c005c4df99f2e5cf1e268e217d331d175e445297b1d2991932f\",\n  \"inscription\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2afi0\",\n  \"reveal\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2af\",\n  \"fees\": 185514\n}\n`
	jsonStr := strings.ReplaceAll(tmpText, `\n`, "")
	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")
	btcMintResp := &ord_service.MintStdoputRespose{}

	bytes := []byte(jsonStr)
	err = json.Unmarshal(bytes, btcMintResp)
	if err != nil {
		u.Logger.Error(err)
		return nil, nil, err
	}

	if mintype == entity.BIT {
		btc.MintResponse = entity.MintStdoputResponse(*btcMintResp)
		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(btc.OrdAddress, btc)
		if err != nil {
			u.Logger.Error(err)
			return nil, nil, err
		}
		u.Logger.Info("updated", updated)

	} else {
		eth.MintResponse = entity.MintStdoputResponse(*btcMintResp)
		updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(eth.OrdAddress, eth)
		if err != nil {
			u.Logger.Error(err)
			return nil, nil, err
		}
		u.Logger.Info("updated", updated)
	}

	u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
		TokenID:       btcMintResp.Inscription,
		Commit:        btcMintResp.Commit,
		Reveal:        btcMintResp.Reveal,
		Fees:          btcMintResp.Fees,
		MinterAddress: os.Getenv("ORD_MASTER_ADDRESS"),
		Owner:         "",
		ProjectID:     btc.ProjectID,
		Action:        entity.MINT,
		Type:          mintype,
		Balance:       btc.Balance,
		Amount:        btc.Amount,
		ProccessID:    btc.UUID,
	})

	return btcMintResp, &btc.FileURI, nil
}

func (u Usecase) ReadGCSFolder(input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
	u.Logger.Info("input", input)
	u.GCS.ReadFolder("btc-projects/project-1/")
	return nil, nil
}

func (u Usecase) UpdateBtcMintedStatus(btcWallet *entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	u.Logger.Info("input", btcWallet)

	btcWallet.IsMinted = true

	updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(btcWallet.OrdAddress, btcWallet)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("updated", updated)
	return btcWallet, nil
}

func (u Usecase) GetBalanceSegwitBTCWallet(userAddress string) (string, error) {

	u.Logger.Info("userAddress", userAddress)

	_, bs, err := u.buildBTCClient()
	if err != nil {
		u.Logger.Error(err)
		return "", nil
	}
	u.Logger.Info("bs", bs)
	balance, confirm, err := bs.GetBalance(userAddress)
	if err != nil {
		u.Logger.Error(err)
		return "", err
	}
	u.Logger.Info("confirm", confirm)
	u.Logger.Info("balance", balance.String())

	//TODO: @thaibao

	_ = confirm

	return balance.String(), nil
}

func (u Usecase) GetBalanceOrdBTCWallet(userAddress string) (string, error) {
	balanceRequest := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userAddress,
			"wallet",
			"balance",
		},
	}

	u.Logger.Info("balanceRequest", balanceRequest)
	//userWallet := helpers.CreateBTCOrdWallet(btc.UserAddress)
	resp, err := u.OrdService.Exec(balanceRequest)
	if err != nil {
		u.Logger.Error(err)
		return "", err
	}

	u.Logger.Info("balanceResponse", resp)
	balance := strings.ReplaceAll(resp.Stdout, "\n", "")
	return balance, nil
}

func (u Usecase) CheckBalance(btc entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {

	//TODO - removed checking ORD, only Segwit is used
	balance, err := u.GetBalanceSegwitBTCWallet(btc.OrdAddress)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	if balance == "" {
		err := errors.New("balance is empty")
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("balance", balance)
	btc.Balance = strings.ReplaceAll(balance, `\n`, "")
	btc.BalanceCheckTime = btc.BalanceCheckTime + 1
	updated, _ := u.Repo.UpdateBtcWalletAddressByOrdAddr(btc.OrdAddress, &btc)
	u.Logger.Info("updated", btc)
	u.Logger.Info("updated", updated)
	return &btc, nil
}

func (u Usecase) BalanceLogic(btc entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	balance, err := u.CheckBalance(btc)
	if err != nil {
		u.Logger.Error(err)
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
		u.Logger.Error(err)
		return nil, err
	}
	u.Logger.Info("updated", updated)
	return &btc, nil
}

func (u Usecase) MintLogic(btc *entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	var err error

	//if this was minted, skip it
	if btc.IsMinted {
		err = errors.New("This btc was minted")
		u.Logger.Error(err)
		return nil, err
	}

	if !btc.IsConfirm {
		err = errors.New("This btc must be IsConfirmed")
		u.Logger.Error(err)
		return nil, err
	}
	if btc.MintResponse.Inscription != "" {
		err = errors.New(fmt.Sprintf("This btc has Inscription %s", btc.MintResponse.Inscription))
		u.Logger.Error(err)
		return nil, err
	}
	if btc.MintResponse.Reveal != "" {
		err = errors.New(fmt.Sprintf("This btc has revealID %s", btc.MintResponse.Reveal))
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("btc", btc)
	return btc, nil
}

// Mint flow
func (u Usecase) WaitingForBalancing() ([]entity.BTCWalletAddress, error) {
	addreses, err := u.Repo.ListProcessingWalletAddress()
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("addreses", addreses)
	for _, item := range addreses {
		func(item entity.BTCWalletAddress) {

			newItem, err := u.BalanceLogic(item)
			if err != nil {
				u.Logger.Error(err)
				return
			}
			u.Logger.Info(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s", item.OrdAddress), newItem)
			u.Notify(fmt.Sprintf("[WaitingForBalancing][projectID %s]", item.ProjectID), item.UserAddress, fmt.Sprintf("%s checkint BTC %s from [user_address] %s", item.OrdAddress, item.Balance, item.UserAddress))

			updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(item.OrdAddress, newItem)
			if err != nil {
				u.Logger.Error(err)
				return
			}
			u.Logger.Info("updated", updated)
			u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
				MinterAddress: os.Getenv("ORD_MASTER_ADDRESS"),
				Owner:         "",
				ProjectID:     item.ProjectID,
				Action:        entity.BLANCE,
				Type:          entity.BIT,
				Balance:       item.Balance,
				Amount:        item.Amount,
				ProccessID:    item.UUID,
			})
		}(item)

		time.Sleep(2 * time.Second)
	}

	return nil, nil
}

func (u Usecase) WaitingForMinting() ([]entity.BTCWalletAddress, error) {
	addreses, err := u.Repo.ListMintingWalletAddress()
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("addreses", addreses)
	for _, item := range addreses {
		func(item entity.BTCWalletAddress) {

			if item.MintResponse.Inscription != "" {
				err = errors.New("Token is being minted")
				u.Logger.Error(err)
				return
			}

			_, _, err := u.BTCMint(structure.BctMintData{Address: item.OrdAddress})
			if err != nil {
				u.Notify(fmt.Sprintf("[Error][MintFor][projectID %s]", item.ProjectID), item.OrdAddress, err.Error())
				u.Logger.Error(err)
				return
			}

		}(item)

		time.Sleep(2 * time.Second)
	}

	return nil, nil
}

func (u Usecase) WaitingForMinted() ([]entity.BTCWalletAddress, error) {

	_, bs, err := u.buildBTCClient()

	addreses, err := u.Repo.ListBTCAddress()
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("addreses", addreses)
	for _, item := range addreses {
		func(item entity.BTCWalletAddress) {

			addr := item.OriginUserAddress
			if addr == "" {
				addr = item.UserAddress
			}

			//check token is created or not via BlockcypherService
			txInfo, err := bs.CheckTx(item.MintResponse.Reveal)
			if err != nil {
				u.Logger.Error(err)
				u.Notify(fmt.Sprintf("[Error][BTC][SendToken.bs.CheckTx][projectID %s]", item.ProjectID), item.InscriptionID, fmt.Sprintf("%s, object: %s", err.Error(), item.UUID))
				return
			}
			u.Logger.Info("txInfo", txInfo)
			if txInfo.Confirmations > 1 {
				sentTokenResp, err := u.SendToken(addr, item.MintResponse.Inscription)
				if err != nil {
					u.Notify(fmt.Sprintf("[Error][BTC][SendToken][projectID %s]", item.ProjectID), item.InscriptionID, fmt.Sprintf("%s, object: %s", err.Error(), item.UUID))
					u.Logger.Error(err)
					return
				}

				u.Logger.Info(fmt.Sprintf("ListenTheMintedBTC.execResp.%s", item.OrdAddress), sentTokenResp)

				u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
					TokenID:       item.MintResponse.Inscription,
					Commit:        item.MintResponse.Commit,
					Reveal:        item.MintResponse.Reveal,
					Fees:          item.MintResponse.Fees,
					MinterAddress: os.Getenv("ORD_MASTER_ADDRESS"),
					Owner:         item.UserAddress,
					Action:        entity.SENT,
					ProjectID:     item.ProjectID,
					Type:          entity.BIT,
					Balance:       item.Balance,
					Amount:        item.Amount,
					ProccessID:    item.UUID,
				})

				u.Notify(fmt.Sprintf("[SendToken][ProjectID: %s]", item.ProjectID), addr, item.MintResponse.Inscription)

				// u.Logger.Info("fundResp", fundResp
				item.MintResponse.IsSent = true
				updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(item.OrdAddress, &item)
				if err != nil {
					u.Logger.Error(err)
					return
				}

				//TODO: - create entity.TokenURI
				_, err = u.CreateBTCTokenURI(item.ProjectID, item.MintResponse.Inscription, item.FileURI, entity.BIT)
				if err != nil {
					u.Logger.Error(err)
					return
				}
				u.Logger.Info("updated", updated)
				err = u.Repo.UpdateTokenOnchainStatusByTokenId(item.MintResponse.Inscription)
				if err != nil {
					u.Logger.Error(err)
					return
				}
				go u.CreateMintActivity(item.InscriptionID, item.Amount)
				go u.NotifyNFTMinted(item.OriginUserAddress, item.InscriptionID, item.MintResponse.Fees)

			} else {
				u.Logger.Info("checkTx.Inscription.Existed", false)
			}
		}(item)

		time.Sleep(5 * time.Second)
	}

	return nil, nil
}

func (u Usecase) NotifyNFTMinted(btcUserAddr string, inscriptionID string, networkFee int) {
	domain := os.Getenv("DOMAIN")
	webhook := os.Getenv("DISCORD_NFT_MINTED_WEBHOOK")
	u.Logger.Info(
		"NotifyNFTMinted",
		zap.String("btcUserAddr", btcUserAddr),
		zap.String("inscriptionID", inscriptionID),
		zap.Int("networkFee", networkFee),
	)

	tokenUri, err := u.Repo.FindTokenByTokenID(inscriptionID)
	if err != nil {
		u.Logger.ErrorAny("NotifyNFTMinted.FindTokenByTokenID failed", zap.Any("err", err.Error()))
		return
	}

	var minterDisplayName string
	minterAddress := btcUserAddr
	{
		minter, err := u.Repo.FindUserByBtcAddress(btcUserAddr)
		if err == nil {
			minterDisplayName = minter.DisplayName
		} else {
			u.Logger.ErrorAny("NotifyNFTMinted.FindUserByBtcAddress for minter failed", zap.Any("err", err.Error()))
		}
	}

	if tokenUri.Creator == nil {
		u.Logger.ErrorAny("NotifyNFTMinted.tokenUri.CreatorIsEmpty", zap.Any("tokenID", tokenUri.TokenID))
		return
	}

	project, err := u.GetProjectByGenNFTAddr(tokenUri.ProjectID)
	if err != nil {
		u.Logger.ErrorAny("NotifyNFTMinted.GetProjectByGenNFTAddr failed", zap.Any("err", err))
		return
	}
	var category, description string
	if len(project.Categories) > 0 {
		// we assume that there are only one category
		categoryEntity, err := u.GetCategory(project.Categories[0])
		if err != nil {
			u.Logger.ErrorAny("NotifyNFTMinted.GetCategory failed", zap.Any("err", err))
			return
		}
		category = categoryEntity.Name
		description = fmt.Sprintf("Category: %s\n", category)
	}

	ownerName := u.resolveShortName(tokenUri.Creator.DisplayName, tokenUri.Creator.WalletAddress)
	collectionName := project.Name
	// itemCount := project.MaxSupply
	mintedCount := tokenUri.OrderInscriptionIndex

	fields := make([]discordclient.Field, 0)
	addFields := func(fields []discordclient.Field, name string, value string, inline bool) []discordclient.Field {
		if value == "" {
			return fields
		}
		return append(fields, discordclient.Field{
			Name:   name,
			Value:  value,
			Inline: inline,
		})
	}
	fields = addFields(fields, "", project.Description, false)
	fields = addFields(fields, "Mint Price", u.resolveMintPriceBTC(project.MintPrice), true)
	fields = addFields(fields, "Collector", fmt.Sprintf("[%s](%s)",
		u.resolveShortName(minterDisplayName, btcUserAddr),
		fmt.Sprintf("%s/profile/%s", domain, minterAddress),
	), true)

	// fields = addFields(fields, "Minted", fmt.Sprintf("%d/%d", mintedCount, itemCount), true)
	//fields = addFields(fields, "Network Fee", strconv.FormatFloat(float64(networkFee)/1e8, 'f', -1, 64)+" BTC")

	discordMsg := discordclient.Message{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW MINT**",
		Embeds: []discordclient.Embed{{
			Title:       fmt.Sprintf("%s\n***%s #%d***", ownerName, collectionName, mintedCount),
			Url:         fmt.Sprintf("%s/generative/%s/%s", domain, project.GenNFTAddr, tokenUri.TokenID),
			Description: description,
			//Author: discordclient.Author{
			//	Name:    u.resolveShortName(minter.DisplayName, minter.WalletAddress),
			//	Url:     fmt.Sprintf("%s/profile/%s", domain, minter.WalletAddress),
			//	IconUrl: minter.Avatar,
			//},
			Fields: fields,
			Image: discordclient.Image{
				Url: tokenUri.Thumbnail,
			},
		}},
	}
	sendCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	u.Logger.Info("sending message to discord", discordMsg)

	if err := u.DiscordClient.SendMessage(sendCtx, webhook, discordMsg); err != nil {
		u.Logger.Error("error sending message to discord", err)
	}
}

//End Mint flow

func (u Usecase) SendToken(receiveAddr string, inscriptionID string) (*ord_service.ExecRespose, error) {

	sendTokenReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			os.Getenv("ORD_MASTER_ADDRESS"),
			"wallet",
			"send",
			receiveAddr,
			inscriptionID,
			"--fee-rate",
			fmt.Sprintf("%d", entity.DEFAULT_FEE_RATE),
		}}

	u.Logger.Info("sendTokenReq", sendTokenReq)
	resp, err := u.OrdService.Exec(sendTokenReq)

	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("sendTokenRes", resp)
	return resp, err

}

func (u Usecase) Notify(title string, userAddress string, content string) {

	//slack
	preText := fmt.Sprintf("[App: %s][traceID %s] - User address: %s, ", os.Getenv("JAEGER_SERVICE_NAME"), "", userAddress)
	c := fmt.Sprintf("%s", content)

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, c); err != nil {
		u.Logger.Error(err)
	}
}

func (u Usecase) NotifyWithChannel(channelID string, title string, userAddress string, content string) {
	//slack
	preText := fmt.Sprintf("[App: %s][traceID %s] - User address: %s, ", os.Getenv("JAEGER_SERVICE_NAME"), "", userAddress)
	c := fmt.Sprintf("%s", content)

	if _, _, err := u.Slack.SendMessageToSlackWithChannel(channelID, preText, title, c); err != nil {
		u.Logger.Error(err)
	}
}

// phuong:
// send btc from segwit address to master address - it does not call our ORD server
func (u Usecase) JobBtcSendBtcToMaster() error {

	addreses, err := u.Repo.ListWalletAddressToClaimBTC()
	if err != nil {
		u.Logger.Error(err)
		return err
	}
	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	u.Logger.Info("addreses", addreses)
	for _, item := range addreses {

		// send master now:
		tx, err := bs.SendTransactionWithPreferenceFromSegwitAddress(item.Mnemonic, item.OrdAddress, utils.MASTER_ADDRESS, -1, btc.PreferenceMedium)
		if err != nil {
			u.Logger.Error(err)
			continue
		}
		// save tx:
		item.TxSendMaster = tx
		item.IsSentMaster = true
		_, err = u.Repo.UpdateBtcWalletAddress(&item)
		if err != nil {
			u.Logger.Error(err)
			continue
		}

		time.Sleep(3 * time.Second)
	}

	return nil
}

func (u Usecase) GetCurrentMintingByWalletAddress(address string) ([]structure.MintingInscription, error) {
	result := []structure.MintingInscription{}
	listBTC, err := u.Repo.ListMintingWaitingForFundByWalletAddress(address)
	if err != nil {
		return nil, err
	}

	listBTC1, err := u.Repo.ListMintingByWalletAddress(address)
	if err != nil {
		return nil, err
	}

	listBTC2, err := u.Repo.ListMintingWaitingToSendByWalletAddress(address)
	if err != nil {
		return nil, err
	}
	listBTC = append(listBTC, listBTC1...)
	listBTC = append(listBTC, listBTC2...)

	listETH, err := u.Repo.ListMintingETHByWalletAddress(address)
	if err != nil {
		return nil, err
	}

	listETH1, err := u.Repo.ListMintingETHByWalletAddress(address)
	if err != nil {
		return nil, err
	}

	listETH2, err := u.Repo.ListMintingWaitingToSendETHByWalletAddress(address)
	if err != nil {
		return nil, err
	}
	listETH = append(listETH, listETH1...)
	listETH = append(listETH, listETH2...)

	listMintV2, err := u.Repo.ListMintNftBtcByStatusAndAddress(address, []entity.StatusMint{entity.StatusMint_Pending, entity.StatusMint_WaitingForConfirms, entity.StatusMint_ReceivedFund, entity.StatusMint_Minting, entity.StatusMint_Minted, entity.StatusMint_SendingNFTToUser, entity.StatusMint_NeedToRefund, entity.StatusMint_Refunding, entity.StatusMint_TxRefundFailed, entity.StatusMint_TxMintFailed})
	if err != nil {
		return nil, err
	}

	itemIDMap := make(map[string]struct{})

	for _, item := range listBTC {
		projectInfo, err := u.Repo.FindProjectByTokenID(item.ProjectID)
		if err != nil {
			return nil, err
		}
		if _, ok := itemIDMap[item.UUID]; ok {
			continue
		}
		var minting *structure.MintingInscription
		if time.Since(*item.CreatedAt) >= 2*time.Hour {
			continue // timeout if  waited for 2 hours
		}
		if !item.IsConfirm {
			minting = &structure.MintingInscription{
				ID:           item.UUID,
				CreatedAt:    item.CreatedAt,
				Status:       "Waiting for payment",
				FileURI:      item.FileURI,
				ProjectID:    item.ProjectID,
				ProjectImage: projectInfo.Thumbnail,
				ProjectName:  projectInfo.Name,
			}
		} else {
			if !item.IsMinted {
				minting = &structure.MintingInscription{
					ID:           item.UUID,
					CreatedAt:    item.CreatedAt,
					Status:       "Minting",
					FileURI:      item.FileURI,
					ProjectID:    item.ProjectID,
					ProjectImage: projectInfo.Thumbnail,
					ProjectName:  projectInfo.Name,
				}
			} else {
				minting = &structure.MintingInscription{
					ID:            item.UUID,
					CreatedAt:     item.CreatedAt,
					Status:        "Transferring",
					FileURI:       item.FileURI,
					ProjectID:     item.ProjectID,
					ProjectImage:  projectInfo.Thumbnail,
					ProjectName:   projectInfo.Name,
					InscriptionID: item.InscriptionID,
				}
			}
		}
		itemIDMap[item.UUID] = struct{}{}
		result = append(result, *minting)
	}

	for _, item := range listETH {
		projectInfo, err := u.Repo.FindProjectByTokenID(item.ProjectID)
		if err != nil {
			return nil, err
		}
		if _, ok := itemIDMap[item.UUID]; ok {
			continue
		}
		if time.Since(*item.CreatedAt) >= 2*time.Hour {
			continue // timeout if  waited for 2 hours
		}
		var minting *structure.MintingInscription
		if !item.IsConfirm {
			minting = &structure.MintingInscription{
				ID:           item.UUID,
				CreatedAt:    item.CreatedAt,
				Status:       "Waiting for payment",
				FileURI:      item.FileURI,
				ProjectID:    item.ProjectID,
				ProjectImage: projectInfo.Thumbnail,
				ProjectName:  projectInfo.Name,
			}
		} else {
			if !item.IsMinted {
				minting = &structure.MintingInscription{
					ID:           item.UUID,
					CreatedAt:    item.CreatedAt,
					Status:       "Minting",
					FileURI:      item.FileURI,
					ProjectID:    item.ProjectID,
					ProjectImage: projectInfo.Thumbnail,
					ProjectName:  projectInfo.Name,
				}
			} else {
				minting = &structure.MintingInscription{
					ID:            item.UUID,
					CreatedAt:     item.CreatedAt,
					Status:        "Transferring",
					FileURI:       item.FileURI,
					ProjectID:     item.ProjectID,
					ProjectImage:  projectInfo.Thumbnail,
					ProjectName:   projectInfo.Name,
					InscriptionID: item.InscriptionID,
				}
			}
		}
		itemIDMap[item.UUID] = struct{}{}
		result = append(result, *minting)
	}

	for _, item := range listMintV2 {
		projectInfo, err := u.Repo.FindProjectByTokenID(item.ProjectID)
		if err != nil {
			return nil, err
		}

		status := ""
		if time.Since(item.ExpiredAt) >= 1*time.Second && item.Status == entity.StatusMint_Pending {
			continue
		}
		if (item.Status) == -1 {
			continue
		}
		switch item.Status {
		case entity.StatusMint_NeedToRefund, entity.StatusMint_TxRefundFailed:
			status = entity.StatusMintToText[entity.StatusMint_Refunding]
		case entity.StatusMint_TxMintFailed:
			status = entity.StatusMintToText[entity.StatusMint_Minting]
		default:
			status = entity.StatusMintToText[item.Status]
		}

		if item.PayType == "eth" {
			if item.Status == entity.StatusMint_Refunded {
				status = entity.StatusMintToText[entity.StatusMint_Refunding]
			}
		}

		minting := structure.MintingInscription{
			ID:            item.UUID,
			CreatedAt:     item.CreatedAt,
			Status:        status,
			FileURI:       item.FileURI,
			ProjectID:     item.ProjectID,
			ProjectImage:  projectInfo.Thumbnail,
			ProjectName:   projectInfo.Name,
			InscriptionID: item.InscriptionID,
			IsCancel:      int(item.Status) == 0,
			Quantity:      item.Quantity,
		}
		result = append(result, minting)
	}

	return result, nil
}
