package redis

import (
	"context"
	"github.com/go-redis/redis/v9"
	"hyneo/internal/config"
)

func NewClient(ctx context.Context, sc config.Redis) (client *redis.Client, err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     sc.Host + ":" + sc.Port,
		Password: sc.Pass,
		DB:       0,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
