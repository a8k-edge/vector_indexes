package vops

import (
	"testing"

	"mvdb/internals/utils"

	"golang.org/x/sys/cpu"
)

func BenchmarkL2(b *testing.B) {
	base, _, _, queries := utils.LoadSift()
	_ = base
	_ = queries

	dist := L2SqrDistance{}

	b.Run("With AVX2", func(b *testing.B) {
		for _, query := range queries {
			for _, v := range base {
				dist.Similarity(query, v)
			}
		}
	})

	cpu.X86.HasAVX2 = false
	b.Run("Run baseline", func(b *testing.B) {
		for _, query := range queries {
			for _, v := range base {
				dist.Similarity(query, v)
			}
		}
	})
}
