package redis

import (
	"context"
	"gin-admin-server/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func InitRedis() *redis.Client {
	redisConfig := global.GNA_CONFIG.Redis
	if redisConfig.Addr == "" {
		global.GNA_LOG.Error("在 yaml 文件中配置 Redis")
		return nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password, // no password set
		DB:       redisConfig.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.GNA_LOG.Error("redis 连接失败, err:", zap.Error(err))
		panic(err)
	}
	global.GNA_LOG.Info("redis 连接成功:", zap.String("pong", pong))
	return client
}
