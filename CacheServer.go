package main

import (
	"app/internal/config"
	redisRepository "app/internal/repository/redis"
	"context"
	"fmt"
	"github.com/fatih/structs"
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
	consumersGroup := "CacheService"
	repoCashEntity.Client().XGroupCreate(ctx, repoCashEntity.StreamGet, consumersGroup, "0")
	repoCashEntity.Client().XGroupCreate(ctx, repoCashEntity.StreamSet, consumersGroup, "0")
	if err != nil {
		log.Println(err)
	}

	uniqueID := xid.New().String()
	for {
		entries, err := repoCashEntity.Client().XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumersGroup,
			Consumer: uniqueID,
			Streams:  []string{repoCashEntity.StreamSet, repoCashEntity.StreamGet, ">", ">"},
		}).Result()
		if err != nil {
			log.Fatal(err)
		}

		for _, query := range entries {
			switch query.Stream {

			case repoCashEntity.StreamGet:
				for i := 0; i < len(entries[0].Messages); i++ {
					values := entries[0].Messages[i].Values
					idEntity := fmt.Sprintf("%v", values["idEntities"])
					idStream := fmt.Sprintf("%v", values["idOutput"])

					entity, err := repoCashEntity.CashEntityRepositoryRedis.Get(ctx, idEntity)
					if err != nil {
						log.WithError(err).Error("[GET]sGet in redis")
						err = repoCashEntity.Client().XAdd(ctx, &redis.XAddArgs{
							Stream: idStream,
							Values: map[string]interface{}{"status": 404},
						}).Err()
					}

					mapData := structs.Map(entity)
					mapData["status"] = 0

					err = repoCashEntity.Client().XAdd(ctx, &redis.XAddArgs{
						Stream: idStream,
						Values: mapData,
					}).Err()
					if err != nil {
						log.WithError(err).Error("[GET]Error send entity to the stream out")
						continue
					}

					repoCashEntity.Client().XAck(ctx, repoCashEntity.StreamGet, consumersGroup, entries[0].Messages[i].ID)
					repoCashEntity.Client().XDel(ctx, repoCashEntity.StreamGet, entries[0].Messages[i].ID)

				}

			case repoCashEntity.StreamSet:
				for i := 0; i < len(entries[0].Messages); i++ {
					values := entries[0].Messages[i].Values
					key := fmt.Sprintf("%v", values["key"])
					deathTime, _ := strconv.Atoi(values["DeathTime"].(string))

					if int64(deathTime) < time.Now().Unix() {
						log.WithError(err).Error("[SET]Time end")
						continue
					}

					cacheObj := redisRepository.CacheObj{
						DeathTime: int64(deathTime),
						Data:      []byte(values["Data"].(string)),
					}

					err = repoCashEntity.Client().Set(ctx, key, cacheObj, 0).Err()
					if err != nil {
						log.WithError(err).Error("[SET]Error Set")
						continue
					}

					repoCashEntity.Client().XAck(ctx, repoCashEntity.StreamSet, consumersGroup, entries[0].Messages[i].ID)
					repoCashEntity.Client().XDel(ctx, repoCashEntity.StreamSet, entries[0].Messages[i].ID)

					log.Info("[SET]Successful Set in cache")
				}

			}
		}
	}
}
