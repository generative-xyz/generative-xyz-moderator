package usecase

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/html"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"rederinghub.io/utils/request"
	"sort"
	"strconv"
	"strings"
	"time"

	"rederinghub.io/utils/contracts/bfs"

	"github.com/chromedp/chromedp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"rederinghub.io/external/generativeexplorer"
	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/contracts/generative_nft_contract"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
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
	if token.ThumbnailCapturedAt != nil && token.ParsedImage != nil && !strings.HasSuffix(*token.ParsedImage, "i0") {
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
		logger.AtLog.Logger.Error("RunAndCap", zap.Any("tokenID", token.TokenID), zap.Error(err))
		return nil, err
	}

	var contextOpts = []chromedp.ContextOption{}
	//contextOpts = []chromedp.ContextOption{
	//	chromedp.WithErrorf(log.Printf),
	//	chromedp.WithLogf(log.Printf),
	//	chromedp.WithBrowserOption(),
	//}

	opts := []chromedp.ExecAllocatorOption{}
	if os.Getenv("ENV") == "mainnet" {
		opts = append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath("google-chrome"),
			chromedp.Flag("headless", eCH),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("no-first-run", true),
		)
	} else {
		opts = append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("no-first-run", true),
		)
	}

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	cctx, cancel := chromedp.NewContext(allocCtx, contextOpts...)

	//avoid overlap html
	ackCtx, cancel := context.WithTimeout(cctx, time.Duration(captureTimeout)*5*time.Second)
	defer cancel()

	imageURL := token.AnimationURL
	if len(imageURL) == 0 {
		parsedImage := ""
		if token.ParsedImage != nil {
			parsedImage = *token.ParsedImage
		}

		resp = &structure.TokenAnimationURI{
			ParsedImage: parsedImage,
			Thumbnail:   token.Thumbnail,
			Traits:      token.ParsedAttributes,
			TraitsStr:   token.ParsedAttributesStr,
			CapturedAt:  token.ThumbnailCapturedAt,
			IsUpdated:   false,
		}
		return resp, nil
	}
	if strings.Index(imageURL, "data:text/html;base64,") >= 0 {
		if token.AnimationHtml != nil && *token.AnimationHtml != "" {
			imageURL = *token.AnimationHtml
		} else {
			htmlString := strings.ReplaceAll(token.AnimationURL, "data:text/html;base64,", "")
			uploaded, err := u.GCS.UploadBaseToBucket(htmlString, fmt.Sprintf("btc-projects/%s/index.html", token.ProjectID))
			if err == nil {
				fileURI := fmt.Sprintf("%s/%s?seed=%s", os.Getenv("GCS_DOMAIN"), uploaded.Name, token.TokenID)
				imageURL = fileURI
			}
		}
	}

	traits := make(map[string]interface{})

	actions := []chromedp.Action{
		chromedp.EmulateViewport(960, 960),
		chromedp.Navigate(imageURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(time.Second * time.Duration(captureTimeout)),
		chromedp.CaptureScreenshot(&buf),
		chromedp.EvaluateAsDevTools("window.$generativeTraits", &traits),
	}

	if token.ProjectID == os.Getenv("MODULAR_PROJECT_ID") {
		actions = append(actions, u.elementScreenshot("#defaultCanvas0", &buf))
	}

	err = chromedp.Run(ackCtx,
		actions...,
	)

	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("RunAndCap - %s - %s", token.ProjectID, token.TokenID),
			zap.String("tokenID", token.TokenID),
			zap.String("contractAddress", token.ContractAddress),
			zap.String("cenNFTAddr", token.GenNFTAddr),
			zap.String("projectID", token.ProjectID),
			zap.String("tokenID", token.TokenID),
			zap.String("fileURI", imageURL),
			zap.Error(err))
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
				logger.AtLog.Logger.Error("RunAndCap", zap.Any("tokenID", token.TokenID), zap.Error(err))
			} else {
				thumbnail = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), name)
			}

			_ = uploaded
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

	if token.ProjectID == os.Getenv("MODULAR_PROJECT_ID") {
		attrs1 := []entity.TokenUriAttr{}
		attrs2 := []entity.TokenUriAttrStr{}

		for _, i := range resp.Traits {
			if strings.EqualFold(i.TraitType, "hash") {
				continue
			}
			attr := entity.TokenUriAttr{}
			attr.Value = i.Value
			attr.TraitType = strings.ReplaceAll(i.TraitType, "Modular's ", "")
			attrs1 = append(attrs1, attr)
		}

		for _, i := range resp.TraitsStr {
			if strings.EqualFold(i.TraitType, "hash") {
				continue
			}
			attr := entity.TokenUriAttrStr{}
			attr.Value = i.Value
			attr.TraitType = strings.ReplaceAll(i.TraitType, "Modular's ", "")
			attrs2 = append(attrs2, attr)
		}

		resp.Traits = attrs1
		resp.TraitsStr = attrs2
	}

	logger.AtLog.Logger.Info(fmt.Sprintf("RunAndCap - %s - %s", token.ProjectID, token.TokenID),
		zap.String("contractAddress", token.ContractAddress),
		zap.String("cenNFTAddr", token.GenNFTAddr),
		zap.String("projectID", token.ProjectID),
		zap.String("tokenID", token.TokenID),
		zap.String("fileURI", imageURL),
		zap.String("uploaded", resp.Thumbnail),
	)

	return resp, nil
}

func (u Usecase) GetTokenByTokenID(tokenID string, captureTimeout int) (*entity.TokenUri, error) {

	tokenID = strings.ToLower(tokenID)

	tokenUri, err := u.Repo.FindTokenByTokenID(tokenID)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	///logger.AtLog.Logger.Info("tokenUri", zap.Any("tokenUri", tokenUri))
	logger.AtLog.Logger.Info("tokenID", zap.Any("tokenUri.TokenID", tokenUri.TokenID))
	return tokenUri, nil
}

