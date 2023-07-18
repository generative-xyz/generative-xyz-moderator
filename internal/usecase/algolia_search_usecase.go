package usecase

import (
	"fmt"
	"go.uber.org/zap"
	"rederinghub.io/utils/copier"
	"rederinghub.io/utils/logger"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/algolia"
)

func (uc *Usecase) AlgoliaSearchProjectListing(filter *algolia.AlgoliaFilter) ([]*response.ProjectListing, int, int, error) {
	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)
	filter.FilterStr = "isHidden = 0"
	resp, err := algoliaClient.Search("project-listing", filter)
	if err != nil {
		return nil, 0, 0, err
	}
	listings := []*response.ProjectListing{}
	resp.UnmarshalHits(&listings)
	return listings, resp.NbHits, resp.NbPages, nil
}

func (uc *Usecase) AlgoliaSearchProject(filter *algolia.AlgoliaFilter) ([]entity.Projects, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "project" {
		return nil, 0, 0, nil
	}
	filter.FilterStr = "isHidden = 0 AND isSynced = 1"
	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)

	resp, err := algoliaClient.Search("projects", filter)
	if err != nil {
		return nil, 0, 0, err
	}

	projects := []*response.SearchProject{}
	resp.UnmarshalHits(&projects)
	if len(projects) == 0 {
		return nil, resp.NbHits, resp.NbPages, nil
	}

	ids := []string{}
	for _, i := range projects {
		ids = append(ids, i.ObjectId)
	}

	pe := &entity.FilterProjects{Ids: ids}
	pe.Limit = int64(filter.Limit)
	pe.Page = 1
	hidden := false
	pe.IsHidden = &hidden
	uProjects, err := uc.Repo.GetProjects(*pe)
	if err != nil {
		return nil, 0, 0, err
	}
	iProjects := uProjects.Result
	eProjects := iProjects.([]entity.Projects)
	return eProjects, resp.NbHits, resp.NbPages, nil
}

func (uc *Usecase) AlgoliaSearchInscriptionFromTo(filter *algolia.AlgoliaFilter) ([]*response.SearhcInscription, error) {
	if filter.FromNumber >= 0 && filter.ToNumber > 0 {
		filter.FilterStr += fmt.Sprintf("number:%d TO %d AND sat > 0", filter.FromNumber, filter.ToNumber)
	}

	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)
	resp, err := algoliaClient.Search("inscriptions", filter)
	if err != nil {
		uc.Logger.Error(err)
		return nil, err
	}

	inscriptions := []*response.SearhcInscription{}
	for _, h := range resp.Hits {
		i := &response.SearhcInscription{
			ObjectId:      h["objectID"].(string),
			InscriptionId: h["inscription_id"].(string),
			Number:        int64(h["number"].(float64)),
			Sat:           h["sat"].(float64),
			Chain:         h["chain"].(string),
			GenesisFee:    int64(h["genesis_fee"].(float64)),
			GenesisHeight: int64(h["genesis_height"].(float64)),
			Timestamp:     h["timestamp"].(string),
			ContentType:   h["content_type"].(string),
		}
		inscriptions = append(inscriptions, i)
	}
	return inscriptions, nil
}

