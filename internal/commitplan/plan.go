package commitplan

import (
	"context"
	"errors"

	"cloud.google.com/go/spanner"
)

var ErrConcurrentModification = errors.New("concurrent modification")

type ConditionalUpdate struct {
	Stmt   string
	Params map[string]interface{}
}

type Plan struct {
	muts             []*spanner.Mutation
	conditionalUpdates []ConditionalUpdate
}

func NewPlan() *Plan {
	return &Plan{}
}

func (p *Plan) Add(m *spanner.Mutation) {
	if m == nil {
		return
	}
	p.muts = append(p.muts, m)
}

func (p *Plan) AddConditionalUpdate(stmt string, params map[string]interface{}) {
	if stmt == "" {
		return
	}
	p.conditionalUpdates = append(p.conditionalUpdates, ConditionalUpdate{Stmt: stmt, Params: params})
}

func (p *Plan) Mutations() []*spanner.Mutation {
	return p.muts
}

func (p *Plan) ConditionalUpdates() []ConditionalUpdate {
	return p.conditionalUpdates
}

type Executor interface {
	Execute(ctx context.Context, plan *Plan) error
}

type Applier interface {
	Apply(ctx context.Context, plan *Plan) error
}
