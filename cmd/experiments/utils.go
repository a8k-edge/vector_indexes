package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"mvdb/internals/index"
	"mvdb/internals/utils"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func doSearch(i index.Index, k int, queries [][]float32, truth [][]int) (latency float64, recall float64) {
	start := time.Now()
	_, indexes := index.SearchMany(i, queries, k)
	elapsed := time.Since(start)

	latencyNS := elapsed.Nanoseconds() / int64(len(queries))
	latency = float64(latencyNS) / 1e+6
	fmt.Printf("Search took %.2fms (total %s)\n", latency, elapsed)

	expect, got := 0, 0
	for i := range queries {
		expect += k
		got += utils.IntersectionCount(indexes[i], truth[i][:k])
	}
	recall = float64(got) / float64(expect)
	fmt.Printf("Recall %.4f\n", recall)

	return latency, recall
}

func vizualize(results *Results, xLabel string, expName string) {
	basePath := filepath.Join(OutputPath, expName)
	os.MkdirAll(basePath, os.ModePerm)

	for label, data := range results.Data {
		p := plot.New()
		p.Title.Text = string(label)
		p.X.Label.Text = xLabel

		values := make(plotter.Values, len(data))
		labels := make([]string, len(data))
		for i, value := range data {
			labels[i] = value.Label
			values[i] = value.Value
		}

		bars, err := plotter.NewBarChart(values, vg.Points(20))
		if err != nil {
			panic(err)
		}
		color := color.RGBA{R: 0, G: 128, B: 255, A: 255} // Blue
		bars.Color = color
		p.Add(bars)
		p.NominalX(labels...)
		p.X.Tick.Label.Rotation = -math.Pi / 3 // 60 degree
		p.X.Tick.Label.XAlign = 0.1

		factor := vg.Length(4)
		if len(data) > 10 {
			factor = 14
		}

		fname := fmt.Sprintf("%s.png", label)
		location := filepath.Join(basePath, fname)
		if err := p.Save(factor*vg.Inch, factor*vg.Inch, location); err != nil {
			panic(err)
		}
	}
}

func save(results *Results, expName string) {
	dataList := make([][]string, 0)
	for category, methods := range results.Data {
		dataList = append(dataList, []string{string(category), "Method", "Value"})

		for _, value := range methods {
			dataList = append(dataList, []string{"", value.Label, formatFloat(value.Value)})
		}

		dataList = append(dataList, []string{"", "", ""})
	}

	basePath := filepath.Join(OutputPath, expName)
	os.MkdirAll(basePath, os.ModePerm)
	location := filepath.Join(basePath, "results.csv")
	file, err := os.Create(location)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range dataList {
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
