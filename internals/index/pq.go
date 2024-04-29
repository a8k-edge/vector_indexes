package index

import "container/heap"

type Item struct {
	Distance float32
	ID       int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Top() *Item {
	return pq[0]
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Pop the highest dist
	return pq[i].Distance > pq[j].Distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Item))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func MinHeap() PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&minHeap{pq})
	return pq
}

// maxHeap is a wrapper around PriorityQueue to use it as a max heap.
type minHeap struct {
	PriorityQueue
}

// Less returns true if the element at index i has higher priority than the element at index j.
// For a max heap, it compares based on the negative of distance.
//
//	func (mh maxHeap) Less(i, j int) bool {
//		return mh.PriorityQueue[i].Distance > mh.PriorityQueue[j].Distance
//	}
func (mh minHeap) Less(i, j int) bool {
	return mh.PriorityQueue[i].Distance < mh.PriorityQueue[j].Distance
}
