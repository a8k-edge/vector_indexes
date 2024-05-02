package index

import (
	"unsafe"

	"mvdb/internals/vops"
)

type Flat struct {
	data [][]float32
	dim  int

	metric vops.Provider
}

func NewFlat(dim int, metric vops.Provider) *Flat {
	if metric == nil {
		metric = &vops.L2SqrDistance{}
	}
	return &Flat{
		dim:    dim,
		metric: metric,
	}
}

func (f *Flat) Size() uint64 {
	sliceHeaderSize := unsafe.Sizeof(f.data)
	floatSliceHeaderSize := unsafe.Sizeof(f.data[0])
	floatSize := unsafe.Sizeof(f.data[0][0])
	totalSize := sliceHeaderSize + floatSliceHeaderSize*uintptr(len(f.data)) + floatSize*uintptr(len(f.data)*int(f.dim))

	return uint64(totalSize)
}

func (f *Flat) AddBatch(vectors [][]float32) {
	f.data = append(f.data, vectors...)
}

func (f *Flat) Search(query []float32, k int) ([]float32, []int) {
	futherQ := NewMaxHeapQ(k)

	i := 0
	for ; i < k && i < len(f.data); i++ {
		distance := f.metric.Similarity(query, f.data[i])
		futherQ.Push(i, distance)
	}

	for ; i < len(f.data); i++ {
		distance := f.metric.Similarity(query, f.data[i])
		if distance < futherQ.Top().Distance {
			futherQ.Pop()
			futherQ.Push(i, distance)
		}
	}

	return futherQ.PopAll()
}

func (f *Flat) ComputeDistTo(query []float32, id int) float32 {
	return f.metric.Similarity(query, f.data[id])
}
