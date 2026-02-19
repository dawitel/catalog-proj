package product

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	productv1 "github.com/dawitel/product-catalog-service/proto/product/v1"
)

func (h *Handler) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.GetProductReply, error) {
	if req == nil || req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id required")
	}
	dto, err := h.Get.Execute(ctx, req.GetProductId())
	if err != nil {
		return nil, ToGRPC(err)
	}
	return getProductReplyFromDTO(dto), nil
}
