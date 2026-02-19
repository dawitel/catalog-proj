package domain

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMoney(t *testing.T) {
	m := NewMoney(100, 1)
	require.NotNil(t, m)
	assert.Equal(t, int64(100), m.Numerator())
	assert.Equal(t, int64(1), m.Denominator())

	m = NewMoney(1999, 100)
	require.NotNil(t, m)
	assert.Equal(t, int64(1999), m.Numerator())
	assert.Equal(t, int64(100), m.Denominator())

	assert.Nil(t, NewMoney(1, 0))
}

func TestNewMoneyFromRat(t *testing.T) {
	r := big.NewRat(50, 1)
	m := NewMoneyFromRat(r)
	require.NotNil(t, m)
	assert.Equal(t, int64(50), m.Numerator())
	assert.Equal(t, int64(1), m.Denominator())

	assert.Nil(t, NewMoneyFromRat(nil))
}

func TestMoney_NilReceiver(t *testing.T) {
	var m *Money
	assert.Equal(t, int64(0), m.Numerator())
	assert.Equal(t, int64(0), m.Denominator())
	assert.NotNil(t, m.Rat())
}

func TestMoney_Rat(t *testing.T) {
	m := NewMoney(3, 4)
	require.NotNil(t, m)
	r := m.Rat()
	require.NotNil(t, r)
	assert.True(t, r.Cmp(big.NewRat(3, 4)) == 0)
	r.SetInt64(99)
	assert.Equal(t, int64(3), m.Numerator())
}
