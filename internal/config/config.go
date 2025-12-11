package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type LogConfig struct {
	Level       string `json:"level"`
	Encoding    string `json:"encoding,optional"`
	Development bool   `json:"development,optional"`
}

type DBConfig struct {
	DataSource string `json:"DataSource"`
}

type CacheTTL struct {
	UserStatusSeconds int `json:"UserStatusSeconds"`
}

type Config struct {
	zrpc.RpcServerConf
	LogConf   LogConfig       `json:"LogConf"`
	RedisConf redis.RedisConf `json:"RedisConf"`
	DBConf    DBConfig        `json:"DBConf"`
	CacheTTL  CacheTTL        `json:"CacheTTL"`
}