func (uc *Usecase) AlgoliaSearchInscription(filter *algolia.AlgoliaFilter) ([]*response.SearchResponse, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "inscription" {
		return nil, 0, 0, nil
	}

	if filter.FromNumber > 0 && filter.ToNumber > 0 {
		filter.FilterStr += fmt.Sprintf("number:%d TO %d", filter.FromNumber, filter.ToNumber)
	}

	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)
	resp, err := algoliaClient.Search("inscriptions", filter)
	if err != nil {
		uc.Logger.Error(err)
		return nil, 0, 0, err
	}

	uc.Logger.Infof("%s", uc.Config.GenerativeExplorerApi)

	inscriptions := []*response.SearhcInscription{}
	userAddresses := []string{}
	inscriptionIds := []string{}
	client := resty.New()
	for _, h := range resp.Hits {
		i := &response.SearhcInscription{
			ObjectId:      h["objectID"].(string),
			InscriptionId: h["inscription_id"].(string),
			Number:        int64(h["number"].(float64)),
			Sat:           h["sat"].(float64),
			Chain:         h["chain"].(string),
			GenesisFee:    int64(h["genesis_fee"].(float64)),
			GenesisHeight: int64(h["genesis_height"].(float64)),
			Timestamp:     h["timestamp"].(string),
			ContentType:   h["content_type"].(string),
		}

		inscriptionIds = append(inscriptionIds, i.InscriptionId)
		if v, ok := h["address"]; ok && v.(string) != "" {
			i.Address = v.(string)
		}

		resp := &response.SearhcInscription{}
		_, err := client.R().
			SetResult(&resp).
			Get(fmt.Sprintf("%s/inscription/%s", uc.Config.GenerativeExplorerApi, i.InscriptionId))

		if err != nil {
			uc.Logger.Error(err)
		}

		if resp.Address != "" {
			i.Address = resp.Address
			userAddresses = append(userAddresses, resp.Address)
		}

		inscriptions = append(inscriptions, i)
	}
	uc.Logger.Infof("inscription count %d", len(inscriptions))

	users, err := uc.Repo.ListUserBywalletAddressBtcTaproot(userAddresses)
	mapOwner := make(map[string]*response.ArtistResponse)
	for _, o := range users {
		mapOwner[o.WalletAddressBTCTaproot] = o
	}

	pe := &entity.FilterTokenUris{TokenIDs: inscriptionIds}
	pe.Limit = int64(filter.Limit)
	pe.Page = 1
	tokens, err := uc.Repo.FilterTokenUri(*pe)
	if err != nil {
		uc.Logger.Error(err)
		return nil, 0, 0, err
	}
	iTokens := tokens.Result
	rTokens := iTokens.([]entity.TokenUri)

	mapData := make(map[string]string)
	projectIds := []string{}
	for _, t := range rTokens {
		projectIds = append(projectIds, t.ProjectID)
		mapData[t.TokenID] = t.ProjectID
	}

	projects, err := uc.Repo.FindProjectByTokenIDs(projectIds)
	if err != nil {
		uc.Logger.Error(err)
		return nil, 0, 0, err
	}

	mapProject := make(map[string]*entity.Projects)
	for _, p := range projects {
		mapProject[p.TokenID] = p
	}

	dataResp := []*response.SearchResponse{}
	for _, i := range inscriptions {
		i.Owner = mapOwner[i.Address]
		pId := mapData[i.InscriptionId]
		if pId != "" {
			if p, ok := mapProject[pId]; ok {
				i.ProjectName = p.Name
				i.ProjectTokenId = p.TokenID
			}
		}

		listingInfo, err := uc.Repo.GetDexBTCListingOrderPendingByInscriptionID(i.InscriptionId)
		if err == nil && listingInfo.CancelTx == "" {
			i.Buyable = true
			i.PriceBTC = fmt.Sprintf("%v", listingInfo.Amount)
		}

		obj := &response.SearchResponse{
			ObjectType:  "inscription",
			Inscription: i,
		}
		dataResp = append(dataResp, obj)
	}

	return dataResp, resp.NbHits, resp.NbPages, nil
}

func (uc *Usecase) AlgoliaSearchArtist(filter *algolia.AlgoliaFilter) ([]*response.ArtistResponse, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "artist" {
		return nil, 0, 0, nil
	}
	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)

	resp, err := algoliaClient.Search("users", filter)
	if err != nil {
		return nil, 0, 0, err
	}
	artists := []*response.SearchArtist{}
	resp.UnmarshalHits(&artists)
	if len(artists) == 0 {
		return nil, resp.NbHits, resp.NbPages, nil
	}

	ids := []string{}
	for _, i := range artists {
		ids = append(ids, i.ObjectId)
	}

	req := structure.FilterUsers{Ids: ids}
	req.Limit = int64(filter.Limit)
	req.Page = 1
	uUsers, err := uc.Repo.ListUsers(req)
	iUsers := uUsers.Result
	rUsers := iUsers.([]*response.ArtistResponse)
	return rUsers, resp.NbHits, resp.NbPages, nil
}

