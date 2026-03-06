package contracts

import (
	"CatalogService/internal/app/product/domain"
	"github.com/vladdmytrenkoo/committer"
	"time"
)

type OutboxRepository interface {
	InsertEventMut(event domain.DomainEvent, now time.Time) committer.Mutation
}
