package usecase

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	discordclient "rederinghub.io/utils/discord"
	"rederinghub.io/utils/logger"
)

func (u Usecase) TestSendDiscordNoti(notify entity.DiscordNoti) error {

	//TODO - for testing purpose
	testChannel := make(map[entity.DiscordNotiType]string)
	testChannel[entity.NEW_SALE_PERCEPTRON] = "https://discord.com/api/webhooks/1095169601198887052/XBOyJeb1j70Xmny1oY0gadnwwHTeEjWmENitWwga-W5TB8Jq4a2xzoECkGOPP3TC3HDp"
	testChannel[entity.NEW_MINT_PERCEPTRON] = "https://discord.com/api/webhooks/1095169601198887052/XBOyJeb1j70Xmny1oY0gadnwwHTeEjWmENitWwga-W5TB8Jq4a2xzoECkGOPP3TC3HDp"
	testChannel[entity.NEW_ART_WEBHOOK] = "https://discord.com/api/webhooks/1095169426741014568/7htHmB18SDGc5fcbCdgUCg5WzFLqwytw5ckkDa0bU97iRUFH0ogekzwwfkskKeFlLgsG"
	testChannel[entity.NEW_MINT] = "https://discord.com/api/webhooks/1095169426741014568/7htHmB18SDGc5fcbCdgUCg5WzFLqwytw5ckkDa0bU97iRUFH0ogekzwwfkskKeFlLgsG"
	testChannel[entity.NEW_SALE] = "https://discord.com/api/webhooks/1095169426741014568/7htHmB18SDGc5fcbCdgUCg5WzFLqwytw5ckkDa0bU97iRUFH0ogekzwwfkskKeFlLgsG"
	testChannel[entity.NEW_PROJECT] = "https://discord.com/api/webhooks/1095169426741014568/7htHmB18SDGc5fcbCdgUCg5WzFLqwytw5ckkDa0bU97iRUFH0ogekzwwfkskKeFlLgsG"
	testChannel[entity.NEW_PROJECT_REPORT] = "https://discord.com/api/webhooks/1095169426741014568/7htHmB18SDGc5fcbCdgUCg5WzFLqwytw5ckkDa0bU97iRUFH0ogekzwwfkskKeFlLgsG"
	testChannel[entity.NEW_PROJECT_APPROVED] = "https://discord.com/api/webhooks/1095169426741014568/7htHmB18SDGc5fcbCdgUCg5WzFLqwytw5ckkDa0bU97iRUFH0ogekzwwfkskKeFlLgsG"

	// send inserted notification
	discordMsg := &discordclient.Message{}
	copier.Copy(discordMsg, notify.Message)
	if notify.RequireImage {
		imageURL, err := u.getImageSource(notify.ImageSourceID, notify.ImageSourceType)
		if err != nil || imageURL == "" {
			// skip in case image not founds
			if notify.NumRetried > MaxFindImageRetryTimes {
				//u.Repo.UpdateDiscordStatus(notify.UUID, entity.FAILED, "failed to get image")
				return err
			}
			//u.Repo.UpdateDiscordNotiAddRetry(notify.UUID)
			return err
		}
		if len(notify.Message.Embeds) > 0 {
			switch notify.ImagePosition {
			case entity.ThumbNailPosition:
				discordMsg.Embeds[0].Thumbnail.Url = imageURL
			case entity.FullImagePosition:
				discordMsg.Embeds[0].Image.Url = imageURL
			}
		}
	}
	logger.AtLog.Logger.Info("sending new airdrop message to discord", zap.Any("discordMsg", discordMsg))

	err := u.DiscordClient.SendMessage(context.TODO(), testChannel[notify.Type], *discordMsg)
	if err != nil {
		return err
	}

	return nil
}
