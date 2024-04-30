package main

import (
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"time"

	"mvdb/internals/index"
	"mvdb/internals/utils"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func ExpFlatPQ() {
	dim := 128
	divisors := getDivisors(dim)
	sort.Ints(divisors)
	baseData, truth, learn, queries := utils.LoadSift()
	k := 100

	results := map[string]map[int]float64{
		"Construction time(sec) per m": make(map[int]float64),
		"Query latency(ms) per m":      make(map[int]float64),
		"Recall per m":                 make(map[int]float64),
	}

	for _, m := range divisors {
		if m == 1 || m == dim {
			continue
		}
		fmt.Println("********************")
		fmt.Println("Work on ", m)

		start := time.Now()
		fpq := index.NewFlatPQ(dim, nil)
		fpq.Train(learn, index.NewPQ(128, m))
		fpq.AddBatch(baseData)
		elapsed := time.Since(start)
		fmt.Printf("Construction took %s\n", elapsed)
		results["Construction time(sec) per m"][m] = elapsed.Seconds()

		start = time.Now()
		_, indexes := index.SearchMany(fpq, queries, k)
		elapsed = time.Since(start)
		latencyNS := elapsed.Nanoseconds() / int64(len(queries))
		fmt.Printf("Search took %s %f\n", elapsed, float64(latencyNS)/1e+6)
		results["Query latency(ms) per m"][m] = float64(latencyNS) / 1e+6

		expect, got := 0, 0
		for i := range queries {
			expect += k
			got += utils.IntersectionCount(indexes[i], truth[i][:k])
		}
		recall := float64(got) / float64(expect)
		fmt.Printf("Recall %f\n", recall)
		results["Recall per m"][m] = recall
	}
	fmt.Println(results)

	for label, data := range results {
		p := plot.New()
		p.Title.Text = label
		p.X.Label.Text = "m parameter"

		values := make(plotter.Values, len(data))
		labels := make([]string, len(data))
		i := 0
		for m, value := range data {
			labels[i] = strconv.Itoa(m)
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
		if err := p.Save(4*vg.Inch, 4*vg.Inch, fname); err != nil {
			panic(err)
		}
	}
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
