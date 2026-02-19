package list_products

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

func TestQuery_Execute_Success(t *testing.T) {
	now := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	result := &contracts.ListProductsResult{
		Items: []*contracts.ProductRow{
			{
				ID:                   "id1",
				Name:                 "a",
				Category:             "c1",
				BasePriceNumerator:   50,
				BasePriceDenominator: 1,
				Status:               "active",
			},
		},
		NextToken: "tok",
	}
	rm := contractsmocks.NewMockReadModel(t)
	rm.EXPECT().GetProductByID(mock.Anything, mock.Anything).Maybe()
	rm.EXPECT().ListProducts(mock.Anything, contracts.ListFilter{}, contracts.ListPage{PageSize: 10}).Return(result, nil)

	clock := clockmocks.NewMockClock(t)
	clock.EXPECT().Now().Return(now)

	q := New(rm, clock)
	res, err := q.Execute(context.Background(), contracts.ListFilter{}, contracts.ListPage{PageSize: 10})
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res.Items, 1)
	assert.Equal(t, "id1", res.Items[0].ID)
	assert.Equal(t, int64(50), res.Items[0].BasePriceNumerator)
	assert.Equal(t, int64(50), res.Items[0].EffectivePriceNumerator)
	assert.Equal(t, "tok", res.NextToken)
}

func TestQuery_Execute_WithDiscount(t *testing.T) {
	now := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	start := now.Add(-time.Hour)
	end := now.Add(time.Hour)
	result := &contracts.ListProductsResult{
		Items: []*contracts.ProductRow{
			{
				ID:                   "id1",
				Name:                 "a",
				BasePriceNumerator:   100,
				BasePriceDenominator: 1,
				DiscountPercent:      big.NewRat(25, 100),
				DiscountStartDate:    start,
				DiscountEndDate:      end,
				Status:               "active",
			},
		},
	}
	rm := contractsmocks.NewMockReadModel(t)
	rm.EXPECT().GetProductByID(mock.Anything, mock.Anything).Maybe()
	rm.EXPECT().ListProducts(mock.Anything, contracts.ListFilter{}, contracts.ListPage{}).Return(result, nil)

	clock := clockmocks.NewMockClock(t)
	clock.EXPECT().Now().Return(now)

	q := New(rm, clock)
	res, err := q.Execute(context.Background(), contracts.ListFilter{}, contracts.ListPage{})
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res.Items, 1)
	assert.Equal(t, int64(100), res.Items[0].EffectivePriceNumerator)
	assert.Equal(t, int64(1), res.Items[0].EffectivePriceDenominator)
}

func TestQuery_Execute_ReadModelError(t *testing.T) {
	wantErr := errors.New("db error")
	rm := contractsmocks.NewMockReadModel(t)
	rm.EXPECT().GetProductByID(mock.Anything, mock.Anything).Maybe()
	rm.EXPECT().ListProducts(mock.Anything, mock.Anything, mock.Anything).Return(nil, wantErr)

	clock := clockmocks.NewMockClock(t)

	q := New(rm, clock)
	res, err := q.Execute(context.Background(), contracts.ListFilter{}, contracts.ListPage{})
	assert.ErrorIs(t, err, wantErr)
	assert.Nil(t, res)
}
