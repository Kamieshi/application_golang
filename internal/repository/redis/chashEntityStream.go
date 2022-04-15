package repository

import (
	"app/internal/models"
	"context"
	"encoding/json"
	"github.com/fatih/structs"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CashSteamEntityRep struct {
	StreamName string
	CashEntityRepositoryRedis
}

func NewCashSteamEntityRep(addr string) *CashSteamEntityRep {
	return &CashSteamEntityRep{
		StreamName:                "entitySet",
		CashEntityRepositoryRedis: *NewCashEntityRepository(addr),
	}
}

func (rr CashSteamEntityRep) Set(c context.Context, entity *models.Entity) error {
	key := entity.Id.(primitive.ObjectID).Hex()
	value, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	cacheObj := NewCache(value)
	data := structs.Map(cacheObj)
	data["key"] = key

	arg := redis.XAddArgs{
		Stream: rr.StreamName,
		MaxLen: 0,
		ID:     "",
		Values: data,
	}

	err = rr.client.XAdd(c, &arg).Err()
	if err != nil {
		logrus.WithError(err).Error("Set in redis")
		return err
	}
	return nil
}
