package get_product

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/dawitel/product-catalog-service/internal/app/product/contracts"
	clockmocks "github.com/dawitel/product-catalog-service/mocks/clock"
	contractsmocks "github.com/dawitel/product-catalog-service/mocks/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestQuery_Execute_SuccessNoDiscount(t *testing.T) {
	now := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	row := &contracts.ProductRow{
		ID:                   "id1",
		Name:                 "p",
		Description:          "d",
		Category:             "c",
		BasePriceNumerator:   10000,
		BasePriceDenominator: 100,
		Status:               "active",
	}
	rm := contractsmocks.NewMockReadModel(t)
	rm.EXPECT().GetProductByID(mock.Anything, "id1").Return(row, nil)
	rm.EXPECT().ListProducts(mock.Anything, mock.Anything, mock.Anything).Maybe()

	clock := clockmocks.NewMockClock(t)
	clock.EXPECT().Now().Return(now)

	q := New(rm, clock)
	dto, err := q.Execute(context.Background(), "id1")
	require.NoError(t, err)
	require.NotNil(t, dto)
	assert.Equal(t, "id1", dto.ID)
	assert.Equal(t, int64(10000), dto.BasePriceNumerator)
	assert.Equal(t, int64(100), dto.BasePriceDenominator)
	assert.Equal(t, int64(100), dto.EffectivePriceNumerator)
	assert.Equal(t, int64(1), dto.EffectivePriceDenominator)
	assert.Equal(t, "active", dto.Status)
}

func TestQuery_Execute_SuccessWithDiscount(t *testing.T) {
	now := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	start := now.Add(-time.Hour)
	end := now.Add(time.Hour)
	rat := new(big.Rat).SetFrac(big.NewInt(20), big.NewInt(100))
	row := &contracts.ProductRow{
		ID:                   "id1",
		Name:                 "p",
		Category:             "c",
		BasePriceNumerator:   10000,
		BasePriceDenominator: 100,
		DiscountPercent:      rat,
		DiscountStartDate:    start,
		DiscountEndDate:      end,
		Status:               "active",
	}
	rm := contractsmocks.NewMockReadModel(t)
	rm.EXPECT().GetProductByID(mock.Anything, "id1").Return(row, nil)
	rm.EXPECT().ListProducts(mock.Anything, mock.Anything, mock.Anything).Maybe()

	clock := clockmocks.NewMockClock(t)
	clock.EXPECT().Now().Return(now)

	q := New(rm, clock)
	dto, err := q.Execute(context.Background(), "id1")
	require.NoError(t, err)
	require.NotNil(t, dto)
	assert.Equal(t, "id1", dto.ID)
	assert.Equal(t, int64(100), dto.EffectivePriceNumerator)
	assert.Equal(t, int64(1), dto.EffectivePriceDenominator)
}

func TestQuery_Execute_ReadModelError(t *testing.T) {
	wantErr := errors.New("not found")
	rm := contractsmocks.NewMockReadModel(t)
	rm.EXPECT().GetProductByID(mock.Anything, "id1").Return(nil, wantErr)
	rm.EXPECT().ListProducts(mock.Anything, mock.Anything, mock.Anything).Maybe()

	clock := clockmocks.NewMockClock(t)

	q := New(rm, clock)
	dto, err := q.Execute(context.Background(), "id1")
	assert.ErrorIs(t, err, wantErr)
	assert.Nil(t, dto)
}
