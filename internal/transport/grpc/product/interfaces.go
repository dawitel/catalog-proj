package product

import (
	"context"

	"github.com/dawitel/product-catalog-service/internal/app/product/contracts"
	"github.com/dawitel/product-catalog-service/internal/app/product/queries/get_product"
	"github.com/dawitel/product-catalog-service/internal/app/product/queries/list_products"
	activate_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/activate_product"
	apply_discount "github.com/dawitel/product-catalog-service/internal/app/product/usecases/apply_discount"
	archive_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/archive_product"
	create_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/create_product"
	deactivate_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/deactivate_product"
	remove_discount "github.com/dawitel/product-catalog-service/internal/app/product/usecases/remove_discount"
	update_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/update_product"
)

type CreateRunner interface {
	Execute(ctx context.Context, req create_product.Request) (string, error)
}

type UpdateRunner interface {
	Execute(ctx context.Context, req update_product.Request) error
}

type ActivateRunner interface {
	Execute(ctx context.Context, req activate_product.Request) error
}

type DeactivateRunner interface {
	Execute(ctx context.Context, req deactivate_product.Request) error
}

type ArchiveRunner interface {
	Execute(ctx context.Context, req archive_product.Request) error
}

type ApplyDiscountRunner interface {
	Execute(ctx context.Context, req apply_discount.Request) error
}

type RemoveDiscountRunner interface {
	Execute(ctx context.Context, req remove_discount.Request) error
}

type GetProductRunner interface {
	Execute(ctx context.Context, productID string) (*get_product.DTO, error)
}

type ListProductsRunner interface {
	Execute(ctx context.Context, filter contracts.ListFilter, page contracts.ListPage) (*list_products.Result, error)
}
