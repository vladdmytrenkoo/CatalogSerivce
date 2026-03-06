package m_outbox

import (
	"time"

	"cloud.google.com/go/spanner"
)

const (
	StatusPending   = "pending"
	StatusProcessed = "processed"
)

type Data struct {
	EventID     string
	EventType   string
	AggregateID string
	Payload     string
	Status      string
	CreatedAt   time.Time
	ProcessedAt *time.Time
}

func InsertMut(d *Data) *spanner.Mutation {
	return spanner.InsertMap(Table, map[string]interface{}{
		EventID:     d.EventID,
		EventType:   d.EventType,
		AggregateID: d.AggregateID,
		Payload:     d.Payload,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		ProcessedAt: d.ProcessedAt,
	})
}
