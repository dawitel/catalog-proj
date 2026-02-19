package product

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	activate_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/activate_product"
	productv1 "github.com/dawitel/product-catalog-service/proto/product/v1"
)

func (h *Handler) ActivateProduct(ctx context.Context, req *productv1.ActivateProductRequest) (*productv1.ActivateProductReply, error) {
	if req == nil || req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id required")
	}
	if err := h.Activate.Execute(ctx, activate_product.Request{ProductID: req.GetProductId()}); err != nil {
		return nil, LogRPCError("ActivateProduct", err)
	}
	return &productv1.ActivateProductReply{}, nil
}
