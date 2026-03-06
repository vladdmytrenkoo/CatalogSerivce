package product

import (
	"CatalogService/internal/app/product/domain"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapDomainErrorToGRPC(err error) error {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrProductNotActive),
		errors.Is(err, domain.ErrProductAlreadyActive),
		errors.Is(err, domain.ErrProductArchived),
		errors.Is(err, domain.ErrDiscountAlreadyActive),
		errors.Is(err, domain.ErrNoActiveDiscount),
		errors.Is(err, domain.ErrInvalidDiscountPeriod):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrInvalidProductName),
		errors.Is(err, domain.ErrInvalidCategory),
		errors.Is(err, domain.ErrInvalidPrice),
		errors.Is(err, domain.ErrInvalidDiscountPercent):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
