package index

import (
	"unsafe"

	"mvdb/internals/vops"
)

type ProductQuantizer struct {
	dim     int
	m       int
	subsize int
	k       int

	centroids [][][]float32
	metric    vops.Provider
}

func NewPQ(dim int, m int) *ProductQuantizer {
	if dim%m != 0 {
		panic("Invalid dim or m")
	}
	// k := int(math.Pow(2, float64(m)))
	k := 1000
	return &ProductQuantizer{
		dim:       dim,
		m:         m,
		subsize:   dim / m,
		k:         k,
		centroids: make([][][]float32, m),
	}
}

func (pq *ProductQuantizer) Size() uint64 {
	sliceHeaderSize := unsafe.Sizeof(pq.centroids)
	float32SliceHeaderSize := unsafe.Sizeof(pq.centroids[0])
	float32Size := unsafe.Sizeof(pq.centroids[0][0][0])
	totalSize := sliceHeaderSize + float32SliceHeaderSize*uintptr(len(pq.centroids)) + float32Size*uintptr(len(pq.centroids)*pq.m*pq.subsize*pq.k)
	return uint64(totalSize)
}

func (p *ProductQuantizer) Train(data [][]float32) {
	for j := 0; j < p.m; j++ {
		p.centroids[j], _ = LloydKmeansBySection(p.subsize, j, data, p.k, 10, p.metric)
	}
}

func (p *ProductQuantizer) Encode(vector []float32) []int {
	codes := make([]int, p.m)
	for j := 0; j < p.m; j++ {
		codes[j] = NClosestTo(1, vector[j*p.subsize:j*p.subsize+p.subsize], p.centroids[j], p.metric)[0]
	}
	return codes
}

func (p *ProductQuantizer) DistTable(query []float32) [][]float32 {
	distTable := make([][]float32, p.m)
	for j := 0; j < p.m; j++ {
		distTable[j] = make([]float32, p.k)
		for i, centroid := range p.centroids[j] {
			distTable[j][i] = p.metric.Similarity(query[j*p.subsize:j*p.subsize+p.subsize], centroid)
		}
	}
	return distTable
}
