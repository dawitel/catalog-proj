package product

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	productv1 "github.com/dawitel/product-catalog-service/proto/product/v1"
)

func (h *Handler) UpdateProduct(ctx context.Context, req *productv1.UpdateProductRequest) (*productv1.UpdateProductReply, error) {
	if req == nil || req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id required")
	}
	appReq := updateRequestFromProto(req)
	if err := h.Update.Execute(ctx, appReq); err != nil {
		return nil, ToGRPC(err)
	}
	return &productv1.UpdateProductReply{}, nil
}
