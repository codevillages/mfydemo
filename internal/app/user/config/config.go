package config

import (
	"time"

	"github.com/mfyai/mfydemo/internal/pkg/logger"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config defines service configuration.
type Config struct {
	zrpc.RpcServerConf `yaml:",inline"`

	DB struct {
		DataSource string `yaml:"dataSource"`
	} `yaml:"mysql"`

	Redis redis.RedisConf `yaml:"redis"`

	Cache struct {
		StatusTTL time.Duration `yaml:"statusTTL"`
	} `yaml:"cache"`

	Log logger.Config `yaml:"log"`
}
