package contracts

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
)

type ProductRepo interface {
	InsertMut(p *domain.Product) *spanner.Mutation
	UpdateMut(p *domain.Product) *spanner.Mutation
	Get(ctx context.Context, id string) (*domain.Product, error)
}
