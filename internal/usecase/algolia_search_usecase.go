package usecase

import (
	"fmt"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/algolia"
)

func (uc *Usecase) AlgoliaSearchProjectV1(filter *algolia.AlgoliaFilter) (*entity.Pagination, error) {
	if filter.ObjType != "" && filter.ObjType != "project" {
		return nil, nil
	}
	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)
	resp, err := algoliaClient.FetchObjIdsBySearch("projects", filter)
	if err != nil {
		return nil, err
	}

	pe := &entity.FilterProjects{Ids: resp}
	pe.Page = int64(filter.Page)
	pe.Limit = int64(filter.Limit)
	return uc.Repo.GetProjects(*pe)
}

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
		inscriptions = append(inscriptions, i)
	}
	resp.UnmarshalHits(&inscriptions)

	dataResp := []*response.SearchResponse{}
	for _, i := range inscriptions {
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
	tokens, err := uc.Repo.FilterTokenUri(*pe)
	if err != nil {
		return nil, 0, 0, err
	}
	iTokens := tokens.Result
	rTokens := iTokens.([]entity.TokenUri)

	return rTokens, resp.NbHits, resp.NbPages, nil
}
