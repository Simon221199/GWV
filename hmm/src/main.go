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

// Build map with prior Probabilities π : S --> R+
func priorProbabilities(sentencesTags [][]string) probabilityMap {

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

// Build map with transition probabilities T: S x S --> R+
func transitionMatrix(tagsCount map[string]int, sentencesTags [][]string) matrix {

	//
	// "Implement a hidden markov model where the probability distribution
	// of Word[i] only depends on the state of Tag[i] and Tag[i] only depends on Tag[i-1]..."
	//
	// This function calculates that Tag[i] only depends on Tag[i-1]
	//

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

	//
	// Init transition matrix
	//

	transitionMatrix := make(matrix)
	for tag := range tagsCount {
		transitionMatrix[tag] = make(probabilityMap)
	}

	// Calculate probabilities
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

// Build map with emission probabilities E: S x O --> R+
func emissionMatrix(sentencesWords [][]string, sentencesTags [][]string) matrix {

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

// Make sure that your tagger can cope with input that includes words that are not in the training data.
// This gets called if word is not in model.words
func (model hmm) tagProbabilityHeuristic(tag, word string) float {

	//
	// Word looks like a number!
	//

	numReg := regexp.MustCompile(`^[0-9,.:]+$`)
	if numReg.MatchString(word) {

		if tag == "CARD" {
			return float(1)
		}

		return float(0)
	}

	//
	// Word looks like a noun!
	//

	nn := regexp.MustCompile(`^[A-Z]`)
	if nn.MatchString(word) || strings.Contains(word, "-") {

		if tag == "NN" {
			return float(1)
		}

		return float(0)
	}

	//
	// Take default probability for tag
	//

	return model.tagProbability[tag]
}

// Create a function that takes a list of words (possibly from the command line) and
// uses filtering to produce a corresponding list of PoS tags.
func (model hmm) forwardAlgorithm(phrase []string) []string {

	//
	// Init Forward algorithm
	//

	startWord := phrase[0]
	initResults := make(probabilityMap)

	for s := range model.transitionMatrix {

		if val := model.words[startWord]; !val {
			initResults[s] = model.priorProbabilities[s] * model.tagProbabilityHeuristic(s, startWord)
			continue
		}

		initResults[s] = model.priorProbabilities[s] * model.emissionsMatrix[s][startWord]
	}

	aResults := make([]probabilityMap, len(phrase))
	aResults[0] = initResults

	//
	// repeat, for k = 1 to k = t - 1 and for all s ∈ S
	// t = len(phrase)
	//

	for k := 0; k < len(phrase)-1; k++ {

		result := make(probabilityMap)

		for s := range model.transitionMatrix {

			akSum := float(0)

			for q := range model.transitionMatrix {
				akSum += aResults[k][q] * model.transitionMatrix[q][s]
			}

			eVal := float(0)

			if val := model.words[phrase[k+1]]; val {
				eVal = model.emissionsMatrix[s][phrase[k+1]]
			} else {
				eVal = model.tagProbabilityHeuristic(s, phrase[k+1])
			}

			result[s] = eVal * akSum
		}

		aResults[k+1] = result
	}

	//
	// Select best tag for every word in phrase
	//

	resultTag := make([]string, len(phrase))

	for inx := range phrase {

		result := aResults[inx]

		bestTag := "###"
		bestScore := float(0)

		for tag, score := range result {
			if bestScore < score {
				bestScore = score
				bestTag = tag
			}
		}

		resultTag[inx] = bestTag
	}

	//
	// aggregate
	//

	// aggregate := 0.0
	// for _, a := range aResults[len(phrase)-1] {
	// 	aggregate += a
	// }
	//
	// fmt.Println("aggregate", aggregate)

	return resultTag
}

// Backward algorithm, not working!
func (model hmm) backwardAlgorithm(phrase []string) []string {

	//
	// Init
	//

	initResults := make(probabilityMap)

	for s := range model.tagProbability {
		initResults[s] = 1
	}

	beta := make([]probabilityMap, len(phrase))
	beta[len(phrase)-1] = initResults

	//
	// repeat, for k = 1 to k = t - 1 and for all s ∈ S
	// t = len(phrase)
	//

	for k := len(phrase) - 2; k > 0; k++ {

		result := make(probabilityMap)

		for s := range model.transitionMatrix {

			val := float(0)

			for q := range model.transitionMatrix {

				eVal := float(0)

				if val := model.words[phrase[k+1]]; val {
					eVal = model.emissionsMatrix[q][phrase[k+1]]
				} else {
					eVal = model.tagProbabilityHeuristic(q, phrase[k+1])
				}

				val += beta[k+1][q] * model.transitionMatrix[q][s] * eVal
			}

			result[s] = val
		}

		beta[k+1] = result
	}

	return nil
}

func evaluate(model hmm) {

	content, err := ioutil.ReadFile("hdt-10001-12000-test.tags")
	if err != nil {
		panic(err)
	}

	text := strings.TrimSpace(string(content))
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

	evalText := ""

	for inx, phrase := range trainPhrases {

		fmt.Printf("\rEval phrase: %d/%d", len(trainPhrases), inx+1)

		tags := model.forwardAlgorithm(phrase)

		// fmt.Println("phrase", phrase)
		// fmt.Println("tags", tags)

		for iny := range phrase {
			evalText += phrase[iny] + "\t" + tags[iny] + "\n"
		}

		evalText += "\n"
	}

	fmt.Println()

	err = ioutil.WriteFile("results.tags", []byte(evalText), 0755)
	if err != nil {
		panic(err)
	}
}

func main() {

	//
	// Train with hdt-1-10000-train.tags
	//

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

	fmt.Println("tagsSum", tagsSum)

	tagProbability := make(probabilityMap)
	for tag, count := range tagsCount {
		tagProbability[tag] = float(count) / float(tagsSum)
	}

	priorProbabilities := priorProbabilities(sentencesTags)
	fmt.Println("priorProbabilities", priorProbabilities)

	emissionsMatrix := emissionMatrix(sentencesWords, sentencesTags)
	transitionMatrix := transitionMatrix(tagsCount, sentencesTags)

	// fmt.Println("emissionsMatrix", emissionsMatrix)
	// fmt.Println("transitionMatrix", transitionMatrix)

	model := hmm{
		words:              words,
		tagProbability:     tagProbability,
		priorProbabilities: priorProbabilities,
		emissionsMatrix:    emissionsMatrix,
		transitionMatrix:   transitionMatrix,
	}

	//
	// Build results.tags
	//

	evaluate(model)

	// Fail Sentences
	// phrase := "Sie begründeten ihren Pessimismus unter anderem mit dem Umsatzrückgang nach den Attentaten am 11. September ."

	// Unknown Word
	// phrase := "Pro Monat sind dafür 2,99 Euro fällig ."

	// Simple Example
	// phrase := "Dazu kommen zehn statt bisher fünf E-Mail-Adressen sowie zehn MByte Webspace ."

	// phraseParts := strings.Split(phrase, " ")
	// tags := model.forwardAlgorithm(phraseParts)
	// tags := model.viterbiAlgorithm(phraseParts)
	// fmt.Println(phraseParts)
	// fmt.Println(tags)

	// total 176311
	// diff -u hdt-10001-12000-test.tags results.tags | grep '^+' | wc -l
	// 8768
	// 5315
	// 5128 --> Number heuristic
	// 4968 --> NN heuristic --> 97.18% right
}
