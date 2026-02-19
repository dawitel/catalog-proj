package get_product

import (
	"context"

	"github.com/dawitel/product-catalog-service/internal/app/product/contracts"
	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/dawitel/product-catalog-service/internal/app/product/domain/services"
	"github.com/dawitel/product-catalog-service/internal/pkg/clock"
)

type Query struct {
	readModel contracts.ReadModel
	clock     clock.Clock
}

func New(readModel contracts.ReadModel, c clock.Clock) *Query {
	return &Query{readModel: readModel, clock: c}
}

func (q *Query) Execute(ctx context.Context, productID string) (*DTO, error) {
	row, err := q.readModel.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	basePrice := domain.NewMoney(row.BasePriceNumerator, row.BasePriceDenominator)
	var discount *domain.Discount
	if row.DiscountPercent != nil && row.DiscountPercent.Denom().IsInt64() && row.DiscountPercent.Denom().Int64() == 100 && row.DiscountPercent.Num().IsInt64() {
		if pct := row.DiscountPercent.Num().Int64(); pct >= 0 && pct <= 100 {
			discount = domain.NewDiscount(pct, row.DiscountStartDate, row.DiscountEndDate)
		}
	}
	now := q.clock.Now()
	effective := services.EffectivePrice(basePrice, discount, now)
	if effective == nil {
		effective = basePrice
	}
	return &DTO{
		ID:                        row.ID,
		Name:                      row.Name,
		Description:               row.Description,
		Category:                  row.Category,
		BasePriceNumerator:        row.BasePriceNumerator,
		BasePriceDenominator:      row.BasePriceDenominator,
		EffectivePriceNumerator:   effective.Numerator(),
		EffectivePriceDenominator: effective.Denominator(),
		Status:                    row.Status,
	}, nil
}
