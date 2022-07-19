package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

const algSign = "HS256"

func TestRegisterUser(t *testing.T) {
	dataJSON, _ := json.Marshal(map[string]string{
		"password": "string",
		"username": "string",
	})
	buf := bytes.NewReader(dataJSON)
	resp, err := http.Post(URLCreateUser, "application/json", buf)
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.WithError(err).Error()
		}
	}()
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

	userRep := repository.NewRepoUsersPostgres(connPullDB)
	t.Cleanup(func() {
		if err = userRep.Delete(ctx, actualUser.UserName); err != nil {
			log.WithError(err).Error()
		}
	})

	expectUser, _ := userRep.Get(ctx, actualUser.UserName)
	assert.Equal(t, expectUser.ID, actualUser.ID)
}

func TestLogin(t *testing.T) {
	userRep := repository.NewRepoUsersPostgres(connPullDB)
	userServ := service.NewUserService(userRep)
	_, err := userServ.Create(ctx, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	dataJSON, _ := json.Marshal(map[string]string{
		"password": "test",
		"username": "test",
	})

	buf := bytes.NewReader(dataJSON)
	resp, err := http.Post(urlLogin, "application/json", buf)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.WithError(err).Error()
		}
	}()

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
		log.WithError(errors.New("error create token")).Error()
	}

	var actualUser models.User
	err = json.NewDecoder(resp.Body).Decode(&actualUser)
	if err != nil {
		t.Fatal(err)
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != algSign {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	}

	token, err := jwt.Parse(AccessData.AccessTk, keyFunc)
	if err != nil {
		t.Fatal(err)
	}
	claims := token.Claims.(jwt.MapClaims)
	idSession := claims["id_session"].(string)
	repAuth := repository.NewRepoAuthPostgres(connPullDB)

	session, err := repAuth.Get(ctx, uuid.MustParse(idSession))
	if err != nil {
		t.Fatalf("Not found session :%s", idSession)
	}
	t.Cleanup(func() {
		if err = repAuth.Delete(ctx, uuid.MustParse(idSession)); err != nil {
			log.WithError(err).Error()
		}
		if err = userRep.Delete(ctx, actualUser.UserName); err != nil {
			log.WithError(err).Error()
		}
	})
	assert.Equal(t, session.RfToken, AccessData.RefreshTk)
}

func TestLogOut(t *testing.T) {
	userRep := repository.NewRepoUsersPostgres(connPullDB)
	userServ := service.NewUserService(userRep)
	_, err := userServ.Create(ctx, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	dataJSON, _ := json.Marshal(map[string]string{
		"password": "test",
		"username": "test",
	})

	buf := bytes.NewReader(dataJSON)
	resp, err := http.Post(urlLogin, "application/json", buf)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.WithError(err).Error()
		}
	}()

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
		t.Fatalf("Error create token, response cod %d", resp.StatusCode)
	}

	var actualUser models.User
	err = json.NewDecoder(resp.Body).Decode(&actualUser)
	if err != nil {
		t.Fatal(err)
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != algSign {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	}

	token, err := jwt.Parse(AccessData.AccessTk, keyFunc)
	if err != nil {
		t.Fatal(err)
	}
	claims := token.Claims.(jwt.MapClaims)
	idSession := claims["id_session"].(string)
	repAuth := repository.NewRepoAuthPostgres(connPullDB)

	t.Cleanup(func() {
		if err = repAuth.Delete(ctx, uuid.MustParse(idSession)); err != nil {
			log.WithError(err).Error()
		}
		if err = userRep.Delete(ctx, actualUser.UserName); err != nil {
			log.WithError(err).Error()
		}
	})

	req, _ = http.NewRequest("GET", urlLogOut, nil)
	req.Header.Add("Authorization", "Bearer "+AccessData.AccessTk)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	userSession, err := repAuth.Get(ctx, uuid.MustParse(idSession))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, userSession.Disabled, true)
}

func TestRefresh(t *testing.T) {
	userRep := repository.NewRepoUsersPostgres(connPullDB)
	userServ := service.NewUserService(userRep)
	_, err := userServ.Create(ctx, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	dataJSON, _ := json.Marshal(map[string]string{
		"password": "test",
		"username": "test",
	})

	buf := bytes.NewReader(dataJSON)
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
		t.Fatalf("Error create token, response cod %d", resp.StatusCode)
	}

	var actualUser models.User
	err = json.NewDecoder(resp.Body).Decode(&actualUser)
	if err != nil {
		t.Fatal(err)
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != algSign {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	}

	token, err := jwt.Parse(AccessData.AccessTk, keyFunc)
	if err != nil {
		t.Fatal(err)
	}
	claims := token.Claims.(jwt.MapClaims)
	idSession := claims["id_session"].(string)
	repAuth := repository.NewRepoAuthPostgres(connPullDB)

	t.Cleanup(func() {
		if err = repAuth.Delete(ctx, uuid.MustParse(idSession)); err != nil {
			log.WithError(err).Error()
		}
		if err = userRep.Delete(ctx, actualUser.UserName); err != nil {
			log.WithError(err).Error()
		}
	})

	sessionBeforeRefresh, err := repAuth.Get(ctx, uuid.MustParse(idSession))
	if err != nil {
		t.Fatal("Get session object before refresh")
	}

	dataJson, _ := json.Marshal(map[string]string{
		"refresh": AccessData.RefreshTk,
	})
	buf = bytes.NewReader(dataJson)
	req, _ = http.NewRequest("POST", urlRefresh, buf)
	req.Header.Add("Authorization", "Bearer "+AccessData.AccessTk)
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatal("Can't get new RfToken")
	}
	sessionAfterRefresh, err := repAuth.Get(ctx, uuid.MustParse(idSession))
	if err != nil {
		t.Fatal("Get session object after refresh")
	}

	assert.NotEqual(t, sessionBeforeRefresh.RfToken, sessionAfterRefresh)
	assert.Equal(t, sessionBeforeRefresh.ID, sessionAfterRefresh.ID)

	type rfTokenFromResponse struct {
		RefreshTk string `json:"refresh"`
	}
	var rfFromResponce rfTokenFromResponse
	err = json.NewDecoder(resp.Body).Decode(&rfFromResponce)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rfFromResponce.RefreshTk, sessionAfterRefresh.RfToken)

}
