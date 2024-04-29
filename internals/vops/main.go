package vops

import (
	"math"
	"unsafe"

	"golang.org/x/sys/cpu"
)

type Provider interface {
	Similarity(a, b []float32) float32
}

type L2SqrDistance struct{}

func (d *L2SqrDistance) Similarity(a, b []float32) float32 {
	// TODO: set function at init state
	if cpu.X86.HasAVX2 {
		return L2AVX256(a, b)
	}

	var r float32
	for i := 0; i < len(a); i++ {
		diff := b[i] - a[i]
		r += diff * diff
	}
	return r
}

type L2Distance struct{}

func (d *L2Distance) Similarity(a, b []float32) float32 {
	var r float32
	for i := 0; i < len(a); i++ {
		diff := b[i] - a[i]
		r += diff * diff
	}
	return float32(math.Sqrt(float64(r)))
}

type CosineSimilarity struct{}

func (c *CosineSimilarity) Similarity(a, b []float32) float32 {
	return dotProd(a, b)
	// aNorm, bNorm := l2norm(a), l2norm(b)
	// if aNorm == 0 || bNorm == 0 {
	// 	return 0
	// }
	// return dotProd(a, b) / (aNorm * bNorm)
}

func NormalizeV(vector []float32) []float32 {
	nv := make([]float32, len(vector))
	copy(nv, vector)

	norm := l2norm(nv)
	if norm == 0 {
		return nv
	}
	for i := range vector {
		nv[i] = nv[i] / norm
	}
	return nv
}

func dotProd(a, b []float32) float32 {
	var r float32
	for i := 0; i < len(a); i++ {
		r += a[i] * b[i]
	}
	return r
}

func l2norm(a []float32) float32 {
	var r float32

	for i := 0; i < len(a); i++ {
		r += a[i] * a[i]
	}

	return float32(math.Sqrt(float64(r)))
}

func l2_256(a, b, res, len unsafe.Pointer)

// TODO: prefetch
func L2AVX256(x []float32, y []float32) float32 {
	var res float32

	l := len(x)
	l2_256(
		unsafe.Pointer(unsafe.SliceData(x)),
		unsafe.Pointer(unsafe.SliceData(y)),
		unsafe.Pointer(&res),
		unsafe.Pointer(&l))

	return res
}
