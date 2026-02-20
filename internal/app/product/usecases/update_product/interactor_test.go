package update_product

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/dawitel/catalog-proj/internal/app/product/domain"
	"github.com/dawitel/catalog-proj/internal/commitplan"
	clockmocks "github.com/dawitel/catalog-proj/mocks/clock"
	commitplanmocks "github.com/dawitel/catalog-proj/mocks/commitplan"
	contractsmocks "github.com/dawitel/catalog-proj/mocks/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateProduct_Success(t *testing.T) {
	now := time.Now()
	base := domain.NewMoney(100, 1)
	product := domain.NewProduct("id1", "old", "oldd", "oldc", base, now)

	productRepo := contractsmocks.NewMockProductRepo(t)
	productRepo.EXPECT().Get(mock.Anything, "id1").Return(product, nil)
	productRepo.EXPECT().UpdateConditional(mock.Anything).RunAndReturn(func(p *domain.Product) *commitplan.ConditionalUpdate {
		return &commitplan.ConditionalUpdate{Stmt: "UPDATE products SET name = @name WHERE product_id = @product_id AND version = @version", Params: map[string]interface{}{}}
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
	err := it.Execute(context.Background(), Request{ProductID: "id1", Name: "new", Description: "newd", Category: "newc"})
	require.NoError(t, err)
	require.NotNil(t, capturedPlan)
	assert.GreaterOrEqual(t, len(capturedPlan.ConditionalUpdates()), 1)
	assert.GreaterOrEqual(t, len(capturedPlan.Mutations()), 1)
	assert.Equal(t, "new", product.Name())
}

func TestUpdateProduct_GetError(t *testing.T) {
	productRepo := contractsmocks.NewMockProductRepo(t)
	productRepo.EXPECT().Get(mock.Anything, "id1").Return(nil, errors.New("not found"))

	outboxRepo := contractsmocks.NewMockOutboxRepo(t)
	applier := commitplanmocks.NewMockApplier(t)
	clock := clockmocks.NewMockClock(t)

	it := New(productRepo, outboxRepo, applier, clock)
	err := it.Execute(context.Background(), Request{ProductID: "id1", Name: "n", Category: "c"})
	assert.Error(t, err)
}

func TestUpdateProduct_Archived(t *testing.T) {
	now := time.Now()
	base := domain.NewMoney(100, 1)
	product := domain.NewProduct("id1", "n", "d", "c", base, now)
	product.Archive(now)

	productRepo := contractsmocks.NewMockProductRepo(t)
	productRepo.EXPECT().Get(mock.Anything, "id1").Return(product, nil)

	outboxRepo := contractsmocks.NewMockOutboxRepo(t)
	applier := commitplanmocks.NewMockApplier(t)
	clock := clockmocks.NewMockClock(t)
	clock.EXPECT().Now().Return(now)

	it := New(productRepo, outboxRepo, applier, clock)
	err := it.Execute(context.Background(), Request{ProductID: "id1", Name: "x", Category: "y"})
	assert.ErrorIs(t, err, domain.ErrProductArchived)
}
