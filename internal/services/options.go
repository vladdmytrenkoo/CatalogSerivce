package services

import (
	"CatalogService/internal/app/product/queries/get_product"
	"CatalogService/internal/app/product/queries/list_products"
	"CatalogService/internal/app/product/repo"
	"CatalogService/internal/app/product/usecases/activate_product"
	"CatalogService/internal/app/product/usecases/apply_discount"
	"CatalogService/internal/app/product/usecases/create_product"
	"CatalogService/internal/app/product/usecases/update_product"
	"CatalogService/internal/pkg/clock"
	"CatalogService/internal/transport/grpc/product"

	"cloud.google.com/go/spanner"
	spanner_committer "github.com/vladdmytrenkoo/committer/spanner"
)

type Options struct {
	Commands product.Commands
	Queries  product.Queries
}

func NewOptions(client *spanner.Client) *Options {
	realClock := clock.RealClock{}
	productRepo := repo.NewProductRepo(client)
	productReadRepo := repo.NewProductReadRepo(client)
	outboxRepo := repo.NewOutboxRepo()
	committer := spanner_committer.New(client)

	commands := product.Commands{
		CreateProduct:  create_product.New(productRepo, outboxRepo, committer, realClock),
		UpdateProduct:  update_product.New(productRepo, outboxRepo, committer, realClock),
		Activate:       activate_product.NewActivate(productRepo, outboxRepo, committer, realClock),
		Deactivate:     activate_product.NewDeactivate(productRepo, outboxRepo, committer, realClock),
		ApplyDiscount:  apply_discount.NewApply(productRepo, outboxRepo, committer, realClock),
		RemoveDiscount: apply_discount.NewRemove(productRepo, outboxRepo, committer, realClock),
	}

	queries := product.Queries{
		GetProduct:   get_product.New(productReadRepo),
		ListProducts: list_products.New(productReadRepo),
	}

	return &Options{
		Commands: commands,
		Queries:  queries,
	}
}
