package usecase

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	discordclient "rederinghub.io/utils/discord"
	"rederinghub.io/utils/logger"
)

const (
	PerceptronProjectID = "1002573"
	PFPsCategory        = "PFPs"
)

type addUserDiscordFieldReq struct {
	Fields  []entity.Field
	Key     string
	Address string
	UserID  string
	Inline  bool
	Domain  string
}

func addDiscordField(fields []entity.Field, name string, value string, inline bool) []entity.Field {
	if value == "" {
		return fields
	}
	return append(fields, entity.Field{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
}

func (u Usecase) addUserDiscordField(req addUserDiscordFieldReq) []entity.Field {
	var user *entity.Users
	var err error
	if req.Address != "" {
		user, err = u.Repo.FindUserByAddress(req.Address)
	} else if req.UserID != "" {
		user, err = u.Repo.FindUserByID(req.UserID)
	}
	var userStr string
	if err == nil && user != nil {
		address := user.WalletAddressBTC
		if address == "" {
			address = user.WalletAddress
		}
		if address == "" {
			address = user.WalletAddressBTCTaproot
		}
		userStr = fmt.Sprintf("[%s](%s)",
			u.resolveShortName(user.DisplayName, address),
			fmt.Sprintf("%s/profile/%s", req.Domain, address),
		)
	} else {
		logger.AtLog.Logger.Error("NotifyNewSale.FindUserByAddress")
		userStr = fmt.Sprintf("[%s](%s)",
			u.resolveShortName("", req.Address),
			fmt.Sprintf("%s/profile/%s", req.Domain, req.Address),
		)
	}
	if userStr != "" {
		return addDiscordField(req.Fields, req.Key, userStr, req.Inline)
	} else {
		return req.Fields
	}
}

func (u Usecase) NotifyNewAirdrop(airdrop *entity.Airdrop) error {
	domain := os.Getenv("DOMAIN")
	fields := make([]entity.Field, 0)
	file := strings.Replace(airdrop.File, "html", "png", 1)

	fields = u.addUserDiscordField(addUserDiscordFieldReq{
		Fields: fields,
		Key:    "Key holder",
		UserID: airdrop.Receiver,
		Inline: false,
		Domain: domain,
	})

	// fields = addFields(fields, "File", file, false)

	var title string
	if airdrop.File == utils.AIRDROP_MAGIC {
		title = "MAGIC KEY"
	} else if airdrop.File == utils.AIRDROP_GOLDEN {
		title = "GOLDEN KEY"
	} else {
		title = "SILVER KEY"
	}

	inscriptionNumTitle := ""

	inscriptionInfo, err := u.GetInscribeInfo(airdrop.InscriptionId)

	if err == nil && inscriptionInfo != nil {
		inscriptionNumTitle = fmt.Sprintf(" #%v", inscriptionInfo.Index)
	} else {
		logger.AtLog.Logger.Error("ErrorWhenGetInscribeInfo", zap.Any("inscriptionId", airdrop.InscriptionId))
	}

	parsedThumbnailUrl, err := url.Parse(file)
	if err != nil {
		logger.AtLog.Logger.Error("ErrorParseProjectThumbnailURL", zap.Error(err))
	}
	parsedThumbnail := parsedThumbnailUrl.String()

	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW KEY**",
		Embeds: []entity.Embed{{
			Title:  fmt.Sprintf("%s%s", title, inscriptionNumTitle),
			Url:    "https://generative.xyz",
			Fields: fields,
			Image: entity.Image{
				Url: parsedThumbnail,
			},
		}},
	}

	logger.AtLog.Logger.Info("sending new airdrop message to discord", zap.Any("message", zap.Any("discordMsg)", discordMsg)))

	noti := entity.DiscordNoti{
		Message:    discordMsg,
		NumRetried: 0,
		Status:     entity.PENDING,
		Type:       entity.NEW_AIRDROP,
		Meta: entity.DiscordNotiMeta{
			InscriptionID: airdrop.InscriptionId,
		},
	}

	// create discord message
	err = u.CreateDiscordNoti(noti)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNewAirdrop.CreateDiscordNoti", zap.Error(err))
		return err
	}
	return nil
}

