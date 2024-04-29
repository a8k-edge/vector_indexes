package index

import (
	"math/rand"

	"mvdb/internals/vops"
)

func LloydKmeansBySection(
	dim int,
	sectionIndex int,
	data [][]float32,
	k int,
	maxIter int,
	metric vops.Provider,
) ([][]float32, []int) {
	l := len(data)
	if l < k {
		panic("Kmeans: Not enough data to fit")
	}
	start, end := sectionIndex*dim, sectionIndex*dim+dim

	// kmeans++ initialization
	distances := make([]float32, l)
	centroids := make([][]float32, k)
	centroids[0] = data[rand.Intn(l)][start:end]

	for i := 1; i < k; i++ {
		var totalDistances float32
		for vid, vector := range data {
			distance := metric.Similarity(vector[start:end], centroids[i-1])
			if distance < distances[vid] {
				distances[vid] = distance
			}
			totalDistances += distances[vid]
		}

		target := rand.Float32() * totalDistances
		for vid, distance := range distances {
			target -= distance
			if target <= 0 {
				centroids[i] = data[vid]
				break
			}
		}
	}

	labels := make([]int, l)
	changes := 0
	clustCount := make([]int, k)
	newCentroids := make([][]float32, k)
	strictConvergence := false

	// converge
	for iter := 0; iter < maxIter; iter++ {
		clustCount = make([]int, k)

		for vid, vector := range data {
			cid := nearest(vector[start:end], centroids, metric)
			clustCount[cid] += 1
			if cid != labels[vid] {
				labels[vid] = cid
				changes += 1
			}
		}

		// resort if empty
		// TODO: assign new centroid to furthest vector(dist to it centroid)
		// for i := empty centroids
		// labels[]
		for cid := range centroids {
			if clustCount[cid] != 0 {
				continue
			}
			var vid int
			for {
				vid = rand.Intn(l)
				if clustCount[labels[vid]] > 1 {
					clustCount[labels[vid]] -= 1
					break
				}
			}

			labels[vid] = cid
			clustCount[cid] += 1
			changes = l
		}

		for i := 0; i < k; i++ {
			newCentroids[i] = make([]float32, dim)
		}
		for vid, vector := range data {
			cid := labels[vid]
			for j := 0; j < dim; j++ {
				newCentroids[cid][j] += vector[start+j]
			}
		}
		for i := 0; i < k; i++ {
			for j := 0; j < dim; j++ {
				newCentroids[i][j] /= float32(clustCount[i])
			}
		}

		centroids, newCentroids = newCentroids, centroids

		// TODO: consider centroids shift
		if changes > int(float32(l)*0.01) {
			strictConvergence = true
			break
		}
	}

	if strictConvergence {
		for vid, vector := range data {
			cid := nearest(vector[start:end], centroids, metric)
			clustCount[cid] += 1
			if cid != labels[vid] {
				labels[vid] = cid
				changes += 1
			}
		}
	}

	return centroids, labels
}

func LloydKmeans(data [][]float32, k int, maxIter int, metric vops.Provider) ([][]float32, []int) {
	dim := len(data[0])
	return LloydKmeansBySection(dim, 0, data, k, maxIter, metric)
}

func nearest(vector []float32, centroids [][]float32, metric vops.Provider) int {
	nearest := 0
	minDist := metric.Similarity(vector, centroids[0])
	for j := 1; j < len(centroids); j++ {
		dist := metric.Similarity(vector, centroids[j])
		if dist < minDist {
			minDist = dist
			nearest = j
		}
	}
	return nearest
}
