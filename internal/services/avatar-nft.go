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
		player.Tattoo,
		player.Glasses,
		player.Captain,
		player.FacePaint,
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
			TraitType: "DNA",
			Value:     player.Dna,
		},
		{
			TraitType: "Team",
			Value:     player.Nation,
		},
		{
			TraitType: "Number",
			Value:     player.Number.String(),
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
			TraitType: "Hair",
			Value:     player.Hair,
		},
		{
			TraitType: "Beard",
			Value:     player.Beard,
		},
		{
			TraitType: "Facial Paint",
			Value:     player.FacePaint,
		},
		{
			TraitType: "Glasses",
			Value:     player.Glasses,
		},
		{
			TraitType: "Boots",
			Value:     player.Shoes,
		},
		{
			TraitType: "Tattoos",
			Value:     player.Tattoo,
		},
		{
			TraitType: "Reacts",
			Value:     player.EmotionTime,
		},
		{
			TraitType: "Captain",
			Value:     player.Captain,
		},
		{
			TraitType: "Undershirt",
			Value:     player.Undershirt,
		},
	}

	s.getAvatarTraitsResolvedValue(ret)

	return ret
}

func (s *service) getAvatarTraitsResolvedValue(attrs []*model.OpenSeaAttribute) {
	data := map[string]map[string]string{
		"DNA": {
			"1": "Male",
			"2": "Female",
			"3": "Robot",
			"4": "Ape",
			"5": "Alien",
			"6": "Golden Ballhead",
		},
		"Top": {
			"1": "Shirt",
			"2": "Hoodie",
			"3": "Topless",
			"4": "Camisole",
		},
		"Bottom": {
			"1": "Shorts",
			"2": "Joggers",
			"3": "Shorts",
			"4": "Skirt",
		},
		"Hair": {
			"0": "None",
			"1": "Short",
			"2": "Long",
			"3": "Wild",
			"4": "Beanie",
			"5": "Sombrero",
			"6": "Cockscomb",
			"7": "Headdress",
			"8": "Infamous",
			"9": "Mohawk",
		},
		"Boots": {
			"1": "Classic",
			"2": "Futuristic",
			"3": "Golden",
		},
		"Tattoos": {
			"0": "None",
			"1": "Bitcoin",
			"2": "Ethereum",
			"3": "Soccer Ball",
		},
		"Beard": {
			"0": "None",
			"1": "Trimmed",
			"2": "Bushy",
		},
		"Glasses": {
			"0": "None",
			"1": "3D",
			"2": "VR",
			"3": "AR",
			"4": "Goggles",
			"5": "Gold",
		},
		"Captain": {
			"0": "No",
			"1": "Yes",
		},
		"Undershirt": {
			"0":  "None",
			"1":  "GM",
			"2":  "WAGMI",
			"3":  "Probably Nothing",
			"4":  "Up Only",
			"5":  "RightClick Save As",
			"6":  "JPGs",
			"7":  "LFG",
			"8":  "Wen Moon",
			"9":  "Shelling Point",
			"10": "Valhalla",
			"11": "McDonald’s",
			"12": "NFA",
			"13": "MFERS",
			"14": "XCOPY",
			"15": "Moonbirds",
			"16": "Nouns",
			"17": "Cryptoadz",
		},
		"Facial Paint": {
			"0": "None",
			"1": "Flag",
		},
		"Reacts": {
			"1": "Daily",
			"2": "Weekly",
			"3": "Monthly",
		},
	}

	for _, attr := range attrs {
		if attr.TraitType == "Team" {
			continue
		} else if attr.TraitType == "Number" {
			if attr.Value == "0" {
				attr.Value = "None"
			}
		} else {
			if item, ok := data[attr.TraitType]; ok {
				if val, ok := item[attr.Value]; ok {
					attr.Value = val
				}
			}
		}
	}
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
		Image:       "https://cdn.rove.to/events/world-cup/wc-thumbnail.png",
		Description: pointerutil.MakePointer("Build your World Cup dream team with PLAYΞRS, a limited collection of 10,000 living 3D generative avatars. Every PLAYΞRS is algorithmically generated from 16 different attributes, with expressions that react to their team’s success on the field. Start your collection today."),
	}, nil
}
