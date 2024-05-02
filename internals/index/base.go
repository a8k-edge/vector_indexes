package index

import (
	"math"
	"runtime"
	"sync"
)

type Index interface {
	AddBatch(vectors [][]float32)
	Search(query []float32, k int) ([]float32, []int)
}

func SearchMany(index Index, queries [][]float32, k int) ([][]float32, [][]int) {
	distances := make([][]float32, len(queries))
	indexes := make([][]int, len(queries))

	ParallelFor(len(queries), func(qi int) {
		d, i := index.Search(queries[qi], k)
		distances[qi] = d
		indexes[qi] = i
	})

	return distances, indexes
}

func ParallelFor(n int, action func(taskIndex int)) {
	n64 := float64(n)
	workerCount := runtime.GOMAXPROCS(0)
	wg := &sync.WaitGroup{}
	split := int(math.Ceil(n64 / float64(workerCount)))
	for worker := 0; worker < workerCount; worker++ {
		workerID := worker

		wg.Add(1)
		func() {
			defer wg.Done()
			for i := workerID * split; i < int(math.Min(float64((workerID+1)*split), n64)); i++ {
				action(i)
			}
		}()
	}
	wg.Wait()
}
