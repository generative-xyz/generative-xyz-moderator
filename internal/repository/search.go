package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

// Search projects
type SearchProjectChan struct {
	Data *[]entity.Projects
	Err  error
}

type SearchProjectListedChan struct {
	Data *map[string]entity.ProjectListed
	Err  error
}

type SearchProjectFloorPriceChan struct {
	Data *map[string]entity.ProjectFloorPrice
	Err  error
}

type SearchProjectVolumeChan struct {
	Data *map[string]entity.ProjectVolume
	Err  error
}

type SearchProjectProtabChan struct {
	Data *[]*entity.ProjectsProtab
	Err  error
}

type SearchProjectProtabAPIChan struct {
	Data *[]*entity.ProjectsProtabAPI
	Err  error
}

type SearchTotalProjectsChan struct {
	Data *int64
	Err  error
}

func (r Repository) SearchProjects(filter entity.FilterProjects) ([]entity.Projects, int, int, error) {
	match := bson.D{
		{"$text",
			bson.D{
				{"$search", fmt.Sprintf("\"%s\"", *filter.Search)},

				{"$caseSensitive", false},
				{"$diacriticSensitive", false},
			},
		},
		{"isHidden", false},
		{"isSynced", true},
	}

	skip := (filter.Page - 1) * filter.Limit
	pipeline := bson.A{
		bson.D{{"$match", match}},
		bson.D{
			{"$project",
				bson.D{
					{"txhash", 0},
					{"txHex", 0},
					{"commitTxHash", 0},
					{"revealTxHash", 0},
					{"htmlFile", 0},
					{"animation_html", 0},
					{"images", 0},
					{"processingImages", 0},
					{"mintedImages", 0},
					{"reportUsers", 0},
					{"whiteListEthContracts", 0},
					{"reservers", 0},
					{"scripts", 0},
					{"thirdPartyScripts", 0},
					{"nftTokenUri", 0},
					{"score", bson.D{{"$meta", "textScore"}}},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"score", -1}}}},
		bson.D{{"$skip", skip}},
		bson.D{{"$limit", filter.Limit}},
	}

	projectsChan := make(chan SearchProjectChan, 1)
	totalProjectsChan := make(chan SearchTotalProjectsChan, 1)

	go func(projectsChan chan SearchProjectChan) {
		p := &[]entity.Projects{}
		var err error

		defer func() {
			projectsChan <- SearchProjectChan{
				Err:  err,
				Data: p,
			}
		}()

		c, err := r.DB.Collection(utils.COLLECTION_PROJECTS).Aggregate(context.TODO(), pipeline)
		if err != nil {
			return
		}

		p1 := []entity.Projects{}
		p = &p1
		err = c.All(context.TODO(), &p1)
	}(projectsChan)
	go func(totalProjectsChan chan SearchTotalProjectsChan) {
		p := new(int64)
		var err error

		defer func() {
			totalProjectsChan <- SearchTotalProjectsChan{
				Err:  err,
				Data: p,
			}
		}()

		totalItems, err := r.DB.Collection(utils.COLLECTION_PROJECTS).CountDocuments(context.TODO(), match)
		if err != nil {
			return
		}

		p = &totalItems
	}(totalProjectsChan)

	pFChan := <-projectsChan
	tpFChan := <-totalProjectsChan

	if pFChan.Err != nil {
		return nil, 0, 0, pFChan.Err
	}

	if tpFChan.Err != nil {
		return nil, 0, 0, tpFChan.Err
	}

	totalItems := *tpFChan.Data
	projects := *pFChan.Data
	totalPages := math.Ceil(float64(totalItems) / float64(filter.Limit))
	return projects, int(totalItems), int(totalPages), nil
}

