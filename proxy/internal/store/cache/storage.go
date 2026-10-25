package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	queryTimeoutDuration = time.Duration(time.Millisecond * 500)
)

type Storage struct {
	Blocklist interface {
		IsBlocked(ctx context.Context, ip string) (string, bool, error)
		List(ctx context.Context) ([]BlockedEntry, error)
		Add(ctx context.Context, ip, reason string, ttlSec int) error
		Delete(ctx context.Context, ip string) error
	}
	Activity interface {
		IncrWindow(ctx context.Context, key string, windowSec int) (int64, error)
	}
}

func NewStorage(rdb *redis.Client) *Storage {
	return &Storage{
		Blocklist: &BlocklistStore{
			rdb: rdb,
		},
		Activity: &ActivityStore{
			rdb: rdb,
		},
	}
}
