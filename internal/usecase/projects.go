package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/copier"

	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/contracts/generative_nft_contract"
	"rederinghub.io/utils/contracts/generative_project_contract"
	discordclient "rederinghub.io/utils/discord"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/redis"
)

type uploadFileChan struct {
	FileURL *string
	Err     error
}

func (u Usecase) CreateProject(req structure.CreateProjectReq) (*entity.Projects, error) {

	pe := &entity.Projects{}
	err := copier.Copy(pe, req)
	if err != nil {
		u.Logger.ErrorAny("CreateProject", zap.Any("err", err))
		return nil, err
	}

	err = u.Repo.CreateProject(pe)
	if err != nil {
		u.Logger.ErrorAny("CreateProject", zap.Any("err", err))
		return nil, err
	}

	u.Logger.ErrorAny("CreateProject", zap.Any("project", pe))
	return pe, nil
}

func (u Usecase) networkFeeBySize(size int64) int64 {
	response, err := http.Get("https://mempool.space/api/v1/fees/recommended")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	feeRateValue := int64(entity.DEFAULT_FEE_RATE)
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		u.Logger.Error(err)
		return 0
	} else {
		type feeRate struct {
			fastestFee  int
			halfHourFee int
			hourFee     int
			economyFee  int
			minimumFee  int
		}

		feeRateObj := feeRate{}

		err = json.Unmarshal(responseData, &feeRateObj)
		if err != nil {
			u.Logger.Error(err)
			return 0
		}
		if feeRateObj.fastestFee > 0 {
			feeRateValue = int64(feeRateObj.fastestFee)
		}
	}

	return size * feeRateValue

}

func (u Usecase) CreateBTCProject(req structure.CreateBtcProjectReq) (*entity.Projects, error) {
	u.Logger.LogAny("CreateBTCProject", zap.Any("CreateBtcProjectReq", req))

	pe := &entity.Projects{}
	err := copier.Copy(pe, req)
	if err != nil {
		u.Logger.ErrorAny("CreateBTCProject", zap.Any("copier.Copy", err))
		return nil, err
	}

	mPrice := helpers.StringToBTCAmount(req.MintPrice)
	maxID, err := u.Repo.GetMaxBtcProjectID()
	if err != nil {
		u.Logger.ErrorAny("CreateBTCProject", zap.Any("err.GetMaxBtcProjectID", err))
		return nil, err
	}
	maxID = maxID + 1
	pe.TokenIDInt = maxID
	pe.TokenID = fmt.Sprintf("%d", maxID)
	pe.ContractAddress = os.Getenv("GENERATIVE_PROJECT")
	pe.MintPrice = mPrice.String()
	pe.NetworkFee = big.NewInt(u.networkFeeBySize(int64(300000 / 4))).String() // will update after unzip and check data or check from animation url
	pe.IsHidden = false
	pe.Status = true
	pe.IsSynced = true
	nftTokenURI := make(map[string]interface{})
	nftTokenURI["name"] = pe.Name
	nftTokenURI["description"] = pe.Description
	nftTokenURI["image"] = pe.Thumbnail
	nftTokenURI["animation_url"] = ""
	nftTokenURI["attributes"] = []string{}
	creatorAddrr, err := u.Repo.FindUserByWalletAddress(req.CreatorAddrr)
	if err != nil {
		u.Logger.ErrorAny("CreateBTCProject", zap.Any("err.FindUserByWalletAddress", err))
		return nil, err
	}

	if creatorAddrr.WalletAddressBTC == "" {
		creatorAddrr.WalletAddressBTC = req.CreatorAddrrBTC
		updated, err := u.Repo.UpdateUserByID(creatorAddrr.UUID, creatorAddrr)
		if err != nil {
			u.Logger.ErrorAny("CreateBTCProject", zap.Any("err.UpdateUserByID", err))

		} else {
			u.Logger.Info("updated.creatorAddrr", creatorAddrr)
			u.Logger.Info("updated", updated)
		}
	}

	isPubsub := false
	animationURL := ""
	zipLink := req.ZipLink
	if zipLink != nil && *zipLink != "" {
		pe.IsHidden = true
		isPubsub = true
		pe.Status = false
	} else {
		if req.AnimationURL != nil {
			animationURL = *req.AnimationURL
			maxSize := helpers.CalcOrigBinaryLength(animationURL)
			pe.NetworkFee = big.NewInt(u.networkFeeBySize(int64(maxSize / 4))).String()
			htmlContent, err := helpers.Base64Decode(strings.ReplaceAll(animationURL, "data:text/html;base64,", ""))
			if err == nil {
				isFullChain, err := helpers.IsFullChain(string(htmlContent))
				if err == nil {
					pe.IsFullChain = isFullChain
					u.Logger.LogAny("CreateBTCProject", zap.Any("isFullChain", isFullChain))
				} else {
					u.Logger.ErrorAny("CreateBTCProject", zap.Any("isFullChain", err))
				}
			} else {
				u.Logger.ErrorAny("CreateBTCProject", zap.Any("isFullChain", err))
			}
			nftTokenURI["animation_url"] = animationURL
		}
	}

	bytes, err := json.Marshal(nftTokenURI)
	if err != nil {
		u.Logger.ErrorAny("CreateBTCProject", zap.Any("marshal", err))
		return nil, err
	}
	nftToken := helpers.Base64Encode(bytes)
	now := time.Now().UTC()

	pe.NftTokenUri = fmt.Sprintf("data:application/json;base64,%s", nftToken)
	pe.ProcessingImages = []string{}
	pe.MintedImages = nil
	pe.MintedTime = &now
	pe.CreatorProfile = *creatorAddrr
	pe.CreatorAddrrBTC = req.CreatorAddrrBTC
	pe.LimitSupply = 0
	pe.GenNFTAddr = pe.TokenID
	if len(req.Categories) != 0 {
		pe.Categories = []string{req.Categories[0]}
	}

	if pe.Categories == nil || len(pe.Categories) == 0 {
		pe.Categories = []string{u.Config.OtherCategoryID}
	}

	u.Logger.LogAny("CreateBTCProject", zap.Any("project", pe))
	err = u.Repo.CreateProject(pe)
	if err != nil {
		u.Logger.ErrorAny("CreateBTCProject", zap.Any("CreateProject", err))
		return nil, err
	}

	if isPubsub {
		err = u.PubSub.Producer(utils.PUBSUB_PROJECT_UNZIP, redis.PubSubPayload{Data: structure.ProjectUnzipPayload{ProjectID: pe.TokenID, ZipLink: *zipLink}})
		if err != nil {
			u.Logger.Error("u.Repo.CreateProject", err.Error(), err)
			//return nil, err
		}
	}

	u.NotifyWithChannel(os.Getenv("SLACK_PROJECT_CHANNEL_ID"), fmt.Sprintf("[Project is created][projectID %s]", pe.TokenID), fmt.Sprintf("TraceID: %s", pe.TraceID), fmt.Sprintf("Project %s has been created by user %s", pe.Name, pe.CreatorAddrr))
	u.NotifyCreateNewProjectToDiscord(pe, creatorAddrr)

	return pe, nil
}

