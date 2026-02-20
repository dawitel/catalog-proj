package product

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	deactivate_product "github.com/dawitel/catalog-proj/internal/app/product/usecases/deactivate_product"
	productv1 "github.com/dawitel/catalog-proj/proto/product/v1"
)

func (h *Handler) DeactivateProduct(ctx context.Context, req *productv1.DeactivateProductRequest) (*productv1.DeactivateProductReply, error) {
	if req == nil || req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id required")
	}
	if err := h.Deactivate.Execute(ctx, deactivate_product.Request{ProductID: req.GetProductId()}); err != nil {
		return nil, LogRPCError("DeactivateProduct", err)
	}
	return &productv1.DeactivateProductReply{}, nil
}
