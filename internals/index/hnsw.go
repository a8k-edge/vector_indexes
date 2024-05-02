package index

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"

	"mvdb/internals/vops"
)

type HNSW struct {
	vectors [][]float32

	neighbors []map[int][]int
	levels    []int
	probs     []float64

	M              int
	efSearch       int
	efConstruction int
	Rng            *rand.Rand

	ep     int
	maxLvl int
	metric vops.Provider
}

func NewHNSW(M int, ef int, efConstruction int, metric vops.Provider) *HNSW {
	if metric == nil {
		metric = &vops.L2SqrDistance{}
	}
	levelMult := 1 / math.Log(float64(M))
	probs := make([]float64, 0)
	for level := 0; ; level++ {
		prob := math.Exp(-float64(level)/levelMult) * (1 - math.Exp(-1/levelMult))
		if prob < 1e-9 {
			break
		}
		probs = append(probs, prob)
	}

	return &HNSW{
		neighbors: make([]map[int][]int, 0),
		probs:     probs,

		M:              M,
		efSearch:       ef,
		efConstruction: efConstruction,
		Rng:            rand.New(rand.NewSource(12345)),

		ep:     -1,
		maxLvl: -1,
		metric: metric,
	}
}

func (h *HNSW) PrintLvlDistribution() {
	dist := make(map[int]int)

	for _, lvl := range h.levels {
		dist[lvl]++
	}

	fmt.Println(dist)
}

func (h *HNSW) SetEf(ef int) {
	h.efSearch = ef
}

func (h *HNSW) Print() {
	fmt.Println("neighbors", h.neighbors)
	fmt.Println("vectors", h.vectors)
	fmt.Println("levels", h.levels)
	fmt.Println("ep", h.ep)
	fmt.Println("maxLvl", h.maxLvl)
}

func (h *HNSW) Size() uint64 {
	var size uintptr

	size += reflect.TypeOf(h.vectors).Size()
	for _, v := range h.vectors {
		size += reflect.TypeOf(v).Size()
		for _, f := range v {
			size += reflect.TypeOf(f).Size()
		}
	}

	// size += reflect.TypeOf(h.neighbors).Size()
	// for _, m := range h.neighbors {
	// 	size += reflect.TypeOf(m).Size()
	// 	for _, s := range m {
	// 		size += reflect.TypeOf(s).Size()
	// 		for _, i := range s {
	// 			size += reflect.TypeOf(i).Size()
	// 		}
	// 	}
	// }

	size += reflect.TypeOf(h.levels).Size()
	for _, l := range h.levels {
		size += reflect.TypeOf(l).Size()
	}

	size += reflect.TypeOf(h.probs).Size()
	for _, p := range h.probs {
		size += reflect.TypeOf(p).Size()
	}

	size += reflect.TypeOf(h.M).Size()
	size += reflect.TypeOf(h.efSearch).Size()
	size += reflect.TypeOf(h.efConstruction).Size()
	size += reflect.TypeOf(h.Rng).Size()
	size += reflect.TypeOf(h.ep).Size()
	size += reflect.TypeOf(h.maxLvl).Size()

	return uint64(size)
}

func (h *HNSW) AddBatch(vectors [][]float32) {
	ntotal := len(h.vectors)
	for i := 0; i < len(vectors); i++ {
		lvl := h.randomLevel()
		h.levels = append(h.levels, lvl)
		h.vectors = append(h.vectors, vectors[i])
		neighbors := make(map[int][]int)
		for i := 0; i <= lvl; i++ {
			neighbors[i] = make([]int, 0)
		}
		h.neighbors = append(h.neighbors, neighbors)
	}

	// TODO: paralellize
	for i, vector := range vectors {
		id := ntotal + i
		lvl := h.levels[id]

		nearest := h.ep
		if nearest == -1 {
			h.maxLvl = lvl
			h.ep = id
		}
		if nearest < 0 {
			continue
		}

		curLvl := h.maxLvl
		nearestD := h.metric.Similarity(vector, h.vectors[nearest])

		for ; curLvl > lvl; curLvl-- {
			nearest, nearestD = h.greedyNearest(vector, curLvl, nearest, nearestD)
		}

		for ; curLvl >= 0; curLvl-- {
			candidates := NewMinHeapQ(0)
			candidates.Push(nearest, nearestD)
			toAdd := h.searchNeighborToAdd(vector, nearest, nearestD, curLvl)

			m := h.M
			if curLvl == 0 {
				m *= 2
			}

			neighbors := h.shrinkNeighbors(toAdd, m)

			for _, nid := range neighbors {
				h.addLink(id, nid, curLvl, m)
			}
			for _, nid := range h.neighbors[id][curLvl] {
				h.addLink(nid, id, curLvl, m)
			}
		}
	}
}

