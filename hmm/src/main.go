package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type matrixList map[string]float64
type matrix map[string]matrixList

func transitionMatrixMap(tagsCount map[string]int, sentencesTags [][]string) matrix {

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

	// init transition matrix
	transitionMatrix := make(matrix)
	for tag := range tagsCount {
		transitionMatrix[tag] = make(matrixList)
	}

	for tag, data := range transitionCount {

		sum := 0

		for _, val := range data {
			sum += val
		}

		for tag2, val := range data {
			transitionMatrix[tag][tag2] = float64(val) / float64(sum)
		}
	}

	return transitionMatrix
}

func emissionMatrixMap(sentencesWords [][]string, sentencesTags [][]string) matrix {

	emissionCount := make(map[string]map[string]int)

	for inx, sentence := range sentencesWords {
		for iny := range sentence {

			word := sentencesWords[inx][iny]
			tag := sentencesTags[inx][iny]

			if emissionCount[tag] == nil {
				emissionCount[tag] = make(map[string]int)
			}

			emissionCount[tag][word]++
		}
	}

	emissionMatrix := make(matrix)

	for tag, words := range emissionCount {

		emissionMatrix[tag] = make(matrixList)

		wordsSum := 0

		for _, count := range words {
			wordsSum += count
		}

		for word, count := range words {
			emissionMatrix[tag][word] = float64(count) / float64(wordsSum)
		}
	}

	return emissionMatrix
}

type hmm struct {
	priorProbabilities map[string]float64
	emissionsMatrix    matrix
	transitionMatrix   matrix
}

func (model hmm) forwardAlgorithm(phrase []string) {

	// phrase --> Observation
	startWord := phrase[0]

	initResults := make(matrixList)

	for s := range model.transitionMatrix {
		initResults[s] = model.priorProbabilities[s] * model.emissionsMatrix[s][startWord]
	}

	fmt.Println("startWord", startWord)
	fmt.Println("initResults", initResults)

	aResults := make([]matrixList, 1)
	aResults[0] = initResults

	for k := 0; k < len(phrase)-1; k++ {

		result := make(matrixList)

		for s := range model.transitionMatrix {

			akSum := 0.0

			for q := range model.transitionMatrix {
				akSum += aResults[k][q] * model.transitionMatrix[q][s]
			}

			eVal := model.emissionsMatrix[s][phrase[k+1]]
			result[s] = eVal * akSum
		}

		aResults = append(aResults, result)
	}

	fmt.Println("aResults", aResults)

	for inx, word := range phrase {

		result := aResults[inx]

		bestTag := ""
		bestScore := 0.0

		for tag, score := range result {
			if score > bestScore {
				bestScore = score
				bestTag = tag
			}
		}

		fmt.Println(word, bestTag)
	}

	aggregate := 0.0
	for _, a := range aResults[len(phrase)-1] {
		aggregate += a
	}
	fmt.Println("aggregate", aggregate)
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

	priorProbabilities := make(matrixList)
	for tag, count := range startTagsCount {
		priorProbabilities[tag] = float64(count) / float64(len(sentencesTags))
	}

	fmt.Println("priorProbabilities", priorProbabilities)

	emissionsMatrix := emissionMatrixMap(sentencesWords, sentencesTags)
	transitionMatrix := transitionMatrixMap(tagsCount, sentencesTags)

	// fmt.Println("emissionsMatrix", emissionsMatrix)
	// fmt.Println("transitionMatrix", transitionMatrix)

	// phrase := "Pro Monat sind dafür 2,99 Euro fällig ."
	phrase := "Dazu kommen zehn statt bisher fünf E-Mail-Adressen sowie zehn MByte Webspace ."
	phraseParts := strings.Split(phrase, " ")

	model := hmm{
		priorProbabilities: priorProbabilities,
		emissionsMatrix:    emissionsMatrix,
		transitionMatrix:   transitionMatrix,
	}

	model.forwardAlgorithm(phraseParts)
}
