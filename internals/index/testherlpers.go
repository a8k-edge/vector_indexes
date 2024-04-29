package index

import "mvdb/internals/utils"

func loadSift() (base [][]float32, truth [][]int, learn [][]float32, queries [][]float32) {
	base, err := utils.FVecsRead("dataset/sift/sift_base.fvecs")
	if err != nil {
		panic(err)
	}

	truth, err = utils.IVecsRead("dataset/sift/sift_groundtruth.ivecs")
	if err != nil {
		panic(err)
	}

	learn, err = utils.FVecsRead("dataset/sift/sift_learn.fvecs")
	if err != nil {
		panic(err)
	}

	queries, err = utils.FVecsRead("dataset/sift/sift_query.fvecs")
	if err != nil {
		panic(err)
	}

	return
}

func LoadSift() (base [][]float32, truth [][]int, learn [][]float32, queries [][]float32) {
	return loadSift()
}
