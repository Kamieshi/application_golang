package main

import (
	"app/internal/config"
	redisRepository "app/internal/repository/redis"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func main() {
	configuration, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	config.InitLogger()
	ctx := context.Background()
	repoCashEntity := redisRepository.NewCashSteamEntityRep(configuration.REDIS_URL)
	consumersGroup := "tickets-consumer-group"
	err = repoCashEntity.Client().XGroupCreate(ctx, repoCashEntity.StreamName, consumersGroup, "0").Err()
	if err != nil {
		log.Println(err)
	}

	uniqueID := xid.New().String()
	for {
		entries, err := repoCashEntity.Client().XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumersGroup,
			Consumer: uniqueID,
			Streams:  []string{repoCashEntity.StreamName, ">"},
			Count:    2,
			Block:    0,
			NoAck:    false,
		}).Result()
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(entries[0].Messages); i++ {
			values := entries[0].Messages[i].Values
			key := fmt.Sprintf("%v", values["key"])
			deathTime, _ := strconv.Atoi(values["DeathTime"].(string))

			if int64(deathTime) < time.Now().Unix() {
				log.WithError(err).Error("Time end")
				break
			}

			cacheObj := redisRepository.CacheObj{
				DeathTime: int64(deathTime),
				Data:      []byte(values["Data"].(string)),
			}

			err = repoCashEntity.Client().Set(ctx, key, cacheObj, 0).Err()
			if err != nil {
				log.WithError(err).Error("Error Set")
			}
			log.Info("Successful Set in cache")
		}
	}

}
