package repository

import "context"

type CacheRepository interface {
	Set(c context.Context, key string, value string) error
	Get(c context.Context, key string) (string, error)
	Delete(c context.Context, key string) error
}
