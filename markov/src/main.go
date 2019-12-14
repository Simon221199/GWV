package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

func main() {

	// test := []string{"A", "B", "C", "D"}
	// fmt.Println(test)
	// fmt.Println(findDescendantsMany(test, []string{"A", "B"}))
	// fmt.Println(findDescendantsMany(test, []string{"A", "C"}))
	// fmt.Println(findDescendantsMany(test, []string{"B", "C"}))
	// fmt.Println(findDescendantsMany(test, []string{"C", "D"}))
	// fmt.Println(test[1:])
	//
	// for iny, word2 := range test[1:] {
	// 	fmt.Println(word2)
	// 	fmt.Println(iny)
	// }

	rand.Seed(time.Now().UnixNano())

	bytes, err := ioutil.ReadFile("ggcc-one-word-per-line.txt")
	if err != nil {
		log.Fatal(err)
	}

	content := string(bytes)
	words := strings.Split(content, "\n")

	// for inx, word := range words {
	// 	words[ inx ] = strings.ToLower(word)
	// }

	fmt.Printf("words=%d\n", len(words))

	// bigrams(words)
	trigrams(words)
}
