package repository

import (
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"app/internal/models"
	"app/internal/service"
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
	sessionFromDB, err := repAuth.Get(ctx, fakeSession.ID)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		if errRepAuthDelete := repAuth.Delete(ctx, sessionFromDB.ID); errRepAuthDelete != nil {
			logrus.WithError(errRepAuthDelete).Error()
		}
		if errRepUserDelete := repUser.Delete(ctx, user.UserName); errRepUserDelete != nil {
			logrus.WithError(errRepUserDelete).Error()
		}
	})
	sessionFromDB.CreatedAt = fakeSession.CreatedAt

	if !reflect.DeepEqual(*sessionFromDB, fakeSession) {
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
			logrus.WithError(errRepAuthDelete).Error()
		}
		if errRepUserDelete := repUser.Delete(ctx, user.UserName); errRepUserDelete != nil {
			logrus.WithError(errRepUserDelete).Error()
		}
	})
	sessionFromDB, err := repAuth.Get(ctx, fakeSession.ID)
	if err != nil {
		t.Error(err)
	}
	fakeSession.RfToken = "new token"
	err = repAuth.Update(ctx, &fakeSession)
	if err != nil {
		t.Error(err)
	}
	sessionFromDBAfterUpdate, err := repAuth.Get(ctx, fakeSession.ID)
	if err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(sessionFromDB, sessionFromDBAfterUpdate) {
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
	sessionFromDB, err := repAuth.Get(ctx, fakeSession.ID)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		if errRepAuthDelete := repAuth.Delete(ctx, sessionFromDB.ID); errRepAuthDelete != nil {
			logrus.WithError(errRepAuthDelete).Error()
		}
		if errRepUserDelete := repUser.Delete(ctx, user.UserName); errRepUserDelete != nil {
			logrus.WithError(errRepUserDelete).Error()
		}
	})
	sessionFromDB.CreatedAt = fakeSession.CreatedAt

	if !reflect.DeepEqual(*sessionFromDB, fakeSession) {
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
	sessionFromDB, err := repAuth.Get(ctx, fakeSession.ID)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		if errRepAuthDelete := repAuth.Delete(ctx, sessionFromDB.ID); errRepAuthDelete != nil {
			logrus.WithError(errRepAuthDelete).Error()
		}
		if errRepUserDelete := repUser.Delete(ctx, user.UserName); errRepUserDelete != nil {
			logrus.WithError(errRepUserDelete).Error()
		}
	})
	sessionFromDB.CreatedAt = fakeSession.CreatedAt

	if !reflect.DeepEqual(*sessionFromDB, fakeSession) {
		t.Error("Not equal")
	}
	if err = repAuth.Delete(ctx, sessionFromDB.ID); err != nil {
		t.Error(err)
	}
	sessionFromDBAfterDelete, _ := repAuth.Get(ctx, fakeSession.ID)
	if sessionFromDBAfterDelete != nil {
		t.Error("Session didn't delete")
	}
}