func (u Usecase) GetToken(req structure.GetTokenMessageReq, captureTimeout int) (*entity.TokenUri, error) {

	logger.AtLog.Logger.Info("GetToken", zap.Any("req", zap.Any("req)", req)))
	tokenID := strings.ToLower(req.TokenID)
	tokenUri, err := u.Repo.FindTokenByTokenID(tokenID)
	if err != nil {
		if !helpers.IsOrdinalProjectToken(req.TokenID) {
			//this was used for ETH (old flow), try to get DB
			if errors.Is(err, mongo.ErrNoDocuments) {
				token, err2 := u.getTokenInfo(req)
				if err2 != nil {
					logger.AtLog.Logger.Error("GetToken", zap.Any("req", req), zap.String("tokenID", tokenID), zap.String("action", "getProjectDetailFromChain"), zap.Error(err2))
					return nil, err2
				}
				return token, nil
			}
		}
		logger.AtLog.Logger.Error("GetToken", zap.Any("req", req), zap.String("tokenID", tokenID), zap.String("action", "FindTokenBy"), zap.Error(err))
		return nil, err
	}
	if tokenUri.Project != nil && tokenUri.InscribedBy != "" {
		tokenUri.Project.InscribedBy = tokenUri.InscribedBy
	}
	if tokenUri.NftTokenId != "" {
		inscribeBtc := &entity.InscribeBTC{}
		if err = u.Repo.FindOneBy(context.Background(), inscribeBtc.TableName(), bson.M{"inscriptionID": tokenUri.TokenID}, inscribeBtc); err == nil {
			tokenUri.Project.OrdinalsTx = inscribeBtc.OrdinalsTx
			tokenUri.Project.OwnerOf = inscribeBtc.OwnerOf
			tokenUri.Project.TokenAddress = inscribeBtc.TokenAddress
			tokenUri.Project.TokenId = inscribeBtc.TokenId
		}
	}

	if helpers.IsOrdinalProjectToken(req.TokenID) {
		client := resty.New()
		resp := &response.SearhcInscription{}
		_, err = client.R().
			EnableTrace().
			SetResult(&resp).
			Get(fmt.Sprintf("%s/inscription/%s", u.Config.GenerativeExplorerApi, tokenUri.TokenID))
		logger.AtLog.Logger.Info("incriptionData", zap.Any("data", zap.Any("resp)", resp)))
		if err != nil {
			logger.AtLog.Logger.Error("GetToken.Inscription", zap.Any("req", req), zap.String("action", "Inscription"), zap.Error(err))
		} else {
			tokenUri.Owner = nil
			if resp.Address != "" {
				tokenUri.OwnerAddr = resp.Address
				user, err := u.Repo.FindUserByBtcAddressTaproot(resp.Address)
				if err == nil {
					tokenUri.Owner = user
				}
			}
		}
	} else {
		//client, err1 := helpers.ChainDialer(os.Getenv("TC_ENDPOINT_PUBLIC"))
		//if err1 != nil {
		//	logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", tokenID), zap.Any("req", req), zap.String("action", "EthDialer"), zap.String("dial", os.Getenv("TC_ENDPOINT_PUBLIC")), zap.Error(err1))
		//} else {
		//	addr, err2 := u.ownerOf(client, common.HexToAddress(tokenUri.GenNFTAddr), tokenID)
		//	if err2 != nil {
		//		logger.AtLog.Logger.Error("getTokenInfo get ownerOf", zap.String("tokenID", tokenID), zap.Any("req", req), zap.String("action", "ownerOf"), zap.Error(err2))
		//	} else {
		//		if addr != nil {
		//			if tokenUri.OwnerAddr != addr.String() {
		//				tokenUri.Owner = nil
		//				tokenUri.OwnerAddr = addr.String()
		//				user, err := u.Repo.FindUserByWalletAddress(addr.String())
		//				if err == nil && user != nil {
		//					tokenUri.Owner = user
		//				}
		//			}
		//		}
		//	}
		//}

		user, err := u.Repo.FindUserByWalletAddress(tokenUri.OwnerAddr)
		if err == nil && user != nil {
			tokenUri.Owner = user
		}

		mkl, err := u.Repo.FindActivateListingByTokenID(tokenID)
		if err == nil && mkl != nil {
			tokenUri.Buyable = true
			tokenUri.PriceBrc20 = entity.PriceBRC20Obj{
				Value:      mkl.Price,
				Address:    mkl.Erc20Token,
				OfferingID: mkl.OfferingId,
			}
		} else {
			tokenUri.Buyable = false
			tokenUri.PriceBrc20 = entity.PriceBRC20Obj{
				Value:      "",
				Address:    "",
				OfferingID: "",
			}
		}

	}

	go func() {
		if tokenUri.ThumbnailCapturedAt == nil {
			payload := redis.PubSubPayload{Data: structure.TokenImagePayload{
				TokenID:         tokenUri.TokenID,
				ContractAddress: tokenUri.ContractAddress,
			}}

			err = u.PubSub.Producer(utils.PUBSUB_TOKEN_THUMBNAIL, payload)
			if err != nil {
				logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", tokenID), zap.Any("req", req), zap.String("action", "u.PubSub.Producer"), zap.Error(err))
			}
		}
	}()

	go func() {
		//upload animation URL
		if tokenUri.AnimationHtml == nil {
			p, err := u.Repo.FindProjectByTokenID(tokenUri.ProjectID)
			if err != nil {
				logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", tokenID), zap.Any("req", req), zap.Error(err))
				return
			}

			htmlUrl, err := u.parseAnimationURL(*p)
			if err != nil {
				logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", tokenID), zap.Any("req", req), zap.Error(err))
				return
			}

			animationHtml := fmt.Sprintf("%s?seed=%s", *htmlUrl, tokenUri.TokenID)
			tokenUri.AnimationHtml = &animationHtml

			_, err = u.Repo.UpdateOrInsertTokenUri(tokenUri.ContractAddress, tokenUri.TokenID, tokenUri)
			if err != nil {
				logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", tokenID), zap.Any("req", req), zap.Error(err))
				return
			}
		}

	}()

	///logger.AtLog.Logger.Info("tokenUri", zap.Any("tokenUri", tokenUri))
	return tokenUri, nil
}

