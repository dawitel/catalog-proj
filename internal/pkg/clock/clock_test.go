package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReal_Now(t *testing.T) {
	before := time.Now()
	got := Real{}.Now()
	after := time.Now()
	assert.False(t, got.IsZero())
	assert.True(t, !got.Before(before) || got.Equal(before), "Now should be >= before")
	assert.True(t, !got.After(after) || got.Equal(after), "Now should be <= after")
}
