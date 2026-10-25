package store

import (
	"context"
	"time"

	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/db"
	"go.uber.org/zap"
)

const (
	QueryTimeoutDuration = time.Duration(time.Second * 15)
)

type Storage struct {
	Request interface {
		List(ctx context.Context, query PaginatedQuery) (Paginated[[]Request], error)
		Stats(ctx context.Context) (Stats, error)
		Timeseries(ctx context.Context) ([]Timeserie, error)
		Categories(ctx context.Context) ([]Category, error)
		Insert(req *Request)
		Close()
	}
}

func NewStorage(l *zap.SugaredLogger, pg *db.PG) *Storage {
	return &Storage{
		Request: newRequestStore(l, pg),
	}
}
