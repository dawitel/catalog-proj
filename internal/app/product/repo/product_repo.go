package repo

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/dawitel/product-catalog-service/internal/commitplan"
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

func (r *ProductRepo) UpdateConditional(p *domain.Product) *commitplan.ConditionalUpdate {
	if p == nil || p.Changes() == nil {
		return nil
	}
	var setParts []string
	params := map[string]interface{}{
		"product_id": p.ID(),
		"version":    p.Version(),
	}
	if p.Changes().Dirty(domain.FieldName) {
		setParts = append(setParts, "name = @name")
		params["name"] = p.Name()
	}
	if p.Changes().Dirty(domain.FieldDescription) {
		setParts = append(setParts, "description = @description")
		params["description"] = p.Description()
	}
	if p.Changes().Dirty(domain.FieldCategory) {
		setParts = append(setParts, "category = @category")
		params["category"] = p.Category()
	}
	if p.Changes().Dirty(domain.FieldBasePrice) {
		if bp := p.BasePrice(); bp != nil {
			setParts = append(setParts, "base_price_numerator = @base_price_numerator", "base_price_denominator = @base_price_denominator")
			params["base_price_numerator"] = bp.Numerator()
			params["base_price_denominator"] = bp.Denominator()
		}
	}
	if p.Changes().Dirty(domain.FieldDiscount) {
		setParts = append(setParts, "discount_percent = @discount_percent", "discount_start_date = @discount_start_date", "discount_end_date = @discount_end_date")
		if d := p.Discount(); d != nil {
			params["discount_percent"] = big.NewRat(d.Percentage(), 100)
			params["discount_start_date"] = d.StartDate()
			params["discount_end_date"] = d.EndDate()
		} else {
			params["discount_percent"] = nil
			params["discount_start_date"] = time.Time{}
			params["discount_end_date"] = time.Time{}
		}
	}
	if p.Changes().Dirty(domain.FieldStatus) {
		setParts = append(setParts, "status = @status")
		params["status"] = string(p.Status())
	}
	if p.Changes().Dirty(domain.FieldUpdatedAt) {
		setParts = append(setParts, "updated_at = @updated_at")
		params["updated_at"] = p.UpdatedAt()
	}
	if p.Changes().Dirty(domain.FieldArchivedAt) {
		setParts = append(setParts, "archived_at = @archived_at")
		params["archived_at"] = p.ArchivedAt()
	}
	if len(setParts) == 0 {
		return nil
	}
	setParts = append(setParts, "version = version + 1")
	stmt := fmt.Sprintf("UPDATE products SET %s WHERE product_id = @product_id AND version = @version", strings.Join(setParts, ", "))
	return &commitplan.ConditionalUpdate{Stmt: stmt, Params: params}
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
		Version:              1,
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
	var version int64
	var createdAt, updatedAt, archivedAt spanner.NullTime
	if err := row.Columns(&id, &name, &desc, &category, &baseNum, &baseDen,
		&discountPercent, &discountStart, &discountEnd, &status, &version, &createdAt, &updatedAt, &archivedAt); err != nil {
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
	var createdAtTime, updatedAtTime, archivedAtTime time.Time
	if createdAt.Valid {
		createdAtTime = createdAt.Time
	}
	if updatedAt.Valid {
		updatedAtTime = updatedAt.Time
	}
	if archivedAt.Valid {
		archivedAtTime = archivedAt.Time
	}
	return domain.RestoreProduct(id, name, description, category, basePrice, discount,
		domain.ProductStatus(status), version, createdAtTime, updatedAtTime, archivedAtTime), nil
}
