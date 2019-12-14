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

	rand.Seed(time.Now().UnixNano())

	bytes, err := ioutil.ReadFile("ggcc-one-word-per-line.txt")
	if err != nil {
		log.Fatal(err)
	}

	content := string(bytes)
	words := strings.Split(content, "\n")

	fmt.Printf("words=%d\n", len(words))

	// bigrams(words)

	// last := []string{"fick", "dich"}
	// last := []string{"angela", "merkel"}
	// last := []string{"ich", "bin"}
	// last := []string{"du", "bist"}

	for inx := 0; inx < 10; inx++ {
		startWord := words[rand.Intn(len(words))]
		descendants := findDescendants(words, startWord)
		start := []string{startWord, descendants[rand.Intn(len(descendants))]}
		trigrams(words, start)
	}
}
