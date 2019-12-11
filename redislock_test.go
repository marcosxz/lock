package lock

import (
	"github.com/go-redis/redis"
	"sync"
	"testing"
)

func TestNewRedisLock(t *testing.T) {

	conn := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"127.0.0.1:6379"},
		DB:       0,
		Password: "",
	})

	wg := &sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// get lock
			lock := NewRedisLock(WithRedisConn(conn))
			if err := lock.TryLock(); err != nil {
				t.Error(err)
				t.FailNow()
			}

			// successful
			t.Log("get redis lock successful")

			// handle
			defer func() {
				if err := lock.Unlock(); err != nil {
					t.Error(err)
				}
			}()

			t.Log("handle", i)
		}(i)
	}

	wg.Wait()
}
