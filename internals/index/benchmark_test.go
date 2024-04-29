package index

import (
	"testing"
)

func BenchmarkFlatIndex(b *testing.B) {
	baseData, _, _, queries := loadSift()

	b.Run("Index Construction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := NewFlat(128, nil)
			f.AddBatch(baseData)
		}
	})

	b.Run("Query", func(b *testing.B) {
		f := NewFlat(128, nil)
		f.AddBatch(baseData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f.SearchMany(queries, 100)
		}
	})
}

func BenchmarkFlatPQIndex(b *testing.B) {
	baseData, _, learn, queries := loadSift()

	b.Run("Index Construction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := NewFlatPQ(128, nil)
			f.Train(learn, NewPQ(128, 8))
			f.AddBatch(baseData)
		}
	})

	b.Run("Query", func(b *testing.B) {
		f := NewFlatPQ(128, nil)
		f.Train(learn, NewPQ(128, 8))
		f.AddBatch(baseData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, query := range queries {
				f.Search(query, 10)
			}
		}
	})
}

func BenchmarkIVF(b *testing.B) {
	baseData, _, learn, queries := loadSift()
	_ = queries

	b.Run("Index Construction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ivf := NewIVF(128, 100, 3, nil)
			ivf.Train(learn)
			ivf.AddBatch(baseData)
		}
	})

	b.Run("Query", func(b *testing.B) {
		ivf := NewIVF(128, 100, 3, nil)
		ivf.Train(learn)
		ivf.AddBatch(baseData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ivf.SearchMany(queries, 100)
			// for _, query := range queries {
			// 	ivf.Search(query, 10)
			// }
		}
	})
}

func BenchmarkIVFPQ(b *testing.B) {
	baseData, _, learn, queries := loadSift()
	_ = queries

	b.Run("Index Construction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ivf := NewIVFPQ(128, 100, 3, nil)
			ivf.Train(learn, NewPQ(128, 8))
			ivf.AddBatch(baseData)
		}
	})

	b.Run("Query", func(b *testing.B) {
		ivf := NewIVFPQ(128, 100, 1, nil)
		ivf.Train(learn, NewPQ(128, 8))
		ivf.AddBatch(baseData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, query := range queries {
				ivf.Search(query, 10)
			}
		}
	})
}

func BenchmarkHNSW(b *testing.B) {
	baseData, _, _, queries := loadSift()

	b.Run("Index Construction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := NewHNSW(10, 20, 20, nil)
			f.AddBatch(baseData)
		}
	})

	b.Run("Query", func(b *testing.B) {
		f := NewFlat(128, nil)
		f.AddBatch(baseData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, query := range queries {
				f.Search(query, 10)
			}
		}
	})
}

func BenchmarkVamana(b *testing.B) {
	baseData, _, _, queries := loadSift()

	b.Run("Index Construction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f := NewVamana(10, 120, 1.3, nil)
			f.AddBatch(baseData)
		}
	})

	b.Run("Query", func(b *testing.B) {
		f := NewVamana(10, 120, 1.3, nil)
		f.AddBatch(baseData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, query := range queries {
				f.Search(query, 10)
			}
		}
	})
}
