package index

import (
	"sync"
	"unsafe"

	"mvdb/internals/vops"
)

type Flat struct {
	Data [][]float32
	dim  uint

	metric vops.Provider
}

func NewFlat(dim uint, metric vops.Provider) *Flat {
	if metric == nil {
		metric = &vops.L2SqrDistance{}
	}
	return &Flat{
		dim:    dim,
		metric: metric,
	}
}

func (f *Flat) Size() uint64 {
	sliceHeaderSize := unsafe.Sizeof(f.Data)
	floatSliceHeaderSize := unsafe.Sizeof(f.Data[0])
	floatSize := unsafe.Sizeof(f.Data[0][0])
	totalSize := sliceHeaderSize + floatSliceHeaderSize*uintptr(len(f.Data)) + floatSize*uintptr(len(f.Data)*int(f.dim))

	return uint64(totalSize)
}

func (f *Flat) Add(vector []float32) {
	f.AddBatch([][]float32{vector})
}

func (f *Flat) AddBatch(vectors [][]float32) {
	// Several maginitude slower than memcpy
	// TODO: try bench with copy
	f.Data = append(f.Data, vectors...)
}

func (f *Flat) Search(query []float32, k int) ([]float32, []int) {
	futherQ := NewMaxHeapQ(k)

	i := 0
	for ; i < k && i < len(f.Data); i++ {
		distance := f.metric.Similarity(query, f.Data[i])
		futherQ.Push(i, distance)
	}

	for ; i < len(f.Data); i++ {
		distance := f.metric.Similarity(query, f.Data[i])
		if distance < futherQ.Top().Distance {
			futherQ.Pop()
			futherQ.Push(i, distance)
		}
	}

	return futherQ.PopAll()
}

func (f *Flat) SearchMany(queries [][]float32, k int) ([][]float32, [][]int) {
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
