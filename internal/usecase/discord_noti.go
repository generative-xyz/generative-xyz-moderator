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

func (u Usecase) NotifyNewAirdrop(airdrop *entity.Airdrop) error {
	domain := os.Getenv("DOMAIN")
	webhook := os.Getenv("DISCORD_AIRDROP_WEBHOOK")
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

	userId := airdrop.Receiver
	file := strings.Replace(airdrop.File, "html", "png", 1)

	user, err := u.Repo.FindUserByID(userId)
	var userStr string
	if err == nil && user != nil {
		address := user.WalletAddressBTC
		if address == "" {
			address = user.WalletAddress
		}
		userStr = fmt.Sprintf("[%s](%s)",
			u.resolveShortName(user.DisplayName, address),
			fmt.Sprintf("%s/profile/%s", domain, address),
		)
	} else {
		u.Logger.Error("NotifyNewAirdrop.findUser", err)
		userStr = userId
	}

	key := "Key holder"
	fields = addFields(fields, key, userStr, false)
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

			//Author: discordclient.Author{
			//	Name:    u.resolveShortName(minter.DisplayName, minter.WalletAddress),
			//	Url:     fmt.Sprintf("%s/profile/%s", domain, minter.WalletAddress),
			//	IconUrl: minter.Avatar,
			//},
			Fields: fields,
			Image: discordclient.Image{
				Url: file,
			},
		}},
	}

	sendCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	u.Logger.Info("sending message to discord", discordMsg)

	if err := u.DiscordClient.SendMessage(sendCtx, webhook, discordMsg); err != nil {
		u.Logger.Error("NotifyNewAirdrop.errorSendingMessageToDiscord", err)
	}
	return nil
}
