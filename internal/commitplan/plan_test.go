package commitplan

import (
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/stretchr/testify/assert"
)

func TestNewPlan(t *testing.T) {
	p := NewPlan()
	assert.NotNil(t, p)
	muts := p.Mutations()
	assert.Nil(t, muts)
}

func TestPlan_Add(t *testing.T) {
	m := spanner.InsertMap("t", map[string]interface{}{"id": "x"})
	p := NewPlan()
	p.Add(m)
	muts := p.Mutations()
	assert.Len(t, muts, 1)
	p.Add(spanner.InsertMap("t", map[string]interface{}{"id": "y"}))
	muts = p.Mutations()
	assert.Len(t, muts, 2)
}

func TestPlan_AddNil(t *testing.T) {
	p := NewPlan()
	m := spanner.InsertMap("t", map[string]interface{}{"id": "x"})
	p.Add(m)
	p.Add(nil)
	muts := p.Mutations()
	assert.Len(t, muts, 1)
}
