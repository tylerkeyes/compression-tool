package internal

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Decoding functions

func Decompress_file(fileName string) {
	prefixCodes, binaryString, binaryLen := extract_stored_values(fileName)
	invertedCodes := invert_prefix_codes(prefixCodes)
	decodedString := convert_binary_to_char(invertedCodes, binaryString, binaryLen)
	write_to_file(fileName, decodedString)
}

func extract_stored_values(fileName string) (*map[rune]string, string, int) {
	bytes, err := os.ReadFile(fileName)
	assert_msg(err, fmt.Sprintf("error reading file: %v\n", err))
	data := string(bytes)

	var prefixCodesLen int
	var prefixCodes map[rune]string
	var bitNumber int
	var encodedString strings.Builder

	// find all lines of the file, using the starting and ending indexes in the string
	var regex = regexp.MustCompile(`(?m)^.*$`)
	fileLines := regex.FindAllStringIndex(data, -1)
	firstThree := fileLines[:3]

	prefixCodesLen, err = strconv.Atoi(data[firstThree[0][0]:firstThree[0][1]])
	prefixCodes = *create_prefix_codes_map(data[firstThree[1][0]:firstThree[1][1]], prefixCodesLen)
	bitNumber, err = strconv.Atoi(data[firstThree[2][0]:firstThree[2][1]])
	assert_msg(err, "problem reading binary data length")

	dataStart := firstThree[2][1] + 1 // adding 1 to move past the last '\n'
	compressedData := data[dataStart:]

	for i := 0; i < len(compressedData); i++ {
		binaryChar := fmt.Sprintf("%08b", compressedData[i])
		_, err = encodedString.WriteString(binaryChar)
		assert(err)
	}

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
