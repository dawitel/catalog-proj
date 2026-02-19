package services

import (
	"math/big"
	"time"

	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
)

// EffectivePrice returns base price with discount applied at time at if discount is valid; otherwise returns base unchanged.
func EffectivePrice(base *domain.Money, discount *domain.Discount, at time.Time) *domain.Money {
	if base == nil {
		return nil
	}
	if discount == nil || !discount.IsValidAt(at) {
		return domain.NewMoneyFromRat(base.Rat())
	}
	price := base.Rat()
	factor := big.NewRat(100-discount.Percentage(), 100)
	effective := new(big.Rat).Mul(price, factor)
	return domain.NewMoneyFromRat(effective)
}
