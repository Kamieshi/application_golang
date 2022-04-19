package repository

import (
	"app/internal/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/fatih/structs"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CashSteamEntityRep struct {
	StreamSet string
	StreamGet string
	GroupName string
	CashEntityRepositoryRedis
}

func NewCashSteamEntityRep(addr string) *CashSteamEntityRep {
	return &CashSteamEntityRep{
		StreamGet:                 "entityGet",
		StreamSet:                 "entitySet",
		GroupName:                 "Reader",
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
		Stream: rr.StreamSet,
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

func (rr CashSteamEntityRep) Get(c context.Context, id string) (models.Entity, error) {
	ent := models.Entity{}
	logrus.Info("Try get with cache")

	idOutStream := xid.New().String()

	err := rr.client.XAdd(c, &redis.XAddArgs{
		Stream: rr.StreamGet,
		Values: map[string]interface{}{"idEntities": id, "idOutput": idOutStream},
	}).Err()
	if err != nil {
		logrus.Error(err)
		return ent, err
	}

	value, err := rr.client.XRead(c, &redis.XReadArgs{
		Block:   1 * time.Millisecond,
		Streams: []string{idOutStream, "0"},
	}).Result()
	if err != nil {
		logrus.Error(err)
		return ent, err
	}

	rr.Client().Del(c, idOutStream)

	if value[0].Messages[0].Values["status"] != "0" {
		logrus.WithFields(logrus.Fields{
			"status":       value[0].Messages[0].Values["status"],
			"id_entity":    id,
			"id_OutStream": idOutStream,
			"result":       value,
		}).Error()
		err = errors.New("Not fount in cashe")

		return ent, err
	}
	err = ent.InitForMap(value[0].Messages[0].Values)
	return ent, err
}
