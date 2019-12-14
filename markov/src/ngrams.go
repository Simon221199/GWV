package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

func findDescendantsMany(words []string, predecessor []string) []string {

	descendants := make([]string, 0)

	for inx := range words {

		descendantIndex := -1

		for iny, word := range predecessor {

			index := inx + iny

			if len(words) <= index {
				break
			}

			if words[ index ] != word {
				break
			}

			if len(predecessor) - 1 == iny && len(words) > index + 1 {
				descendantIndex = index + 1
			}
		}

		if descendantIndex >= 0 {
			descendants = append(descendants, words[ descendantIndex ])
		}
	}

	sort.Strings(descendants)

	return descendants
}

func trigrams(words []string) {

	// last := []string{"fick", "dich"}
	// last := []string{"angela", "merkel"}
	// last := []string{"ich", "bin"}
	last := []string{"du", "bist"}
	sentences := strings.Join(last, " ")

	descendantsStart := findDescendantsMany(words, last)
	fmt.Println(descendantsStart)
	// descendantsStart := findDescendantsMany(words, []string{"dich", "und"})
	// fmt.Println(descendantsStart)

	for inx := 0; inx < 8; inx++ {

		descendants := findDescendantsMany(words, last)

		// fmt.Println(last)
		// fmt.Println(descendants)

		if len(descendants) <= 0 {
			descendants = words
		}

		last[ 0 ] = last[ 1 ]
		last[ 1 ] = descendants[ rand.Intn(len(descendants)) ]

		sentences += " " + last[ 1 ]
	}

	fmt.Println(sentences)
}
