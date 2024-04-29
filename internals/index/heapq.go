package index

import (
	"fmt"

	"mvdb/internals/vops"
)

// TODO: have 2 different heapQ types (min, max)

type HeapQ struct {
	items []Item
	less  func(items []Item, i, j int) bool
}

func NewMinHeapQ(cap int) *HeapQ {
	return &HeapQ{
		items: make([]Item, 0, cap),
		less: func(items []Item, i, j int) bool {
			return items[i].Distance < items[j].Distance
		},
	}
}

func NewMaxHeapQ(cap int) *HeapQ {
	return &HeapQ{
		items: make([]Item, 0, cap),
		less: func(items []Item, i, j int) bool {
			return items[i].Distance > items[j].Distance
		},
	}
}

func (h *HeapQ) Print() {
	for i := range h.items {
		fmt.Printf("%f ", h.items[i].Distance)
	}
	fmt.Println()
}

func (h *HeapQ) Top() Item {
	return h.items[0]
}

func (h *HeapQ) Len() int {
	return len(h.items)
}

func (h *HeapQ) Pop() Item {
	n := h.Len() - 1
	h.swap(0, n)
	h.down(0, n)

	item := h.items[n]
	h.items = h.items[:n]
	return item
}

func (h *HeapQ) Push(id int, distance float32) {
	item := Item{
		ID:       id,
		Distance: distance,
	}
	h.items = append(h.items, item)
	h.up(h.Len() - 1)
}

func (h *HeapQ) up(i int) {
	for {
		parent := (i - 1) / 2
		if parent == i || !h.less(h.items, i, parent) {
			break
		}
		h.swap(parent, i)
		i = parent
	}
}

func (h *HeapQ) down(start int, n int) {
	i := start
	for {
		lc := 2*i + 1
		if lc >= n {
			break
		}
		child := lc
		if rc := lc + 1; rc < n && h.less(h.items, rc, lc) {
			child = rc
		}
		if !h.less(h.items, child, i) {
			break
		}
		h.swap(child, i)
		i = child
	}
}

func (h *HeapQ) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (q *HeapQ) PopAll() ([]float32, []int) {
	size := q.Len()
	topKD := make([]float32, size)
	topKI := make([]int, size)

	for i := size - 1; i >= 0; i-- {
		item := q.Pop()
		topKD[i] = item.Distance
		topKI[i] = item.ID
	}

	return topKD, topKI
}

// Should be used for MaxHeap
func (q *HeapQ) PopMin() Item {
	last := len(q.items) - 1
	i := last

	minItem := q.items[i]
	minIndex := i
	i--
	for i >= 0 {
		if minItem.Distance > q.items[i].Distance {
			minItem = q.items[i]
			minIndex = i
		}
		i--
	}

	if minIndex != last {
		q.swap(minIndex, last)
		for j := len(q.items) / 2; j < last; j++ {
			q.up(j)
		}
	}
	q.items = q.items[:last]

	return minItem
}

// TODO: test it
func NClosestTo(n int, query []float32, vectors [][]float32, metric vops.Provider) []int {
	if n == 1 {
		clst := 0
		minDist := metric.Similarity(query, vectors[0])
		for i, vector := range vectors {
			distance := metric.Similarity(query, vector)
			if distance < minDist {
				clst = i
				minDist = distance
			}
		}
		return []int{clst}
	}

	futherQ := NewMaxHeapQ(n)

	for i := 0; i < len(vectors); i++ {
		distance := metric.Similarity(query, vectors[i])

		if i < n {
			futherQ.Push(i, distance)
		} else if distance < futherQ.Top().Distance {
			futherQ.Pop()
			futherQ.Push(i, distance)
		}
	}

	clst := make([]int, futherQ.Len())
	for i := futherQ.Len() - 1; i >= 0; i-- {
		clst[i] = futherQ.Pop().ID
	}

	return clst
}
