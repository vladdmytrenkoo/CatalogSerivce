package list_products

import (
	"math/big"
	"time"
)

type ProductItemDTO struct {
	ID             string
	Name           string
	Category       string
	BasePrice      *big.Rat
	EffectivePrice *big.Rat
	Status         string
	CreatedAt      time.Time
}

type ResultDTO struct {
	Products      []*ProductItemDTO
	NextPageToken string
}