func (u Usecase) getTokenInfo(req structure.GetTokenMessageReq) (*entity.TokenUri, error) {

	logger.AtLog.Logger.Info("req", zap.Any("req", req))
	addr := common.HexToAddress(req.ContractAddress)
	isUpdated := false

	dataObject, err := u.Repo.FindTokenByTokenID(req.TokenID)
	if err != nil {
		logger.AtLog.Logger.Error("getTokenInfo", zap.Any("req", req), zap.String("action", "FindTokenByTokenID"), zap.Error(err))
		if errors.Is(err, mongo.ErrNoDocuments) {
			dataObject = &entity.TokenUri{}
			isUpdated = true
		} else {
			logger.AtLog.Logger.Error("getTokenInfo", zap.Any("req", req), zap.String("action", "FindTokenByTokenID"), zap.Error(err))
			return nil, err
		}
	}

	mftMintedTimeChan := make(chan structure.NftMintedTimeChan, 1)
	tokendatachan := make(chan structure.TokenDataChan, 1)
	//tokenImageChan := make(chan structure.TokenAnimationURIChan, 1)

	// call to contract to get emotion
	client := u.TcClientPublicNode.GetClient()
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(req.TokenID, 10)
	if !ok {
		err := errors.New("cannot convert tokenID")
		logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", req.TokenID), zap.Error(err))
		return nil, err
	}
	projectID := new(big.Int).Div(tokenID, big.NewInt(1000000))
	nftProjectDetail, err := u.getProjectDetailFromChain(structure.GetProjectDetailMessageReq{
		ContractAddress: addr.String(),
		ProjectID:       projectID.String(),
	})
	if err != nil {
		logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", req.TokenID), zap.String("action", "getProjectDetailFromChain"), zap.Error(err))
		return nil, err
	}

	nftProject := nftProjectDetail.ProjectDetail
	parentAddr := nftProject.GenNFTAddr

	//get getNftProjectTokenUri
	go func(tokenDataChan chan structure.TokenDataChan, parentAddr common.Address, tokenID string) {
		var err error
		tok := &entity.TokenUri{}
		fromBFS := false

		tokeBFS := entity.TokenFromBase64{}

		defer func() {
			tokenDataChan <- structure.TokenDataChan{
				Data: tok,
				Err:  err,
			}
		}()

		tokenUriData, err := u.getNftProjectTokenUri(client, parentAddr, req.TokenID)
		if err != nil {
			logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", req.TokenID), zap.Error(err))
			return
		}
		seed, err := u.getSeedFromTokenId(client, parentAddr, tokenID)
		if err != nil {
			u.Logger.ErrorAny("getTokenInfo not valid seed", zap.String("tokenID", req.TokenID), zap.Any("error", err))
			return
		}
		tok.Seed = *seed

		if strings.Index(*tokenUriData, "data:application/json;base64,") == -1 {
			if strings.Index(*tokenUriData, "bfs://") > -1 {
				bfsContract := common.HexToAddress(os.Getenv("BFS_CONTRACT"))
				tokenUriData, err = u.getBFSData(client, bfsContract, parentAddr, tok.Seed)
				if err != nil {
					u.Logger.ErrorAny("getTokenInfo not valid seed", zap.Any("BFS_CONTRACT", os.Getenv("BFS_CONTRACT")), zap.String("tokenID", req.TokenID), zap.Any("error", err))
					return
				}

				fromBFS = true

			} else {
				u.Logger.ErrorAny("getTokenInfo not valid", zap.String("tokenID", req.TokenID))
				return
			}
		}

		base64Str := strings.ReplaceAll(*tokenUriData, "data:application/json;base64,", "")
		data, err := helpers.Base64Decode(base64Str)
		if err != nil {
			logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", req.TokenID), zap.Error(err))
			return
		}

		stringData := string(data)
		//stringData = strings.ReplaceAll(stringData, "\n", "")
		//stringData = strings.ReplaceAll(stringData, "\b", "")
		//stringData = strings.ReplaceAll(stringData, "\f", "")
		//stringData = strings.ReplaceAll(stringData, "\r", "")
		//stringData = strings.ReplaceAll(stringData, "\t", "")
		//stringData = strings.ReplaceAll(stringData, "\\'", "'")

		if fromBFS {
			err = json.Unmarshal([]byte(stringData), &tokeBFS)
			if err != nil {
				logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", req.TokenID), zap.String("action", "json.Unmarshal"), zap.Error(err))
				return
			}

			imageURL := ""
			imageArr := strings.Split(tokeBFS.Image, ",")
			if len(imageArr) == 2 {
				ext := helpers.FileType(imageArr[0])
				now := time.Now().UTC().UnixMicro()

				//fName := fmt.Sprintf("btc-projects/%d/%d-%s%s", projectID, now, tokenID, ext)
				fName := fmt.Sprintf("thumb/%s-%d%s", tokenID, now, ext)
				uploaded, err := u.GCS.UploadBaseToBucket(imageArr[1], fName)
				if err == nil {
					imageURL = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
				}
			}

			tok.Name = tokeBFS.Name
			tok.Description = tokeBFS.Description
			tok.Image = imageURL
			tok.Thumbnail = imageURL
			tok.ParsedImage = &imageURL
			tok.AnimationURL = tokeBFS.AnimationURL
			now := time.Now().UTC()
			tok.ThumbnailCapturedAt = &now

			attrs := []entity.TokenUriAttr{}
			for _, attr := range tokeBFS.Attributes {
				tmp := entity.TokenUriAttr{
					TraitType: attr.TraitType,
					Value:     attr.Value,
				}

				attrs = append(attrs, tmp)
			}

			tok.ParsedAttributes = attrs
			tok.ParsedAttributesStr = tokeBFS.Attributes
			tok.Attributes = ""

			return
		}
		err = json.Unmarshal([]byte(stringData), tok)
		tok.Name = ""
		// try base64 decode description
		if len(tok.Description) > 0 {
			desc, errDecode := helpers.Base64Decode(tok.Description)
			if errDecode == nil {
				tok.Description = string(desc)
			}
		}
		if err != nil {
			logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", req.TokenID), zap.String("action", "json.Unmarshal"), zap.Error(err))
			return
		}

		//TODO - upload the base64 image into GCS

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

		nftMintedTime, err = u.GetNftMintedTime(client, structure.GetNftMintedTimeReq{
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

	logger.AtLog.Logger.Info("dataObject.ContractAddress", zap.Any("dataObject.ContractAddress", dataObject.ContractAddress), zap.Any("dataObject.Creator", dataObject.Creator), zap.String("tokenID", req.TokenID), zap.Any("dataObject.ProjectID", dataObject.ProjectID))

	project, err := u.Repo.FindProjectBy(dataObject.ContractAddress, dataObject.ProjectID)
	if err != nil {
		logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", req.TokenID), zap.Error(err))
		return nil, err
	}

	dataObject.Project = project
	/*if nftProjectDetail.NftProjectDetail.IndexReserve != nil && nftProjectDetail.NftProjectDetail.Index != nil {
		if project.MintingInfo.Index != nftProjectDetail.NftProjectDetail.Index.Int64() || project.MintingInfo.IndexReverse != nftProjectDetail.NftProjectDetail.IndexReserve.Int64() {
			project.MintingInfo.Index = nftProjectDetail.NftProjectDetail.Index.Int64()
			project.MintingInfo.IndexReverse = nftProjectDetail.NftProjectDetail.IndexReserve.Int64()
			u.Repo.UpdateProjectIndexData(project.UUID, project.MintingInfo.Index, project.MintingInfo.IndexReverse)
		}
	}*/
	creator, err := u.Repo.FindUserByWalletAddress(dataObject.CreatorAddr)
	if err != nil {
		logger.AtLog.Logger.Error("getTokenInfo", zap.String("tokenID", req.TokenID), zap.String("action", "FindUserByWalletAddress"), zap.Error(err))
		creator = &entity.Users{}
	}
	dataObject.Creator = creator
	mftMintedTime := <-mftMintedTimeChan

	if mftMintedTime.Err == nil {
		nft := mftMintedTime.NftMintedTime.Nft
		//onwer
		if nft.Owner != dataObject.OwnerAddr || (dataObject.Owner != nil && nft.Owner != dataObject.Owner.WalletAddress) {

			ownerAddr := strings.ToLower(nft.Owner)

			logger.AtLog.Logger.Info("dataObject.OwnerAddr.old", zap.String("tokenID", req.TokenID), zap.Any("dataObject.OwnerAddr", dataObject.OwnerAddr), zap.Any("ownerAddr", ownerAddr))
			owner, err := u.Repo.FindUserByWalletAddress(ownerAddr)
			if err != nil {
				logger.AtLog.Logger.Error("err", zap.Error(err))
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
		logger.AtLog.Logger.Error("u.GetNftMintedTime", zap.String("tokenID", req.TokenID), zap.Error(mftMintedTime.Err))
	}

	tokenFChan := <-tokendatachan
	if tokenFChan.Err == nil {
		dataObject.Name = tokenFChan.Data.Name
		if dataObject.Name == "" {
			dataObject.Name = dataObject.TokenID
		}
		dataObject.Description = tokenFChan.Data.Description

		dataObject.Thumbnail = tokenFChan.Data.Image
		if tokenFChan.Data.Image != "" {
			dataObject.Image = tokenFChan.Data.Image
		}

		if tokenFChan.Data.Thumbnail != "" {
			dataObject.Thumbnail = tokenFChan.Data.Thumbnail
		}

		if tokenFChan.Data.ParsedImage != nil && *tokenFChan.Data.ParsedImage != "" {
			dataObject.ParsedImage = tokenFChan.Data.ParsedImage
		}

		if tokenFChan.Data.ThumbnailCapturedAt != nil {
			dataObject.ThumbnailCapturedAt = tokenFChan.Data.ThumbnailCapturedAt
		}

		dataObject.AnimationURL = tokenFChan.Data.AnimationURL
		dataObject.Attributes = tokenFChan.Data.Attributes
		dataObject.Seed = tokenFChan.Data.Seed

		if len(tokenFChan.Data.ParsedAttributes) > 0 {
			dataObject.ParsedAttributes = tokenFChan.Data.ParsedAttributes
		}

		if len(tokenFChan.Data.ParsedAttributesStr) > 0 {
			dataObject.ParsedAttributesStr = tokenFChan.Data.ParsedAttributesStr
		}

	} else {
		logger.AtLog.Logger.Error("tokenFChan.Err", zap.String("tokenID", req.TokenID), zap.Error(tokenFChan.Err))

		if project.IsBigFile {
			dataObject.Thumbnail = project.Thumbnail
			jsonFile := project.ProcessingImages[0]
			fileName := strings.ReplaceAll(jsonFile, "https://cdn.generative.xyz/", "")
			bytes, err := u.GCS.ReadFile(fileName)

			if err == nil {
				tokeBFS := entity.TokenFromBase64{}
				err := json.Unmarshal(bytes, &tokeBFS)
				if err == nil {
					base64Image := tokeBFS.Image
					i := strings.Index(base64Image, ",")
					if i >= 0 {
						now := time.Now().UTC().Unix()
						name := fmt.Sprintf("thumb/%s-%d.png", dataObject.TokenID, now)
						base64Image = base64Image[i+1:]
						uploaded, err := u.GCS.UploadBaseToBucket(base64Image, name)
						if err != nil {
							logger.AtLog.Logger.Error("RunAndCap", zap.Any("tokenID", dataObject.TokenID), zap.Error(err))
						} else {
							logger.AtLog.Logger.Info("RunAndCap", zap.Any("tokenID", dataObject.TokenID), zap.Any("uploaded", uploaded))
							thumbnail := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), name)
							dataObject.Thumbnail = thumbnail
							isUpdated = true
						}
					}

				} else {
					logger.AtLog.Logger.Error("json.Unmarshal", zap.Error(err))
				}

			} else {
				logger.AtLog.Logger.Error("u.GCS.ReadFile", zap.String("fileName", fileName), zap.Error(err))
			}
		}
	}

	tokIdMini := dataObject.TokenIDInt % 100000
	dataObject.TokenIDMini = &tokIdMini

	logger.AtLog.Logger.Info(fmt.Sprintf("Data for minter address %v and OwnerAddr %v", dataObject.MinterAddress, dataObject.OwnerAddr), zap.String("tokenID", req.TokenID), zap.Bool("true", true))

	isAddress := func(s *string) bool {
		if s == nil {
			return false
		}
		return strings.HasPrefix(*s, "0x")
	}

	if dataObject.MinterAddress != nil {
		logger.AtLog.Logger.Info(fmt.Sprintf("Minter address %s", *dataObject.MinterAddress), zap.String("tokenID", req.TokenID), zap.Bool("true", true))
	}

	if !isAddress(dataObject.MinterAddress) && dataObject.OwnerAddr != "" {
		logger.AtLog.Logger.Info("Update minter address", zap.Any("true", true))
		dataObject.MinterAddress = &dataObject.OwnerAddr
		isUpdated = true
	}

	if isUpdated {
		if project.IsBigFile {
			dataObject.Thumbnail = project.Thumbnail
			// TODO:
			//linkToJson := project.ProcessingImages[0]
			//_ = linkToJson
			//tokeBFS := entity.TokenFromBase64{}
			//err = json.Unmarshal([]byte(linkToJson), &tokeBFS)
			// store google cloud
		}
		updated, err := u.Repo.UpdateOrInsertTokenUri(dataObject.ContractAddress, dataObject.TokenID, dataObject)
		if err != nil {
			logger.AtLog.Logger.Error("getTokenInfo", zap.Any("req", req), zap.String("action", "UpdateOrInsertTokenUri"), zap.Error(err))
			return nil, err
		}
		logger.AtLog.Logger.Info("getTokenInfo", zap.Any("req", req), zap.Any("updated", updated))
	}

	if !project.IsBigFile {
		//capture image
		payload := redis.PubSubPayload{Data: structure.TokenImagePayload{
			TokenID:         dataObject.TokenID,
			ContractAddress: dataObject.ContractAddress,
		}}

		err = u.PubSub.Producer(utils.PUBSUB_TOKEN_THUMBNAIL, payload)
		if err != nil {
			logger.AtLog.Logger.Error("getTokenInfo", zap.Any("req", req), zap.String("action", "u.PubSub.Producer"), zap.Error(err))
		}
	} else {
		u.Repo.CreateFragmentJob(context.TODO(), &entity.TokenFragmentJob{
			TokenId:  req.TokenID,
			FilePath: project.ProcessingImages[0],
		})
	}

	return dataObject, nil
}

func (u Usecase) getBFSData(client *ethclient.Client, bfsContract common.Address, gNft common.Address, seed string) (*string, error) {
	bfsC, err := bfs.NewBfs(bfsContract, client)
	if err != nil {
		return nil, err
	}
	/*value, err := bfsC.Count(nil, gNft, seed)
	if err != nil {
		return nil, err
	}*/

	var bytesArr []byte
	//if value.Cmp(big.NewInt(0)) > 0 {
	nextChunks := big.NewInt(0)
	for {
		bytes, nextChunks, err := bfsC.Load(nil, gNft, seed, nextChunks)
		if err != nil {
			return nil, err
		}
		bytesArr = append(bytesArr, bytes...)
		if nextChunks.Int64() == -1 {
			break
		}
	}
	//}

	if len(bytesArr) > 0 {
		result := "data:application/json;base64," + helpers.Base64Encode(bytesArr)
		return &result, nil
	}
	return nil, errors.New("Invalid bfs data")
}

func (u Usecase) getSeedFromTokenId(client *ethclient.Client, contractAddr common.Address, tokenIDStr string) (*string, error) {
	gNft, err := generative_nft_contract.NewGenerativeNftContract(contractAddr, client)
	if err != nil {
		return nil, err
	}
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(tokenIDStr, 10)
	if !ok {
		return nil, errors.New("cannot convert tokenID")
	}
	val, err := gNft.TokenIdToHash(nil, tokenID)
	if err != nil {
		return nil, err
	}
	result := "0x" + strings.ToUpper(hex.EncodeToString(val[:]))
	return &result, nil
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

	value, err := gNft.TokenURI(nil, tokenID)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (u Usecase) UpdateTokensFromChain() error {

	//TODO - we will use pagination instead of all
	tokens, err := u.Repo.GetAllTokens()
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return err
	}

	logger.AtLog.Logger.Info("tokens.Count", zap.Any("len(tokens)", len(tokens)))
	for _, token := range tokens {

		_, err := u.GetToken(structure.GetTokenMessageReq{ContractAddress: token.ContractAddress, TokenID: token.TokenID}, 5)
		if err != nil {
			logger.AtLog.Logger.Error("err", zap.Error(err))
			return err
		}
	}

	return nil
}

func (u Usecase) GetTokensByContract(contractAddress string, filter nfts.MoralisFilter) (*entity.Pagination, error) {

	client, err := helpers.EthDialer()
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	contractAddr := common.HexToAddress(contractAddress)
	gNft, err := generative_nft_contract.NewGenerativeNftContract(contractAddr, client)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	project, err := gNft.Project(nil)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}
	parentAddr := project.ProjectAddr

	resp, err := u.MoralisNft.GetNftByContract(contractAddress, filter)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}
	parentAddrStr := parentAddr.String()
	result := []entity.TokenUri{}
	for _, item := range resp.Result {
		tokenID := item.TokenID
		token, err := u.GetToken(structure.GetTokenMessageReq{ContractAddress: parentAddrStr, TokenID: tokenID}, 5)
		if err != nil {
			logger.AtLog.Logger.Error("err", zap.Error(err))
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
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	tokens, err := u.Repo.FilterTokenUri(*pe)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("tokens", zap.Any("tokens.Total", tokens.Total))
	return tokens, nil
}

func (u Usecase) FilterTokensNew(filter structure.FilterTokens) (*entity.Pagination, error) {
	pe := &entity.FilterTokenUris{}
	p, err := u.Repo.FindProjectByTokenIDOrGenNFTAddr(*filter.GenNFTAddr)
	if err != nil {
		logger.AtLog.Error("FilterTokensNew", zap.Any("filter", filter), zap.Error(err))
		return nil, err
	}

	isOrdinal := helpers.IsOrdinalProject(p.TokenId)

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

		groupTraits := make(map[string][]string)
		traits := p.TraitsStat
		for _, trait := range traits {
			values := trait.TraitValuesStat

			for _, value := range values {
				if value.Rarity >= int32(minInt) && value.Rarity <= int32(maxInt) {
					groupTraits[trait.TraitName] = append(groupTraits[trait.TraitName], value.Value)

				}
			}
		}

		for key, groupTrait := range groupTraits {
			r := structure.TokenUriAttrReq{}
			r.TraitType = key
			r.Values = groupTrait
			filter.RarityAttributes = append(filter.Attributes, r)
		}
	}

	err = copier.Copy(pe, filter)
	if err != nil {
		logger.AtLog.Error("FilterTokensNew", zap.Any("filter", filter), zap.Error(err))
		return nil, err
	}

	tokens := &entity.Pagination{}
	if isOrdinal {
		tokens, err = u.Repo.FilterTokenUriNew(*pe)
		if err != nil {
			logger.AtLog.Error("FilterTokensNew", zap.Any("filter", filter), zap.Error(err))
			return nil, err
		}
	} else {
		tokens, err = u.Repo.FilterTokenUriTCNew(*pe)
		if err != nil {
			logger.AtLog.Error("FilterTokensNew", zap.Any("filter", filter), zap.Error(err))
			return nil, err
		}
	}

	genService := generativeexplorer.NewGenerativeExplorer(u.Cache)
	ctx := context.Background()
	resp := []entity.TokenUriListingFilter{}
	for _, item := range tokens.Result.([]entity.TokenUriListingFilter) {
		if helpers.IsOrdinalProjectToken(item.TokenID) {
			iResp, err := genService.Inscription(item.TokenID)
			if err == nil && iResp != nil {
				item.OwnerAddress = iResp.Address
				if iResp.Address != item.Owner.WalletAddressBTCTaproot {
					item.Owner = entity.TokenURIListingOwner{
						WalletAddressBTCTaproot: iResp.Address,
						WalletAddress:           "",
						DisplayName:             "",
						Avatar:                  "",
					}
				}

				item.Royalty = item.Project.Royalty
				item.IsMinting = false
				item.MintingInfo = entity.MintingInfo{Done: 0, All: 0, Pending: 0}
				item.IsOnChain = true
			}
		} else {
			item.Owner = entity.TokenURIListingOwner{
				WalletAddressBTCTaproot: item.OnwerInternal.WalletAddressBTCTaproot,
				WalletAddress:           item.OnwerInternal.WalletAddress,
				DisplayName:             item.OnwerInternal.DisplayName,
				Avatar:                  item.OnwerInternal.Avatar,
			}

			item.PriceBRC20Obj = entity.PriceBRC20Obj{
				Value:      item.PriceBRC20,
				Address:    item.PriceBRC20Address,
				OfferingID: item.OfferingID,
			}

			item.Royalty = item.Project.Royalty
			item.MintingInfo = entity.MintingInfo{Done: 0, All: 0, Pending: 0}
			item.IsMinting = false

			mintingInfo, _ := u.Repo.AggregateMintingInfo(ctx, item.TokenID)
			if len(mintingInfo) >= 1 {
				mtinfo := mintingInfo[0]
				item.MintingInfo = entity.MintingInfo{
					All:     mtinfo.All,
					Done:    mtinfo.Done,
					Pending: mtinfo.Pending,
				}
				if mtinfo.Done < mtinfo.All {
					item.IsMinting = true
					item.IsOnChain = false
				} else {
					item.IsOnChain = true
				}
			} else {
				item.IsOnChain = true
			}

		}

		//spew.Dump(iResp)
		resp = append(resp, item)
	}

	tokens.Result = resp
	return tokens, nil
}

func (u Usecase) UpdateToken(req structure.UpdateTokenReq) (*entity.TokenUri, error) {

	p, err := u.Repo.FindTokenBy(req.ContracAddress, req.TokenID)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	if req.Priority != nil {
		p.Priority = req.Priority
	}

	updated, err := u.Repo.UpdateOrInsertTokenUri(req.ContracAddress, req.TokenID, p)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
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
		logger.AtLog.Logger.Error("err", zap.Error(err))
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

func (u Usecase) CreateBTCTokenURI(ownerAddress, projectID, tokenID, mintedURL string, paidType entity.TokenPaidType, opts ...string) (*entity.TokenUri, error) {

	// find project by projectID
	project, err := u.Repo.FindProjectByTokenID(projectID)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
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
	tokenUri.OwnerAddr = ownerAddress
	if len(opts) > 0 {
		tokenUri.NftTokenId = opts[0]
	}
	if len(opts) > 1 {
		tokenUri.InscribedBy = opts[1]
	}
	if len(opts) > 2 {
		tokenUri.OriginalInscribedBy = opts[2]
	}
	nftTokenUri := project.NftTokenUri
	logger.AtLog.Logger.Info("nftTokenUri", zap.Any("nftTokenUri", nftTokenUri))

	projectNftTokenUri := &structure.ProjectAnimationUrl{}
	err = helpers.Base64DecodeRaw(project.NftTokenUri, projectNftTokenUri)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	imageURI := ""
	if projectNftTokenUri.AnimationUrl != "" {
		logger.AtLog.Logger.Info("nftTokenUri", zap.Any("len(nftTokenUri)", len(nftTokenUri)))
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
		logger.AtLog.Logger.Info("mintedURL", zap.Any("mintedURL", mintedURL))
	}

	tokenUri.OrderInscriptionIndex = int(project.MintingInfo.Index + 1)
	_, err = u.Repo.UpdateOrInsertTokenUri(tokenUri.ContractAddress, tokenUri.TokenID, &tokenUri)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
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
		logger.AtLog.Logger.Error("err", zap.Error(err))
	}

	return pTokenUri, nil
}

func (u Usecase) TriggerPubsubTokenThumbnail(tokenId string) (*entity.TokenUri, error) {
	pTokenUri, err := u.Repo.FindTokenByTokenID(tokenId)
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
		logger.AtLog.Logger.Error("err", zap.Error(err))
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
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	if strings.Index(token.Image, ".glb") == -1 {
		err = errors.New("Token's image is not a glb file")
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}
	now := time.Now().Unix()

	base64Data := strings.ReplaceAll(req.Thumbnail, "data:image/png;base64,", "")
	uploaded, err := u.GCS.UploadBaseToBucket(base64Data, fmt.Sprintf("btc-projects/%s/thumb/token-%s-%d.png", token.ProjectID, token.TokenID, now))
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}
	logger.AtLog.Logger.Info("uploaded", zap.Any("uploaded", uploaded))
	thumb := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	//spew.Dump(thumb)
	token.Thumbnail = thumb

	updated, err := u.Repo.UpdateOrInsertTokenUri(token.ContractAddress, token.TokenID, token)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
	return token, nil
}

// When go to this, you need to make sure that meta's project is created
func (u Usecase) CreateBTCTokenURIFromCollectionInscription(meta entity.CollectionMeta, inscription entity.CollectionInscription) (*entity.TokenUri, error) {
	// find project by projectID
	project, err := u.Repo.FindProjectByTokenID(meta.ProjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.AtLog.Logger.Error("CanNotFindProjectByInscriptionIcon", zap.Any("inscriptionIcon", meta.InscriptionIcon))
			return nil, repository.ErrNoProjectsFound
		} else {
			return nil, err
		}
	}

	tokenUri := entity.TokenUri{}
	tokenUri.ContractAddress = os.Getenv("GENERATIVE_BTC_PROJECT")
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
	logger.AtLog.Logger.Info("nftTokenUri", zap.Any("nftTokenUri", nftTokenUri))

	projectNftTokenUri := &structure.ProjectAnimationUrl{}
	err = helpers.Base64DecodeRaw(project.NftTokenUri, projectNftTokenUri)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	now := time.Now().UTC()
	imageURI := fmt.Sprintf("https://generativeexplorer.com/content/%s", inscription.ID)

	resp := &entity.InscriptionDetail{}
	_, err = resty.New().R().
		SetResult(&resp).
		Get(fmt.Sprintf("%s/inscription/%s", u.Config.GenerativeExplorerApi, inscription.ID))
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	if strings.Contains(resp.ContentType, "text/html") {
		tokenUri.AnimationURL = imageURI
	}

	tokenUri.Thumbnail = imageURI
	tokenUri.Image = imageURI
	tokenUri.ParsedImage = &imageURI
	tokenUri.ThumbnailCapturedAt = &now
	tokenUri.Source = inscription.Source
	logger.AtLog.Logger.Info("mintedURL", zap.Any("imageURI", imageURI))

	_, err = u.Repo.UpdateOrInsertTokenUri(tokenUri.ContractAddress, tokenUri.TokenID, &tokenUri)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	pTokenUri, err := u.Repo.FindTokenBy(tokenUri.ContractAddress, tokenUri.TokenID)
	if err != nil {
		return nil, err
	}

	// update project index and max supply
	index := project.MintingInfo.Index + 1
	maxSupply := project.MaxSupply
	if index > maxSupply {
		maxSupply = index
	}
	err = u.Repo.UpdateProjectIndexAndMaxSupply(project.TokenID, maxSupply, index)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	go u.TriggerPubsubTokenThumbnail(pTokenUri.TokenID)

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
	//spew.Dump(link)
	return &link, nil

}