func (uc *Usecase) AlgoliaSearchTokenUri(filter *algolia.AlgoliaFilter) ([]entity.TokenUri, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "token" {
		return nil, 0, 0, nil
	}

	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)

	resp, err := algoliaClient.Search("token-uris", filter)
	if err != nil {
		return nil, 0, 0, err
	}
	data := []*response.SearchTokenUri{}
	resp.UnmarshalHits(&data)
	if len(data) == 0 {
		return nil, resp.NbHits, resp.NbPages, nil
	}

	ids := []string{}
	for _, i := range data {
		ids = append(ids, i.ObjectId)
	}
	pe := &entity.FilterTokenUris{Ids: ids}
	pe.Limit = int64(filter.Limit)
	pe.Page = 1
	tokens, err := uc.Repo.FilterTokenUri(*pe)
	if err != nil {
		return nil, 0, 0, err
	}
	iTokens := tokens.Result
	rTokens := iTokens.([]entity.TokenUri)
	return rTokens, resp.NbHits, resp.NbPages, nil
}

func (uc *Usecase) DBSearchProject(filter *algolia.AlgoliaFilter) ([]entity.Projects, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "project" {
		return nil, 0, 0, nil
	}

	pe := &entity.FilterProjects{}
	pe.Limit = int64(filter.Limit)
	pe.Page = 1

	pe.Search = &filter.SearchStr
	hidden := false
	isSynced := true
	pe.IsHidden = &hidden
	pe.IsSynced = &isSynced

	uProjects, total, totalPages, err := uc.Repo.SearchProjects(*pe)
	if err != nil {
		return nil, 0, 0, err
	}

	return uProjects, total, totalPages, nil
}

func (uc *Usecase) DBSearchArtists(filter *algolia.AlgoliaFilter) ([]*response.ArtistResponse, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "artist" {
		return nil, 0, 0, nil
	}

	pe := &entity.FilterProjects{}
	pe.Limit = int64(filter.Limit)
	pe.Page = 1

	pe.Search = &filter.SearchStr
	hidden := false
	isSynced := true
	pe.IsHidden = &hidden
	pe.IsSynced = &isSynced

	as1, total, totalPages, err := uc.Repo.SearchArtists(*pe)
	if err != nil {
		return nil, 0, 0, err
	}

	as := []*response.ArtistResponse{}
	for _, item := range as1 {
		respItem := &response.ArtistResponse{}
		err := copier.Copy(respItem, item)
		if err != nil {
			return nil, 0, 0, err
		}

		as = append(as, respItem)
	}

	return as, total, totalPages, nil
}

// Project - protab
type ProjectUniqueOwnersChan struct {
	ProjectID       string
	ContractAddress string
	Owners          int
	Err             error
}

func (uc *Usecase) JobProjectProtab() error {
	key := fmt.Sprintf("CrontabDBProjectProtab")
	page := 1
	limit := 100

	for {
		filter := &algolia.AlgoliaFilter{
			Page:  page,
			Limit: limit,
		}

		projects, totalItems, totalPages, err := uc.DBProjectProtab(filter)
		if err != nil {
			logger.AtLog.Logger.Error(key, zap.Error(err))
			return err
		}

		if len(projects) == 0 {
			break
		}

		key := fmt.Sprintf("CrontabDBProjectProtab - page: %d, limit: %d", page, limit)
		logger.AtLog.Logger.Info(key,
			zap.Int("projects", len(projects)),
			zap.Int("totalItems", totalItems),
			zap.Int("totalPages", totalPages),
		)

		//save projects to DB
		var wg sync.WaitGroup
		for i, p := range projects {

			wg.Add(1)
			go func(wg *sync.WaitGroup, p *entity.ProjectsProtab) {
				defer wg.Done()

				err := uc.Repo.InsertProjectProData(p)
				if err != nil {
					logger.AtLog.Logger.Error(fmt.Sprintf("%s - projectID %s", key, p.TokenID), zap.Error(err),
						zap.String("project", p.TokenID),
						//zap.String("contract_address", p.ContractAddress),
					)
					//return err
				} else {
					//logger.AtLog.Logger.Info(fmt.Sprintf("%s - projectID %s", key, p.TokenID), zap.Error(err),
					//	zap.Any("project", p.TokenIDInt),
					//	//zap.String("contract_address", p.ContractAddress),
					//)
				}

			}(&wg, p)

			if (i > 0 && i%20 == 0) || (i == len(projects)-1) {
				wg.Wait()
			}
		}

		page++
	}

	return nil
}

