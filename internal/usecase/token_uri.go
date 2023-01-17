package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
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
			log.Error("u.getTokenInfo.mongo.ErrNoDocuments", err.Error(), err)
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

		log.SetData("u.getTokenInfo.tokenUri.ProjectID.Token.Empty", "true")
	}

	isUpdate := false
	project, err := u.Repo.FindProjectBy(contractAddress, tokenUri.ProjectID)
	if err != nil {
		log.Error("u.GetProjectDetail", err.Error(), err)
		return nil, err
	}

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
			log.Error("u.GetUserProfileByWalletAddress",err.Error(), err)
			return
		}
	}

	creatorProfileChan := make(chan structure.ProfileChan, 1) 
	// find owner address on moralis
	nft, err := u.MoralisNft.GetNftByContractAndTokenIDNoCahe(project.GenNFTAddr, tokenID)
	if err == nil {
		if tokenUri.OwnerAddr == "" || strings.ToLower(nft.Owner)  != strings.ToLower(tokenUri.OwnerAddr)  {
			tokenUri.OwnerAddr = strings.ToLower(nft.Owner)
			isUpdate = true
			log.SetData("tokenUri.OwnerAddr.Updated", tokenUri.OwnerAddr)
		}

		ownerProfileChan := make(chan structure.ProfileChan, 1) 
		go getProfile(span, ownerProfileChan, nft.Owner)
		ownerProfileResp := <-ownerProfileChan
		
		log.SetData("ownerProfileResp", ownerProfileResp)
		if tokenUri.Owner == nil || !reflect.DeepEqual(ownerProfileResp.Data, tokenUri.Owner)  {
			isUpdate = true
			tokenUri.Owner = ownerProfileResp.Data
			log.SetData("tokenUri.GenNFTAddr.Owner", tokenUri.Owner)
		}
		//return nil, err
	}else{
		log.Error(" u.MoralisNft.GetNftByContractAndTokenIDNoCahe", err.Error(), err)
	}

	log.SetData("nft", nft)
	go getProfile(span, creatorProfileChan, project.CreatorAddrr)
	
	creatorProfileResp := <-creatorProfileChan
	tokenIDInt, _ := strconv.Atoi(tokenID)
	log.SetData("creatorProfileResp", creatorProfileResp)

	tokenUri.TokenIDInt = tokenIDInt
	if tokenUri.ParsedAttributes == nil {
		isUpdate = true
		cctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		traits := make(map[string]interface{})
		err = chromedp.Run(cctx,
			chromedp.Navigate(tokenUri.AnimationURL),
			chromedp.EvaluateAsDevTools("window.$generativeTraits", &traits),
		)

		log.SetData("traits",traits)
		if err == nil {
			attrs := []entity.TokenUriAttr{}
			for key, item := range traits {
				attr := entity.TokenUriAttr{}
				attr.TraitType = key
				attr.Value = item
	
				attrs = append(attrs, attr)
			}
	
			tokenUri.ParsedAttributes = attrs
			log.SetData("tokenUri.ParsedAttributes.Updated", tokenUri.ParsedAttributes)
		}else{
			log.Error("chromedp.Run.err.generativeTraits",err.Error(), err)
		}
		
	}

	//if true {
	if tokenUri.ParsedImage == nil {
		isUpdate = true
		var buf []byte
		cctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		err = chromedp.Run(cctx,
			chromedp.EmulateViewport(960, 960),
			chromedp.Navigate(tokenUri.AnimationURL),
			chromedp.Sleep(time.Second*time.Duration(captureTimeout)),
			chromedp.CaptureScreenshot(&buf),
		)

		if err == nil {
			image := helpers.Base64Encode(buf)
			image = fmt.Sprintf("%s,%s", "data:image/png;base64", image)
			tokenUri.ParsedImage = &image
			log.SetData("tokenUri.ParsedImage.Updated", "true")
		}else{
			log.Error("chromedp.Run.err", err.Error(), err)
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
				log.Error(" u.GetNftMintedTime", err.Error(), err)
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
		log.SetData("tokenUri.GenNFTAddr.Updated", tokenUri.GenNFTAddr)
	}
	
	if tokenUri.Creator == nil  {
		isUpdate = true
		tokenUri.Creator = creatorProfileResp.Data
		log.SetData("tokenUri.GenNFTAddr.Creator", tokenUri.Creator)
	}
	
	if tokenUri.Project == nil || tokenUri.Project.Stats != project.Stats  {
		isUpdate = true
		tokenUri.Project = project
		log.SetData("tokenUri.GenNFTAddr.project", project.GenNFTAddr)
	}

	if  tokenUri.Thumbnail == "" {
		if  tokenUri.ParsedImage != nil {
			base64Image := *tokenUri.ParsedImage
			i := strings.Index(base64Image, ",")
			if i >= 0 {
				name := fmt.Sprintf("thumb/%s-%s.png", tokenUri.ContractAddress, tokenUri.TokenID)
				base64Image = base64Image[i+1:]
				uploaded, err := u.GCS.UploadBaseToBucket(base64Image,  name)
				if err != nil {
					
					log.Error("u.GCS.UploadBaseToBucket", err.Error(), err)
				}else{
					log.SetData("uploaded", uploaded)
					tokenUri.Thumbnail = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), name)
					isUpdate = true
				}
			}
			// pass reader to NewDecoder
		}
	}
	
	log.SetData("isUpdate", isUpdate)
	if tokenUri.Priority ==  nil {
		priority  := 0
		tokenUri.Priority = &priority
		isUpdate =  true
	}
	
	if tokenUri.TokenIDMini ==  nil {
		tokIdMini  := tokenUri.TokenIDInt %  100000
		tokenUri.TokenIDMini = &tokIdMini
		isUpdate =  true
	}

	if isUpdate {
		updated, err := u.Repo.UpdateOrInsertTokenUri(contractAddress, tokenID, tokenUri)
		if err != nil {
			log.Error("u.Repo.UpdateOne", err.Error(), err)
			return nil, err
		}
		log.SetData("updated", updated)
	}

	return tokenUri, nil
}

