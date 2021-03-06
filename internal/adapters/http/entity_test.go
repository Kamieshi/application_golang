package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"app/internal/models"
	repository "app/internal/repository/posgres"
	"app/internal/service"
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
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(bodyForAuth)
	resp, err := http.Post(urlLogin, "application/json", buf)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.WithError(err).Error()
		}
	}()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		if err = repUser.Delete(ctx, username); err != nil {
			log.WithError(err).Error()
		}
		return nil, fmt.Errorf("Trouble with auth %d", resp.StatusCode)
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
		if err = repAuth.Delete(ctx, idSession); err != nil {
			log.WithError(err).Error()
		}
		if err = repUser.Delete(ctx, username); err != nil {
			log.WithError(err).Error()
		}
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
	req, _ := http.NewRequest("GET", url, http.NoBody)
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
	req, _ := http.NewRequest("DELETE", url, http.NoBody)
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
	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.WithError(err).Error()
		}
	}()
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
		if err = repEntity.Delete(ctx, entityFormRepository.ID.String()); err != nil {
			log.WithError(err).Error()
		}
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
		if err = repEntity.Delete(ctx, entity1.ID.String()); err != nil {
			log.WithError(err).Error()
		}
		if err = repEntity.Delete(ctx, entity2.ID.String()); err != nil {
			log.WithError(err).Error()
		}
	})
	req := maker.GetAuthGet(urlGetAllEntity)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.WithError(err).Error()
		}
	}()

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
		Price:    3,
		IsActive: false,
	}
	repEntity := repository.NewRepoEntityPostgres(connPullDB)
	if err = repEntity.Add(ctx, entityExpected); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err = repEntity.Delete(ctx, entityExpected.ID.String()); err != nil {
			log.WithError(err).Error()
		}
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
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.WithError(err).Error()
		}
	}()
	assert.Equal(t, resp.StatusCode, http.StatusCreated)
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
		if err = repEntity.Delete(ctx, entityExpected.ID.String()); err != nil {
			log.WithError(err).Error()
		}
	})

	req := maker.GetAuthGet(urlGetByIdEntity + entityExpected.ID.String())
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.WithError(err).Error()
		}
	}()
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
		if err = repEntity.Delete(ctx, entityExpected.ID.String()); err != nil {
			log.WithError(err).Error()
		}
	})

	req := maker.GetAuthDelete(urlDeleteEntity + entityExpected.ID.String())
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.WithError(err).Error()
		}
	}()
	assert.Equal(t, resp.StatusCode, http.StatusNoContent)
	ent, err := repEntity.GetForID(ctx, entityExpected.ID.String())
	assert.Error(t, err)
	assert.Nil(t, ent)
}
