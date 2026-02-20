package product

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	productv1 "github.com/dawitel/catalog-proj/proto/product/v1"
)

func (h *Handler) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.CreateProductReply, error) {
	if req == nil || req.GetName() == "" || req.GetCategory() == "" || req.GetBasePriceDenominator() == 0 || req.GetBasePriceNumerator() < 0 {
		return nil, status.Error(codes.InvalidArgument, "name, category and valid base price required")
	}
	appReq := createRequestFromProto(req)
	productID, err := h.Create.Execute(ctx, appReq)
	if err != nil {
		return nil, LogRPCError("CreateProduct", err)
	}
	return &productv1.CreateProductReply{ProductId: productID}, nil
}
