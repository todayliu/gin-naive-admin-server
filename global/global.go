package global

import (
	"gin-admin-server/config"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var (
	GNA_VIPER               *viper.Viper
	GNA_CONFIG              config.Server
	GNA_DB                  *gorm.DB
	GNA_LOG                 *zap.Logger
	GNA_REDIS               *redis.Client
	GNA_Concurrency_Control = &singleflight.Group{}
)
