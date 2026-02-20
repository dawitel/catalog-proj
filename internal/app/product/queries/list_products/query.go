package list_products

import (
	"context"

	"github.com/dawitel/catalog-proj/internal/app/product/contracts"
	"github.com/dawitel/catalog-proj/internal/app/product/domain"
	"github.com/dawitel/catalog-proj/internal/app/product/domain/services"
	"github.com/dawitel/catalog-proj/internal/pkg/clock"
)

type Query struct {
	readModel contracts.ReadModel
	clock     clock.Clock
}

func New(readModel contracts.ReadModel, c clock.Clock) *Query {
	return &Query{readModel: readModel, clock: c}
}

func (q *Query) Execute(ctx context.Context, filter contracts.ListFilter, page contracts.ListPage) (*Result, error) {
	rmResult, err := q.readModel.ListProducts(ctx, filter, page)
	if err != nil {
		return nil, err
	}
	now := q.clock.Now()
	items := make([]Item, 0, len(rmResult.Items))
	for _, row := range rmResult.Items {
		basePrice := domain.NewMoney(row.BasePriceNumerator, row.BasePriceDenominator)
		var discount *domain.Discount
		if row.DiscountPercent != nil && row.DiscountPercent.Denom().IsInt64() && row.DiscountPercent.Denom().Int64() == 100 && row.DiscountPercent.Num().IsInt64() {
			if pct := row.DiscountPercent.Num().Int64(); pct >= 0 && pct <= 100 {
				discount = domain.NewDiscount(pct, row.DiscountStartDate, row.DiscountEndDate)
			}
		}
		effective := services.EffectivePrice(basePrice, discount, now)
		if effective == nil {
			effective = basePrice
		}
		items = append(items, Item{
			ID:                        row.ID,
			Name:                      row.Name,
			Description:               row.Description,
			Category:                  row.Category,
			BasePriceNumerator:        row.BasePriceNumerator,
			BasePriceDenominator:      row.BasePriceDenominator,
			EffectivePriceNumerator:   effective.Numerator(),
			EffectivePriceDenominator: effective.Denominator(),
			Status:                    row.Status,
		})
	}
	return &Result{Items: items, NextToken: rmResult.NextToken}, nil
}
