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
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/contracts/generative_nft_contract"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/redis"
)

func (u Usecase) RunAndCap(token *entity.TokenUri) (*structure.TokenAnimationURI, error) {
	captureTimeout := entity.DEFAULT_CAPTURE_TIME
	p, err := u.Repo.FindProjectByTokenID(token.ProjectID)
	if err == nil && p != nil && p.CatureThumbnailDelayTime != nil && *p.CatureThumbnailDelayTime != 0 {
		captureTimeout = *p.CatureThumbnailDelayTime
	}

	var buf []byte
	attrs := []entity.TokenUriAttr{}
	strAttrs := []entity.TokenUriAttrStr{}
	if token == nil {
		return nil, errors.New("Token is empty")
	}
	resp := &structure.TokenAnimationURI{}
	u.Logger.LogAny("RunAndCap", zap.Any("token", token))
	if token.ThumbnailCapturedAt != nil && token.ParsedImage != nil {
		resp = &structure.TokenAnimationURI{
			ParsedImage: *token.ParsedImage,
			Thumbnail:   token.Thumbnail,
			Traits:      token.ParsedAttributes,
			TraitsStr:   token.ParsedAttributesStr,
			CapturedAt:  token.ThumbnailCapturedAt,
			IsUpdated:   false,
		}
		return resp, nil
	}

	eCH, err := strconv.ParseBool(os.Getenv("ENABLED_CHROME_HEADLESS"))
	if err != nil {
		u.Logger.ErrorAny("RunAndCap", zap.Error(err))
		return nil, err
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("google-chrome"),
		chromedp.Flag("headless", eCH),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-first-run", true),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	cctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	imageURL := token.AnimationURL
	if strings.Index(imageURL, "data:text/html;base64,") >= 0 {
		htmlString := strings.ReplaceAll(token.AnimationURL, "data:text/html;base64,", "")
		uploaded, err := u.GCS.UploadBaseToBucket(htmlString, fmt.Sprintf("btc-projects/%s/index.html", token.ProjectID))
		if err == nil {
			fileURI := fmt.Sprintf("%s/%s?seed=%s", os.Getenv("GCS_DOMAIN"), uploaded.Name, token.TokenID)
			imageURL = fileURI
		}
		u.Logger.LogAny("RunAndCap", zap.Any("token", token), zap.Any("fileURI", imageURL), zap.Any("uploaded", uploaded))
	}

	traits := make(map[string]interface{})
	err = chromedp.Run(cctx,
		chromedp.EmulateViewport(960, 960),
		chromedp.Navigate(imageURL),
		chromedp.Sleep(time.Second*time.Duration(captureTimeout)),
		chromedp.CaptureScreenshot(&buf),
		chromedp.EvaluateAsDevTools("window.$generativeTraits", &traits),
	)

	if err != nil {
		u.Logger.Error("RunAndCap", zap.Any("chromedp.Run", zap.Error(err)))
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

	thumbnail := ""
	now := time.Now().UTC()
	if image != "" {
		base64Image := image
		i := strings.Index(base64Image, ",")
		if i >= 0 {
			now := time.Now().UTC().Unix()
			name := fmt.Sprintf("thumb/%s-%d.png", token.TokenID, now)
			base64Image = base64Image[i+1:]
			uploaded, err := u.GCS.UploadBaseToBucket(base64Image, name)
			if err != nil {
				u.Logger.ErrorAny("RunAndCap", zap.Any("UploadBaseToBucket", zap.Error(err)))
			} else {
				u.Logger.LogAny("RunAndCap", zap.Any("uploaded", uploaded))
				thumbnail = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), name)
			}
		}
	}

	resp = &structure.TokenAnimationURI{
		ParsedImage: thumbnail,
		Thumbnail:   thumbnail,
		Traits:      attrs,
		TraitsStr:   strAttrs,
		CapturedAt:  &now,
		IsUpdated:   true,
	}

	u.Logger.LogAny("RunAndCap", zap.Any("token", token), zap.Any("fileURI", imageURL), zap.Any("resp", resp))
	return resp, nil
}

