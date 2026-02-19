package product

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/dawitel/product-catalog-service/internal/pkg/logger"
)

func ToGRPC(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrProductNotActive), errors.Is(err, domain.ErrInvalidDiscountPeriod):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrProductArchived):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func LogRPCError(method string, err error) error {
	if err != nil {
		logger.Warn("RPC error", "method", method, "err", err)
	}
	return ToGRPC(err)
}
