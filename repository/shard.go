package repository

import (
	"application/service/models"
	"log"
)

type Shard interface {
	GetAllItems() ([]interface{}, error)
	GetItemForID(id int) (interface{}, error)
	AddItem(interface{}) (bool, error)
	UpdateItem(id int, obj interface{}) (bool, error)
	DeleteItem(id int) error
}

type ShardEntity interface {
	GetAllItems() ([]models.Entity, error)
	GetItemForID(id int) (models.Entity, error)
	AddItem(obj models.Entity) (bool, error)
	UpdateItem(id int, obj models.Entity) error
	DeleteItem(id int) error
}


func GetAllItems(shard ShardEntity) ([]models.Entity, error){
	 entitys,err := shard.GetAllItems()
	 if err!= nil{
		 log.Println(err)
		 return nil, err
	 }
	 return entitys, nil
}

func 	GetItemForID(shard ShardEntity,id int) (models.Entity, error){
	entity,err := shard.GetItemForID(id)
	if err!= nil{
		log.Println(err)
		return models.Entity{}, err
	}
	return entity, nil
}

func AddItem(shard ShardEntity, obj models.Entity) (bool, error){
	status,err := shard.AddItem(obj)
	if err!= nil{
		log.Println(err)
		return status, err
	}
	return status, nil
}

func UpdateItem(shard ShardEntity,id int, obj models.Entity) error{
	err:= shard.UpdateItem(id,obj)
	return err
}

func DeleteItem(shard ShardEntity, id int) error{
	err := shard.DeleteItem(id)
	return err
}