package remove_discount

import (
	"context"
	"encoding/json"

	"github.com/dawitel/product-catalog-service/internal/app/product/contracts"
	"github.com/dawitel/product-catalog-service/internal/commitplan"
	"github.com/dawitel/product-catalog-service/internal/pkg/clock"
	"github.com/google/uuid"
)

type Request struct {
	ProductID string
}

type Interactor struct {
	productRepo contracts.ProductRepo
	outboxRepo  contracts.OutboxRepo
	applier     commitplan.Applier
	clock       clock.Clock
}

func New(productRepo contracts.ProductRepo, outboxRepo contracts.OutboxRepo, applier commitplan.Applier, c clock.Clock) *Interactor {
	return &Interactor{productRepo: productRepo, outboxRepo: outboxRepo, applier: applier, clock: c}
}

func (it *Interactor) Execute(ctx context.Context, req Request) error {
	product, err := it.productRepo.Get(ctx, req.ProductID)
	if err != nil {
		return err
	}
	now := it.clock.Now()
	if err := product.RemoveDiscount(now); err != nil {
		return err
	}
	plan := commitplan.NewPlan()
	plan.Add(it.productRepo.UpdateMut(product))
	for _, ev := range product.DomainEvents() {
		payload, _ := json.Marshal(map[string]interface{}{"product_id": ev.AggregateID(), "occurred_at": ev.OccurredAt()})
		plan.Add(it.outboxRepo.InsertMut(uuid.New().String(), ev.EventType(), ev.AggregateID(), string(payload), "pending", ev.OccurredAt()))
	}
	return it.applier.Apply(ctx, plan)
}
