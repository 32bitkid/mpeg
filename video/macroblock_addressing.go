package video

import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/huffman"
import "errors"

type macroblockAddressIncrementDecoder func(bitreader.BitReader) (uint32, error)

var ErrUnexpectedDecodedValueType = errors.New("unexpected decoded value type")
var macroblockAddressIncrementHuffmanDecoder = huffman.NewHuffmanDecoder(huffman.HuffmanTable{
	"1":             uint32(1),
	"011":           uint32(2),
	"010":           uint32(3),
	"0011":          uint32(4),
	"0010":          uint32(5),
	"0001 1":        uint32(6),
	"0001 0":        uint32(7),
	"0000 111":      uint32(8),
	"0000 110":      uint32(9),
	"0000 1011":     uint32(10),
	"0000 1010":     uint32(11),
	"0000 1001":     uint32(12),
	"0000 1000":     uint32(13),
	"0000 0111":     uint32(14),
	"0000 0110":     uint32(15),
	"0000 0101 11":  uint32(16),
	"0000 0101 10":  uint32(17),
	"0000 0101 01":  uint32(18),
	"0000 0101 00":  uint32(19),
	"0000 0100 11":  uint32(20),
	"0000 0100 10":  uint32(21),
	"0000 0100 011": uint32(22),
	"0000 0100 010": uint32(23),
	"0000 0100 001": uint32(24),
	"0000 0100 000": uint32(25),
	"0000 0011 111": uint32(26),
	"0000 0011 110": uint32(27),
	"0000 0011 101": uint32(28),
	"0000 0011 100": uint32(29),
	"0000 0011 011": uint32(30),
	"0000 0011 010": uint32(31),
	"0000 0011 001": uint32(32),
	"0000 0011 000": uint32(33),
})

func decodeMacroblockAddressIncrement(br bitreader.BitReader) (uint32, error) {
	val, err := macroblockAddressIncrementHuffmanDecoder.Decode(br)
	if err != nil {
		return 0, err
	} else if i, ok := val.(uint32); ok {
		return i, nil
	} else {
		return 0, ErrUnexpectedDecodedValueType
	}
}

var MacroblockAddressIncrementDecoder = struct {
	Decode macroblockAddressIncrementDecoder
}{
	decodeMacroblockAddressIncrement,
}
