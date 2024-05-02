package num

import (
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

func RandMatrix(dim int) [][]float32 {
	m := make([][]float32, dim)
	for i := 0; i < dim; i++ {
		m[i] = RandVec(dim)
	}
	return m
}

func RandVec(dim int) []float32 {
	v := make([]float32, dim)
	for i := 0; i < dim; i++ {
		v[i] = rand.Float32() * float32(rand.Intn(300))
	}
	return v
}

func RandVec64(dim int) []float64 {
	v := make([]float64, dim)
	for i := 0; i < dim; i++ {
		// TODO: consider NormFloat64
		v[i] = rand.Float64() * float64(rand.Intn(300))
	}
	return v
}

func RandRotationFMatrix(dim int) []float32 {
	A := mat.NewDense(dim, dim, RandVec64(dim*dim))
	var qr mat.QR
	qr.Factorize(A)
	Q := mat.NewDense(dim, dim, nil)
	qr.QTo(Q)

	return toFloat32(Q.RawMatrix().Data)
}

func toFloat32(s []float64) []float32 {
	r := make([]float32, len(s))
	for i, val := range s {
		r[i] = float32(val)
	}
	return r
}