func (u Usecase) NotifyCreateNewProjectToDiscord(pe *entity.Projects, creatorAddrr *entity.Users) {
	domain := os.Getenv("DOMAIN")
	webhook := os.Getenv("DISCORD_NEW_PROJECT_WEBHOOK")
	fields := make([]discordclient.Field, 0)
	addFields := func(fields []discordclient.Field, name string, value string) []discordclient.Field {
		if value == "" {
			return fields
		}
		return append(fields, discordclient.Field{
			Name:  name,
			Value: value,
		})
	}

	fields = addFields(fields, "Mint Price", u.resolveMintPriceBTC(pe.MintPrice))
	fields = addFields(fields, "Max Supply", fmt.Sprintf("%d", pe.MaxSupply))
	discordMsg := discordclient.Message{
		Username: "Satoshi 27",
		Embeds: []discordclient.Embed{{
			Title:       fmt.Sprintf("just launched %s", pe.Name),
			Url:         fmt.Sprintf("%s/generative/%s", domain, pe.GenNFTAddr),
			Description: pe.Description,
			Author: discordclient.Author{
				Name:    u.resolveShortName(creatorAddrr.DisplayName, creatorAddrr.WalletAddress),
				Url:     fmt.Sprintf("%s/profile/%s", domain, creatorAddrr.WalletAddress),
				IconUrl: creatorAddrr.Avatar,
			},
			Fields: fields,
			Image: discordclient.Image{
				Url: pe.Thumbnail,
			},
		}},
	}
	sendCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	u.Logger.Info("sending message to discord", discordMsg)

	if err := u.DiscordClient.SendMessage(sendCtx, webhook, discordMsg); err != nil {
		u.Logger.Error("error sending message to discord", err)
	}
}

func (u Usecase) resolveMintPriceBTC(priceStr string) string {
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return priceStr
	}
	return strconv.FormatFloat(float64(price)/1e8, 'f', -1, 64) + " BTC"

}

func (u Usecase) resolveShortName(userName string, userAddr string) string {
	if userName != "" {
		return userName
	}

	return userAddr[:4] + "..." + userAddr[len(userAddr)-4:]
}

