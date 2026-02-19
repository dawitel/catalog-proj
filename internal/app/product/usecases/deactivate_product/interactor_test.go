package deactivate_product

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

func TestDeactivateProduct_Success(t *testing.T) {
	now := time.Now()
	base := domain.NewMoney(100, 1)
	product := domain.NewProduct("id1", "p", "d", "c", base, now)

	productRepo := contractsmocks.NewMockProductRepo(t)
	productRepo.EXPECT().Get(mock.Anything, "id1").Return(product, nil)
	productRepo.EXPECT().UpdateMut(mock.Anything).RunAndReturn(func(p *domain.Product) *spanner.Mutation {
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
	err := it.Execute(context.Background(), Request{ProductID: "id1"})
	require.NoError(t, err)
	require.NotNil(t, capturedPlan)
	assert.GreaterOrEqual(t, len(capturedPlan.Mutations()), 2)
	assert.Equal(t, domain.ProductStatusInactive, product.Status())
}

func TestDeactivateProduct_GetError(t *testing.T) {
	productRepo := contractsmocks.NewMockProductRepo(t)
	productRepo.EXPECT().Get(mock.Anything, "id1").Return(nil, errors.New("not found"))

	outboxRepo := contractsmocks.NewMockOutboxRepo(t)
	applier := commitplanmocks.NewMockApplier(t)
	clock := clockmocks.NewMockClock(t)

	it := New(productRepo, outboxRepo, applier, clock)
	err := it.Execute(context.Background(), Request{ProductID: "id1"})
	assert.Error(t, err)
}

func TestDeactivateProduct_Archived(t *testing.T) {
	now := time.Now()
	base := domain.NewMoney(100, 1)
	product := domain.NewProduct("id1", "p", "d", "c", base, now)
	product.Archive(now)

	productRepo := contractsmocks.NewMockProductRepo(t)
	productRepo.EXPECT().Get(mock.Anything, "id1").Return(product, nil)

	outboxRepo := contractsmocks.NewMockOutboxRepo(t)
	applier := commitplanmocks.NewMockApplier(t)
	clock := clockmocks.NewMockClock(t)
	clock.EXPECT().Now().Return(now)

	it := New(productRepo, outboxRepo, applier, clock)
	err := it.Execute(context.Background(), Request{ProductID: "id1"})
	assert.ErrorIs(t, err, domain.ErrProductArchived)
}