func (u Usecase) GetTokenByTokenID(tokenID string, captureTimeout int) (*entity.TokenUri, error) {

	tokenID = strings.ToLower(tokenID)

	tokenUri, err := u.Repo.FindTokenByTokenID(tokenID)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	///u.Logger.Info("tokenUri", tokenUri)
	u.Logger.Info("tokenID", tokenUri.TokenID)
	return tokenUri, nil
}

func (u Usecase) GetToken(req structure.GetTokenMessageReq, captureTimeout int) (*entity.TokenUri, error) {
	u.Logger.LogAny("GetToken", zap.Any("req", req))
	//contractAddress := strings.ToLower(req.ContractAddress)
	tokenID := strings.ToLower(req.TokenID)

	tokenUri, err := u.Repo.FindTokenByTokenID(tokenID)
	if err != nil {
		u.Logger.ErrorAny("GetToken", zap.Any("req", req), zap.String("action", "FindTokenBy"), zap.Error(err))
		return nil, err
	}
	if tokenUri.Project != nil && tokenUri.InscribedBy != "" {
		tokenUri.Project.InscribedBy = tokenUri.InscribedBy
	}
	client := resty.New()
	resp := &response.SearhcInscription{}
	_, err = client.R().
		EnableTrace().
		SetResult(&resp).
		Get(fmt.Sprintf("%s/inscription/%s", u.Config.GenerativeExplorerApi, tokenUri.TokenID))
	u.Logger.Info("incriptionData", zap.Any("data", resp))
	if err != nil {
		u.Logger.ErrorAny("GetToken.Inscription", zap.Any("req", req), zap.String("action", "Inscription"), zap.Error(err))
		// return nil, err
	}

	tokenUri.Owner = nil
	if resp.Address != "" {
		tokenUri.OwnerAddr = resp.Address
		user, err := u.Repo.FindUserByBtcAddressTaproot(resp.Address)
		if err == nil {
			tokenUri.Owner = user
		}
	}

	//this was used for ETH (old flow)
	// if err != nil {
	// 	u.Logger.ErrorAny("GetToken", zap.Any("req", req), zap.String("action", "FindTokenBy"), zap.Error(err))
	// 	if errors.Is(err, mongo.ErrNoDocuments) {
	// 		token, err := u.getTokenInfo(req)
	// 		if err != nil {
	// 			u.Logger.ErrorAny("GetToken", zap.Any("req", req), zap.String("action", "getProjectDetailFromChain"), zap.Error(err))
	// 			return nil, err
	// 		}
	// 		return token, nil
	// 	} else {
	// 		u.Logger.ErrorAny("GetToken", zap.Any("req", req), zap.String("action", "FindTokenBy"), zap.Error(err))
	// 		return nil, err
	// 	}
	// }

	go func() {
		if tokenUri.Thumbnail == "" {
			payload := redis.PubSubPayload{Data: structure.TokenImagePayload{
				TokenID:         tokenUri.TokenID,
				ContractAddress: tokenUri.ContractAddress,
			}}

			u.Logger.LogAny("GetToken.Thumbnail", zap.Any("payload", payload))
			err = u.PubSub.Producer(utils.PUBSUB_TOKEN_THUMBNAIL, payload)
			if err != nil {
				u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "u.PubSub.Producer"), zap.Error(err))
			}
		}
	}()

	go func() {
		//upload animation URL
		if tokenUri.AnimationHtml == nil {
			p, err := u.Repo.FindProjectByTokenID(tokenUri.ProjectID)
			if err != nil {
				return
			}

			htmlUrl, err := u.parseAnimationURL(*p)
			if err != nil {
				return
			}

			animationHtml := fmt.Sprintf("%s?seed=%s", *htmlUrl, tokenUri.TokenID)
			tokenUri.AnimationHtml = &animationHtml

			_, err = u.Repo.UpdateOrInsertTokenUri(tokenUri.ContractAddress, tokenUri.TokenID, tokenUri)
			if err != nil {
				return
			}
		}

	}()

	///u.Logger.Info("tokenUri", tokenUri)
	u.Logger.LogAny("GetToken", zap.Any("req", req), zap.Any("tokenUri", tokenUri))
	return tokenUri, nil
}

