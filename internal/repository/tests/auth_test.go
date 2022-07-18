package tests

import (
	"app/internal/models"
	"app/internal/repository/posgres"
	"app/internal/service"
	"reflect"
	"testing"
	"time"
)

func TestRepositoryAuthCreate(t *testing.T) {
	repAuth := repository.NewRepoAuthPostgres(pgPool)
	repUser := repository.NewRepoUsersPostgres(pgPool)
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
		repAuth.Delete(ctx, sessionFromDb.ID)
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
		repAuth.Delete(ctx, fakeSession.ID)
		repUser.Delete(ctx, user.UserName)

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
	repAuth := repository.NewRepoAuthPostgres(pgPool)
	repUser := repository.NewRepoUsersPostgres(pgPool)
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
		repAuth.Delete(ctx, sessionFromDb.ID)
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
		repAuth.Delete(ctx, sessionFromDb.ID)
		repUser.Delete(ctx, user.UserName)

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
