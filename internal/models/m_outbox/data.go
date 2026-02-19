package m_outbox

import (
	"time"

	"cloud.google.com/go/spanner"
)

type Row struct {
	EventID     string
	EventType   string
	AggregateID string
	Payload     string
	Status      string
	CreatedAt   time.Time
	ProcessedAt time.Time
}

func InsertMut(r *Row) *spanner.Mutation {
	return spanner.Insert(TableName, []string{
		EventID, EventType, AggregateID, Payload, Status, CreatedAt,
	}, []interface{}{
		r.EventID, r.EventType, r.AggregateID, r.Payload, r.Status, r.CreatedAt,
	})
}
