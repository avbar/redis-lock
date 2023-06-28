package locker

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ttl = 60 * time.Second
)

var (
	ErrReleaseLock = errors.New("lock cannot be released")
)

type RedisLock struct {
	redisClient *redis.Client
}

func NewRedisLocker(redisClient *redis.Client) *RedisLock {
	return &RedisLock{
		redisClient: redisClient,
	}
}

func (l *RedisLock) Lock(ctx context.Context, key, token string) (bool, error) {
	return l.redisClient.SetNX(ctx, key, token, ttl).Result()
}

func (l *RedisLock) Unlock(ctx context.Context, key, token string) error {
	releaseLockScript := redis.NewScript(`
		if redis.call("get",KEYS[1]) == ARGV[1] then
			return redis.call("del",KEYS[1])
		else
			return 0
		end
	`)

	res, err := releaseLockScript.Run(ctx, l.redisClient, []string{key}, token).Result()
	if err != nil {
		return err
	}

	if delCount, ok := res.(int64); !ok || delCount == 0 {
		return ErrReleaseLock
	}

	return nil
}
