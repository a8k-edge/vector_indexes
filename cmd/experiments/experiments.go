package main

import (
	"fmt"
	"time"
)

// TODO: verbose option
// TODO: save results to csv
// TODO: dir for files output
func main() {
	var start time.Time
	var took time.Duration

	start = time.Now()
	// ExpFlatPQ()
	// ExpIVF()
	// ExpIVFPQ()
	// ExpHNSW()
	// ExpVamana()
	ExpPreTransform()

	took = time.Since(start)
	fmt.Printf("Exp toook %s\n", took)
}
