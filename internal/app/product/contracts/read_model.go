package contracts

import (
	"context"
	"math/big"
	"time"
)

type ProductReadModel interface {
	GetByID(ctx context.Context, id string) (*ProductView, error)
	ListActive(ctx context.Context, filter ListFilter) (*ProductListResult, error)
}

type ProductView struct {
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

type ListFilter struct {
	Category  string
	PageSize  int32
	PageToken string
}

type ProductListResult struct {
	Products      []*ProductView
	NextPageToken string
}
