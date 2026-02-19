package product

import (
	"time"

	"github.com/dawitel/product-catalog-service/internal/app/product/queries/get_product"
	"github.com/dawitel/product-catalog-service/internal/app/product/queries/list_products"
	create_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/create_product"
	update_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/update_product"
	productv1 "github.com/dawitel/product-catalog-service/proto/product/v1"
)

func createRequestFromProto(req *productv1.CreateProductRequest) create_product.Request {
	return create_product.Request{
		Name:                 req.GetName(),
		Description:          req.GetDescription(),
		Category:             req.GetCategory(),
		BasePriceNumerator:   req.GetBasePriceNumerator(),
		BasePriceDenominator: req.GetBasePriceDenominator(),
	}
}

func updateRequestFromProto(req *productv1.UpdateProductRequest) update_product.Request {
	return update_product.Request{
		ProductID:   req.GetProductId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Category:    req.GetCategory(),
	}
}

func applyDiscountRequestFromProto(req *productv1.ApplyDiscountRequest) (productID string, percent int64, start, end time.Time) {
	return req.GetProductId(), req.GetPercent(),
		time.Unix(req.GetStartDateUnix(), 0),
		time.Unix(req.GetEndDateUnix(), 0)
}

func getProductReplyFromDTO(d *get_product.DTO) *productv1.GetProductReply {
	if d == nil {
		return nil
	}
	return &productv1.GetProductReply{
		ProductId:                 d.ID,
		Name:                      d.Name,
		Description:               d.Description,
		Category:                  d.Category,
		BasePriceNumerator:        d.BasePriceNumerator,
		BasePriceDenominator:      d.BasePriceDenominator,
		EffectivePriceNumerator:   d.EffectivePriceNumerator,
		EffectivePriceDenominator: d.EffectivePriceDenominator,
		Status:                    d.Status,
	}
}

func listProductsReplyFromResult(r *list_products.Result) *productv1.ListProductsReply {
	if r == nil {
		return &productv1.ListProductsReply{}
	}
	items := make([]*productv1.ProductItem, 0, len(r.Items))
	for _, it := range r.Items {
		items = append(items, &productv1.ProductItem{
			ProductId:                 it.ID,
			Name:                      it.Name,
			Description:               it.Description,
			Category:                  it.Category,
			BasePriceNumerator:        it.BasePriceNumerator,
			BasePriceDenominator:      it.BasePriceDenominator,
			EffectivePriceNumerator:   it.EffectivePriceNumerator,
			EffectivePriceDenominator: it.EffectivePriceDenominator,
			Status:                    it.Status,
		})
	}
	return &productv1.ListProductsReply{
		Items:         items,
		NextPageToken: r.NextToken,
	}
}
