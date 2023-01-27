package usecase

import (
	"encoding/json"

	"rederinghub.io/internal/usecase/structure"
)

func (u *Usecase) PubSubCreateTokenThumbnail(tracingInjection map[string]string, channelName string, payload interface{}) {
	span, log := u.StartSpanFromInjecttion(tracingInjection, "PubSubCreateTokenThumbnail")
	defer u.Tracer.FinishSpan(span, log)

	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Error("PubSubCreateTokenThumbnai.json.Marshal", err.Error(), err)
		return
	}

	tokenURI := &structure.TokenImagePayload{}
	err = json.Unmarshal(bytes, tokenURI)
	if err != nil {
		log.Error("PubSubCreateTokenThumbnai.json.Unmarshal", err.Error(), err)
		return
	}

	log.SetData("payload", tokenURI)
	log.SetData("tokenID", tokenURI.TokenID)
	log.SetTag("tokenID", tokenURI.TokenID)

	token, err := u.Repo.FindTokenByWithoutCache(tokenURI.ContractAddress, tokenURI.TokenID)
	if err != nil {
		log.Error("PubSubCreateTokenThumbnai.FindTokenBy", err.Error(), err)
		return
	}
	
	log.SetTag("found.tokenID", token.TokenID)
	log.SetTag("found.tokenID.thumbnail", token.Thumbnail)
	
	resp, err := u.RunAndCap(span, token, 20)
	if err != nil {
		log.Error("PubSubCreateTokenThumbnai.RunAndCap", err.Error(), err)
		return
	}

	log.SetData("resp",resp.CapturedAt)
	log.SetData("IsUpdated",resp.IsUpdated)
	log.SetData("Thumbnail",resp.Thumbnail)
	log.SetData("IsUpdated",resp.IsUpdated)

	if resp.IsUpdated {
		token.ParsedImage = &resp.ParsedImage
		token.Thumbnail = resp.Thumbnail
		token.ParsedAttributes = resp.Traits
		token.ParsedAttributesStr = resp.TraitsStr
		token.ThumbnailCapturedAt = resp.CapturedAt

		updated, err := u.Repo.UpdateOrInsertTokenUri(tokenURI.ContractAddress, tokenURI.TokenID, token)
		if err != nil {
			log.Error("PubSubCreateTokenThumbnai.UpdateOrInsertTokenUri", err.Error(), err)
		}
		log.SetData("updated", updated)
	}

}