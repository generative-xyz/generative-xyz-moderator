package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/copier"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"rederinghub.io/api"
	"rederinghub.io/internal/adapter"
	"rederinghub.io/internal/dto"
	"rederinghub.io/pkg/config"
	"rederinghub.io/pkg/contracts/generative_boilerplate"
	"rederinghub.io/pkg/logger"
	"rederinghub.io/pkg/utils/constants/contract"
)

const (
	NftInfo = "nftInfo"
)

func (s *service) GetTemplate(ctx context.Context, req *api.GetTemplateRequest) (*api.GetTemplateResponse, error) {
	appConfig := config.AppConfig()
	moralisResp, err := s.moralisAdapter.ListNFTs(appConfig.GenerativeBoilerplateContract, appConfig.ChainID)
	if err != nil {
		log.Errorf("moralis get nft error", err)
		return nil, err
	}
	resp := api.GetTemplateResponse{
		Template: make([]*api.Template, 0, len(moralisResp.Result)),
	}
	for _, nft := range moralisResp.Result {
		var (
			image    = ""
			metadata = make(map[string]interface{})
		)
		if err := json.Unmarshal([]byte(nft.Metadata), &metadata); err == nil {
			if v, ok := metadata["image"]; ok {
				image = v.(string)
			}
		}
		resp.Template = append(resp.Template, &api.Template{
			Name:          nft.Name,
			TokenId:       nft.TokenID,
			Symbol:        nft.Symbol,
			MetadataImage: image,
		})
	}
	resp.Total = int32(len(moralisResp.Result))

	return &resp, nil
}

func (s *service) GetTemplateDetail(ctx context.Context, req *api.GetTemplateDetailRequest) (*api.GetTemplateDetailResponse, error) {
	chainURL, ok := GetRPCURLFromChainID(req.ChainId)
	if !ok {
		return nil, errors.New("missing config chain_config from server")
	}

	var templateDTOFromMongo bson.M
	if err := s.templateRepository.FindOne(context.Background(), map[string]interface{}{
		NftInfo: dto.NftInfo{
			NetworkType:     int(contract.EVM_NetworkType),
			ChainId:         req.ChainId,
			TokenId:         req.TokenId,
			ContractAddress: req.ContractAddress,
		},
	}, &templateDTOFromMongo); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
	} else {
		bytes, _ := json.Marshal(templateDTOFromMongo)
		var templateDTO dto.TemplateDTO
		if err := json.Unmarshal(bytes, &templateDTO); err != nil {
			return nil, err
		}

		return templateDTO.ToProto(), nil
	}

	// Get data from blockchain if not exist
	client, err := ethclient.Dial(chainURL)
	addr := common.HexToAddress(req.ContractAddress)

	instance, err := generative_boilerplate.NewGenerativeBoilerplate(addr, client)
	if err != nil {
		return nil, err
	}

	tokenID, ok := new(big.Int).SetString(req.TokenId, 10)
	if !ok {
		return nil, errors.New("invalid token id")
	}

	resp, err := instance.Projects(&bind.CallOpts{Context: context.Background(), Pending: false}, tokenID)
	if err != nil {
		return nil, err
	}

	if resp.Script == "" && resp.ScriptType == 0 && resp.ProjectName == "" {
		return nil, status.Errorf(codes.NotFound, "token_id %v is not found", req.TokenId)
	}

	var templateDTO = dto.TemplateDTO{NftInfo: dto.NftInfo{
		NetworkType:     int(contract.EVM_NetworkType),
		ChainId:         req.ChainId,
		TokenId:         req.TokenId,
		ContractAddress: req.ContractAddress,
	}}
	if err := copier.Copy(&templateDTO, &resp); err != nil {
		return nil, err
	}

	var templateModel bson.M
	if _bytes, err := json.Marshal(&templateDTO); err != nil {
		return nil, err
	} else {
		if err = json.Unmarshal(_bytes, &templateModel); err != nil {
			return nil, err
		}
	}

	_, err = s.templateRepository.Create(context.Background(), &templateModel)
	if err != nil {
		logger.AtLog.Errorf("[TemplateDetail] Create error %v", err)
		return nil, err
	}

	protoResp := templateDTO.ToProto()
	protoResp.NftInfo = &api.NftInfo{
		NetworkType:     int32(contract.EVM_NetworkType),
		ChainId:         req.ChainId,
		TokenId:         req.TokenId,
		ContractAddress: req.ContractAddress,
	}

	return templateDTO.ToProto(), nil
}

func (s *service) TemplateRendering(_ctx context.Context, request *api.TemplateRenderingRequest) (*api.TemplateRenderingResponse, error) {
	ctx, cancel := context.WithTimeout(_ctx, 30*time.Minute)
	defer cancel()
	var (
		templateDTOFromMongo bson.M
		templateDTO          dto.TemplateDTO
	)
	if err := s.templateRepository.FindOne(context.Background(), map[string]interface{}{
		NftInfo: dto.NftInfo{
			NetworkType:     int(contract.EVM_NetworkType),
			ChainId:         request.ChainId,
			TokenId:         request.TokenId,
			ContractAddress: request.ContractAddress,
		},
	}, &templateDTOFromMongo); err != nil {
		return nil, err
	}

	bytes, _ := json.Marshal(templateDTOFromMongo)
	if err := json.Unmarshal(bytes, &templateDTO); err != nil {
		return nil, err
	}

	resp, err := s.renderMachineAdapter.Render(ctx, &adapter.RenderRequest{
		Script: templateDTO.Script,
		Params: func(_request *api.TemplateRenderingRequest) []string {
			params := make([]string, 0, len(request.Params.Params))
			for _, item := range request.Params.Params {
				params = append(params, item.Value)
			}

			return params
		}(request),
		Seed: request.Params.Seed,
	})

	if err != nil {
		return nil, err
	}

	_resp := &api.TemplateRenderingResponse{
		Glb:   fmt.Sprintf("ipfs://%v", resp.Glb),
		Image: fmt.Sprintf("ipfs://%v", resp.Image),
		//Video: resp.Video,
	}

	if resp.Video != "" {
		_resp.Video = fmt.Sprintf("ipfs://%v", resp.Video)
	}

	return _resp, nil
}
