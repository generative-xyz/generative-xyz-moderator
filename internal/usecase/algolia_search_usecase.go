package usecase

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/algolia"
)

func (uc *Usecase) AlgoliaSearchProject(filter *algolia.AlgoliaFilter) ([]entity.Projects, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "project" {
		return nil, 0, 0, nil
	}
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
	uProjects, err := uc.Repo.GetProjects(*pe)
	if err != nil {
		return nil, 0, 0, err
	}
	iProjects := uProjects.Result
	eProjects := iProjects.([]entity.Projects)
	return eProjects, resp.NbHits, resp.NbPages, nil
}

func (uc *Usecase) AlgoliaSearchInscription(filter *algolia.AlgoliaFilter) ([]*response.SearchResponse, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "inscription" {
		return nil, 0, 0, nil
	}

	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)
	resp, err := algoliaClient.Search("inscriptions", filter)
	if err != nil {
		return nil, 0, 0, err
	}

	inscriptions := []*response.SearhcInscription{}
	userAddresses := []string{}
	inscriptionIds := []string{}
	client := resty.New()
	for _, h := range resp.Hits {
		i := &response.SearhcInscription{
			ObjectId:      h["objectID"].(string),
			InscriptionId: h["inscription_id"].(string),
			Number:        int64(h["number"].(float64)),
			Sat:           fmt.Sprintf("%d", int64(h["sat"].(float64))),
			Chain:         h["chain"].(string),
			GenesisFee:    int64(h["genesis_fee"].(float64)),
			GenesisHeight: int64(h["genesis_height"].(float64)),
			Timestamp:     h["timestamp"].(string),
			ContentType:   h["content_type"].(string),
		}

		inscriptionIds = append(inscriptionIds, i.InscriptionId)
		if v, ok := h["address"]; ok && v.(string) != "" {
			i.Address = v.(string)
			resp := &response.SearhcInscription{}
			_, err := client.R().
				SetResult(&resp).
				Get(fmt.Sprintf("%s/inscription/%s", uc.Config.GenerativeExplorerApi, i.InscriptionId))
			if err != nil {
				continue
			}
			userAddresses = append(userAddresses, resp.Address)
		}
		inscriptions = append(inscriptions, i)
	}

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