func (h *HNSW) Search(query []float32, k int) ([]float32, []int) {
	nearest := h.ep
	nearestD := h.metric.Similarity(query, h.vectors[nearest])

	for lvl := h.maxLvl; lvl > 0; lvl-- {
		nearest, nearestD = h.greedyNearest(query, lvl, nearest, nearestD)
	}

	ef := max(h.efSearch, k)
	candidates := NewMinHeapQ(0)
	candidates.Push(nearest, nearestD)
	heapResult := h.searchFromCandidates(query, candidates, ef, 0)

	for heapResult.Len() > k {
		heapResult.Pop()
	}
	return heapResult.PopAll()
}

func (h *HNSW) greedyNearest(query []float32, lvl int, nearest int, nearestD float32) (int, float32) {
	for {
		prevNearest := nearest

		neighbors := h.neighbors[nearest][lvl]
		for _, id := range neighbors {
			d := h.metric.Similarity(query, h.vectors[id])
			if d < nearestD {
				nearest = id
				nearestD = d
			}
		}
		if prevNearest == nearest {
			break
		}
	}

	return nearest, nearestD
}

func (h *HNSW) searchFromCandidates(query []float32, candidates *HeapQ, ef int, lvl int) *HeapQ {
	visited := make(map[int]bool)
	futherQ := NewMaxHeapQ(0)
	for _, item := range candidates.items {
		futherQ.Push(item.ID, item.Distance)
		visited[item.ID] = true
	}

	for candidates.Len() > 0 {
		clstItem := candidates.Pop()
		v := clstItem.ID
		furthestItem := futherQ.Top()

		// TODO: What if candidates already the closests?
		if futherQ.Len() >= ef && furthestItem.Distance < clstItem.Distance {
			break
		}

		for _, child := range h.neighbors[v][lvl] {
			if visited[child] {
				continue
			}
			visited[child] = true

			dist := h.metric.Similarity(query, h.vectors[child])
			if dist < furthestItem.Distance || futherQ.Len() < ef {
				candidates.Push(child, dist)
				futherQ.Push(child, dist)
				if futherQ.Len() > ef {
					futherQ.Pop()
				}
			}

		}
	}

	return futherQ
}

func (h *HNSW) searchNeighborToAdd(query []float32, ep int, epDist float32, lvl int) *HeapQ {
	clstQ := NewMinHeapQ(h.efConstruction)
	futherQ := NewMaxHeapQ(h.efConstruction)
	visited := make(map[int]bool)
	clstQ.Push(ep, epDist)
	futherQ.Push(ep, epDist)
	visited[ep] = true

	for clstQ.Len() > 0 {
		clst := clstQ.Pop()
		if clst.Distance > futherQ.Top().Distance {
			break
		}
		for _, nid := range h.neighbors[clst.ID][lvl] {
			if visited[nid] {
				continue
			}
			visited[nid] = true

			dist := h.metric.Similarity(query, h.vectors[nid])
			if futherQ.Len() < h.efConstruction || futherQ.Top().Distance > dist {
				clstQ.Push(nid, dist)
				futherQ.Push(nid, dist)
				if futherQ.Len() > h.efConstruction {
					futherQ.Pop()
				}
			}
		}
	}

	return futherQ
}

func (h *HNSW) shrinkNeighbors(candidates *HeapQ, m int) []int {
	clstQ := NewMinHeapQ(candidates.Len())
	for candidates.Len() > 0 {
		item := candidates.Pop()
		clstQ.Push(item.ID, item.Distance)
	}
	result := []int{}

	for clstQ.Len() > 0 {
		item := clstQ.Pop()

		good := true
		for _, vid := range result {
			dist := h.metric.Similarity(h.vectors[vid], h.vectors[item.ID])
			if dist < item.Distance {
				good = false
				break
			}
		}

		if good {
			result = append(result, item.ID)
			if len(result) == m {
				break
			}
		}
	}
	return result
}

func (h *HNSW) addLink(src int, dest int, lvl int, m int) {
	if len(h.neighbors[src][lvl]) < m {
		h.neighbors[src][lvl] = append(h.neighbors[src][lvl], dest)
		return
	}

	neighbors := h.neighbors[src][lvl]
	futherQ := NewMaxHeapQ(len(neighbors))
	for _, id := range neighbors {
		futherQ.Push(id, h.metric.Similarity(h.vectors[src], h.vectors[id]))
	}
	newNeighbors := h.shrinkNeighbors(futherQ, m)
	h.neighbors[src][lvl] = newNeighbors
}

func (h *HNSW) randomLevel() int {
	f := h.Rng.Float64()
	for level := 0; level < len(h.probs); level++ {
		if f < h.probs[level] {
			return level
		}
		f -= h.probs[level]
	}
	return len(h.probs) - 1
}