func (u Usecase) NotifyNewSale(order entity.DexBTCListing, buyerAddress string) error {
	logger.AtLog.Logger.Info("NotifyNewSale.Start", zap.Any("order", order), zap.Any("buyerAddress", zap.Any("buyerAddress)", buyerAddress)))
	domain := os.Getenv("DOMAIN")

	tokenUri, err := u.Repo.FindTokenByTokenID(order.InscriptionID)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNFTMinted.FindTokenByTokenID failed", zap.Any("err", err.Error()))
		return err
	}

	project, err := u.GetProjectByGenNFTAddr(tokenUri.ProjectID)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNFTMinted.GetProjectByGenNFTAddr failed", zap.Any("err", err))
		return err
	}

	var category, description string
	if len(project.Categories) > 0 {
		// we assume that there are only one category
		categoryEntity, err := u.GetCategory(project.Categories[0])
		if err != nil {
			logger.AtLog.Logger.Error("NotifyNFTMinted.GetCategory failed", zap.Any("err", err))
			return err
		}
		category = categoryEntity.Name
		description = fmt.Sprintf("Category: %s\n", category)
	}

	ownerName := u.resolveShortName(project.CreatorProfile.DisplayName, project.CreatorProfile.WalletAddress)
	collectionName := project.Name
	mintedCount := tokenUri.OrderInscriptionIndex

	fields := make([]entity.Field, 0)

	fields = addDiscordField(fields, "Sale Price", u.resolveMintPriceBTC(fmt.Sprintf("%v", order.Amount)), true)

	if buyerAddress != "" {
		fields = u.addUserDiscordField(addUserDiscordFieldReq{
			Fields:  fields,
			Key:     "Buyer",
			Address: buyerAddress,
			Inline:  true,
			Domain:  domain,
		})
	}

	fields = u.addUserDiscordField(addUserDiscordFieldReq{
		Fields:  fields,
		Key:     "Seller",
		Address: order.SellerAddress,
		Inline:  true,
		Domain:  domain,
	})

	parsedThumbnailUrl, err := url.Parse(tokenUri.Thumbnail)
	if err != nil {
		logger.AtLog.Logger.Error("ErrorParseProjectThumbnailURL", zap.Error(err))
	}
	parsedThumbnail := parsedThumbnailUrl.String()

	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW SALE**",
		Embeds: []entity.Embed{{
			Title:       fmt.Sprintf("%s\n***%s #%d***", ownerName, collectionName, mintedCount),
			Url:         fmt.Sprintf("%s/generative/%s/%s", domain, project.GenNFTAddr, tokenUri.TokenID),
			Description: description,
			Fields:      fields,
			Image: entity.Image{
				Url: parsedThumbnail,
			},
		}},
	}

	logger.AtLog.Logger.Info("sending new sale message to discord", zap.Any("message", zap.Any("discordMsg)", discordMsg)))

	channels := []entity.DiscordNotiType{
		entity.DISCORD_SALE_CHANNEL,
	}

	if tokenUri.ProjectID == PerceptronProjectID {
		channels = append(channels, entity.DISCORD_PERCEPTRON_CHANNEL)
		time.Sleep(1 * time.Minute)
	}

	if order.Amount > 0 {
		if category == PFPsCategory {
			channels = append(channels, entity.DISCORD_PFPS_CHANNEL)
		} else {
			channels = append(channels, entity.DISCORD_ART_CHANNEL)
		}
	}

	for _, c := range channels {
		noti := entity.DiscordNoti{
			Message:    discordMsg,
			NumRetried: 0,
			Status:     entity.PENDING,
			Type:       c,
			Meta: entity.DiscordNotiMeta{
				ProjectID:     project.TokenID,
				InscriptionID: tokenUri.TokenID,
			},
		}

		// create discord message
		err = u.CreateDiscordNoti(noti)
		if err != nil {
			logger.AtLog.Logger.Error("NotifyNewSale.CreateDiscordNoti", zap.Error(err))
		}
	}
	return nil
}

