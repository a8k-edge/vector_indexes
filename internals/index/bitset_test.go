package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitSet_Set(t *testing.T) {
	bs := NewBitSet(65)
	bs.Set(65)
	assert.True(t, bs.IsSet(65), "Bit 65 should be set")
	assert.False(t, bs.IsSet(64), "Bit 64 should not be set")
}

func TestBitSet_Reset(t *testing.T) {
	bs := NewBitSet(65)
	bs.Set(65)
	bs.Reset(65)
	assert.False(t, bs.IsSet(65), "Bit 65 should be reset")
}

func TestBitSet_IsSet(t *testing.T) {
	bs := NewBitSet(65)
	bs.Set(65)
	assert.True(t, bs.IsSet(65), "Bit 65 should be set")
	assert.False(t, bs.IsSet(64), "Bit 64 should not be set")
}
