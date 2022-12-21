package redis

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"rederinghub.io/utils/config"

	"github.com/go-redis/redis"
)

type IRedisCache interface {
	SetData(key string, value interface{}) error
	SetStringData(key string, value string) error
	SetStringDataWithExpTime(key string, value string,  exipredIn int) error
	GetData(key string) (*string, error)
	Delete(key string) error
	SetDataWithExpireTime(key string, value interface{}, exipredIn int) error //exipredIn second
}

type redisCache struct {
	cfg    config.RedisConfig
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(cfg config.RedisConfig) *redisCache {
	r := new(redisCache)
	ctx := context.Background()
	redisDB, err := strconv.Atoi( cfg.DB)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password, // no password set
		DB:      redisDB,       // use default DB
	})

	r.cfg = cfg
	r.client = rdb
	r.ctx = ctx
	return r
}

func (r *redisCache) SetStringData(key string, value string) error {
	err := r.client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) SetStringDataWithExpTime(key string, value string,  exipredIn int) error {
	timeD := time.Duration(rand.Int31n(int32(exipredIn))) * time.Second
	err := r.client.Set(key, value, timeD).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) SetData(key string, value interface{}) error {
	valueByte, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.client.Set(key, valueByte, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) SetDataWithExpireTime(key string, value interface{}, exipredIn int) error {
	valueByte, err := json.Marshal(value)
	if err != nil {
		return err
	}
	timeD := time.Duration(rand.Int31n(int32(exipredIn))) * time.Second
	err = r.client.Set(key, valueByte, timeD).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) GetData(key string) (*string, error) {
	value, err := r.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func (r *redisCache) Delete(key string) error {
	err := r.client.Del(key).Err()
	if err != nil {
		return err
	}
	return nil
}
