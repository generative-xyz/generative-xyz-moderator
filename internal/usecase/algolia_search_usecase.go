package usecase

import (
	"fmt"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils/algolia"
)

func (uc *Usecase) AlgoliaSearchProject(filter *algolia.AlgoliaFilter) ([]*response.SearchResponse, error) {
	if filter.ObjType != "" && filter.ObjType != "project" {
		return nil, nil
	}
	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)

	resp, err := algoliaClient.Search("projects", filter)
	if err != nil {
		return nil, err
	}

	projects := []*response.SearchProject{}
	resp.UnmarshalHits(&projects)
	dataResp := []*response.SearchResponse{}
	for _, i := range projects {
		obj := &response.SearchResponse{
			ObjectType: "project",
			Project:    i,
		}
		dataResp = append(dataResp, obj)
	}

	return dataResp, nil
}

func (uc *Usecase) AlgoliaSearchInscription(filter *algolia.AlgoliaFilter) ([]*response.SearchResponse, error) {
	if filter.ObjType != "" && filter.ObjType != "inscription" {
		return nil, nil
	}

	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)

	resp, err := algoliaClient.Search("inscriptions", filter)
	if err != nil {
		return nil, err
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

	return dataResp, nil
}

func (uc *Usecase) AlgoliaSearchArtist(filter *algolia.AlgoliaFilter) ([]*response.SearchResponse, error) {
	if filter.ObjType != "" && filter.ObjType != "artist" {
		return nil, nil
	}
	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)

	resp, err := algoliaClient.Search("users", filter)
	if err != nil {
		return nil, err
	}
	artists := []*response.SearchArtist{}
	resp.UnmarshalHits(&artists)

	dataResp := []*response.SearchResponse{}
	for _, i := range artists {
		obj := &response.SearchResponse{
			ObjectType: "artist",
			Artist:     i,
		}
		dataResp = append(dataResp, obj)
	}
	return dataResp, nil
}

func (uc *Usecase) AlgoliaSearchTokenUri(filter *algolia.AlgoliaFilter) ([]*response.SearchResponse, error) {
	if filter.ObjType != "" && filter.ObjType != "token" {
		return nil, nil
	}

	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)

	resp, err := algoliaClient.Search("token-uris", filter)
	if err != nil {
		return nil, err
	}
	inscriptions := []*response.SearchTokenUri{}
	resp.UnmarshalHits(&inscriptions)

	dataResp := []*response.SearchResponse{}
	for _, i := range inscriptions {
		obj := &response.SearchResponse{
			ObjectType: "token",
			TokenUri:   i,
		}
		dataResp = append(dataResp, obj)
	}

	return dataResp, nil
}
