package index

import (
	"sync"
	"unsafe"

	"mvdb/internals/vops"
)

type FlatPQ struct {
	Codes [][]int
	dim   int

	isTrained bool
	pq        *ProductQuantizer

	metric vops.Provider
}

func NewFlatPQ(dim int, metric vops.Provider) *FlatPQ {
	if metric == nil {
		metric = &vops.L2SqrDistance{}
	}
	return &FlatPQ{
		dim:    dim,
		metric: metric,
	}
}

func (fpq *FlatPQ) SizePQ() uint64 {
	return fpq.pq.Size()
}

func (fpq *FlatPQ) Size() uint64 {
	// Size of slice header
	sliceHeaderSize := unsafe.Sizeof(fpq.Codes)

	// Size of each int slice header
	intSliceHeaderSize := unsafe.Sizeof(fpq.Codes[0])

	// Size of each int value
	intSize := unsafe.Sizeof(fpq.Codes[0][0])
	codesDim := len(fpq.Codes[0])

	totalSize := sliceHeaderSize + intSliceHeaderSize*uintptr(len(fpq.Codes)) + intSize*uintptr(len(fpq.Codes)*codesDim)

	return uint64(totalSize)
}

func (f *FlatPQ) Train(data [][]float32, pq *ProductQuantizer) {
	f.pq = pq
	f.pq.metric = f.metric

	f.pq.Train(data)
	f.isTrained = true
}

func (f *FlatPQ) Add(vector []float32) {
	f.AddBatch([][]float32{vector})
}

func (f *FlatPQ) AddBatch(vectors [][]float32) {
	if !f.isTrained {
		panic("Can't add to untrained index")
	}
	for _, v := range vectors {
		codes := f.pq.Encode(v)
		f.Codes = append(f.Codes, codes)
	}
}

func (f *FlatPQ) Delete(id int) {
	f.Codes = append(f.Codes[:id], f.Codes[id+1:]...)
}

func (f *FlatPQ) Search(query []float32, k int) ([]float32, []int) {
	distTable := f.pq.DistTable(query)
	futherQ := NewMaxHeapQ(k)
	i := 0

	for ; i < k && i < len(f.Codes); i++ {
		var distance float32
		for j, code := range f.Codes[i] {
			distance += distTable[j][code]
		}
		futherQ.Push(i, distance)
	}

	for ; i < len(f.Codes); i++ {
		var distance float32
		for j, code := range f.Codes[i] {
			distance += distTable[j][code]
		}
		if distance < futherQ.Top().Distance {
			futherQ.Pop()
			futherQ.Push(i, distance)
		}
	}

	return futherQ.PopAll()
}

func (f *FlatPQ) SearchMany(queries [][]float32, k int) ([][]float32, [][]int) {
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
