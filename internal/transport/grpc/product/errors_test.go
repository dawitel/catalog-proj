package product

import (
	"errors"
	"testing"

	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToGRPC_Nil(t *testing.T) {
	assert.Nil(t, ToGRPC(nil))
}

func TestToGRPC_ErrProductNotFound(t *testing.T) {
	err := ToGRPC(domain.ErrProductNotFound)
	assert.NotNil(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestToGRPC_ErrInvalidProduct(t *testing.T) {
	err := ToGRPC(domain.ErrInvalidProduct)
	assert.NotNil(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestToGRPC_ErrProductNotActive(t *testing.T) {
	err := ToGRPC(domain.ErrProductNotActive)
	assert.NotNil(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
}

func TestToGRPC_ErrInvalidDiscountPeriod(t *testing.T) {
	err := ToGRPC(domain.ErrInvalidDiscountPeriod)
	assert.NotNil(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
}

func TestToGRPC_ErrProductArchived(t *testing.T) {
	err := ToGRPC(domain.ErrProductArchived)
	assert.NotNil(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
}

func TestToGRPC_UnknownError(t *testing.T) {
	unknown := errors.New("something failed")
	err := ToGRPC(unknown)
	assert.NotNil(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}
