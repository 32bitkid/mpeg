package video

import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/huffman"
import "errors"

type macroblockAddressIncrementDecoderFn func(bitreader.BitReader) (int, error)

var ErrUnexpectedDecodedValueType = errors.New("unexpected decoded value type")
var macroblockAddressIncrementHuffmanDecoder = huffman.NewHuffmanDecoder(huffman.HuffmanTable{
	"1":             1,
	"011":           2,
	"010":           3,
	"0011":          4,
	"0010":          5,
	"0001 1":        6,
	"0001 0":        7,
	"0000 111":      8,
	"0000 110":      9,
	"0000 1011":     10,
	"0000 1010":     11,
	"0000 1001":     12,
	"0000 1000":     13,
	"0000 0111":     14,
	"0000 0110":     15,
	"0000 0101 11":  16,
	"0000 0101 10":  17,
	"0000 0101 01":  18,
	"0000 0101 00":  19,
	"0000 0100 11":  20,
	"0000 0100 10":  21,
	"0000 0100 011": 22,
	"0000 0100 010": 23,
	"0000 0100 001": 24,
	"0000 0100 000": 25,
	"0000 0011 111": 26,
	"0000 0011 110": 27,
	"0000 0011 101": 28,
	"0000 0011 100": 29,
	"0000 0011 011": 30,
	"0000 0011 010": 31,
	"0000 0011 001": 32,
	"0000 0011 000": 33,
})

var macroblockAddressIncrementDecoder = struct {
	Decode macroblockAddressIncrementDecoderFn
}{
	func(br bitreader.BitReader) (int, error) {
		val, err := macroblockAddressIncrementHuffmanDecoder.Decode(br)
		if err != nil {
			return 0, err
		} else if i, ok := val.(int); ok {
			return i, nil
		} else {
			return 0, ErrUnexpectedDecodedValueType
		}
	},
}
