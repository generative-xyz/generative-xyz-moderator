package usecase

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"
	_req "rederinghub.io/utils/request"
)

//Queue functions

// Processing functions
func (u *Usecase) PubSubCreateTokenThumbnail(tracingInjection map[string]string, channelName string, payload interface{}) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		logger.AtLog.Error("PubSubCreateTokenThumbnail", zap.Any("json.Marshal", zap.Error(err)))
		return
	}

	tokenURI := &structure.TokenImagePayload{}
	err = json.Unmarshal(bytes, tokenURI)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("PubSubCreateTokenThumbnail - %s - %s", tokenURI.ContractAddress, tokenURI.TokenID), zap.Any("json.Unmarshal", zap.Error(err)))
		return
	}

	//Save to DB

	//move this to another Job
	u.Capture(tokenURI)
}

func (u *Usecase) Capture(tokenURI *structure.TokenImagePayload) {
	var err error
	token := new(entity.TokenUri)
	resp := new(structure.TokenAnimationURI)

	defer func() {
		if err == nil {
			logger.AtLog.Logger.Info(fmt.Sprintf("PubSubCreateTokenThumbnail - %s", tokenURI.TokenID), zap.Any("request", tokenURI), zap.String("tokenURI", tokenURI.TokenID), zap.String("contract_address", tokenURI.ContractAddress), zap.String("gen_nft_addrress", token.GenNFTAddr), zap.String("Thumbnail", resp.Thumbnail), zap.Bool("Thumbnail.IsUpdated", resp.IsUpdated))

		} else {
			logger.AtLog.Logger.Error(fmt.Sprintf("PubSubCreateTokenThumbnail - %s", tokenURI.TokenID), zap.Any("request", tokenURI), zap.String("tokenURI", tokenURI.TokenID), zap.String("contract_address", tokenURI.ContractAddress), zap.Error(err))
		}
	}()

	token, err = u.Repo.FindTokenForCaptureThumbnail(tokenURI.ContractAddress, tokenURI.TokenID)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("PubSubCreateTokenThumbnail - %s - %s", tokenURI.ContractAddress, tokenURI.TokenID), zap.Error(err), zap.String("tokenURI", tokenURI.TokenID), zap.String("contract_address", tokenURI.ContractAddress))
		return
	}

	resp, err = u.RunAndCap(token)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("PubSubCreateTokenThumbnail - %s - %s", tokenURI.ContractAddress, tokenURI.TokenID), zap.Error(err), zap.String("tokenURI", tokenURI.TokenID), zap.String("contract_address", tokenURI.ContractAddress))
		return
	}

	if resp.IsUpdated {
		err = u.Repo.UpdateTokenThumbnail(tokenURI.ContractAddress, tokenURI.TokenID, resp.Thumbnail, resp.ParsedImage, resp.Traits, resp.TraitsStr, resp.CapturedAt)
		if err != nil {
			logger.AtLog.Error(fmt.Sprintf("PubSubCreateTokenThumbnail - %s - %s", tokenURI.ContractAddress, tokenURI.TokenID), zap.Error(err), zap.String("tokenURI", tokenURI.TokenID), zap.String("contract_address", tokenURI.ContractAddress))
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

func (u *Usecase) PubSubEthProjectUnzip(tracingInjection map[string]string, channelName string, payload interface{}) {

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

	_, err = u.UnzipETHProjectFile(pz)
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

func (u *Usecase) PubsubTestChannel(tracingInjection map[string]string, channelName string, payload interface{}) {

	logger.AtLog.Logger.Info("PubsubTestChannel", zap.String("channelName", channelName), zap.Any("tracingInjection", tracingInjection))
}
