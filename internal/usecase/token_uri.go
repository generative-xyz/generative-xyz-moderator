package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/copier"
	"github.com/oliamb/cutter"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"

	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/contracts/generative_nft_contract"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) GetLiveToken(rootSpan opentracing.Span, req structure.GetTokenMessageReq, captureTimeout int) (*entity.TokenUri, error) {
	span, log := u.StartSpan("GetLiveToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	log.SetTag("contractAddress",req.ContractAddress)
	log.SetTag("tokenID",req.TokenID)

	contractAddress := strings.ToLower(req.ContractAddress) 
	tokenID := strings.ToLower(req.TokenID)
	
	tokenUri, err := u.Repo.FindTokenBy(contractAddress, tokenID)
	if err != nil {
		log.Error("u.Repo.FindTokenBy", err.Error(), err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.SetData("u.getTokenInfo.mongo.ErrNoDocuments", "true")
			tokenUri, err = u.getTokenInfo(span, req)
			if err != nil {
				log.Error("u.getTokenInfo", err.Error(), err)
				return nil, err
			}

		} else {
			log.Error("u.Repo.FindTokenBy", err.Error(), err)
			return nil, err
		}
	}

	log.SetData("tokenUri", tokenUri)
	if tokenUri.ProjectID == "" {
		tokenUri, err = u.getTokenInfo(span, req)
		if err != nil {
			log.Error("u.getTokenInfo", err.Error(), err)
			return nil, err
		}

		log.SetData("u.getTokenInfo.tokenUri.ProjectID.Token.Empty", "true")
	}
	isUpdate := false

	// find project by projectID and contract address
	project, err := u.GetProjectDetail(span, structure.GetProjectDetailMessageReq{ContractAddress: contractAddress, ProjectID: tokenUri.ProjectID})
	if err != nil {
		log.Error("u.GetProjectDetail", err.Error(), err)
		return nil, err
	}

	// find owner address on moralis
	nft, err := u.MoralisNft.GetNftByContractAndTokenID(project.GenNFTAddr, tokenID)
	if err != nil {
		log.Error(" u.MoralisNft.GetNftByContractAndTokenID", err.Error(), err)
		return nil, err
	}

	if nft.Owner != tokenUri.OwnerAddr {
		tokenUri.OwnerAddr =  nft.Owner
		isUpdate = true
	}

	log.SetData("nft", nft)
	getProfile := func(rootSpan opentracing.Span, c chan structure.ProfileChan, address string) {
		span, log := u.StartSpan("GetToken.Profile", rootSpan)
		defer u.Tracer.FinishSpan(span, log)
		
		var profile *entity.Users
		var err error
		log.SetTag(utils.WALLET_ADDRESS_TAG, address)
		
		defer func() {
			response :=  structure.ProfileChan{
				Data:  profile,
				Err:  err,
			}

			log.SetData("response", response)
			c <- response
		}()

		profile, err = u.GetUserProfileByWalletAddress(span, strings.ToLower(address))
		if err != nil {
			return
		}
	}

	ownerProfileChan := make(chan structure.ProfileChan, 1) 
	creatorProfileChan := make(chan structure.ProfileChan, 1) 

	go getProfile(span, ownerProfileChan, nft.Owner)
	go getProfile(span, creatorProfileChan, project.CreatorAddrr)

	ownerProfileResp := <-ownerProfileChan
	creatorProfileResp := <-creatorProfileChan
	tokenIDInt, _ := strconv.Atoi(tokenID)

	log.SetData("ownerProfileResp", ownerProfileResp)
	log.SetData("creatorProfileResp", creatorProfileResp)

	tokenUri.TokenIDInt = tokenIDInt

	if tokenUri.OwnerAddr == "" {
		tokenUri.OwnerAddr = strings.ToLower(nft.Owner)
		isUpdate = true
	}

	if tokenUri.ParsedAttributes == nil {
		isUpdate = true
		cctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		traits := make(map[string]interface{})
		err = chromedp.Run(cctx,
			chromedp.Navigate(tokenUri.AnimationURL),
			chromedp.EvaluateAsDevTools("window.$generativeTraits", &traits),
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


	//if true {
	if tokenUri.ParsedImage != nil {
		isUpdate = true
		var buf []byte
		cctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		err = chromedp.Run(cctx,
			chromedp.Navigate(tokenUri.AnimationURL),
			chromedp.Sleep(time.Second*time.Duration(captureTimeout)),
			chromedp.CaptureScreenshot(&buf),
		)

		img, _, err := image.Decode(bytes.NewReader(buf))
		if err == nil {
			croppedImg, err := cutter.Crop(img, cutter.Config{
				Width:  960,
				Height: 960,
				Mode: cutter.Centered,
			})
	
			buf1 := new(bytes.Buffer)
			err = png.Encode(buf1, croppedImg)
			
			if err == nil {
				bytesArr := buf1.Bytes()
				image := helpers.Base64Encode(bytesArr)
				image = fmt.Sprintf("%s,%s", "data:image/png;base64", image)
				// if err != nil {
				// 	log.Error("chromedp.ParsedImage.Run", err.Error(), err)
				// 	return nil, err
				// }
				
				tokenUri.ParsedImage = &image
			}else{
				log.Error("image.Decode", err.Error(), err)
				return nil, err
			}
	
		}else{
			log.Error("image.Decode", err.Error(), err)
			//return nil, err
		}

		
	}

	if tokenUri.ProjectID == "" {
		err := func() error {
			tokenID := new(big.Int)
			tokenID, ok := tokenID.SetString(req.TokenID, 10)
			if !ok {
				return errors.New("cannot convert tokenID to big int")
			}
			projectID := new(big.Int).Div(tokenID, big.NewInt(1000000))

			tokenUri.ProjectID = projectID.String()
			isUpdate = true
			return nil
		}()
		if err != nil {
			log.Error("error update token uri project id", err.Error(), err)
		}
	}

	if tokenUri.BlockNumberMinted == nil || tokenUri.MintedTime == nil {
		err := func() error {
			project, err := u.GetProjectDetail(span, structure.GetProjectDetailMessageReq{ContractAddress: req.ContractAddress, ProjectID: tokenUri.ProjectID})
			if err != nil {
				return err
			}

			log.SetData("project", project)

			nftMintedTime, err := u.GetNftMintedTime(span, structure.GetNftMintedTimeReq{
				ContractAddress: project.GenNFTAddr,
				TokenID: req.TokenID,
			})
			if err != nil {
				return err
			}

			tokenUri.BlockNumberMinted = nftMintedTime.BlockNumberMinted
			tokenUri.MintedTime = nftMintedTime.MintedTime
			isUpdate = true
			return nil
		}()
		if err != nil {
			log.Error("error update token uri block number minted", err.Error(), err)
		}
	}
	
	if tokenUri.GenNFTAddr == ""  {
		isUpdate = true
		tokenUri.GenNFTAddr = project.GenNFTAddr
	}
	
	if tokenUri.Owner == nil  {
		isUpdate = true
		tokenUri.Owner = ownerProfileResp.Data
	}
	
	if tokenUri.Creator == nil  {
		isUpdate = true
		tokenUri.Creator = creatorProfileResp.Data
	}
	
	if tokenUri.Project == nil  {
		isUpdate = true
		tokenUri.Project = project
	}

	log.SetData("isUpdate", isUpdate)

	spew.Dump(tokenUri)
	//isUpdate = true
	if isUpdate {
		updated, err := u.Repo.UpdateOrInsertTokenUri(contractAddress, tokenID, tokenUri)
		if err != nil {
			log.Error("u.Repo.UpdateOne", err.Error(), err)
			return nil, err
		}
		log.SetData("updated", updated)
	}

	//log.SetData("tokenUri.Inserted", tokenUri)
	// err = u.Repo.CreateTokenURI(dataObject)
	// if err != nil {
	// 	log.Error("u.Repo.CreateTokenURI", err.Error(), err)
	// 	return nil, err
	// }
	return tokenUri, nil
}

func (u Usecase) GetToken(rootSpan opentracing.Span, req structure.GetTokenMessageReq, captureTimeout int) (*entity.TokenUri, error) {
	span, log := u.StartSpan("GetToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	go u.GetLiveToken(span, req, captureTimeout)

	contractAddress := strings.ToLower(req.ContractAddress) 
	tokenID := strings.ToLower(req.TokenID)

	tokenUri, err := u.Repo.FindTokenBy(contractAddress, tokenID)
	if err != nil {
		log.Error("u.Repo.FindTokenBy", err.Error(), err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			token, err := u.GetLiveToken(span, req, captureTimeout)
			if err != nil { 
				log.Error("u.GetLiveToken", err.Error(), err)
				return nil, err
			}
			log.SetData("live.tokenUri", tokenUri)
			return token, nil
		}
	}

	log.SetData("tokenUri", tokenUri)
	return tokenUri, nil
}

func (u Usecase) getTokenInfo(rootSpan opentracing.Span, req structure.GetTokenMessageReq) (*entity.TokenUri, error) {
	span, log := u.StartSpan("Usecase.getTokenInfo", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("req", req)
	addr := common.HexToAddress(req.ContractAddress)

	// call to contract to get emotion
	client, err := helpers.EthDialer()
	if err != nil {
		log.Error("ethclient.Dial", err.Error(), err)
		return nil, err
	}

	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(req.TokenID, 10)
	if !ok {
		return nil, errors.New("cannot convert tokenID")
	}
	projectID := new(big.Int).Div(tokenID, big.NewInt(1000000))
	nftProjectDetail, err := u.getProjectDetailFromChain(span, structure.GetProjectDetailMessageReq{
		ContractAddress: addr.String(),
		ProjectID: projectID.String(),
	})
	if err != nil {
		log.Error("u.getNftContractDetail", err.Error(), err)
		return nil, err
	}

	nftProject := nftProjectDetail.ProjectDetail
	parentAddr := nftProject.GenNFTAddr
	
	log.SetData("nftProject", nftProject)
	log.SetData("parentAddr", parentAddr)
	
	tokenUriData, err := u.getNftProjectTokenUri(client, parentAddr, req.TokenID)
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

	dataObject.ContractAddress = strings.ToLower(req.ContractAddress)
	dataObject.CreatorAddr = strings.ToLower(nftProject.Creator)
	dataObject.OwnerAddr = strings.ToLower(nftProject.Creator)

	dataObject.TokenID = req.TokenID
	dataObject.ProjectID = projectID.String()

	log.SetData("dataObject", dataObject)
	return dataObject, nil
}

func (u Usecase) getNftProjectTokenUri(client *ethclient.Client, contractAddr common.Address, tokenIDStr string) (*string, error) {
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(tokenIDStr, 10)
	if !ok {
		return nil, errors.New("cannot convert tokenID")
	}

	gNft, err := generative_nft_contract.NewGenerativeNftContract(contractAddr, client)
	if err != nil {
		return nil, err
	}

	value, err := gNft.TokenGenerativeURI(nil, tokenID)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (u Usecase) UpdateTokensFromChain(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("Usecase.UpdateTokensFromChain", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	
	tokens, err := u.Repo.GetAllTokens()
	if err != nil {
		log.Error("GetAllTokens", err.Error(), err)
		return err
	}

	for _, token := range tokens {
		
		_, err := u.GetToken(span, structure.GetTokenMessageReq{ContractAddress: token.ContractAddress, TokenID: token.TokenID}, 5)
		if err != nil {
			log.Error(fmt.Sprintf("u.GetToken_%s_%s", token.ContractAddress,  token.TokenID), err.Error(), err)
			return err
		}
	}

	return nil
}

func (u Usecase) GetTokensByContract(rootSpan opentracing.Span, contractAddress string, filter nfts.MoralisFilter) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetTokensByContract", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	client, err := helpers.EthDialer()
	if err != nil {
		log.Error("ethclient.Dial", err.Error(), err)
		return nil, err
	}

	contractAddr := common.HexToAddress(contractAddress)
	gNft, err := generative_nft_contract.NewGenerativeNftContract(contractAddr, client)
	if err != nil {
		log.Error("generative_nft_contract.NewGenerativeNftContract", err.Error(), err)
		return nil, err
	}

	project, err := gNft.Project(nil)
	if err != nil {
		log.Error("gNft.Project", err.Error(), err)
		return nil, err
	}
	parentAddr := project.ProjectAddr

	resp, err := u.MoralisNft.GetNftByContract(contractAddress, filter)
	if err != nil {
		log.Error("u.MoralisNft.GetNftByContract", err.Error(), err)
		return nil, err
	}
	parentAddrStr := parentAddr.String()
	result := []entity.TokenUri{}
	for _, item := range resp.Result {
		tokenID := item.TokenID
		token, err := u.GetToken(span, structure.GetTokenMessageReq{ContractAddress: parentAddrStr, TokenID: tokenID}, 5)
		if err != nil {
			log.Error("u.getTokenInfo", err.Error(), err)
			return nil, err
		}
		result = append(result, *token)
	}

	p := &entity.Pagination{}
	p.Result = result
	p.Currsor = resp.Cursor
	p.Total = int64(resp.Total)
	p.Page = int64(resp.Page)
	p.PageSize = int64(resp.PageSize)
	return p, nil
}

func (u Usecase) FilterTokens(rootSpan opentracing.Span,  filter structure.FilterTokens) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetTokensByContract", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	
	if filter.GenNFTAddr != nil {
		go func (rootSpan opentracing.Span, genNftAddress string) {
			span, log := u.StartSpan("GetTokensByContract.Live.Process", rootSpan)
			defer u.Tracer.FinishSpan(span, log)
			
			u.GetTokensByContract(span, genNftAddress, nfts.MoralisFilter{})

		}(span, *filter.GenNFTAddr)
	}

	pe := &entity.FilterTokenUris{}
	err := copier.Copy(pe, filter)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	tokens, err := u.Repo.FilterTokenUri(*pe)
	if err != nil {
		log.Error("u.Repo.FilterTokenUri", err.Error(), err)
		return nil, err
	}
	
	return tokens, nil
}
