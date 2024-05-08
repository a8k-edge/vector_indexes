package main

import (
	"fmt"
	"strconv"
	"time"

	"mvdb/internals/index"
	"mvdb/internals/utils"
)

func ExpIVF() {
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

		start := time.Now()
		ivf := index.NewIVF(128, cellCount, 1, nil)
		ivf.Train(learn)
		ivf.AddBatch(baseData)
		elapsed := time.Since(start)
		fmt.Printf("Construction took %s\n", elapsed)
		results.Add(ConstructionTime, strconv.Itoa(cellCount), elapsed.Seconds())

		nprobesTrials := []int{
			int(float64(cellCount) * 0.1),
			int(float64(cellCount) * 0.3),
			int(float64(cellCount) * 0.5),
			int(float64(cellCount) * 0.8),
			int(float64(cellCount) * 0.9),
		}
		for _, nprobes := range nprobesTrials {
			fmt.Println("********************")
			label := fmt.Sprintf("%d-%d", cellCount, nprobes)
			fmt.Printf("cells-nprobes = %s  \n", label)
			ivf.SetNprobes(nprobes)

			latency, recall := doSearch(ivf, k, queries, truth)
			results.Add(QueryLatency, label, latency)
			results.Add(Recall, label, recall)
		}
	}

	fmt.Println("==============")
	fmt.Println(results)

	vizualize(results, "cells-nprobes", "IVF")
	save(results, "IVF")
}
