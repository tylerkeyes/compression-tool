package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tylerkeyes/compression-tool/huffman"
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

	freqMap, _, dataLen := count_frequency(fileName)

	huffmanTree := huffman.BuildTree(freqMap)

	prefixCodes := huffman.CreatePrefixCodes(huffmanTree)

	encodedBinary := create_binary_encoding(fileName, prefixCodes)

	log.Printf("compression improvement: %v\n", float64(len(*encodedBinary))/float64(dataLen)*100)

	write_compressed_data(encodedBinary, prefixCodes, fileName)
}

func count_frequency(fileName string) (*map[rune]int, *[]rune, int) {
	charFreq := make(map[rune]int, 0) // start with size for each UTF-8 character
	charList := make([]rune, 0)
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
			charList = append(charList, char)
		}
	}

	return &charFreq, &charList, len(data)
}

func create_binary_encoding(fileName string, prefixCodes *map[rune]string) *[]byte {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		// this technically should not ever happen, since the input file was already read
		log.Fatalf("error reading file: %v\n", err)
	}

	data := string(bytes)
	var sb strings.Builder

	for _, char := range data {
		sb.WriteString((*prefixCodes)[char])
	}

	encodedString := sb.String()
	binaryEncoding := convert_string_to_bytes(encodedString)

	return binaryEncoding
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

func write_compressed_data(encodedBytes *[]byte, prefixCodes *map[rune]string, fileName string) {
	newFileName := fileName + ".zip"
	outputFile, err := os.Create(newFileName)
	if err != nil {
		log.Fatalf("problem writing to file %v\n", newFileName)
	}
	defer outputFile.Close()

	dataSize := len(*encodedBytes)
	codeSize := len(*prefixCodes)

	_, err = outputFile.WriteString(fmt.Sprintf("%v\n", codeSize))
	if err != nil {
		log.Fatalf("problem writing to file %v\n", newFileName)
	}

	prefixCodeString := convert_map_to_string(prefixCodes)
	_, err = outputFile.WriteString(prefixCodeString)
	if err != nil {
		log.Fatalf("problem writing to file %v\n", newFileName)
	}

	_, err = outputFile.WriteString(fmt.Sprintf("\n%v\n", dataSize))
	if err != nil {
		log.Fatalf("problem writing to file %v\n", newFileName)
	}

	_, err = outputFile.Write(*encodedBytes)
	if err != nil {
		log.Fatalf("problem writing to file %v\n", newFileName)
	}
}

func convert_map_to_string(prefixCodes *map[rune]string) string {
	var sb strings.Builder

	for key, val := range *prefixCodes {
		codedString := fmt.Sprintf("%v:%v,", key, val)
		sb.WriteString(codedString)
	}

	return sb.String()
}