func (u Usecase) NotifyNewListing(order entity.DexBTCListing) error {
	logger.AtLog.Logger.Info("NotifyNewListing.Start", zap.Any("order", zap.Any("order)", order)))
	domain := os.Getenv("DOMAIN")

	tokenUri, err := u.Repo.FindTokenByTokenID(order.InscriptionID)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNFTMinted.FindTokenByTokenID failed", zap.Any("err", err.Error()))
		return err
	}

	project, err := u.GetProjectByGenNFTAddr(tokenUri.ProjectID)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNFTMinted.GetProjectByGenNFTAddr failed", zap.Any("err", err))
		return err
	}

	var category, description string
	if len(project.Categories) > 0 {
		// we assume that there are only one category
		categoryEntity, err := u.GetCategory(project.Categories[0])
		if err != nil {
			logger.AtLog.Logger.Error("NotifyNFTMinted.GetCategory failed", zap.Any("err", err))
			return err
		}
		category = categoryEntity.Name
		description = fmt.Sprintf("Category: %s\n", category)
	}

	ownerName := u.resolveShortName(project.CreatorProfile.DisplayName, project.CreatorProfile.WalletAddress)
	collectionName := project.Name
	mintedCount := tokenUri.OrderInscriptionIndex

	fields := make([]entity.Field, 0)

	fields = addDiscordField(fields, "", u.resolveShortDescription(project.Description), false)

	fields = addDiscordField(fields, "List Price", u.resolveMintPriceBTC(fmt.Sprintf("%v", order.Amount)), true)

	fields = u.addUserDiscordField(addUserDiscordFieldReq{
		Fields:  fields,
		Key:     "Seller",
		Address: order.SellerAddress,
		Inline:  true,
		Domain:  domain,
	})

	parsedThumbnailUrl, err := url.Parse(tokenUri.Thumbnail)
	if err != nil {
		logger.AtLog.Logger.Error("ErrorParseProjectThumbnailURL", zap.Error(err))
	}
	parsedThumbnail := parsedThumbnailUrl.String()

	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW LISTING**",
		Embeds: []entity.Embed{{
			Title:       fmt.Sprintf("%s\n***%s #%d***", ownerName, collectionName, mintedCount),
			Url:         fmt.Sprintf("%s/generative/%s/%s", domain, project.GenNFTAddr, tokenUri.TokenID),
			Description: description,
			Fields:      fields,
			Image: entity.Image{
				Url: parsedThumbnail,
			},
		}},
	}

	logger.AtLog.Logger.Info("sending new new listing message to discord", zap.Any("message", zap.Any("discordMsg)", discordMsg)))
	noti := entity.DiscordNoti{
		Message:    discordMsg,
		NumRetried: 0,
		Status:     entity.PENDING,
		Type:       entity.NEW_LISTING,
		Meta: entity.DiscordNotiMeta{
			ProjectID:     project.TokenID,
			InscriptionID: tokenUri.TokenID,
		},
	}

	// create discord message
	err = u.CreateDiscordNoti(noti)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNewListing.CreateDiscordNoti", zap.Error(err))
		return err
	}
	return nil
}

