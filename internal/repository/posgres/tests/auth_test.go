package tests

import (
	"app/internal/models"
	repository "app/internal/repository/posgres"
	"app/internal/service"
	"reflect"
	"testing"
	"time"
)

func TestRepositoryAuthCreate(t *testing.T) {
	repAuth := repository.NewRepoAuthPostgres(pgPool)
	repUser := repository.NewRepoUsersPostgres(pgPool)
	servUser := service.NewUserService(repUser)
	user, err := servUser.Create(ctx, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	fakeSession := models.Session{
		UserId:          user.ID,
		RfToken:         "Test",
		UniqueSignature: "test",
		CreatedAt:       time.Now(),
		Disabled:        false,
	}
	err = repAuth.Create(ctx, &fakeSession)
	if err != nil {
		t.Error(err)
	}
	sessionFromDb, err := repAuth.Get(ctx, fakeSession.ID.String())
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		repAuth.Delete(ctx, sessionFromDb.ID.String())
		repUser.Delete(ctx, user.UserName)

	})
	sessionFromDb.CreatedAt = fakeSession.CreatedAt

	if !reflect.DeepEqual(*sessionFromDb, fakeSession) {
		t.Error("Not equal")
	}

}

func TestRepositoryAuthUpdate(t *testing.T) {
	repAuth := repository.NewRepoAuthPostgres(pgPool)
	repUser := repository.NewRepoUsersPostgres(pgPool)
	servUser := service.NewUserService(repUser)
	user, err := servUser.Create(ctx, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	fakeSession := models.Session{
		UserId:          user.ID,
		RfToken:         "Test",
		UniqueSignature: "test",
		CreatedAt:       time.Now(),
		Disabled:        false,
	}
	err = repAuth.Create(ctx, &fakeSession)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		repAuth.Delete(ctx, fakeSession.ID.String())
		repUser.Delete(ctx, user.UserName)

	})
	sessionFromDb, err := repAuth.Get(ctx, fakeSession.ID.String())
	if err != nil {
		t.Error(err)
	}
	fakeSession.RfToken = "new token"
	err = repAuth.Update(ctx, &fakeSession)
	if err != nil {
		t.Error(err)
	}
	sessionFromDbAfterUpdate, err := repAuth.Get(ctx, fakeSession.ID.String())
	if err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(sessionFromDb, sessionFromDbAfterUpdate) {
		t.Error("Not updated")
	}
}

func TestRepositoryAuthGet(t *testing.T) {
	repAuth := repository.NewRepoAuthPostgres(pgPool)
	repUser := repository.NewRepoUsersPostgres(pgPool)
	servUser := service.NewUserService(repUser)
	user, err := servUser.Create(ctx, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	fakeSession := models.Session{
		UserId:          user.ID,
		RfToken:         "Test",
		UniqueSignature: "test",
		CreatedAt:       time.Now(),
		Disabled:        false,
	}
	err = repAuth.Create(ctx, &fakeSession)
	if err != nil {
		t.Error(err)
	}
	sessionFromDb, err := repAuth.Get(ctx, fakeSession.ID.String())
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		repAuth.Delete(ctx, sessionFromDb.ID.String())
		repUser.Delete(ctx, user.UserName)

	})
	sessionFromDb.CreatedAt = fakeSession.CreatedAt

	if !reflect.DeepEqual(*sessionFromDb, fakeSession) {
		t.Error("Not equal")
	}

}

func TestRepositoryAuthDelete(t *testing.T) {
	repAuth := repository.NewRepoAuthPostgres(pgPool)
	repUser := repository.NewRepoUsersPostgres(pgPool)
	servUser := service.NewUserService(repUser)
	user, err := servUser.Create(ctx, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	fakeSession := models.Session{
		UserId:          user.ID,
		RfToken:         "Test",
		UniqueSignature: "test",
		CreatedAt:       time.Now(),
		Disabled:        false,
	}
	err = repAuth.Create(ctx, &fakeSession)
	if err != nil {
		t.Error(err)
	}
	sessionFromDb, err := repAuth.Get(ctx, fakeSession.ID.String())
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		repAuth.Delete(ctx, sessionFromDb.ID.String())
		repUser.Delete(ctx, user.UserName)

	})
	sessionFromDb.CreatedAt = fakeSession.CreatedAt

	if !reflect.DeepEqual(*sessionFromDb, fakeSession) {
		t.Error("Not equal")
	}

	err = repAuth.Delete(ctx, sessionFromDb.ID.String())

	sessionFromDbAfterDelete, _ := repAuth.Get(ctx, fakeSession.ID.String())
	if sessionFromDbAfterDelete != nil {
		t.Error("Session didn't delete")
	}

}
