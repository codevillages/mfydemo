package service

import (
	"github.com/mfyai/mfydemo/internal/app/user/config"
	"github.com/mfyai/mfydemo/internal/app/user/usecase"
	"github.com/mfyai/mfydemo/internal/infra/cache"
	"github.com/mfyai/mfydemo/internal/infra/persistence"
	"github.com/mfyai/mfydemo/internal/pkg/logger"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// ServiceContext wires dependencies.
type ServiceContext struct {
	Config     config.Config
	Logger     *zap.Logger
	UserUsecase *usecase.UserUsecase
}

// NewServiceContext builds service dependencies.
func NewServiceContext(c config.Config) (*ServiceContext, error) {
	log, err := logger.New(c.Log)
	if err != nil {
		return nil, err
	}
	logger.SetupGoZero(log)

	dbConn := sqlx.NewMysql(c.DB.DataSource)
	redisClient := redis.MustNewRedis(c.Redis)
	userRepo := persistence.NewUserRepo(dbConn)
	userCache := cache.NewUserCacheRepo(redisClient, c.Cache.StatusTTL)

	hashFn := func(raw string) (string, error) {
		hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
		return string(hashed), nil
	}
	compareFn := func(raw, hashed string) error {
		return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
	}

	uc := usecase.NewUserUsecase(userRepo, userCache, log, hashFn, compareFn)

	return &ServiceContext{
		Config:     c,
		Logger:     log,
		UserUsecase: uc,
	}, nil
}
