package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type ActivityStore struct {
	rdb *redis.Client
}

var incrWindowScript = redis.NewScript(`
	local current = redis.call('INCR', KEYS[1])
	if current == 1 then redis.call('EXPIRE', KEYS[1], ARGV[1]) end
	return current
	`)

func (s *ActivityStore) IncrWindow(ctx context.Context, key string, windowSec int) (int64, error) {
	return incrWindowScript.Run(ctx, s.rdb, []string{key}, windowSec).Int64()
}
