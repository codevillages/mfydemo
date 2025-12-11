package constant

import (
	"errors"
	"fmt"

	model "github.com/mfyai/mfydemo/mos/tidb/model"
)

const (
	RedisUserStatusPrefix = "user:status:"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrVersionConflict = errors.New("user version conflict")
	ErrInvalidArgument = errors.New("invalid argument")
)

const (
	UserStatusUnknown  = model.UserStatusUnknown
	UserStatusActive   = model.UserStatusActive
	UserStatusDisabled = model.UserStatusDisabled
)

func StatusCacheKey(id int64) string {
	return fmt.Sprintf("%s%d", RedisUserStatusPrefix, id)
}
