package index

import (
	"fmt"
	"testing"
	"time"

	"mvdb/internals/utils"

	"github.com/stretchr/testify/assert"
)

func TestFlatRecall(t *testing.T) {
	baseData, truth, _, queries := loadSift()
	k := 100

	f := NewFlat(128, nil)
	f.AddBatch(baseData)

	_, indexes := f.SearchMany(queries, k)
	expect, got := 0, 0
	for i := range queries {
		expect += k
		got += utils.IntersectionCount(indexes[i], truth[i][:k])
	}

	// got, expect := 0, 0
	// for qi := range queries {
	// 	_, indexes := f.Search(queries[qi], k)
	// 	expect += k
	// 	got += utils.IntersectionCount(indexes, truth[qi][:k])
	// }

	recall := float64(got) / float64(expect)
	fmt.Println("Recall:", recall)
	assert.True(t, recall >= 0.99)
}

func TestFlatPQRecall(t *testing.T) {
	baseData, truth, learn, queries := loadSift()
	_ = learn

	f := NewFlatPQ(128, nil)
	f.Train(learn, NewPQ(128, 8))
	f.AddBatch(baseData)

	got, expect := 0, 0
	k := 100
	for qi := range queries {
		_, indexes := f.Search(queries[qi], k)
		expect += k
		got += utils.IntersectionCount(indexes, truth[qi][:k])
	}

	recall := float64(got) / float64(expect)
	fmt.Println("Recall:", recall)
	assert.True(t, recall >= 0.8)
}

func TestIVFRecall(t *testing.T) {
	baseData, truth, learn, queries := loadSift()

	ivf := NewIVF(128, 100, 10, nil)
	ivf.Train(learn)
	ivf.AddBatch(baseData)

	k := 100
	_, indexes := ivf.SearchMany(queries, k)
	expect, got := 0, 0
	for i := range queries {
		expect += k
		got += utils.IntersectionCount(indexes[i], truth[i][:k])
	}

	// for qi := range queries {
	// 	_, indexes := ivf.Search(queries[qi], k)
	// 	expect += k
	// 	got += utils.IntersectionCount(indexes, truth[qi][:k])
	// }

	recall := float64(got) / float64(expect)
	fmt.Println("Recall:", recall)
	assert.True(t, recall >= 0.99)
}

func TestIVFPQRecall(t *testing.T) {
	baseData, truth, learn, queries := loadSift()

	ivf := NewIVFPQ(128, 100, 10, nil)
	ivf.Train(learn, NewPQ(128, 8))
	ivf.AddBatch(baseData)

	got, expect := 0, 0
	k := 100
	for qi := range queries {
		_, indexes := ivf.Search(queries[qi], k)
		expect += k
		got += utils.IntersectionCount(indexes, truth[qi][:k])
	}

	recall := float64(got) / float64(expect)
	fmt.Println("Recall:", recall)
	assert.True(t, recall >= 0.99)
}

func TestHNSWRecall(t *testing.T) {
	baseData, truth, _, queries := loadSift()
	_ = truth
	_ = queries

	hnsw := NewHNSW(32, 150, 40, nil)
	hnsw.AddBatch(baseData)

	hnsw.PrintLvlDistribution()

	got, expect := 0, 0
	k := 100
	start := time.Now()
	for qi := range queries {
		_, indexes := hnsw.Search(queries[qi], k)
		expect += k
		got += utils.IntersectionCount(indexes, truth[qi][:k])
	}
	elapsed := time.Since(start)
	fmt.Printf("Search took %s\n", elapsed)

	recall := float64(got) / float64(expect)
	fmt.Println("Recall:", recall)
	assert.True(t, recall >= 0.96)
}

func TestVamanaRecall(t *testing.T) {
	baseData, truth, _, queries := loadSift()
	_ = truth
	_ = queries

	start := time.Now()
	vmn := NewVamana(32, 150, 1.2, nil)
	vmn.AddBatch(baseData[:100000])
	elapsed := time.Since(start)
	fmt.Printf("Build took %s\n", elapsed)
	vmn.Print()

	got, expect := 0, 0
	k := 100
	start = time.Now()
	for qi := range queries {
		if qi%1000 == 0 {
			fmt.Println(time.Now(), "Qu", qi)
		}
		_, indexes := vmn.Search(queries[qi], k)
		expect += k
		got += utils.IntersectionCount(indexes, truth[qi][:k])
	}
	elapsed = time.Since(start)
	fmt.Println(len(queries))
	fmt.Printf("Search took %s\n", elapsed)

	recall := float64(got) / float64(expect)
	fmt.Println("Recall:", recall)
	assert.True(t, recall >= 0.96)
}
