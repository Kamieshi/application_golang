package repository

import (
	"app/internal/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strconv"
	"time"
)

type CacheObj struct {
	Data      []byte `json:"data"`
	DeathTime int64  `json:"deathTime"`
}

func NewCache(value []byte) CacheObj {
	timeLive, _ := strconv.Atoi(os.Getenv("TIME_EXPIRED_CACHE_MINUTE"))
	return CacheObj{
		Data:      value,
		DeathTime: time.Now().Add(time.Duration(timeLive) * time.Minute).Unix(),
	}
}

func (co CacheObj) MarshalBinary() ([]byte, error) {
	return json.Marshal(co)
}

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

	cacheObj := NewCache(value)
	err = rr.client.Set(c, key, cacheObj, 0).Err()
	if err != nil {
		logrus.WithError(err).Error("Set in redis")
		return err
	}
	return nil
}

func (rr CashEntityRepositoryRedis) Get(c context.Context, id string) (models.Entity, error) {
	var cacheObj CacheObj
	ent := models.Entity{}

	val, err := rr.client.Get(c, id).Result()

	if err != nil {
		logrus.WithError(err).Error("Get in redis")
		return ent, err
	}

	err = json.Unmarshal([]byte(val), &cacheObj)
	if err != nil {
		return ent, err
	}

	if cacheObj.DeathTime < time.Now().Unix() {
		_ = rr.Delete(c, id)
		logrus.WithError(err).Info("Get in redis")
		return ent, errors.New("time expired")
	}

	err = json.Unmarshal(cacheObj.Data, &ent)
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
