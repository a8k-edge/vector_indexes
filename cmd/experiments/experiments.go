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

type Results struct {
	Data map[ResultType][]ResultItem
}

type ResultItem struct {
	Label string
	Value float64
}

func NewResult() *Results {
	return &Results{
		Data: make(map[ResultType][]ResultItem),
	}
}

func (r *Results) Add(key ResultType, label string, value float64) {
	r.Data[key] = append(r.Data[key], ResultItem{label, value})
}

func (r *Results) String() string {
	var sb strings.Builder

	for resultType, order := range r.Data {
		sb.WriteString(fmt.Sprintf("%s:\n", resultType))
		for _, item := range order {
			sb.WriteString(fmt.Sprintf("\t%s: %f\n", item.Label, item.Value))
		}
	}

	return sb.String()
}
