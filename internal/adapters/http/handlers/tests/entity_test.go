package tests

import (
	"app/internal/models"
	repository "app/internal/repository/posgres"
	"app/internal/service"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

type MakerAuthRequest struct {
	username    string
	password    string
	accessToken string
	rfToken     string
	cleaner     func()
}

func NewMaker(username, password string) (*MakerAuthRequest, error) {
	repUser := repository.NewRepoUsersPostgres(connPullDB)
	servUser := service.NewUserService(repUser)
	_, err := servUser.Create(ctx, username, password)
	if err != nil {
		return nil, err
	}
	bodyForAuth, err := json.Marshal(map[string]string{
		"password": password,
		"username": username,
	},
	)
	buf := bytes.NewReader(bodyForAuth)
	resp, err := http.Post(urlLogin, "application/json", buf)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		repUser.Delete(ctx, username)
		return nil, errors.New(fmt.Sprintf("Trouble with auth %d", resp.StatusCode))
	}

	type responseData struct {
		AccessTk  string `json:"access"`
		RefreshTk string `json:"refresh"`
	}
	var AccessData responseData
	err = json.NewDecoder(resp.Body).Decode(&AccessData)
	if err != nil {
		return nil, err
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	}

	token, err := jwt.Parse(AccessData.AccessTk, keyFunc)
	if err != nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)
	idSession := uuid.MustParse(claims["id_session"].(string))

	repAuth := repository.NewRepoAuthPostgres(connPullDB)

	cleanerFunc := func() {
		repAuth.Delete(ctx, idSession)
		repUser.Delete(ctx, username)
	}
	return &MakerAuthRequest{
		username:    username,
		password:    password,
		accessToken: AccessData.AccessTk,
		rfToken:     AccessData.RefreshTk,
		cleaner:     cleanerFunc,
	}, err
}

func (m *MakerAuthRequest) GetAuthPOST(url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Authorization", "Bearer "+m.accessToken)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func (m *MakerAuthRequest) GetAuthGet(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+m.accessToken)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func (m *MakerAuthRequest) GetAuthPUT(url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest("PUT", url, body)
	req.Header.Add("Authorization", "Bearer "+m.accessToken)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func (m *MakerAuthRequest) GetAuthDelete(url string) *http.Request {
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+m.accessToken)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func TestCreate(t *testing.T) {
	maker, err := NewMaker("test", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer maker.cleaner()

	entity := models.Entity{
		Name:     "test",
		Price:    100,
		IsActive: false,
	}
	dataForWrite, err := json.Marshal(entity)
	if err != nil {
		t.Fatal(err)
	}
	bufferWrite := bytes.NewReader(dataForWrite)

	req := maker.GetAuthPOST(urlCreateEntity, bufferWrite)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Response status : %d", resp.StatusCode)
	}
	var actualEntity models.Entity
	err = json.NewDecoder(resp.Body).Decode(&actualEntity)
	if err != nil {
		t.Fatal("Failed decode response body")
	}
	repEntity := repository.NewRepoEntityPostgres(connPullDB)
	entityFormRepository, err := repEntity.GetForID(ctx, actualEntity.ID.String())
	if err != nil {
		t.Fatal("Failed get entity from Repository")
	}
	t.Cleanup(func() {
		repEntity.Delete(ctx, entityFormRepository.ID.String())
	})
	assert.Equal(t, entityFormRepository.ID, actualEntity.ID)
	assert.Equal(t, entityFormRepository.Name, actualEntity.Name)
}

func TestGetAll(t *testing.T) {
	maker, err := NewMaker("test", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer maker.cleaner()

	entity1 := &models.Entity{
		Name:     "ent1",
		Price:    0,
		IsActive: false,
	}
	entity2 := &models.Entity{
		Name:     "ent12",
		Price:    0,
		IsActive: false,
	}

	repEntity := repository.NewRepoEntityPostgres(connPullDB)

	if err = repEntity.Add(ctx, entity1); err != nil {
		t.Fatal(err)
	}
	if err = repEntity.Add(ctx, entity2); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		repEntity.Delete(ctx, entity1.ID.String())
		repEntity.Delete(ctx, entity2.ID.String())
	})
	req := maker.GetAuthGet(urlGetAllEntity)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Error request ", resp.StatusCode)
	}
	var entities []models.Entity
	if err = json.NewDecoder(resp.Body).Decode(&entities); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(entities), 2)
}

func TestUpdate(t *testing.T) {
	maker, err := NewMaker("test", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer maker.cleaner()

	entityExpected := &models.Entity{
		Name:     "ent1",
		Price:    0,
		IsActive: false,
	}
	repEntity := repository.NewRepoEntityPostgres(connPullDB)
	if err = repEntity.Add(ctx, entityExpected); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		repEntity.Delete(ctx, entityExpected.ID.String())
	})

	entityBeforeUpdate, err := repEntity.GetForID(ctx, entityExpected.ID.String())

	entityExpected.Name = "New value"
	dataForWrite, err := json.Marshal(entityExpected)
	if err != nil {
		t.Fatal(err)
	}
	bufferRead := bytes.NewReader(dataForWrite)

	req := maker.GetAuthPUT(urlGetByIdEntity+entityExpected.ID.String(), bufferRead)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Error request ", resp.StatusCode)
	}
	entityAfterUpdate, err := repEntity.GetForID(ctx, entityExpected.ID.String())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, entityAfterUpdate.ID, entityBeforeUpdate.ID)
	assert.NotEqual(t, entityAfterUpdate.Name, entityBeforeUpdate.Name)
}

func TestGetByID(t *testing.T) {
	maker, err := NewMaker("test", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer maker.cleaner()

	entityExpected := &models.Entity{
		Name:     "ent1",
		Price:    0,
		IsActive: false,
	}
	repEntity := repository.NewRepoEntityPostgres(connPullDB)
	if err = repEntity.Add(ctx, entityExpected); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		repEntity.Delete(ctx, entityExpected.ID.String())
	})

	req := maker.GetAuthGet(urlGetByIdEntity + entityExpected.ID.String())
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Error request ", resp.StatusCode)
	}

	var entitiesActualy models.Entity
	if err = json.NewDecoder(resp.Body).Decode(&entitiesActualy); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, entitiesActualy.ID, entityExpected.ID)
	assert.Equal(t, entitiesActualy.Name, entityExpected.Name)
}

func TestDelete(t *testing.T) {
	maker, err := NewMaker("test", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer maker.cleaner()

	entityExpected := &models.Entity{
		Name:     "ent1",
		Price:    0,
		IsActive: false,
	}
	repEntity := repository.NewRepoEntityPostgres(connPullDB)
	if err = repEntity.Add(ctx, entityExpected); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		repEntity.Delete(ctx, entityExpected.ID.String())
	})

	req := maker.GetAuthDelete(urlDeleteEntity + entityExpected.ID.String())
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Error request ", resp.StatusCode)
	}

	ent, err := repEntity.GetForID(ctx, entityExpected.ID.String())
	assert.Error(t, err)
	assert.Nil(t, ent)

}
