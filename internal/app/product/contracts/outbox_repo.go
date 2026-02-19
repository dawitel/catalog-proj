package contracts

import (
	"time"

	"cloud.google.com/go/spanner"
)

type OutboxRepo interface {
	InsertMut(eventID, eventType, aggregateID, payload, status string, createdAt time.Time) *spanner.Mutation
}
