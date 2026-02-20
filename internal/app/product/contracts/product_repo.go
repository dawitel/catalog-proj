package contracts

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/dawitel/catalog-proj/internal/app/product/domain"
	"github.com/dawitel/catalog-proj/internal/commitplan"
)

type ProductRepo interface {
	InsertMut(p *domain.Product) *spanner.Mutation
	UpdateConditional(p *domain.Product) *commitplan.ConditionalUpdate
	Get(ctx context.Context, id string) (*domain.Product, error)
}
