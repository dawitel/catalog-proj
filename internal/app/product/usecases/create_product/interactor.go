package create_product

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dawitel/catalog-proj/internal/app/product/contracts"
	"github.com/dawitel/catalog-proj/internal/app/product/domain"
	"github.com/dawitel/catalog-proj/internal/commitplan"
	"github.com/dawitel/catalog-proj/internal/pkg/clock"
	"github.com/google/uuid"
)

type Request struct {
	ID                   string
	Name                 string
	Description          string
	Category             string
	BasePriceNumerator   int64
	BasePriceDenominator int64
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

func (it *Interactor) Execute(ctx context.Context, req Request) (string, error) {
	// Golden Mutation: (1) create aggregate, (2) domain validation in constructor,
	// (3) build plan, (4) repo mutations (insert), (5) outbox events, (6) apply plan.
	id := req.ID
	if id == "" {
		id = uuid.New().String()
	}

	now := it.clock.Now()

	basePrice := domain.NewMoney(req.BasePriceNumerator, req.BasePriceDenominator)

	product := domain.NewProduct(id, req.Name, req.Description, req.Category, basePrice, now)
	if product == nil {
		return "", errors.New("invalid product")
	}

	plan := commitplan.NewPlan()
	if mut := it.productRepo.InsertMut(product); mut != nil {
		plan.Add(mut)
	}
	for _, ev := range product.DomainEvents() {
		payload, _ := json.Marshal(eventPayload(ev))
		if mut := it.outboxRepo.InsertMut(uuid.New().String(), ev.EventType(), ev.AggregateID(), string(payload), "pending", ev.OccurredAt()); mut != nil {
			plan.Add(mut)
		}
	}

	if err := it.applier.Apply(ctx, plan); err != nil {
		return "", err
	}

	return product.ID(), nil
}

func eventPayload(ev domain.DomainEvent) map[string]interface{} {
	e, ok := ev.(*domain.ProductCreatedEvent)
	if !ok {
		return map[string]interface{}{"aggregate_id": ev.AggregateID(), "occurred_at": ev.OccurredAt()}
	}

	m := map[string]interface{}{
		"product_id": e.ProductID, "name": e.Name, "description": e.Description, "category": e.Category,
		"occurred_at": e.At,
	}

	if e.BasePrice != nil {
		m["base_price_numerator"] = e.BasePrice.Numerator()
		m["base_price_denominator"] = e.BasePrice.Denominator()
	}

	return m
}
