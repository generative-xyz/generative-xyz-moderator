package algolia

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

type AlgoliaFilter struct {
	Page      int
	Limit     int
	SearchStr string
}

type IAlgolia interface {
	FetchObjIdsBySearch(string, *AlgoliaFilter) ([]string, error)
	Search(string, *AlgoliaFilter) (search.QueryRes, error)
}

type algolia struct {
	client *search.Client
}

func NewAlgoliaClient(algoliaApplicationId, algoliaApiKey string) IAlgolia {
	client := search.NewClient(algoliaApplicationId, algoliaApiKey)
	return &algolia{
		client: client,
	}
}

func (al *algolia) Search(indexName string, builder *AlgoliaFilter) (search.QueryRes, error) {
	index := al.client.InitIndex(indexName)
	if builder.Limit == 0 {
		builder.Limit = 10
	}

	if builder.Page == 0 {
		builder.Page = 1
	}

	opts := []interface{}{
		opt.Page(builder.Page - 1),
		opt.HitsPerPage(builder.Limit),
	}
	res, err := index.Search(builder.SearchStr, opts...)
	return res, err
}

func (al *algolia) FetchObjIdsBySearch(indexName string, builder *AlgoliaFilter) ([]string, error) {
	index := al.client.InitIndex(indexName)
	if builder.Limit == 0 {
		builder.Limit = 10
	}

	opts := []interface{}{
		opt.Page(builder.Page),
		opt.HitsPerPage(builder.Limit),
	}
	res, err := index.Search(builder.SearchStr, opts)
	if err != nil {
		return nil, err
	}
	data := []string{}
	for _, p := range res.Hits {
		data = append(data, p["objectID"].(string))
	}
	return data, nil
}
