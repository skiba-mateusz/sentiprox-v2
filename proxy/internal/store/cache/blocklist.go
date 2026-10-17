package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type BlockedEntry struct {
	IP  string `json:"ip"`
	TTL int    `josn:"ttl"`
}

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

func (s *BlocklistStore) List(ctx context.Context) ([]BlockedEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	results := []BlockedEntry{}

	iter := s.rdb.Scan(ctx, 0, "block:*", 500).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		if key == "" {
			continue
		}

		ttl, err := s.rdb.TTL(ctx, key).Result()
		if err != nil {
			return results, err
		}

		results = append(results, BlockedEntry{
			IP:  strings.TrimPrefix(key, "block:"),
			TTL: int(ttl.Seconds()),
		})
	}

	return results, iter.Err()
}

func (s *BlocklistStore) Add(ctx context.Context, ip, reason string, ttlSec int) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	return s.rdb.Set(ctx, "block:"+ip, reason, time.Duration(ttlSec)*time.Second).Err()
}
