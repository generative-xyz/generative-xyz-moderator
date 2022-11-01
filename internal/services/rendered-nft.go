package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
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
	"rederinghub.io/pkg/contracts/candy_contract"
	"rederinghub.io/pkg/contracts/generative_param"
	"rederinghub.io/pkg/logger"
	"rederinghub.io/pkg/utils"
)

func GetRenderedNft(
	chainID string,
	contractAddress string,
	projectID string,
	tokenID string,
	template dto.TemplateDTO,
	attributes []generative_param.TraitInfoTrait,
	image string,
	glb string,
) (*model.RenderedNft, error) {
	protoAttributes := make([]*model.OpenSeaAttribute, len(attributes))
	for i := 0; i < len(attributes); i++ {
		var value string
		if len(attributes[i].AvailableValues) > 0 {
			id, err := strconv.ParseInt(attributes[i].Value.String(), 10, 32)
			if err != nil {
				return nil, err
			}
			value = attributes[i].AvailableValues[id]
		} else {
			value = attributes[i].Value.String()
		}
		protoAttributes[i] = &model.OpenSeaAttribute{
			TraitType: attributes[i].Name,
			Value:     value,
		}
	}

	return &model.RenderedNft{
		ChainID:         chainID,
		ContractAddress: contractAddress,
		ProjectID:       projectID,
		TokenID:         tokenID,
		Name:            fmt.Sprintf("%s #%s", template.ProjectName, tokenID),
		Image:           &image,
		Glb:             &glb,
		ExternalLink:    utils.MakeStringPointer("https://rove.to"),
		Attributes:      protoAttributes,
		Description:     template.Description,
	}, nil
}

func (s *service) GetRenderedNft(ctx context.Context, req *api.GetRenderedNftRequest) (*api.GetRenderedNftResponse, error) {
	logger.AtLog.Infof("Handle [GetRenderedNft] %s %s %s %s", req.ChainId, req.ContractAddress, req.ProjectId, req.TokenId)

	// lowercase contract address
	req.ContractAddress = strings.ToLower(req.ContractAddress)

	var templateDTOFromMongo bson.M
	if err := s.templateRepository.FindOne(ctx, map[string]interface{}{
		"nftInfo.tokenId":         req.ProjectId,
		"nftInfo.chainId":         req.ChainId,
		"nftInfo.contractAddress": req.ContractAddress,
	}, &templateDTOFromMongo); err != nil {
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
		"contractAddress": req.ContractAddress,
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

		return renderedNft.ToProto(), nil
	}

	return nil, TemplateNotFoundError{
		TokenID: req.TokenId,
		ChainID: req.ChainId,
	}
}

