package create_product

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/dawitel/product-catalog-service/internal/commitplan"
	clockmocks "github.com/dawitel/product-catalog-service/mocks/clock"
	commitplanmocks "github.com/dawitel/product-catalog-service/mocks/commitplan"
	contractsmocks "github.com/dawitel/product-catalog-service/mocks/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct_Success(t *testing.T) {
	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	productRepo := contractsmocks.NewMockProductRepo(t)
	productRepo.EXPECT().InsertMut(mock.Anything).RunAndReturn(func(p *domain.Product) *spanner.Mutation {
		return spanner.InsertMap("products", map[string]interface{}{"id": p.ID()})
	})

	outboxRepo := contractsmocks.NewMockOutboxRepo(t)
	outboxRepo.EXPECT().InsertMut(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(_, _, _, _, _ string, _ time.Time) *spanner.Mutation {
			return spanner.InsertMap("outbox_events", map[string]interface{}{"event_id": "x"})
		})

	var capturedPlan *commitplan.Plan
	applier := commitplanmocks.NewMockApplier(t)
	applier.EXPECT().Apply(mock.Anything, mock.Anything).Run(func(_ context.Context, plan *commitplan.Plan) {
		capturedPlan = plan
	}).Return(nil)

	clock := clockmocks.NewMockClock(t)
	clock.EXPECT().Now().Return(now)

	it := New(productRepo, outboxRepo, applier, clock)
	id, err := it.Execute(context.Background(), Request{
		Name:                 "p",
		Description:          "d",
		Category:             "c",
		BasePriceNumerator:   100,
		BasePriceDenominator: 1,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, id)
	require.NotNil(t, capturedPlan)
	assert.Len(t, capturedPlan.Mutations(), 2)
}

func TestCreateProduct_SuccessWithID(t *testing.T) {
	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	productRepo := contractsmocks.NewMockProductRepo(t)
	productRepo.EXPECT().InsertMut(mock.Anything).RunAndReturn(func(p *domain.Product) *spanner.Mutation {
		return spanner.InsertMap("products", map[string]interface{}{"id": p.ID()})
	})

	outboxRepo := contractsmocks.NewMockOutboxRepo(t)
	outboxRepo.EXPECT().InsertMut(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(_, _, _, _, _ string, _ time.Time) *spanner.Mutation {
			return spanner.InsertMap("outbox_events", map[string]interface{}{"event_id": "x"})
		})

	applier := commitplanmocks.NewMockApplier(t)
	applier.EXPECT().Apply(mock.Anything, mock.Anything).Return(nil)

	clock := clockmocks.NewMockClock(t)
	clock.EXPECT().Now().Return(now)

	it := New(productRepo, outboxRepo, applier, clock)
	id, err := it.Execute(context.Background(), Request{
		ID:                   "custom-id",
		Name:                 "p",
		Category:             "c",
		BasePriceNumerator:   100,
		BasePriceDenominator: 1,
	})
	require.NoError(t, err)
	assert.Equal(t, "custom-id", id)
}

func TestCreateProduct_InvalidBasePrice(t *testing.T) {
	productRepo := contractsmocks.NewMockProductRepo(t)
	outboxRepo := contractsmocks.NewMockOutboxRepo(t)
	applier := commitplanmocks.NewMockApplier(t)
	clock := clockmocks.NewMockClock(t)
	clock.EXPECT().Now().Return(time.Now())

	it := New(productRepo, outboxRepo, applier, clock)
	id, err := it.Execute(context.Background(), Request{
		Name:                 "p",
		Category:             "c",
		BasePriceNumerator:   100,
		BasePriceDenominator: 0,
	})
	assert.Error(t, err)
	assert.Empty(t, id)
}

func TestCreateProduct_CommitterError(t *testing.T) {
	applyErr := errors.New("apply failed")
	now := time.Now()

	productRepo := contractsmocks.NewMockProductRepo(t)
	productRepo.EXPECT().InsertMut(mock.Anything).RunAndReturn(func(p *domain.Product) *spanner.Mutation {
		return spanner.InsertMap("products", map[string]interface{}{"id": p.ID()})
	})

	outboxRepo := contractsmocks.NewMockOutboxRepo(t)
	outboxRepo.EXPECT().InsertMut(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(_, _, _, _, _ string, _ time.Time) *spanner.Mutation {
			return spanner.InsertMap("outbox_events", map[string]interface{}{"event_id": "x"})
		})

	applier := commitplanmocks.NewMockApplier(t)
	applier.EXPECT().Apply(mock.Anything, mock.Anything).Return(applyErr)

	clock := clockmocks.NewMockClock(t)
	clock.EXPECT().Now().Return(now)

	it := New(productRepo, outboxRepo, applier, clock)
	id, err := it.Execute(context.Background(), Request{
		Name:                 "p",
		Category:             "c",
		BasePriceNumerator:   100,
		BasePriceDenominator: 1,
	})
	assert.ErrorIs(t, err, applyErr)
	assert.Empty(t, id)
}
