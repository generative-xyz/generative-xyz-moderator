package usecase

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/helpers"
)

func (u *Usecase) PubSubCreateTokenThumbnail(tracingInjection map[string]string, channelName string, payload interface{}) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		u.Logger.ErrorAny("PubSubCreateTokenThumbnai", zap.Any("json.Marshal", zap.Error(err)))
		return
	}

	tokenURI := &structure.TokenImagePayload{}
	err = json.Unmarshal(bytes, tokenURI)
	if err != nil {
		u.Logger.ErrorAny("PubSubCreateTokenThumbnai", zap.Any("json.Unmarshal", zap.Error(err)))
		return
	}

	u.Logger.LogAny("PubSubCreateTokenThumbnai", zap.Any("tokenURI", tokenURI))
	token, err := u.Repo.FindTokenByWithoutCache(tokenURI.ContractAddress, tokenURI.TokenID)
	if err != nil {
		u.Logger.ErrorAny("PubSubCreateTokenThumbnai", zap.Any("FindTokenByWithoutCache", zap.Error(err)))
		return
	}

	
   resp, err := u.RunAndCap(token, 20)
	if err != nil {
		u.Logger.ErrorAny("PubSubCreateTokenThumbnai", zap.Any("RunAndCap", zap.Error(err)))
		return
	}

	u.Logger.LogAny("PubSubCreateTokenThumbnai", zap.Any("RunAndCap.resp", resp))
	if resp.IsUpdated {
		token.ParsedImage = &resp.ParsedImage
		token.Thumbnail = resp.Thumbnail
		token.ParsedAttributes = resp.Traits
		token.ParsedAttributesStr = resp.TraitsStr
		token.ThumbnailCapturedAt = resp.CapturedAt

		updated, err := u.Repo.UpdateOrInsertTokenUri(tokenURI.ContractAddress, tokenURI.TokenID, token)
		if err != nil {
			u.Logger.ErrorAny("PubSubCreateTokenThumbnai",zap.Any("UpdateOrInsertTokenUri", err))
		}
		u.Logger.LogAny("PubSubCreateTokenThumbnai", zap.Any("updatedp", updated), zap.String("tokenID", token.TokenID))
		u.NotifyWithChannel(os.Getenv("SLACK_PROJECT_CHANNEL_ID"), fmt.Sprintf("[Token's thumnail is captured][token %s]", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name)),"", fmt.Sprintf("%s", helpers.CreateTokenImageLink(token.Thumbnail) ))

	}
}

func (u *Usecase) PubSubProjectUnzip(tracingInjection map[string]string, channelName string, payload interface{}) {


	bytes, err := json.Marshal(payload)
	if err != nil {
		u.Logger.ErrorAny("PubSubProjectUnzip", zap.Any("payload", payload), zap.Error(err))
		return
	}

	pz := &structure.ProjectUnzipPayload{}
	err = json.Unmarshal(bytes, pz)
	if err != nil {
		u.Logger.ErrorAny("PubSubProjectUnzip", zap.Any("payload", payload), zap.Error(err))
		return
	}

	_, err = u.UnzipProjectFile(pz)
	if err != nil {
		u.Logger.ErrorAny("PubSubProjectUnzip", zap.Any("payload", payload), zap.Error(err))
		return
	}


}