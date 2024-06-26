package index

import "mvdb/internals/utils"

func loadSift() (base [][]float32, truth [][]int, learn [][]float32, queries [][]float32) {
	base, err := utils.FVecsRead("dataset/siftsmall/siftsmall_base.fvecs")
	if err != nil {
		panic(err)
	}

	truth, err = utils.IVecsRead("dataset/siftsmall/siftsmall_groundtruth.ivecs")
	if err != nil {
		panic(err)
	}

	learn, err = utils.FVecsRead("dataset/siftsmall/siftsmall_learn.fvecs")
	if err != nil {
		panic(err)
	}

	queries, err = utils.FVecsRead("dataset/siftsmall/siftsmall_query.fvecs")
	if err != nil {
		panic(err)
	}

	return
}

func LoadSift() (base [][]float32, truth [][]int, learn [][]float32, queries [][]float32) {
	return loadSift()
}
