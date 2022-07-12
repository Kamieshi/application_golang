package repository

import (
	"app/internal/models"
	"app/internal/repository"
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
	"sync"
	"time"
)

type LocalStorage struct {
	M       sync.Mutex
	storage map[string]*CacheEntityObject
}

type CacheEntityObject struct {
	EntityObj *models.Entity
	DeathTime time.Time
}

type CashSteamEntityRep struct {
	StreamCommand string
	GroupName     string
	CashEntityRepositoryRedis
	LocalStorage *LocalStorage
}

func NewCashSteamEntityRep(addr string, repEnt *repository.RepoEntity) *CashSteamEntityRep {
	localStorage := LocalStorage{storage: make(map[string]*CacheEntityObject)}
	return &CashSteamEntityRep{
		StreamCommand:             "StreamCommand",
		GroupName:                 "Reader",
		LocalStorage:              &localStorage,
		CashEntityRepositoryRedis: *NewCashEntityRepository(addr, repEnt),
	}
}

type Command struct {
	Type     string
	EntityId string
}

func (c *Command) Marshal() map[string]string {
	return map[string]string{
		"Type":     c.Type,
		"EntityId": c.EntityId,
	}
}

func creatCacheEntity(ent *models.Entity) *CacheEntityObject {
	TLC, _ := strconv.Atoi(os.Getenv("TLC"))
	return &CacheEntityObject{
		EntityObj: ent,
		DeathTime: time.Now().Add(time.Duration(TLC) * time.Minute),
	}
}

func unMarshalCommand(data map[string]interface{}) *Command {
	var com = Command{
		Type:     data["Type"].(string),
		EntityId: data["EntityId"].(string),
	}
	return &com
}

func (r *CashSteamEntityRep) sendCommand(ctx context.Context, command Command) error {
	arg := redis.XAddArgs{
		Stream: r.StreamCommand,
		MaxLen: 0,
		ID:     "",
		Values: command,
	}
	res := r.client.XAdd(ctx, &arg)
	return res.Err()
}

func (r *CashSteamEntityRep) Set(ctx context.Context, entity *models.Entity) error {
	cacheObj := creatCacheEntity(entity)
	r.LocalStorage.M.Lock()
	r.LocalStorage.storage[entity.ID.String()] = cacheObj
	r.LocalStorage.M.Unlock()
	writeCommand := Command{Type: "write", EntityId: entity.ID.String()}
	err := r.sendCommand(ctx, writeCommand)
	return err
}

func (r *CashSteamEntityRep) Get(ctx context.Context, idEntity string) (*models.Entity, bool) {
	r.LocalStorage.M.Lock()
	cacheObj := r.LocalStorage.storage[idEntity]
	r.LocalStorage.M.Unlock()
	if cacheObj != nil {
		if !cacheObj.DeathTime.After(time.Now()) {
			r.Delete(ctx, idEntity)
			return nil, false
		}
		return cacheObj.EntityObj, true
	}
	return nil, false
}

func (r *CashSteamEntityRep) Delete(ctx context.Context, idEntity string) {
	r.LocalStorage.M.Lock()
	delete(r.LocalStorage.storage, idEntity)
	r.LocalStorage.M.Unlock()
	deleteCommand := Command{Type: "delete", EntityId: idEntity}
	r.sendCommand(ctx, deleteCommand)
}

func (r *CashSteamEntityRep) Listener(ctx context.Context) {
	r.client.XGroupDestroy(ctx, r.StreamCommand, r.GroupName)
	r.client.XGroupCreate(ctx, r.StreamCommand, r.GroupName, "$")
	args := redis.XReadGroupArgs{
		Group:    r.GroupName,
		Consumer: "Reader",
		Streams:  []string{r.StreamCommand, ">"},
	}
	for {
		messages := r.client.XReadGroup(ctx, &args).Val()
		for _, message := range messages {
			for _, comm := range message.Messages {
				command := unMarshalCommand(comm.Values)
				if command.Type == "write" {
					entity, err := r.entityRep.GetForID(ctx, command.EntityId)
					if err == nil {
						cacheObj := creatCacheEntity(entity)
						r.LocalStorage.M.Lock()
						r.LocalStorage.storage[entity.ID.String()] = cacheObj
						r.LocalStorage.M.Unlock()
						continue
					}
					r.LocalStorage.M.Lock()
					delete(r.LocalStorage.storage, entity.ID.String())
					r.LocalStorage.M.Unlock()
				}
			}
		}
	}
}
