package repository

import (
	"app/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	var cacheObj CacheObj
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

	fmt.Println(rr.client.XGroupCreateMkStream(c, idOutStream, rr.GroupName, "0"))

	newContext, cancelFunc := context.WithTimeout(c, time.Millisecond*3)
	defer cancelFunc()

	value, err := rr.client.XReadGroup(newContext, &redis.XReadGroupArgs{
		Group:   rr.GroupName,
		Streams: []string{idOutStream, ">"},
		Count:   1,
	}).Result()
	if err != nil {
		logrus.Error(err)
		return ent, err
	}

	fmt.Println(value)
	//TODO Unpars from map[string]interface to model.Entity
	err = json.Unmarshal([]byte("das"), &cacheObj)
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
	logrus.Info("From Cache")
	return ent, nil
}
