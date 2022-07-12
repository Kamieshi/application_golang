package tests

import (
	"app/internal/models"
	repository "app/internal/repository/posgres"
	"app/internal/service"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	dataJson, _ := json.Marshal(map[string]string{
		"password": "string",
		"username": "string",
	})
	buf := bytes.NewReader(dataJson)
	resp, err := http.Post(urlCreateUser, "application/json", buf)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 202, resp.StatusCode)
	if err != nil {
		t.Fatalf("Send response error :%s", err)
	}

	var actualUser models.User
	err = json.NewDecoder(resp.Body).Decode(&actualUser)
	if err != nil {
		t.Fatal(err)
	}

	userRep := repository.NewRepoUsersPostgres(connPullDb)
	t.Cleanup(func() {
		userRep.Delete(ctx, actualUser.UserName)
	})

	expectUser, _ := userRep.Get(ctx, actualUser.UserName)
	assert.Equal(t, expectUser.ID, actualUser.ID)
}

func TestLogin(t *testing.T) {

	userRep := repository.NewRepoUsersPostgres(connPullDb)
	userServ := service.NewUserService(userRep)
	_, err := userServ.Create(ctx, "test", "test")
	if err != nil {
		t.Fatal(err)
	}
	dataJson, _ := json.Marshal(map[string]string{
		"password": "test",
		"username": "test",
	})

	buf := bytes.NewReader(dataJson)
	resp, err := http.Post(urlLogin, "application/json", buf)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		t.Fatalf("Response code: %d", resp.StatusCode)
	}
	type responseData struct {
		AccessTk  string `json:"access"`
		RefreshTk string `json:"refresh"`
	}
	var AccessData responseData
	err = json.NewDecoder(resp.Body).Decode(&AccessData)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", urlCheckAuth, nil)
	req.Header.Add("Authorization", "Bearer "+AccessData.AccessTk)

	client := http.DefaultClient
	resp, err = client.Do(req)
	if resp.StatusCode != http.StatusAccepted {
		log.Fatal("Error create token")
	}

	var actualUser models.User
	err = json.NewDecoder(resp.Body).Decode(&actualUser)
	if err != nil {
		t.Fatal(err)
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	}

	token, err := jwt.Parse(AccessData.AccessTk, keyFunc)
	if err != nil {
		t.Fatal(err)
	}
	claims := token.Claims.(jwt.MapClaims)
	id_session := claims["id_session"].(string)
	repAuth := repository.NewRepoAuthPostgres(connPullDb)

	session, err := repAuth.Get(ctx, id_session)
	if err != nil {
		t.Fatalf("Not found session :%s", id_session)
	}
	assert.Equal(t, session.RfToken, AccessData.RefreshTk)

}