func (r Repository) AggregateForProjectsProtab(filter entity.FilterProjects) ([]*entity.ProjectsProtab, int, int, error) {
	match := bson.D{
		{"isHidden", false},
		{"isSynced", true},
		{"tokenid", bson.M{"$nin": []string{"999998", "999999", "999997"}}},
	}

	skip := (filter.Page - 1) * filter.Limit
	pipeline := bson.A{
		bson.D{{"$match", match}},
		bson.D{
			{"$project",
				bson.D{
					{"tokenid", 1},
					{"tokenIDInt", 1},
					{"contractAddress", 1},
					{"creatorAddress", 1},
					{"creatorAddrrBTC", 1},
					{"thumbnail", 1},
					{"name", 1},
					{"stats", 1},
					{"maxSupply", 1},
				},
			},
		},
		bson.D{{"$skip", skip}},
		bson.D{{"$limit", filter.Limit}},
	}

	projectsChan := make(chan SearchProjectProtabChan, 1)
	totalProjectsChan := make(chan SearchTotalProjectsChan, 1)

	go func(projectsChan chan SearchProjectProtabChan) {
		p := &[]*entity.ProjectsProtab{}
		var err error

		defer func() {
			projectsChan <- SearchProjectProtabChan{
				Err:  err,
				Data: p,
			}
		}()

		//mint - volume
		pipeline = append(pipeline,
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "dex_project_mint_volume"},
						{"localField", "tokenid"},
						{"foreignField", "_id"},
						{"as", "mint_volume"},
					},
				},
			},
			bson.D{
				{"$unwind",
					bson.D{
						{"path", "$mint_volume"},
						{"preserveNullAndEmptyArrays", true},
					},
				},
			},
			bson.D{{"$addFields", bson.D{{"mint_volume", "$mint_volume.amount"}}}},
		)

		//cex - volume
		pipeline = append(pipeline,
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "dex_cex_volume"},
						{"localField", "tokenid"},
						{"foreignField", "_id"},
						{"as", "cex_volume"},
					},
				},
			},
			bson.D{
				{"$unwind",
					bson.D{
						{"path", "$cex_volume"},
						{"preserveNullAndEmptyArrays", true},
					},
				},
			},
			bson.D{{"$addFields", bson.D{{"cex_volume", "$cex_volume.amount"}}}},
		)

		c, err := r.DB.Collection(utils.COLLECTION_PROJECTS).Aggregate(context.TODO(), pipeline)
		if err != nil {
			return
		}

		p1 := []*entity.ProjectsProtab{}
		p = &p1
		err = c.All(context.TODO(), &p1)
	}(projectsChan)
	go func(totalProjectsChan chan SearchTotalProjectsChan) {
		p := new(int64)
		var err error

		defer func() {
			totalProjectsChan <- SearchTotalProjectsChan{
				Err:  err,
				Data: p,
			}
		}()

		totalItems, err := r.DB.Collection(utils.COLLECTION_PROJECTS).CountDocuments(context.TODO(), match)
		if err != nil {
			return
		}

		p = &totalItems
	}(totalProjectsChan)

	pFChan := <-projectsChan
	tpFChan := <-totalProjectsChan

	if pFChan.Err != nil {
		return nil, 0, 0, pFChan.Err
	}

	if tpFChan.Err != nil {
		return nil, 0, 0, tpFChan.Err
	}

	totalItems := *tpFChan.Data
	projects := *pFChan.Data
	totalPages := math.Ceil(float64(totalItems) / float64(filter.Limit))

	projectIDs := []string{}
	for _, p := range projects {
		projectIDs = append(projectIDs, p.TokenID)
	}

	projectListedChan := make(chan SearchProjectListedChan, 1)
	projectFloorPrice := make(chan SearchProjectFloorPriceChan, 1)
	projectVolume := make(chan SearchProjectVolumeChan, 1)

	go func(projectIDs []string, projectListedChan chan SearchProjectListedChan) {
		p := &[]entity.ProjectListed{}
		var err error

		defer func() {
			if err != nil {
				projectListedChan <- SearchProjectListedChan{
					Err:  err,
					Data: nil,
				}
			} else {

				resp := make(map[string]entity.ProjectListed)
				for _, item := range *p {
					resp[item.ID] = item
				}

				projectListedChan <- SearchProjectListedChan{
					Err:  nil,
					Data: &resp,
				}
			}

		}()

		//cex-volume
		pipeline1 := bson.A{
			bson.D{
				{"$match",
					bson.D{
						{"_id",
							bson.D{
								{"$in",
									projectIDs,
								},
							},
						},
					},
				},
			},
		}

		c, err := r.DB.Collection("dex_project_mkp_listed").Aggregate(context.TODO(), pipeline1)
		if err != nil {
			return
		}

		p1 := []entity.ProjectListed{}
		p = &p1
		err = c.All(context.TODO(), &p1)
	}(projectIDs, projectListedChan)
	go func(projectIDs []string, projectFloorPrice chan SearchProjectFloorPriceChan) {
		p := &[]entity.ProjectFloorPrice{}
		var err error

		defer func() {
			if err != nil {
				projectFloorPrice <- SearchProjectFloorPriceChan{
					Err:  err,
					Data: nil,
				}
			} else {

				resp := make(map[string]entity.ProjectFloorPrice)
				for _, item := range *p {
					resp[item.ID] = item
				}

				projectFloorPrice <- SearchProjectFloorPriceChan{
					Err:  nil,
					Data: &resp,
				}
			}

		}()

		//cex-volume
		pipeline1 := bson.A{
			bson.D{
				{"$match",
					bson.D{
						{"_id",
							bson.D{
								{"$in",
									projectIDs,
								},
							},
						},
					},
				},
			},
		}

		c, err := r.DB.Collection("dex_btc_listing_floor_price").Aggregate(context.TODO(), pipeline1)
		if err != nil {
			return
		}

		p1 := []entity.ProjectFloorPrice{}
		p = &p1
		err = c.All(context.TODO(), &p1)
	}(projectIDs, projectFloorPrice)
	go func(projectIDs []string, projectVolume chan SearchProjectVolumeChan) {
		p := &[]entity.ProjectVolume{}
		var err error

		defer func() {
			if err != nil {
				projectVolume <- SearchProjectVolumeChan{
					Err:  err,
					Data: nil,
				}
			} else {

				resp := make(map[string]entity.ProjectVolume)
				for _, item := range *p {
					resp[item.ID] = item
				}

				projectVolume <- SearchProjectVolumeChan{
					Err:  nil,
					Data: &resp,
				}
			}

		}()

		//cex-volume
		pipeline1 := bson.A{
			bson.D{
				{"$match",
					bson.D{
						{"_id",
							bson.D{
								{"$in",
									projectIDs,
								},
							},
						},
					},
				},
			},
		}

		c, err := r.DB.Collection("dex_btc_volume").Aggregate(context.TODO(), pipeline1)
		if err != nil {
			return
		}

		p1 := []entity.ProjectVolume{}
		p = &p1
		err = c.All(context.TODO(), &p1)
	}(projectIDs, projectVolume)

	pListedFChan := <-projectListedChan
	pFloorPriceFChan := <-projectFloorPrice
	pVolumeFChan := <-projectVolume

	if pListedFChan.Err != nil {
		return nil, 0, 0, pListedFChan.Err
	}

	if pFloorPriceFChan.Err != nil {
		return nil, 0, 0, pFloorPriceFChan.Err
	}

	if pVolumeFChan.Err != nil {
		return nil, 0, 0, pVolumeFChan.Err
	}

	for _, p := range projects {

		l := *pListedFChan.Data
		f := *pFloorPriceFChan.Data
		pv := *pVolumeFChan.Data

		p.Listed = l[p.TokenID].Listed
		p.FloorPrice = f[p.TokenID].Amount
		p.Volume = pv[p.TokenID].Amount

	}

	return projects, int(totalItems), int(totalPages), nil
}

