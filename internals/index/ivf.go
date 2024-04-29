package index

import (
	"sync"

	"mvdb/internals/vops"
)

type IVF struct {
	vectors     [][]float32
	vidxByCells [][]int
	centroids   [][]float32

	dim        uint
	cellsCount int
	nprobes    int

	isTrained bool
	metric    vops.Provider
}

func NewIVF(dim uint, cellsCount int, nprobes int, metric vops.Provider) *IVF {
	if metric == nil {
		metric = &vops.L2SqrDistance{}
	}
	return &IVF{
		dim:        dim,
		cellsCount: cellsCount,
		nprobes:    nprobes,
		metric:     metric,
	}
}

func (ivf *IVF) SetNprobes(nprobes int) {
	ivf.nprobes = nprobes
}

func (ivf *IVF) Train(data [][]float32) {
	ivf.centroids, _ = LloydKmeans(data, ivf.cellsCount, 10, ivf.metric)
	ivf.vidxByCells = make([][]int, len(ivf.centroids))
	ivf.isTrained = true
}

func (ivf *IVF) Add(vector []float32) {
	ivf.AddBatch([][]float32{vector})
}

func (ivf *IVF) AddBatch(vectors [][]float32) {
	if !ivf.isTrained {
		panic("Can't add to untrained index")
	}
	for _, v := range vectors {
		cell := NClosestTo(1, v, ivf.centroids, ivf.metric)[0]
		ivf.vectors = append(ivf.vectors, v)
		ivf.vidxByCells[cell] = append(ivf.vidxByCells[cell], len(ivf.vectors)-1)
	}
}

func (ivf *IVF) Search(query []float32, k int) ([]float32, []int) {
	cells := NClosestTo(ivf.nprobes, query, ivf.centroids, ivf.metric)
	futherQ := NewMaxHeapQ(k)

	for _, cellIdx := range cells {
		for _, vid := range ivf.vidxByCells[cellIdx] {
			v := ivf.vectors[vid]
			distance := ivf.metric.Similarity(query, v)

			if futherQ.Len() < k {
				futherQ.Push(vid, distance)
			} else if distance < futherQ.Top().Distance {
				futherQ.Pop()
				futherQ.Push(vid, distance)
			}
		}
	}

	return futherQ.PopAll()
}

func (f *IVF) SearchMany(queries [][]float32, k int) ([][]float32, [][]int) {
	// TODO: somehow throttle concurrency to min(len(queries), max procs)
	var wg sync.WaitGroup

	distances := make([][]float32, len(queries))
	indexes := make([][]int, len(queries))

	for qi := range queries {
		wg.Add(1)
		go func(qi int) {
			defer wg.Done()

			d, i := f.Search(queries[qi], k)
			distances[qi] = d
			indexes[qi] = i
		}(qi)
	}
	wg.Wait()

	return distances, indexes
}
