package contracts

import (
	"context"
	"math/big"
	"time"
)

type ProductRow struct {
	ID                   string
	Name                 string
	Description          string
	Category             string
	BasePriceNumerator   int64
	BasePriceDenominator int64
	DiscountPercent      *big.Rat
	DiscountStartDate    time.Time
	DiscountEndDate      time.Time
	Status               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	ArchivedAt           time.Time
}

type ListFilter struct {
	Category *string
}

type ListPage struct {
	PageSize int
	Token    string
}

type ListProductsResult struct {
	Items     []*ProductRow
	NextToken string
}

type ReadModel interface {
	GetProductByID(ctx context.Context, id string) (*ProductRow, error)
	ListProducts(ctx context.Context, filter ListFilter, page ListPage) (*ListProductsResult, error)
}
