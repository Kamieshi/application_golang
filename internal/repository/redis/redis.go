package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(addr string) *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisRepository{
		client: client,
	}
}

func (rr RedisRepository) Set(c context.Context, key string, value string) error {
	err := rr.client.Set(c, key, value, 0).Err()
	if err != nil {
		logrus.WithError(err).Error("Set in redis")
		return err
	}
	return nil
}
func (rr RedisRepository) Get(c context.Context, key string) (string, error) {
	val, err := rr.client.Get(c, key).Result()
	if err != nil {
		logrus.WithError(err).Error("Get in redis")
		return "", err
	}
	return val, nil
}

func (rr RedisRepository) Delete(c context.Context, key string) error {
	rr.client.Del(c, key)
	return nil
}
