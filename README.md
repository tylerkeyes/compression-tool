## Compression Tool

A go implementation of a compression tool, using Huffman Encoding/Decoding.

### Usage

Compress a file:
`go run compression.go -c input.txt`

Specify the name of the compressed file:
`go run compression.go -c -d input.zip input.txt`
If the name of the compressed file is not specified, the name will be in the form: <file name>.zip

Decompress a file:
`go run compression.go -x input.zip`
