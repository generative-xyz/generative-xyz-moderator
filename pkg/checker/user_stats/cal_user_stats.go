package user_stats

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/utils/logger"
)

type checkerUser struct {
	batchSize int
	ch        chan *entity.Users
	updateCh  chan *entity.Users
	wg        sync.WaitGroup
	updated   int
	repo      repository.Repository
}

func NewChecker(repo repository.Repository, batchSize int) *checkerUser {
	return &checkerUser{
		batchSize: batchSize,
		ch:        make(chan *entity.Users, batchSize),
		updateCh:  make(chan *entity.Users, batchSize),
		repo:      repo,
	}
}

func (c *checkerUser) Start(ctx context.Context) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		if err := c.updateUsers(ctx); err != nil {
			logger.AtLog.Logger.Error("updateUsers failed", zap.Error(err))
		}
	}()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer close(c.updateCh)

		users := make([]*entity.Users, 0, c.batchSize)

		for user := range c.ch {
			users = append(users, user)
			if len(users) < c.batchSize {
				continue
			}
			if err := c.calUsers(ctx, users); err != nil {
				logger.AtLog.Logger.Error("calUsers failed", zap.Error(err))
			}
			users = users[:0]
		}

		if len(users) > 0 {
			if err := c.calUsers(ctx, users); err != nil {
				logger.AtLog.Logger.Error("calUsers failed", zap.Error(err))
			}
			users = nil
		}
	}()
}

func (c *checkerUser) Add(user *entity.Users) {
	c.ch <- user
}

func (c *checkerUser) Stop() int {
	close(c.ch)
	c.wg.Wait()
	return c.updated
}

func (c *checkerUser) updateUsers(ctx context.Context) error {
	count := 0
	for user := range c.updateCh {
		logger.AtLog.Infof("updating %s ... ", user.WalletAddress)
		if _, err := c.repo.UpdateByID(ctx, user.TableName(),
			user.ID,
			bson.M{
				"$set": bson.M{
					"stats": user.Stats,
				},
			}); err != nil {
			logger.AtLog.Logger.Info(fmt.Sprintf("update %s ... failed", user.WalletAddress), zap.Error(err))
			continue
		} else {
			c.updated++
			count++
			if count < c.batchSize {
				continue
			}
		}
		count = 0
		logger.AtLog.Info("waiting 5s for next batchSize...")
		time.Sleep(5 * time.Second)
	}
	return nil
}

type CountProject struct {
	WalletAddress   string `bson:"_id"`
	TotalCollection int64  `bson:"total_collection"`
	TotalMint       int64  `bson:"total_mint"`
	TotalMinted     int64  `bson:"total_minted"`
}

func (c *checkerUser) calUsers(ctx context.Context, users []*entity.Users) error {
	userWallets := make([]string, 0, len(users))
	for _, v := range users {
		userWallets = append(userWallets, v.WalletAddress)
	}
	match := bson.M{
		"$match": bson.M{
			"creatorAddress": bson.M{"$in": userWallets},
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": "$creatorAddress",
			"total_collection": bson.M{
				"$sum": 1,
			},
			"total_mint": bson.M{
				"$sum": "$maxSupply",
			},
			"total_minted": bson.M{
				"$sum": "$index",
			},
		},
	}

	cur, err := c.repo.DB.Collection(entity.Projects{}.TableName()).Aggregate(ctx, bson.A{match, group})
	if err != nil {
		return err
	}

	defer cur.Close(ctx)

	countProjectMap := make(map[string]CountProject, len(users))
	for cur.Next(ctx) {
		var countProject CountProject
		if err := cur.Decode(&countProject); err != nil {
			return err
		}
		countProjectMap[countProject.WalletAddress] = countProject
	}
	if err := cur.Err(); err != nil {
		return err
	}

	for _, user := range users {
		countProject, ok := countProjectMap[user.WalletAddress]
		if !ok {
			continue
		}
		needUpdate := false
		if countProject.TotalCollection != user.Stats.CollectionCreated {
			needUpdate = true
			user.Stats.CollectionCreated = countProject.TotalCollection
		}
		if countProject.TotalMint != user.Stats.TotalMint {
			needUpdate = true
			user.Stats.TotalMint = countProject.TotalMint
		}
		if countProject.TotalMinted != user.Stats.TotalMinted {
			needUpdate = true
			user.Stats.TotalMinted = countProject.TotalMinted
		}
		if !needUpdate {
			continue
		}
		c.updateCh <- user
	}
	return nil
}
