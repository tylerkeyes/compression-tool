package internal

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Decoding functions

func Decompress_file(fileName string) {
	prefixCodes, binaryString, binaryLen := extract_stored_values(fileName)
	log.Printf("prefix codes: %+v\n", prefixCodes)
	invertedCodes := invert_prefix_codes(prefixCodes)
	log.Printf("inverted prefix codes: %+v\n", invertedCodes)
	decodedString := convert_binary_to_char(invertedCodes, binaryString, binaryLen)
	fmt.Printf("decoded string: %+v\n", decodedString)
	log.Printf("prefixCodes len: %v, binaryString: %v, binaryLen: %v, decodedString: %v\n", len(*invertedCodes), len(binaryString), binaryLen, len(decodedString))
	write_to_file(fileName, decodedString)
}

func extract_stored_values(fileName string) (*map[rune]string, string, int) {
	bytes, err := os.ReadFile(fileName)
	assert_msg(err, fmt.Sprintf("error reading file: %v\n", err))
	data := string(bytes)
	newlineCount := 0
	newlineLimit := 3
	trail := 0

	var prefixCodesLen int
	var prefixCodes map[rune]string
	var bitNumber int

	var encodedString strings.Builder

	for pos, char := range data {
		if newlineCount < newlineLimit {
			if char == '\n' {
				if newlineCount == 0 {
					prefixCodesLen, err = strconv.Atoi(data[trail:pos])
					assert_msg(err, "problem reading prefix codes map size")
					newlineCount += 1
					trail = pos + 1
				} else if newlineCount == 1 {
					prefixCodes = *create_prefix_codes_map(data[trail:pos], prefixCodesLen)
					trail = pos + 1
					newlineCount += 1
				} else if newlineCount == 2 {
					bitNumber, err = strconv.Atoi(data[trail:pos])
					assert_msg(err, "problem reading binary data length")
					newlineCount += 1
					trail = pos + 1
				}
			}
		} else {
			binaryChar := fmt.Sprintf("%08b", char)
			_, err = encodedString.WriteString(binaryChar)
			assert(err)
		}
	}

	log.Printf("binary string: %v\n", encodedString.String())

	return &prefixCodes, encodedString.String(), bitNumber
}

func create_prefix_codes_map(data string, mapSize int) *map[rune]string {
	prefixCodes := make(map[rune]string)

	mapElements := strings.Split(data, ",")
	for _, pair := range mapElements {
		if len(pair) > 0 {
			keyPair := strings.Split(pair, ":")
			key, err := strconv.Atoi(keyPair[0])
			assert(err)
			keyRune := rune(key)
			prefixCodes[keyRune] = keyPair[1]
		}
	}

	if len(prefixCodes) != mapSize {
		log.Fatalf("prefix codes size does not match the expected size: %v != %v\n", len(prefixCodes), mapSize)
	}

	return &prefixCodes
}

// Invert the prefix codes map to match binary strings to a rune
func invert_prefix_codes(prefixCodes *map[rune]string) *map[string]rune {
	newPrefixCodes := make(map[string]rune, 0)

	for k, v := range *prefixCodes {
		newPrefixCodes[v] = k
	}

	return &newPrefixCodes
}

func convert_binary_to_char(prefixCodes *map[string]rune, binaryString string, binaryLen int) string {
	var result strings.Builder
	var word strings.Builder
	var err error

	for i := 0; i < binaryLen; i++ {
		c := binaryString[i]

		_, err = word.WriteString(string(c))
		assert(err)

		if char, ok := (*prefixCodes)[word.String()]; ok {
			_, err = result.WriteString(string(char))
			assert(err)
			word.Reset()
		}
	}

	return result.String()
}

func write_to_file(fileName string, decodedString string) {
	splitName := strings.Split(fileName, ".zip")
	if len(splitName) != 2 {
		log.Fatalf("should not be here")
	}
	outputName := splitName[0]
	err := os.WriteFile(outputName, []byte(decodedString), 0644)
	assert(err)
}