func (u Usecase) GetTokensMap(tokenIDs []string) (map[string]entity.TokenUri, error) {
	tokens, err := u.Repo.FindTokenByTokenIds(tokenIDs)
	if err != nil {
		return nil, err
	}

	tokenIdToToken := map[string]entity.TokenUri{}
	for id, token := range tokens {
		tokenIdToToken[tokens[id].TokenID] = token
	}
	return tokenIdToToken, nil
}

type tokenOwner struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Count   int    `json:"count"`
}

func (u Usecase) AnalyticsTokenUriOwner(f structure.FilterTokens) (interface{}, error) {
	filter := &entity.FilterTokenUris{}
	err := copier.Copy(filter, f)
	if err != nil {
		return nil, err
	}

	owners := make(map[string]*tokenOwner)
	tokenIDs := make(map[string]*string)
	tokens, err := u.Repo.AnalyticsTokenUriOwner(*filter)
	if err != nil {
		return nil, err
	}

	genService := generativeexplorer.NewGenerativeExplorer(u.Cache)
	for _, token := range tokens {
		tokenID := token.TokenID
		tokenIDs[tokenID] = &token.OwnerAddress

		if helpers.IsOrdinalProject(token.TokenID) {
			var address, name, avatar string

			iResp, _ := genService.Inscription(tokenID)
			if iResp != nil {
				address = iResp.Address
			} else {
				if token.Owner.WalletAddressBTCTaproot != "" {
					user, _ := u.Repo.FindUserByBtcAddressTaproot(token.Owner.WalletAddressBTCTaproot)
					if user != nil {
						name = user.DisplayName
						avatar = user.Avatar
						address = user.WalletAddressBTCTaproot
					}
				} else if token.OwnerAddress != "" {
					user, _ := u.Repo.FindUserByWalletAddress(token.OwnerAddress)
					if user != nil {
						name = user.DisplayName
						avatar = user.Avatar
						address = user.WalletAddressBTCTaproot
					}
				} else {
					continue
				}
			}

			if _, has := owners[address]; !has {
				owners[address] = &tokenOwner{
					Address: address,
					Name:    name,
					Avatar:  avatar,
					Count:   0,
				}
			}
			owners[address].Count++
		}
	}

	var result []*tokenOwner
	for _, v := range owners {
		result = append(result, v)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result, nil
}

func (u Usecase) GetTokenMintingInfo(tokenID string) ([]repository.AggregateTokenMintingInfo, error) {

	tokenID = strings.ToLower(tokenID)
	info, err := u.Repo.AggregateMintingInfo(context.Background(), tokenID)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	///logger.AtLog.Logger.Info("tokenUri", zap.Any("tokenUri", tokenUri))
	logger.AtLog.Logger.Info("tokenID", zap.Any("tokenID", tokenID), zap.Any("info", info))
	return info, nil
}

func (u Usecase) InsertDataTokenUriFromNFT(jsonFilepath string, projectId string) {
	// uc.InsertDataTokenUriFromNFT("nfts.json", "999998")
	jsonFile, err := os.Open(jsonFilepath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	nfts := []entity.NFT{}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &nfts)
	inputChan := make(chan entity.NFT)
	outputChan := make(chan *entity.TokenUri)

	workerNumber := 100
	if err == nil {
		contractAddress := "0xda00b6a8b521113501bb98fd0a7ffcfe756d9962"

		//init worker
		for i := 0; i < len(nfts); i++ {
			go u.InsertDataTokenWoker(inputChan, outputChan, projectId, contractAddress)
			if i > 0 && i%workerNumber == 0 {
				time.Sleep(time.Millisecond * 500)
			}
		}

		for i := 0; i < len(nfts); i++ {
			inputChan <- nfts[i]
		}

		for i := 0; i < len(nfts); i++ {
			item := <-outputChan
			u.Repo.UpdateOrInsertTokenUri(contractAddress, item.TokenID, item)
		}
	}
}

func (u Usecase) InsertDataTokenWoker(nftInput chan entity.NFT, outputChan chan *entity.TokenUri, projectId string, contractAddress string) {
	animationURL := ""
	isPubsub := false
	nft := <-nftInput
	if nft.Metadata != nil {
		if nft.Metadata.AnimationURL != nil {
			animationURL = *nft.Metadata.AnimationURL
			isPubsub = true
		}
	} else {
		filename := fmt.Sprintf("metadata-%s-%s.txt", nft.CollectionAddress, nft.TokenId)
		data, err := helpers.Openfile(filename)
		if err != nil {
			url := fmt.Sprintf("https://api-nft-explorer.trustless.computer/api/v1/collection/%s/nft/%s/tokenuri", nft.CollectionAddress, nft.TokenId)

			statusCode, body, err := request.GetRequest(url)
			if err != nil {
				logger.AtLog.Error("PubSubCreateTokenThumbnail - GetRequest", zap.String("url", url), zap.Error(err))
				return
			}

			if statusCode != 200 {
				err := errors.New(fmt.Sprintf("Status code: %s", statusCode))
				logger.AtLog.Error("PubSubCreateTokenThumbnail - statusCode", zap.String("url", url), zap.Error(err))
				return
			}

			metadata := &entity.NFTMetadata{}
			err = json.Unmarshal(body, metadata)
			if err != nil {
				logger.AtLog.Error("PubSubCreateTokenThumbnail - Unmarshal", zap.String("url", url), zap.Error(err))
				return
			}
			if metadata.AnimationURL != nil && *metadata.AnimationURL != "" {
				animationURL = *metadata.AnimationURL
				isPubsub = true

				err = helpers.CreateFile(filename, *metadata.AnimationURL)
				if err != nil {
					logger.AtLog.Error("PubSubCreateTokenThumbnail - CreateFile", zap.String("url", url), zap.Error(err))
					return
				}
			}
		}

		data, err = helpers.Openfile(filename)
		if err != nil {
			logger.AtLog.Error("PubSubCreateTokenThumbnail - Unmarshal", zap.String("filename", filename), zap.Error(err))
			return
		}
		//reopen file
		base64str := bytes.NewBuffer(data).String()
		base64str = strings.ReplaceAll(base64str, "data:text/html;base64,", "")
		base64str = strings.ReplaceAll(base64str, `"`, ``)
		bytesData, err := helpers.Base64Decode(base64str)
		if err == nil {
			html := bytes.NewBuffer(bytesData).String()
			html = strings.ReplaceAll(html, `Web3.givenProvider`, `"https://tc-node.trustless.computer"`)
			html = strings.ReplaceAll(html, `isFakeData=!1`, `isFakeData=1`)
			html = strings.ReplaceAll(html, `document.getElementById("btn-close").style.display="none",console.log(window.ethereum),window.ethereum`, `document.getElementById("btn-close").style.display="none",console.log(window.ethereum),window.ethereum || isFakeData`)
			bytesData := []byte(html)
			err := helpers.CreateTxtFile(fmt.Sprintf("index-%s-%s.html", nft.CollectionAddress, nft.TokenId), bytesData)
			if err == nil {
				now := time.Now().UTC().UnixNano()
				encoded := helpers.Base64Encode(bytesData)
				uploaded, err := u.GCS.UploadBaseToBucket(encoded, fmt.Sprintf("index-%d-%s-%s.html", now, nft.CollectionAddress, nft.TokenId))

				if err == nil {
					fileURI := fmt.Sprintf("%s/%s?seed=%s", os.Getenv("GCS_DOMAIN"), uploaded.Name, nft.TokenId)

					animationURL = fileURI
					isPubsub = true
				}

			}
		}
	}

	fmt.Sprintf("Inserting %s", nft.Name)
	owner, _ := u.Repo.FindUserByAddress(nft.Owner)
	project, _ := u.Repo.FindProjectByTokenID(projectId)
	item := &entity.TokenUri{
		ContractAddress:     contractAddress,
		GenNFTAddr:          nft.CollectionAddress,
		OwnerAddr:           nft.Owner,
		MinterAddress:       &nft.Owner,
		Owner:               owner,
		ProjectID:           projectId,
		Project:             project,
		ParsedAttributesStr: nft.Attributes,
		//Name:                nft.Name
	}

	tempProjectId, _ := new(big.Int).SetString(projectId, 10)
	item.ProjectIDInt = tempProjectId.Int64()
	tempTokenId, _ := new(big.Int).SetString(nft.TokenId, 10)
	fakeTokenId := item.ProjectIDInt*1000000 + tempTokenId.Int64()
	item.TokenIDInt = int(fakeTokenId)
	item.TokenID = fmt.Sprintf("%d", fakeTokenId)

	if !isPubsub {
		var buf []byte
		resp, _ := http.Get(nft.TokenUri)
		buf, _ = ioutil.ReadAll(resp.Body)
		image := helpers.Base64Encode(buf)
		image = fmt.Sprintf("%s,%s", "data:image/svg+xml;base64", image)

		thumbnail := ""
		if image != "" {
			base64Image := image
			i := strings.Index(base64Image, ",")
			if i >= 0 {
				now := time.Now().UTC().Unix()

				name := fmt.Sprintf("thumb/%d-%d.svg", fakeTokenId, now)
				base64Image = base64Image[i+1:]
				uploaded, err := u.GCS.UploadBaseToBucket(base64Image, name)
				if err != nil {
					logger.AtLog.Logger.Error("InsertDataTokenUriFromNFT", zap.Any("tokenID", fakeTokenId), zap.Error(err))
				} else {
					logger.AtLog.Logger.Info("InsertDataTokenUriFromNFT", zap.Any("tokenID", fakeTokenId), zap.Any("uploaded", uploaded))
					thumbnail = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), name)
				}
			}
		}
		item.ParsedImage = &thumbnail
		item.Thumbnail = thumbnail
		item.Image = thumbnail

	} else {
		item.AnimationURL = animationURL

		//open chrome and capture image
		resp, err := u.RunAndCap(item)
		if err != nil {
			logger.AtLog.Error("PubSubCreateTokenThumbnai", zap.Any("RunAndCap", zap.Error(err)))
			return
		}

		if resp.IsUpdated {
			item.ParsedImage = &resp.ParsedImage
			item.Thumbnail = resp.Thumbnail
			item.ParsedAttributes = resp.Traits
			item.ParsedAttributesStr = resp.TraitsStr
			item.ThumbnailCapturedAt = resp.CapturedAt
		}
	}

	outputChan <- item

}

