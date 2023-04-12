package usecase

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/url"
	"os"
	"rederinghub.io/utils/constants/dao_project_voted"
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
	MaxSendDiscordRetryTimes = 3
	PerceptronProjectID      = "1002573"
	PFPsCategory             = "PFPs"
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

func (u Usecase) NotifyNewSale(order entity.DexBTCListing) error {

	domain := os.Getenv("DOMAIN")
	tokenUri, err := u.Repo.FindTokenByTokenID(order.InscriptionID)
	if err != nil {
		return err
	}

	project, err := u.Repo.FindProjectByTokenID(tokenUri.ProjectID)
	if err != nil {
		return err
	}

	var category string
	if len(project.Categories) > 0 {
		categoryEntity, _ := u.GetCategory(project.Categories[0])
		if categoryEntity != nil {
			category = categoryEntity.Name
		}
	}

	owner, err := u.Repo.FindUserByAddress(project.CreatorProfile.WalletAddress)
	if err != nil {
		return err
	}
	ownerName := owner.GetDisplayNameByWalletAddress()

	fields := make([]entity.Field, 0)
	fields = addDiscordField(fields, "Sale Price", u.resolveMintPriceBTC(fmt.Sprintf("%v", order.Amount)), true)
	fields = u.addUserDiscordField(addUserDiscordFieldReq{
		Fields:  fields,
		Key:     "Buyer",
		Address: order.Buyer,
		Inline:  true,
		Domain:  domain,
	})
	fields = u.addUserDiscordField(addUserDiscordFieldReq{
		Fields:  fields,
		Key:     "Seller",
		Address: order.SellerAddress,
		Inline:  true,
		Domain:  domain,
	})

	parsedThumbnail := ""
	parsedThumbnailUrl, _ := url.Parse(tokenUri.Thumbnail)

	if parsedThumbnailUrl != nil {
		parsedThumbnail = parsedThumbnailUrl.String()
	}

	embed := entity.Embed{
		Url:    fmt.Sprintf("%s/generative/%s/%s", domain, project.GenNFTAddr, tokenUri.TokenID),
		Fields: fields,
	}

	if order.Amount == 0 {
		embed.Title = fmt.Sprintf("%s\n***%s #%d***", ownerName, project.Name, tokenUri.OrderInscriptionIndex)
		embed.Thumbnail = entity.Thumbnail{
			Url: parsedThumbnail,
		}
	} else {
		embed.Title = fmt.Sprintf("%s\n***%s #%d***", ownerName, project.Name, tokenUri.OrderInscriptionIndex)
		embed.Image = entity.Image{
			Url: parsedThumbnail,
		}
	}

	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW SALE**",
		Embeds:    []entity.Embed{embed},
	}

	logger.AtLog.Logger.Info("sending new sale message to discord", zap.Any("message", zap.Any("discordMsg)", discordMsg)))
	types := []entity.DiscordNotiType{entity.NEW_SALE}
	if order.Amount > 0 {
		if tokenUri.ProjectID == PerceptronProjectID {
			types = append(types, entity.NEW_SALE_PERCEPTRON)
		} else if category == PFPsCategory {
			types = append(types, entity.NEW_SALE_PFPS)
		} else {
			types = append(types, entity.NEW_SALE_ART)
		}
	}
	for _, t := range types {
		noti := entity.DiscordNoti{
			Message:    discordMsg,
			NumRetried: 0,
			Status:     entity.PENDING,
			Type:       t,
			Meta: entity.DiscordNotiMeta{
				ProjectID:     project.TokenID,
				InscriptionID: tokenUri.TokenID,
				Amount:        order.Amount,
				Category:      category,
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

	domain := os.Getenv("DOMAIN")
	tokenUri, err := u.Repo.FindTokenByTokenID(order.InscriptionID)
	if err != nil {
		return err
	}

	project, err := u.GetProjectByGenNFTAddr(tokenUri.ProjectID)
	if err != nil {
		return err
	}

	var category, description string
	if len(project.Categories) > 0 {
		// we assume that there are only one category
		categoryEntity, _ := u.GetCategory(project.Categories[0])
		if categoryEntity != nil {
			category = categoryEntity.Name
			description = fmt.Sprintf("Category: %s\n", category)
		}
	}

	ownerName := u.resolveShortName(project.CreatorProfile.DisplayName, project.CreatorProfile.WalletAddress)
	collectionName := project.Name
	mintedCount := tokenUri.OrderInscriptionIndex

	fields := make([]entity.Field, 0)
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

	notify := entity.DiscordNoti{
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
	err = u.CreateDiscordNoti(notify)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) NotifyNFTMinted(inscriptionID string) error {
	time.Sleep(3 * time.Minute)
	domain := os.Getenv("DOMAIN")
	tokenUri, err := u.Repo.FindTokenByTokenID(inscriptionID)
	if err != nil {
		return err
	}
	project, err := u.Repo.FindProjectByTokenID(tokenUri.ProjectID)
	if err != nil {
		return err
	}

	owner, err := u.Repo.FindUserByWalletAddress(project.CreatorProfile.WalletAddress)
	if err != nil {
		return err
	}

	var category string
	if len(project.Categories) > 0 {
		categoryEntity, _ := u.GetCategory(project.Categories[0])
		if categoryEntity != nil {
			category = categoryEntity.Name
		}
	}

	fields := make([]entity.Field, 0)
	mintPrice := project.MintPrice
	mintNftBtc, _ := u.Repo.FindMintNftBtcByInscriptionID(inscriptionID)
	if mintNftBtc != nil {
		if v, ok := mintNftBtc.EstFeeInfo["btc"]; ok {
			mintPrice = v.MintPrice
		}
	}

	fields = addDiscordField(fields, "Mint Price", u.resolveMintPriceBTC(mintPrice), true)
	mintPriceInNum, _ := strconv.Atoi(mintPrice)

	fields = u.addUserDiscordField(addUserDiscordFieldReq{
		Fields:  fields,
		Key:     "Collector",
		Address: tokenUri.OwnerAddr,
		Inline:  true,
		Domain:  domain,
	})

	parsedThumbnail := ""
	parsedThumbnailUrl, _ := url.Parse(tokenUri.Thumbnail)
	if parsedThumbnailUrl != nil {
		parsedThumbnail = parsedThumbnailUrl.String()
	}

	embed := entity.Embed{
		Url:    fmt.Sprintf("%s/generative/%s/%s", domain, project.GenNFTAddr, tokenUri.TokenID),
		Fields: fields,
	}

	if mintPriceInNum == 0 {
		embed.Title = fmt.Sprintf("%s\n***%v #%v***", owner.GetDisplayNameByWalletAddress(), project.Name, tokenUri.OrderInscriptionIndex)
		embed.Thumbnail = entity.Thumbnail{
			Url: parsedThumbnail,
		}
	} else {
		embed.Title = fmt.Sprintf("%s\n***%v #%v***", owner.GetDisplayNameByWalletAddress(), project.Name, tokenUri.OrderInscriptionIndex)
		embed.Image = entity.Image{
			Url: parsedThumbnail,
		}
	}

	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW MINT**",
		Embeds:    []entity.Embed{embed},
	}

	types := []entity.DiscordNotiType{entity.NEW_MINT}
	if mintPriceInNum > 0 {
		if tokenUri.ProjectID == PerceptronProjectID {
			types = append(types, entity.NEW_MINT_PERCEPTRON)
		} else if category == PFPsCategory {
			types = append(types, entity.NEW_MINT_PFPS)
		} else {
			types = append(types, entity.NEW_MINT_ART)
		}
	}

	for _, t := range types {
		notify := entity.DiscordNoti{
			Message:    discordMsg,
			NumRetried: 0,
			Status:     entity.PENDING,
			Type:       t,
			Meta: entity.DiscordNotiMeta{
				ProjectID:     project.TokenID,
				InscriptionID: tokenUri.TokenID,
				Amount:        uint64(mintPriceInNum),
				Category:      category,
			},
		}
		err = u.CreateDiscordNoti(notify)
	}

	return err
}

func (u Usecase) NotifyNewProject(project *entity.Projects, owner *entity.Users, proposed bool, proposalID string) {

	domain := os.Getenv("DOMAIN")

	var category string
	collectionName := project.Name

	thumbnail := ""
	parsedThumbnailUrl, _ := url.Parse(project.Thumbnail)
	if parsedThumbnailUrl != nil {
		thumbnail = parsedThumbnailUrl.String()
	}

	fields := make([]entity.Field, 0)
	var msgType entity.DiscordNotiType
	if len(project.Categories) > 0 {
		categoryEntity, _ := u.GetCategory(project.Categories[0])
		if categoryEntity != nil {
			category = categoryEntity.Name
		}
	}

	discordMsg := entity.DiscordMessage{
		Username: "Satoshi 27",
	}
	embed := entity.Embed{
		Title: fmt.Sprintf("%s\n***%s***", owner.GetDisplayNameByWalletAddress(), collectionName),
	}

	if proposed {
		embed.Url = fmt.Sprintf("%s/dao?tab=0&id=%s", domain, proposalID)
		msgType = entity.NEW_PROJECT_PROPOSED
		discordMsg.Content = "**NEW DROP PROPOSED ✋**"
		fields = addDiscordField(fields, "Category", category, false)
		fields = addDiscordField(fields, "", u.resolveShortDescription(project.Description), false)
		embed.Image = entity.Image{
			Url: thumbnail,
		}
	} else {
		embed.Url = fmt.Sprintf("%s/generative/%s", domain, project.GenNFTAddr)
		msgType = entity.NEW_PROJECT_APPROVED
		discordMsg.Content = "**NEW DROP APPROVED ✅**"
		embed.Thumbnail = entity.Thumbnail{
			Url: thumbnail,
		}
	}

	fields = addDiscordField(fields, "Mint Price", u.resolveMintPriceBTC(project.MintPrice), true)
	fields = addDiscordField(fields, "Max Supply", fmt.Sprintf("%d", project.MaxSupply), true)
	embed.Fields = fields
	discordMsg.Embeds = []entity.Embed{embed}

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
	err := u.CreateDiscordNoti(noti)
	if err != nil {
		logger.AtLog.Logger.Error("NotifyNewProject.CreateDiscordNoti", zap.Error(err))
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

func (u Usecase) NotifiNewProjectReport(project *entity.Projects, reportLink, reporterWalletAddress string) error {

	domain := os.Getenv("DOMAIN")
	reporter, err := u.Repo.FindUserByWalletAddress(reporterWalletAddress)
	if err != nil {
		return err
	}
	owner, err := u.Repo.FindUserByWalletAddress(project.CreatorAddrr)
	if err != nil {
		return err
	}

	reporterName := reporter.GetDisplayNameByTapRootAddress()
	catName := ""
	parsedThumbnail := ""
	ownerName := owner.GetDisplayNameByTapRootAddress()

	if len(project.Categories) > 0 {
		category, _ := u.Repo.FindCategory(project.Categories[0])
		if category != nil {
			catName = category.Name
		}
	}

	parsedThumbnailUrl, _ := url.Parse(project.Thumbnail)
	if parsedThumbnailUrl != nil {
		parsedThumbnail = parsedThumbnailUrl.String()
	}

	fields := make([]entity.Field, 0)
	fields = addDiscordField(fields, "", u.resolveShortDescription(project.Description), false)
	fields = addDiscordField(fields, "Reporter", fmt.Sprintf("[%s](%s)", reporterName, domain+"/profile/"+reporter.WalletAddressBTCTaproot), false)
	fields = addDiscordField(fields, "Original Work", fmt.Sprintf("[%s](%s)", reportLink, reportLink), false)
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
		Type:       entity.NEW_PROJECT_REPORT,
		Meta: entity.DiscordNotiMeta{
			ProjectID: project.TokenID,
			Category:  catName,
		},
	}

	// create discord message
	err = u.CreateDiscordNoti(noti)
	if err != nil {
		logger.AtLog.Logger.Error("NotifiNewProjectReport.CreateDiscordNoti", zap.Error(err))
	}
	return nil
}

func (u Usecase) NotifyNewProjectVote(daoProject *entity.DaoProject, vote *entity.DaoProjectVoted) error {
	project := &entity.Projects{}
	if err := u.Repo.FindOneBy(context.TODO(), project.TableName(), bson.M{"_id": daoProject.ProjectId}, project); err != nil {
		return err
	}
	voter, err := u.Repo.FindUserByWalletAddress(vote.CreatedBy)
	if err != nil {
		return err
	}

	owner, err := u.Repo.FindUserByWalletAddress(project.CreatorAddrr)

	if err != nil {
		return err
	}

	domain := os.Getenv("DOMAIN")
	countVote := u.Repo.CountDAOProjectVoteByStatus(context.TODO(), daoProject.ID, dao_project_voted.Voted)
	totalVote := u.Config.CountVoteDAO
	if totalVote <= 0 {
		totalVote = 2
	}
	thumbnail := ""

	parsedThumbnailUrl, _ := url.Parse(project.Thumbnail)
	if parsedThumbnailUrl != nil {
		thumbnail = parsedThumbnailUrl.String()
	}

	fields := make([]entity.Field, 0)
	fields = addDiscordField(fields, "Voter", fmt.Sprintf("[%s](%s)", voter.GetDisplayNameByWalletAddress(), domain+"/profile/"+voter.WalletAddress), true)
	fields = addDiscordField(fields, "Voted", fmt.Sprintf("%d/%d", countVote, totalVote), true)

	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   fmt.Sprintf("**NEW VOTE :thumbsup:**"),
		Embeds: []entity.Embed{{
			Title:  fmt.Sprintf("%v\n***%s***", owner.GetDisplayNameByWalletAddress(), project.Name),
			Url:    fmt.Sprintf("%v/generative/%s", domain, project.TokenID),
			Fields: fields,
			Thumbnail: entity.Thumbnail{
				Url: thumbnail,
			},
		}},
	}

	noti := entity.DiscordNoti{
		Message:    discordMsg,
		NumRetried: 0,
		Status:     entity.PENDING,
		Type:       entity.NEW_PROJECT_VOTE,
		Meta: entity.DiscordNotiMeta{
			ProjectID: project.TokenID,
		},
	}

	// create discord message
	err = u.CreateDiscordNoti(noti)
	if err != nil {
		logger.AtLog.Logger.Error("NotifiNewProjectReport.CreateDiscordNoti", zap.Error(err))
	}
	return nil
}

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
				if noti.NumRetried+1 == MaxSendDiscordRetryTimes {
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
		if partner.MatchProject(noti.Meta.ProjectID) && partner.MatchCategory(noti.Meta.Category) && partner.MatchAmountGreaterThanZero(noti.Meta.Amount) {
			tmpNoti := &entity.DiscordNoti{}
			copier.Copy(tmpNoti, noti)
			tmpNoti.Webhook = webhook
			tmpNoti.Meta.SentTo = partner.Name
			logger.AtLog.Logger.Info("Create Discord Notifications", zap.Any("event", noti.Type), zap.Any("partner", partner.Name))
			u.Repo.CreateDiscordNoti(*tmpNoti)
		}
	}

	return nil
}

func (u Usecase) TestSendNoti() {
	domain := os.Getenv("DOMAIN")
	if domain == "https://devnet.generative.xyz" {
		//project, _ := u.Repo.FindProjectByTokenID("1001001")
		incriptionID := "48ace75807885f42ade68e2bc532ec8affd7870c58c4529108e05435e9ea40aai0"

		//user, _ := u.Repo.FindUserByWalletAddress(project.CreatorAddrr)
		//daoProject := &entity.DaoProject{}
		//u.Repo.FindOneBy(context.TODO(), daoProject.TableName(), bson.M{"_id": "642b921af46c66bdf68c2d82"}, daoProject)
		//vote := &entity.DaoProjectVoted{
		//	CreatedBy:    project.CreatorAddrr,
		//	DaoProjectId: daoProject.ProjectId,
		//	Status:       1,
		//}

		//u.NotifiNewProjectReport(project, domain, project.CreatorAddrr)
		//u.Notifynewsale(entity.DexBTCListing{
		//	SellerAddress: project.ContractAddress,
		//	Buyer:         project.ContractAddress,
		//	Amount:        100000000,
		//	InscriptionID: incriptionID,
		//})
		//u.NotifyNewSale(entity.DexBTCListing{
		//	SellerAddress: project.ContractAddress,
		//	Buyer:         project.ContractAddress,
		//	Amount:        0,
		//	InscriptionID: incriptionID,
		//})
		u.NotifyNFTMinted(incriptionID)
		//u.NotifyNewProject(project, user, true, "proposalID")
		//u.NotifyNewProject(project, user, false, "proposalID")
		//u.NotifyNewProjectVote(daoProject, vote)
		//u.NotifyNewListing(entity.DexBTCListing{
		//	SellerAddress: project.ContractAddress,
		//	Buyer:         project.ContractAddress,
		//	Amount:        100000000,
		//	InscriptionID: incriptionID,
		//})
		u.JobSendDiscordNoti()
		fmt.Println("done")
	}
}
