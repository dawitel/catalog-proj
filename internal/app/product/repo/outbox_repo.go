package repo

import (
	"time"

	"cloud.google.com/go/spanner"
	"github.com/dawitel/catalog-proj/internal/models/m_outbox"
)

type OutboxRepo struct{}

func NewOutboxRepo() *OutboxRepo {
	return &OutboxRepo{}
}

func (r *OutboxRepo) InsertMut(eventID, eventType, aggregateID, payload, status string, createdAt time.Time) *spanner.Mutation {
	row := &m_outbox.Row{
		EventID:     eventID,
		EventType:   eventType,
		AggregateID: aggregateID,
		Payload:     payload,
		Status:      status,
		CreatedAt:   createdAt,
	}
	return m_outbox.InsertMut(row)
}
