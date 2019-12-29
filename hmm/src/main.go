package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	fmt.Println("hallo")

	content, err := ioutil.ReadFile("hdt-1-10000-train.tagsCount")
	if err != nil {
		panic(err)
	}

	text := string(content)
	parts := strings.Split(text, "\n")

	startTagsCount := 0
	startTags := make(map[string]int)
	tagsCount := make(map[string]int)

	isFirstWord := true

	for _, line := range parts {

		words := strings.Split(line, "\t")

		// sentence ended
		if len(words) < 2 {
			isFirstWord = true
			continue
		}

		if isFirstWord {
			startTags[words[1]]++
			startTagsCount++
			isFirstWord = false
		}

		tagsCount[words[1]]++
	}

	fmt.Println("tagsCount", tagsCount)
	fmt.Println("startTags", startTags)

	startProbability := make(map[string]float32)

	for key, val := range startTags {
		startProbability[key] = float32(val) / float32(startTagsCount)
	}

	fmt.Println("startProbability", startProbability)

	inx := 0
	tags := make(map[string]int)

	for tag := range tagsCount {
		tags[tag] = inx
		inx++
	}

	// init transition matrix
	transitionMatrix := make([][]int, len(tags))
	for inx := range transitionMatrix {
		transitionMatrix[inx] = make([]int, len(tags))
	}

}
