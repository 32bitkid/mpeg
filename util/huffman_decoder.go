package util

import "errors"

type HuffmanDecoder interface {
	Decode(BitReader32) (interface{}, error)
}

type HuffmanTable map[string]interface{}

var ErrMissingHuffmanValue = errors.New("missing huffman value")

func NewHuffmanDecoder(init HuffmanTable) HuffmanDecoder {
	return NewBinaryTreeHuffmanDecoder(init)
}
