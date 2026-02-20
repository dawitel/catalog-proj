package product

import (
	productv1 "github.com/dawitel/catalog-proj/proto/product/v1"
)

type Handler struct {
	productv1.UnimplementedProductServiceServer
	Create     CreateRunner
	Update     UpdateRunner
	Activate   ActivateRunner
	Deactivate DeactivateRunner
	Archive    ArchiveRunner
	ApplyDisc  ApplyDiscountRunner
	RemoveDisc RemoveDiscountRunner
	Get        GetProductRunner
	List       ListProductsRunner
}

func NewHandler(
	create CreateRunner,
	update UpdateRunner,
	activate ActivateRunner,
	deactivate DeactivateRunner,
	archive ArchiveRunner,
	applyDiscount ApplyDiscountRunner,
	removeDiscount RemoveDiscountRunner,
	getProduct GetProductRunner,
	listProducts ListProductsRunner,
) *Handler {
	return &Handler{
		Create:     create,
		Update:     update,
		Activate:   activate,
		Deactivate: deactivate,
		Archive:    archive,
		ApplyDisc:  applyDiscount,
		RemoveDisc: removeDiscount,
		Get:        getProduct,
		List:       listProducts,
	}
}
