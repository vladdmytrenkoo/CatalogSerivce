package contracts

import (
	"CatalogService/internal/app/product/domain"
	"context"

	committer "github.com/vladdmytrenkoo/commiter"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	InsertMut(product *domain.Product) committer.Mutation
	UpdateMut(product *domain.Product) committer.Mutation
}