func (u Usecase) getTokenInfo(req structure.GetTokenMessageReq) (*entity.TokenUri, error) {

	u.Logger.Info("req", req)
	addr := common.HexToAddress(req.ContractAddress)
	//fAddr := strings.ToLower(req.ContractAddress)
	isUpdated := false

	dataObject, err := u.Repo.FindTokenByTokenID(req.TokenID)
	if err != nil {
		u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "FindTokenByTokenID"), zap.Error(err))
		if errors.Is(err, mongo.ErrNoDocuments) {
			dataObject = &entity.TokenUri{}
			isUpdated = true
		} else {
			u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "FindTokenByTokenID"), zap.Error(err))
			return nil, err
		}
	}

	mftMintedTimeChan := make(chan structure.NftMintedTimeChan, 1)
	tokendatachan := make(chan structure.TokenDataChan, 1)
	//tokenImageChan := make(chan structure.TokenAnimationURIChan, 1)

	// call to contract to get emotion
	client, err := helpers.EthDialer()
	if err != nil {
		u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "EthDialer"), zap.Error(err))
		return nil, err
	}

	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(req.TokenID, 10)
	if !ok {
		err := errors.New("cannot convert tokenID")
		u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "tokenID.SetString"), zap.Error(err))
		return nil, err
	}
	projectID := new(big.Int).Div(tokenID, big.NewInt(1000000))
	nftProjectDetail, err := u.getProjectDetailFromChain(structure.GetProjectDetailMessageReq{
		ContractAddress: addr.String(),
		ProjectID:       projectID.String(),
	})
	if err != nil {
		u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "getProjectDetailFromChain"), zap.Error(err))
		return nil, err
	}

	nftProject := nftProjectDetail.ProjectDetail
	parentAddr := nftProject.GenNFTAddr

	//get getNftProjectTokenUri
	go func(tokenDataChan chan structure.TokenDataChan, parentAddr common.Address, tokenID string) {
		var err error
		tok := &entity.TokenUri{}

		defer func() {
			tokenDataChan <- structure.TokenDataChan{
				Data: tok,
				Err:  err,
			}
		}()

		tokenUriData, err := u.getNftProjectTokenUri(client, parentAddr, req.TokenID)
		if err != nil {
			u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "getNftProjectTokenUri"), zap.Error(err))
			return
		}

		base64Str := strings.ReplaceAll(*tokenUriData, "data:application/json;base64,", "")
		data, err := helpers.Base64Decode(base64Str)
		if err != nil {
			u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "helpers.Base64Decode"), zap.Error(err))
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
			u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "json.Unmarshal"), zap.Error(err))
			return
		}

	}(tokendatachan, parentAddr, req.TokenID)

	//get minted time
	go func(mftMintedTimeChan chan structure.NftMintedTimeChan, genNFTAddr string) {
		nftMintedTime := &structure.NftMintedTime{}
		var err error

		defer func() {
			mftMintedTimeChan <- structure.NftMintedTimeChan{
				NftMintedTime: nftMintedTime,
				Err:           err,
			}
		}()

		nftMintedTime, err = u.GetNftMintedTime(structure.GetNftMintedTimeReq{
			ContractAddress: genNFTAddr,
			TokenID:         req.TokenID,
		})
	}(mftMintedTimeChan, strings.ToLower(parentAddr.String()))
	dataObject.ContractAddress = strings.ToLower(req.ContractAddress)
	dataObject.CreatorAddr = strings.ToLower(nftProject.Creator)
	dataObject.GenNFTAddr = strings.ToLower(parentAddr.String())

	tokenIDint, _ := strconv.Atoi(req.TokenID)

	dataObject.TokenID = req.TokenID
	dataObject.TokenIDInt = tokenIDint
	dataObject.ProjectID = projectID.String()
	dataObject.ProjectIDInt = projectID.Int64()

	u.Logger.Info("dataObject.ContractAddress", dataObject.ContractAddress)
	u.Logger.Info("dataObject.Creator", dataObject.Creator)
	u.Logger.Info("dataObject.TokenID", dataObject.TokenID)
	u.Logger.Info("dataObject.ProjectID", dataObject.ProjectID)

	project, err := u.Repo.FindProjectBy(dataObject.ContractAddress, dataObject.ProjectID)
	if err != nil {
		u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "findProjectBy"), zap.Error(err))
		return nil, err
	}

	dataObject.Project = project
	creator, err := u.Repo.FindUserByWalletAddress(dataObject.CreatorAddr)
	if err != nil {
		u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "FindUserByWalletAddress"), zap.Error(err))
		creator = &entity.Users{}
	}
	dataObject.Creator = creator
	mftMintedTime := <-mftMintedTimeChan

	if mftMintedTime.Err == nil {
		nft := mftMintedTime.NftMintedTime.Nft
		//onwer
		if nft.Owner != dataObject.OwnerAddr || (dataObject.Owner != nil && nft.Owner != dataObject.Owner.WalletAddress) {

			ownerAddr := strings.ToLower(nft.Owner)

			u.Logger.Info("dataObject.OwnerAddr.old", dataObject.OwnerAddr)
			u.Logger.Info("dataObject.OwnerAddr.new", ownerAddr)
			owner, err := u.Repo.FindUserByWalletAddress(ownerAddr)
			if err != nil {
				u.Logger.Error(err)
				//return nil, err
				owner = &entity.Users{}
			}

			dataObject.Owner = owner
			dataObject.OwnerAddr = ownerAddr
			isUpdated = true
		}

		if mftMintedTime.NftMintedTime.MintedTime != dataObject.MintedTime {
			dataObject.BlockNumberMinted = mftMintedTime.NftMintedTime.BlockNumberMinted
			dataObject.MintedTime = mftMintedTime.NftMintedTime.MintedTime
			isUpdated = true
		}

	} else {
		u.Logger.Error(" u.GetNftMintedTime", mftMintedTime.Err.Error(), mftMintedTime.Err)
	}

	tokenFChan := <-tokendatachan
	if tokenFChan.Err == nil {
		dataObject.Name = tokenFChan.Data.Name
		dataObject.Description = tokenFChan.Data.Description
		dataObject.Image = tokenFChan.Data.Image
		dataObject.AnimationURL = tokenFChan.Data.AnimationURL
		dataObject.Attributes = tokenFChan.Data.Attributes
		dataObject.Image = tokenFChan.Data.Image

	} else {
		u.Logger.Error("tokenFChan.Err", tokenFChan.Err.Error(), tokenFChan.Err)
	}

	tokIdMini := dataObject.TokenIDInt % 100000
	dataObject.TokenIDMini = &tokIdMini

	u.Logger.Info(fmt.Sprintf("Data for minter address %v and OwnerAddr %v", dataObject.MinterAddress, dataObject.OwnerAddr), true)

	isAddress := func(s *string) bool {
		if s == nil {
			return false
		}
		return strings.HasPrefix(*s, "0x")
	}

	if dataObject.MinterAddress != nil {
		u.Logger.Info(fmt.Sprintf("Minter address %s", *dataObject.MinterAddress), true)
	}

	if !isAddress(dataObject.MinterAddress) && dataObject.OwnerAddr != "" {
		u.Logger.Info("Update minter address", true)
		dataObject.MinterAddress = &dataObject.OwnerAddr
		isUpdated = true
	}

	if isUpdated {
		updated, err := u.Repo.UpdateOrInsertTokenUri(dataObject.ContractAddress, dataObject.TokenID, dataObject)
		if err != nil {
			u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "UpdateOrInsertTokenUri"), zap.Error(err))
			return nil, err
		}
		u.Logger.LogAny("getTokenInfo", zap.Any("req", req), zap.Any("updated", updated), zap.String("action", "UpdateOrInsertTokenUri"))
	}

	//capture image
	payload := redis.PubSubPayload{Data: structure.TokenImagePayload{
		TokenID:         dataObject.TokenID,
		ContractAddress: dataObject.ContractAddress,
	}}

	err = u.PubSub.Producer(utils.PUBSUB_TOKEN_THUMBNAIL, payload)
	if err != nil {
		u.Logger.ErrorAny("getTokenInfo", zap.Any("req", req), zap.String("action", "u.PubSub.Producer"), zap.Error(err))
	}

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

