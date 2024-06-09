package internal

import (
	"container/heap"
	"fmt"
	_ "log"
)

type HuffBaseNode interface {
	IsLeaf() bool
	Weight() int
}

func isPowerOf2(n int) bool {
	if n == 0 {
		return false
	}
	for n != 1 {
		if n%2 != 0 {
			return false
		}
		n = n / 2
	}
	return true
}

func VisualizeTree(tree *HuffBaseNode) {
	queue := make([]HuffBaseNode, 0)

	queue = append(queue, *tree)
	count := 1
	output := ""
	var node HuffBaseNode

	for len(queue) > 0 {
		node = queue[0]
		queue = queue[1:]

		output = output + fmt.Sprint(node.Weight()) + " | "
		if isPowerOf2(count) {
			output = output + "\n"
		}

		if count == 4 {
			break
		}
		count += 1
		if node.IsLeaf() {
			node = node.(HuffLeafNode)
		} else {
			leafNode := node.(HuffInternalNode)
			queue = append(queue, *leafNode.left)
			queue = append(queue, *leafNode.right)
		}
	}
}

// Huffman Leaf Node
type HuffLeafNode struct {
	element rune
	weight  int
}

func NewHuffLeafNode(element rune, weight int) HuffLeafNode {
	return HuffLeafNode{
		element,
		weight,
	}
}

func (node HuffLeafNode) IsLeaf() bool {
	return true
}

func (node HuffLeafNode) Weight() int {
	return node.weight
}

// Huffman Internal Node
type HuffInternalNode struct {
	weight int
	left   *HuffBaseNode
	right  *HuffBaseNode
}

func NewHuffInternalNode(weight int, left *HuffBaseNode, right *HuffBaseNode) HuffInternalNode {
	return HuffInternalNode{
		weight,
		left,
		right,
	}
}

func (node HuffInternalNode) IsLeaf() bool {
	return false
}

func (node HuffInternalNode) Weight() int {
	return node.weight
}

// Heap used for creating Huffman Tree
type HuffTreeHeap []HuffBaseNode

func (h HuffTreeHeap) Len() int {
	return len(h)
}

func (h HuffTreeHeap) Less(i, j int) bool {
	return h[i].Weight() < h[j].Weight()
}

func (h HuffTreeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *HuffTreeHeap) Push(x any) {
	*h = append(*h, x.(HuffBaseNode))
}

func (h *HuffTreeHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	*h = old[0 : n-1]
	return x
}

// Build the Huffman Encoding Tree
func BuildTree(freqMap *map[rune]int) *HuffBaseNode {
	huffHeap := &HuffTreeHeap{}
	heap.Init(huffHeap)
	// create the input min-heap of leaf nodes
	for k, v := range *freqMap {
		node := NewHuffLeafNode(k, v)
		heap.Push(huffHeap, node)
	}

	var result HuffBaseNode
	// construct the Huffman Encoding tree using the min-heap
	for huffHeap.Len() > 1 {
		left := heap.Pop(huffHeap).(HuffBaseNode)
		right := heap.Pop(huffHeap).(HuffBaseNode)
		top := NewHuffInternalNode(left.Weight()+right.Weight(), &left, &right)
		heap.Push(huffHeap, top)
		result = top
	}

	return &result
}

// Generate the prefix codes table
func CreatePrefixCodes(tree *HuffBaseNode) *map[rune]string {
	prefixCodes := make(map[rune]string, 0)

	traverseHuffmanTree(*tree, &prefixCodes, "")

	return &prefixCodes
}

// recursively traverse the tree, assigning encoded values when an end is reached
func traverseHuffmanTree(tree HuffBaseNode, prefixCodes *map[rune]string, encoding string) {
	if tree.IsLeaf() {
		leaf := tree.(HuffLeafNode)
		(*prefixCodes)[leaf.element] = encoding
		return
	}

	leftNode := tree.(HuffInternalNode).left
	rightNode := tree.(HuffInternalNode).right

	encodingL := encoding + "0"
	encodingR := encoding + "1"
	traverseHuffmanTree(*leftNode, prefixCodes, encodingL)
	traverseHuffmanTree(*rightNode, prefixCodes, encodingR)
}
