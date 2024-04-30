package index

type RerankIndex struct {
	base  Index
	store *Flat

	expandFactor float64
}

// TODO: encapsulate base's & store's Train, Add, etc.?
// Train signature may not be generic for future

func NewRerankIndex(base Index, store *Flat, ef float64) *RerankIndex {
	return &RerankIndex{
		base:         base,
		store:        store,
		expandFactor: ef,
	}
}

func (rr *RerankIndex) SetExpandFactor(ef float64) {
	rr.expandFactor = ef
}

func (rr *RerankIndex) Search(query []float32, k int) ([]float32, []int) {
	expandedK := int(float64(k) * rr.expandFactor)
	_, I := rr.base.Search(query, expandedK)

	futherQ := NewMaxHeapQ(k)
	for i := 0; i < len(I); i++ {
		id := I[i]
		dist := rr.store.ComputeDistTo(query, id)
		if i < k {
			futherQ.Push(id, dist)
		} else if dist < futherQ.Top().Distance {
			futherQ.Pop()
			futherQ.Push(id, dist)
		}
	}
	return futherQ.PopAll()
}
