package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type BlocklistStore struct {
	rdb *redis.Client
}

func (s *BlocklistStore) IsBlocked(ctx context.Context, ip string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	reason, err := s.rdb.Get(ctx, "block:"+ip).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return reason, true, nil
}

func (s *BlocklistStore) Add(ctx context.Context, ip, reason string, ttlSec int) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	return s.rdb.Set(ctx, "block:"+ip, reason, time.Duration(ttlSec)*time.Second).Err()
}
