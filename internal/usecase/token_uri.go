package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
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
	"rederinghub.io/utils/helpers"
)

func (u Usecase) GetLiveToken(rootSpan opentracing.Span, req structure.GetTokenMessageReq, captureTimeout int) (*entity.TokenUri, error) {
	span, log := u.StartSpan("GetLiveToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetTag("tokenID", req.TokenID)
	log.SetTag("contractAddress", req.ContractAddress)
	
	tokenUri, err := u.getTokenInfo(span, req)
	if err != nil {
		log.Error("u.getTokenInfo", err.Error(), err)
		return nil, err
	}

	updated, err := u.Repo.UpdateOrInsertTokenUri(tokenUri.ContractAddress, tokenUri.TokenID, tokenUri)
	if err != nil {
		log.Error("u.Repo.UpdateOne", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return tokenUri, nil
}

func (u Usecase) RunAndCap(rootSpan opentracing.Span, token *entity.TokenUri, captureTimeout int) {
	span, log := u.StartSpan("RunAndCap", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	if token == nil {
		return
	}

	log.SetTag("tokenID", token.TokenID)
	log.SetTag("contractAddress", token.ContractAddress)
	
	var buf []byte
	attrs := []entity.TokenUriAttr{}
	strAttrs := []entity.TokenUriAttrStr{}

	cctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	traits := make(map[string]interface{})
	err := chromedp.Run(cctx,
		chromedp.EmulateViewport(960, 960),
		chromedp.Navigate(token.AnimationURL),
		chromedp.Sleep(time.Second*time.Duration(captureTimeout)),
		chromedp.CaptureScreenshot(&buf),
		chromedp.EvaluateAsDevTools("window.$generativeTraits", &traits),
	)

	if err != nil {
		log.Error("chromedp.Run.err.generativeTraits",err.Error(), err)
		return 
	}

	for key, item := range traits {
		attr := entity.TokenUriAttr{}
		attr.TraitType = key
		attr.Value = item
				
		strAttr := entity.TokenUriAttrStr{}
		strAttr.TraitType = key
		strAttr.Value = fmt.Sprintf("%v", item)

		attrs = append(attrs, attr)
		strAttrs = append(strAttrs, strAttr)
	}

	image := helpers.Base64Encode(buf)
	image = fmt.Sprintf("%s,%s", "data:image/png;base64", image)

	now := time.Now().UTC()
	token.ParsedImage = &image
	token.ParsedAttributes = attrs
	token.ParsedAttributesStr = strAttrs
	token.ThumbnailCapturedAt = &now

	if  token.ParsedImage != nil {
		base64Image := *token.ParsedImage
		i := strings.Index(base64Image, ",")
		if i >= 0 {
			name := fmt.Sprintf("thumb/%s-%s.png", token.ContractAddress, token.TokenID)
			base64Image = base64Image[i+1:]
			uploaded, err := u.GCS.UploadBaseToBucket(base64Image,  name)
			if err != nil {
				log.Error("u.GCS.UploadBaseToBucket", err.Error(), err)
			}else{
				log.SetData("uploaded", uploaded)
				token.Thumbnail = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), name)
			}
		}
		// pass reader to NewDecoder
	}

	updated, err := u.Repo.UpdateOrInsertTokenUri(token.ContractAddress, token.TokenID, token)
	if err != nil {
		log.Error("runAndCap.UpdateOrInsertTokenUri",err.Error(), err)
		return 
	}

	log.SetData("updated", updated)
}

func (u Usecase) CaptureAnimationURI(rootSpan opentracing.Span, token *entity.TokenUri, captureTimeout int) (*structure.TokenAnimationURI,  error) {
	span, log := u.StartSpan("CaptureImageToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	attrs := []entity.TokenUriAttr{}
	strAttrs := []entity.TokenUriAttrStr{}
	resp := &structure.TokenAnimationURI{
		Thumbnail: "",
		Traits:  attrs,
		TraitsStr: strAttrs,
	}

	if token.ParsedImage == nil {
		go u.RunAndCap(span, token, captureTimeout)
		return resp, nil
	}
	
	if token.ThumbnailCapturedAt == nil {
		go u.RunAndCap(span, token, captureTimeout)
		return resp, nil
	}
	
	if token.ThumbnailCapturedAt !=nil {
	 	if token.ThumbnailCapturedAt.Add(time.Hour * 6).After(time.Now()) {
			go u.RunAndCap(span, token, captureTimeout)
			return resp, nil
		}
	}
	
	return resp, nil
}

func (u Usecase) GetToken(rootSpan opentracing.Span, req structure.GetTokenMessageReq, captureTimeout int) (*entity.TokenUri, error) {
	span, log := u.StartSpan("GetToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("req", req)

	defer func() {
		go u.GetLiveToken(span, req, captureTimeout)
	}()

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

	dataObject, err := u.Repo.FindTokenBy(req.ContractAddress, req.TokenID)
	if err != nil {
		log.Error("u.Repo.FindTokenBy", err.Error(), err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			dataObject = &entity.TokenUri{}

		} else {
			log.Error("u.Repo.FindTokenBy", err.Error(), err)
			return nil, err
		}
	} 

	mftMintedTimeChan := make(chan structure.NftMintedTimeChan, 1)
	tokendatachan := make(chan structure.TokenDataChan, 1)

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


	go func(tokenDataChan chan structure.TokenDataChan, parentAddr common.Address, tokenID string) {
		var err error
		tok := &entity.TokenUri{}

		tokenUriData, err := u.getNftProjectTokenUri(client, parentAddr, req.TokenID)
		
		defer func ()  {
			tokenDataChan <- structure.TokenDataChan{
				Data:  tok,
				Err:  err,
			}
		}()

		base64Str := strings.ReplaceAll(*tokenUriData, "data:application/json;base64,", "")
		data, err := helpers.Base64Decode(base64Str)
		if err != nil {
			return 
		}

		stringData := string(data)
		stringData = strings.ReplaceAll(stringData, "\n", "\\n")
		stringData = strings.ReplaceAll(stringData, "\b", "\\b")
		stringData = strings.ReplaceAll(stringData, "\f", "\\f")
		stringData = strings.ReplaceAll(stringData, "\r", "\\r")
		stringData = strings.ReplaceAll(stringData, "\t", "\\t")

		err = json.Unmarshal([]byte(stringData), tok)
		if err != nil {
			return 
		}

	}(tokendatachan, parentAddr, req.TokenID)

	go func(mftMintedTimeChan chan structure.NftMintedTimeChan, genNFTAddr string) {
		nftMintedTime :=  &structure.NftMintedTime{}
		var err error

		defer func ()  {
			mftMintedTimeChan <- structure.NftMintedTimeChan{
				NftMintedTime:  nftMintedTime,
				Err:  err,
			}
		}()

		nftMintedTime, err = u.GetNftMintedTime(span, structure.GetNftMintedTimeReq{
			ContractAddress: genNFTAddr,
			TokenID: req.TokenID,
		})
	}(mftMintedTimeChan, strings.ToLower(parentAddr.String()))
	
	
	log.SetData("nftProject", nftProject)
	log.SetData("parentAddr", parentAddr)
	//log.SetData("tokenUriData", tokenUriData)

	dataObject.ContractAddress = strings.ToLower(req.ContractAddress)
	dataObject.CreatorAddr = strings.ToLower(nftProject.Creator)
	dataObject.OwnerAddr = strings.ToLower(nftProject.Creator)
	dataObject.GenNFTAddr = strings.ToLower(parentAddr.String())

	tokenIDint, _ := strconv.Atoi(req.TokenID)

	dataObject.TokenID = req.TokenID
	dataObject.TokenIDInt = tokenIDint
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

	project, err := u.Repo.FindProjectBy(dataObject.ContractAddress, dataObject.ProjectID)
	if err != nil {
		log.Error("u.GetProjectDetail", err.Error(), err)
		return nil, err
	}

	dataObject.Project = project
	creator, err := u.Repo.FindUserByWalletAddress(dataObject.CreatorAddr)
	if err != nil {
		log.Error("u.Repo.FindUserByWalletAddress.creator", err.Error(), err)
		return nil, err
	}
	dataObject.Creator = creator
	mftMintedTime := <- mftMintedTimeChan
	if mftMintedTime.Err == nil {
		dataObject.BlockNumberMinted = mftMintedTime.NftMintedTime.BlockNumberMinted
		dataObject.MintedTime = mftMintedTime.NftMintedTime.MintedTime
		nft := mftMintedTime.NftMintedTime.Nft
		owner, err := u.Repo.FindUserByWalletAddress(nft.Owner)
		if err != nil {
			log.Error("u.Repo.FindUserByWalletAddress.owner", err.Error(), err)
			return nil, err
		}
		dataObject.Owner = owner

	}else{
		log.Error(" u.GetNftMintedTime",  mftMintedTime.Err.Error(),  mftMintedTime.Err)
	}

	tokenFChan := <- tokendatachan 
	if tokenFChan.Err == nil {
		dataObject.Name = tokenFChan.Data.Name
		dataObject.Description = tokenFChan.Data.Description
		dataObject.Image = tokenFChan.Data.Image
		dataObject.AnimationURL = tokenFChan.Data.AnimationURL
		dataObject.Attributes = tokenFChan.Data.Attributes
		dataObject.Image = tokenFChan.Data.Image
		
	}

	tokIdMini  := dataObject.TokenIDInt %  100000
	dataObject.TokenIDMini = &tokIdMini

	u.CaptureAnimationURI(span, dataObject, 6)
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
