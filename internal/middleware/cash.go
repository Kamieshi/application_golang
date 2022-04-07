package middleware

import (
	"app/internal/repository"
	"crypto/sha256"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CashMiddleware struct {
	RepCash repository.CacheRepository
}

func createHashSHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (cm CashMiddleware) ReadNewCache(c echo.Context) func() {
	cont := c
	//TODO Get Body data from response
	val := *c.Response()

	return func() {
		key := createHashSHA256(c.Request().RequestURI)

		cm.RepCash.Set(c.Request().Context(), key, "val")
		logrus.Info("Run func ReadNewCashe")
		logrus.WithFields(logrus.Fields{
			"context": cont,
			"value":   val,
			"key":     key,
		}).Info("values")
	}
}

func (cm CashMiddleware) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Info("run cash middleware")
		key := createHashSHA256(c.Request().RequestURI)
		//_ = cm.RepCash.Delete(c.Request().Context(),key)
		data, err := cm.RepCash.Get(c.Request().Context(), key)
		if err == nil {
			logrus.Info("Get from cache")
			return c.String(http.StatusAccepted, data)
		}

		afterFunc := cm.ReadNewCache(c)
		c.Response().After(afterFunc)
		return next(c)
	}
}
