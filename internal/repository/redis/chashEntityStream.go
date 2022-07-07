package repository

import (
	"app/internal/models"
	"app/internal/repository"
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
)

type LocalStorage struct {
	storage map[string]models.Entity
	sync.Mutex
}

type CashSteamEntityRep struct {
	StreamCommand string
	GroupName     string
	CashEntityRepositoryRedis
	LocalStorage *LocalStorage
}

func NewCashSteamEntityRep(addr string, repEnt *repository.RepoEntity) *CashSteamEntityRep {
	localStorage := LocalStorage{storage: make(map[string]models.Entity)}
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
	r.LocalStorage.storage[entity.Id] = *entity
	writeCommand := Command{Type: "write", EntityId: entity.Id}
	err := r.sendCommand(ctx, writeCommand)
	return err
}

func (r *CashSteamEntityRep) Get(ctx context.Context, idEntity string) (*models.Entity, bool) {
	entity, exist := r.LocalStorage.storage[idEntity]
	return &entity, exist
}

func (r *CashSteamEntityRep) Delete(ctx context.Context, idEntity string) {
	delete(r.LocalStorage.storage, idEntity)
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
						r.LocalStorage.storage[entity.Id] = *entity
						continue
					}
					delete(r.LocalStorage.storage, entity.Id)
				}
			}
		}
	}
}
