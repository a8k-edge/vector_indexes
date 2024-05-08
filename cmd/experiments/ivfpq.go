package main

import (
	"fmt"
	"sort"
	"time"

	"mvdb/internals/index"
	"mvdb/internals/utils"
)

func ExpIVFPQ() {
	dim := 128
	divisors := getDivisors(dim)
	sort.Ints(divisors)
	baseData, truth, learn, queries := utils.LoadSift()
	k := 100

	dataSize := float64(len(baseData))
	cellsCountTrials := []int{
		int(dataSize * 0.01),
		int(dataSize * 0.1),
		int(dataSize * 0.3),
		int(dataSize * 0.5),
		int(dataSize * 0.8),
		int(dataSize * 0.9),
	}

	results := NewResult()

	for _, cellCount := range cellsCountTrials {
		for _, m := range divisors {
			if m == 1 || m == dim {
				continue
			}
			baseLabel := fmt.Sprintf("%d-%d", cellCount, m)
			start := time.Now()
			ivf := index.NewIVFPQ(dim, cellCount, 1, nil)
			ivf.Train(learn, index.NewPQ(dim, m))
			ivf.AddBatch(baseData)
			elapsed := time.Since(start)
			fmt.Printf("Construction took %s\n", elapsed)
			results.Add(ConstructionTime, baseLabel, elapsed.Seconds())

			nprobesTrials := []int{
				int(float64(cellCount) * 0.1),
				int(float64(cellCount) * 0.3),
			}
			for _, nprobes := range nprobesTrials {
				fmt.Println("********************")
				label := fmt.Sprintf("%s-%d", baseLabel, nprobes)
				fmt.Printf("cells-m-nprobes = %s \n", label)
				ivf.SetNprobes(nprobes)

				latency, recall := doSearch(ivf, k, queries, truth)
				results.Add(QueryLatency, label, latency)
				results.Add(Recall, label, recall)
			}
		}
	}

	fmt.Println("==============")
	fmt.Println(results)

	vizualize(results, "cells-nprobes", "IVFPQ")
	save(results, "IVFPQ")
}
