package apply_discount

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dawitel/product-catalog-service/internal/app/product/contracts"
	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/dawitel/product-catalog-service/internal/commitplan"
	"github.com/dawitel/product-catalog-service/internal/pkg/clock"
	"github.com/google/uuid"
)

type Request struct {
	ProductID  string
	Percent    int64
	StartDate  time.Time
	EndDate    time.Time
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
	discount := domain.NewDiscount(req.Percent, req.StartDate, req.EndDate)
	if discount == nil {
		return domain.ErrInvalidDiscountPeriod
	}
	if err := product.ApplyDiscount(discount, now); err != nil {
		return err
	}
	plan := commitplan.NewPlan()
	plan.Add(it.productRepo.UpdateMut(product))
	for _, ev := range product.DomainEvents() {
		payload, _ := json.Marshal(eventPayload(ev))
		plan.Add(it.outboxRepo.InsertMut(uuid.New().String(), ev.EventType(), ev.AggregateID(), string(payload), "pending", ev.OccurredAt()))
	}
	return it.applier.Apply(ctx, plan)
}

func eventPayload(ev domain.DomainEvent) map[string]interface{} {
	e, ok := ev.(*domain.DiscountAppliedEvent)
	if !ok {
		return map[string]interface{}{"aggregate_id": ev.AggregateID(), "occurred_at": ev.OccurredAt()}
	}
	return map[string]interface{}{
		"product_id": e.ProductID, "percent": e.Percent,
		"start_date": e.StartDate, "end_date": e.EndDate, "occurred_at": e.At,
	}
}
