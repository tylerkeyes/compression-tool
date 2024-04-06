package main

import (
	"log"
	"os"
)

/*
1. Read the text and determine the frequency of each character occurring.
2. Build the binary tree from the frequencies.
3. Generate the prefix-code table from the tree.
4. Encode the text using the code table.
5. Encode the tree - we’ll need to include this in the output file so we can decode it.
6. Write the encoded tree and text to an output field
*/
func main() {
	if len(os.Args) < 2 {
		log.Fatal("file name not given")
	}
	fileName := os.Args[1]
	log.Printf("file name: %v\n", fileName)

	freqMap := count_frequency(fileName)
	log.Printf("X: %v\n", (*freqMap)['X'])
	log.Printf("t: %v\n", (*freqMap)['t'])
	log.Printf("frequency map:\n%+v\n", freqMap)
}

func count_frequency(fileName string) *map[rune]int {
	charFreq := make(map[rune]int, 256) // start with size for each UTF-8 character
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("error reading file: %v\n", err)
	}

	data := string(bytes)

	for _, char := range data {
		v, ok := charFreq[char]
		if ok {
			charFreq[char] = v + 1
		} else {
			charFreq[char] = 1
		}
	}

	return &charFreq
}
