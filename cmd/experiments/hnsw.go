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

func ExpHNSW() {
	baseData, truth, _, queries := utils.LoadSift()
	k := 100

	results := map[string]map[string]float64{
		"Construction time(sec)": make(map[string]float64),
		"Query latency(ms)":      make(map[string]float64),
		"Recall":                 make(map[string]float64),
	}

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
			results["Construction time(sec)"][baseLabel] = elapsed.Seconds()

			for _, ef := range efs {
				fmt.Println("********************")
				label := fmt.Sprintf("%s-%d", baseLabel, ef)
				fmt.Printf("Work on %s (M-efConst.-efSearch) \n", label)
				hnsw.SetEf(ef)

				start = time.Now()
				_, indexes := hnsw.SearchMany(queries, k)
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
		p.X.Label.Text = "M-efConst.-efSearch"

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
		if err := p.Save(50*vg.Inch, 50*vg.Inch, fname); err != nil {
			panic(err)
		}
	}
}
