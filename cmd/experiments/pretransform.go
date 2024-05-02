package main

import (
	"fmt"

	"mvdb/internals/index"
	"mvdb/internals/utils"
)

func ExpPreTransform() {
	dim := 128
	baseData, truth, _, queries := utils.LoadSift()
	k := 100

	f := index.NewFlat(dim, nil)
	transform := index.NewRandomRotationMatrix(dim)

	pretransform := index.NewIndexPreTransform(f, transform)

	pretransform.AddBatch(baseData)

	got, expect := 0, 0
	for qi := range queries {
		_, indexes := pretransform.Search(queries[qi], k)
		expect += k
		got += utils.IntersectionCount(indexes, truth[qi][:k])
	}

	recall := float64(got) / float64(expect)
	fmt.Println("Recall:", recall)
}
