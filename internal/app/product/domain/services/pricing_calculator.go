package services

import (
	"CatalogService/internal/app/product/domain"
	"math/big"
	"time"
)

type PricingCalculator struct{}

func NewPricingCalculator() *PricingCalculator {
	return &PricingCalculator{}
}

func (pc *PricingCalculator) EffectivePrice(product *domain.Product, now time.Time) *big.Rat {
	base := product.BasePrice().Amount()

	discount := product.Discount()
	if discount == nil || !discount.IsValidAt(now) {
		return base
	}

	discountAmount := new(big.Rat).Mul(base, discount.FractionOff())

	return new(big.Rat).Sub(base, discountAmount)
}
