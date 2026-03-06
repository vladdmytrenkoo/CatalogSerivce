package repo

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/vladdmytrenkoo/committer"

	"CatalogService/internal/app/product/contracts"
	"CatalogService/internal/app/product/domain"
	"CatalogService/internal/models/m_outbox"
)

var _ contracts.OutboxRepository = (*OutboxRepo)(nil)

type OutboxRepo struct{}

func NewOutboxRepo() *OutboxRepo {
	return &OutboxRepo{}
}

func (r *OutboxRepo) InsertEventMut(event domain.DomainEvent, now time.Time) committer.Mutation {
	payload, _ := json.Marshal(event)

	data := &m_outbox.Data{
		EventID:     uuid.NewString(),
		EventType:   event.EventType(),
		AggregateID: event.AggregateID(),
		Payload:     string(payload),
		Status:      m_outbox.StatusPending,
		CreatedAt:   now,
	}

	return m_outbox.InsertMut(data)
}
