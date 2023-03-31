package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

func (u *Usecase) PubSubCreateTokenThumbnail(channelName string, payload interface{}) {
	info, err := u.Queue.CreateSchedule(channelName, payload, "default")
	if err != nil {
		logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("taskInfo", info), zap.Any("payload", payload), zap.Error(err))
		return
	}
	logger.AtLog.Info("PubSubCreateTokenThumbnai", zap.Any("taskInfo", info), zap.Any("payload", payload))
}

func (u *Usecase) ProccessCreateTokenThumbnail(ctx context.Context, t *asynq.Task) error {
	//Move them to 
	payload := t.Payload()
	tokenURI := &structure.TokenImagePayload{}
	err := json.Unmarshal(payload, tokenURI)
	if err != nil {
		logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("json.Unmarshal", zap.Error(err)))
		return err
	}

	logger.AtLog.Logger.Info("PubSubCreateTokenThumbnai", zap.Any("tokenURI", zap.Any("tokenURI)", tokenURI)))
	token, err := u.Repo.FindTokenByWithoutCache(tokenURI.ContractAddress, tokenURI.TokenID)
	if err != nil {
		logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("FindTokenByWithoutCache", zap.Error(err)))
		return err
	}

	resp, err := u.RunAndCap(token)
	if err != nil {
		logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("RunAndCap", zap.Error(err)))
		return err
	}

	logger.AtLog.Logger.Info("PubSubCreateTokenThumbnai", zap.Any("RunAndCap.resp", zap.Any("resp)", resp)))
	if resp.IsUpdated {
		token.ParsedImage = &resp.ParsedImage
		token.Thumbnail = resp.Thumbnail
		token.ParsedAttributes = resp.Traits
		token.ParsedAttributesStr = resp.TraitsStr
		token.ThumbnailCapturedAt = resp.CapturedAt

		updated, err := u.Repo.UpdateOrInsertTokenUri(tokenURI.ContractAddress, tokenURI.TokenID, token)
		if err != nil {
			logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("UpdateOrInsertTokenUri", err))
		}
		logger.AtLog.Logger.Info("PubSubCreateTokenThumbnai", zap.Any("updatedp", updated), zap.String("tokenID", token.TokenID))
		u.NotifyWithChannel(os.Getenv("SLACK_PROJECT_CHANNEL_ID"), fmt.Sprintf("[Token's thumnail is captured][token %s]", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name)), "", fmt.Sprintf("%s", helpers.CreateTokenImageLink(token.Thumbnail)))

	}

	return nil
}

func (u *Usecase) PubSubProjectUnzip(channelName string, payload interface{}) {
	info, err := u.Queue.CreateSchedule(channelName, payload,  "critical")
	if err != nil {
		logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("taskInfo", info), zap.Any("payload", payload), zap.Error(err))
		return
	}
	logger.AtLog.Info("PubSubCreateTokenThumbnai", zap.Any("taskInfo", info), zap.Any("payload", payload))
}

func (u *Usecase) ProccessUnzip(ctx context.Context, t *asynq.Task) error {
	payload := t.Payload()
	pz := &structure.ProjectUnzipPayload{}
	err := json.Unmarshal(payload, pz)
	if err != nil {
		logger.AtLog.Error("PubSubProjectUnzip", zap.Any("payload", payload), zap.Error(err))
		return err
	}

	_, err = u.UnzipProjectFile(pz)
	if err != nil {
		logger.AtLog.Error("PubSubProjectUnzip", zap.Any("payload", payload), zap.Error(err))
		return err
	}
	return nil
}