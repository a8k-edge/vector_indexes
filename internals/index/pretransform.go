package index

import "mvdb/internals/num"

type VectorTransform interface {
	Apply(vectors [][]float32) [][]float32
}

type RandomRotationMatrix struct {
	// flat matrix
	matrix []float32
	dim    int
}

func NewRandomRotationMatrix(dim int) *RandomRotationMatrix {
	return &RandomRotationMatrix{
		dim:    dim,
		matrix: num.RandRotationFMatrix(dim),
	}
}

func (rrm *RandomRotationMatrix) Apply(vectors [][]float32) [][]float32 {
	size := len(vectors)
	result := make([][]float32, size)

	for i, vector := range vectors {
		result[i] = make([]float32, rrm.dim)
		for ri := 0; ri < rrm.dim; ri++ {
			for ci := 0; ci < rrm.dim; ci++ {
				result[i][ri] += rrm.matrix[ri*rrm.dim+ci] * vector[ci]
			}
		}
	}

	return result
}

type IndexPreTransform struct {
	index     Index
	transform VectorTransform
}

func NewIndexPreTransform(index Index, transform VectorTransform) *IndexPreTransform {
	return &IndexPreTransform{
		index:     index,
		transform: transform,
	}
}

func (it *IndexPreTransform) Add(vector []float32) {
	it.AddBatch([][]float32{vector})
}

func (it *IndexPreTransform) AddBatch(vectors [][]float32) {
	tvectors := it.transform.Apply(vectors)
	it.index.AddBatch(tvectors)
}

func (it *IndexPreTransform) Search(query []float32, k int) ([]float32, []int) {
	tquery := it.transform.Apply([][]float32{query})[0]
	return it.index.Search(tquery, k)
}
