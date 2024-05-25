package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tylerkeyes/compression-tool/internal"
)

/*
Encode/Compress
1. Read the text and determine the frequency of each character occurring.
2. Build the binary tree from the frequencies.
3. Generate the prefix-code table from the tree.
4. Encode the text using the code table.
5. Encode the tree - we’ll need to include this in the output file so we can decode it.
6. Write the encoded tree and text to an output field

Decode/Restore
1. Read the first 3 lines.
2. Create the prefixCodes map from the string on line 2.
3. Read remaining bytes and create full binary string.
4. Convert binary string to original file.
*/
/*
	TODO: Read compressed file and restore to original format
*/
var (
	encode          = flag.Bool("c", false, "compress the input file")
	decode          = flag.Bool("x", false, "decompress the input file")
	destinationFile = flag.String("d", "", "name of the compressed file")
	help            = flag.Bool("h", false, "help information")
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

func main() {
	if len(os.Args) < 2 {
		log.Fatal("file name not given")
	}
	flag.Parse()

	if *help {
		helpMsg := `fct - File Compression Tool
	Commands:

	-c		compress the input file
	-x		decompress the input file
	-d		name of the compressed file, defaults to <input_file>.zip
	`
		fmt.Printf(helpMsg)
		os.Exit(0)
	}

	fileName := flag.Args()[0]
	log.Printf("command line args: %+v\n", flag.Args())

	if *encode {
		if "" == *destinationFile {
			*destinationFile = fileName + ".zip"
		}
		internal.Compress_file(fileName, *destinationFile)
	} else if *decode {
		internal.Decompress_file(fileName)
	}
}