func (uc *Usecase) JobProjectProtabUniqueOwner() error {
	key := fmt.Sprintf("JobProjectProtabUniqueOwner")
	page := 1
	limit := 10

	for {
		filter := &algolia.AlgoliaFilter{
			Page:  page,
			Limit: limit,
		}

		uProjects, totalItems, totalPages, err := uc.DBProjectProtabAggerateOwner(filter)
		if err != nil {
			logger.AtLog.Logger.Error(key, zap.Error(err))
			return err
		}

		if len(uProjects) == 0 {
			break
		}

		key := fmt.Sprintf("JobProjectProtabUniqueOwner - page: %d, limit: %d", page, limit)
		logger.AtLog.Logger.Info(key,
			zap.Int("projects", len(uProjects)),
			zap.Int("totalItems", totalItems),
			zap.Int("totalPages", totalPages),
		)

		//calculate unique owners
		inputChan := make(chan entity.ProjectsProtab, len(uProjects))
		outputChan := make(chan ProjectUniqueOwnersChan, len(uProjects))

		var wg sync.WaitGroup

		for _, _ = range uProjects {
			go uc.CalculateUniqueOwner(&wg, inputChan, outputChan)
		}

		for i, pe := range uProjects {
			wg.Add(1)
			inputChan <- *pe
			if (i != 0 && i%10 == 0) || (i == len(uProjects)-1) {
				wg.Wait()
			}
		}

		for _, _ = range uProjects {

			dataFChan := <-outputChan
			if dataFChan.Err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("%s - contract: %s, tokenID: %s", key, dataFChan.ContractAddress, dataFChan.ProjectID), zap.Error(err))
				continue
			}

			//skip if owner is zero
			if dataFChan.Owners == 0 {
				continue
			}

			err := uc.Repo.UpdateProjectUniqueOwner(dataFChan.ContractAddress, dataFChan.ProjectID, dataFChan.Owners)
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("%s - contract: %s, tokenID: %s", key, dataFChan.ContractAddress, dataFChan.ProjectID), zap.Error(err))
				continue
			}

			logger.AtLog.Logger.Info(fmt.Sprintf("%s - contract: %s, tokenID: %s", key, dataFChan.ContractAddress, dataFChan.ProjectID), zap.Int("owners", dataFChan.Owners))

		}

		page++
	}

	return nil
}

func (uc *Usecase) DBProjectProtab(filter *algolia.AlgoliaFilter) ([]*entity.ProjectsProtab, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "project" {
		return nil, 0, 0, nil
	}

	pe := &entity.FilterProjects{}
	pe.Limit = int64(filter.Limit)
	pe.Page = int64(filter.Page)

	pe.Search = &filter.SearchStr
	hidden := false
	isSynced := true
	pe.IsHidden = &hidden
	pe.IsSynced = &isSynced

	uProjects, total, totalPages, err := uc.Repo.AggregateForProjectsProtab(*pe)
	if err != nil {
		return nil, 0, 0, err
	}

	for _, p := range uProjects {
		//p.UniqueOwners = uOwners[p.TokenID]
		p.Volume = p.Volume + p.MintVolume + p.CEXVolume
		if p.FloorPrice > 0 {
			p.IsBuyable = true
		}

		//hard code for 2 projects
		if p.TokenID == "1000001" {
			p.Volume = 1407649205
		} else if p.TokenID == "1002573" {
			p.Volume = 866067000
		}
	}

	//calculate unique owners
	return uProjects, total, totalPages, nil
}

func (uc *Usecase) DBProjectProtabAggerateOwner(filter *algolia.AlgoliaFilter) ([]*entity.ProjectsProtab, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "project" {
		return nil, 0, 0, nil
	}

	pe := &entity.FilterProjects{}
	pe.Limit = int64(filter.Limit)
	pe.Page = int64(filter.Page)

	pe.Search = &filter.SearchStr
	hidden := false
	isSynced := true
	pe.IsHidden = &hidden
	pe.IsSynced = &isSynced

	uProjects, total, totalPages, err := uc.Repo.GetProjectsProtab(*pe)
	if err != nil {
		return nil, 0, 0, err
	}

	//calculate unique owners
	return uProjects, total, totalPages, nil
}

