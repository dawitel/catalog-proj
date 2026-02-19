package product

import (
	"context"

	"github.com/dawitel/product-catalog-service/internal/app/product/contracts"
	productv1 "github.com/dawitel/product-catalog-service/proto/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) ListProducts(ctx context.Context, req *productv1.ListProductsRequest) (*productv1.ListProductsReply, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request required")
	}
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}
	filter := contracts.ListFilter{}
	if req.GetCategory() != "" {
		c := req.GetCategory()
		filter.Category = &c
	}
	page := contracts.ListPage{PageSize: pageSize, Token: req.GetPageToken()}
	result, err := h.List.Execute(ctx, filter, page)
	if err != nil {
		return nil, LogRPCError("ListProducts", err)
	}
	return listProductsReplyFromResult(result), nil
}