func (u Usecase) PreviewTokenByTokenID(tokenID string) (string, error) {
	tokenUri, err := u.Repo.FindTokenByTokenID(tokenID)
	if err != nil {
		return "", err
	}

	animationURL := tokenUri.AnimationURL
	_b, _, _, err := helpers.HttpRequest(animationURL, "GET", map[string]string{}, nil)
	if err != nil {
		return "", err
	}

	str := string(_b)

	doc, err2 := htmlquery.Parse(strings.NewReader(str))
	if err2 != nil {
		err = err2
		return "", err
	}

	f := "//script"
	found := u.FindNode(doc, f)

	dataS := ""
	scriptContents := make(map[int]string)
	for k, i := range found {
		for _, attr := range i.Attr {
			if attr.Key == "src" {
				scriptUrl := fmt.Sprintf("https://ordinals.com%s", attr.Val)
				_bS, _, _, err := helpers.HttpRequest(scriptUrl, "GET", map[string]string{}, nil)
				if err != nil {
					return "", err
				}

				scriptContent := string(_bS)
				scriptContent = strings.ReplaceAll(scriptContent, "/content", "https://ordinals.com/content")
				scriptContents[k] = scriptContent
			}

			if attr.Key == "data-s" {
				dataS = attr.Val
			}
		}

	}

	found1 := htmlquery.Find(doc, f)
	for k, node := range found1 {
		newAttrs := []html.Attribute{}

		for _, attr := range node.Attr {
			if attr.Key == "src" {
				continue
			}
			newAttrs = append(newAttrs, attr)
		}

		node.Attr = newAttrs
		//node.FirstChild = &html.Node{}
		//node.FirstChild.Data = scriptContents[k]

		spew.Dump(node.Attr)
		spew.Dump(k)
	}

	var b bytes.Buffer
	err = html.Render(&b, doc)
	if err != nil {
		return "", err
	}

	str1 := b.String()
	str1 = strings.ReplaceAll(str1, fmt.Sprintf(`<script data-s="%s"></script>`, dataS), fmt.Sprintf(`<script data-s="%s">%s</script>`, dataS, scriptContents[0]))

	return str1, nil

}

