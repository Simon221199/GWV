package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

func findDescendantsForMany(words []string, predecessor []string) []string {

	descendants := make([]string, 0)

	for inx := range words {

		descendantIndex := -1

		for iny, word := range predecessor {

			index := inx + iny

			// prevent index overflow
			if len(words) <= index {
				break
			}

			// check if sequent match
			if words[index] != word {
				break
			}

			// prevent overflow
			if len(words) <= (index + 1) {
				break
			}

			// break in end of comment is reached
			if words[index+1] == "" {
				break
			}

			// end reached
			if len(predecessor)-1 == iny {
				descendantIndex = index + 1
			}
		}

		if descendantIndex >= 0 {
			descendants = append(descendants, words[descendantIndex])
		}
	}

	sort.Strings(descendants)

	return descendants
}

func trigrams(words []string, start []string) {

	last := start
	sentences := strings.Join(last, " ")

	// descendantsStart := findDescendantsForMany(words, last)
	// fmt.Println(descendantsStart)

	for inx := 0; inx < 8; inx++ {

		descendants := findDescendantsForMany(words, last)

		if len(descendants) <= 0 {
			descendants = words
		}

		last[0] = last[1]
		last[1] = descendants[rand.Intn(len(descendants))]

		sentences += " " + last[1]
	}

	fmt.Println(sentences)
}
