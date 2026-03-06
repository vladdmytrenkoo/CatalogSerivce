package list_products

import (
	"CatalogService/internal/app/product/contracts"
	"context"
)

type Request struct {
	Category  string
	PageSize  int32
	PageToken string
}

type Query struct {
	readModel contracts.ProductReadModel
}

func New(readModel contracts.ProductReadModel) *Query {
	return &Query{readModel: readModel}
}

func (q *Query) Execute(ctx context.Context, req Request) (*ResultDTO, error) {
	result, err := q.readModel.ListActive(ctx, contracts.ListFilter{
		Category:  req.Category,
		PageSize:  req.PageSize,
		PageToken: req.PageToken,
	})
	if err != nil {
		return nil, err
	}

	items := make([]*ProductItemDTO, 0, len(result.Products))
	for _, v := range result.Products {
		items = append(items, &ProductItemDTO{
			ID:             v.ID,
			Name:           v.Name,
			Category:       v.Category,
			BasePrice:      v.BasePrice,
			EffectivePrice: v.EffectivePrice,
			Status:         v.Status,
			CreatedAt:      v.CreatedAt,
		})
	}

	return &ResultDTO{
		Products:      items,
		NextPageToken: result.NextPageToken,
	}, nil
}
