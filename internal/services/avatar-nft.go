package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"rederinghub.io/api"
	"rederinghub.io/internal/dto"
	"rederinghub.io/internal/model"
	"rederinghub.io/pkg/contracts/avatar_contract"
	"rederinghub.io/pkg/logger"
	"rederinghub.io/pkg/utils/pointerutil"
)

func (s *service) hotfixResponseAvatarPost(avatar model.RenderedNft) *api.GetAvatarMetadataResponse {
	resp := avatar.ToAvatarResponse()
	resp.AnimationUrl = resp.GlbUrl
	return resp
}

func (s *service) getAvatarCacheRedisKey(chainID, contractAddress, tokenID string) string {
	return fmt.Sprintf("Hub_%s_%s_%s", chainID, strings.ToLower(contractAddress), tokenID)
}

func (s *service) GetAvatarMetadataPost(ctx context.Context, req *api.GetAvatarMetadataRequest) (*api.GetAvatarMetadataResponse, error) {
	//req.ContractAddress = ""
	req.ProjectId = "140292"

	logger.AtLog.Infof("Handle [GetAvatarMetadataPost] %s %s %s %s", req.ChainId, req.ContractAddress, req.ProjectId, req.TokenId)

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

	// find avatar emotion
	var emotion string

	redisKey := s.getAvatarCacheRedisKey(req.ChainId, req.ContractAddress, req.TokenId)
	cached := s.redisClient.Get(ctx, redisKey).Val()
	if len(cached) > 0 {
		emotion = cached
		logger.AtLog.Infof("Cache hit %v", emotion)
	} else {
		chainURL, ok := GetRPCURLFromChainID(req.ChainId)
		if !ok {
			return nil, errors.New("missing config chain_config from server")
		}

		// call to contract to get emotion
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

		emotion = player.Emotion

		s.redisClient.Set(ctx, redisKey, emotion, 5*time.Minute)
		logger.AtLog.Infof("Cache missed %v", emotion)
	}

	var emotionInt int
	if i, err := strconv.Atoi(emotion); err == nil {
		emotionInt = i
	}

	// find in mongo
	var renderedNftBson bson.M
	err := s.renderedNftRepository.FindOne(context.Background(), map[string]interface{}{
		"chainId":          req.ChainId,
		"contractAddress":  template.NftInfo.ContractAddress,
		"projectId":        req.ProjectId,
		"tokenId":          req.TokenId,
		"metadata.emotion": emotionInt,
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

		return s.hotfixResponseAvatarPost(renderedNft), nil
	}

	return nil, errors.New("avatar is not rendered")
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
			"6": "Ballon d’Or",
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
			"4": "Gold",
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

func (s *service) getPlayerAvatarCacheKey(tokenID string) string {
	return fmt.Sprintf("player_emotion:%v", tokenID)
}

func (s *service) GetClearCacheInternal(ctx context.Context, req *api.GetClearCacheInternalRequest) (*api.GetClearCacheInternalResponse, error) {
	redisKey := s.getAvatarCacheRedisKey(req.ChainId, req.ContractAddress, req.TokenId)
	if err := s.redisClient.Del(ctx, redisKey).Err(); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &api.GetClearCacheInternalResponse{
		Message: "OK",
	}, nil
}

func (s *service) GetAvatarMetadata(ctx context.Context, req *api.GetAvatarMetadataRequest) (*api.GetAvatarMetadataResponse, error) {
	//req.ContractAddress = ""
	//req.ContractAddress = ""
	req.ProjectId = "140292"

	logger.AtLog.Infof("Handle [GetAvatarMetadata] %s %s %s %s", req.ChainId, req.ContractAddress, req.ProjectId, req.TokenId)

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

	// find avatar emotion
	var emotion string
	redisKey := s.getAvatarCacheRedisKey(req.ChainId, req.ContractAddress, req.TokenId)
	cached := s.redisClient.Get(ctx, redisKey).Val()
	if len(cached) > 0 {
		emotion = cached
		logger.AtLog.Infof("Cache hit %v", emotion)
	} else {
		chainURL, ok := GetRPCURLFromChainID(req.ChainId)
		if !ok {
			return nil, errors.New("missing config chain_config from server")
		}

		// call to contract to get emotion
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

		emotion = player.Emotion

		s.redisClient.Set(ctx, redisKey, emotion, 5*time.Minute)
		logger.AtLog.Infof("Cache missed %v", emotion)
	}

	var emotionInt int
	if i, err := strconv.Atoi(emotion); err == nil {
		emotionInt = i
	}

	// find in mongo
	var renderedNftBson bson.M
	err := s.renderedNftRepository.FindOne(context.Background(), map[string]interface{}{
		"chainId":          req.ChainId,
		"contractAddress":  template.NftInfo.ContractAddress,
		"projectId":        req.ProjectId,
		"tokenId":          req.TokenId,
		"metadata.emotion": emotionInt,
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
		Image:       "https://cdn.rove.to/events/world-cup/football-avatar-thumbnail.gif",
		Description: pointerutil.MakePointer("Build your World Cup dream team with PLAYΞRS, a limited collection of 10,000 living 3D generative avatars. Every PLAYΞRS is algorithmically generated from 16 different attributes, with expressions that react to their team’s success on the field. Start your collection today."),
	}, nil
}
