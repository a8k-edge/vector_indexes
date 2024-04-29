package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmokeHNSWIndex(t *testing.T) {
	f := NewHNSW(10, 20, 20, nil)
	f.AddBatch([][]float32{
		{0, 0},
		{1, 1},
		{2, 2},
	})
	distances, indexes := f.Search([]float32{0, 0}, 3)

	assert.True(t, len(distances) == 3, "Result not equal len 2")

	dExpected := []float32{0, 2, 8}
	for i := range distances {
		assert.InDelta(t, dExpected[i], distances[i], 0.01)
	}
	iExpected := []int{0, 1, 2}
	assert.ElementsMatch(t, iExpected, indexes, "Query result not valid")
}
