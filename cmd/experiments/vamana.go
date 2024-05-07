package main

import (
	"fmt"
	"time"

	"mvdb/internals/index"
	"mvdb/internals/utils"
)

func ExpVamana() {
	baseData, truth, _, queries := utils.LoadSift()
	k := 100

	results := NewResult()

	rs := []int{6, 16, 20, 32, 40}
	ls := []int{40, 150}
	alphas := []float32{1.0, 1.2, 1.4}

	for _, R := range rs {
		for _, L := range ls {
			for _, alpha := range alphas {
				fmt.Println("********************")
				label := fmt.Sprintf("%d-%d-%f", R, L, alpha)
				fmt.Printf("Work on %s (R-L-alpha) \n", label)

				start := time.Now()
				vmn := index.NewVamana(R, L, alpha, nil)
				vmn.AddBatch(baseData)
				elapsed := time.Since(start)
				fmt.Printf("Construction took %s\n", elapsed)
				results[ConstructionTime][label] = elapsed.Seconds()

				latency, recall := doSearch(vmn, k, queries, truth)
				results[QueryLatency][label] = latency
				results[Recall][label] = recall
			}
		}
	}

	fmt.Println("==============")
	fmt.Println(results)

	vizualize(results, "M-efConst.-efSearch", "Vamana")
	save(results, "Vamana")
}
