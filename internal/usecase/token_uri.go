package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"

	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/contracts/generative_nft_contract"
	"rederinghub.io/utils/contracts/generative_project_contract"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) GetToken(rootSpan opentracing.Span, req structure.GetTokenMessageReq, captureTimeout int) (*entity.TokenUri, error) {
	span, log := u.StartSpan("GetToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	contractAddress := req.ContractAddress
	tokenID := req.TokenID
	tokenUri, err := u.Repo.FindTokenBy(req.ContractAddress, tokenID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
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

	if tokenUri.ProjectID == "" {
		tokenUri, err = u.getTokenInfo(span, req)
			if err != nil {
				log.Error("u.getTokenInfo", err.Error(), err)
				return nil, err
			}
	}

	isUpdate := false

	// find project by projectID and contract address
	project, err := u.GetProjectDetail(span, structure.GetProjectDetailMessageReq{ContractAddress: contractAddress, ProjectID: tokenUri.ProjectID})
	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		return nil, err
	}

	// find owner address on moralis
	nft, err := u.MoralisNft.GetNftByContractAndTokenID(project.GenNFTAddr, tokenID)
	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		return nil, err
	}

	getProfile := func(c chan structure.ProfileChan, address string) {
		var profile *entity.Users
		var err error

		defer func() {
			c <- structure.ProfileChan{
				Data:  profile,
				Err:  err,
			}
		}()

		profile, err = u.GetUserProfileByWalletAddress(span, strings.ToLower(address))
		if err != nil {
			return
		}
	}

	ownerProfileChan := make(chan structure.ProfileChan, 1) 
	creatorProfileChan := make(chan structure.ProfileChan, 1) 

	go getProfile(ownerProfileChan, nft.Owner)
	go getProfile(creatorProfileChan, project.CreatorAddrr)

	ownerProfileResp := <-ownerProfileChan
	creatorProfileResp := <-creatorProfileChan
	tokenIDInt, _ := strconv.Atoi(tokenID)

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

	if tokenUri.ParsedImage == nil {
		isUpdate = true
		var buf []byte
		cctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		err = chromedp.Run(cctx,
			chromedp.Navigate(tokenUri.AnimationURL),
			chromedp.Sleep(time.Second*time.Duration(captureTimeout)),
			chromedp.CaptureScreenshot(&buf),
		)

		image := helpers.Base64Encode(buf)
		image = fmt.Sprintf("%s,%s", "data:image/png;base64", image)
		// if err != nil {
		// 	log.Error("chromedp.ParsedImage.Run", err.Error(), err)
		// 	return nil, err
		// }

		tokenUri.ParsedImage = &image
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

	if isUpdate {
		updated, err := u.Repo.UpdateTokenByID(tokenUri.UUID, tokenUri)
		if err != nil {
			log.Error("u.Repo.UpdateOne", err.Error(), err)
			return nil, err
		}
		log.SetData("updated", updated)
	}

	tokenUri.Owner = ownerProfileResp.Data
	tokenUri.Creator = creatorProfileResp.Data
	tokenUri.Project = project
	return tokenUri, nil
}

func (u Usecase) GetTokenTraits(rootSpan opentracing.Span, req structure.GetTokenMessageReq) (*entity.TokenUri, error) {
	span, log := u.StartSpan("GetTokenTraits", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("req", req)
	tokenUri, err := u.Repo.FindTokenBy(req.ContractAddress, req.TokenID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
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

	if tokenUri.ParsedAttributes == nil {
		cctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		traits := make(map[string]interface{})
		err = chromedp.Run(cctx,
			chromedp.Navigate(tokenUri.AnimationURL),
			chromedp.EvaluateAsDevTools("window.$generativeTraits", &traits),
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

func (u Usecase) getTokenInfo(rootSpan opentracing.Span, req structure.GetTokenMessageReq) (*entity.TokenUri, error) {
	span, log := u.StartSpan("UserProfile", rootSpan)
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
	nftProjectDetail, err := u.getNftContractDetail(client, addr, *projectID)
	if err != nil {
		log.Error("u.getNftContractDetail", err.Error(), err)
		return nil, err
	}

	nftProject := nftProjectDetail.ProjectDetail
	parentAddr := nftProject.GenNFTAddr

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
	dataObject.TokenID = req.TokenID
	dataObject.ProjectID = projectID.String()

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
	mintedTimeChan := make (chan structure.NftMintedTimeChan, 1)

	go func(pDchan chan structure.ProjectDetailChan, projectID *big.Int) {
		proDetail := &generative_project_contract.NFTProjectProject{}
		var err error

		defer func() {
			pDchan <- structure.ProjectDetailChan{
				ProjectDetail: proDetail,
				Err:           err,
			}
		}()

		proDetailReps, err := gProject.ProjectDetails(nil, projectID)
		if err != nil {
			return
		}

		proDetail = &proDetailReps

	}(pDchan, &projectID)

	go func(pDchan chan structure.ProjectStatusChan, projectID *big.Int) {
		var status *bool
		var err error

		defer func() {
			pDchan <- structure.ProjectStatusChan{
				Status: status,
				Err:    err,
			}
		}()

		pStatus, err := gProject.ProjectStatus(nil, projectID)
		if err != nil {
			return
		}

		status = &pStatus

	}(pStatuschan, &projectID)

	go func(pDchan chan structure.ProjectNftTokenUriChan, projectID *big.Int) {
		var tokenURI *string
		var err error

		defer func() {
			pDchan <- structure.ProjectNftTokenUriChan{
				TokenURI: tokenURI,
				Err:      err,
			}
		}()

		pTokenUri, err := gProject.TokenURI(nil, projectID)
		if err != nil {
			return
		}

		tokenURI = &pTokenUri

	}(pTokenURIchan, &projectID)

	go func(mintedTimeChan chan structure.NftMintedTimeChan, projectID *big.Int) {
		var mintedTime *structure.NftMintedTime
		var err error
		defer func() {
			mintedTimeChan <- structure.NftMintedTimeChan{
				NftMintedTime: mintedTime,
				Err: err,
			}
		}()
		span, _ := u.StartSpanWithoutRoot("getNftContractDetail.GetNftMintedTime")
		mintedTime, err = u.GetNftMintedTime(span, structure.GetNftMintedTimeReq{
			ContractAddress: contractAddr.String(),
			TokenID: projectID.String(),
		})
	}(mintedTimeChan, &projectID)

	detailFromChain := <-pDchan
	statusFromChain := <-pStatuschan
	tokenFromChain := <-pTokenURIchan
	nftMintedTime := <-mintedTimeChan

	if detailFromChain.Err != nil {
		return nil, detailFromChain.Err
	}

	if statusFromChain.Err != nil {
		return nil, statusFromChain.Err
	}

	if tokenFromChain.Err != nil {
		return nil, tokenFromChain.Err
	}

	if nftMintedTime.Err != nil {
		return nil, nftMintedTime.Err
	}

	gNftProject, err := generative_nft_contract.NewGenerativeNftContract(detailFromChain.ProjectDetail.GenNFTAddr, client)
	if err != nil {
		return nil, err
	}

	//nft project detail chain
	nftProjectDchan := make(chan structure.NftProjectDetailChan, 1)
	go func(nftProjectDchan chan structure.NftProjectDetailChan, gNftProject *generative_nft_contract.GenerativeNftContract) {
		data := &structure.NftProjectDetail{}
		var err error

		defer func() {
			nftProjectDchan <- structure.NftProjectDetailChan{
				Data: data,
				Err:  err,
			}
		}()

		respData, err := gNftProject.Project(nil)
		err = copier.Copy(data, respData)

	}(nftProjectDchan, gNftProject)

	nftRoyaltychan := make(chan structure.RoyaltyChan, 1)
	go func(nftRoyaltychan chan structure.RoyaltyChan, gNftProject *generative_nft_contract.GenerativeNftContract) {
		var data *big.Int
		var err error

		defer func() {
			nftRoyaltychan <- structure.RoyaltyChan{
				Data: data,
				Err:  err,
			}
		}()

		data, err = gNftProject.Royalty(nil)

	}(nftRoyaltychan, gNftProject)

	dataFromNftPChan := <-nftProjectDchan
	dataFromRoyaltyPChan := <-nftRoyaltychan

	resp := &structure.ProjectDetail{
		ProjectDetail: detailFromChain.ProjectDetail,
		Status:        *statusFromChain.Status,
		NftTokenUri:   *tokenFromChain.TokenURI,
		NftMintedTime:  *nftMintedTime.NftMintedTime,
	}

	if dataFromNftPChan.Err == nil && dataFromNftPChan.Data != nil {
		resp.NftProjectDetail = *dataFromNftPChan.Data
	} else {
		resp.NftProjectDetail = structure.NftProjectDetail{}
	}

	if dataFromRoyaltyPChan.Err == nil && dataFromRoyaltyPChan.Data != nil {
		resp.Royalty = structure.ProjectRoyalty{
			Data: *dataFromRoyaltyPChan.Data,
		}
	}

	return resp, nil
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

func (u Usecase) getProjectDetailFromChain(rootSpan opentracing.Span, req structure.GetProjectDetailMessageReq) (*structure.ProjectDetail, error) {
	span, log := u.StartSpan("getProjectDetailFromChain", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	contractDataKey := fmt.Sprintf("detail.%s.%s", req.ContractAddress, req.ProjectID)

	//u.Cache.Delete(contractDataKey)
	data, err := u.Cache.GetData(contractDataKey)
	if err != nil {
		log.SetData("req", req)
		
		addr := common.HexToAddress(req.ContractAddress)
		// call to contract to get emotion
		client, err := helpers.EthDialer()
		if err != nil {
			log.Error("ethclient.Dial", err.Error(), err)
			return nil, err
		}

		projectID := new(big.Int)
		projectID, ok := projectID.SetString(req.ProjectID, 10)
		if !ok {
			return nil, errors.New("cannot convert tokenID")
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

	// call to contract to get emotion
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

func (u Usecase) GetTokensByWalletAddress(rootSpan opentracing.Span, contractAddress string, filter nfts.MoralisFilter) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetTokensByContract", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	// call to contract to get emotion
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