func (u Usecase) LoadContent(tokenID string) (string, error) {
	_b, _, _, err := helpers.HttpRequest(fmt.Sprintf(`https://ordinals.com/content/%s`, tokenID), "GET", map[string]string{}, nil)
	if err != nil {
		return "", err
	}

	return string(_b), nil

}

func (u *Usecase) findDocByName(doc *html.Node, name string) []string {
	content := []string{}
	f := "//script"
	found := htmlquery.Find(doc, f)
	for _, f1 := range found {
		for _, f2 := range f1.Attr {
			if f2.Key == "src" {
				content1 := f2.Val
				if content1 != "" {
					content = append(content, content1)
				}
			}
		}
	}

	return content

}

func (u *Usecase) FindNode(doc *html.Node, find string) []*html.Node {
	found := htmlquery.Find(doc, find)
	if len(found) > 0 {
		return found
	}

	if doc.FirstChild == nil {
		return []*html.Node{}
	}

	if doc.LastChild == nil {
		return []*html.Node{}
	}

	res := []*html.Node{}
	fC := u.FindNode(doc.FirstChild, find)
	lC := u.FindNode(doc.LastChild, find)

	res = append(res, fC...)
	res = append(res, lC...)

	return res
}

func (u Usecase) elementScreenshot(sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
