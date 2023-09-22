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
					{"name", 1},
					{"creatorProfile", 1},
					{"image", 1},
					{"thumbnail", 1},
					{"description", 1},
					{"tokenid", 1},
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