func (u Usecase) UpdateTokensFromChain() error {

	//TODO - we will use pagination instead of all
	tokens, err := u.Repo.GetAllTokens()
	if err != nil {
		u.Logger.Error(err)
		return err
	}

	u.Logger.Info("tokens.Count", len(tokens))
	for _, token := range tokens {

		_, err := u.GetToken(structure.GetTokenMessageReq{ContractAddress: token.ContractAddress, TokenID: token.TokenID}, 5)
		if err != nil {
			u.Logger.Error(err)
			return err
		}
	}

	return nil
}

func (u Usecase) GetTokensByContract(contractAddress string, filter nfts.MoralisFilter) (*entity.Pagination, error) {

	client, err := helpers.EthDialer()
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	contractAddr := common.HexToAddress(contractAddress)
	gNft, err := generative_nft_contract.NewGenerativeNftContract(contractAddr, client)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	project, err := gNft.Project(nil)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}
	parentAddr := project.ProjectAddr

	resp, err := u.MoralisNft.GetNftByContract(contractAddress, filter)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}
	parentAddrStr := parentAddr.String()
	result := []entity.TokenUri{}
	for _, item := range resp.Result {
		tokenID := item.TokenID
		token, err := u.GetToken(structure.GetTokenMessageReq{ContractAddress: parentAddrStr, TokenID: tokenID}, 5)
		if err != nil {
			u.Logger.Error(err)
			return nil, err
		}
		result = append(result, *token)
	}

	p := &entity.Pagination{}
	p.Result = result
	p.Cursor = resp.Cursor
	p.Total = int64(resp.Total)
	p.Page = int64(resp.Page)
	p.PageSize = int64(resp.PageSize)
	return p, nil
}

