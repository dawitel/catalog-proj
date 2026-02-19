package m_product

import (
	"math/big"
	"time"

	"cloud.google.com/go/spanner"
)

type Row struct {
	ProductID            string
	Name                 string
	Description          string
	Category             string
	BasePriceNumerator   int64
	BasePriceDenominator int64
	DiscountPercent      *big.Rat
	DiscountStartDate    time.Time
	DiscountEndDate      time.Time
	Status               string
	Version              int64
	CreatedAt            time.Time
	UpdatedAt            time.Time
	ArchivedAt           time.Time
}

func InsertMut(r *Row) *spanner.Mutation {
	return spanner.Insert(TableName, []string{
		ProductID, Name, Description, Category,
		BasePriceNumerator, BasePriceDenominator,
		DiscountPercent, DiscountStartDate, DiscountEndDate,
		Status, Version, CreatedAt, UpdatedAt, ArchivedAt,
	}, []interface{}{
		r.ProductID, r.Name, r.Description, r.Category,
		r.BasePriceNumerator, r.BasePriceDenominator,
		r.DiscountPercent, r.DiscountStartDate, r.DiscountEndDate,
		r.Status, r.Version, r.CreatedAt, r.UpdatedAt, r.ArchivedAt,
	})
}

func UpdateMut(productID string, updates map[string]interface{}) *spanner.Mutation {
	if len(updates) == 0 {
		return nil
	}
	cols := []string{ProductID}
	vals := []interface{}{productID}
	for _, k := range []string{Name, Description, Category, BasePriceNumerator, BasePriceDenominator, DiscountPercent, DiscountStartDate, DiscountEndDate, Status, CreatedAt, UpdatedAt, ArchivedAt} {
		if v, ok := updates[k]; ok {
			cols = append(cols, k)
			vals = append(vals, v)
		}
	}
	return spanner.Update(TableName, cols, vals)
}
