package usecase

import (
	"encoding/json"
	"fmt"
	"os"
	"rederinghub.io/utils/helpers"

	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"
	_req "rederinghub.io/utils/request"
)

//Queue functions

// Processing functions
func (u *Usecase) PubSubCreateTokenThumbnail(tracingInjection map[string]string, channelName string, payload interface{}) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("json.Marshal", zap.Error(err)))
		return
	}

	tokenURI := &structure.TokenImagePayload{}
	err = json.Unmarshal(bytes, tokenURI)
	if err != nil {
		logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("json.Unmarshal", zap.Error(err)))
		return
	}

	u.Capture(tokenURI)
}

func (u *Usecase) Capture(tokenURI *structure.TokenImagePayload) {
	logger.AtLog.Logger.Info("PubSubCreateTokenThumbnai", zap.Any("tokenURI", zap.Any("tokenURI)", tokenURI)))
	token, err := u.Repo.FindTokenByWithoutCache(tokenURI.ContractAddress, tokenURI.TokenID)
	if err != nil {
		logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("FindTokenByWithoutCache", zap.Error(err)))
		return
	}

	resp, err := u.RunAndCap(token)
	if err != nil {
		logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("RunAndCap", zap.Error(err)))
		return
	}

	logger.AtLog.Logger.Info("PubSubCreateTokenThumbnai", zap.Any("RunAndCap.resp", zap.Any("resp)", resp)))
	if resp.IsUpdated {
		token.ParsedImage = &resp.ParsedImage
		token.Thumbnail = resp.Thumbnail
		token.ParsedAttributes = resp.Traits
		token.ParsedAttributesStr = resp.TraitsStr
		token.ThumbnailCapturedAt = resp.CapturedAt

		if os.Getenv("ENV") == "mainnet" {
			updated, err := u.Repo.UpdateOrInsertTokenUri(tokenURI.ContractAddress, tokenURI.TokenID, token)
			if err != nil {
				logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("UpdateOrInsertTokenUri", err))
			}
			logger.AtLog.Logger.Info("PubSubCreateTokenThumbnai", zap.Any("updatedp", updated), zap.String("tokenID", token.TokenID))
			u.NotifyWithChannel(os.Getenv("SLACK_PROJECT_CHANNEL_ID"), fmt.Sprintf("[Token's thumnail is captured][token %s]", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name)), "", fmt.Sprintf("%s", helpers.CreateTokenImageLink(token.Thumbnail)))
		}

	}
}

func (u *Usecase) PubSubProjectUnzip(tracingInjection map[string]string, channelName string, payload interface{}) {

	bytes, err := json.Marshal(payload)
	if err != nil {
		logger.AtLog.Error("PubSubProjectUnzip", zap.Any("payload", payload), zap.Error(err))
		return
	}

	pz := &structure.ProjectUnzipPayload{}
	err = json.Unmarshal(bytes, pz)
	if err != nil {
		logger.AtLog.Error("PubSubProjectUnzip", zap.Any("payload", payload), zap.Error(err))
		return
	}

	_, err = u.UnzipProjectFile(pz)
	if err != nil {
		logger.AtLog.Error("PubSubProjectUnzip", zap.Any("payload", payload), zap.Error(err))
		return
	}

}

func (u *Usecase) PubSubCaptureThumbnail(tracingInjection map[string]string, channelName string, payload interface{}) {

	bytes, err := json.Marshal(payload)
	if err != nil {
		logger.AtLog.Error("PubSubProjectUnzip", zap.Any("payload", payload), zap.Error(err))
		return
	}

	req := &request.CaptureRequest{}
	err = json.Unmarshal(bytes, req)
	if err != nil {
		logger.AtLog.Error("PubSubProjectUnzip", zap.Any("payload", payload), zap.Error(err))
		return
	}

	url, err := u.CaptureContent(req.ID, req.Url)
	if err != nil {
		logger.AtLog.Error("PubSubProjectUnzip", zap.Any("payload", payload), zap.Error(err))
		return
	}

	renderURL := fmt.Sprintf("%s/api/v1/device/%s/renderer-set-image", os.Getenv("RENDER_DOMAIN"), req.ID)

	postPayload := map[string]string{
		"image_url": url,
	}
	code, result, err := _req.PostToRenderer(renderURL, postPayload)

	if err != nil {
		logger.AtLog.Error("PubSubProjectUnzip", zap.Int("code", code), zap.Any("postPayload", postPayload), zap.Error(err))
		return
	}
	//logger.AtLog.Info("call to renderer-set-image ", zap.Error(err), zap.String("renderURL", renderURL), zap.String("device_id", req.ID), zap.Int("code", code), zap.String("response", string(result)))
	fmt.Printf(`[POST] %s request {"display_url": "%s"}  - resp {"code": "%d", "result":"%v", "error": "%v" } \n`, renderURL, req.Url, code, string(result), err)
}
