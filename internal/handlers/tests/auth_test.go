package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	_assert := assert.New(t)
	dataJson, _ := json.Marshal(map[string]string{
		"password": "string",
		"username": "string",
	})
	buf := bytes.NewReader(dataJson)
	resp, err := http.Post(addrApi+"/user", "application/json", buf)
	_assert.Nil(err)
	_assert.Equal(202, resp.StatusCode)
	if err != nil {
		fmt.Println(err)
	}
}
