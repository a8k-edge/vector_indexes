package vops

import (
	"testing"

	"mvdb/internals/utils"

	"github.com/stretchr/testify/assert"
)

func TestL2Dist(t *testing.T) {
	tests := []struct {
		a []float32
		b []float32
		d float32
	}{
		{[]float32{1, 1}, []float32{2, 2}, 1.41},
		{[]float32{2, 2}, []float32{1, 1}, 1.41},
		{[]float32{1, 1}, []float32{1, 1}, 0},
		{[]float32{1, 1}, []float32{3, 3}, 2.83},
	}
	provider := L2Distance{}

	for _, tc := range tests {
		d := provider.Similarity(tc.a, tc.b)
		assert.InDelta(t, tc.d, d, 0.02)
	}
}

func BenchmarkL2Dist(b *testing.B) {
	dim := 128
	x, y := utils.RVec(dim), utils.RVec(dim)
	provider := L2Distance{}

	for i := 0; i < b.N; i++ {
		provider.Similarity(x, y)
	}
}

func TestL2SqrDist(t *testing.T) {
	tests := []struct {
		a []float32
		b []float32
		d float32
	}{
		{[]float32{1, 1}, []float32{2, 2}, 2},
		{[]float32{2, 2}, []float32{1, 1}, 2},
		{[]float32{1, 1}, []float32{1, 1}, 0},
		{[]float32{1, 1}, []float32{3, 3}, 8},
	}
	provider := L2SqrDistance{}

	for _, tc := range tests {
		d := provider.Similarity(tc.a, tc.b)
		assert.InDelta(t, tc.d, d, 0.02)
	}
}

func BenchmarkL2SqrDist(b *testing.B) {
	dim := 128
	x, y := utils.RVec(dim), utils.RVec(dim)
	provider := L2SqrDistance{}

	for i := 0; i < b.N; i++ {
		provider.Similarity(x, y)
	}
}

func TestCosineSim(t *testing.T) {
	tests := []struct {
		a []float32
		b []float32
		d float32
	}{
		{[]float32{2, 3, 4}, []float32{4, 6, 8}, 1},
		{[]float32{1, -2}, []float32{-1, 2}, -1},
		{[]float32{3, 0}, []float32{0, 4}, 0},
		{[]float32{1, 2, 3}, []float32{4, 5, 6}, 0.96},
	}
	provider := CosineSimilarity{}

	for _, tc := range tests {
		d := provider.Similarity(tc.a, tc.b)
		assert.InDelta(t, tc.d, d, 0.02)
	}
}

func BenchmarkCosineSim(b *testing.B) {
	dim := 128
	x, y := utils.RVec(dim), utils.RVec(dim)
	provider := CosineSimilarity{}

	for i := 0; i < b.N; i++ {
		provider.Similarity(x, y)
	}
}
