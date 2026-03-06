package apply_discount

import (
	"context"
	"math/big"
	"time"

	"CatalogService/internal/app/product/contracts"
	"CatalogService/internal/app/product/domain"
	"CatalogService/internal/pkg/clock"

	committer "github.com/vladdmytrenkoo/commiter"
)

type ApplyRequest struct {
	ProductID  string
	Percentage *big.Rat
	StartDate  time.Time
	EndDate    time.Time
}

type ApplyInteractor struct {
	repo      contracts.ProductRepository
	outbox    contracts.OutboxRepository
	committer committer.Committer
	clock     clock.Clock
}

func NewApply(
	repo contracts.ProductRepository,
	outbox contracts.OutboxRepository,
	committer committer.Committer,
	clock clock.Clock,
) *ApplyInteractor {
	return &ApplyInteractor{repo: repo, outbox: outbox, committer: committer, clock: clock}
}

func (it *ApplyInteractor) Execute(ctx context.Context, req ApplyRequest) error {
	now := it.clock.Now()

	product, err := it.repo.GetByID(ctx, req.ProductID)
	if err != nil {
		return err
	}

	discount, err := domain.NewDiscount(req.Percentage, req.StartDate, req.EndDate)
	if err != nil {
		return err
	}

	if err := product.ApplyDiscount(discount, now); err != nil {
		return err
	}

	plan := committer.NewPlan()
	plan.Add(it.repo.UpdateMut(product))

	for _, event := range product.DomainEvents() {
		plan.Add(it.outbox.InsertEventMut(event, now))
	}

	return it.committer.Apply(ctx, plan)
}

type RemoveRequest struct {
	ProductID string
}

type RemoveInteractor struct {
	repo      contracts.ProductRepository
	outbox    contracts.OutboxRepository
	committer committer.Committer
	clock     clock.Clock
}

func NewRemove(
	repo contracts.ProductRepository,
	outbox contracts.OutboxRepository,
	committer committer.Committer,
	clock clock.Clock,
) *RemoveInteractor {
	return &RemoveInteractor{repo: repo, outbox: outbox, committer: committer, clock: clock}
}

func (it *RemoveInteractor) Execute(ctx context.Context, req RemoveRequest) error {
	now := it.clock.Now()

	product, err := it.repo.GetByID(ctx, req.ProductID)
	if err != nil {
		return err
	}

	if err := product.RemoveDiscount(now); err != nil {
		return err
	}

	plan := committer.NewPlan()
	plan.Add(it.repo.UpdateMut(product))

	for _, event := range product.DomainEvents() {
		plan.Add(it.outbox.InsertEventMut(event, now))
	}

	return it.committer.Apply(ctx, plan)
}
