package repository

import (
	"app/internal/models"
	"app/internal/service"
	log "github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
)

func TestRepositoryAuthCreate(t *testing.T) {
	repAuth := NewRepoAuthPostgres(pgPool)
	repUser := NewRepoUsersPostgres(pgPool)
	servUser := service.NewUserService(repUser)
	user, err := servUser.Create(ctx, "unit_tests", "unit_tests")
	if err != nil {
		t.Fatal(err)
	}

	fakeSession := models.Session{
		UserID:          user.ID,
		RfToken:         "Test",
		UniqueSignature: "unit_tests",
		CreatedAt:       time.Now(),
		Disabled:        false,
	}
	err = repAuth.Create(ctx, &fakeSession)
	if err != nil {
		t.Error(err)
	}
	sessionFromDb, err := repAuth.Get(ctx, fakeSession.ID)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		if errRepAuthDelete := repAuth.Delete(ctx, sessionFromDb.ID); errRepAuthDelete != nil {
			log.WithError(errRepAuthDelete).Error()
		}
		if errRepUserDelete := repUser.Delete(ctx, user.UserName); errRepUserDelete != nil {
			log.WithError(errRepUserDelete).Error()
		}
	})
	sessionFromDb.CreatedAt = fakeSession.CreatedAt

	if !reflect.DeepEqual(*sessionFromDb, fakeSession) {
		t.Error("Not equal")
	}

}

func TestRepositoryAuthUpdate(t *testing.T) {
	repAuth := NewRepoAuthPostgres(pgPool)
	repUser := NewRepoUsersPostgres(pgPool)
	servUser := service.NewUserService(repUser)
	user, err := servUser.Create(ctx, "unit_tests", "unit_tests")
	if err != nil {
		t.Fatal(err)
	}

	fakeSession := models.Session{
		UserID:          user.ID,
		RfToken:         "Test",
		UniqueSignature: "unit_tests",
		CreatedAt:       time.Now(),
		Disabled:        false,
	}
	err = repAuth.Create(ctx, &fakeSession)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		if errRepAuthDelete := repAuth.Delete(ctx, fakeSession.ID); errRepAuthDelete != nil {
			log.WithError(errRepAuthDelete).Error()
		}
		if errRepUserDelete := repUser.Delete(ctx, user.UserName); errRepUserDelete != nil {
			log.WithError(errRepUserDelete).Error()
		}
	})
	sessionFromDb, err := repAuth.Get(ctx, fakeSession.ID)
	if err != nil {
		t.Error(err)
	}
	fakeSession.RfToken = "new token"
	err = repAuth.Update(ctx, &fakeSession)
	if err != nil {
		t.Error(err)
	}
	sessionFromDbAfterUpdate, err := repAuth.Get(ctx, fakeSession.ID)
	if err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(sessionFromDb, sessionFromDbAfterUpdate) {
		t.Error("Not updated")
	}
}

func TestRepositoryAuthGet(t *testing.T) {
	repAuth := NewRepoAuthPostgres(pgPool)
	repUser := NewRepoUsersPostgres(pgPool)
	servUser := service.NewUserService(repUser)
	user, err := servUser.Create(ctx, "unit_tests", "unit_tests")
	if err != nil {
		t.Fatal(err)
	}

	fakeSession := models.Session{
		UserID:          user.ID,
		RfToken:         "Test",
		UniqueSignature: "unit_tests",
		CreatedAt:       time.Now(),
		Disabled:        false,
	}
	err = repAuth.Create(ctx, &fakeSession)
	if err != nil {
		t.Error(err)
	}
	sessionFromDb, err := repAuth.Get(ctx, fakeSession.ID)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		if errRepAuthDelete := repAuth.Delete(ctx, sessionFromDb.ID); errRepAuthDelete != nil {
			log.WithError(errRepAuthDelete).Error()
		}
		if errRepUserDelete := repUser.Delete(ctx, user.UserName); errRepUserDelete != nil {
			log.WithError(errRepUserDelete).Error()
		}
	})
	sessionFromDb.CreatedAt = fakeSession.CreatedAt

	if !reflect.DeepEqual(*sessionFromDb, fakeSession) {
		t.Error("Not equal")
	}

}

func TestRepositoryAuthDelete(t *testing.T) {
	repAuth := NewRepoAuthPostgres(pgPool)
	repUser := NewRepoUsersPostgres(pgPool)
	servUser := service.NewUserService(repUser)
	user, err := servUser.Create(ctx, "unit_tests", "unit_tests")
	if err != nil {
		t.Fatal(err)
	}

	fakeSession := models.Session{
		UserID:          user.ID,
		RfToken:         "Test",
		UniqueSignature: "unit_tests",
		CreatedAt:       time.Now(),
		Disabled:        false,
	}
	err = repAuth.Create(ctx, &fakeSession)
	if err != nil {
		t.Error(err)
	}
	sessionFromDb, err := repAuth.Get(ctx, fakeSession.ID)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		if errRepAuthDelete := repAuth.Delete(ctx, sessionFromDb.ID); errRepAuthDelete != nil {
			log.WithError(errRepAuthDelete).Error()
		}
		if errRepUserDelete := repUser.Delete(ctx, user.UserName); errRepUserDelete != nil {
			log.WithError(errRepUserDelete).Error()
		}
	})
	sessionFromDb.CreatedAt = fakeSession.CreatedAt

	if !reflect.DeepEqual(*sessionFromDb, fakeSession) {
		t.Error("Not equal")
	}

	err = repAuth.Delete(ctx, sessionFromDb.ID)

	sessionFromDbAfterDelete, _ := repAuth.Get(ctx, fakeSession.ID)
	if sessionFromDbAfterDelete != nil {
		t.Error("Session didn't delete")
	}

}