func (u Usecase) NotifyNFTMinted(btcUserAddr string, inscriptionID string) {
	domain := os.Getenv("DOMAIN")
	logger.AtLog.Logger.Info(
		"NotifyNFTMinted",
		zap.String("btcUserAddr", btcUserAddr),
		zap.String("inscriptionID", inscriptionID),
	)

	tokenUri, err := u.Repo.FindTokenByTokenID(inscriptionID)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNFTMinted.FindTokenByTokenID failed", zap.Any("err", err.Error()))
		return
	}

	var minterDisplayName string
	minterAddress := btcUserAddr
	{
		minter, err := u.Repo.FindUserByBtcAddress(btcUserAddr)
		if err == nil {
			minterDisplayName = minter.DisplayName
		} else {
			logger.AtLog.Logger.Error("NotifyNFTMinted.FindUserByBtcAddress for minter failed", zap.Any("err", err.Error()))
		}
	}

	if tokenUri.Creator == nil {
		logger.AtLog.Logger.Error("NotifyNFTMinted.tokenUri.CreatorIsEmpty", zap.Any("tokenID", tokenUri.TokenID))
		return
	}

	project, err := u.GetProjectByGenNFTAddr(tokenUri.ProjectID)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNFTMinted.GetProjectByGenNFTAddr failed", zap.Any("err", err))
		return
	}
	var category, description string
	if len(project.Categories) > 0 {
		// we assume that there are only one category
		categoryEntity, err := u.GetCategory(project.Categories[0])
		if err != nil {
			logger.AtLog.Logger.Error("NotifyNFTMinted.GetCategory failed", zap.Any("err", err))
			return
		}
		category = categoryEntity.Name
		description = fmt.Sprintf("Category: %s\n", category)
	}

	ownerName := u.resolveShortName(tokenUri.Creator.DisplayName, tokenUri.Creator.WalletAddress)
	collectionName := project.Name
	// itemCount := project.MaxSupply
	mintedCount := tokenUri.OrderInscriptionIndex

	fields := make([]entity.Field, 0)
	addFields := func(fields []entity.Field, name string, value string, inline bool) []entity.Field {
		if value == "" {
			return fields
		}
		return append(fields, entity.Field{
			Name:   name,
			Value:  value,
			Inline: inline,
		})
	}

	mintNftBtc, err := u.Repo.FindMintNftBtcByInscriptionID(inscriptionID)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNFTMinted.FindMintNftBtcByInscriptionID failed", zap.Any("err", err))
		return
	}

	mintPrice := project.MintPrice
	if v, ok := mintNftBtc.EstFeeInfo["btc"]; ok {
		mintPrice = v.MintPrice
	}
	fields = addFields(fields, "Mint Price", u.resolveMintPriceBTC(mintPrice), true)

	fields = addFields(fields, "Collector", fmt.Sprintf("[%s](%s)",
		u.resolveShortName(minterDisplayName, btcUserAddr),
		fmt.Sprintf("%s/profile/%s", domain, minterAddress),
	), true)

	// fields = addFields(fields, "Minted", fmt.Sprintf("%d/%d", mintedCount, itemCount), true)

	parsedThumbnailUrl, err := url.Parse(tokenUri.Thumbnail)
	if err != nil {
		logger.AtLog.Logger.Error("ErrorParseProjectThumbnailURL", zap.Error(err))
	}
	parsedThumbnail := parsedThumbnailUrl.String()

	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW MINT**",
		Embeds: []entity.Embed{{
			Title:       fmt.Sprintf("%s\n***%s #%d***", ownerName, collectionName, mintedCount),
			Url:         fmt.Sprintf("%s/generative/%s/%s", domain, project.GenNFTAddr, tokenUri.TokenID),
			Description: description,
			//Author: discordclient.Author{
			//	Name:    u.resolveShortName(minter.DisplayName, minter.WalletAddress),
			//	Url:     fmt.Sprintf("%s/profile/%s", domain, minter.WalletAddress),
			//	IconUrl: minter.Avatar,
			//},
			Fields: fields,
			Image: entity.Image{
				Url: parsedThumbnail,
			},
		}},
	}

	logger.AtLog.Logger.Info("sending new nft minted message to discord", zap.Any("message", zap.Any("discordMsg)", discordMsg)))

	channels := []entity.DiscordNotiType{
		entity.NEW_MINT,
	}

	if project.TokenID == PerceptronProjectID {
		channels = append(channels, entity.DISCORD_PERCEPTRON_CHANNEL)
		// wait for perception capture image before send noti to discord
		time.Sleep(1 * time.Minute)
	}

	if value, err := strconv.Atoi(mintPrice); err == nil && value > 0 {
		if category == PFPsCategory {
			channels = append(channels, entity.DISCORD_PFPS_CHANNEL)
		} else {
			channels = append(channels, entity.DISCORD_ART_CHANNEL)
		}
	}

	for _, c := range channels {
		noti := entity.DiscordNoti{
			Message:    discordMsg,
			NumRetried: 0,
			Status:     entity.PENDING,
			Type:       c,
			Meta: entity.DiscordNotiMeta{
				ProjectID:     project.TokenID,
				InscriptionID: tokenUri.TokenID,
			},
		}
		err = u.CreateDiscordNoti(noti)
	}
	// create discord message
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNFTMinted.CreateDiscordNoti", zap.Error(err))
	}
}

