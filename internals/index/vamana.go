package index

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"mvdb/internals/vops"
)

const GraphSlackFactor = 1.3

type Vamana struct {
	vectors   [][]float32
	neighbors [][]int

	start int
	R     int
	L     int
	alpha float32

	metric vops.Provider
}

func NewVamana(R int, L int, alpha float32, metric vops.Provider) *Vamana {
	if metric == nil {
		metric = &vops.L2SqrDistance{}
	}
	return &Vamana{
		neighbors: make([][]int, 0),
		R:         R,
		L:         L,
		alpha:     alpha,
		metric:    metric,
	}
}

func (v *Vamana) Print() {
	fmt.Println("VAMANA")
	maxDegree := 0
	totalDegree := 0
	minDegree := math.MaxInt
	countLess2Degree := 0

	for i := range v.vectors {
		l := len(v.neighbors[i])
		maxDegree = max(maxDegree, l)
		minDegree = min(minDegree, l)
		totalDegree += l
		if l < 2 {
			countLess2Degree++
		}
	}
	fmt.Println("Max Degree: ", maxDegree)
	fmt.Println("Avg Degree: ", totalDegree/len(v.vectors))
	fmt.Println("Min Degree: ", minDegree)
	fmt.Println("Count < 2 degree: ", countLess2Degree)
}

func (v *Vamana) AddBatch(vectors [][]float32) {
	nstart := len(v.vectors)
	v.vectors = append(v.vectors, vectors...)

	v.start = nstart + v.findMedoid(vectors)

	for id := range v.vectors {
		idxs := make([]int, v.R)
		exists := make(map[int]bool)

		for i := 0; i < min(v.R, len(vectors)-1); {
			randomIndex := rand.Intn(len(v.vectors))
			if !exists[randomIndex] && randomIndex != id {
				exists[randomIndex] = true
				idxs = append(idxs, randomIndex)
				i++
			}
		}

		v.neighbors = append(v.neighbors, idxs)
	}

	vectorsOrder := make([]int, 0, len(vectors))
	for _, id := range rand.Perm(len(vectors)) {
		vectorsOrder = append(vectorsOrder, id)
	}

	for i := range vectorsOrder {
		id := nstart + vectorsOrder[i]
		// if i%10000 == 0 {
		// 	fmt.Println(time.Now(), i)
		// }

		candidates := v.candidatesToAddNew(id, v.vectors[id])
		pruned := v.prune(id, v.vectors[id], candidates)

		v.neighbors[id] = pruned

		for _, nid := range pruned {
			v.neighbors[nid] = append(v.neighbors[nid], id)
			l := len(v.neighbors[nid])
			if l < int(GraphSlackFactor*float32(v.R)) {
				dummyCandidates := make([]Item, l)
				for j, desID := range v.neighbors[nid] {
					dist := v.metric.Similarity(v.vectors[nid], v.vectors[desID])
					dummyCandidates[j] = Item{ID: desID, Distance: dist}
				}
				npruned := v.prune(nid, v.vectors[nid], dummyCandidates)
				v.neighbors[nid] = npruned
			}
		}
	}

	for i := range vectorsOrder {
		id := nstart + vectorsOrder[i]
		l := len(v.neighbors[id])
		if l > v.R {
			dummyCandidates := make([]Item, l)
			for j, desID := range v.neighbors[id] {
				dist := v.metric.Similarity(v.vectors[id], v.vectors[desID])
				dummyCandidates[j] = Item{ID: desID, Distance: dist}
			}
			npruned := v.prune(id, v.vectors[id], dummyCandidates)
			v.neighbors[id] = npruned
		}
	}
}

func (v *Vamana) Search(query []float32, k int) ([]float32, []int) {
	bestLNodes := NewMaxHeapQ(v.L)
	resultHeap := NewMaxHeapQ(v.L)
	visitedB := NewBitSet(len(v.vectors))

	bestLNodes.Push(v.start, v.metric.Similarity(query, v.vectors[v.start]))
	visitedB.Set(v.start)
	expanded := make(map[int]bool)
	expandedCount := 0

	for expandedCount < bestLNodes.Len() {
		clst := bestLNodes.PopMin()
		resultHeap.Push(clst.ID, clst.Distance)
		expandedCount++
		expanded[clst.ID] = true

		for _, nid := range v.neighbors[clst.ID] {
			if visitedB.IsSet(nid) {
				continue
			}
			visitedB.Set(nid)
			dist := v.metric.Similarity(query, v.vectors[nid])
			if bestLNodes.Len() < v.L || dist < bestLNodes.Top().Distance {
				bestLNodes.Push(nid, dist)
				if bestLNodes.Len() >= v.L {
					item := bestLNodes.Pop()
					if expanded[item.ID] {
						expandedCount--
					}
				}
			}
		}
	}

	for resultHeap.Len() > k {
		resultHeap.Pop()
	}

	return resultHeap.PopAll()
}

func (v *Vamana) findMedoid(dataset [][]float32) int {
	dim := len(dataset[0])
	centroid := make([]float32, dim)

	for _, vector := range dataset {
		for i := 0; i < dim; i++ {
			centroid[i] += vector[i]
		}
	}
	for i := 0; i < dim; i++ {
		centroid[i] /= float32(len(dataset))
	}

	clst := 0
	minDist := v.metric.Similarity(centroid, dataset[0])
	for i, vector := range dataset {
		dist := v.metric.Similarity(centroid, vector)
		if dist < minDist {
			clst = i
			minDist = dist
		}
	}
	return clst
}

