package utils

import (
	"encoding/binary"
	"math"
	"os"
)

func FVecsRead(fname string) ([][]float32, error) {
	intVectors, err := IVecsRead(fname)
	if err != nil {
		return nil, err
	}

	floatVectors := make([][]float32, len(intVectors))
	for i, vector := range intVectors {
		floatVectors[i] = make([]float32, len(vector))
		for j, value := range vector {
			floatVectors[i][j] = math.Float32frombits(uint32(value))
		}
	}

	return floatVectors, nil
}

func IVecsRead(fname string) ([][]int, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	data := make([]byte, fileSize)
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	a := make([]int32, len(data)/4)
	for i := 0; i < len(data); i += 4 {
		a[i/4] = int32(binary.LittleEndian.Uint32(data[i : i+4]))
	}

	d := int(a[0])

	reshapedList := make([][]int, 0)
	for i := 0; i < len(a); i += d + 1 {
		reshapedList = append(reshapedList, intSlice(a[i+1:i+d+1]))
	}

	return reshapedList, nil
}

func intSlice(slice []int32) []int {
	result := make([]int, len(slice))
	for i, v := range slice {
		result[i] = int(v)
	}
	return result
}

func IntersectionCount(a, b []int) int {
	c := 0
	m := make(map[int]bool)

	for _, v := range a {
		m[v] = true
	}
	for _, v := range b {
		if m[v] {
			c++
		}
	}

	return c
}

func LoadSift() (base [][]float32, truth [][]int, learn [][]float32, queries [][]float32) {
	base, err := FVecsRead("internals/index/dataset/siftsmall/siftsmall_base.fvecs")
	if err != nil {
		panic(err)
	}

	truth, err = IVecsRead("internals/index/dataset/siftsmall/siftsmall_groundtruth.ivecs")
	if err != nil {
		panic(err)
	}

	learn, err = FVecsRead("internals/index/dataset/siftsmall/siftsmall_learn.fvecs")
	if err != nil {
		panic(err)
	}

	queries, err = FVecsRead("internals/index/dataset/siftsmall/siftsmall_query.fvecs")
	if err != nil {
		panic(err)
	}

	return
}
