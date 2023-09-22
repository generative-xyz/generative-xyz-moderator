package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"time"
)

func (r *Repository) AggregateGMDashboardCachedDataByTime(time *time.Time) ([]entity.AggregatedGMDashBoard, error) {
	resp := []entity.AggregatedGMDashBoard{}
	f := bson.A{
		bson.D{{"$match", bson.D{{"created_at", bson.D{{"$lte", time}}}}}},
		bson.D{
			{"$project",
				bson.D{
					{"usdt", "$value.usdtvalue"},
					{"contributors",
						bson.D{
							{"$size",
								bson.D{
									{"$ifNull",
										bson.A{
											"$value.items",
											bson.A{},
										},
									},
								},
							},
						},
					},
					{"created_at", 1},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$limit", 1}},
	}

	cursor, err := r.DB.Collection(entity.CachedGMDashBoard{}.TableName()).Aggregate(context.Background(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
