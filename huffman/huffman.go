package huffman

import (
	"container/heap"
)

type HuffBaseNode interface {
	IsLeaf() bool
	Weight() int
}

// Huffman Leaf Node
type HuffLeafNode struct {
	element rune
	weight  int
}

func NewHuffLeafNode(element rune, weight int) *HuffLeafNode {
	return &HuffLeafNode{
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

func NewHuffInternalNode(weight int, left *HuffBaseNode, right *HuffBaseNode) *HuffInternalNode {
	return &HuffInternalNode{
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
	var tmp1 HuffBaseNode
	var tmp2 HuffBaseNode
	var tmp3 HuffBaseNode

	huffHeap := &HuffTreeHeap{}
	heap.Init(huffHeap)
	// create the input min-heap of leaf nodes
	for k, v := range *freqMap {
		node := NewHuffLeafNode(k, v)
		heap.Push(huffHeap, node)
	}

	// construct the Huffman Encoding tree using the min-heap
	for huffHeap.Len() > 1 {
		tmp1 = heap.Pop(huffHeap).(HuffBaseNode)
		tmp2 = heap.Pop(huffHeap).(HuffBaseNode)
		tmp3 = *(NewHuffInternalNode(tmp1.Weight()+tmp2.Weight(), &tmp1, &tmp2))
		heap.Push(huffHeap, tmp3)
	}

	return &tmp3
}
