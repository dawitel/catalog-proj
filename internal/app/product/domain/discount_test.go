package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDiscount(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	d := NewDiscount(10, start, end)
	require.NotNil(t, d)
	assert.Equal(t, int64(10), d.Percentage())
	assert.True(t, d.StartDate().Equal(start))
	assert.True(t, d.EndDate().Equal(end))
	assert.Nil(t, NewDiscount(-1, start, end))
	assert.Nil(t, NewDiscount(101, start, end))
	assert.Nil(t, NewDiscount(50, end, start))
}

func TestDiscount_IsValidAt(t *testing.T) {
	start := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC)
	d := NewDiscount(20, start, end)
	require.NotNil(t, d)
	assert.False(t, d.IsValidAt(time.Date(2025, 5, 31, 23, 59, 59, 0, time.UTC)))
	assert.True(t, d.IsValidAt(start))
	assert.True(t, d.IsValidAt(time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)))
	assert.True(t, d.IsValidAt(end))
	assert.False(t, d.IsValidAt(time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)))
}

func TestDiscount_NilReceiver(t *testing.T) {
	var d *Discount
	assert.Equal(t, int64(0), d.Percentage())
	assert.True(t, d.StartDate().IsZero())
	assert.True(t, d.EndDate().IsZero())
	assert.False(t, d.IsValidAt(time.Now()))
}
