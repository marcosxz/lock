package lock

import (
	"context"
	"errors"
	"time"
)

type redisLock struct {
	opts *redisLockOption
}

func NewRedisLock(options ...RedisLockOption) *redisLock {
	return &redisLock{opts: initRedisLockOptions(options...)}
}

func (r *redisLock) TryLock() error {
	if r.opts.err != nil {
		return r.opts.err
	}
	ctx, cancel := context.WithTimeout(context.Background(), r.opts.timeout)
	for {
		select {
		case <-ctx.Done():
			return errors.New("redisLock tryLock timeout")
		default:
			if ok, err := r.tryLock(); err != nil || ok {
				cancel()
				return err
			} else {
				// 没有获取到锁继续获取下一次，直到超时
			}
		}
	}
}

func (r *redisLock) Unlock() error {
	var err error
	if r.opts.err != nil {
		err = r.opts.err
	} else {
		_, err = r.opts.conn.Del(r.opts.resource).Result()
	}
	return err
}

func (r *redisLock) AddTimeout(exTime time.Duration) error {
	if ttlTime, err := r.opts.conn.TTL(r.opts.resource).Result(); err != nil {
		return err
	} else if ttlTime <= 0 {
		return errors.New("redisLock AddTimeout failed")
	} else {
		if _, err := r.opts.conn.Set(r.opts.resource, r.opts.token, ttlTime+exTime).Result(); err != nil {
			return err
		}
		return nil
	}
}

func (r *redisLock) tryLock() (ok bool, err error) {
	return r.opts.conn.SetNX(r.opts.resource, r.opts.token, r.opts.timeout).Result()
}
