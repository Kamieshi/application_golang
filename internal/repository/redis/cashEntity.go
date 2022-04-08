package repository

import (
	"app/internal/models"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CashEntityRepositoryRedis struct {
	client *redis.Client
}

func NewCashEntityRepository(addr string) *CashEntityRepositoryRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &CashEntityRepositoryRedis{
		client: client,
	}
}

func (rr CashEntityRepositoryRedis) Set(c context.Context, entity *models.Entity) error {
	key := entity.Id.(primitive.ObjectID).Hex()
	value, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	err = rr.client.Set(c, key, value, 0).Err()
	if err != nil {
		logrus.WithError(err).Error("Set in redis")
		return err
	}
	return nil
}

func (rr CashEntityRepositoryRedis) Get(c context.Context, id string) (models.Entity, error) {
	val, err := rr.client.Get(c, id).Result()
	ent := models.Entity{}
	if err != nil {
		logrus.WithError(err).Error("Get in redis")
		return ent, err
	}
	err = json.Unmarshal([]byte(val), &ent)
	if err != nil {
		return ent, err
	}
	return ent, nil
}

func (rr CashEntityRepositoryRedis) Delete(c context.Context, id string) error {
	res := rr.client.Del(c, id)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
