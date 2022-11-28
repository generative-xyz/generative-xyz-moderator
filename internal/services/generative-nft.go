package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"rederinghub.io/api"
	"rederinghub.io/internal/adapter"
	"rederinghub.io/internal/dto"
	"rederinghub.io/internal/model"
	"rederinghub.io/pkg/config"
	"rederinghub.io/pkg/contracts/confetti_contract"
	"rederinghub.io/pkg/contracts/generative_param"
	"rederinghub.io/pkg/contracts/horn_contract"
	"rederinghub.io/pkg/logger"
)

func (s *service) GetGenerativeNFTMetadataPost(ctx context.Context, req *api.GetGenerativeNFTMetadataRequest) (*api.GetGenerativeNFTMetadataResponse, error) {
	logger.AtLog.Infof("Handle [GetGenerativeNFTMetadata] %s %s %s", req.ChainId, req.ContractAddress, req.TokenId)

	var templateDTOFromMongo bson.M
	if err := s.templateRepository.FindOne(context.Background(), map[string]interface{}{
		"nftInfo.contractAddress": req.ContractAddress,
		"nftInfo.chainId":         req.ChainId,
	}, &templateDTOFromMongo); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, TemplateNotFoundError{ContractAddress: req.ContractAddress, ChainID: req.ChainId}
		}
		return nil, err
	}

	var template dto.TemplateDTO
	{
		doc, err := json.Marshal(templateDTOFromMongo)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(doc, &template); err != nil {
			return nil, err
		}
	}

	chainURL, ok := GetRPCURLFromChainID(req.ChainId)
	if !ok {
		return nil, errors.New("missing config chain_config from server")
	}

	// find in mongo
	var renderedNftBson bson.M
	err := s.renderedNftRepository.FindOne(context.Background(), map[string]interface{}{
		"chainId":         req.ChainId,
		"contractAddress": template.NftInfo.ContractAddress,
		"tokenId":         req.TokenId,
	}, &renderedNftBson)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	// found in mongo
	if err == nil {
		var renderedNft model.RenderedNft
		{
			doc, err := json.Marshal(renderedNftBson)
			if err != nil {
				return nil, err
			}
			if err = json.Unmarshal(doc, &renderedNft); err != nil {
				return nil, err
			}
		}

		return renderedNft.ToGenerativeProto(template), nil
	}

	params, obj, err := s.getParamFromContract(template, chainURL, req)
	if err != nil {
		logger.AtLog.Errorf("Error getParamFromContract %v", err)
		return nil, err
	}

	// call to render machine
	rendered, err := s.renderMachineAdapter.Render(ctx, &adapter.RenderRequest{
		Script:      template.Script,
		Params:      params,
		BlenderType: template.BlenderType,
		RenderImage: true,
		RenderModel: true,
		Seed:        "1",
	})

	if err != nil {
		return nil, err
	}

	attributes := make([]generative_param.TraitInfoTrait, 0)

	// save to mongo
	renderedNft, err := GetRenderedNft(
		req.ChainId,
		req.ContractAddress,
		template.BlenderType,
		req.TokenId,
		template,
		attributes,
		fmt.Sprintf("ipfs://%v", rendered.Image),
		fmt.Sprintf("ipfs://%v", rendered.Glb),
	)
	if err != nil {
		return nil, err
	}

	renderedNft.Attributes = s.getOpenSeaTraits(template, obj)

	var renderedNftModel bson.M
	if _bytes, err := json.Marshal(&renderedNft); err != nil {
		return nil, err
	} else {
		if err = json.Unmarshal(_bytes, &renderedNftModel); err != nil {
			return nil, err
		}
	}

	_, err = s.renderedNftRepository.Create(context.Background(), &renderedNftModel)
	if err != nil {
		logger.AtLog.Errorf("[GetRenderedNft] Create error %v", err)
		return nil, err
	}

	logger.AtLog.Infof("Done [GetGenerativeNFTMetadataPost] #%s", req.TokenId)
	appConfig := config.AppConfig()
	go func() {
		if err := s.moralisAdapter.ResyncNFTMetadata(req.ContractAddress, appConfig.ChainID, req.TokenId); err != nil {
			logger.AtLog.Errorf("error when call resync nft metadata %v", err)
		}
	}()

	return renderedNft.ToGenerativeProto(template), nil
}