func (u Usecase) FilterTokens(filter structure.FilterTokens) (*entity.Pagination, error) {
	pe := &entity.FilterTokenUris{}
	err := copier.Copy(pe, filter)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	tokens, err := u.Repo.FilterTokenUri(*pe)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("tokens", tokens.Total)
	return tokens, nil
}

func (u Usecase) FilterTokensNew(filter structure.FilterTokens) (*entity.Pagination, error) {
	pe := &entity.FilterTokenUris{}
	

	//filerAttrs := []structure.TokenUriAttrReq{}
	if filter.Rarity != nil && *filter.Rarity != "" {
		r := strings.Split(*filter.Rarity, ",")
		min := "0"
		max := "100"
		if len(r) == 2 {
			min = r[0]
			max = r[1]
		}

		minInt, _ := strconv.Atoi(min)
		maxInt, _ := strconv.Atoi(max)

		groupTraits := make(map [string][]string)
		p, err := u.Repo.FindProjectByTokenID(*filter.GenNFTAddr)
		if err == nil {
			traits := p.TraitsStat
			for _, trait := range traits {
				values := trait.TraitValuesStat
				
				for _, value := range values {
					if  value.Rarity >= int32(minInt) && value.Rarity <= int32(maxInt) {
						groupTraits[trait.TraitName] =   append(groupTraits[trait.TraitName], value.Value)
						
					}
				}
			}
		}

		for key, groupTrait := range  groupTraits {
			r := structure.TokenUriAttrReq{}
			r.TraitType = key
			r.Values = groupTrait
			filter.RarityAttributes = append(filter.Attributes, r)
		}
	}

	err := copier.Copy(pe, filter)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}
	
	tokens, err := u.Repo.FilterTokenUriNew(*pe)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("tokens", tokens.Total)
	return tokens, nil
}

func (u Usecase) UpdateToken(req structure.UpdateTokenReq) (*entity.TokenUri, error) {

	p, err := u.Repo.FindTokenBy(req.ContracAddress, req.TokenID)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	if req.Priority != nil {
		p.Priority = req.Priority
	}

	updated, err := u.Repo.UpdateOrInsertTokenUri(req.ContracAddress, req.TokenID, p)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("updated", updated)
	return p, nil
}

