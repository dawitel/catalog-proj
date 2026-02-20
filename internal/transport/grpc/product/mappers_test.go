package product

import (
	"testing"
	"time"

	"github.com/dawitel/catalog-proj/internal/app/product/queries/get_product"
	"github.com/dawitel/catalog-proj/internal/app/product/queries/list_products"
	create_product "github.com/dawitel/catalog-proj/internal/app/product/usecases/create_product"
	update_product "github.com/dawitel/catalog-proj/internal/app/product/usecases/update_product"
	productv1 "github.com/dawitel/catalog-proj/proto/product/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateRequestFromProto(t *testing.T) {
	req := &productv1.CreateProductRequest{
		Name:                 "p1",
		Description:          "desc",
		Category:             "cat",
		BasePriceNumerator:   1999,
		BasePriceDenominator: 100,
	}
	got := createRequestFromProto(req)
	assert.Equal(t, create_product.Request{
		Name:                 "p1",
		Description:          "desc",
		Category:             "cat",
		BasePriceNumerator:   1999,
		BasePriceDenominator: 100,
	}, got)
}

func TestUpdateRequestFromProto(t *testing.T) {
	req := &productv1.UpdateProductRequest{
		ProductId:   "id1",
		Name:        "n",
		Description: "d",
		Category:    "c",
	}
	got := updateRequestFromProto(req)
	assert.Equal(t, update_product.Request{
		ProductID:   "id1",
		Name:        "n",
		Description: "d",
		Category:    "c",
	}, got)
}

func TestApplyDiscountRequestFromProto(t *testing.T) {
	start := time.Unix(1000, 0)
	end := time.Unix(2000, 0)
	req := &productv1.ApplyDiscountRequest{
		ProductId:     "pid",
		Percent:       25,
		StartDateUnix: 1000,
		EndDateUnix:   2000,
	}
	productID, percent, s, e := applyDiscountRequestFromProto(req)
	assert.Equal(t, "pid", productID)
	assert.Equal(t, int64(25), percent)
	assert.True(t, s.Equal(start))
	assert.True(t, e.Equal(end))
}

func TestGetProductReplyFromDTO_Nil(t *testing.T) {
	assert.Nil(t, getProductReplyFromDTO(nil))
}

func TestGetProductReplyFromDTO_NonNil(t *testing.T) {
	d := &get_product.DTO{
		ID:                        "id1",
		Name:                      "n",
		Description:               "d",
		Category:                  "c",
		BasePriceNumerator:        100,
		BasePriceDenominator:      1,
		EffectivePriceNumerator:   80,
		EffectivePriceDenominator: 1,
		Status:                    "active",
	}
	got := getProductReplyFromDTO(d)
	assert.NotNil(t, got)
	assert.Equal(t, "id1", got.ProductId)
	assert.Equal(t, "n", got.Name)
	assert.Equal(t, "d", got.Description)
	assert.Equal(t, "c", got.Category)
	assert.Equal(t, int64(100), got.BasePriceNumerator)
	assert.Equal(t, int64(1), got.BasePriceDenominator)
	assert.Equal(t, int64(80), got.EffectivePriceNumerator)
	assert.Equal(t, int64(1), got.EffectivePriceDenominator)
	assert.Equal(t, "active", got.Status)
}

func TestListProductsReplyFromResult_Nil(t *testing.T) {
	got := listProductsReplyFromResult(nil)
	assert.NotNil(t, got)
	assert.Empty(t, got.Items)
	assert.Empty(t, got.NextPageToken)
}

func TestListProductsReplyFromResult_NonNil(t *testing.T) {
	r := &list_products.Result{
		Items: []list_products.Item{
			{
				ID:                        "id1",
				Name:                      "a",
				Description:               "ad",
				Category:                  "ac",
				BasePriceNumerator:        50,
				BasePriceDenominator:      1,
				EffectivePriceNumerator:   40,
				EffectivePriceDenominator: 1,
				Status:                    "active",
			},
		},
		NextToken: "tok",
	}
	got := listProductsReplyFromResult(r)
	assert.Len(t, got.Items, 1)
	assert.Equal(t, "id1", got.Items[0].ProductId)
	assert.Equal(t, "a", got.Items[0].Name)
	assert.Equal(t, int64(40), got.Items[0].EffectivePriceNumerator)
	assert.Equal(t, int64(1), got.Items[0].EffectivePriceDenominator)
	assert.Equal(t, "tok", got.NextPageToken)
}
