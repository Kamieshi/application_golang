package adapters

import (
	"context"
	"sync"
)

type AdapterInterface interface {
	Start(wg *sync.WaitGroup, ctx context.Context) error
}
