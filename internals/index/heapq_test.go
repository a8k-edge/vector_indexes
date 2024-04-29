package index

import (
	"cmp"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinHeapQ(t *testing.T) {
	minHeapQ := NewMinHeapQ(5)

	minHeapQ.Push(1, 10.0)
	minHeapQ.Push(2, 5.0)
	minHeapQ.Push(3, 15.0)
	minHeapQ.Push(4, 3.0)

	assert.True(t, minHeapQ.Len() == 4)
	assert.True(t, minHeapQ.Top().Distance == 3.0)

	items := make([]float32, 0)
	for minHeapQ.Len() > 0 {
		items = append(items, minHeapQ.Pop().Distance)
	}

	assert.True(t, slices.IsSorted(items))
}

func TestMaxHeapQ(t *testing.T) {
	maxHeapQ := NewMaxHeapQ(5)

	maxHeapQ.Push(1, 10.0)
	maxHeapQ.Push(2, 5.0)
	maxHeapQ.Push(3, 15.0)
	maxHeapQ.Push(4, 3.0)

	assert.True(t, maxHeapQ.Len() == 4)
	assert.True(t, maxHeapQ.Top().Distance == 15.0)

	items := make([]float32, 0)
	for maxHeapQ.Len() > 0 {
		items = append(items, maxHeapQ.Pop().Distance)
	}

	isSorted := slices.IsSortedFunc(items, func(a, b float32) int {
		return cmp.Compare(a, b) * -1
	})
	assert.True(t, isSorted)
}
