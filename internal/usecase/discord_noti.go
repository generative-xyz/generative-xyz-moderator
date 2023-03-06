package usecase

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	discordclient "rederinghub.io/utils/discord"
)

type addUserDiscordFieldReq struct {
	Fields []discordclient.Field
	Key string
	Address string
	UserID string
	Inline bool
	Domain string
}

func addDiscordField(fields []discordclient.Field, name string, value string, inline bool) []discordclient.Field {
	if value == "" {
		return fields
	}
	return append(fields, discordclient.Field{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
}

func (u Usecase) addUserDiscordField(req addUserDiscordFieldReq) []discordclient.Field {
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
		userStr = req.Address
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
	fields := make([]discordclient.Field, 0)
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

	discordMsg := discordclient.Message{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW KEY**",
		Embeds: []discordclient.Embed{{
			Title: fmt.Sprintf("%s%s", title, inscriptionNumTitle),
			Url: fmt.Sprintf("https://generativeexplorer.com/inscription/%s", airdrop.InscriptionId),
			Fields: fields,
			Image: discordclient.Image{
				Url: file,
			},
		}},
	}

	sendCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	u.Logger.Info("sending new airdrop message to discord", discordMsg)

	if err := u.DiscordClient.SendMessage(sendCtx, webhook, discordMsg); err != nil {
		u.Logger.Error("NotifyNewAirdrop.errorSendingMessageToDiscord", err)
	}
	return nil
}

func (u Usecase) NotifyNewSale(order entity.DexBTCListing, buyerAddress string) error {
	u.Logger.Info("NotifyNewSale.Start")
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

	fields := make([]discordclient.Field, 0)

	fields = addDiscordField(fields, "", project.Description, false)

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

	discordMsg := discordclient.Message{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW SALE**",
		Embeds: []discordclient.Embed{{
			Title: fmt.Sprintf("%s\n***%s #%d***", ownerName, collectionName, mintedCount),
			Url: fmt.Sprintf("%s/generative/%s/%s", domain, project.GenNFTAddr, tokenUri.TokenID),
			Description: description,
			Fields: fields,
			Image: discordclient.Image{
				Url: tokenUri.Thumbnail,
			},
		}},
	}
	sendCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	u.Logger.Info("sending new sale message to discord", discordMsg)

	if err := u.DiscordClient.SendMessage(sendCtx, webhook, discordMsg); err != nil {
		u.Logger.Error("NotifyNewSale.errorSendingMessageToDiscord", err)
	}
	return nil	
}


// Nguoi ban
// Token inscriptionId
// Price
func (u Usecase) NotifyNewListing(order entity.DexBTCListing) error {
	u.Logger.Info("NotifyNewListing.Start")
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


	fields := make([]discordclient.Field, 0)

	fields = addDiscordField(fields, "", project.Description, false)

	fields = addDiscordField(fields, "List Price", u.resolveMintPriceBTC(fmt.Sprintf("%v", order.Amount)), true)

	fields = u.addUserDiscordField(addUserDiscordFieldReq{
		Fields: fields,
		Key: "Seller",
		Address: order.SellerAddress,
		Inline: true,
		Domain: domain,
	})

	discordMsg := discordclient.Message{
		Username:  "Satoshi 27",
		AvatarUrl: "",
		Content:   "**NEW LISTING**",
		Embeds: []discordclient.Embed{{
			Title: fmt.Sprintf("%s\n***%s #%d***", ownerName, collectionName, mintedCount),
			Url: fmt.Sprintf("%s/generative/%s/%s", domain, project.GenNFTAddr, tokenUri.TokenID),
			Description: description,
			Fields: fields,
			Image: discordclient.Image{
				Url: tokenUri.Thumbnail,
			},
		}},
	}
	sendCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	u.Logger.Info("sending new listing message to discord", discordMsg)

	if err := u.DiscordClient.SendMessage(sendCtx, webhook, discordMsg); err != nil {
		u.Logger.Error("NotifyNewListing.errorSendingMessageToDiscord", err)
	}
	return nil	
}
