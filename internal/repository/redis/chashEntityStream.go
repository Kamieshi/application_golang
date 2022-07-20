// Package repository work with cache stream
package repository

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	rds "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"app/internal/models"
	"app/internal/repository"
)

// LocalStorage storage
type LocalStorage struct {
	M       sync.Mutex
	storage map[string]*CacheEntityObject
}

// CacheEntityObject Cache obj
type CacheEntityObject struct {
	EntityObj *models.Entity
	DeathTime time.Time
}

// CashSteamEntityRep rep
type CashSteamEntityRep struct {
	StreamCommand string
	GroupName     string
	CashEntityRepositoryRedis
	LocalStorage *LocalStorage
}

// NewCashSteamEntityRep Constructor
func NewCashSteamEntityRep(addr string, repEnt repository.RepoEntity) *CashSteamEntityRep {
	localStorage := LocalStorage{storage: make(map[string]*CacheEntityObject)}
	return &CashSteamEntityRep{
		StreamCommand:             "StreamCommand",
		GroupName:                 "Reader",
		LocalStorage:              &localStorage,
		CashEntityRepositoryRedis: *NewCashEntityRepository(addr, repEnt),
	}
}

// Command type command
type Command struct {
	Type     string
	EntityID string
}

// Marshal command
func (c *Command) Marshal() map[string]string {
	return map[string]string{
		"Type":     c.Type,
		"EntityID": c.EntityID,
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
		EntityID: data["EntityID"].(string),
	}
	return &com
}

func (r *CashSteamEntityRep) sendCommand(ctx context.Context, command Command) error {
	arg := rds.XAddArgs{
		Stream: r.StreamCommand,
		MaxLen: 0,
		ID:     "",
		Values: command,
	}
	res := r.client.XAdd(ctx, &arg)
	if res.Err() != nil {
		return fmt.Errorf("cashEntityStream.go/sendCommand : %v", res.Err())
	}
	return nil
}

// Set new cache item
func (r *CashSteamEntityRep) Set(ctx context.Context, entity *models.Entity) error {
	cacheObj := creatCacheEntity(entity)
	r.LocalStorage.M.Lock()
	r.LocalStorage.storage[entity.ID.String()] = cacheObj
	r.LocalStorage.M.Unlock()
	writeCommand := Command{Type: "write", EntityID: entity.ID.String()}
	err := r.sendCommand(ctx, writeCommand)
	if err != nil {
		return fmt.Errorf("cashEntityStream.go/Set : %v", err)
	}
	return nil
}

// Get entity from cache
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

// Delete value from cache
func (r *CashSteamEntityRep) Delete(ctx context.Context, idEntity string) {
	r.LocalStorage.M.Lock()
	delete(r.LocalStorage.storage, idEntity)
	r.LocalStorage.M.Unlock()
	deleteCommand := Command{Type: "delete", EntityID: idEntity}
	err := r.sendCommand(ctx, deleteCommand)
	if err != nil {
		log.WithError(err).Error("Delete haven't worked")
	}
}

// Listener goroutine from listening redis stream
func (r *CashSteamEntityRep) Listener(ctx context.Context) {
	r.client.XGroupDestroy(ctx, r.StreamCommand, r.GroupName)
	r.client.XGroupCreate(ctx, r.StreamCommand, r.GroupName, "$")
	args := rds.XReadGroupArgs{
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
					entity, err := r.entityRep.GetForID(ctx, command.EntityID)
					if err == nil {
						cacheObj := creatCacheEntity(entity)
						r.LocalStorage.M.Lock()
						r.LocalStorage.storage[entity.ID.String()] = cacheObj
						r.LocalStorage.M.Unlock()
						continue
					}
					log.WithError(err)
					r.LocalStorage.M.Lock()
					delete(r.LocalStorage.storage, entity.ID.String())
					r.LocalStorage.M.Unlock()
				}
			}
		}
	}
}
