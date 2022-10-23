package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"rederinghub.io/pkg/logger"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	mongo_tracer "gopkg.in/DataDog/dd-trace-go.v1/contrib/go.mongodb.org/mongo-driver/mongo"
)

const (
	GrpcContextTimeOutInSec = 10
)

func Init() (*mongo.Database, error) {
	db, err := connectDb(DefaultConnectionFromConfig())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectDb(conn *Connection) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(conn.Uri)
	clientOptions.SetMaxPoolSize(conn.MaxPoolSize)
	clientOptions.SetMinPoolSize(conn.MaxPoolSize)
	clientOptions.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	clientOptions.SetReadConcern(readconcern.Majority())
	clientOptions.SetMonitor(mongo_tracer.NewMonitor())
	ctx, cancel := context.WithTimeout(context.Background(), GrpcContextTimeOutInSec*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client.Database(conn.DbName), nil
}

type Connection struct {
	Uri         string
	DbName      string
	MaxPoolSize uint64
	MinPoolSize uint64
}

func CreateMongoDbConnection(conn *Connection) *mongo.Database {
	if conn == nil {
		conn = DefaultConnectionFromConfig()
	}
	db, err := connectDb(conn)
	if err != nil {
		logger.AtLog.Fatalf("mongodb connected failed: %v", err)
	}
	return db
}

func DefaultConnectionFromConfig() *Connection {
	conn := &Connection{
		Uri:         viper.GetString(`MONGODB_URI`),
		DbName:      viper.GetString(`MONGODB_DBNAME`),
		MaxPoolSize: viper.GetUint64(`MONGODB_MAX_POOL_SIZE`),
		MinPoolSize: viper.GetUint64(`MONGODB_MIN_POOL_SIZE`),
	}
	if conn.MaxPoolSize <= 0 {
		conn.MaxPoolSize = 100
	}
	if conn.MinPoolSize <= 0 {
		conn.MinPoolSize = 4
	}
	return conn
}
