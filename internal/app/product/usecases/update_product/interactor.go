package update_product

import (
	"context"
	"encoding/json"

	"github.com/dawitel/product-catalog-service/internal/app/product/contracts"
	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/dawitel/product-catalog-service/internal/commitplan"
	"github.com/dawitel/product-catalog-service/internal/pkg/clock"
	"github.com/google/uuid"
)

type Request struct {
	ProductID   string
	Name        string
	Description string
	Category    string
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
	if err := product.UpdateDetails(req.Name, req.Description, req.Category, now); err != nil {
		return err
	}
	plan := commitplan.NewPlan()
	if mut := it.productRepo.UpdateMut(product); mut != nil {
		plan.Add(mut)
	}
	for _, ev := range product.DomainEvents() {
		payload, _ := json.Marshal(eventPayload(ev))
		if mut := it.outboxRepo.InsertMut(uuid.New().String(), ev.EventType(), ev.AggregateID(), string(payload), "pending", ev.OccurredAt()); mut != nil {
			plan.Add(mut)
		}
	}
	return it.applier.Apply(ctx, plan)
}

func eventPayload(ev domain.DomainEvent) map[string]interface{} {
	e, ok := ev.(*domain.ProductUpdatedEvent)
	if !ok {
		return map[string]interface{}{"aggregate_id": ev.AggregateID(), "occurred_at": ev.OccurredAt()}
	}
	return map[string]interface{}{
		"product_id": e.ProductID, "name": e.Name, "description": e.Description, "category": e.Category,
		"occurred_at": e.At,
	}
}
