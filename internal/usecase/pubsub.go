package usecase

import (
	"encoding/json"

	"go.uber.org/zap"
	"rederinghub.io/internal/usecase/structure"
)

func (u *Usecase) PubSubCreateTokenThumbnail(tracingInjection map[string]string, channelName string, payload interface{}) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		u.Logger.ErrorAny("PubSubCreateTokenThumbnai", zap.Any("json.Marshal", err))
		return
	}

	tokenURI := &structure.TokenImagePayload{}
	err = json.Unmarshal(bytes, tokenURI)
	if err != nil {
		u.Logger.ErrorAny("PubSubCreateTokenThumbnai", zap.Any("json.Unmarshal", err))
		return
	}

	u.Logger.LogAny("PubSubCreateTokenThumbnai", zap.Any("tokenURI", tokenURI))
	token, err := u.Repo.FindTokenByWithoutCache(tokenURI.ContractAddress, tokenURI.TokenID)
	if err != nil {
		u.Logger.Error("PubSubCreateTokenThumbnai.FindTokenBy", err.Error(), err)
		return
	}

	
   resp, err := u.RunAndCap(token, 20)
	if err != nil {
		u.Logger.Error("PubSubCreateTokenThumbnai.RunAndCap", err.Error(), err)
		return
	}

	u.Logger.Info("resp",resp.CapturedAt)
	u.Logger.Info("IsUpdated",resp.IsUpdated)
	u.Logger.Info("Thumbnail",resp.Thumbnail)
	u.Logger.Info("IsUpdated",resp.IsUpdated)

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
		u.Logger.Info("updated", updated)
	}
}

func (u *Usecase) PubSubProjectUnzip(tracingInjection map[string]string, channelName string, payload interface{}) {


	bytes, err := json.Marshal(payload)
	if err != nil {
		u.Logger.Error("PubSubProjectUnzipl.json.Marshal", err.Error(), err)
		return
	}

	pz := &structure.ProjectUnzipPayload{}
	err = json.Unmarshal(bytes, pz)
	if err != nil {
		u.Logger.Error("PubSubProjectUnzipl.json.Unmarshal", err.Error(), err)
		return
	}

	u.Logger.Info("payload", pz)
	u.Logger.Info("projectID", pz.ProjectID)
	


	pe, err := u.UnzipProjectFile(pz)
	if err != nil {
		u.Logger.Error("PubSubProjectUnzipl.json.Unmarshal", err.Error(), err)
		return
	}

	u.Logger.Info("project", pe)

}