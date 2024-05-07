package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const OutputPath = "output"

func init() {
	os.Mkdir(OutputPath, 0o777)
}

// TODO: consider ordered map for results item
// TODO: write csv output to a single file
func main() {
	var start time.Time
	var took time.Duration

	start = time.Now()
	ExpFlatPQ()
	ExpIVF()
	ExpIVFPQ()
	ExpHNSW()
	ExpVamana()
	// ExpPreTransform()

	took = time.Since(start)
	fmt.Printf("Took %s\n", took)
}

type ResultType string

const (
	ConstructionTime ResultType = "Construction time(sec)"
	QueryLatency     ResultType = "Query latency(ms)"
	Recall           ResultType = "Recall"
)

type Results map[ResultType]map[string]float64

func NewResult() Results {
	results := make(Results)
	results[ConstructionTime] = make(map[string]float64)
	results[QueryLatency] = make(map[string]float64)
	results[Recall] = make(map[string]float64)
	return results
}

func (r Results) String() string {
	var sb strings.Builder

	for resultType, data := range r {
		sb.WriteString(fmt.Sprintf("%s:\n", resultType))
		for key, value := range data {
			sb.WriteString(fmt.Sprintf("\t%s: %f\n", key, value))
		}
	}

	return sb.String()
}
