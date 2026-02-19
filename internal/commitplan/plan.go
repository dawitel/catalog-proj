package commitplan

import (
	"context"

	"cloud.google.com/go/spanner"
)

type Plan struct {
	muts []*spanner.Mutation
}

func NewPlan() *Plan {
	return &Plan{muts: nil}
}

func (p *Plan) Add(m *spanner.Mutation) {
	if m == nil {
		return
	}
	p.muts = append(p.muts, m)
}

func (p *Plan) Mutations() []*spanner.Mutation {
	return p.muts
}

type Executor interface {
	Execute(ctx context.Context, plan *Plan) error
}
