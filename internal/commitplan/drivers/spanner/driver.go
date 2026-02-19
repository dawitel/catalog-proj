package spanner

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/dawitel/product-catalog-service/internal/commitplan"
)

type Executor struct {
	client *spanner.Client
}

func NewExecutor(client *spanner.Client) *Executor {
	return &Executor{client: client}
}

func (e *Executor) Execute(ctx context.Context, plan *commitplan.Plan) error {
	_, err := e.client.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		for _, cu := range plan.ConditionalUpdates() {
			stmt := spanner.Statement{SQL: cu.Stmt, Params: cu.Params}
			n, err := tx.Update(ctx, stmt)
			if err != nil {
				return err
			}
			if n != 1 {
				return commitplan.ErrConcurrentModification
			}
		}
		return tx.BufferWrite(plan.Mutations())
	})
	return err
}
