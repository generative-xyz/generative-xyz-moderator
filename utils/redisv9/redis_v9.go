package redisv9

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"rederinghub.io/utils/config"
)

const (
	// DefaultMaxRetries
	DefaultMaxRetries int = 3
	// DefaultPoolSize --
	DefaultPoolSize = 50
	// DefaultTimeout--
	DefaultTimeout = 3 * time.Second
)

type Client interface {
	Get(ctx context.Context, key string, result interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Dels(ctx context.Context, keys ...string) error
	DelPrefix(ctx context.Context, prefix string) error
	HGet(ctx context.Context, key, field string, result interface{}) error
	HSet(ctx context.Context, key string, values ...interface{}) error
	HDel(ctx context.Context, key string) error
	GetRedisClient() *goredislib.Client
	GetRedSyncClient() *redsync.Redsync
	// ... more
}

type clientImpl struct {
	redisClient *goredislib.Client
	redSync     *redsync.Redsync
}

func NewClient(cfg config.RedisConfig) Client {

	redisClient := goredislib.NewClient(&goredislib.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		MaxRetries:   DefaultMaxRetries,
		ReadTimeout:  DefaultTimeout,
		WriteTimeout: DefaultTimeout,
		PoolSize:     DefaultPoolSize,
		PoolTimeout:  DefaultTimeout,
	})
	redSync := redsync.New(goredis.NewPool(redisClient))
	return &clientImpl{redisClient, redSync}
}

func (s *clientImpl) GetRedisClient() *goredislib.Client {
	return s.redisClient
}

func (s *clientImpl) GetRedSyncClient() *redsync.Redsync {
	return s.redSync
}

func (s *clientImpl) Get(ctx context.Context, key string, result /*must have method MarshalBinary and UnmarshalBinary*/ interface{}) error {
	return s.redisClient.Get(ctx, key).Scan(result)
}

func (s *clientImpl) Set(ctx context.Context, key string, value /*must have method MarshalBinary and UnmarshalBinary*/ interface{}, expiration time.Duration) error {
	return s.redisClient.Set(ctx, key, value, expiration).Err()
}

func (s *clientImpl) Del(ctx context.Context, key string) error {
	return s.redisClient.Del(ctx, key).Err()
}

func (s *clientImpl) Dels(ctx context.Context, keys ...string) error {
	return s.redisClient.Del(ctx, keys...).Err()
}

func (s *clientImpl) DelPrefix(ctx context.Context, prefix string) error {
	// TODO: move to worker
	return s.delKeys(ctx, prefix, s.redisClient)
}

func (s *clientImpl) delKeys(ctx context.Context, prefix string, client goredislib.UniversalClient) error {
	keys, err := client.Keys(ctx, prefix+"*").Result() // TODO: remove this code
	if err == nil && len(keys) > 0 {
		if err := client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (s *clientImpl) HGet(ctx context.Context, key, field string, result interface{}) error {
	return s.redisClient.HGet(ctx, key, field).Scan(result)
}

func (s *clientImpl) HSet(ctx context.Context, key string, values ...interface{}) error {
	return s.redisClient.HSet(ctx, key, values).Err()
}

func (s *clientImpl) HDel(ctx context.Context, key string) error {
	return s.redisClient.HDel(ctx, key).Err()
}
