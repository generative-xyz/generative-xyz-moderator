package usecase

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	discordclient "rederinghub.io/utils/discord"
)

type addUserDiscordFieldReq struct {
	Fields []entity.Field
	Key string
	Address string
	UserID string
	Inline bool
	Domain string
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
	} else if req.UserID != ""  {
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
		u.Logger.ErrorAny("NotifyNewSale.FindUserByAddress")
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
	webhook := os.Getenv("DISCORD_AIRDROP_WEBHOOK")
	fields := make([]entity.Field, 0)
	file := strings.Replace(airdrop.File, "html", "png", 1)

	fields = u.addUserDiscordField(addUserDiscordFieldReq{
		Fields: fields,
		Key: "Key holder",
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
		u.Logger.Error("ErrorWhenGetInscribeInfo", zap.Any("inscriptionId", airdrop.InscriptionId))
	}
	
	parsedThumbnailUrl, err := url.Parse(file)
	if err != nil {
		u.Logger.ErrorAny("ErrorParseProjectThumbnailURL", zap.Error(err))
	}
	parsedThumbnail := parsedThumbnailUrl.String()


	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW KEY**",
		Embeds: []entity.Embed{{
			Title: fmt.Sprintf("%s%s", title, inscriptionNumTitle),
			Url: fmt.Sprintf("https://generativeexplorer.com/inscription/%s", airdrop.InscriptionId),
			Fields: fields,
			Image: entity.Image{
				Url: parsedThumbnail,
			},
		}},
	}

	u.Logger.Info("sending new airdrop message to discord", zap.Any("message", discordMsg))

	// create discord message
	err = u.Repo.CreateDiscordNoti(discordMsg, webhook)
	if err != nil {
		u.Logger.ErrorAny("NotifyNewAirdrop.CreateDiscordNoti", zap.Error(err))
		return err
	}
	return nil
}

func (u Usecase) NotifyNewSale(order entity.DexBTCListing, buyerAddress string) error {
	u.Logger.Info("NotifyNewSale.Start", zap.Any("order", order), zap.Any("buyerAddress", buyerAddress))
	domain := os.Getenv("DOMAIN")
	webhook := os.Getenv("DISCORD_NEW_SALE_WEBHOOK")

	tokenUri, err := u.Repo.FindTokenByTokenID(order.InscriptionID)
	if err != nil {
		u.Logger.ErrorAny("NotifyNFTMinted.FindTokenByTokenID failed", zap.Any("err", err.Error()))
		return err
	}

	project, err := u.GetProjectByGenNFTAddr(tokenUri.ProjectID)
	if err != nil {
		u.Logger.ErrorAny("NotifyNFTMinted.GetProjectByGenNFTAddr failed", zap.Any("err", err))
		return err
	}

	var category, description string
	if len(project.Categories) > 0 {
		// we assume that there are only one category
		categoryEntity, err := u.GetCategory(project.Categories[0])
		if err != nil {
			u.Logger.ErrorAny("NotifyNFTMinted.GetCategory failed", zap.Any("err", err))
			return err
		}
		category = categoryEntity.Name
		description = fmt.Sprintf("Category: %s\n", category)
	}

	ownerName := u.resolveShortName(tokenUri.Creator.DisplayName, tokenUri.Creator.WalletAddress)
	collectionName := project.Name
	mintedCount := tokenUri.OrderInscriptionIndex

	fields := make([]entity.Field, 0)

	fields = addDiscordField(fields, "", u.resolveShortDescription(project.Description), false)

	fields = addDiscordField(fields, "Sale Price", u.resolveMintPriceBTC(fmt.Sprintf("%v", order.Amount)), true)

	if buyerAddress != "" {
		fields = u.addUserDiscordField(addUserDiscordFieldReq{
			Fields: fields,
			Key: "Buyer",
			Address: buyerAddress,
			Inline: true,
			Domain: domain,
		})
	}

	fields = u.addUserDiscordField(addUserDiscordFieldReq{
		Fields: fields,
		Key: "Seller",
		Address: order.SellerAddress,
		Inline: true,
		Domain: domain,
	})

	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW SALE**",
		Embeds: []entity.Embed{{
			Title: fmt.Sprintf("%s\n***%s #%d***", ownerName, collectionName, mintedCount),
			Url: fmt.Sprintf("%s/generative/%s/%s", domain, project.GenNFTAddr, tokenUri.TokenID),
			Description: description,
			Fields: fields,
			Image: entity.Image{
				Url: tokenUri.Thumbnail,
			},
		}},
	}

	u.Logger.Info("sending new sale message to discord", zap.Any("message", discordMsg))

	// create discord message
	err = u.Repo.CreateDiscordNoti(discordMsg, webhook)
	if err != nil {
		u.Logger.ErrorAny("NotifyNewSale.CreateDiscordNoti", zap.Error(err))
		return err
	}
	return nil
}