func (u Usecase) NotifyCreateNewProjectToDiscord(project *entity.Projects, owner *entity.Users, proposed bool, proposalID string) {
	domain := os.Getenv("DOMAIN")

	var category, description string
	if len(project.Categories) > 0 {
		// we assume that there are only one category
		categoryEntity, err := u.GetCategory(project.Categories[0])
		if err != nil {
			logger.AtLog.Logger.Error("NotifyCreateNewProjectToDiscord.GetCategory failed", zap.Any("err", err))
			return
		}
		category = categoryEntity.Name
		description = fmt.Sprintf("Category: %s\n", category)
	}
	address := owner.WalletAddressBTC
	if address == "" {
		address = owner.WalletAddress
	}
	ownerName := u.resolveShortName(owner.DisplayName, address)
	collectionName := project.Name

	fields := make([]entity.Field, 0)
	addFields := func(fields []entity.Field, name string, value string, inline bool) []entity.Field {
		if value == "" {
			return fields
		}
		return append(fields, entity.Field{
			Name:   name,
			Value:  value,
			Inline: inline,
		})
	}
	fields = addFields(fields, "", u.resolveShortDescription(project.Description), false)
	fields = addFields(fields, "Mint Price", u.resolveMintPriceBTC(project.MintPrice), true)
	fields = addFields(fields, "Max Supply", fmt.Sprintf("%d", project.MaxSupply), true)

	parsedThumbnailUrl, err := url.Parse(project.Thumbnail)
	if err != nil {
		logger.AtLog.Logger.Error("ErrorParseProjectThumbnailURL", zap.Error(err))
	}
	parsedThumbnail := parsedThumbnailUrl.String()

	var content string
	if proposed {
		content = "**NEW DROP PROPOSED ✋**"
	} else {
		content = "**NEW DROP APPROVED ✅**"
	}

	var url string
	if proposed {
		url = fmt.Sprintf("%s/dao?tab=0&id=%s", domain, proposalID)
	} else {
		url = fmt.Sprintf("%s/generative/%s", domain, project.GenNFTAddr)
	}

	discordMsg := entity.DiscordMessage{
		Username: "Satoshi 27",
		Content:  content,
		Embeds: []entity.Embed{{
			Title:       fmt.Sprintf("%s\n***%s***", ownerName, collectionName),
			Url:         url,
			Description: description,
			//Author: discordclient.Author{
			//	Name:    u.resolveShortName(owner.DisplayName, owner.WalletAddress),
			//	Url:     fmt.Sprintf("%s/profile/%s", domain, owner.WalletAddress),
			//	IconUrl: owner.Avatar,
			//},
			Fields: fields,
			Image: entity.Image{
				Url: parsedThumbnail,
			},
		}},
	}
	logger.AtLog.Logger.Info("sending new create new project message to discord", zap.Any("message", zap.Any("discordMsg)", discordMsg)))

	var msgType entity.DiscordNotiType
	if proposed {
		msgType = entity.NEW_PROJECT_PROPOSED
	} else {
		msgType = entity.NEW_PROJECT_APPROVED
	}

	noti := entity.DiscordNoti{
		Message:    discordMsg,
		NumRetried: 0,
		Status:     entity.PENDING,
		Type:       msgType,
		Meta: entity.DiscordNotiMeta{
			ProjectID: project.TokenID,
		},
	}

	// create discord message
	err = u.CreateDiscordNoti(noti)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyCreateNewProjectToDiscord.CreateDiscordNoti", zap.Error(err))
	}
}

