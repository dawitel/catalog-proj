package repo

import (
	"context"
	"math/big"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/dawitel/product-catalog-service/internal/models/m_product"
)

type ProductRepo struct {
	client *spanner.Client
}

func NewProductRepo(client *spanner.Client) *ProductRepo {
	return &ProductRepo{client: client}
}

func (r *ProductRepo) InsertMut(p *domain.Product) *spanner.Mutation {
	if p == nil {
		return nil
	}
	row := domainProductToRow(p)
	return m_product.InsertMut(row)
}

// UpdateMut builds an update mutation only from dirty fields in the change tracker.
// Returns nil when there are no changes (avoids no-op writes).
func (r *ProductRepo) UpdateMut(p *domain.Product) *spanner.Mutation {
	if p == nil || p.Changes() == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if p.Changes().Dirty(domain.FieldName) {
		updates[m_product.Name] = p.Name()
	}
	if p.Changes().Dirty(domain.FieldDescription) {
		updates[m_product.Description] = p.Description()
	}
	if p.Changes().Dirty(domain.FieldCategory) {
		updates[m_product.Category] = p.Category()
	}
	if p.Changes().Dirty(domain.FieldBasePrice) {
		if bp := p.BasePrice(); bp != nil {
			updates[m_product.BasePriceNumerator] = bp.Numerator()
			updates[m_product.BasePriceDenominator] = bp.Denominator()
		}
	}
	if p.Changes().Dirty(domain.FieldDiscount) {
		if d := p.Discount(); d != nil {
			updates[m_product.DiscountPercent] = big.NewRat(d.Percentage(), 100)
			updates[m_product.DiscountStartDate] = d.StartDate()
			updates[m_product.DiscountEndDate] = d.EndDate()
		} else {
			updates[m_product.DiscountPercent] = nil
			updates[m_product.DiscountStartDate] = time.Time{}
			updates[m_product.DiscountEndDate] = time.Time{}
		}
	}
	if p.Changes().Dirty(domain.FieldStatus) {
		updates[m_product.Status] = string(p.Status())
	}
	if p.Changes().Dirty(domain.FieldUpdatedAt) {
		updates[m_product.UpdatedAt] = p.UpdatedAt()
	}
	if p.Changes().Dirty(domain.FieldArchivedAt) {
		updates[m_product.ArchivedAt] = p.ArchivedAt()
	}
	if len(updates) == 0 {
		return nil
	}
	return m_product.UpdateMut(p.ID(), updates)
}

func (r *ProductRepo) Get(ctx context.Context, id string) (*domain.Product, error) {
	row, err := r.client.Single().ReadRow(ctx, m_product.TableName, spanner.Key{id}, m_product.Columns())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return rowToDomainProduct(row)
}

func domainProductToRow(p *domain.Product) *m_product.Row {
	row := &m_product.Row{
		ProductID:            p.ID(),
		Name:                 p.Name(),
		Description:          p.Description(),
		Category:             p.Category(),
		BasePriceNumerator:   p.BasePrice().Numerator(),
		BasePriceDenominator: p.BasePrice().Denominator(),
		Status:               string(p.Status()),
		CreatedAt:            p.CreatedAt(),
		UpdatedAt:            p.UpdatedAt(),
		ArchivedAt:           p.ArchivedAt(),
	}
	if d := p.Discount(); d != nil {
		row.DiscountPercent = big.NewRat(d.Percentage(), 100)
		row.DiscountStartDate = d.StartDate()
		row.DiscountEndDate = d.EndDate()
	}
	return row
}

func rowToDomainProduct(row *spanner.Row) (*domain.Product, error) {
	var id, name, category, status string
	var desc spanner.NullString
	var baseNum, baseDen int64
	var discountPercent *big.Rat
	var discountStart, discountEnd spanner.NullTime
	var createdAt, updatedAt time.Time
	var archivedAt spanner.NullTime
	if err := row.Columns(&id, &name, &desc, &category, &baseNum, &baseDen,
		&discountPercent, &discountStart, &discountEnd, &status, &createdAt, &updatedAt, &archivedAt); err != nil {
		return nil, err
	}
	description := ""
	if desc.Valid {
		description = desc.StringVal
	}
	basePrice := domain.NewMoney(baseNum, baseDen)
	var discount *domain.Discount
	if discountPercent != nil && discountStart.Valid && discountEnd.Valid {
		num := discountPercent.Num()
		if num.IsInt64() && num.Int64() >= 0 && num.Int64() <= 100 && discountPercent.Denom().IsInt64() && discountPercent.Denom().Int64() == 1 {
			discount = domain.NewDiscount(num.Int64(), discountStart.Time, discountEnd.Time)
		}
	}
	var archivedAtTime time.Time
	if archivedAt.Valid {
		archivedAtTime = archivedAt.Time
	}
	return domain.RestoreProduct(id, name, description, category, basePrice, discount,
		domain.ProductStatus(status), createdAt, updatedAt, archivedAtTime), nil
}
