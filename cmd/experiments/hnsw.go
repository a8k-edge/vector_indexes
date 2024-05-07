package main

import (
	"fmt"
	"time"

	"mvdb/internals/index"
	"mvdb/internals/utils"
)

func ExpHNSW() {
	baseData, truth, _, queries := utils.LoadSift()
	k := 100

	results := NewResult()

	ms := []int{6, 16, 20, 32, 40}
	efConstructions := []int{40, 150}
	efs := []int{10, 40, 100, 200}

	for _, M := range ms {
		for _, efConstruction := range efConstructions {
			baseLabel := fmt.Sprintf("%d-%d", M, efConstruction)
			start := time.Now()
			hnsw := index.NewHNSW(M, 1, efConstruction, nil)
			hnsw.AddBatch(baseData)
			elapsed := time.Since(start)
			fmt.Printf("Construction took %s\n", elapsed)
			results[ConstructionTime][baseLabel] = elapsed.Seconds()

			for _, ef := range efs {
				fmt.Println("********************")
				label := fmt.Sprintf("%s-%d", baseLabel, ef)
				fmt.Printf("M-efConst.-efSearch = %s \n", label)
				hnsw.SetEf(ef)

				latency, recall := doSearch(hnsw, k, queries, truth)
				results[QueryLatency][label] = latency
				results[Recall][label] = recall
			}
		}
	}

	fmt.Println("==============")
	fmt.Println(results)

	vizualize(results, "M-efConst.-efSearch", "HNSW")
	save(results, "HNSW")
}
