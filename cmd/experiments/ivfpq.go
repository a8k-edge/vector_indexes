package main

import (
	"fmt"
	"image/color"
	"sort"
	"time"

	"mvdb/internals/index"
	"mvdb/internals/utils"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// TODO: better UX
// TODO: more relevant experiment variables
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

	results := map[string]map[string]float64{
		"Construction time(sec) per cells count": make(map[string]float64),
		"Query latency(ms)":                      make(map[string]float64),
		"Recall":                                 make(map[string]float64),
	}

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
			results["Construction time(sec) per cells count"][baseLabel] = elapsed.Seconds()

			nprobesTrials := []int{
				int(float64(cellCount) * 0.1),
				int(float64(cellCount) * 0.3),
			}
			for _, nprobes := range nprobesTrials {
				fmt.Println("********************")
				label := fmt.Sprintf("%s-%d", baseLabel, nprobes)
				fmt.Printf("Work on %s (cells-m-nprobes) \n", label)
				ivf.SetNprobes(nprobes)

				start = time.Now()
				_, indexes := index.SearchMany(ivf, queries, k)
				elapsed = time.Since(start)
				latencyNS := elapsed.Nanoseconds() / int64(len(queries))
				fmt.Printf("Search took %s %f\n", elapsed, float64(latencyNS)/1e+6)
				results["Query latency(ms)"][label] = float64(latencyNS) / 1e+6

				expect, got := 0, 0
				for i := range queries {
					expect += k
					got += utils.IntersectionCount(indexes[i], truth[i][:k])
				}
				recall := float64(got) / float64(expect)
				fmt.Printf("Recall %f\n", recall)
				results["Recall"][label] = recall
			}
		}
	}
	fmt.Println(results)

	for label, data := range results {
		p := plot.New()
		p.Title.Text = label
		p.X.Label.Text = "cells-nprobes"

		values := make(plotter.Values, len(data))
		labels := make([]string, len(data))
		i := 0
		for key, value := range data {
			labels[i] = key
			values[i] = value
			i++
		}

		bars, err := plotter.NewBarChart(values, vg.Points(20))
		if err != nil {
			panic(err)
		}
		color := color.RGBA{R: 0, G: 128, B: 255, A: 255} // Blue
		bars.Color = color
		p.Add(bars)
		p.NominalX(labels...)

		fname := fmt.Sprintf("%s.png", label)
		if err := p.Save(60*vg.Inch, 60*vg.Inch, fname); err != nil {
			panic(err)
		}
	}
}
