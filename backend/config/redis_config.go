package config

import (
	"context"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// InitRedisClient 初始化Redis客户端
func InitRedisClient() (*redis.Client, error) {
	addr := GetEnv("REDIS_ADDR", "47.121.210.209:16379")
	password := GetEnv("REDIS_PASSWORD", "123456")
	db, _ := strconv.Atoi(GetEnv("REDIS_DB", "0"))

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Redis连接成功")
	return client, nil
}
