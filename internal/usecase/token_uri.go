package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/davecgh/go-spew/spew"
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

func (u Usecase) GetToken(rootSpan opentracing.Span,  req structure.GetTokenMessageReq) (*entity.TokenUri, error) {
	span, log := u.StartSpan("GetToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	
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

	// tokenUri, err := u.getTokenInfo(span, req)
	// if err != nil {
	// 	log.Error("u.getTokenInfo", err.Error(), err)
	// 	return nil, err
	// }

	isUpdate := false
	if tokenUri.ParsedAttributes == nil  {
		isUpdate = true
		cctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		
		traits := make(map[string]interface{})
		err = chromedp.Run(cctx,
			chromedp.Navigate(tokenUri.AnimationURL),
			chromedp.EvaluateAsDevTools("window.$generativeTraits",&traits),
		)

		// if err != nil {
		// 	log.Error("chromedp.Run", err.Error(), err)
		// 	return nil, err
		// }

		attrs := []entity.TokenUriAttr{}
		for key, item := range traits {
			attr := entity.TokenUriAttr{}
			attr.TraitType = key
			attr.Value = item
			
			attrs = append(attrs, attr)
		}
		tokenUri.ParsedAttributes = attrs
	}
	
	if tokenUri.ParsedImage == nil  {
		isUpdate = true
		var buf []byte
		cctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
	
		err = chromedp.Run(cctx,
			chromedp.Navigate(tokenUri.AnimationURL),
			chromedp.CaptureScreenshot(&buf),
		)

		image := helpers.Base64Eecode(buf)
		image = fmt.Sprintf("%s,%s","data:image/png;base64",image)
		// if err != nil {
		// 	log.Error("chromedp.ParsedImage.Run", err.Error(), err)
		// 	return nil, err
		// }

		tokenUri.ParsedImage = &image
		
	}

	if isUpdate {
		updated, err := u.Repo.UpdateTokenByID(tokenUri.UUID, tokenUri)
		if err != nil {
			log.Error("u.Repo.UpdateOne", err.Error(), err)
			return nil, err
		}
		log.SetData("updated", updated)
	}
	
	return tokenUri, nil
}

func (u Usecase) GetTokenTraits(rootSpan opentracing.Span,  req structure.GetTokenMessageReq) (*entity.TokenUri, error) {
	span, log := u.StartSpan("GetTokenTraits", rootSpan)
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

	if tokenUri.ParsedAttributes == nil  {
		cctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		
		traits := make(map[string]interface{})
		err = chromedp.Run(cctx,
			chromedp.Navigate(tokenUri.AnimationURL),
			chromedp.EvaluateAsDevTools("window.$generativeTraits",&traits),
		)

		if err != nil {
			log.Error("chromedp.Run", err.Error(), err)
			return nil, err
		}

		attrs := []entity.TokenUriAttr{}
		for key, item := range traits {
			attr := entity.TokenUriAttr{}
			attr.TraitType = key
			attr.Value = item
			
			attrs = append(attrs, attr)
		}
		tokenUri.ParsedAttributes = attrs

		updated, err := u.Repo.UpdateTokenByID(tokenUri.UUID, tokenUri)
		if err != nil {
			log.Error("u.Repo.UpdateOne", err.Error(), err)
			return nil, err
		}
		log.SetData("updated", updated)

		return tokenUri, nil
	}
	
	return tokenUri, nil
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

	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(req.TokenID, 10)
	if !ok {
		err := errors.New("Cannot convert tokenID")
		return nil, err
	}
	projectID := new(big.Int).Div(tokenID, big.NewInt(1000000))
	nftProjectDetail, err := u.getNftContractDetail(client, addr, *projectID)
	if err != nil {
		log.Error("u.getNftContractDetail", err.Error(), err)
		return nil, err
	}

	nftProject := nftProjectDetail.ProjectDetail
	parentAddr := nftProject.GenNFTAddr
	spew.Dump(parentAddr.String())
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

func (u Usecase) getNftContractDetail(client *ethclient.Client, contractAddr common.Address, projectID big.Int) (*structure.ProjectDetail, error) {
	gProject, err := generative_project_contract.NewGenerativeProjectContract(contractAddr, client)
	if err != nil {
		return nil, err
	}

	pDchan := make(chan structure.ProjectDetailChan, 1)
	pStatuschan := make(chan structure.ProjectStatusChan, 1)
	pTokenURIchan := make(chan structure.ProjectNftTokenUriChan, 1)

	go func(pDchan chan structure.ProjectDetailChan, projectID *big.Int) {
		proDetail := &generative_project_contract.NFTProjectProject{}
		var err error

		defer func ()  {
			pDchan <- structure.ProjectDetailChan{
				ProjectDetail: proDetail,
				Err:  err,
			}
		}()

		proDetailReps, err := gProject.ProjectDetails(nil,  projectID)
		if err != nil {
			return 
		}

		proDetail = &proDetailReps

	}(pDchan, &projectID)

	go func(pDchan chan structure.ProjectStatusChan, projectID *big.Int) {
		var status *bool
		var err error

		defer func ()  {
			pDchan <- structure.ProjectStatusChan{
				Status: status,
				Err:  err,
			}
		}()

		pStatus, err := gProject.ProjectStatus(nil,  projectID)
		if err != nil {
			return 
		}

		status = &pStatus

	}(pStatuschan, &projectID)

	go func(pDchan chan structure.ProjectNftTokenUriChan, projectID *big.Int) {
		var tokenURI *string
		var err error

		defer func ()  {
			pDchan <- structure.ProjectNftTokenUriChan{
				TokenURI: tokenURI,
				Err:  err,
			}
		}()

		pTokenUri, err := gProject.TokenURI(nil,  projectID)
		if err != nil {
			return 
		}

		tokenURI = &pTokenUri

	}(pTokenURIchan, &projectID)


	detailFromChain := <-  pDchan
	statusFromChain := <-  pStatuschan
	tokenFromChain := <-  pTokenURIchan

	if detailFromChain.Err != nil {
		return nil, detailFromChain.Err
	}
	
	if statusFromChain.Err != nil {
		return nil, statusFromChain.Err
	}
	
	if tokenFromChain.Err != nil {
		return nil, tokenFromChain.Err
	}

	resp := &structure.ProjectDetail{
		ProjectDetail: detailFromChain.ProjectDetail,
		Status: *statusFromChain.Status,
		NftTokenUri: *tokenFromChain.TokenURI,
	}
		
	return resp, nil
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

func (u Usecase) GetProjectDetail(rootSpan opentracing.Span,  req structure.GetProjectDetailMessageReq) (*structure.ProjectDetail, error) {
	span, log := u.StartSpan(fmt.Sprintf("GetProjectDetail.%s.%s", req.ContractAddress, req.ProjectID), rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	contractDataKey := fmt.Sprintf("detail.%s.%s", req.ContractAddress, req.ProjectID)
	
	data, err := u.Cache.GetData(contractDataKey)
	if err != nil {
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

		projectID := new(big.Int)
		projectID, ok := projectID.SetString(req.ProjectID, 10)
		if !ok {
			err := errors.New("Cannot convert tokenID")
			return nil, err
		}
		contractDetail, err := u.getNftContractDetail(client, addr, *projectID)
		if err != nil {
			log.Error("u.getNftContractDetail", err.Error(), err)
			return nil, err
		}
		log.SetData("contractDetail", contractDetail)
		u.Cache.SetData(contractDataKey, contractDetail)
		return contractDetail, nil
	}

	bytes := []byte(*data)
	contractDetail := &structure.ProjectDetail{}
	err = json.Unmarshal(bytes, contractDetail)
	if err != nil {
		log.Error("json.Unmarshal", err.Error(), err)
		return nil, err
	}
	log.SetData("cached.ContractDetail", contractDetail)
	return contractDetail, nil
}