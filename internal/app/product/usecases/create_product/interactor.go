package create_product

import (
	"context"

	"CatalogService/internal/app/product/contracts"
	"CatalogService/internal/app/product/domain"
	"CatalogService/internal/pkg/clock"

	"github.com/google/uuid"
	committer "github.com/vladdmytrenkoo/commiter"
)

type Request struct {
	Name        string
	Description string
	Category    string
	PriceNum    int64
	PriceDenom  int64
}

type Interactor struct {
	repo      contracts.ProductRepository
	outbox    contracts.OutboxRepository
	committer committer.Committer
	clock     clock.Clock
}

func New(
	repo contracts.ProductRepository,
	outbox contracts.OutboxRepository,
	committer committer.Committer,
	clock clock.Clock,
) *Interactor {
	return &Interactor{repo: repo, outbox: outbox, committer: committer, clock: clock}
}

func (it *Interactor) Execute(ctx context.Context, req Request) (string, error) {
	now := it.clock.Now()

	basePrice, err := domain.NewMoney(req.PriceNum, req.PriceDenom)
	if err != nil {
		return "", err
	}

	product, err := domain.NewProduct(uuid.NewString(), req.Name, req.Description, req.Category, basePrice, now)
	if err != nil {
		return "", err
	}

	plan := committer.NewPlan()
	plan.Add(it.repo.InsertMut(product))

	for _, event := range product.DomainEvents() {
		plan.Add(it.outbox.InsertEventMut(event, now))
	}

	if err := it.committer.Apply(ctx, plan); err != nil {
		return "", err
	}

	return product.ID(), nil
}
