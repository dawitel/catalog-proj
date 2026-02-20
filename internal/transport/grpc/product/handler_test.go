package product

import (
	"context"
	"errors"
	"testing"

	"github.com/dawitel/catalog-proj/internal/app/product/contracts"
	"github.com/dawitel/catalog-proj/internal/app/product/domain"
	"github.com/dawitel/catalog-proj/internal/app/product/queries/get_product"
	"github.com/dawitel/catalog-proj/internal/app/product/queries/list_products"
	productmocks "github.com/dawitel/catalog-proj/mocks/transport/grpc/product"
	productv1 "github.com/dawitel/catalog-proj/proto/product/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func makeInertMocks(t *testing.T) (
	*productmocks.MockCreateRunner,
	*productmocks.MockUpdateRunner,
	*productmocks.MockActivateRunner,
	*productmocks.MockDeactivateRunner,
	*productmocks.MockArchiveRunner,
	*productmocks.MockApplyDiscountRunner,
	*productmocks.MockRemoveDiscountRunner,
	*productmocks.MockGetProductRunner,
	*productmocks.MockListProductsRunner,
) {
	create := productmocks.NewMockCreateRunner(t)
	update := productmocks.NewMockUpdateRunner(t)
	activate := productmocks.NewMockActivateRunner(t)
	deactivate := productmocks.NewMockDeactivateRunner(t)
	archive := productmocks.NewMockArchiveRunner(t)
	applyDisc := productmocks.NewMockApplyDiscountRunner(t)
	removeDisc := productmocks.NewMockRemoveDiscountRunner(t)
	getProduct := productmocks.NewMockGetProductRunner(t)
	listProducts := productmocks.NewMockListProductsRunner(t)
	return create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts
}

func TestHandler_CreateProduct_Validation(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.CreateProduct(ctx, nil)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.CreateProduct(ctx, &productv1.CreateProductRequest{Name: "x"})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.CreateProduct(ctx, &productv1.CreateProductRequest{Name: "x", Category: "c", BasePriceDenominator: 0})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestHandler_CreateProduct_Success(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	create.EXPECT().Execute(mock.Anything, mock.Anything).Return("new-id", nil)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	rep, err := h.CreateProduct(ctx, &productv1.CreateProductRequest{
		Name:                 "p",
		Category:             "c",
		BasePriceNumerator:   100,
		BasePriceDenominator: 1,
	})
	require.NoError(t, err)
	assert.Equal(t, "new-id", rep.ProductId)
}

func TestHandler_CreateProduct_Error(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	create.EXPECT().Execute(mock.Anything, mock.Anything).Return("", errors.New("db error"))
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.CreateProduct(ctx, &productv1.CreateProductRequest{
		Name:                 "p",
		Category:             "c",
		BasePriceNumerator:   100,
		BasePriceDenominator: 1,
	})
	require.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestHandler_UpdateProduct_Validation(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.UpdateProduct(ctx, nil)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.UpdateProduct(ctx, &productv1.UpdateProductRequest{})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestHandler_UpdateProduct_Success(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	update.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.UpdateProduct(ctx, &productv1.UpdateProductRequest{ProductId: "id1", Name: "n", Category: "c"})
	require.NoError(t, err)
}

func TestHandler_UpdateProduct_DomainError(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	update.EXPECT().Execute(mock.Anything, mock.Anything).Return(domain.ErrProductArchived)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.UpdateProduct(ctx, &productv1.UpdateProductRequest{ProductId: "id1", Name: "n", Category: "c"})
	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestHandler_GetProduct_Validation(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.GetProduct(ctx, nil)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.GetProduct(ctx, &productv1.GetProductRequest{})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestHandler_GetProduct_Success(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	getProduct.EXPECT().Execute(mock.Anything, "id1").Return(&get_product.DTO{ID: "id1", Name: "p", Category: "c", BasePriceNumerator: 100, BasePriceDenominator: 1, EffectivePriceNumerator: 100, EffectivePriceDenominator: 1, Status: "active"}, nil)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	rep, err := h.GetProduct(ctx, &productv1.GetProductRequest{ProductId: "id1"})
	require.NoError(t, err)
	assert.Equal(t, "id1", rep.ProductId)
	assert.Equal(t, "p", rep.Name)
	assert.Equal(t, int64(100), rep.EffectivePriceNumerator)
}

func TestHandler_GetProduct_NotFound(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	getProduct.EXPECT().Execute(mock.Anything, "id1").Return(nil, domain.ErrProductNotFound)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.GetProduct(ctx, &productv1.GetProductRequest{ProductId: "id1"})
	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestHandler_ListProducts_Validation(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.ListProducts(ctx, nil)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestHandler_ListProducts_Success(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	listProducts.EXPECT().Execute(mock.Anything, contracts.ListFilter{}, contracts.ListPage{PageSize: 10, Token: ""}).
		Return(&list_products.Result{Items: []list_products.Item{{ID: "id1", Name: "p", Status: "active"}}, NextToken: "tok"}, nil)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	rep, err := h.ListProducts(ctx, &productv1.ListProductsRequest{PageSize: 10})
	require.NoError(t, err)
	require.Len(t, rep.Items, 1)
	assert.Equal(t, "id1", rep.Items[0].ProductId)
	assert.Equal(t, "tok", rep.NextPageToken)
}

func TestHandler_ActivateProduct_Validation(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.ActivateProduct(ctx, nil)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.ActivateProduct(ctx, &productv1.ActivateProductRequest{})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestHandler_ActivateProduct_Success(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	activate.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.ActivateProduct(ctx, &productv1.ActivateProductRequest{ProductId: "id1"})
	require.NoError(t, err)
}

func TestHandler_DeactivateProduct_Validation(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.DeactivateProduct(ctx, nil)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestHandler_DeactivateProduct_Success(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	deactivate.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.DeactivateProduct(ctx, &productv1.DeactivateProductRequest{ProductId: "id1"})
	require.NoError(t, err)
}

func TestHandler_ApplyDiscount_Validation(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.ApplyDiscount(ctx, nil)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.ApplyDiscount(ctx, &productv1.ApplyDiscountRequest{})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestHandler_ApplyDiscount_Success(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	applyDisc.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.ApplyDiscount(ctx, &productv1.ApplyDiscountRequest{ProductId: "id1", Percent: 10, StartDateUnix: 0, EndDateUnix: 1})
	require.NoError(t, err)
}

func TestHandler_ApplyDiscount_DomainError(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	applyDisc.EXPECT().Execute(mock.Anything, mock.Anything).Return(domain.ErrProductNotActive)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.ApplyDiscount(ctx, &productv1.ApplyDiscountRequest{ProductId: "id1", Percent: 10, StartDateUnix: 0, EndDateUnix: 1})
	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestHandler_RemoveDiscount_Validation(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.RemoveDiscount(ctx, nil)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestHandler_RemoveDiscount_Success(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	removeDisc.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.RemoveDiscount(ctx, &productv1.RemoveDiscountRequest{ProductId: "id1"})
	require.NoError(t, err)
}

func TestHandler_ArchiveProduct_Validation(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.ArchiveProduct(ctx, nil)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = h.ArchiveProduct(ctx, &productv1.ArchiveProductRequest{})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestHandler_ArchiveProduct_Success(t *testing.T) {
	create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts := makeInertMocks(t)
	archive.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil)
	h := NewHandler(create, update, activate, deactivate, archive, applyDisc, removeDisc, getProduct, listProducts)
	ctx := context.Background()
	_, err := h.ArchiveProduct(ctx, &productv1.ArchiveProductRequest{ProductId: "id1"})
	require.NoError(t, err)
}
