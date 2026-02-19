package product

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	archive_product "github.com/dawitel/product-catalog-service/internal/app/product/usecases/archive_product"
	productv1 "github.com/dawitel/product-catalog-service/proto/product/v1"
)

func (h *Handler) ArchiveProduct(ctx context.Context, req *productv1.ArchiveProductRequest) (*productv1.ArchiveProductReply, error) {
	if req == nil || req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id required")
	}
	if err := h.Archive.Execute(ctx, archive_product.Request{ProductID: req.GetProductId()}); err != nil {
		return nil, LogRPCError("ArchiveProduct", err)
	}
	return &productv1.ArchiveProductReply{}, nil
}