// Search Artists
type SearchArtistChan struct {
	Data *[]*entity.FilteredUser
	Err  error
}

type SearchTotalArtists struct {
	Total int `bson:"total" json:"total"`
}

type SearchTotalArtistsChan struct {
	Data *SearchTotalArtists
	Err  error
}

func (r Repository) SearchArtists(filter entity.FilterProjects) ([]*entity.FilteredUser, int, int, error) {
	match := bson.D{
		{"$text",
			bson.D{
				{"$search", fmt.Sprintf("\"%s\"", *filter.Search)},

				{"$caseSensitive", false},
				{"$diacriticSensitive", false},
			},
		},
	}

	skip := (filter.Page - 1) * filter.Limit
	pipeline := bson.A{
		bson.D{{"$match", match}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "projects"},
					{"localField", "wallet_address"},
					{"foreignField", "creatorAddress"},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									bson.D{
										{"isHidden", false},
										{"isSynced", true},
									},
								},
							},
							bson.D{
								{"$addFields",
									bson.D{
										{"id", "$_id"},
										{"walletAddress", "$creatorAddress"},
									},
								},
							},
							bson.D{
								{"$project",
									bson.D{
										{"_id", 0},
										{"id", 1},
										{"name", 1},
										{"walletAddress", 1},
									},
								},
							},
						},
					},
					{"as", "projects"},
				},
			},
		},
		bson.D{{"$addFields", bson.D{{"count_projects", bson.D{{"$size", "$projects"}}}}}},
		bson.D{{"$match", bson.D{{"count_projects", bson.D{{"$gte", 1}}}}}},
	}

	pipeline1 := append(pipeline, bson.D{{"$sort", bson.D{{"score", entity.SORT_DESC}}}})
	pipeline1 = append(pipeline, bson.D{{"$skip", skip}})
	pipeline1 = append(pipeline1, bson.D{{"$limit", filter.Limit}})

	artistChan := make(chan SearchArtistChan, 1)
	totalArtistChan := make(chan SearchTotalArtistsChan, 1)

	go func(artistChan chan SearchArtistChan) {
		p := &[]*entity.FilteredUser{}
		var err error

		defer func() {
			artistChan <- SearchArtistChan{
				Err:  err,
				Data: p,
			}
		}()

		c, err := r.DB.Collection(utils.COLLECTION_USERS).Aggregate(context.TODO(), pipeline1)
		if err != nil {
			return
		}

		p1 := []*entity.FilteredUser{}
		p = &p1
		err = c.All(context.TODO(), &p1)
	}(artistChan)
	go func(totalArtistChan chan SearchTotalArtistsChan) {
		p := new(SearchTotalArtists)
		var err error

		defer func() {
			totalArtistChan <- SearchTotalArtistsChan{
				Err:  err,
				Data: p,
			}
		}()

		pipeline = append(pipeline, bson.D{
			{"$group",
				bson.D{
					{"_id", bson.D{{"wallet_address", "$wallet_address"}}},
					{"total", bson.D{{"$sum", 1}}},
				},
			},
		})

		totalItems, err := r.DB.Collection(utils.COLLECTION_USERS).Aggregate(context.TODO(), pipeline)
		if err != nil {
			return
		}

		t := []SearchTotalArtists{}
		err = totalItems.All(context.TODO(), &t)
		if err != nil {
			return
		}

		if len(t) == 0 {
			return
		}

		t1 := t[0]

		p = &t1
	}(totalArtistChan)

	pFChan := <-artistChan
	tpFChan := <-totalArtistChan

	if pFChan.Err != nil {
		return nil, 0, 0, pFChan.Err
	}

	if tpFChan.Err != nil {
		return nil, 0, 0, tpFChan.Err
	}

	totalItems := *tpFChan.Data
	projects := *pFChan.Data
	totalPages := math.Ceil(float64(totalItems.Total) / float64(filter.Limit))
	return projects, int(totalItems.Total), int(totalPages), nil
}
