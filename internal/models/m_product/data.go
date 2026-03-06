package m_product

import (
	"math/big"
	"time"

	"cloud.google.com/go/spanner"
)

type Data struct {
	ProductID            string
	Name                 string
	Description          *string
	Category             string
	BasePriceNumerator   int64
	BasePriceDenominator int64
	DiscountPercent      *big.Rat
	DiscountStartDate    *time.Time
	DiscountEndDate      *time.Time
	Status               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	ArchivedAt           *time.Time
}

func InsertMut(d *Data) *spanner.Mutation {
	return spanner.InsertMap(Table, map[string]interface{}{
		ProductID:            d.ProductID,
		Name:                 d.Name,
		Description:          d.Description,
		Category:             d.Category,
		BasePriceNumerator:   d.BasePriceNumerator,
		BasePriceDenominator: d.BasePriceDenominator,
		DiscountPercent:      d.DiscountPercent,
		DiscountStartDate:    d.DiscountStartDate,
		DiscountEndDate:      d.DiscountEndDate,
		Status:               d.Status,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
		ArchivedAt:           d.ArchivedAt,
	})
}

func UpdateMut(productID string, cols map[string]interface{}) *spanner.Mutation {
	cols[ProductID] = productID
	return spanner.UpdateMap(Table, cols)
}
