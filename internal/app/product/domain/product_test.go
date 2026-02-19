package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProduct(t *testing.T) {
	now := time.Date(2025, 2, 1, 12, 0, 0, 0, time.UTC)
	base := NewMoney(9999, 100)
	p := NewProduct("id1", "name", "desc", "cat", base, now)
	require.NotNil(t, p)
	assert.Equal(t, "id1", p.ID())
	assert.Equal(t, "name", p.Name())
	assert.Equal(t, "desc", p.Description())
	assert.Equal(t, "cat", p.Category())
	assert.Equal(t, ProductStatusActive, p.Status())
	assert.Equal(t, int64(1), p.Version())
	assert.True(t, p.CreatedAt().Equal(now))
	assert.True(t, p.UpdatedAt().Equal(now))
	assert.Nil(t, p.Discount())
	assert.NotNil(t, p.Changes())
	assert.Len(t, p.DomainEvents(), 1)

	assert.Nil(t, NewProduct("id", "n", "d", "c", nil, now))
}

func TestRestoreProduct(t *testing.T) {
	base := NewMoney(100, 1)
	created := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	p := RestoreProduct("id2", "n", "d", "c", base, nil, ProductStatusInactive, 1, created, updated, time.Time{})
	require.NotNil(t, p)
	assert.Equal(t, "id2", p.ID())
	assert.Equal(t, ProductStatusInactive, p.Status())
	assert.Equal(t, int64(1), p.Version())
	assert.True(t, p.UpdatedAt().Equal(updated))
	assert.Nil(t, p.DomainEvents())

	assert.Nil(t, RestoreProduct("id", "n", "d", "c", nil, nil, ProductStatusActive, 1, created, updated, time.Time{}))
}

func TestProduct_UpdateDetails(t *testing.T) {
	now := time.Now()
	base := NewMoney(100, 1)
	p := NewProduct("id", "old", "oldd", "oldc", base, now)
	require.NotNil(t, p)
	err := p.UpdateDetails("new", "newd", "newc", now.Add(time.Hour))
	require.NoError(t, err)
	assert.Equal(t, "new", p.Name())
	assert.Equal(t, "newd", p.Description())
	assert.Equal(t, "newc", p.Category())
	assert.True(t, p.Changes().Dirty(FieldName))

	p2 := RestoreProduct("a", "n", "d", "c", base, nil, ProductStatusArchived, 1, now, now, now)
	err = p2.UpdateDetails("x", "y", "z", now)
	assert.ErrorIs(t, err, ErrProductArchived)
}

func TestProduct_Activate(t *testing.T) {
	now := time.Now()
	base := NewMoney(100, 1)
	p := RestoreProduct("id", "n", "d", "c", base, nil, ProductStatusInactive, 1, now, now, time.Time{})
	err := p.Activate(now.Add(time.Hour))
	require.NoError(t, err)
	assert.Equal(t, ProductStatusActive, p.Status())

	err = p.Activate(now)
	require.NoError(t, err)
	assert.Equal(t, ProductStatusActive, p.Status())

	p2 := RestoreProduct("a", "n", "d", "c", base, nil, ProductStatusArchived, 1, now, now, now)
	err = p2.Activate(now)
	assert.ErrorIs(t, err, ErrProductArchived)
}

func TestProduct_Deactivate(t *testing.T) {
	now := time.Now()
	base := NewMoney(100, 1)
	p := NewProduct("id", "n", "d", "c", base, now)
	err := p.Deactivate(now.Add(time.Hour))
	require.NoError(t, err)
	assert.Equal(t, ProductStatusInactive, p.Status())

	err = p.Deactivate(now)
	require.NoError(t, err)
	assert.Equal(t, ProductStatusInactive, p.Status())

	p2 := RestoreProduct("a", "n", "d", "c", base, nil, ProductStatusArchived, 1, now, now, now)
	err = p2.Deactivate(now)
	assert.ErrorIs(t, err, ErrProductArchived)
}

func TestProduct_Archive(t *testing.T) {
	now := time.Now()
	base := NewMoney(100, 1)
	p := NewProduct("id", "n", "d", "c", base, now)
	err := p.Archive(now.Add(time.Hour))
	require.NoError(t, err)
	assert.Equal(t, ProductStatusArchived, p.Status())
	assert.False(t, p.ArchivedAt().IsZero())

	err = p.Archive(now)
	require.NoError(t, err)
}

func TestProduct_ApplyDiscount(t *testing.T) {
	now := time.Now()
	base := NewMoney(100, 1)
	p := NewProduct("id", "n", "d", "c", base, now)
	start := now.Add(-time.Hour)
	end := now.Add(time.Hour)
	disc := NewDiscount(25, start, end)
	require.NotNil(t, disc)
	err := p.ApplyDiscount(disc, now)
	require.NoError(t, err)
	assert.NotNil(t, p.Discount())
	assert.Equal(t, int64(25), p.Discount().Percentage())

	p2 := RestoreProduct("a", "n", "d", "c", base, nil, ProductStatusInactive, 1, now, now, time.Time{})
	err = p2.ApplyDiscount(disc, now)
	assert.ErrorIs(t, err, ErrProductNotActive)

	expired := NewDiscount(10, now.Add(-24*time.Hour), now.Add(-1*time.Hour))
	err = p.ApplyDiscount(expired, now)
	assert.ErrorIs(t, err, ErrInvalidDiscountPeriod)

	err = p.ApplyDiscount(nil, now)
	assert.ErrorIs(t, err, ErrInvalidDiscountPeriod)
}

func TestProduct_RemoveDiscount(t *testing.T) {
	now := time.Now()
	base := NewMoney(100, 1)
	disc := NewDiscount(10, now.Add(-time.Hour), now.Add(time.Hour))
	p := NewProduct("id", "n", "d", "c", base, now)
	_ = p.ApplyDiscount(disc, now)
	err := p.RemoveDiscount(now.Add(time.Hour))
	require.NoError(t, err)
	assert.Nil(t, p.Discount())

	p2 := NewProduct("id2", "n", "d", "c", base, now)
	err = p2.RemoveDiscount(now)
	require.NoError(t, err)

	p3 := RestoreProduct("a", "n", "d", "c", base, nil, ProductStatusArchived, 1, now, now, now)
	err = p3.RemoveDiscount(now)
	assert.ErrorIs(t, err, ErrProductArchived)
}

func TestProduct_NilReceiver(t *testing.T) {
	var p *Product
	assert.Empty(t, p.ID())
	assert.Empty(t, p.Name())
	assert.Equal(t, int64(0), p.Version())
	assert.Nil(t, p.BasePrice())
	assert.Nil(t, p.Changes())
	assert.Nil(t, p.DomainEvents())
	assert.ErrorIs(t, p.UpdateDetails("", "", "", time.Time{}), ErrInvalidProduct)
	assert.ErrorIs(t, p.Activate(time.Time{}), ErrInvalidProduct)
	assert.ErrorIs(t, p.Deactivate(time.Time{}), ErrInvalidProduct)
	assert.ErrorIs(t, p.Archive(time.Time{}), ErrInvalidProduct)
	assert.ErrorIs(t, p.ApplyDiscount(nil, time.Time{}), ErrInvalidProduct)
	assert.ErrorIs(t, p.RemoveDiscount(time.Time{}), ErrInvalidProduct)
}
