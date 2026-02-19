package services

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"

	"github.com/dawitel/product-catalog-service/internal/app/product/queries/get_product"
	"github.com/dawitel/product-catalog-service/internal/app/product/queries/list_products"
	activate_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/activate_product"
	apply_discount "github.com/dawitel/product-catalog-service/internal/app/product/usecases/apply_discount"
	create_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/create_product"
	deactivate_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/deactivate_product"
	remove_discount "github.com/dawitel/product-catalog-service/internal/app/product/usecases/remove_discount"
	update_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/update_product"
	"github.com/dawitel/product-catalog-service/internal/app/product/repo"
	"github.com/dawitel/product-catalog-service/internal/pkg/clock"
	"github.com/dawitel/product-catalog-service/internal/pkg/committer"
	"github.com/dawitel/product-catalog-service/internal/transport/grpc/product"
)

type Config struct {
	SpannerProject  string
	SpannerInstance string
	SpannerDatabase string
}

func (c Config) DatabasePath() string {
	return fmt.Sprintf("projects/%s/instances/%s/databases/%s", c.SpannerProject, c.SpannerInstance, c.SpannerDatabase)
}

func NewProductHandler(ctx context.Context, cfg Config) (*product.Handler, *spanner.Client, error) {
	client, err := spanner.NewClient(ctx, cfg.DatabasePath())
	if err != nil {
		return nil, nil, err
	}
	comm := committer.New(client)
	clk := clock.Real{}
	productRepo := repo.NewProductRepo(client)
	outboxRepo := repo.NewOutboxRepo()
	readModel := repo.NewReadModel(client)

	createUC := create_product.New(productRepo, outboxRepo, comm, clk)
	updateUC := update_product.New(productRepo, outboxRepo, comm, clk)
	activateUC := activate_product.New(productRepo, outboxRepo, comm, clk)
	deactivateUC := deactivate_product.New(productRepo, outboxRepo, comm, clk)
	applyDiscUC := apply_discount.New(productRepo, outboxRepo, comm, clk)
	removeDiscUC := remove_discount.New(productRepo, outboxRepo, comm, clk)
	getQuery := get_product.New(readModel, clk)
	listQuery := list_products.New(readModel, clk)

	handler := product.NewHandler(createUC, updateUC, activateUC, deactivateUC, applyDiscUC, removeDiscUC, getQuery, listQuery)
	return handler, client, nil
}
