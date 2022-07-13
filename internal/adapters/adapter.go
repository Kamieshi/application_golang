package adapters

import (
	"context"
	"sync"
)

type AdapterInterface interface {
	Start(ctx context.Context, wg *sync.WaitGroup) error
}
