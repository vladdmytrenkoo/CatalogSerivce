package get_product

import (
	"math/big"
	"time"
)

type ProductDTO struct {
	ID              string
	Name            string
	Description     string
	Category        string
	BasePrice       *big.Rat
	EffectivePrice  *big.Rat
	DiscountPercent *big.Rat
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
