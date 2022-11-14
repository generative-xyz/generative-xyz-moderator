package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/pkg/contracts/avatar_contract"
	"rederinghub.io/pkg/utils/pointerutil"

	"rederinghub.io/api"
	"rederinghub.io/internal/adapter"
	"rederinghub.io/internal/dto"
	"rederinghub.io/internal/model"
	"rederinghub.io/pkg/contracts/generative_param"
	"rederinghub.io/pkg/logger"
)

func (s *service) GetAvatarMetadataPost(ctx context.Context, req *api.GetAvatarMetadataRequest) (*api.GetAvatarMetadataResponse, error) {

	//req.ContractAddress = ""
	req.ProjectId = "140292"

	logger.AtLog.Infof("Handle [GetCandyMetadataPost] %s %s %s %s", req.ChainId, req.ContractAddress, req.ProjectId, req.TokenId)

	chainURL, ok := GetRPCURLFromChainID(req.ChainId)
	if !ok {
		return nil, errors.New("missing config chain_config from server")
	}

	var templateDTOFromMongo bson.M
	if err := s.templateRepository.FindOne(context.Background(), map[string]interface{}{
		"nftInfo.tokenId":         req.ProjectId,
		"nftInfo.contractAddress": req.ContractAddress,
		"nftInfo.chainId":         req.ChainId,
	}, &templateDTOFromMongo); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, TemplateNotFoundError{TokenID: req.ProjectId, ChainID: req.ChainId}
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
		"projectId":       req.ProjectId,
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

		return renderedNft.ToAvatarResponse(), nil
	}

	client, err := ethclient.Dial(chainURL)

	if err != nil {
		return nil, err
	}
	addr := common.HexToAddress(template.NftInfo.ContractAddress)

	instance, err := avatar_contract.NewAvatarContract(addr, client)

	if err != nil {
		return nil, err
	}

	tokenID, ok := new(big.Int).SetString(req.TokenId, 10)
	if !ok {
		return nil, InvalidTokenIDError{TokenID: req.TokenId}
	}

	player, err := instance.GetParamValues(&bind.CallOpts{Context: context.Background(), Pending: false}, tokenID)
	if err != nil {
		return nil, err
	}

	params := []string{
		player.Emotion,
		player.Nation,
		player.Dna,
		player.Beard,
		player.Hair,
		player.Undershirt,
		player.Shoes,
		player.Top,
		player.Bottom,
		player.Number.String(),
		player.Tatoo,
		player.Glasses,
		player.Captain,
	}

	// call to render machine
	rendered, err := s.renderMachineAdapter.Render(ctx, &adapter.RenderRequest{
		Script:      template.Script,
		Params:      params,
		Seed:        "1",
		BlenderType: "avatar",
	})

	if err != nil {
		return nil, err
	}

	attributes := make([]generative_param.TraitInfoTrait, 0)

	// save to mongo
	renderedNft, err := GetRenderedNft(
		req.ChainId,
		req.ContractAddress,
		req.ProjectId,
		req.TokenId,
		template,
		attributes,
		fmt.Sprintf("ipfs://%v", rendered.Image),
		fmt.Sprintf("ipfs://%v", rendered.Glb),
	)
	if err != nil {
		return nil, err
	}
	renderedNft.Attributes = s.getAvatarOpenSeaTraits(&player)
	renderedNft.EmotionTime = player.EmotionTime
	if rendered.Metadata.BackgroundColor != "" {
		renderedNft.Metadata = &model.RenderedNftMetadata{BackgroundColor: &rendered.Metadata.BackgroundColor}
	}

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
		logger.AtLog.Errorf("[GetAvatarMetadataPost] Create error %v", err)
		return nil, err
	}

	logger.AtLog.Infof("Done [GetAvatarMetadataPost] #%s", req.TokenId)

	return renderedNft.ToAvatarResponse(), nil
}

func (s *service) getAvatarOpenSeaTraits(player *avatar_contract.AVATARSPlayer) []*model.OpenSeaAttribute {
	ret := []*model.OpenSeaAttribute{
		{
			TraitType: "Emotion",
			Value:     player.Emotion,
		},
		{
			TraitType: "Nation",
			Value:     player.Nation,
		},
		{
			TraitType: "Dna",
			Value:     player.Dna,
		},
		{
			TraitType: "Beard",
			Value:     player.Beard,
		},
		{
			TraitType: "Hair",
			Value:     player.Hair,
		},
		{
			TraitType: "Undershirt",
			Value:     player.Undershirt,
		},
		{
			TraitType: "Shoes",
			Value:     player.Shoes,
		},
		{
			TraitType: "Top",
			Value:     player.Top,
		},
		{
			TraitType: "Bottom",
			Value:     player.Bottom,
		},
		{
			TraitType: "Number",
			Value:     player.Number.String(),
		},
		{
			TraitType: "Tatoo",
			Value:     player.Tatoo,
		},
		{
			TraitType: "Glasses",
			Value:     player.Glasses,
		},
		{
			TraitType: "Captain",
			Value:     player.Captain,
		},
	}
	return ret
}

func (s *service) GetAvatarMetadata(ctx context.Context, req *api.GetAvatarMetadataRequest) (*api.GetAvatarMetadataResponse, error) {
	//req.ContractAddress = ""
	//req.ContractAddress = ""
	req.ProjectId = "140292"

	logger.AtLog.Infof("Handle [GetCandyMetadataPost] %s %s %s %s", req.ChainId, req.ContractAddress, req.ProjectId, req.TokenId)

	var templateDTOFromMongo bson.M
	if err := s.templateRepository.FindOne(context.Background(), map[string]interface{}{
		"nftInfo.tokenId":         req.ProjectId,
		"nftInfo.contractAddress": req.ContractAddress,
		"nftInfo.chainId":         req.ChainId,
	}, &templateDTOFromMongo); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, TemplateNotFoundError{TokenID: req.ProjectId, ChainID: req.ChainId}
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
		"projectId":       req.ProjectId,
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

		return renderedNft.ToAvatarResponse(), nil
	}

	return &api.GetAvatarMetadataResponse{
		Name:        fmt.Sprintf("Rendering on #%s", req.TokenId),
		Image:       "https://i.seadn.io/gae/iFdea-Nd80jLGqELcfuBuygkdzlwekgyXxoWnh6z7oSUfapJgsfgQs2HABhWTU1xsbyGhLRhkXneprVAG40OWKhB2YzQbW_69UmMuw?auto=format&w=256",
		Description: pointerutil.MakePointer("FOOTBALLΞR is a living character that reacts to the games as you do. Through Chainlink oracles connected to match results, your avatar will smile, frown, and tire as its team wins or loses on the field."),
	}, nil
}
