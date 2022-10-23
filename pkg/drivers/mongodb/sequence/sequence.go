package sequence

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Counter struct {
	ID  string `json:"id" bson:"id"`
	Seq int    `json:"seq" bson:"seq"`
}

func GetNextID(collection *mongo.Collection, sequenceName string) (int, error) {
	findOptions := options.FindOneAndUpdate()
	findOptions.SetUpsert(true)
	var counter Counter
	err := collection.FindOneAndUpdate(context.TODO(),
		bson.M{"id": sequenceName},
		bson.M{"$inc": bson.M{"seq": 1}}, // bson.D{{"$inc", bson.D{{"seq", 1}}}},
		findOptions).Decode(&counter)

	if err != nil {
		return 0, err
	}
	return counter.Seq, nil
}

func GetNextRepeatID(collection *mongo.Collection, sequenceName string, maxId int) (int, error) {
	findOptions := options.FindOneAndUpdate()
	findOptions.SetUpsert(true)
	var counter Counter
	err := collection.FindOneAndUpdate(context.TODO(),
		bson.M{"id": sequenceName},
		bson.M{"$inc": bson.M{"seq": 1}}, // bson.D{{"$inc", bson.D{{"seq", 1}}}},
		findOptions).Decode(&counter)
	if err != nil {
		return 0, err
	}

	if counter.Seq == maxId+1 {
		err := collection.FindOneAndUpdate(context.TODO(),
			bson.M{"id": sequenceName},
			bson.M{"$set": bson.M{"seq": 1}},
			findOptions).Decode(&counter)
		if err != nil {
			return 0, err
		}
		return 0, nil
	}

	return counter.Seq, nil
}
