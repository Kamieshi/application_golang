package repository

import (
	"app/internal/service"
	log "github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

func TestRepositoryUserAdd(t *testing.T) {
	repUser := NewRepoUsersPostgres(pgPool)
	serviceUser := service.NewUserService(repUser)
	user, err := serviceUser.Create(ctx, "testUsername", "TestPassword")
	if err != nil {
		t.Fatal("Error in service user", err)
	}
	t.Cleanup(func() {
		if err = repUser.Delete(ctx, user.UserName); err != nil {
			log.WithError(err).Error()
		}
	})
	userFromDB, err := repUser.Get(ctx, user.UserName)
	if err != nil {
		t.Error(err, " Error Get")
	}
	if !reflect.DeepEqual(userFromDB, user) {
		t.Error("Not equal")
	}
}

func TestRepositoryUserGet(t *testing.T) {
	repUser := NewRepoUsersPostgres(pgPool)
	serviceUser := service.NewUserService(repUser)
	user, err := serviceUser.Create(ctx, "testUsername", "TestPassword")
	if err != nil {
		t.Fatal("Error in service user", err)
	}
	t.Cleanup(func() {
		if err = repUser.Delete(ctx, user.UserName); err != nil {
			log.WithError(err).Error()
		}
	})
	userFromDB, err := repUser.Get(ctx, user.UserName)
	if err != nil {
		t.Error(err, " Error Get")
	}
	if !reflect.DeepEqual(userFromDB, user) {
		t.Error("Not equal")
	}
}

func TestRepositoryUserDelete(t *testing.T) {
	repUser := NewRepoUsersPostgres(pgPool)
	serviceUser := service.NewUserService(repUser)
	user, err := serviceUser.Create(ctx, "testUsername", "TestPassword")
	if err != nil {
		t.Fatal("Error in service user", err)
	}
	userFromDB, err := repUser.Get(ctx, user.UserName)
	if err != nil {
		t.Error(err, " Error Get")
	}
	err = repUser.Delete(ctx, user.UserName)
	if err != nil {
		t.Error(err)
	}
	NewQueryUser, err := repUser.Get(ctx, user.UserName)
	if err == nil {
		t.Error("User has not deleted")
	}

	if reflect.DeepEqual(userFromDB, NewQueryUser) {
		t.Error("User has not deleted")
	}

}

func TestRepositoryUserGetAll(t *testing.T) {
	repUser := NewRepoUsersPostgres(pgPool)
	serviceUser := service.NewUserService(repUser)
	user1, err := serviceUser.Create(ctx, "testUsername1", "TestPassword")
	if err != nil {
		t.Fatal("Error in service user", err)
	}
	user2, err := serviceUser.Create(ctx, "testUsername2", "TestPassword")
	if err != nil {
		t.Fatal("Error in service user", err)
	}
	t.Cleanup(func() {
		if err = repUser.Delete(ctx, user2.UserName); err != nil {
			log.WithError(err).Error()
		}
		if err = repUser.Delete(ctx, user1.UserName); err != nil {
			log.WithError(err).Error()
		}
	})
	allUsersFromDB, err := repUser.GetAll(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(allUsersFromDB) != 2 {
		t.Error("Count users not equal")
	}
}

func TestRepositoryUserUpdate(t *testing.T) {
	repUser := NewRepoUsersPostgres(pgPool)
	serviceUser := service.NewUserService(repUser)
	user1, err := serviceUser.Create(ctx, "testUsername1", "TestPassword")
	if err != nil {
		t.Fatal("Error in service user", err)
	}
	user1.UserName = "New name"
	userFromDbBeforeUpdate, _ := repUser.Get(ctx, user1.UserName)
	err = repUser.Update(ctx, user1)
	if err != nil {
		t.Error(err)
	}
	userFromDbAfterUpdate, _ := repUser.Get(ctx, user1.UserName)
	if reflect.DeepEqual(userFromDbAfterUpdate, userFromDbBeforeUpdate) {
		t.Error("Update don't work")
	}
	t.Cleanup(func() {
		if err = repUser.Delete(ctx, user1.UserName); err != nil {
			log.WithError(err).Error()
		}
	})
}
