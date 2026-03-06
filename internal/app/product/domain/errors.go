package domain

import "errors"

var (
	ErrProductNotFound        = errors.New("product not found")
	ErrProductNotActive       = errors.New("product is not active")
	ErrProductAlreadyActive   = errors.New("product is already active")
	ErrProductArchived        = errors.New("product is archived")
	ErrInvalidDiscountPeriod  = errors.New("discount period is invalid")
	ErrDiscountAlreadyActive  = errors.New("product already has an active discount")
	ErrNoActiveDiscount       = errors.New("product has no active discount")
	ErrInvalidProductName     = errors.New("product name is required")
	ErrInvalidCategory        = errors.New("product category is required")
	ErrInvalidPrice           = errors.New("product price must be positive")
	ErrInvalidDiscountPercent = errors.New("discount percentage must be between 0 and 100 exclusive")
)