func (u Usecase) GetTokensOfAProjectFromChain(project entity.Projects) error {

	contractAddres := project.ContractAddress
	genAddress := project.GenNFTAddr
	// projectID := project.TokenID
	// ProjectIDInt := project.TokenIDInt

	chain := os.Getenv("MORALIS_CHAIN")
	nfts, err := u.MoralisNft.GetNftByContract(genAddress, nfts.MoralisFilter{Chain: &chain})
	if err != nil {
		u.Logger.Error(err)
		return err
	}

	processed := 0
	tokens := nfts.Result
	for _, token := range tokens {
		if processed%5 == 0 {
			time.Sleep(10 * time.Second)
		}

		go func(contractAddres string, tokenID string) {
			u.GetToken(structure.GetTokenMessageReq{
				ContractAddress: contractAddres,
				TokenID:         tokenID,
			}, 20)
		}(contractAddres, token.TokenID)

		processed++
	}

	return nil
}

func (u Usecase) CreateBTCTokenURI(projectID string, tokenID string, mintedURL string, paidType entity.TokenPaidType, nftTokenIds ...string) (*entity.TokenUri, error) {

	// find project by projectID
	u.Logger.Info(utils.TOKEN_ID_TAG, tokenID)
	u.Logger.Info(utils.PROJECT_ID_TAG, projectID)
	u.Logger.Info("mintedURL", mintedURL)

	project, err := u.Repo.FindProjectByTokenID(projectID)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	tokenUri := entity.TokenUri{}
	tokenUri.ContractAddress = project.ContractAddress
	tokenUri.TokenID = tokenID
	blockNumberMinted := "31012412"
	tokenUri.BlockNumberMinted = &blockNumberMinted
	tokenUri.Creator = &project.CreatorProfile
	tokenUri.CreatorAddr = project.CreatorAddrr
	tokenUri.Description = project.Description
	tokenUri.GenNFTAddr = project.GenNFTAddr

	mintedTime := time.Now()
	tokenUri.MintedTime = &mintedTime
	tokenUri.Name = tokenID
	tokenUri.Project = project
	tokenUri.ProjectID = project.TokenID
	tokenUri.ProjectIDInt = project.TokenIDInt
	tokenUri.PaidType = paidType
	tokenUri.IsOnchain = false
	if len(nftTokenIds) > 0 {
		tokenUri.NftTokenId = nftTokenIds[0]
	}
	if len(nftTokenIds) > 1 {
		tokenUri.InscribedBy = nftTokenIds[1]
	}

	nftTokenUri := project.NftTokenUri
	u.Logger.Info("nftTokenUri", nftTokenUri)

	projectNftTokenUri := &structure.ProjectAnimationUrl{}
	err = helpers.Base64DecodeRaw(project.NftTokenUri, projectNftTokenUri)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	imageURI := ""
	if projectNftTokenUri.AnimationUrl != "" {
		u.Logger.Info("nftTokenUri", len(nftTokenUri))
		base64Data := strings.Replace(nftTokenUri, "data:application/json;base64,", "", 1)

		type Data struct {
			AnimationUrl string `bson:"animation_url" json:"animation_url"`
		}

		var data Data

		err = helpers.Base64DecodeRaw(base64Data, &data)

		if err != nil {
			return nil, err
		}
		imageURI = data.AnimationUrl
		tokenUri.AnimationURL = imageURI
	} else if strings.Index(mintedURL, ".html") != -1 {
		imageURI = mintedURL
		tokenUri.AnimationURL = mintedURL
	} else {
		now := time.Now().UTC()
		imageURI = mintedURL
		tokenUri.AnimationURL = ""
		tokenUri.Thumbnail = mintedURL
		tokenUri.Image = mintedURL
		tokenUri.ParsedImage = &mintedURL
		tokenUri.ThumbnailCapturedAt = &now
		u.Logger.Info("mintedURL", mintedURL)
	}

	tokenUri.OrderInscriptionIndex = int(project.MintingInfo.Index + 1)
	_, err = u.Repo.UpdateOrInsertTokenUri(tokenUri.ContractAddress, tokenUri.TokenID, &tokenUri)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}
	// after update, increase index field in project to 1
	err = u.Repo.IncreaseProjectIndex(projectID)
	if err != nil {
		return nil, err
	}
	pTokenUri, err := u.Repo.FindTokenBy(tokenUri.ContractAddress, tokenUri.TokenID)
	if err != nil {
		return nil, err
	}

	//capture image
	payload := redis.PubSubPayload{Data: structure.TokenImagePayload{
		TokenID:         pTokenUri.TokenID,
		ContractAddress: pTokenUri.ContractAddress,
	}}

	err = u.PubSub.Producer(utils.PUBSUB_TOKEN_THUMBNAIL, payload)
	if err != nil {
		u.Logger.Error(err)
	}

	return pTokenUri, nil
}

