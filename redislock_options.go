package lock

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

const (
	defaultRedisLockTimeout  = 10 * time.Second
	defaultRedisLockResource = "redisLock:default"
	defaultRedisLockToken    = "default"
)

type redisLockOption struct {
	conn     redis.UniversalClient
	resource string
	token    string
	timeout  time.Duration
	err      error
}

type RedisLockOption func(*redisLockOption)

func initRedisLockOptions(options ...RedisLockOption) *redisLockOption {

	option := &redisLockOption{
		resource: defaultRedisLockResource,
		token:    defaultRedisLockToken,
		timeout:  defaultRedisLockTimeout,
	}

	for _, opt := range options {
		opt(option)
	}

	if option.conn == nil {
		option.err = errors.New("'initRedisLockOptions' redis conn is nil")
	}

	return option
}

func WithRedisConn(conn redis.UniversalClient) RedisLockOption {
	return func(option *redisLockOption) {
		option.conn = conn
	}
}

func WithResourceToken(resource, token string) RedisLockOption {
	return func(option *redisLockOption) {
		option.resource = resource
		option.token = token
	}
}

func WithTimeout(timeout time.Duration) RedisLockOption {
	return func(option *redisLockOption) {
		option.timeout = timeout
	}
}
