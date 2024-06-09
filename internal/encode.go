package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func assert_msg(err error, msg string) {
	if err != nil {
		log.Fatalf("%v: %v\n", msg, err)
	}
}

// Encoding Functions

func Compress_file(fileName string, archiveName string) {
	freqMap, _, _ := count_frequency(fileName)

	huffmanTree := BuildTree(freqMap)

	prefixCodes := CreatePrefixCodes(huffmanTree)

	encodedBinary, bitLen := create_binary_encoding(fileName, prefixCodes)

	write_compressed_data(encodedBinary, prefixCodes, archiveName, bitLen)
}

func count_frequency(fileName string) (*map[rune]int, *[]rune, int) {
	charFreq := make(map[rune]int, 0) // start with size for each UTF-8 character
	charList := make([]rune, 0)

	bytes, err := os.ReadFile(fileName)
	assert_msg(err, "error reading file")
	data := string(bytes)

	for _, char := range data {
		v, ok := charFreq[char]
		if ok {
			charFreq[char] = v + 1
		} else {
			charFreq[char] = 1
			charList = append(charList, char)
		}
	}

	return &charFreq, &charList, len(data)
}

func create_binary_encoding(fileName string, prefixCodes *map[rune]string) (*[]byte, int) {
	bytes, err := os.ReadFile(fileName)
	// this technically should not ever happen, since the input file was already read
	assert_msg(err, "error reading file")

	data := string(bytes)
	var sb strings.Builder

	for _, char := range data {
		sb.WriteString((*prefixCodes)[char])
	}

	encodedString := sb.String()
	bitLen := len(encodedString)
	binaryEncoding := convert_string_to_bytes(encodedString)

	return binaryEncoding, bitLen
}

func convert_string_to_bytes(encodedString string) *[]byte {
	byteEncoding := make([]byte, 0)
	dataLength := len(encodedString)

	window := 0
	var nextByte byte
	for window < dataLength {
		var currByte byte
		windowEnd := min(window+8, dataLength)
		for c := window; c < windowEnd; c++ {
			if encodedString[c] == '0' {
				nextByte = 0
			} else {
				nextByte = 1
			}
			currByte = currByte << 1
			currByte = currByte + nextByte
		}
		if window+8 > dataLength {
			currByte = currByte << (byte(dataLength - window))
		}
		byteEncoding = append(byteEncoding, currByte)
		window += 8
	}

	return &byteEncoding
}

func write_compressed_data(encodedBytes *[]byte, prefixCodes *map[rune]string, newFileName string, bitLen int) {
	outputFile, err := os.Create(newFileName)
	assert_msg(err, fmt.Sprintf("problem writing to file %v\n", newFileName))
	defer outputFile.Close()

	codeSize := len(*prefixCodes)

	_, err = outputFile.WriteString(fmt.Sprintf("%v\n", codeSize))
	assert_msg(err, fmt.Sprintf("problem writing to file %v\n", newFileName))

	prefixCodeString := convert_map_to_string(prefixCodes)
	_, err = outputFile.WriteString(prefixCodeString)
	assert_msg(err, fmt.Sprintf("problem writing to file %v\n", newFileName))

	_, err = outputFile.WriteString(fmt.Sprintf("\n%v\n", bitLen))
	assert_msg(err, fmt.Sprintf("problem writing to file %v\n", newFileName))

	_, err = outputFile.Write(*encodedBytes)
	assert_msg(err, fmt.Sprintf("problem writing to file %v\n", newFileName))
}

func convert_map_to_string(prefixCodes *map[rune]string) string {
	var sb strings.Builder

	for key, val := range *prefixCodes {
		codedString := fmt.Sprintf("%v:%v,", key, val)
		sb.WriteString(codedString)
	}

	return sb.String()
}