func (s *service) GetCandyMetadataPost(ctx context.Context, req *api.GetCandyMetadataRequest) (*api.GetCandyMetadataResponse, error) {

	//req.ContractAddress = ""
	req.ProjectId = "291199"

	logger.AtLog.Infof("Handle [GetCandyMetadataPost] %s %s %s %s", req.ChainId, req.ContractAddress, req.ProjectId, req.TokenId)

	chainURL, ok := GetRPCURLFromChainID(req.ChainId)
	if !ok {
		return nil, errors.New("missing config chain_config from server")
	}

	var templateDTOFromMongo bson.M
	if err := s.templateRepository.FindOne(context.Background(), map[string]interface{}{
		"nftInfo.tokenId": req.ProjectId,
		"nftInfo.chainId": req.ChainId,
	}, &templateDTOFromMongo); err != nil {
		return nil, err
	}

	var template dto.TemplateDTO
	{
		doc, err := json.Marshal(templateDTOFromMongo)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(doc, &template)
	}

	// find in mongo
	var renderedNftBson bson.M
	err := s.renderedNftRepository.FindOne(context.Background(), map[string]interface{}{
		"chainId": req.ChainId,
		//"contractAddress": template.NftInfo.ContractAddress,
		"projectId": req.ProjectId,
		"tokenId":   req.TokenId,
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
			json.Unmarshal(doc, &renderedNft)
		}

		return renderedNft.ToCandyResponse(), nil
	}

	client, err := ethclient.Dial(chainURL)

	if err != nil {
		return nil, err
	}
	addr := common.HexToAddress(template.MinterNFTInfo.Hex())

	instance, err := candy_contract.NewCandyContract(addr, client)

	if err != nil {
		return nil, err
	}

	tokenID, ok := new(big.Int).SetString(req.TokenId, 10)
	if !ok {
		return nil, InvalidTokenIDError{TokenID: req.TokenId}
	}

	colorPallete, shape, size, surface, err := instance.GetParamValues(&bind.CallOpts{Context: context.Background(), Pending: false}, tokenID)
	if err != nil {
		return nil, err
	}

	params := make([]string, 0, 7)
	for i := 0; i < 4; i++ {
		params = append(params, colorPallete[i])
	}
	params = append(params, shape)
	params = append(params, size)
	params = append(params, surface)

	// call to render machine
	rendered, err := s.renderMachineAdapter.Render(ctx, &adapter.RenderRequest{
		Script: template.Script,
		Params: params,
		Seed:   "1",
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
	renderedNft.Attributes = []*model.OpenSeaAttribute{
		{
			TraitType: "Color 1",
			Value:     colorPallete[0],
		},
		{
			TraitType: "Color 2",
			Value:     colorPallete[1],
		},
		{
			TraitType: "Color 3",
			Value:     colorPallete[2],
		},
		{
			TraitType: "Color 4",
			Value:     colorPallete[3],
		},
		{
			TraitType: "Shape",
			Value:     shape,
		},
		{
			TraitType: "Size",
			Value:     size,
		},
		{
			TraitType: "Surface",
			Value:     surface,
		},
	}
	if err != nil {
		return nil, err
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
		logger.AtLog.Errorf("[GetRenderedNft] Create error %v", err)
		return nil, err
	}

	logger.AtLog.Infof("Done [GenerateCandy] #%s", req.TokenId)

	return renderedNft.ToCandyResponse(), nil
}

func (s *service) GetCandyMetadata(ctx context.Context, req *api.GetCandyMetadataRequest) (*api.GetCandyMetadataResponse, error) {
	//req.ContractAddress = ""
	req.ProjectId = "291199"

	logger.AtLog.Infof("Handle [GetCandyMetadata] %s %s %s %s", req.ChainId, req.ContractAddress, req.ProjectId, req.TokenId)
	var templateDTOFromMongo bson.M
	if err := s.templateRepository.FindOne(context.Background(), map[string]interface{}{
		"nftInfo.tokenId": req.ProjectId,
		"nftInfo.chainId": req.ChainId,
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
		"chainId": req.ChainId,
		//"contractAddress": template.NftInfo.ContractAddress,
		"projectId": req.ProjectId,
		"tokenId":   req.TokenId,
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

		return renderedNft.ToCandyResponse(), nil
	}

	return nil, TemplateNotFoundError{TokenID: req.ProjectId, ChainID: req.ChainId}
}

func (s *service) GetRenderedNftPost(ctx context.Context, req *api.GetRenderedNftRequest) (*api.GetRenderedNftResponse, error) {
	logger.AtLog.Infof("Handle [GetRenderedNftPost] %s %s %s %s", req.ChainId, req.ContractAddress, req.ProjectId, req.TokenId)

	// lowercase contract address
	req.ContractAddress = strings.ToLower(req.ContractAddress)

	chainURL, ok := GetRPCURLFromChainID(req.ChainId)
	if !ok {
		return nil, errors.New("missing config chain_config from server")
	}

	var templateDTOFromMongo bson.M
	if err := s.templateRepository.FindOne(context.Background(), map[string]interface{}{
		"nftInfo.tokenId":         req.ProjectId,
		"nftInfo.chainId":         req.ChainId,
		"nftInfo.contractAddress": req.ContractAddress,
	}, &templateDTOFromMongo); err != nil {
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
		"contractAddress": req.ContractAddress,
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

		return renderedNft.ToProto(), nil
	}

	client, err := ethclient.Dial(chainURL)

	if err != nil {
		return nil, err
	}
	addr := common.HexToAddress(template.MinterNFTInfo.Hex())

	instance, err := generative_param.NewGenerativeParam(addr, client)
	if err != nil {
		return nil, err
	}

	tokenID, ok := new(big.Int).SetString(req.TokenId, 10)
	if !ok {
		return nil, InvalidTokenIDError{TokenID: req.TokenId}
	}

	seed, resp, err := instance.GetParamValues(&bind.CallOpts{Context: context.Background(), Pending: false}, tokenID)
	if err != nil {
		return nil, err
	}

	attributes, err := instance.GetTokenTraits(&bind.CallOpts{Context: context.Background(), Pending: false}, tokenID)
	if err != nil {
		return nil, err
	}

	params := make([]string, len(resp))
	for i := 0; i < len(resp); i++ {
		if len(template.ParamsTemplate.Params[i].AvailableValues) > 0 {
			id, err := strconv.ParseInt(resp[i].Value.String(), 10, 32)
			if err != nil {
				return nil, err
			}
			params[i] = template.ParamsTemplate.Params[i].AvailableValues[id]
		} else {
			params[i] = resp[i].Value.String()
		}
	}

	// call to render machine
	rendered, err := s.renderMachineAdapter.Render(ctx, &adapter.RenderRequest{
		Script: template.Script,
		Params: params,
		Seed:   string(seed[:]),
	})

	if err != nil {
		return nil, err
	}

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

	return renderedNft.ToProto(), nil
}
