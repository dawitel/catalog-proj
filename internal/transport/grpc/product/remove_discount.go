package product

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	remove_discount "github.com/dawitel/product-catalog-service/internal/app/product/usecases/remove_discount"
	productv1 "github.com/dawitel/product-catalog-service/proto/product/v1"
)

func (h *Handler) RemoveDiscount(ctx context.Context, req *productv1.RemoveDiscountRequest) (*productv1.RemoveDiscountReply, error) {
	if req == nil || req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id required")
	}
	if err := h.RemoveDisc.Execute(ctx, remove_discount.Request{ProductID: req.GetProductId()}); err != nil {
		return nil, ToGRPC(err)
	}
	return &productv1.RemoveDiscountReply{}, nil
}
