// Package repository cache
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	rds "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"app/internal/models"
	"app/internal/repository"
)

// CacheObj struct Cache object
type CacheObj struct {
	Data      []byte `json:"data"`
	DeathTime int64  `json:"deathTime"`
}

// NewCache Constructor
func NewCache(value []byte) (CacheObj, error) {
	var cacheObj CacheObj
	timeLive, err := strconv.Atoi(os.Getenv("TIME_EXPIRED_CACHE_MINUTE"))
	if err != nil {
		return cacheObj, err
	}
	cacheObj.Data = value
	cacheObj.DeathTime = time.Now().Add(time.Duration(timeLive) * time.Minute).Unix()
	return cacheObj, err
}

// MarshalBinary object -> []byte
func (co *CacheObj) MarshalBinary() ([]byte, error) {
	return json.Marshal(co)
}

// CashEntityRepositoryRedis Implement repoEntityCache
type CashEntityRepositoryRedis struct {
	entityRep repository.RepoEntity
	client    *rds.Client
}

// Client Redis client
func (c CashEntityRepositoryRedis) Client() *rds.Client {
	return c.client
}

// NewCashEntityRepository Constructor
func NewCashEntityRepository(addr string, entityRep repository.RepoEntity) *CashEntityRepositoryRedis {
	client := rds.NewClient(&rds.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &CashEntityRepositoryRedis{
		client:    client,
		entityRep: entityRep,
	}
}

// Set set cache value
func (c CashEntityRepositoryRedis) Set(ctx context.Context, entity *models.Entity) error {
	key := entity.ID.String()
	value, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	cacheObj, err := NewCache(value)
	if err != nil {
		return err
	}
	err = c.client.Set(ctx, key, cacheObj, 0).Err()
	if err != nil {
		logrus.WithError(err).Error("Set in redis")
		return err
	}
	return nil
}

// Get get cache value
func (c CashEntityRepositoryRedis) Get(ctx context.Context, id string) (*models.Entity, error) {
	var cacheObj CacheObj
	ent := models.Entity{}
	logrus.Info("Try get with cache")
	val, err := c.client.Get(ctx, id).Result()

	if err != rds.Nil {
		logrus.WithError(err).Error("Get in redis")
		return &ent, err
	}

	err = json.Unmarshal([]byte(val), &cacheObj)
	if err != nil {
		return &ent, err
	}

	if cacheObj.DeathTime < time.Now().Unix() {
		c.Delete(ctx, id)
		logrus.WithError(err).Info("Get in redis")
		return &ent, errors.New("time expired")
	}

	err = json.Unmarshal(cacheObj.Data, &ent)
	if err != nil {
		return &ent, err
	}
	logrus.Info("From Cache")
	return &ent, nil
}

// Delete delete value from redis
func (c CashEntityRepositoryRedis) Delete(ctx context.Context, id string) {
	c.client.Del(ctx, id)
}