func (u Usecase) GetAllListListingWithRule() ([]structure.MarketplaceNFTDetail, error) {

	result := []structure.MarketplaceNFTDetail{}
	var nftList []entity.MarketplaceBTCListingFilterPipeline
	var err error

	nftList, err = u.Repo.RetrieveBTCNFTListingsUnsold(9999999, 0)
	if err != nil {
		return nil, err
	}
	for _, listing := range nftList {

		buyOrders, err := u.Repo.GetBTCListingHaveOngoingOrder(listing.UUID)
		if err != nil {
			continue

		}
		currentTime := time.Now()
		isAvailable := true
		for _, order := range buyOrders {
			expireTime := order.ExpiredAt
			// not expired yet still waiting for btc
			if currentTime.Before(expireTime) && (order.Status == entity.StatusBuy_Pending || order.Status == entity.StatusBuy_NotEnoughBalance) {
				isAvailable = false
				break
			}
			// could be expired but received btc
			if order.Status != entity.StatusBuy_Pending && order.Status != entity.StatusBuy_NotEnoughBalance {
				isAvailable = false
				break
			}
		}

		nftInfo := structure.MarketplaceNFTDetail{
			InscriptionID: listing.InscriptionID,
			Name:          listing.Name,
			Description:   listing.Description,
			Price:         listing.Price,
			OrderID:       listing.UUID,
			IsConfirmed:   listing.IsConfirm,
			Buyable:       isAvailable,
			IsCompleted:   listing.IsSold,
			CreatedAt:     listing.CreatedAt,
			PayType:       listing.PayType,
		}
		result = append(result, nftInfo)
	}
	return result, nil
}

func (u Usecase) GetListingDetail(inscriptionID string) (*structure.MarketplaceNFTDetail, error) {
	// addon for check isBuyable (contact Phuong)

	isBuyable := true
	nft, err := u.Repo.FindBtcNFTListingUnsoldByNFTID(inscriptionID)
	if err != nil {
		return nil, err
	}

	if !nft.IsSold {
		buyOrders, err := u.Repo.GetBTCListingHaveOngoingOrder(nft.UUID)
		if err != nil {
			return nil, err
		}
		currentTime := time.Now()
		for _, order := range buyOrders {
			expireTime := order.ExpiredAt
			// not expired yet still waiting for btc
			if currentTime.Before(expireTime) && (order.Status == entity.StatusBuy_Pending || order.Status == entity.StatusBuy_NotEnoughBalance) {
				isBuyable = false
				break
			}
			// could be expired but received btc
			if order.Status != entity.StatusBuy_Pending && order.Status != entity.StatusBuy_NotEnoughBalance {
				isBuyable = false
				break
			}
		}
	}

	listingPaymentInfo, err := u.getListingPaymentInfo(nft.PayType, nft.Price)
	if err != nil {
		return nil, err
	}

	nftInfo := structure.MarketplaceNFTDetail{
		InscriptionID:      nft.InscriptionID,
		Name:               nft.Name,
		Description:        nft.Description,
		Price:              nft.Price,
		OrderID:            nft.UUID,
		IsConfirmed:        nft.IsConfirm,
		Buyable:            isBuyable,
		IsCompleted:        nft.IsSold,
		PaymentListingInfo: listingPaymentInfo,
	}
	return &nftInfo, nil

}

func (u Usecase) UpdateTokenThumbnail(req structure.UpdateTokenThumbnailReq) (*entity.TokenUri, error) {

	token, err := u.Repo.FindTokenByTokenID(req.TokenID)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	if strings.Index(token.Image, ".glb") == -1 {
		err = errors.New("Token's image is not a glb file")
		u.Logger.Error(err)
		return nil, err
	}
	now := time.Now().Unix()

	base64Data := strings.ReplaceAll(req.Thumbnail, "data:image/png;base64,", "")
	uploaded, err := u.GCS.UploadBaseToBucket(base64Data, fmt.Sprintf("btc-projects/%s/thumb/token-%s-%d.png", token.ProjectID, token.TokenID, now))
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}
	u.Logger.Info("uploaded", uploaded)
	thumb := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	spew.Dump(thumb)
	token.Thumbnail = thumb

	updated, err := u.Repo.UpdateOrInsertTokenUri(token.ContractAddress, token.TokenID, token)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("updated", updated)
	return token, nil
}