func (u Usecase) GetToken(rootSpan opentracing.Span, req structure.GetTokenMessageReq, captureTimeout int) (*entity.TokenUri, error) {
	span, log := u.StartSpan("GetToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("req", req)

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
			log.SetData("live.tokenUri", token.TokenID)
			log.SetData("tokenID", token.TokenID)
			return token, nil
		}else{
			return nil, err
		}
	}

	///log.SetData("tokenUri", tokenUri)
	log.SetData("tokenID", tokenUri.TokenID)
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
	log.SetData("tokenUriData", tokenUriData)
	
	base64Str := strings.ReplaceAll(*tokenUriData, "data:application/json;base64,", "")
	data, err := helpers.Base64Decode(base64Str)
	if err != nil {
		log.Error("helpers.Base64Decode", err.Error(), err)
		return nil, err
	}

	stringData := string(data)
	stringData = strings.ReplaceAll(stringData, "\n", "\\n")
	stringData = strings.ReplaceAll(stringData, "\b", "\\b")
	stringData = strings.ReplaceAll(stringData, "\f", "\\f")
	stringData = strings.ReplaceAll(stringData, "\r", "\\r")
	stringData = strings.ReplaceAll(stringData, "\t", "\\t")

	log.SetData("base64Str", base64Str)
	log.SetData("stringData", stringData)

	dataObject := &entity.TokenUri{}
	err = json.Unmarshal([]byte(stringData), dataObject)
	if err != nil {
		log.Error("json.Unmarshal", err.Error(), err)
		return nil, err
	}

	dataObject.ContractAddress = strings.ToLower(req.ContractAddress)
	dataObject.CreatorAddr = strings.ToLower(nftProject.Creator)
	dataObject.OwnerAddr = strings.ToLower(nftProject.Creator)

	dataObject.TokenID = req.TokenID
	dataObject.ProjectID = projectID.String()
	dataObject.ProjectIDInt = projectID.Int64()

	log.SetData("dataObject.ContractAddress", dataObject.ContractAddress)
	log.SetData("dataObject.Creator", dataObject.Creator)
	log.SetData("dataObject.OwnerAddr", dataObject.OwnerAddr)
	log.SetData("dataObject.TokenID", dataObject.TokenID)
	log.SetData("dataObject.ProjectID", dataObject.ProjectID)
	
	log.SetTag("contractAddress", dataObject.ContractAddress)
	log.SetTag("creator", dataObject.Creator)
	log.SetTag("ownerAddr", dataObject.OwnerAddr)
	log.SetTag("tokenID", dataObject.TokenID)
	log.SetTag("projectID", dataObject.ProjectID)
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

	log.SetData("tokens.Count", len(tokens))
	for _, token := range tokens {
		span, log := u.StartSpan("UpdateTokensFromChain.loop", span)

		log.SetTag("tokenID", token.TokenID)
		log.SetTag("genNFTAddr", token.GenNFTAddr)
		
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
	span, log := u.StartSpan("FilterTokens", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	//TODO use redis schedule instead of crontab or routine to get data.
	if filter.GenNFTAddr != nil {
		defer func ()  {
			go func (rootSpan opentracing.Span, genNftAddress string) {
				span, log := u.StartSpan("GetTokensByContract.Live.Process", rootSpan)
				defer u.Tracer.FinishSpan(span, log)
				
				u.GetTokensByContract(span, genNftAddress, nfts.MoralisFilter{})
		
			}(span, *filter.GenNFTAddr)
		}()
	}
	
	pe := &entity.FilterTokenUris{}

	log.SetData("log", log)
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
	log.SetData("tokens", tokens.Total)
	return tokens, nil
}

func (u Usecase) UpdateToken(rootSpan opentracing.Span, req structure.UpdateTokenReq) (*entity.TokenUri, error) {
	span, log := u.StartSpan("UpdateToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	p, err := u.Repo.FindTokenBy(req.ContracAddress, req.TokenID)
	if err != nil {
		log.Error("UpdateProject.FindTokenBy", err.Error(), err)
		return nil, err
	}

	if req.Priority != nil {
		p.Priority = req.Priority
	}
	
	updated, err := u.Repo.UpdateOrInsertTokenUri(req.ContracAddress, req.TokenID, p)
	if err != nil {
		log.Error("UpdateProject.UpdateOrInsertTokenUri", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return p, nil
}