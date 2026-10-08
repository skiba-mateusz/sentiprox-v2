package store

import (
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/db"
	"go.uber.org/zap"
)

type Storage struct {
	Request interface {
		Insert(req *Request)
		Close()
	}
}

func NewStorage(l *zap.SugaredLogger, pg *db.PG) *Storage {
	return &Storage{
		Request: newRequestStore(l, pg),
	}
}
