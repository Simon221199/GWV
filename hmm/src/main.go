package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func transitionMatrixMap(tagsCount map[string]int, sentencesTags [][]string) map[string]map[string]float32 {

	transitionCount := make(map[string]map[string]int)
	for tag := range tagsCount {
		transitionCount[tag] = make(map[string]int)
	}

	for _, tags := range sentencesTags {
		for inx := 1; inx < len(tags); inx++ {
			preTag := tags[inx-1]
			tag := tags[inx]

			transitionCount[preTag][tag]++
		}
	}

	fmt.Println("transitionCount", transitionCount)

	// init transition matrix
	transitionMatrix := make(map[string]map[string]float32)
	for tag := range tagsCount {
		transitionMatrix[tag] = make(map[string]float32)
	}

	for tag, data := range transitionCount {

		sum := 0

		for _, val := range data {
			sum += val
		}

		for tag2, val := range data {
			transitionMatrix[tag][tag2] = float32(val) / float32(sum)
		}
	}

	fmt.Println("transitionMatrix", transitionMatrix)

	return transitionMatrix
}

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

func emissionMatrixMap(tagsCount map[string]int, sentencesWords [][]string, sentencesTags [][]string) map[string]map[string]float32 {

	emissionCount := make(map[string]map[string]int)

	for inx, sentence := range sentencesWords {
		for iny := range sentence {

			word := sentencesWords[inx][iny]
			tag := sentencesTags[inx][iny]

			if emissionCount[word] == nil {
				emissionCount[word] = make(map[string]int)
			}

			emissionCount[word][tag]++
		}
	}

	emissionMatrix := make(map[string]map[string]float32)

	for word, tags := range emissionCount {

		emissionMatrix[word] = make(map[string]float32)

		tagsSum := 0

		for _, count := range tags {
			tagsSum += count
		}

		for tag, count := range tags {
			emissionMatrix[word][tag] = float32(count) / float32(tagsSum)
		}
	}

	return emissionMatrix
}

func main() {
	fmt.Println("hallo")

	content, err := ioutil.ReadFile("hdt-1-10000-train.tags")
	if err != nil {
		panic(err)
	}

	text := strings.TrimSpace(string(content))
	sentences := strings.Split(text, "\n\n")

	tagsCount := make(map[string]int)

	sentencesWords := make([][]string, len(sentences))
	sentencesTags := make([][]string, len(sentences))

	for inx, sentence := range sentences {

		lines := strings.Split(sentence, "\n")

		sentencesWords[inx] = make([]string, len(lines))
		sentencesTags[inx] = make([]string, len(lines))

		for iny, line := range lines {

			wordTag := strings.Split(line, "\t")

			sentencesWords[inx][iny] = wordTag[0]
			sentencesTags[inx][iny] = wordTag[1]

			tagsCount[wordTag[1]]++
		}
	}

	fmt.Println("tagsCount", tagsCount)

	startTagsCount := make(map[string]int)
	for inx := range sentencesTags {
		tag := sentencesTags[inx][0]
		startTagsCount[tag]++
	}

	priorProbabilities := make(map[string]float32)
	for tag, count := range startTagsCount {
		priorProbabilities[tag] = float32(count) / float32(len(sentencesTags))
	}

	fmt.Println("priorProbabilities", priorProbabilities)

	emissionsMatrix := emissionMatrixMap(tagsCount, sentencesWords, sentencesTags)
	// transitionMatrix := transitionMatrixMap(tagsCount, sentencesTags)

	fmt.Println("emissionsMatrix", emissionsMatrix)
	// fmt.Println("transitionMatrix", transitionMatrix)

	// tags, matrix := transitionMatrix(tagsCount, sentencesTags)
	// fmt.Println("tags", tags)
	// fmt.Println("matrix", matrix[ tags[ "$(" ] ][ tags[ "$(" ] ])

	// matrix := transitionMatrixMap(tagsCount, sentencesTags)
	// fmt.Println("matrix", matrix)
}
