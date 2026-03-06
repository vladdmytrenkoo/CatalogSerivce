package get_product

import (
	"CatalogService/internal/app/product/contracts"
	"context"
)

type Query struct {
	readModel contracts.ProductReadModel
}

func New(readModel contracts.ProductReadModel) *Query {
	return &Query{readModel: readModel}
}

func (q *Query) Execute(ctx context.Context, productID string) (*ProductDTO, error) {
	view, err := q.readModel.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return &ProductDTO{
		ID:              view.ID,
		Name:            view.Name,
		Description:     view.Description,
		Category:        view.Category,
		BasePrice:       view.BasePrice,
		EffectivePrice:  view.EffectivePrice,
		DiscountPercent: view.DiscountPercent,
		Status:          view.Status,
		CreatedAt:       view.CreatedAt,
		UpdatedAt:       view.UpdatedAt,
	}, nil
}