// When go to this, you need to make sure that meta's project is created
func (u Usecase) CreateBTCTokenURIFromCollectionInscription(meta entity.CollectionMeta, inscription entity.CollectionInscription) (*entity.TokenUri, error) {
	// find project by projectID
	project, err := u.Repo.FindProjectByInscriptionIcon(meta.InscriptionIcon)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			u.Logger.ErrorAny("CanNotFindProjectByInscriptionIcon", zap.Any("inscriptionIcon", meta.InscriptionIcon))
			return nil, repository.ErrNoProjectsFound
		} else {
			return nil, err
		}
	}

	tokenUri := entity.TokenUri{}
	tokenUri.ContractAddress = project.ContractAddress
	tokenUri.TokenID = inscription.ID
	blockNumberMinted := "31012412"
	tokenUri.BlockNumberMinted = &blockNumberMinted
	tokenUri.Creator = &project.CreatorProfile
	tokenUri.CreatorAddr = project.CreatorAddrr
	tokenUri.Description = project.Description
	tokenUri.GenNFTAddr = project.GenNFTAddr

	mintedTime := time.Now()
	tokenUri.MintedTime = &mintedTime
	tokenUri.Name = inscription.Meta.Name
	tokenUri.Project = project
	tokenUri.ProjectID = project.TokenID
	tokenUri.ProjectIDInt = project.TokenIDInt
	tokenUri.IsOnchain = false
	tokenUri.CreatedByCollectionInscription = true
	count, err := u.Repo.CountTokenUriByProjectId(tokenUri.ProjectID)
	if err == nil && count != nil {
		tokenUri.OrderInscriptionIndex = int(*count + 1)
	}

	nftTokenUri := project.NftTokenUri
	u.Logger.Info("nftTokenUri", nftTokenUri)

	projectNftTokenUri := &structure.ProjectAnimationUrl{}
	err = helpers.Base64DecodeRaw(project.NftTokenUri, projectNftTokenUri)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	now := time.Now().UTC()
	imageURI := fmt.Sprintf("https://ordinals-explorer.generative.xyz/content/%s", inscription.ID)
	tokenUri.AnimationURL = ""
	tokenUri.Thumbnail = imageURI
	tokenUri.Image = imageURI
	tokenUri.ParsedImage = &imageURI
	tokenUri.ThumbnailCapturedAt = &now
	tokenUri.Source = inscription.Source
	u.Logger.Info("mintedURL", imageURI)

	_, err = u.Repo.UpdateOrInsertTokenUri(tokenUri.ContractAddress, tokenUri.TokenID, &tokenUri)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}
	pTokenUri, err := u.Repo.FindTokenBy(tokenUri.ContractAddress, tokenUri.TokenID)
	if err != nil {
		return nil, err
	}

	return pTokenUri, nil
}

func (u Usecase) parseAnimationURL(project entity.Projects) (*string, error) {
	base64 := strings.ReplaceAll(project.NftTokenUri, "data:application/json;base64,", "")
	jsonData, err := helpers.Base64Decode(base64)
	if err != nil {
		return nil, err
	}
	resp := &structure.ProjectAnimationUrl{}
	err = json.Unmarshal(jsonData, resp)
	if err != nil {
		return nil, err
	}

	if resp.AnimationUrl == "" {
		return nil, errors.New("This project doesn't contain html")
	}

	fName := fmt.Sprintf("btc-projects/%s/index.html", project.TokenID)
	htmlString := strings.ReplaceAll(resp.AnimationUrl, "data:text/html;base64,", "")
	uploaded, err := u.GCS.UploadBaseToBucket(htmlString, fName)
	if err != nil {
		return nil, err
	}

	link := fmt.Sprintf("%s/%s/%s", "https://storage.googleapis.com", os.Getenv("GCS_BUCKET"), uploaded.Name)
	spew.Dump(link)
	return &link, nil

}