func (s *service) getParamFromContract(template dto.TemplateDTO, chainURL string, req *api.GetGenerativeNFTMetadataRequest) ([]string, interface{}, error) {
	client, err := ethclient.Dial(chainURL)

	if err != nil {
		return nil, nil, err
	}

	addr := common.HexToAddress(template.NftInfo.ContractAddress)
	tokenID, ok := new(big.Int).SetString(req.TokenId, 10)
	if !ok {
		return nil, nil, InvalidTokenIDError{TokenID: req.TokenId}
	}
	switch template.BlenderType {
	case "confetti":
		instance, err := confetti_contract.NewConfettiContract(addr, client)
		if err != nil {
			return nil, nil, err
		}
		confetti, err := instance.GetParamValues(&bind.CallOpts{Context: context.Background(), Pending: false}, tokenID)
		if err != nil {
			return nil, nil, err
		}
		params := []string{
			confetti.ShapeCanon,
			confetti.ShapeConfetti,
		}
		for _, v := range confetti.PalletteCanon {
			params = append(params, v)
		}
		for _, v := range confetti.PalletteConfetti {
			params = append(params, v)
		}

		return params, confetti, nil
	case "horn":
		instance, err := horn_contract.NewHornContract(addr, client)
		if err != nil {
			return nil, nil, err
		}
		horn, err := instance.GetParamValues(&bind.CallOpts{Context: context.Background(), Pending: false}, tokenID)
		if err != nil {
			return nil, nil, err
		}
		params := []string{
			horn.Nation,
			horn.PalletTop,
			horn.PalletBottom,
		}

		return params, horn, nil
	default:
		return nil, nil, errors.New(fmt.Sprintf("dont know blender_type = %v from template that have contract_address %v", template.BlenderType, template.NftInfo.ContractAddress))
	}
}

func (s *service) getOpenSeaTraits(template dto.TemplateDTO, paramObj interface{}) []*model.OpenSeaAttribute {
	var ret []*model.OpenSeaAttribute
	switch template.BlenderType {
	case "confetti":
		confetti := paramObj.(confetti_contract.CONFETTIConfetti)
		ret = []*model.OpenSeaAttribute{
			{
				TraitType: "CANON SHAPE",
				Value:     confetti.ShapeCanon,
			},
			{
				TraitType: "CONFETTI SHAPE",
				Value:     confetti.ShapeConfetti,
			},
			{
				TraitType: "CANON COLOR",
				Value:     strings.Join(confetti.PalletteCanon[0:], ", "),
			},
			{
				TraitType: "CONFETTI COLOR",
				Value:     strings.Join(confetti.PalletteConfetti[0:], ", "),
			},
		}
	case "horn":
		horn := paramObj.(horn_contract.HORNSHorn)
		ret = []*model.OpenSeaAttribute{
			{
				TraitType: "Nation",
				Value:     horn.Nation,
			},
			{
				TraitType: "TOP COLOR",
				Value:     horn.PalletTop,
			},
			{
				TraitType: "BOTTOM COLOR",
				Value:     horn.PalletBottom,
			},
		}
	}

	s.getOpenseaTraitsResolvedValue(ret, template)

	return ret
}

func (s *service) getOpenseaTraitsResolvedValue(attrs []*model.OpenSeaAttribute, template dto.TemplateDTO) {
	for _, attr := range attrs {
		if item, ok := template.MappingAttribute[attr.TraitType]; ok {
			if val, ok := item[attr.Value]; ok {
				attr.Value = val
			}
		}
	}
}

func (s *service) GetGenerativeNFTMetadata(ctx context.Context, req *api.GetGenerativeNFTMetadataRequest) (*api.GetGenerativeNFTMetadataResponse, error) {

	logger.AtLog.Infof("Handle [GetGenerativeNFTMetadata] %s %s %s", req.ChainId, req.ContractAddress, req.TokenId)

	var templateDTOFromMongo bson.M
	if err := s.templateRepository.FindOne(context.Background(), map[string]interface{}{
		"nftInfo.contractAddress": req.ContractAddress,
		"nftInfo.chainId":         req.ChainId,
	}, &templateDTOFromMongo); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, TemplateNotFoundError{ContractAddress: req.ContractAddress, ChainID: req.ChainId}
		}
		return nil, err
	}

	var template dto.TemplateDTO
	{
		doc, err := json.Marshal(templateDTOFromMongo)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(doc, &template); err != nil {
			return nil, err
		}
	}

	// find in mongo
	var renderedNftBson bson.M
	err := s.renderedNftRepository.FindOne(context.Background(), map[string]interface{}{
		"chainId":         req.ChainId,
		"contractAddress": template.NftInfo.ContractAddress,
		"tokenId":         req.TokenId,
	}, &renderedNftBson)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	// found in mongo
	if err == nil {
		var renderedNft model.RenderedNft
		{
			doc, err := json.Marshal(renderedNftBson)
			if err != nil {
				return nil, err
			}
			if err = json.Unmarshal(doc, &renderedNft); err != nil {
				return nil, err
			}
		}

		return renderedNft.ToGenerativeProto(template), nil
	}

	return &api.GetGenerativeNFTMetadataResponse{
		Name:        fmt.Sprintf("Rendering on #%s", req.TokenId),
		Image:       "https://cdn.rove.to/metaverse/rove/Rove_logo.png",
		Description: template.Description,
	}, nil
}
