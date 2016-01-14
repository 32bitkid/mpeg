package video

import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/huffman"

type codedBlockPattern int
type patternCode [12]bool

func coded_block_pattern(br bitreader.BitReader, chroma_format ChromaFormat) (codedBlockPattern, error) {
	val, err := decodeCpb(br)

	if ChromaFormat422 == chroma_format {
		panic("unsupported: cbp 4:2:2")
	}
	if ChromaFormat444 == chroma_format {
		panic("unsupported: cbp 4:4:4")
	}

	return val, err
}

func (cbp codedBlockPattern) decode(intra, pattern bool, chroma_format ChromaFormat) (pattern_code patternCode) {
	for i := 0; i < 12; i++ {
		if intra {
			pattern_code[i] = true
		} else {
			pattern_code[i] = false
		}
	}

	if pattern {
		for i := 0; i < 6; i++ {
			if mask := 1 << uint(5-i); (int(cbp) & mask) == mask {
				pattern_code[i] = true
			}
		}

		if chroma_format == ChromaFormat422 || chroma_format == ChromaFormat444 {
			panic("unsupported: coded block pattern chroma format")
		}
	}

	return
}

var cbpDecoder = huffman.NewHuffmanDecoder(huffman.HuffmanTable{
	"111":         codedBlockPattern(60),
	"1101":        codedBlockPattern(4),
	"1100":        codedBlockPattern(8),
	"1011":        codedBlockPattern(16),
	"1010":        codedBlockPattern(32),
	"1001 1":      codedBlockPattern(12),
	"1001 0":      codedBlockPattern(48),
	"1000 1":      codedBlockPattern(20),
	"1000 0":      codedBlockPattern(40),
	"0111 1":      codedBlockPattern(28),
	"0111 0":      codedBlockPattern(44),
	"0110 1":      codedBlockPattern(52),
	"0110 0":      codedBlockPattern(56),
	"0101 1":      codedBlockPattern(1),
	"0101 0":      codedBlockPattern(61),
	"0100 1":      codedBlockPattern(2),
	"0100 0":      codedBlockPattern(62),
	"0011 11":     codedBlockPattern(24),
	"0011 10":     codedBlockPattern(36),
	"0011 01":     codedBlockPattern(3),
	"0011 00":     codedBlockPattern(63),
	"0010 111":    codedBlockPattern(5),
	"0010 110":    codedBlockPattern(9),
	"0010 101":    codedBlockPattern(17),
	"0010 100":    codedBlockPattern(33),
	"0010 011":    codedBlockPattern(6),
	"0010 010":    codedBlockPattern(10),
	"0010 001":    codedBlockPattern(18),
	"0010 000":    codedBlockPattern(34),
	"0001 1111":   codedBlockPattern(7),
	"0001 1110":   codedBlockPattern(11),
	"0001 1101":   codedBlockPattern(19),
	"0001 1100":   codedBlockPattern(35),
	"0001 1011":   codedBlockPattern(13),
	"0001 1010":   codedBlockPattern(49),
	"0001 1001":   codedBlockPattern(21),
	"0001 1000":   codedBlockPattern(41),
	"0001 0111":   codedBlockPattern(14),
	"0001 0110":   codedBlockPattern(50),
	"0001 0101":   codedBlockPattern(22),
	"0001 0100":   codedBlockPattern(42),
	"0001 0011":   codedBlockPattern(15),
	"0001 0010":   codedBlockPattern(51),
	"0001 0001":   codedBlockPattern(23),
	"0001 0000":   codedBlockPattern(43),
	"0000 1111":   codedBlockPattern(25),
	"0000 1110":   codedBlockPattern(37),
	"0000 1101":   codedBlockPattern(26),
	"0000 1100":   codedBlockPattern(38),
	"0000 1011":   codedBlockPattern(29),
	"0000 1010":   codedBlockPattern(45),
	"0000 1001":   codedBlockPattern(53),
	"0000 1000":   codedBlockPattern(57),
	"0000 0111":   codedBlockPattern(30),
	"0000 0110":   codedBlockPattern(46),
	"0000 0101":   codedBlockPattern(54),
	"0000 0100":   codedBlockPattern(58),
	"0000 0011 1": codedBlockPattern(31),
	"0000 0011 0": codedBlockPattern(47),
	"0000 0010 1": codedBlockPattern(55),
	"0000 0010 0": codedBlockPattern(59),
	"0000 0001 1": codedBlockPattern(27),
	"0000 0001 0": codedBlockPattern(39),
	"0000 0000 1": codedBlockPattern(0),
})

func decodeCpb(br bitreader.BitReader) (codedBlockPattern, error) {
	val, err := cbpDecoder.Decode(br)
	if err != nil {
		return 0, err
	} else if i, ok := val.(codedBlockPattern); ok {
		return i, nil
	}

	return 0, huffman.ErrMissingHuffmanValue
}
