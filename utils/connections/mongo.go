package connections

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	mongotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go.mongodb.org/mongo-driver/mongo"
)

type mongoCN struct {
	Cnn *mongo.Client
}

func NewMongo(dsn string) (*mongoCN, error) {
	p := new(mongoCN)

	isLoadbalanced := false
	maxPoolSize := uint64(10000)
	minPoolSize := uint64(1)
	clientOption := &options.ClientOptions{
		LoadBalanced: &isLoadbalanced,
		MaxPoolSize:  &maxPoolSize,
		MinPoolSize:  &minPoolSize,
	}

	clientOption.ApplyURI(dsn)
	clientOption.Monitor = mongotrace.NewMonitor()

	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, err
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
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
