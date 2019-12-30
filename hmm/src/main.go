package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type float float64
type probabilityMap map[string]float
type matrix map[string]probabilityMap

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
		transitionMatrix[tag] = make(probabilityMap)
	}

	for tag, data := range transitionCount {

		sum := 0

		for _, val := range data {
			sum += val
		}

		for tag2, val := range data {
			transitionMatrix[tag][tag2] = float(val) / float(sum)
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

		emissionMatrix[tag] = make(probabilityMap)

		wordsSum := 0

		for _, count := range words {
			wordsSum += count
		}

		for word, count := range words {
			emissionMatrix[tag][word] = float(count) / float(wordsSum)
		}
	}

	return emissionMatrix
}

type hmm struct {
	words              map[string]bool
	tagProbability     probabilityMap
	priorProbabilities probabilityMap
	emissionsMatrix    matrix
	transitionMatrix   matrix
}

// Word is not in model.words
func (model hmm) tagProbabilityHeuristic(tag, word string) float {

	numReg := regexp.MustCompile(`^[0-9,.:]+$`)

	// word looks like a number!
	if numReg.MatchString(word) {

		if tag == "CARD" {
			return float(1)
		}

		return float(0)
	}

	nn := regexp.MustCompile(`^[A-Z]`)

	// word looks like a number!
	if nn.MatchString(word) {

		if tag == "NN" {
			return float(1)
		}

		return float(0)
	}

	return model.tagProbability[tag]
}

func (model hmm) forwardAlgorithm(phrase []string) []string {

	// phrase --> Observation
	startWord := phrase[0]

	initResults := make(probabilityMap)

	for s := range model.transitionMatrix {

		if val := model.words[startWord]; !val {
			initResults[s] = model.priorProbabilities[s] * model.tagProbabilityHeuristic(s, startWord)
			continue
		}

		initResults[s] = model.priorProbabilities[s] * model.emissionsMatrix[s][startWord]
	}

	// fmt.Println("startWord", startWord)
	// fmt.Println("initResults", initResults)

	aResults := make([]probabilityMap, 1)
	aResults[0] = initResults

	for k := 0; k < len(phrase)-1; k++ {

		result := make(probabilityMap)

		for s := range model.transitionMatrix {

			akSum := float(0)

			for q := range model.transitionMatrix {
				akSum += aResults[k][q] * model.transitionMatrix[q][s]
			}

			eVal := float(0.0)

			if val := model.words[phrase[k+1]]; val {
				eVal = model.emissionsMatrix[s][phrase[k+1]]
			} else {
				eVal = model.tagProbabilityHeuristic(s, phrase[k+1])
			}

			result[s] = eVal * akSum
		}

		aResults = append(aResults, result)
	}

	// fmt.Println("aResults", aResults)

	resultTag := make([]string, len(phrase))

	for inx := range phrase {

		result := aResults[inx]

		bestTag := "###"
		bestScore := float(0)

		for tag, score := range result {
			if score > bestScore {
				bestScore = score
				bestTag = tag
			}
		}

		resultTag[inx] = bestTag
	}

	// aggregate := 0.0
	// for _, a := range aResults[len(phrase)-1] {
	// 	aggregate += a
	// }
	//
	// fmt.Println("aggregate", aggregate)

	return resultTag
}

func priorProbabilitiesList(sentencesTags [][]string) probabilityMap {

	startTagsCount := make(map[string]int)
	for inx := range sentencesTags {
		tag := sentencesTags[inx][0]
		startTagsCount[tag]++
	}

	priorProbabilities := make(probabilityMap)
	for tag, count := range startTagsCount {
		priorProbabilities[tag] = float(count) / float(len(sentencesTags))
	}

	return priorProbabilities
}

func main() {

	// word := "2,99"
	//
	// var digitCheck = regexp.MustCompile(`^[0-9,.:]+$`)
	//
	// if digitCheck.MatchString(word) {
	// 	fmt.Printf("%q looks like a number.\n", word)
	// } else {
	// 	fmt.Printf("%q looks NOT like a number.\n", word)
	// }
	//
	// if 1 > 0 {
	// 	return
	// }

	content, err := ioutil.ReadFile("hdt-1-10000-train.tags")
	if err != nil {
		panic(err)
	}

	text := strings.TrimSpace(string(content))
	sentences := strings.Split(text, "\n\n")

	tagsCount := make(map[string]int)
	words := make(map[string]bool)

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

			words[wordTag[0]] = true
			tagsCount[wordTag[1]]++
		}
	}

	fmt.Println("tagsCount", tagsCount)

	tagsSum := 0
	for _, count := range tagsCount {
		tagsSum += count
	}

	tagProbability := make(probabilityMap)
	for tag, count := range tagsCount {
		tagProbability[tag] = float(count) / float(tagsSum)
	}

	priorProbabilities := priorProbabilitiesList(sentencesTags)
	fmt.Println("priorProbabilities", priorProbabilities)

	emissionsMatrix := emissionMatrixMap(sentencesWords, sentencesTags)
	transitionMatrix := transitionMatrixMap(tagsCount, sentencesTags)

	// fmt.Println("emissionsMatrix", emissionsMatrix)
	// fmt.Println("transitionMatrix", transitionMatrix)

	model := hmm{
		words:              words,
		tagProbability:     tagProbability,
		priorProbabilities: priorProbabilities,
		emissionsMatrix:    emissionsMatrix,
		transitionMatrix:   transitionMatrix,
	}

	content, err = ioutil.ReadFile("hdt-10001-12000-test.tags")
	if err != nil {
		panic(err)
	}

	text = strings.TrimSpace(string(content))
	text = strings.Trim(text, "\n")
	trainSentences := strings.Split(text, "\n\n")

	trainPhrases := make([][]string, len(trainSentences))

	for inx, sentence := range trainSentences {

		lines := strings.Split(sentence, "\n")
		trainPhrases[inx] = make([]string, len(lines))

		for iny, line := range lines {

			wordTag := strings.Split(line, "\t")
			trainPhrases[inx][iny] = wordTag[0]
		}
	}

	// evalText := ""
	//
	// for _, phrase := range trainPhrases {
	//
	// 	fmt.Println("phrase", phrase)
	//
	// 	tags := model.forwardAlgorithm(phrase)
	//
	// 	fmt.Println("tags", tags)
	//
	// 	for iny := range phrase {
	// 		evalText += phrase[iny] + "\t" + tags[iny] + "\n"
	// 	}
	//
	// 	evalText += "\n"
	// }
	//
	// err = ioutil.WriteFile("results.tags", []byte(evalText), 0755)
	// if err != nil {
	// 	panic(err)
	// }

	phrase := "Sie begründeten ihren Pessimismus unter anderem mit dem Umsatzrückgang nach den Attentaten am 11. September ."
	// phrase := "Pro Monat sind dafür 2,99 Euro fällig ."
	// phrase := "Dazu kommen zehn statt bisher fünf E-Mail-Adressen sowie zehn MByte Webspace ."
	phraseParts := strings.Split(phrase, " ")

	tags := model.forwardAlgorithm(phraseParts)

	fmt.Println(phraseParts)
	fmt.Println(tags)

	// diff -u hdt-10001-12000-test.tags results.tags | grep '^+' | wc -l
	// 8768
	// 5315
	// 5128 --> Number heuristic
	// 4991 --> NN heuristic
}
