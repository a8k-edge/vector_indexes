package index

type Index interface {
	Add(vector []float32)
	AddBatch(vectors [][]float32)
	Search(query []float32, k int) ([]float32, []int)
}
