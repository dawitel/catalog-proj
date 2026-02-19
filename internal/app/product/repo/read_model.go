package repo

import (
	"context"
	"math/big"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dawitel/product-catalog-service/internal/app/product/contracts"
	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/dawitel/product-catalog-service/internal/models/m_product"
)

type ReadModel struct {
	client *spanner.Client
}

func NewReadModel(client *spanner.Client) *ReadModel {
	return &ReadModel{client: client}
}

func (r *ReadModel) GetProductByID(ctx context.Context, id string) (*contracts.ProductRow, error) {
	row, err := r.client.Single().ReadRow(ctx, m_product.TableName, spanner.Key{id}, m_product.Columns())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return spannerRowToProductRow(row)
}

const defaultListPageSize = 20

func (r *ReadModel) ListProducts(ctx context.Context, filter contracts.ListFilter, page contracts.ListPage) (*contracts.ListProductsResult, error) {
	pageSize := page.PageSize
	if pageSize <= 0 {
		pageSize = defaultListPageSize
	}
	limit := pageSize + 1
	stmt := spanner.Statement{
		SQL: `SELECT product_id, name, description, category, base_price_numerator, base_price_denominator,
			discount_percent, discount_start_date, discount_end_date, status, version, created_at, updated_at, archived_at
			FROM products WHERE status = @status AND (@category IS NULL OR category = @category)
			AND (@token = '' OR product_id > @token) ORDER BY product_id LIMIT @limit`,
		Params: map[string]interface{}{
			"status":   "active",
			"category": filter.Category,
			"token":    page.Token,
			"limit":    limit,
		},
	}
	iter := r.client.Single().Query(ctx, stmt)
	defer iter.Stop()
	var items []*contracts.ProductRow
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		pr, err := spannerRowToProductRow(row)
		if err != nil {
			return nil, err
		}
		items = append(items, pr)
	}
	result := &contracts.ListProductsResult{}
	if len(items) > pageSize {
		result.NextToken = items[pageSize-1].ID
		result.Items = items[:pageSize]
	} else {
		result.Items = items
	}
	return result, nil
}

func spannerRowToProductRow(row *spanner.Row) (*contracts.ProductRow, error) {
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
	_ = version
	description := ""
	if desc.Valid {
		description = desc.StringVal
	}
	pr := &contracts.ProductRow{
		ID:                   id,
		Name:                 name,
		Description:          description,
		Category:             category,
		BasePriceNumerator:   baseNum,
		BasePriceDenominator: baseDen,
		DiscountPercent:      discountPercent,
		Status:               status,
	}
	if createdAt.Valid {
		pr.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		pr.UpdatedAt = updatedAt.Time
	}
	if discountStart.Valid {
		pr.DiscountStartDate = discountStart.Time
	}
	if discountEnd.Valid {
		pr.DiscountEndDate = discountEnd.Time
	}
	if archivedAt.Valid {
		pr.ArchivedAt = archivedAt.Time
	}
	return pr, nil
}
