package model

import (
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/pkg/drivers/mongodb"
	"rederinghub.io/pkg/logger"
)

type Database interface {
	DB() *mongo.Database
}

type database struct {
	db *mongo.Database
}

func NewDatabase() Database {
	mongoDB, err := mongodb.Init()
	if err != nil {
		logger.AtLog.Fatalf("connect mongodb failed: %v", err)
	}

	return &database{db: mongoDB}
}

func (r *database) DB() *mongo.Database {
	if r == nil {
		return nil
	}

	return r.db
}
