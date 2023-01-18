package usecase

import (
	"encoding/json"

	"rederinghub.io/internal/entity"
)

func (u *Usecase) PubSubCreateTokenThumbnail(tracingInjection map[string]string, channelName string, payload interface{}) {
	span, log := u.StartSpanFromInjecttion(tracingInjection, "PubSubCreateTokenThumbnail")
	defer u.Tracer.FinishSpan(span, log)

	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Error("PubSubCreateTokenThumbnai.json.Marshal", err.Error(), err)
		return
	}

	tokenURI := &entity.TokenUri{}
	err = json.Unmarshal(bytes, tokenURI)
	if err != nil {
		log.Error("PubSubCreateTokenThumbnai.json.Unmarshal", err.Error(), err)
		return
	}
	log.SetData("tokenURI", tokenURI.TokenID)
	log.SetTag("tokenURI", tokenURI.TokenID)

	resp, err := u.RunAndCap(span, tokenURI, 20)
	if err != nil {
		log.Error("PubSubCreateTokenThumbnai.RunAndCap", err.Error(), err)
		return
	}

	token, err := u.Repo.FindTokenBy(tokenURI.ContractAddress, tokenURI.TokenID)
	if err != nil {
		log.Error("PubSubCreateTokenThumbnai.FindTokenBy", err.Error(), err)
		return
	}

	token.ParsedImage = &resp.ParsedImage
	token.Thumbnail = resp.Thumbnail
	token.ParsedAttributes = resp.Traits
	token.ParsedAttributesStr = resp.TraitsStr

	log.SetData("resp",resp.CapturedAt)
	log.SetData("IsUpdated",resp.IsUpdated)
	log.SetData("Thumbnail",resp.Thumbnail)

	updated, err := u.Repo.UpdateOrInsertTokenUri(tokenURI.ContractAddress, tokenURI.TokenID, token)
	if err != nil {
		log.Error("PubSubCreateTokenThumbnai.UpdateOrInsertTokenUri", err.Error(), err)
	}
	log.SetData("updated", updated)
}