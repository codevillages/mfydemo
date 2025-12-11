package user

import (
	"errors"

	"github.com/mfyai/mfydemo/pkg/constant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, constant.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, constant.ErrVersionConflict):
		return status.Error(codes.Aborted, err.Error())
	case errors.Is(err, constant.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
