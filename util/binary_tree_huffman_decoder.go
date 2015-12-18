package util

func NewBinaryTreeHuffmanDecoder(init HuffmanTable) HuffmanDecoder {
	root, depth := parseInitIntoTree(init)
	return &binaryTreeHuffmanDecoder{root, depth}
}

type binaryTreeHuffmanDecoder struct {
	root  *binaryHuffmanNode
	depth uint
}

type binaryHuffmanNode struct {
	left  interface{}
	right interface{}
}

func parseInitIntoTree(init HuffmanTable) (*binaryHuffmanNode, uint) {
	root := &binaryHuffmanNode{}
	var depth uint = 0

	for _, value := range init {
		currentNode := root

		bitStringLength := len(value.BitString)
		if uint(bitStringLength) > depth {
			depth = uint(bitStringLength)
		}

		for i := 0; i < bitStringLength; i++ {

			bit := value.BitString[i : i+1]

			if i < bitStringLength-1 {
				// Descending
				if bit == "1" {
					if currentNode.left == nil {
						nextNode := &binaryHuffmanNode{}
						currentNode.left = nextNode
						currentNode = nextNode
					} else if nextNode, ok := currentNode.left.(*binaryHuffmanNode); ok {
						currentNode = nextNode
					} else {
						panic("Invalid huffman tree")
					}
				} else {
					if currentNode.right == nil {
						nextNode := &binaryHuffmanNode{}
						currentNode.right = nextNode
						currentNode = nextNode
					} else if nextNode, ok := currentNode.right.(*binaryHuffmanNode); ok {
						currentNode = nextNode
					} else {
						panic("Invalid huffman tree")
					}
				}
			} else {
				// Ending

				if bit == "1" {
					currentNode.left = value.Value
				} else {
					currentNode.right = value.Value
				}
			}
		}
	}

	return root, depth
}

func (self *binaryTreeHuffmanDecoder) Decode(br BitReader32) (interface{}, error) {

	nextbits, err := br.Peek32(self.depth)

	if err != nil {
		return nil, err
	}

	currentNode := self.root

	for i := uint(1); i <= self.depth; i++ {
		var val interface{}

		mask := uint32(1) << (self.depth - i)
		bit := nextbits&mask == mask

		if bit {
			val = currentNode.left
		} else {
			val = currentNode.right
		}

		if nextNode, ok := val.(*binaryHuffmanNode); ok {
			currentNode = nextNode
		} else {
			br.Trash(i)
			return val, nil
		}
	}

	return nil, ErrMissingHuffmanValue
}