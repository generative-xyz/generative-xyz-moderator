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
		ExternalLink:    utils.MakeStringPointer("rove.to"),
		Attributes:      protoAttributes,
	}, nil
}

func (s *service) GetRenderedNft(ctx context.Context, req *api.GetRenderedNftRequest) (*api.GetRenderedNftResponse, error) {
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
		json.Unmarshal(doc, &template)
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
			json.Unmarshal(doc, &renderedNft)
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
		return nil, errors.New("invalid token id")
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
		fmt.Sprintf("https://ipfs.rove.to/ipfs/%v", rendered.Image),
		fmt.Sprintf("https://ipfs.rove.to/ipfs/%v", rendered.Glb),
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
