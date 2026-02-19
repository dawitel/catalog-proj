package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChangeTracker(t *testing.T) {
	c := NewChangeTracker()
	assert.NotNil(t, c)
	assert.False(t, c.Dirty(FieldName))
}

func TestChangeTracker_MarkDirty(t *testing.T) {
	c := NewChangeTracker()
	c.MarkDirty(FieldName)
	assert.True(t, c.Dirty(FieldName))
	assert.False(t, c.Dirty(FieldDescription))
	c.MarkDirty(FieldStatus)
	assert.True(t, c.Dirty(FieldStatus))
}

func TestChangeTracker_NilReceiver(t *testing.T) {
	var c *ChangeTracker
	c.MarkDirty(FieldName)
	assert.False(t, c.Dirty(FieldName))
}
