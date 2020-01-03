package main

import "fmt"

func (model hmm) viterbiAlgorithm(phrase []string) []string {

	alpha := make([]probabilityMap, len(phrase))
	pred := make([]map[string]*string, len(phrase))

	alpha[0] = make(probabilityMap)
	pred[0] = make(map[string]*string)

	for s := range model.tagProbability {

		eVal := float(0)

		if val := model.words[phrase[0]]; val {
			eVal = model.emissionsMatrix[s][phrase[0]]
		} else {
			eVal = model.tagProbabilityHeuristic(s, phrase[0])
		}

		alpha[0][s] = model.priorProbabilities[s] * eVal
		pred[0][s] = nil
	}

	fmt.Println("alpha", alpha[0])

	for k := 0; k < len(phrase)-1; k++ {

		alpha[k+1] = make(probabilityMap)
		pred[k+1] = make(map[string]*string)

		for s := range model.tagProbability {

			alphaMax := float(0)
			var predMax *string

			for q := range model.tagProbability {

				alphaQ := alpha[k][q] * model.transitionMatrix[q][s]

				if alphaMax < alphaQ {
					alphaMax = alphaQ
					predMax = &q
				}
			}

			eVal := float(0)

			if val := model.words[phrase[k+1]]; val {
				eVal = model.emissionsMatrix[s][phrase[k+1]]
			} else {
				eVal = model.tagProbabilityHeuristic(s, phrase[k+1])
			}

			alpha[k+1][s] = alphaMax * eVal
			pred[k+1][s] = predMax
		}
	}

	terminalState := "###"
	best := float(0)

	for s, score := range alpha[len(phrase)-1] {

		if best < score {
			terminalState = s
			best = score
		}
	}

	if terminalState == "###" {
		panic("terminalState == ###")
	}

	// fmt.Println("terminalState", terminalState)
	// fmt.Println("alpha", alpha[ len(phrase)-1 ])
	// fmt.Println("pred", pred)

	tags := make([]*string, len(phrase))
	tags[len(phrase)-1] = &terminalState

	for inx := len(phrase) - 2; inx > 0; inx-- {
		tag := pred[inx][*tags[inx+1]]

		if tag == nil {
			continue
		}

		tags[inx] = tag
	}

	xxx := make([]string, len(phrase))

	for inx, val := range tags {

		if val == nil {
			xxx[inx] = "###"
			continue
		}

		xxx[inx] = *val
	}

	return xxx
}
