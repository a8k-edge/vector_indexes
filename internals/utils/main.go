package utils

import (
	"fmt"
	"math/rand"
)

func RVec(l int) []float32 {
	v := make([]float32, l)
	for i := 0; i < l; i++ {
		v[i] = rand.Float32() * 100
	}
	return v
}

func ByteCountSI(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
