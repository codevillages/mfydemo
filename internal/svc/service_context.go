package svc

import (
	"os"
	"time"

	"github.com/mfyai/mfydemo/internal/config"
	domainrepo "github.com/mfyai/mfydemo/internal/domain/repository"
	"github.com/mfyai/mfydemo/internal/domain/service"
	"github.com/mfyai/mfydemo/internal/infra/cache"
	"github.com/mfyai/mfydemo/internal/infra/dao"
	infraRepo "github.com/mfyai/mfydemo/internal/infra/repository"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ServiceContext struct {
	Config      config.Config
	Logger      *zap.Logger
	Redis       *redis.Redis
	DB          sqlx.SqlConn
	UserModel   dao.UserModel
	UserCache   *cache.UserStatusCache
	UserRepo    domainrepo.UserRepository
	UserService *service.UserService
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	logger, err := buildLogger(c.LogConf)
	if err != nil {
		return nil, err
	}

	redisClient := redis.MustNewRedis(c.RedisConf)
	dbConn := sqlx.NewMysql(c.DBConf.DataSource)

	userModel := dao.NewUserModel(dbConn)
	userCache := cache.NewUserStatusCache(redisClient, logger)
	userRepo := infraRepo.NewUserRepository(userModel, userCache, logger)
	cacheTTL := time.Duration(c.CacheTTL.UserStatusSeconds) * time.Second
	userService := service.NewUserService(userRepo, cacheTTL, logger)

	return &ServiceContext{
		Config:      c,
		Logger:      logger,
		Redis:       redisClient,
		DB:          dbConn,
		UserModel:   userModel,
		UserCache:   userCache,
		UserRepo:    userRepo,
		UserService: userService,
	}, nil
}

func (s *ServiceContext) SyncLogger() {
	if s.Logger != nil {
		_ = s.Logger.Sync()
	}
}

func buildLogger(conf config.LogConfig) (*zap.Logger, error) {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	level := zap.NewAtomicLevel()
	if conf.Level != "" {
		if err := level.UnmarshalText([]byte(conf.Level)); err != nil {
			return nil, err
		}
	} else {
		level.SetLevel(zap.InfoLevel)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)

	options := []zap.Option{zap.AddCaller()}
	if conf.Development {
		options = append(options, zap.Development())
	}

	return zap.New(core, options...), nil
}
