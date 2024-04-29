package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmokeVamanaIndex(t *testing.T) {
	vmn := NewVamana(10, 100, 1.2, nil)
	vmn.AddBatch([][]float32{
		{0, 0},
		{1, 1},
		{2, 2},
	})

	distances, indexes := vmn.Search([]float32{0, 0}, 2)
	assert.True(t, len(distances) == 2, "Result not equal len 2")

	dExpected := []float32{0, 2}
	for i := range distances {
		assert.InDelta(t, dExpected[i], distances[i], 0.01)
	}
	iExpected := []int{0, 1}
	assert.ElementsMatch(t, iExpected, indexes, "Query result not valid")
}