func (u Usecase) NotifyNewBid(ETHWalletAddress string, unitPrice float64, quantity int, collectorRedictTo string) error {
	logger.AtLog.Logger.Info(
		"NotifyNewBid",
		zap.Any("price", unitPrice),
		zap.Any("quantity", quantity),
		zap.Any("ETHWalletAddress", ETHWalletAddress),
	)

	domain := os.Getenv("DOMAIN")

	bidder, err := u.Repo.FindUserByWalletAddress(ETHWalletAddress)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNewBid.FindUserByBtcAddress", zap.Error(err))
		return err
	}

	fields := make([]entity.Field, 0)
	addFields := func(fields []entity.Field, name string, value string, inline bool) []entity.Field {
		if value == "" {
			return fields
		}
		return append(fields, entity.Field{
			Name:   name,
			Value:  value,
			Inline: inline,
		})
	}
	fields = addFields(fields, "", "Category: AI", false)

	bidderName := bidder.DisplayName
	if bidderName == "" {
		bidderName = bidder.WalletAddress[:4] + "..." + bidder.WalletAddress[len(bidder.WalletAddress)-4:]
	}

	CollectorUrl := ""
	switch collectorRedictTo {
	case "opensea":
		CollectorUrl = "https://opensea.io/" + bidder.WalletAddress
		if bidderName == "" {
			bidderName = bidder.WalletAddress[:4] + "..." + bidder.WalletAddress[len(bidder.WalletAddress)-4:]
		}
	default:
		CollectorUrl = domain + "/profile/" + bidder.WalletAddressBTCTaproot
		if bidderName == "" {
			bidderName = bidder.WalletAddressBTCTaproot[:4] + "..." + bidder.WalletAddressBTCTaproot[len(bidder.WalletAddressBTCTaproot)-4:]
		}
	}

	fields = addFields(fields, "Collector", fmt.Sprintf("[%s](%s)", bidderName, CollectorUrl), true)
	fields = addFields(fields, "Bid Amount", fmt.Sprintf("%.3f ETH", unitPrice), true)
	fields = addFields(fields, "Quantity", fmt.Sprintf("%d", quantity), true)
	fields = addFields(fields, "", "Perceptrons is an experimental collection of on-chain AI models. While many projects have stored outputs from AI models on-chain, Perceptrons attempts to store the actual AI models themselves, allowing users to query the artwork and run live image recognition tasks.", false)

	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW BID**",
		Embeds: []entity.Embed{{
			Title:  fmt.Sprintf("Generative\n***AI Series: Perceptrons***"),
			Url:    fmt.Sprintf("%v/ai", domain),
			Fields: fields,
			Image: entity.Image{
				Url: "https://storage.googleapis.com/generative-static-prod/btc-projects/perceptrons.gif",
			},
		}},
	}
	noti := entity.DiscordNoti{
		Message:    discordMsg,
		NumRetried: 0,
		Status:     entity.PENDING,
		Type:       entity.NEW_BID,
	}

	// create discord message
	err = u.CreateDiscordNoti(noti)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNFTMinted.CreateDiscordNoti", zap.Error(err))
	}
	return nil
}

func (u Usecase) NewReportNoti(project *entity.Projects, reportLink, walletAddress string) error {

	logger.AtLog.Logger.Info("NewReportNoti", zap.Any("projectID", project.TokenID), zap.Any("walletAddress", walletAddress))
	domain := os.Getenv("DOMAIN")

	reporter, err := u.Repo.FindUserByWalletAddress(walletAddress)
	if err != nil {
		logger.AtLog.Logger.Error("NewReportNoti.FindUserByWalletAddress", zap.Error(err))
		return err
	}

	owner, err := u.Repo.FindUserByWalletAddress(project.CreatorAddrr)
	if err != nil {
		logger.AtLog.Logger.Error("NewReportNoti.FindUserByWalletAddress", zap.Error(err))
		return err
	}

	reporterName := reporter.GetDisplayNameByTapRootAddress()

	fields := make([]entity.Field, 0)
	addFields := func(fields []entity.Field, name string, value string, inline bool) []entity.Field {
		if value == "" {
			return fields
		}
		return append(fields, entity.Field{
			Name:   name,
			Value:  value,
			Inline: inline,
		})
	}
	if len(project.Categories) > 0 {
		category, _ := u.Repo.FindCategory(project.Categories[0])
		if category != nil {
			fields = addFields(fields, "", fmt.Sprintf("Category: %s", category.Name), false)
		}
	}

	fields = addFields(fields, "Reporter", fmt.Sprintf("[%s](%s)", reporterName, domain+"/profile/"+reporter.WalletAddressBTCTaproot), false)
	fields = addFields(fields, "Original Work", fmt.Sprintf("[%s](%s)", reportLink, reportLink), false)
	fields = addFields(fields, "", project.Description, false)

	parsedThumbnailUrl, err := url.Parse(project.Thumbnail)
	if err != nil {
		logger.AtLog.Logger.Error("ErrorParseProjectThumbnailURL", zap.Error(err))
	}
	parsedThumbnail := parsedThumbnailUrl.String()
	ownerName := owner.GetDisplayNameByTapRootAddress()
	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   fmt.Sprintf("**:sos: NEW REPORT: COPYMINT :sos:**"),
		Embeds: []entity.Embed{{
			Title:  fmt.Sprintf("%v\n***%s***", ownerName, project.Name),
			Url:    fmt.Sprintf("%v/generative/%s", domain, project.TokenID),
			Fields: fields,
			Image: entity.Image{
				Url: parsedThumbnail,
			},
		}},
	}

	noti := entity.DiscordNoti{
		Message:    discordMsg,
		NumRetried: 0,
		Status:     entity.PENDING,
		Type:       entity.DISCORD_DAO_CHANNEL,
		Meta: entity.DiscordNotiMeta{
			ProjectID: project.TokenID,
		},
	}

	// create discord message
	err = u.CreateDiscordNoti(noti)
	if err != nil {
		logger.AtLog.Logger.Error("NewReportNoti.CreateDiscordNoti", zap.Error(err))
	}
	return nil
}

