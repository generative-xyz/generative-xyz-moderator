package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/contracts/generative_nft_contract"
	"rederinghub.io/utils/contracts/generative_project_contract"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) GetToken(rootSpan opentracing.Span,  req structure.GetTokenMessageReq) (*structure.GetTokenMessageResp, error) {
	span, log := u.StartSpan("GetToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("req", req)
	tokenUri, err := u.Repo.FindTokenBy(req.ContractAddress, req.TokenID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			tokenUri, err = u.getTokenInfo(span, req)
			if err != nil {
				log.Error("u.getTokenInfo", err.Error(), err)
				return nil, err
			}

		}else{
			log.Error("u.Repo.FindTokenBy", err.Error(), err)
			return nil, err
		}
	}
	
	cctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res interface{}
	err = chromedp.Run(cctx,
		chromedp.Navigate(tokenUri.AnimationURL),
		chromedp.EvaluateAsDevTools("window.$generativeTraits",&res),
	)

	if err != nil {
		log.Error("chromedp.Run", err.Error(), err)
		return nil, err
	}

	log.SetData("res", res)
	resp := &structure.GetTokenMessageResp{}
	return resp, nil
}

func (u Usecase) getTokenInfo(rootSpan opentracing.Span,  req structure.GetTokenMessageReq) (*entity.TokenUri, error) {
	span, log := u.StartSpan("UserProfile", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	
	log.SetData("req", req)
	chainURL := os.Getenv("CHAIN_URL")
	addr := common.HexToAddress(req.ContractAddress)

	log.SetData("chainURL", chainURL)
	// call to contract to get emotion
	client, err := ethclient.Dial(chainURL)
	if err != nil {
		log.Error("ethclient.Dial", err.Error(), err)
		return nil, err
	}

	nftProjectDetail, err := u.getNftContractDetail(client, addr, req.TokenID)
	if err != nil {
		log.Error("u.getNftContractDetail", err.Error(), err)
		return nil, err
	}

	parentAddr := nftProjectDetail.GenNFTAddr
	tokenUriData, err := u.getNftProjectTokenUri(client, parentAddr,  req.TokenID)
	if err != nil {
		log.Error("u.getNftProjectTokenUri", err.Error(), err)
		return nil, err
	}

	log.SetData("parentAddr", parentAddr)
	
	base64Str := strings.ReplaceAll(*tokenUriData, "data:application/json;base64,", "")
	data, err := helpers.Base64Decode(base64Str)
	if err != nil {
		log.Error("helpers.Base64Decode", err.Error(), err)
		return nil, err
	}

	dataObject := &entity.TokenUri{}
	err = json.Unmarshal(data, dataObject)
	if err != nil {
		log.Error("json.Unmarshal", err.Error(), err)
		return nil, err
	}

	dataObject.ContractAddress = req.ContractAddress
	dataObject.TokenID = req.TokenID
	
	err = u.Repo.CreateTokenURI(dataObject)
	if err != nil {
		log.Error("u.Repo.CreateTokenURI", err.Error(), err)
		return nil, err
	}

	log.SetData("dataObject", dataObject)
	return dataObject, nil
}

func (u Usecase) getNftContractDetail(client *ethclient.Client, contractAddr common.Address, tokenIDStr string) (*generative_project_contract.NFTProjectProject, error) {
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(tokenIDStr, 10)
	if !ok {
		err := errors.New("Cannot convert tokenID")
		return nil, err
	}

	gProject, err := generative_project_contract.NewGenerativeProjectContract(contractAddr, client)
	if err != nil {
		return nil, err
	}

	projectID := new(big.Int).Div(tokenID, big.NewInt(1000000))
	proDetail, err := gProject.ProjectDetails(nil,  projectID)
	if err != nil {
		return nil, err
	}
	return &proDetail, nil
}

func (u Usecase) getNftProjectTokenUri(client *ethclient.Client, contractAddr common.Address, tokenIDStr string) (*string, error) {
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(tokenIDStr, 10)
	if !ok {
		err := errors.New("Cannot convert tokenID")
		return nil, err
	}

	gNft, err := generative_nft_contract.NewGenerativeNftContract(contractAddr, client)
	if err != nil {
		return nil, err
	}

	value, err := gNft.TokenGenerativeURI(nil, tokenID)
	if err !=nil {
		return nil, err
	}

	return &value, nil
}