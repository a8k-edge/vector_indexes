package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"mvdb/internals/index"
	"mvdb/internals/utils"
)

func ExpFlatPQ() {
	dim := 128
	divisors := getDivisors(dim)
	sort.Ints(divisors)
	baseData, truth, learn, queries := utils.LoadSift()
	k := 100

	results := NewResult()

	for _, m := range divisors {
		if m == dim {
			continue
		}
		fmt.Println("********************")
		fmt.Println("m =", m)
		mstr := strconv.Itoa(m)

		start := time.Now()
		fpq := index.NewFlatPQ(dim, nil)
		fpq.Train(learn, index.NewPQ(128, m))
		fpq.AddBatch(baseData)
		elapsed := time.Since(start)
		fmt.Printf("Construction took %s\n", elapsed)
		results.Add(ConstructionTime, mstr, elapsed.Seconds())

		latency, recall := doSearch(fpq, k, queries, truth)
		results.Add(QueryLatency, mstr, latency)
		results.Add(Recall, mstr, recall)
	}

	fmt.Println("==============")
	fmt.Println(results)

	vizualize(results, "m parameter", "FlatPQ")
	save(results, "FlatPQ")
}

func getDivisors(n int) []int {
	divisors := []int{}

	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			divisors = append(divisors, i)
			if i != n/i {
				divisors = append(divisors, n/i)
			}
		}
	}

	return divisors
}
