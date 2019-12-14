package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func findDescendants(words []string, predecessor string) []string {

	descendants := make([]string, 0)

	for inx, word := range words {

		if word != predecessor {
			continue
		}

		// prevent index overflow
		if len(words) <= (inx + 1) {
			continue
		}

		// break in end of comment is reached
		if words[inx+1] == "" {
			continue
		}

		descendants = append(descendants, words[inx+1])
	}

	sort.Strings(descendants)

	return descendants
}

func bigrams(words []string) {

	length := 8
	// start := "angela"
	// start := "haha"
	// start := "fick"
	start := "ich"
	last := start
	sentences := last

	descendantsStart := findDescendants(words, last)
	fmt.Println(descendantsStart)

	for inx := 0; inx < length; inx++ {

		descendants := findDescendants(words, last)
		last = descendants[rand.Intn(len(descendants))]
		sentences += " " + last
	}

	fmt.Println(sentences)
}