const MAX_SEND_DISCORD_RETRY_TIMES = 3

func (u Usecase) JobSendDiscordNoti() error {
	logger.AtLog.Logger.Info("JobSendDiscordNoti.Start")
	for page := int64(1); ; page++ {
		status := entity.PENDING
		logger.AtLog.Logger.Info("JobSendDiscordNoti.StartGetPendingMessages", zap.Any("page", zap.Any("page)", page)))
		resp, err := u.Repo.GetDiscordNoties(entity.GetDiscordNotiReq{
			Page:   page,
			Limit:  10,
			Status: &status,
		})
		if err != nil {
			logger.AtLog.Logger.Error("JobSendDiscordNoti.ErrorWhenGetPendingNoties", zap.Any("page", page))
			return errors.WithStack(err)
		}
		uNoties := resp.Result
		noties := uNoties.([]entity.DiscordNoti)
		logger.AtLog.Logger.Info("JobSendDiscordNoti.DoneGetPendingMessages", zap.Any("page", page), zap.Any("lenNoties", zap.Any("len(noties))", len(noties))))
		if len(noties) == 0 {
			break
		}

		processed := 0

		for _, noti := range noties {
			processed += 1
			discordMsg := &discordclient.Message{}
			copier.Copy(discordMsg, noti.Message)
			logger.AtLog.Logger.Info("sending new airdrop message to discord", zap.Any("discordMsg", discordMsg))

			if err := u.DiscordClient.SendMessage(context.TODO(), noti.Webhook, *discordMsg); err != nil {
				logger.AtLog.Logger.Error("JobSendDiscordNoti.errorSendingMessageToDiscord", zap.Error(err))
				u.Repo.UpdateDiscordNotiAddRetry(noti.UUID)
				if noti.NumRetried+1 == MAX_SEND_DISCORD_RETRY_TIMES {
					u.Repo.UpdateDiscordStatus(noti.UUID, entity.FAILED)
				}
			} else {
				u.Repo.UpdateDiscordStatus(noti.UUID, entity.DONE)
			}
			if processed%5 == 0 {
				time.Sleep(1 * time.Second)
			}
		}
	}

	return nil
}

func (u Usecase) CreateDiscordNoti(noti entity.DiscordNoti) error {
	partners, err := u.Repo.GetAllDiscordPartner()
	if err != nil {
		return errors.WithStack(err)
	}
	for _, partner := range partners {
		webhook := partner.Webhooks[string(noti.Type)]
		if webhook == "" {
			continue
		}
		toCreate := (len(partner.ProjectIDs) == 0)
		for _, projectID := range partner.ProjectIDs {
			if projectID == noti.Meta.ProjectID {
				toCreate = true
			}
		}

		if toCreate {
			tmpNoti := &entity.DiscordNoti{}
			copier.Copy(tmpNoti, noti)
			tmpNoti.Webhook = webhook
			tmpNoti.Meta.SentTo = partner.Name
			logger.AtLog.Logger.Info("CreateDiscordNoti.SendToPartner", zap.Any("tmpNoti", zap.Any("tmpNoti)", tmpNoti)))
			u.Repo.CreateDiscordNoti(*tmpNoti)
		}
	}

	return nil
}