func (u Usecase) UpdateBTCProject(req structure.UpdateBTCProjectReq) (*entity.Projects, error) {

	if req.ProjectID == nil {
		err := errors.New("ProjectID is requeried")
		u.Logger.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	if req.CreatetorAddress == nil {
		err := errors.New("CreatorAddress is requeried")
		u.Logger.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	p, err := u.Repo.FindProjectByTokenID(*req.ProjectID)
	if err != nil {
		u.Logger.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	if strings.ToLower(p.CreatorAddrr) != strings.ToLower(*req.CreatetorAddress) {
		err := errors.New("Only owner can update this project")
		u.Logger.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	if req.Name != nil && *req.Name != "" {
		p.Name = *req.Name
	}

	if req.Description != nil && *req.Description != "" {
		p.Description = *req.Description
	}

	if req.Thumbnail != nil && *req.Thumbnail != "" {
		p.Thumbnail = *req.Thumbnail
	}

	if req.IsHidden != nil && *req.IsHidden != p.IsHidden {
		p.IsHidden = *req.IsHidden
	}

	if len(req.Categories) > 0 {
		p.Categories = []string{req.Categories[0]}
	}

	if req.MaxSupply != nil && *req.MaxSupply != 0 && *req.MaxSupply != p.MaxSupply {
		// if p.MintingInfo.Index > 0 {
		// 	err := errors.New("Project is minted, cannot update max supply")
		// 	u.Logger.Error("pjID.minted", err.Error(), err)
		// 	return nil, err
		// }

		p.MaxSupply = *req.MaxSupply
	}

	if req.Royalty != nil {
		// if *req.Royalty > 2500 {
		// 	err := errors.New("Royalty must be less than 25")
		// 	u.Logger.Error("pjID.empty", err.Error(), err)
		// 	return nil, err
		// }

		// if *req.Royalty != p.Royalty && p.MintingInfo.Index > 0 {
		// 	err := errors.New("Project is minted, cannot update max supply")
		// 	u.Logger.Error("pjID.minted", err.Error(), err)
		// 	return nil, err
		// }

		p.Royalty = *req.Royalty
	}

	if req.MintPrice != nil {
		// mFStr := p.MintPrice
		reqMfFStr := helpers.StringToBTCAmount(*req.MintPrice)
		// if p.MintingInfo.Index > 0 && mFStr != reqMfFStr.String() {
		// 	err := errors.New("Project is minted, cannot update mint price")
		// 	u.Logger.Error("pjID.minted", err.Error(), err)
		// 	return nil, err
		// }
		p.MintPrice = reqMfFStr.String()
	}

	updated, err := u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		u.Logger.Error("updated", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("updated", updated)
	return p, nil
}

func (u Usecase) DeleteBTCProject(req structure.UpdateBTCProjectReq) (*entity.Projects, error) {

	p, err := u.Repo.FindProjectByTokenID(*req.ProjectID)
	if err != nil {
		u.Logger.ErrorAny("DeleteProject", zap.Any("err.FindProjectBy", err))
		return nil, err
	}
	whitelist := make(map[string]bool)
	whitelist["0xe23fcb129d6ea1b847202b14a56f957e5a464f64"] = true // andy
	whitelist["0x668ea0470396138acd0b9ccf6fbdb8a845b717b0"] = true // thaibao
	whitelist["0xe55eade1b17bba28a80a71633af8c15dc2d556a5"] = true // thaibao
	whitelist["0x9ef2cf140a51f87d266121409304399f0d93820f"] = true // ken
	whitelist["0xe10db08ab370eb3173ad8b0396a63f3af010364d"] = true // della
	whitelist["0xd77f54424cc2bd2a7315b1018e53548f62f690c0"] = true // anne
	if strings.ToLower(p.CreatorAddrr) != strings.ToLower(*req.CreatetorAddress) && !whitelist[strings.ToLower(*req.CreatetorAddress)] {
		u.Logger.ErrorAny("DeleteProject", zap.Any("err.CreatorAddrr", err))
		return nil, err
	}

	p.IsSynced = false
	p.Status = false
	p.IsHidden = true

	updated, err := u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		u.Logger.ErrorAny("UpdateProject", zap.Any("err.UpdateProject", err))
		return nil, err
	}

	u.Logger.Info("updated", updated)
	u.Logger.LogAny("UpdateProject", zap.Any("project", p))
	return p, nil
}

func (u Usecase) SetCategoriesForBTCProject(req structure.UpdateBTCProjectReq) (*entity.Projects, error) {

	if req.ProjectID == nil {
		err := errors.New("ProjectID is requeried")
		u.Logger.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	p, err := u.Repo.FindProjectByTokenID(*req.ProjectID)
	if err != nil {
		u.Logger.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	if len(req.Categories) > 0 {
		p.Categories = []string{req.Categories[0]}
	}

	updated, err := u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		u.Logger.Error("updated", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("updated", updated)
	return p, nil
}

func (u Usecase) UpdateProject(req structure.UpdateProjectReq) (*entity.Projects, error) {
	p, err := u.Repo.FindProjectBy(req.ContracAddress, req.TokenID)
	if err != nil {
		u.Logger.ErrorAny("UpdateProject", zap.Any("err.FindProjectBy", err))
		return nil, err
	}

	if req.Priority != nil {
		priority := 0
		p.Priority = &priority
	}

	if len(p.ReportUsers) >= u.Config.MaxReportCount {
		p.IsHidden = true
	}
	updated, err := u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		u.Logger.ErrorAny("UpdateProject", zap.Any("err.UpdateProject", err))
		return nil, err
	}

	u.Logger.Info("updated", updated)
	u.Logger.LogAny("UpdateProject", zap.Any("project", p))
	return p, nil
}

func (u Usecase) ReportProject(tokenId, iWalletAddress, originalLink string) (*entity.Projects, error) {
	p, err := u.Repo.FindProjectByTokenID(tokenId)
	if err != nil {
		u.Logger.Error("ReportProject.FindProjectBy", err.Error(), err)
		return nil, err
	}

	for _, r := range p.ReportUsers {
		if r.ReportUserAddress == iWalletAddress {
			return nil, errors.New("You have already reported before.")
		}
	}

	rep := &entity.ReportProject{
		ReportUserAddress: iWalletAddress,
		OriginalLink:      originalLink,
	}

	p.ReportUsers = append(p.ReportUsers, rep)
	if len(p.ReportUsers) >= u.Config.MaxReportCount {
		p.IsHidden = true
	}
	updated, err := u.Repo.UpdateProject(p.UUID, p)

	if err != nil {
		u.Logger.Error("UpdateProject.ReportProject", err.Error(), err)
		return nil, err
	}
	u.Logger.Info("updated", updated)

	return p, nil
}

func (u Usecase) GetProjectByGenNFTAddr(genNFTAddr string) (*entity.Projects, error) {
	project, err := u.Repo.FindProjectByGenNFTAddr(genNFTAddr)
	return project, err
}

func (u Usecase) GetProjects(req structure.FilterProjects) (*entity.Pagination, error) {

	pe := &entity.FilterProjects{}
	err := copier.Copy(pe, req)
	if err != nil {
		u.Logger.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	projects, err := u.Repo.GetProjects(*pe)
	if err != nil {
		u.Logger.Error("u.Repo.GetProjects", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("projects", projects.Total)
	return projects, nil
}

func (u Usecase) GetRandomProject() (*entity.Projects, error) {

	caddr := os.Getenv("RANDOM_PR_CONTRACT")
	pID := os.Getenv("RANDOM_PR_PROJECT")

	if caddr != "" && pID != "" {
		return u.GetProjectDetail(structure.GetProjectDetailMessageReq{
			ContractAddress: caddr,
			ProjectID:       pID,
		})
	}

	key := helpers.ProjectRandomKey()

	//always reload data
	go func() {
		p, err := u.Repo.GetAllProjects(entity.FilterProjects{})
		if err != nil {
			return
		}
		u.Cache.SetData(key, p)
	}()

	cached, err := u.Cache.GetData(key)
	if err != nil {
		p, err := u.Repo.GetAllProjects(entity.FilterProjects{})
		if err != nil {
			u.Logger.Error("u.Repo.GetProjects", err.Error(), err)
			return nil, err
		}
		u.Cache.SetData(key, p)
	}

	cached, err = u.Cache.GetData(key)
	projects := []entity.Projects{}
	bytes := []byte(*cached)
	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		u.Logger.Error("json.Unmarshal", err.Error(), err)
		return nil, err
	}

	if len(projects) == 0 {
		err := errors.New("Project are not found")
		u.Logger.Error("Projects.are.not.found", err.Error(), err)
		return nil, err
	}

	timeNow := time.Now().UTC().Nanosecond()
	rand := int(timeNow) % len(projects)

	//TODO - cache will be applied here

	projectRand := projects[rand]
	return u.GetProjectDetail(structure.GetProjectDetailMessageReq{
		ContractAddress: projectRand.ContractAddress,
		ProjectID:       projectRand.TokenID,
	})
}

func (u Usecase) GetMintedOutProjects(req structure.FilterProjects) (*entity.Pagination, error) {

	pe := &entity.FilterProjects{}
	err := copier.Copy(pe, req)
	if err != nil {
		u.Logger.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	pe.WalletAddress = req.WalletAddress
	projects, err := u.Repo.GetMintedOutProjects(*pe)
	if err != nil {
		u.Logger.Error("u.Repo.GetMintedOutProjects", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("projects", projects.Total)
	return projects, nil
}

func (u Usecase) GetProjectDetail(req structure.GetProjectDetailMessageReq) (*entity.Projects, error) {
	u.Logger.LogAny("GetProjectDetail", zap.Any("req", req))
	c, _ := u.Repo.FindProjectWithoutCache(req.ContractAddress, req.ProjectID)
	if (c == nil) || (c != nil && !c.IsSynced) || c.MintedTime == nil {
		// p, err := u.UpdateProjectFromChain(req.ContractAddress, req.ProjectID)
		// if err != nil {
		// 	u.Logger.Error("u.Repo.FindProjectBy", err.Error(), err)
		// 	return nil, err
		// }
		// return p, nil
		return nil, errors.New("project is not found")
	}
	mintPriceInt, err := strconv.ParseInt(c.MintPrice, 10, 64)
	if err != nil {
		u.Logger.ErrorAny("GetProjectDetail", zap.Any("strconv.ParseInt", err))
		return nil, err
	}
	ethPrice, err := u.convertBTCToETH(fmt.Sprintf("%f", float64(mintPriceInt)/1e8))
	if err != nil {
		u.Logger.ErrorAny("GetProjectDetail", zap.Any("convertBTCToETH", err))
		return nil, err
	}
	c.MintPriceEth = ethPrice

	networkFeeInt, err := strconv.ParseInt(c.NetworkFee, 10, 64)
	if err == nil {
		ethNetworkFeePrice, err := u.convertBTCToETH(fmt.Sprintf("%f", float64(networkFeeInt)/1e8))
		if err != nil {
			u.Logger.ErrorAny("GetProjectDetail", zap.Any("convertBTCToETH", err))
			return nil, err
		}
		c.NetworkFeeEth = ethNetworkFeePrice
	}
	u.Logger.LogAny("GetProjectDetail", zap.Any("project", c))
	return c, nil
}

func (u Usecase) GetRecentWorksProjects(req structure.FilterProjects) (*entity.Pagination, error) {

	pe := &entity.FilterProjects{}
	err := copier.Copy(pe, req)
	if err != nil {
		u.Logger.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	pe.WalletAddress = req.WalletAddress
	projects, err := u.Repo.GetRecentWorksProjects(*pe)
	if err != nil {
		u.Logger.Error("u.Repo.GetRecentWorksProjects", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("projects", projects.Total)
	return projects, nil
}

func (u Usecase) GetUpdatedProjectStats(req structure.GetProjectReq) (*entity.ProjectStat, []entity.TraitStat, error) {

	project, err := u.Repo.FindProjectBy(req.ContractAddr, req.TokenID)
	if err != nil {
		return nil, nil, err
	}

	// do not resync
	if project.Stats.LastTimeSynced != nil && project.Stats.LastTimeSynced.Unix()+int64(u.Config.TimeResyncProjectStat) > time.Now().Unix() {
		return &project.Stats, project.TraitsStat, nil
	}

	allTokenFromDb, err := u.Repo.GetAllTokensByProjectID(project.TokenID)
	if err != nil {
		return nil, nil, err
	}
	owners := make(map[string]bool)
	for _, token := range allTokenFromDb {
		owners[token.OwnerAddr] = true
	}

	var allListings []entity.MarketplaceListings
	var allOffers []entity.MarketplaceOffers

	allListings, err = u.Repo.GetAllListingByCollectionContract(project.GenNFTAddr)
	if err != nil {
		u.Logger.Error("u.Repo.GetAllListingByCollectionContract", err.Error(), err)
		return nil, nil, err
	}

	allOffers, err = u.Repo.GetAllOfferByCollectionContract(project.GenNFTAddr)
	if err != nil {
		u.Logger.Error("u.Repo.GetAllOfferByCollectionContract", err.Error(), err)
		return nil, nil, err
	}

	var totalTradingVolumn *big.Int
	var floorPrice *big.Int
	var bestMakeOfferPrice *big.Int
	var listedPercent int32
	listingSet := make(map[string]bool)

	for _, listing := range allListings {
		if listing.Erc20Token != utils.EVM_NULL_ADDRESS {
			continue
		}
		price := new(big.Int)
		price, ok := price.SetString(listing.Price, 10)
		if !ok {
			err := errors.New("fail to convert price to big int")
			u.Logger.Error("fail to convert price to big int", err.Error(), err)
			continue
		}
		durationTime, err := strconv.ParseInt(listing.DurationTime, 10, 64)
		if err != nil {
			u.Logger.Error("fail to parse duration time", err.Error(), err)
			continue
		}

		// update total volumn trading
		if listing.Finished {
			if totalTradingVolumn == nil {
				totalTradingVolumn = new(big.Int)
			}
			totalTradingVolumn.Add(totalTradingVolumn, price)
		}
		// update listing percent
		if !listing.Closed && (time.Now().Unix() < durationTime || durationTime == 0) {
			listingSet[listing.TokenId] = true
		}

		// update floor price
		if listing.Finished {
			if floorPrice == nil {
				floorPrice = price
			} else {
				if floorPrice.Cmp(price) > 0 {
					floorPrice = price
				}
			}
		}
	}

	for _, offer := range allOffers {
		price := new(big.Int)
		price, ok := price.SetString(offer.Price, 10)
		if !ok {
			err := errors.New("fail to convert price to big int")
			u.Logger.Error("fail to convert price to big int", err.Error(), err)
			continue
		}
		durationTime, err := strconv.ParseInt(offer.DurationTime, 10, 64)
		if err != nil {
			u.Logger.Error("fail to parse duration time", err.Error(), err)
			continue
		}

		// update total volumn trading
		if offer.Finished {
			if totalTradingVolumn == nil {
				totalTradingVolumn = new(big.Int)
			}
			totalTradingVolumn.Add(totalTradingVolumn, price)
		}

		// update floor price
		if !offer.Closed && (time.Now().Unix() < durationTime || durationTime == 0) {
			if bestMakeOfferPrice == nil {
				bestMakeOfferPrice = price
			} else {
				if bestMakeOfferPrice.Cmp(price) < 0 {
					bestMakeOfferPrice = price
				}
			}
		}

		// update floor price
		if offer.Finished {
			if floorPrice == nil {
				floorPrice = price
			} else {
				if floorPrice.Cmp(price) > 0 {
					floorPrice = price
				}
			}
		}
	}

	if len(allTokenFromDb) > 0 {
		listedPercent = int32(len(listingSet) * 100 / len(allTokenFromDb))
	} else {
		listedPercent = 0
	}

	if totalTradingVolumn == nil {
		totalTradingVolumn = new(big.Int)
	}
	if floorPrice == nil {
		floorPrice = new(big.Int)
	}
	if bestMakeOfferPrice == nil {
		bestMakeOfferPrice = new(big.Int)
	}

	// update trait stats
	traitToCnt := make(map[string]int32)
	traitValueToCnt := make(map[string]map[string]int32)
	for _, token := range allTokenFromDb {
		for _, attribute := range token.ParsedAttributes {
			traitToCnt[attribute.TraitType] += 1
			if traitValueToCnt[attribute.TraitType] == nil {
				traitValueToCnt[attribute.TraitType] = make(map[string]int32)
			}
			traitValueToCnt[attribute.TraitType][fmt.Sprintf("%v", attribute.Value)] += 1
		}
	}

	traitsStat := make([]entity.TraitStat, 0)
	for k, cnt := range traitToCnt {
		traitValueStat := make([]entity.TraitValueStat, 0)
		for value, cntValue := range traitValueToCnt[k] {
			traitValueStat = append(traitValueStat, entity.TraitValueStat{
				Value:  value,
				Rarity: int32(cntValue * 100 / cnt),
			})
		}
		traitsStat = append(traitsStat, entity.TraitStat{
			TraitName:       k,
			TraitValuesStat: traitValueStat,
		})
	}

	now := time.Now()

	project, err = u.Repo.FindProjectBy(req.ContractAddr, req.TokenID)
	if err != nil {
		return nil, nil, err
	}

	return &entity.ProjectStat{
		LastTimeSynced:     &now,
		UniqueOwnerCount:   uint32(len(owners)),
		TotalTradingVolumn: totalTradingVolumn.String(),
		FloorPrice:         floorPrice.String(),
		BestMakeOfferPrice: bestMakeOfferPrice.String(),
		ListedPercent:      listedPercent,
		MintedCount:        project.Stats.MintedCount,
	}, traitsStat, nil
}

func (u Usecase) getProjectDetailFromChainWithoutCache(req structure.GetProjectDetailMessageReq) (*structure.ProjectDetail, error) {

	contractDataKey := fmt.Sprintf("detail.%s.%s", req.ContractAddress, req.ProjectID)

	u.Logger.Info("req", req)

	addr := common.HexToAddress(req.ContractAddress)
	// call to contract to get emotion
	client, err := helpers.EthDialer()
	if err != nil {
		u.Logger.Error("ethclient.Dial", err.Error(), err)
		return nil, err
	}

	projectID := new(big.Int)
	projectID, ok := projectID.SetString(req.ProjectID, 10)
	if !ok {
		return nil, errors.New("cannot convert tokenID")
	}
	contractDetail, err := u.getNftContractDetailInternal(client, addr, *projectID)
	if err != nil {
		u.Logger.Error("u.getNftContractDetailInternal", err.Error(), err)
		return nil, err
	}
	//u.Logger.Info("contractDetail", contractDetail)
	u.Cache.SetData(contractDataKey, contractDetail)
	return contractDetail, nil
}

// Get from chain with cache
func (u Usecase) getProjectDetailFromChain(req structure.GetProjectDetailMessageReq) (*structure.ProjectDetail, error) {

	contractDataKey := helpers.ProjectDetailKey(req.ContractAddress, req.ProjectID)

	//u.Cache.Delete(contractDataKey)
	data, err := u.Cache.GetData(contractDataKey)
	if err != nil {
		u.Logger.Info("req", req)

		addr := common.HexToAddress(req.ContractAddress)
		// call to contract to get emotion
		client, err := helpers.EthDialer()
		if err != nil {
			u.Logger.Error("ethclient.Dial", err.Error(), err)
			return nil, err
		}

		projectID := new(big.Int)
		projectID, ok := projectID.SetString(req.ProjectID, 10)
		if !ok {
			return nil, errors.New("cannot convert tokenID")
		}
		contractDetail, err := u.getNftContractDetailInternal(client, addr, *projectID)
		if err != nil {
			u.Logger.Error("u.getNftContractDetail", err.Error(), err)
			return nil, err
		}
		u.Logger.Info("contractDetail", contractDetail)
		u.Cache.SetData(contractDataKey, contractDetail)
		return contractDetail, nil
	}

	contractDetail := &structure.ProjectDetail{}
	err = helpers.ParseCache(data, contractDetail)
	if err != nil {
		u.Logger.Error("helpers.ParseCache", err.Error(), err)
		return nil, err
	}

	return contractDetail, nil
}

// Internal get project detail
func (u Usecase) getNftContractDetailInternal(client *ethclient.Client, contractAddr common.Address, projectID big.Int) (*structure.ProjectDetail, error) {

	gProject, err := generative_project_contract.NewGenerativeProjectContract(contractAddr, client)
	if err != nil {
		u.Logger.Error("generative_project_contract.NewGenerativeProjectContract", err.Error(), err)
		return nil, err
	}

	pDchan := make(chan structure.ProjectDetailChan, 1)
	pStatuschan := make(chan structure.ProjectStatusChan, 1)
	pTokenURIchan := make(chan structure.ProjectNftTokenUriChan, 1)

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

	detailFromChain := <-pDchan
	statusFromChain := <-pStatuschan
	tokenFromChain := <-pTokenURIchan

	if detailFromChain.Err != nil {
		return nil, detailFromChain.Err
	}

	if statusFromChain.Err != nil {
		u.Logger.Error("statusFromChain.Err", statusFromChain.Err.Error(), statusFromChain.Err)
		return nil, statusFromChain.Err
	}

	if tokenFromChain.Err != nil {
		u.Logger.Error("tokenFromChain.Err", tokenFromChain.Err.Error(), tokenFromChain.Err)
		return nil, tokenFromChain.Err
	}

	gNftProject, err := generative_nft_contract.NewGenerativeNftContract(detailFromChain.ProjectDetail.GenNFTAddr, client)
	if err != nil {
		u.Logger.Error("generative_nft_contract.NewGenerativeNftContract", err.Error(), err)
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
	}

	u.Logger.Info("resp", resp)
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

	u.Logger.Info("resp", resp)
	return resp, nil
}

func (u Usecase) UnzipProjectFile(zipPayload *structure.ProjectUnzipPayload) (*entity.Projects, error) {
	pe, err := u.Repo.FindProjectByTokenID(zipPayload.ProjectID)
	if err != nil {
		u.Logger.Error("http.Get", err.Error(), err)
		return nil, err
	}
	u.Logger.Info("project", pe)
	nftTokenURI := make(map[string]interface{})
	nftTokenURI["name"] = pe.Name
	nftTokenURI["description"] = pe.Description
	nftTokenURI["image"] = pe.Thumbnail
	nftTokenURI["animation_url"] = ""
	nftTokenURI["attributes"] = []string{}

	u.Logger.Info("zipPayload", zipPayload)
	images := []string{}
	zipLink := zipPayload.ZipLink

	spew.Dump(os.Getenv("GCS_DOMAIN"))
	groupIndex := strings.Index(zipLink, "btc-projects/")
	strLen := len(zipLink)
	zipLink = zipLink[groupIndex:strLen]
	spew.Dump(zipLink)
	err = u.GCS.UnzipFile(zipLink)
	if err != nil {
		u.Logger.Error("http.Get", err.Error(), err)
		return nil, err
	}

	unzipFoler := zipLink + "_unzip"
	files, err := u.GCS.ReadFolder(unzipFoler)
	if err != nil {
		u.Logger.Error("http.Get", err.Error(), err)
		return nil, err
	}
	maxSize := uint64(0)
	for _, f := range files {
		if strings.Index(strings.ToLower(f.Name), strings.ToLower("__MACOSX")) > -1 {
			continue
		}
		if strings.Index(strings.ToLower(f.Name), strings.ToLower(".DS_Store")) > -1 {
			continue
		}

		temp := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), f.Name)
		images = append(images, temp)
		nftTokenURI["image"] = temp
		if uint64(f.Size) > maxSize {
			maxSize = uint64(f.Size)
		}
	}
	//

	pe.Images = images
	if len(images) > 0 {
		pe.IsFullChain = true
		if len(pe.Images) == 1 {
			// edition
			if pe.MaxSupply > 1 {
				//-> clone for maxsupply - 1 files
				for i := 1; i < int(pe.MaxSupply); i++ {
					pe.Images = append(pe.Images, pe.Images[0])
				}
			}
		} else {
			// list file
			if len(pe.Images) < int(pe.MaxSupply) {
				// max supply need to equal max files
				pe.MaxSupply = int64(len(pe.Images))
			}
		}

	}
	pe.IsHidden = false
	pe.Status = true
	pe.IsSynced = true

	networkFee := big.NewInt(u.networkFeeBySize(int64(maxSize / 4))) // will update after unzip and check data
	pe.NetworkFee = networkFee.String()

	updated, err := u.Repo.UpdateProject(pe.UUID, pe)
	if err != nil {
		u.Logger.Error("u.Repo.UpdateProject", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("updated", updated)
	return pe, nil
}

func (u Usecase) UploadFileZip(fc []byte, uploadChan chan uploadFileChan, peName string, fileName string, wg *sync.WaitGroup) {

	var err error
	var uploadedUrl *string

	defer func() {
		uploadChan <- uploadFileChan{
			FileURL: uploadedUrl,
			Err:     err,
		}
		wg.Done()
	}()

	base64Data := helpers.Base64Encode(fc)

	key := helpers.GenerateSlug(peName)
	key = fmt.Sprintf("btc-projects/%s/unzip", key)

	uploadFileName := fmt.Sprintf("%s/%s", key, fileName)
	uploaded, err := u.GCS.UploadBaseToBucket(base64Data, uploadFileName)
	if err != nil {
		u.Logger.Error("u.GCS.UploadBaseToBucket", err.Error(), err)
		return
	}

	u.Logger.Info("uploaded", uploaded)
	cdnURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	uploadedUrl = &cdnURL

}

func (u Usecase) CreateProjectFromCollectionMeta(meta entity.CollectionMeta) (*entity.Projects, error) {
	u.Logger.Info(fmt.Sprintf("Start create project from collection meta %s %s", meta.Name, meta.InscriptionIcon))
	pe := &entity.Projects{}

	mPrice := helpers.StringToBTCAmount("0")

	thumbnail := fmt.Sprintf("https://ordinals-explorer.generative.xyz/content/%s", meta.InscriptionIcon)

	pe.ContractAddress = os.Getenv("GENERATIVE_PROJECT")
	pe.MintPrice = mPrice.String()
	pe.NetworkFee = big.NewInt(u.networkFeeBySize(int64(300000 / 4))).String() // will update after unzip and check data or check from animation url
	pe.IsHidden = true
	pe.Status = true
	pe.IsSynced = true
	nftTokenURI := make(map[string]interface{})
	nftTokenURI["name"] = meta.Name
	nftTokenURI["description"] = meta.Description
	nftTokenURI["image"] = thumbnail
	nftTokenURI["animation_url"] = ""
	nftTokenURI["attributes"] = []string{}

	pe.CreatorAddrr = "0x0000000000000000000000000000000000000000"
	creatorAddrr, err := u.Repo.FindUserByWalletAddress(pe.CreatorAddrr)
	if err != nil {
		u.Logger.Error("u.Repo.FindUserByWalletAddress", err.Error(), err)
		return nil, err
	}

	pe.CreatorName = creatorAddrr.DisplayName

	bytes, err := json.Marshal(nftTokenURI)
	if err != nil {
		u.Logger.Error("json.Marshal.nftTokenURI", err.Error(), err)
		return nil, err
	}
	nftToken := helpers.Base64Encode(bytes)
	now := time.Now().UTC()

	pe.NftTokenUri = fmt.Sprintf("data:application/json;base64,%s", nftToken)
	pe.ProcessingImages = []string{}
	pe.MintedImages = nil
	pe.MintedTime = &now
	pe.CreatorProfile = *creatorAddrr
	pe.CreatorAddrrBTC = creatorAddrr.WalletAddressBTC
	pe.LimitSupply = 0
	pe.GenNFTAddr = pe.TokenID

	pe.Name = meta.Name
	pe.Description = meta.Description
	maxSupply, err := strconv.ParseInt(meta.Supply, 10, 64)
	if err != nil {
		maxSupply = 0
	}
	pe.MaxSupply = maxSupply
	countIndex, err := u.Repo.CountCollectionInscriptionByInscriptionIcon(meta.InscriptionIcon)
	var index int64
	if err != nil {
		index = 0
	} else {
		index = *countIndex
	}
	pe.MintingInfo.Index = index

	if pe.Categories == nil || len(pe.Categories) == 0 {
		pe.Categories = []string{u.Config.OtherCategoryID}
	}

	pe.Royalty = meta.Royalty
	pe.SocialTwitter = meta.TwitterLink
	pe.SocialDiscord = meta.DiscordLink
	pe.SocialWeb = meta.WebsiteLink
	pe.Thumbnail = thumbnail

	maxID, err := u.Repo.GetMaxBtcProjectID()
	if err != nil {
		u.Logger.Error("u.Repo.GetMaxBtcProjectID", err.Error(), err)
		return nil, err
	}
	maxID = maxID + 1
	pe.TokenIDInt = maxID
	pe.TokenID = fmt.Sprintf("%d", maxID)
	pe.GenNFTAddr = pe.TokenID
	pe.InscriptionIcon = meta.InscriptionIcon
	pe.CreatedByCollectionMeta = true
	blockNumberMinted := "0"
	pe.BlockNumberMinted = &blockNumberMinted

	err = u.Repo.CreateProject(pe)
	if err != nil {
		u.Logger.Error("u.Repo.CreateProjectFromInscription", err.Error(), err)
		return nil, err
	}

	u.Logger.Info(fmt.Sprintf("Done create project from collection meta %s %s", meta.Name, meta.InscriptionIcon))

	return pe, nil
}
