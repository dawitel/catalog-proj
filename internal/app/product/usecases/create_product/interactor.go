package create_product

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dawitel/product-catalog-service/internal/app/product/contracts"
	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/dawitel/product-catalog-service/internal/commitplan"
	"github.com/dawitel/product-catalog-service/internal/pkg/clock"
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
	plan.Add(it.productRepo.InsertMut(product))
	for _, ev := range product.DomainEvents() {
		payload, _ := json.Marshal(eventPayload(ev))
		plan.Add(it.outboxRepo.InsertMut(uuid.New().String(), ev.EventType(), ev.AggregateID(), string(payload), "pending", ev.OccurredAt()))
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
