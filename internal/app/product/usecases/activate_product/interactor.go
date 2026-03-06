package activate_product

import (
	"context"

	"CatalogService/internal/app/product/contracts"
	"CatalogService/internal/pkg/clock"

	"github.com/vladdmytrenkoo/commiter"
)

type Request struct {
	ProductID string
}

type ActivateInteractor struct {
	repo      contracts.ProductRepository
	outbox    contracts.OutboxRepository
	committer committer.Committer
	clock     clock.Clock
}

func NewActivate(
	repo contracts.ProductRepository,
	outbox contracts.OutboxRepository,
	committer committer.Committer,
	clock clock.Clock,
) *ActivateInteractor {
	return &ActivateInteractor{repo: repo, outbox: outbox, committer: committer, clock: clock}
}

func (it *ActivateInteractor) Execute(ctx context.Context, req Request) error {
	now := it.clock.Now()

	product, err := it.repo.GetByID(ctx, req.ProductID)
	if err != nil {
		return err
	}

	if err := product.Activate(now); err != nil {
		return err
	}

	plan := committer.NewPlan()
	plan.Add(it.repo.UpdateMut(product))

	for _, event := range product.DomainEvents() {
		plan.Add(it.outbox.InsertEventMut(event, now))
	}

	return it.committer.Apply(ctx, plan)
}

type DeactivateInteractor struct {
	repo      contracts.ProductRepository
	outbox    contracts.OutboxRepository
	committer committer.Committer
	clock     clock.Clock
}

func NewDeactivate(
	repo contracts.ProductRepository,
	outbox contracts.OutboxRepository,
	committer committer.Committer,
	clock clock.Clock,
) *DeactivateInteractor {
	return &DeactivateInteractor{repo: repo, outbox: outbox, committer: committer, clock: clock}
}

func (it *DeactivateInteractor) Execute(ctx context.Context, req Request) error {
	now := it.clock.Now()

	product, err := it.repo.GetByID(ctx, req.ProductID)
	if err != nil {
		return err
	}

	if err := product.Deactivate(now); err != nil {
		return err
	}

	plan := committer.NewPlan()
	plan.Add(it.repo.UpdateMut(product))

	for _, event := range product.DomainEvents() {
		plan.Add(it.outbox.InsertEventMut(event, now))
	}

	return it.committer.Apply(ctx, plan)
}

type ArchiveInteractor struct {
	repo      contracts.ProductRepository
	outbox    contracts.OutboxRepository
	committer committer.Committer
	clock     clock.Clock
}

func NewArchive(
	repo contracts.ProductRepository,
	outbox contracts.OutboxRepository,
	committer committer.Committer,
	clock clock.Clock,
) *ArchiveInteractor {
	return &ArchiveInteractor{repo: repo, outbox: outbox, committer: committer, clock: clock}
}

func (it *ArchiveInteractor) Execute(ctx context.Context, req Request) error {
	now := it.clock.Now()

	product, err := it.repo.GetByID(ctx, req.ProductID)
	if err != nil {
		return err
	}

	if err := product.Archive(now); err != nil {
		return err
	}

	plan := committer.NewPlan()
	plan.Add(it.repo.UpdateMut(product))

	for _, event := range product.DomainEvents() {
		plan.Add(it.outbox.InsertEventMut(event, now))
	}

	return it.committer.Apply(ctx, plan)
}
