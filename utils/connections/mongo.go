package connections

import (
	"context"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


type mongoCN struct {
	Cnn *mongo.Client
}

func NewMongo(dsn string) (*mongoCN, error) {
	p := new(mongoCN)
	cmdMonitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			//spew.Dump(evt.Command)
		},
		Succeeded:  func(ctx context.Context, evt *event.CommandSucceededEvent) {
			//spew.Dump(evt.DurationNanos)
		},
		Failed:    func(ctx context.Context, evt *event.CommandFailedEvent) {
			//spew.Dump(evt.Failure)
		},
	}

	isLoadbalanced := false
	maxPoolSize := uint64(10000)
	minPoolSize := uint64(1)
	clientOption := &options.ClientOptions{
		LoadBalanced: &isLoadbalanced,
		MaxPoolSize:  &maxPoolSize,
		MinPoolSize:  &minPoolSize,
	}

	clientOption.ApplyURI(dsn)
	clientOption.Monitor = cmdMonitor

	client, err := mongo.Connect(context.TODO(),clientOption)
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
