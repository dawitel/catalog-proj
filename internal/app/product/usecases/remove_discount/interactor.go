package remove_discount

import (
	"context"
	"encoding/json"

	"github.com/dawitel/catalog-proj/internal/app/product/contracts"
	"github.com/dawitel/catalog-proj/internal/commitplan"
	"github.com/dawitel/catalog-proj/internal/pkg/clock"
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
	if cu := it.productRepo.UpdateConditional(product); cu != nil {
		plan.AddConditionalUpdate(cu.Stmt, cu.Params)
	}
	for _, ev := range product.DomainEvents() {
		payload, _ := json.Marshal(map[string]interface{}{"product_id": ev.AggregateID(), "occurred_at": ev.OccurredAt()})
		if mut := it.outboxRepo.InsertMut(uuid.New().String(), ev.EventType(), ev.AggregateID(), string(payload), "pending", ev.OccurredAt()); mut != nil {
			plan.Add(mut)
		}
	}
	return it.applier.Apply(ctx, plan)
}