func (uc *Usecase) DBProjectProtabAPI(filter *algolia.AlgoliaFilter) ([]*entity.ProjectsProtabAPI, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "project" {
		return nil, 0, 0, nil
	}

	pe := &entity.FilterProjects{}
	pe.Limit = int64(filter.Limit)
	pe.Page = int64(filter.Page)

	pe.Search = &filter.SearchStr
	hidden := false
	isSynced := true
	pe.IsHidden = &hidden
	pe.IsSynced = &isSynced

	uProjects, total, totalPages, err := uc.Repo.AggregateProjectsProtab(*pe)
	if err != nil {
		return nil, 0, 0, err
	}

	//calculate unique owners
	return uProjects, total, totalPages, nil
}

func (uc *Usecase) DBProjectProtabAPIFormatData(filter *algolia.AlgoliaFilter) ([]*response.ProjectListing, int, int, error) {
	uProjects, total, totalPages, err := uc.DBProjectProtabAPI(filter)
	if err != nil {
		return nil, 0, 0, err
	}

	listings := []*response.ProjectListing{}
	for _, p := range uProjects {

		respItem := &response.ProjectListing{
			ObjectID:        p.TokenID,
			ContractAddress: p.ContractAddress,
			TotalSupply:     p.MaxSupply,
			NumberOwners:    int64(p.UniqueOwners),
			ProjectMarketplaceData: &response.ProjectMarketplaceData{
				Listed:          uint64(p.Listed),
				FloorPrice:      uint64(p.FloorPrice),
				TotalVolume:     uint64(p.Volume),
				CEXVolume:       uint64(p.CEXVolume),
				MintVolume:      uint64(p.MintVolume),
				FirstSaleVolume: float64(0), //TODO - update this field
			},
			IsHidden:    false, // hidden projects are filtered out of this collection
			TotalVolume: uint64(p.Volume),
			IsBuyable:   p.IsBuyable,
			Project: &response.ProjectInfo{
				Name:            p.Project.Name,
				ContractAddress: p.Project.ContractAddress,
				TokenId:         p.Project.TokenId,
				Thumbnail:       p.Project.Thumbnail,
				CreatorAddress:  p.Project.CreatorAddress,
				MaxSupply:       int64(p.Project.MaxSupply),
				IsMintedOut:     p.Project.IsMintedOut,
				MintingInfo: response.ProjectMintingInfo{
					Index:        int64(p.Project.MintingInfo.Index),
					IndexReverse: int64(p.Project.MintingInfo.IndexReverse),
				},
			},
			Owner: &response.OwnerInfo{
				WalletAddress:           p.Owner.WalletAddress,
				WalletAddressPayment:    p.Owner.WalletAddressPayment,
				WalletAddressBTC:        p.Owner.WalletAddressBtc,
				WalletAddressBTCTaproot: p.Owner.WalletAddressBtcTaproot,
				DisplayName:             p.Owner.DisplayName,
				Avatar:                  p.Owner.Avatar,
			},
			Priority: 1,
		}
		listings = append(listings, respItem)
	}

	//calculate unique owners
	return listings, total, totalPages, nil
}

func (uc *Usecase) CalculateUniqueOwner(wg *sync.WaitGroup, inputChan chan entity.ProjectsProtab, outputChan chan ProjectUniqueOwnersChan) {
	defer wg.Done()
	in := <-inputChan
	owner := 0
	var err error

	key := fmt.Sprintf("CalculateUniqueOwner - projectID: %s ", in.TokenID)

	pID := strings.ToLower(in.TokenID)
	owners, err := uc.AnalyticsTokenUriOwner(structure.FilterTokens{
		GenNFTAddr: &pID,
	})

	if err == nil {
		owners1 := owners.([]*tokenOwner)
		owner = len(owners1)
		//logger.AtLog.Logger.Info(key, zap.Int("owner", owner), zap.String("projectID", in.TokenID))
	} else {
		logger.AtLog.Logger.Error(key, zap.Error(err), zap.String("projectID", in.TokenID))
	}

	outputChan <- ProjectUniqueOwnersChan{
		ProjectID:       in.TokenID,
		ContractAddress: strings.ToLower(in.ContractAddress),
		Err:             err,
		Owners:          owner,
	}
}
