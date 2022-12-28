package connections

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


type mongoCN struct {
	Cnn *mongo.Client
}

func NewMongo(dsn string) (*mongoCN, error) {
	p := new(mongoCN)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))
    if err != nil {
		return nil, err
	}	

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	
	p.Cnn = client
	return p, nil
}

func (s *mongoCN) Connect() interface{} {
	return nil
}

func (s *mongoCN) Disconnect() error {
	err := s.Cnn.Disconnect(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

func (s *mongoCN) GetType() interface{} {
	return s.Cnn
}