func (u Usecase) NotifyNewListing(order entity.DexBTCListing) error {
	u.Logger.Info("NotifyNewListing.Start", zap.Any("order", order))
	domain := os.Getenv("DOMAIN")
	webhook := os.Getenv("DISCORD_NEW_LISTING_WEBHOOK")

	tokenUri, err := u.Repo.FindTokenByTokenID(order.InscriptionID)
	if err != nil {
		u.Logger.ErrorAny("NotifyNFTMinted.FindTokenByTokenID failed", zap.Any("err", err.Error()))
		return err
	}

	project, err := u.GetProjectByGenNFTAddr(tokenUri.ProjectID)
	if err != nil {
		u.Logger.ErrorAny("NotifyNFTMinted.GetProjectByGenNFTAddr failed", zap.Any("err", err))
		return err
	}

	var category, description string
	if len(project.Categories) > 0 {
		// we assume that there are only one category
		categoryEntity, err := u.GetCategory(project.Categories[0])
		if err != nil {
			u.Logger.ErrorAny("NotifyNFTMinted.GetCategory failed", zap.Any("err", err))
			return err
		}
		category = categoryEntity.Name
		description = fmt.Sprintf("Category: %s\n", category)
	}

	ownerName := u.resolveShortName(tokenUri.Creator.DisplayName, tokenUri.Creator.WalletAddress)
	collectionName := project.Name
	mintedCount := tokenUri.OrderInscriptionIndex

	fields := make([]entity.Field, 0)

	fields = addDiscordField(fields, "", u.resolveShortDescription(project.Description), false)

	fields = addDiscordField(fields, "List Price", u.resolveMintPriceBTC(fmt.Sprintf("%v", order.Amount)), true)

	fields = u.addUserDiscordField(addUserDiscordFieldReq{
		Fields: fields,
		Key: "Seller",
		Address: order.SellerAddress,
		Inline: true,
		Domain: domain,
	})

	parsedThumbnailUrl, err := url.Parse(tokenUri.Thumbnail)
	if err != nil {
		u.Logger.ErrorAny("ErrorParseProjectThumbnailURL", zap.Error(err))
	}
	parsedThumbnail := parsedThumbnailUrl.String()

	discordMsg := entity.DiscordMessage{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW LISTING**",
		Embeds: []entity.Embed{{
			Title: fmt.Sprintf("%s\n***%s #%d***", ownerName, collectionName, mintedCount),
			Url: fmt.Sprintf("%s/generative/%s/%s", domain, project.GenNFTAddr, tokenUri.TokenID),
			Description: description,
			Fields: fields,
			Image: entity.Image{
				Url: parsedThumbnail,
			},
		}},
	}

	u.Logger.Info("sending new new listing message to discord", zap.Any("message", discordMsg))

	// create discord message
	err = u.Repo.CreateDiscordNoti(discordMsg, webhook)
	if err != nil {
		u.Logger.ErrorAny("NotifyNewListing.CreateDiscordNoti", zap.Error(err))
		return err
	}
	return nil
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
	fields = addFields(fields, "Collector", fmt.Sprintf("[%s](%s)",
		u.resolveShortName(minterDisplayName, btcUserAddr),
		fmt.Sprintf("%s/profile/%s", domain, minterAddress),
	), true)

	// fields = addFields(fields, "Minted", fmt.Sprintf("%d/%d", mintedCount, itemCount), true)
	//fields = addFields(fields, "Network Fee", strconv.FormatFloat(float64(networkFee)/1e8, 'f', -1, 64)+" BTC")

	parsedThumbnailUrl, err := url.Parse(tokenUri.Thumbnail)
	if err != nil {
		u.Logger.ErrorAny("ErrorParseProjectThumbnailURL", zap.Error(err))
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

	u.Logger.Info("sending new nft minted message to discord", zap.Any("message", discordMsg))

	// create discord message
	err = u.Repo.CreateDiscordNoti(discordMsg, webhook)
	if err != nil {
		u.Logger.ErrorAny("NotifyNFTMinted.CreateDiscordNoti", zap.Error(err))
	}
}

func (u Usecase) NotifyCreateNewProjectToDiscord(project *entity.Projects, owner *entity.Users) {
	domain := os.Getenv("DOMAIN")
	webhook := os.Getenv("DISCORD_NEW_PROJECT_WEBHOOK")

	var category, description string
	if len(project.Categories) > 0 {
		// we assume that there are only one category
		categoryEntity, err := u.GetCategory(project.Categories[0])
		if err != nil {
			u.Logger.ErrorAny("NotifyCreateNewProjectToDiscord.GetCategory failed", zap.Any("err", err))
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
	fields = addFields(fields, "", project.Description, false)
	fields = addFields(fields, "Mint Price", u.resolveMintPriceBTC(project.MintPrice), true)
	fields = addFields(fields, "Max Supply", fmt.Sprintf("%d", project.MaxSupply), true)

	parsedThumbnailUrl, err := url.Parse(project.Thumbnail)
	if err != nil {
		u.Logger.ErrorAny("ErrorParseProjectThumbnailURL", zap.Error(err))
	}
	parsedThumbnail := parsedThumbnailUrl.String()

	discordMsg := entity.DiscordMessage{
		Username: "Satoshi 27",
		Content:  "**NEW DROP**",
		Embeds: []entity.Embed{{
			Title:       fmt.Sprintf("%s\n***%s***", ownerName, collectionName),
			Url:         fmt.Sprintf("%s/generative/%s", domain, project.GenNFTAddr),
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
	u.Logger.Info("sending new create new project message to discord", zap.Any("message", discordMsg))
	// create discord message
	err = u.Repo.CreateDiscordNoti(discordMsg, webhook)
	if err != nil {
		u.Logger.ErrorAny("NotifyCreateNewProjectToDiscord.CreateDiscordNoti", zap.Error(err))
	}
}

const MAX_SEND_DISCORD_RETRY_TIMES = 3

func (u Usecase) JobSendDiscordNoti() error {
	u.Logger.Info("JobSendDiscordNoti.Start")
	for page := int64(1);; page++ {
		status := entity.PENDING
		u.Logger.Info("JobSendDiscordNoti.StartGetPendingMessages", zap.Any("page", page))
		resp, err := u.Repo.GetDiscordNoties(entity.GetDiscordNotiReq{
			Page: page,
			Limit: 10,
			Status: &status,
		})
		if err != nil {
			u.Logger.ErrorAny("JobSendDiscordNoti.ErrorWhenGetPendingNoties", zap.Any("page", page))
			return errors.WithStack(err)
		}
		uNoties := resp.Result
		noties := uNoties.([]entity.DiscordNoti)
		u.Logger.Info("JobSendDiscordNoti.DoneGetPendingMessages", zap.Any("page", page), zap.Any("lenNoties", len(noties)))
		if len(noties) == 0 {
			break
		}

		processed := 0

		for _, noti := range noties {
			processed += 1
			discordMsg := &discordclient.Message{} 
			copier.Copy(discordMsg, noti.Message)
			u.Logger.Info("sending new airdrop message to discord", discordMsg)

			if err := u.DiscordClient.SendMessage(context.TODO(), noti.Webhook, *discordMsg); err != nil {
				u.Logger.ErrorAny("JobSendDiscordNoti.errorSendingMessageToDiscord", zap.Error(err))
				u.Repo.UpdateDiscordNotiAddRetry(noti.UUID)
				if noti.NumRetried + 1 == MAX_SEND_DISCORD_RETRY_TIMES {
					u.Repo.UpdateDiscordStatus(noti.UUID, entity.FAILED)
				}
			} else {
				u.Repo.UpdateDiscordStatus(noti.UUID, entity.DONE)
			}
			if processed % 5 == 0 {
				time.Sleep(1 * time.Second)
			}
		}
	}

	return nil
}