func (v *Vamana) candidatesToAddNew(id int, query []float32) []Item {
	expandedNodes := make([]Item, 0, int(1.05*GraphSlackFactor*float32(v.L)))
	bestLNodes := NewNeighborPriorityQueue(v.L)
	visitedB := NewBitSet(len(v.vectors))

	dist := v.metric.Similarity(query, v.vectors[v.start])
	bestLNodes.Insert(Neighbor{id: v.start, distance: dist, expanded: false})
	visitedB.Set(v.start)

	for bestLNodes.HasUnexpandedNode() {
		clst := bestLNodes.ClosestUnexpanded()
		expandedNodes = append(expandedNodes, Item{ID: clst.id, Distance: clst.distance})

		for _, nid := range v.neighbors[clst.id] {
			if visitedB.IsSet(nid) {
				continue
			}
			visitedB.Set(nid)
			dist := v.metric.Similarity(query, v.vectors[nid])
			bestLNodes.Insert(Neighbor{id: nid, distance: dist, expanded: false})
		}
	}
	return expandedNodes
}

func (v *Vamana) candidatesToAdd(id int, query []float32) []Item {
	bestLNodes := NewMaxHeapQ(v.L)
	candidates := make([]Item, 0, v.L)
	visited := make(map[int]bool)

	dist := v.metric.Similarity(query, v.vectors[v.start])
	bestLNodes.Push(v.start, dist)
	visited[v.start] = true
	expanded := make(map[int]bool)
	expandedCount := 0

	// for bestLNodes.Len() > 0 {
	for expandedCount < bestLNodes.Len() {
		clst := bestLNodes.PopMin()
		expandedCount++
		expanded[clst.ID] = true
		if clst.ID != id {
			candidates = append(candidates, clst)
		}

		for _, nid := range v.neighbors[clst.ID] {
			if visited[nid] {
				continue
			}
			visited[nid] = true
			dist := v.metric.Similarity(query, v.vectors[nid])
			if bestLNodes.Len() < v.L || dist < bestLNodes.Top().Distance {
				bestLNodes.Push(nid, dist)
				if bestLNodes.Len() >= v.L {
					item := bestLNodes.Pop()
					if expanded[item.ID] {
						expandedCount--
					}
				}
			}
		}
	}

	return candidates
}

func (v *Vamana) prune(id int, query []float32, candidates []Item) []int {
	var curAlpha float32 = 1.0
	result := make([]int, 0, v.R)
	occludeFactor := make([]float32, len(candidates))

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Distance < candidates[j].Distance
	})

	for curAlpha <= v.alpha {
		for i := 0; i < len(candidates) && len(result) < v.R; i++ {
			if occludeFactor[i] > curAlpha {
				continue
			}
			occludeFactor[i] = math.MaxFloat32

			if id != candidates[i].ID {
				result = append(result, candidates[i].ID)
			}

			for j := i + 1; j < len(candidates); j++ {
				if occludeFactor[j] > v.alpha {
					continue
				}
				id1, id2 := candidates[i].ID, candidates[j].ID
				dist := v.metric.Similarity(v.vectors[id1], v.vectors[id2])
				if dist == 0 {
					occludeFactor[j] = math.MaxFloat32
				} else {
					occludeFactor[j] = max(occludeFactor[j], candidates[j].Distance/dist)
				}
			}
		}
		curAlpha *= 1.2
	}

	limit := min(v.R, len(result))
	return result[:limit]
}

type Neighbor struct {
	id       int
	distance float32
	expanded bool // Assuming expanded is a boolean flag
}

type NeighborPriorityQueue struct {
	size     int
	capacity int
	cur      int
	data     []Neighbor
}

func NewNeighborPriorityQueue(capacity int) *NeighborPriorityQueue {
	return &NeighborPriorityQueue{
		size:     0,
		capacity: capacity,
		cur:      0,
		data:     make([]Neighbor, capacity+1),
	}
}

func (q *NeighborPriorityQueue) Insert(nbr Neighbor) {
	if q.size == q.capacity && q.data[q.size-1].distance < nbr.distance {
		return
	}

	lo, hi := 0, q.size
	for lo < hi {
		mid := (lo + hi) >> 1
		if nbr.distance < q.data[mid].distance {
			hi = mid
		} else if q.data[mid].id == nbr.id {
			return
		} else {
			lo = mid + 1
		}
	}

	if lo < q.capacity {
		copy(q.data[lo+1:], q.data[lo:q.size])
	}
	q.data[lo] = nbr
	if q.size < q.capacity {
		q.size++
	}
	if lo < q.cur {
		q.cur = lo
	}
}

func (q *NeighborPriorityQueue) ClosestUnexpanded() Neighbor {
	q.data[q.cur].expanded = true
	pre := q.cur
	for q.cur < q.size && q.data[q.cur].expanded {
		q.cur++
	}
	return q.data[pre]
}

func (q *NeighborPriorityQueue) HasUnexpandedNode() bool {
	return q.cur < q.size
}

func (q *NeighborPriorityQueue) Size() int {
	return q.size
}

func (q *NeighborPriorityQueue) Capacity() int {
	return q.capacity
}

func (q *NeighborPriorityQueue) Clear() {
	q.size = 0
	q.cur = 0
}
