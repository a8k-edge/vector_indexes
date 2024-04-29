package main

import (
	"fmt"
	"image/color"
	"time"

	"mvdb/internals/index"
	"mvdb/internals/utils"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func ExpVamana() {
	baseData, truth, _, queries := utils.LoadSift()
	k := 100

	results := map[string]map[string]float64{
		"Construction time(sec)": make(map[string]float64),
		"Query latency(ms)":      make(map[string]float64),
		"Recall":                 make(map[string]float64),
	}

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
				results["Construction time(sec)"][label] = elapsed.Seconds()

				start = time.Now()
				_, indexes := vmn.SearchMany(queries, k)
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
		p.X.Label.Text = "R-L-alpha"

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
		if err := p.Save(40*vg.Inch, 40*vg.Inch, fname); err != nil {
			panic(err)
		}
	}
}
