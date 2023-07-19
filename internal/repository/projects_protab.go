package repository

import (
	"context"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"strings"
)

func (r Repository) InsertProjectProData(input *entity.ProjectsProtab) error {

	f := bson.D{
		{"tokenid", input.TokenID},
	}

	u := &entity.UpdateProjectsProtab{}

	err := copier.Copy(u, input)
	if err != nil {
		return err
	}

	updatedData, err := u.ToBson()
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": *updatedData,
	}

	opts := &options.UpdateOptions{}
	opts.SetUpsert(true)
	_, err = r.DB.Collection(entity.ProjectsProtab{}.TableName()).UpdateOne(context.TODO(), f, update, opts)
	if err != nil {
		return err
	}

	return err
}

func (r Repository) UpdateProjectVolumeBuyable(contractAddress, tokenID string, volume float64, buyable bool) error {

	f := bson.D{
		{"contractAddress", contractAddress},
		{"tokenid", tokenID},
	}

	update := bson.M{
		"$set": bson.D{
			{"stats.volume", volume},
			{"stats.buyable", buyable},
		},
	}

	_, err := r.DB.Collection(entity.Projects{}.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return err
	}

	return err
}

func (r Repository) UpdateProjectUniqueOwner(contractAddress string, tokenID string, owners int) error {

	f := bson.D{
		{"tokenid", tokenID},
		{"contractAddress", strings.ToLower(contractAddress)},
	}

	update := bson.M{
		"$set": bson.D{{"unique_owners", owners}},
	}

	_, err := r.DB.Collection(entity.ProjectsProtab{}.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return err
	}

	return err
}

func (r Repository) GetProjectsProtab(filter entity.FilterProjects) ([]*entity.ProjectsProtab, int, int, error) {

	skip := (filter.Page - 1) * filter.Limit
	pipeline := bson.A{
		bson.D{{"$match", bson.D{
			{"tokenid", bson.M{"$ne": ""}},
		}}},
		bson.D{{"$sort", bson.D{
			{"is_buyable", entity.SORT_DESC},
			{"volume", entity.SORT_DESC},
		}}},
		bson.D{{"$skip", skip}},
		bson.D{{"$limit", filter.Limit}},
	}

	match := bson.D{}
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

		c, err := r.DB.Collection(utils.COLLECTION_PROJECT_PROTAB).Aggregate(context.TODO(), pipeline)
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

		totalItems, err := r.DB.Collection(utils.COLLECTION_PROJECT_PROTAB).CountDocuments(context.TODO(), match)
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

func (r Repository) AggregateProjectsProtab(filter entity.FilterProjects) ([]*entity.ProjectsProtabAPI, int, int, error) {
	match := bson.D{
		{"tokenid", bson.D{{"$ne", ""}}},
		{"isHidden", false},
		{"isSynced", true},
	}

	skip := (filter.Page - 1) * filter.Limit
	pipeline := bson.A{
		bson.D{
			{"$project",
				bson.D{
					{"tokenid", 1},
					{"contractAddress", 1},
					{"name", 1},
					{"description", 1},
					{"thumbnail", 1},
					{"creatorAddress", 1},
					{"isSynced", 1},
					{"isHidden", 1},
					{"index", 1},
					{"indexReverse", 1},
					{"maxSupply", 1},
					{"is_buyable", 1},
					{"stats.buyable", 1},
					{"stats.volume", 1},
					{"creatorProfile.wallet_address", 1},
					{"creatorProfile.wallet_address_payment", 1},
					{"creatorProfile.wallet_address_btc", 1},
					{"creatorProfile.wallet_address_btc_taproot", 1},
					{"creatorProfile.display_name", 1},
					{"creatorProfile.avatar", 1},
				},
			},
		},
		bson.D{
			{"$match", match},
		},
		bson.D{
			{"$sort",
				bson.D{
					{"stats.buyable", entity.SORT_DESC},
					{"stats.volume", entity.SORT_DESC},
				},
			},
		},
		bson.D{{"$skip", skip}},
		bson.D{{"$limit", filter.Limit}},
	}

	projectsChan := make(chan SearchProjectProtabAPIChan, 1)
	totalProjectsChan := make(chan SearchTotalProjectsChan, 1)

	go func(projectsChan chan SearchProjectProtabAPIChan) {
		p := &[]*entity.ProjectsProtabAPI{}
		var err error

		defer func() {
			projectsChan <- SearchProjectProtabAPIChan{
				Err:  err,
				Data: p,
			}
		}()

		pipeline = append(pipeline, bson.D{
			{"$lookup",
				bson.D{
					{"from", utils.COLLECTION_PROJECT_PROTAB},
					{"localField", "tokenid"},
					{"foreignField", "tokenid"},
					{"pipeline",
						bson.A{
							bson.D{
								{"$project",
									bson.D{
										{"cex_volume", 1},
										{"floor_price", 1},
										{"is_buyable", 1},
										{"listed", 1},
										{"mint_volume", 1},
										{"unique_owners", 1},
										{"volume", 1},
									},
								},
							},
						},
					},
					{"as", "projects_detail"},
				},
			},
		},
			bson.D{
				{"$unwind",
					bson.D{
						{"path", "$projects_detail"},
						{"preserveNullAndEmptyArrays", true},
					},
				},
			},
			bson.D{
				{"$addFields",
					bson.D{
						{"cex_volume", "$projects_detail.cex_volume"},
						{"floor_price", "$projects_detail.floor_price"},
						{"is_buyable", "$projects_detail.is_buyable"},
						{"listed", "$projects_detail.listed"},
						{"mint_volume", "$projects_detail.mint_volume"},
						{"unique_owners", "$projects_detail.unique_owners"},
						{"volume", "$projects_detail.volume"},
					},
				},
			})

		c, err := r.DB.Collection(utils.COLLECTION_PROJECTS).Aggregate(context.TODO(), pipeline)
		if err != nil {
			return
		}

		p1 := []*entity.ProjectsProtabAPI{}
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
