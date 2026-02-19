package services

import (
	"testing"
	"time"

	"github.com/dawitel/product-catalog-service/internal/app/product/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEffectivePrice_NoDiscount(t *testing.T) {
	base := domain.NewMoney(10000, 100)
	at := time.Now()
	effective := EffectivePrice(base, nil, at)
	require.NotNil(t, effective)
	assert.True(t, effective.Rat().Cmp(base.Rat()) == 0)
}

func TestEffectivePrice_InvalidDiscountPeriod(t *testing.T) {
	base := domain.NewMoney(10000, 100)
	start := time.Now().Add(-48 * time.Hour)
	end := time.Now().Add(-24 * time.Hour)
	disc := domain.NewDiscount(20, start, end)
	at := time.Now()
	effective := EffectivePrice(base, disc, at)
	require.NotNil(t, effective)
	assert.True(t, effective.Rat().Cmp(base.Rat()) == 0)
}

func TestEffectivePrice_WithDiscount(t *testing.T) {
	base := domain.NewMoney(10000, 100)
	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)
	disc := domain.NewDiscount(20, start, end)
	at := time.Now()
	effective := EffectivePrice(base, disc, at)
	require.NotNil(t, effective)
	expected := domain.NewMoney(8000, 100)
	require.NotNil(t, expected)
	assert.True(t, effective.Rat().Cmp(expected.Rat()) == 0)
}

func TestEffectivePrice_50Percent(t *testing.T) {
	base := domain.NewMoney(100, 1)
	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)
	disc := domain.NewDiscount(50, start, end)
	at := time.Now()
	effective := EffectivePrice(base, disc, at)
	require.NotNil(t, effective)
	assert.Equal(t, int64(50), effective.Numerator())
	assert.Equal(t, int64(1), effective.Denominator())
}

func TestEffectivePrice_NilBase(t *testing.T) {
	at := time.Now()
	effective := EffectivePrice(nil, nil, at)
	assert.Nil(t, effective)
}
