package committer

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/dawitel/catalog-proj/internal/commitplan"
	spannerdriver "github.com/dawitel/catalog-proj/internal/commitplan/drivers/spanner"
)

type Committer struct {
	exec commitplan.Executor
}

func New(client *spanner.Client) *Committer {
	return &Committer{exec: spannerdriver.NewExecutor(client)}
}

func NewWithExecutor(exec commitplan.Executor) *Committer {
	return &Committer{exec: exec}
}

func (c *Committer) Apply(ctx context.Context, plan *commitplan.Plan) error {
	return c.exec.Execute(ctx, plan)
}
