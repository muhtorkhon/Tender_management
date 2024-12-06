package redise

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)



type RedisDB struct {
	Rdb *redis.Client
}

func NewRedis(client *redis.Client) *RedisDB {
	return &RedisDB{
		Rdb: client,
	}
}

func (r *RedisDB) Ping() error {
	return r.Rdb.Ping(context.Background()).Err()
}

func (r *RedisDB) Set(ctx context.Context, key string, value interface{}) error {
	return r.Rdb.Set(ctx, key, value, 0).Err()
}

func (r *RedisDB) Get(ctx context.Context, key string) (string, error) {
	return r.Rdb.Get(ctx, key).Result()
}

func (r *RedisDB) Delete(ctx context.Context, key string) error {
	return r.Rdb.Del(ctx, key).Err()
}

func (r *RedisDB) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *RedisDB) SetEx(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	return r.Rdb.SetEx(ctx, key, value, duration).Err()
}