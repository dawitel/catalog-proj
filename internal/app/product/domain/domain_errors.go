package domain

import "errors"

var (
	ErrProductNotFound       = errors.New("product not found")
	ErrProductNotActive      = errors.New("product not active")
	ErrProductArchived       = errors.New("product is archived")
	ErrInvalidDiscountPeriod = errors.New("invalid discount period")
	ErrInvalidProduct        = errors.New("invalid product")
)
