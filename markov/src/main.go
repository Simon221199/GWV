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
	trigrams(words)
}
