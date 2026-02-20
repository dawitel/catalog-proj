package product

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apply_discount "github.com/dawitel/catalog-proj/internal/app/product/usecases/apply_discount"
	productv1 "github.com/dawitel/catalog-proj/proto/product/v1"
)

func (h *Handler) ApplyDiscount(ctx context.Context, req *productv1.ApplyDiscountRequest) (*productv1.ApplyDiscountReply, error) {
	if req == nil || req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id required")
	}
	productID, percent, start, end := applyDiscountRequestFromProto(req)
	if err := h.ApplyDisc.Execute(ctx, apply_discount.Request{
		ProductID: productID,
		Percent:   percent,
		StartDate: start,
		EndDate:   end,
	}); err != nil {
		return nil, LogRPCError("ApplyDiscount", err)
	}
	return &productv1.ApplyDiscountReply{}, nil
}
