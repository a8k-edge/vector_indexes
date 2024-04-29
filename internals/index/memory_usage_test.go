package index

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryUsage(t *testing.T) {
	baseData, truth, learn, queries := loadSift()
	_ = learn
	_ = truth
	_ = queries

	f := NewFlat(128, nil)
	f.AddBatch(baseData)
	fSize := f.Size()

	fmt.Println("Flat Size", fSize)

	// fpq := NewFlatPQ(128, nil)
	// fpq.Train(learn, NewPQ(128, 8))
	// fpq.AddBatch(baseData)
	// fpqSize := fpq.Size()
	// fmt.Println("FlatPQ Size", fpqSize)
	// fmt.Println("FlatPQ PQ Size", fpq.SizePQ())

	hnsw := NewHNSW(10, 20, 20, nil)
	hnsw.AddBatch(baseData)
	fmt.Println("HNSW Size", hnsw.Size())

	assert.True(t, false)
}
