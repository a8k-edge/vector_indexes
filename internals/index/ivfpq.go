package index

import (
	"mvdb/internals/vops"
)

type IVFPQ struct {
	codes       [][]int
	vidxByCells [][]int
	centroids   [][]float32

	dim        int
	cellsCount int
	nprobes    int

	isTrained bool
	pq        *ProductQuantizer
	metric    vops.Provider
}

func NewIVFPQ(dim int, cellsCount int, nprobes int, metric vops.Provider) *IVFPQ {
	if metric == nil {
		metric = &vops.L2SqrDistance{}
	}
	return &IVFPQ{
		dim:        dim,
		cellsCount: cellsCount,
		nprobes:    nprobes,
		metric:     metric,
	}
}

func (ivf *IVFPQ) SetNprobes(nprobes int) {
	ivf.nprobes = nprobes
}

func (ivf *IVFPQ) Train(data [][]float32, pq *ProductQuantizer) {
	var labels []int
	ivf.centroids, labels = LloydKmeans(data, ivf.cellsCount, 10, ivf.metric)

	ivf.vidxByCells = make([][]int, len(ivf.centroids))

	trainData := make([][]float32, len(data))
	for i, vector := range data {
		centroid := labels[i]
		trainData[i] = residualCalc(ivf.centroids[centroid], vector)
	}

	ivf.pq = pq
	ivf.pq.metric = ivf.metric
	ivf.pq.Train(trainData)
	ivf.isTrained = true
}

func (ivf *IVFPQ) Add(vector []float32) {
	ivf.AddBatch([][]float32{vector})
}

func (ivf *IVFPQ) AddBatch(vectors [][]float32) {
	if !ivf.isTrained {
		panic("Can't add to untrained index")
	}

	for _, vector := range vectors {
		clstCentroid := NClosestTo(1, vector, ivf.centroids, ivf.metric)[0]

		residual := residualCalc(ivf.centroids[clstCentroid], vector)
		codes := ivf.pq.Encode(residual)

		ivf.codes = append(ivf.codes, codes)
		ivf.vidxByCells[clstCentroid] = append(ivf.vidxByCells[clstCentroid], len(ivf.codes)-1)
	}
}

func (ivf *IVFPQ) Search(query []float32, k int) ([]float32, []int) {
	cells := NClosestTo(ivf.nprobes, query, ivf.centroids, ivf.metric)
	futherQ := NewMaxHeapQ(k)

	for _, cellIdx := range cells {
		rquery := residualCalc(ivf.centroids[cellIdx], query)
		distTable := ivf.pq.DistTable(rquery)
		for _, id := range ivf.vidxByCells[cellIdx] {
			var distance float32
			for j, code := range ivf.codes[id] {
				distance += distTable[j][code]
			}

			if futherQ.Len() < k {
				futherQ.Push(id, distance)
			} else if distance < futherQ.Top().Distance {
				futherQ.Pop()
				futherQ.Push(id, distance)
			}
		}
	}

	return futherQ.PopAll()
}

func residualCalc(center []float32, vector []float32) []float32 {
	offset := make([]float32, len(vector))
	for i := range vector {
		offset[i] = center[i] - vector[i]
	}
	return offset
}
