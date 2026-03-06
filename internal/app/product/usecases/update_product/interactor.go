package update_product

import (
	"context"

	"CatalogService/internal/app/product/contracts"
	"CatalogService/internal/pkg/clock"

	committer "github.com/vladdmytrenkoo/commiter"
)

type Request struct {
	ProductID   string
	Name        string
	Description string
	Category    string
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

func (it *Interactor) Execute(ctx context.Context, req Request) error {
	now := it.clock.Now()

	product, err := it.repo.GetByID(ctx, req.ProductID)
	if err != nil {
		return err
	}

	if err := product.Update(req.Name, req.Description, req.Category, now); err != nil {
		return err
	}

	plan := committer.NewPlan()

	if mut := it.repo.UpdateMut(product); mut != nil {
		plan.Add(mut)
	}

	for _, event := range product.DomainEvents() {
		plan.Add(it.outbox.InsertEventMut(event, now))
	}

	if plan.IsEmpty() {
		return nil
	}

	return it.committer.Apply(ctx, plan)
}
