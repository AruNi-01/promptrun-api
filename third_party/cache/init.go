package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"promptrun-api/utils"
	"strconv"
)

var RedisCli *redis.Client

func InitRedis() {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PW"),
		DB:       db,
	})

	_, err := RedisCli.Ping(context.Background()).Result()
	if err != nil {
		utils.Log().Panic("", "Redis 连接失败, errMsg: %s", err.Error())
		panic(err)
	}
}
