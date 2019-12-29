package main

func transitionMatrix(tagsCount map[string]int, sentencesTags [][]string) (map[string]int, [][]float32) {

	inx := 0
	tagIndex := make(map[string]int)

	for tag := range tagsCount {
		tagIndex[tag] = inx
		inx++
	}

	transitionCount := make([][]int, len(tagIndex))
	for inx := range transitionCount {
		transitionCount[inx] = make([]int, len(tagIndex))
	}

	for _, tags := range sentencesTags {
		for inx := 1; inx < len(tags); inx++ {
			preTag := tags[inx-1]
			tag := tags[inx]

			preTagIndex := tagIndex[preTag]
			tagIndex := tagIndex[tag]

			transitionCount[preTagIndex][tagIndex]++
		}
	}

	// init transition matrix
	transitionMatrix := make([][]float32, len(tagIndex))
	for inx := range transitionMatrix {
		transitionMatrix[inx] = make([]float32, len(tagIndex))
	}

	for inx, data := range transitionCount {

		sum := 0

		for _, val := range data {
			sum += val
		}

		for iny, val := range data {
			transitionMatrix[inx][iny] = float32(val) / float32(sum)
		}
	}

	return tagIndex, transitionMatrix
}
