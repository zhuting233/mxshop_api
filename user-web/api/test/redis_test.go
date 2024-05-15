package t_test

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "192.168.57.128", 6379),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//fmt.Println(global.ServerConfig.RedisInfo.Expire)
	rdb.Set(context.Background(), "17384309176", 321123, time.Duration(300)*time.Second)
	val, err := rdb.Get(context.Background(), "17384309176").Result()
	if err != nil {
		panic(err)
	}
	t.Log(val)
}
