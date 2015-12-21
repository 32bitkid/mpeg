package video

import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/huffman"

type DCTCoefficient struct {
	run   int
	level int32
}

type DCTSpecialToken string

var DCTSpecial = struct {
	EndOfBlock DCTSpecialToken
	Escape     DCTSpecialToken
}{
	DCTSpecialToken("end of block"),
	DCTSpecialToken("escape"),
}

type DCTCoefficientDecoderFn func(br bitreader.BitReader, n int) (int, int32, bool, error)

func newDCTCoefficientDecoder(tables [2]huffman.HuffmanTable) DCTCoefficientDecoderFn {
	inital := huffman.NewHuffmanDecoder(tables[0])
	rest := huffman.NewHuffmanDecoder(tables[1])
	return func(br bitreader.BitReader, n int) (run int, level int32, end bool, err error) {

		var decoder huffman.HuffmanDecoder
		if n == 0 {
			decoder = inital
		} else {
			decoder = rest
		}

		val, err := decoder.Decode(br)
		if err != nil {
			return 0, 0, false, err
		}

		if token, ok := val.(DCTSpecialToken); ok {

			if token == DCTSpecial.EndOfBlock {
				return 0, 0, true, nil
			}

			if token == DCTSpecial.Escape {
				val, err := br.Read32(6)
				if err != nil {
					return 0, 0, false, err
				}
				run = int(val)
				sign, err := br.ReadBit()
				if err != nil {
					return 0, 0, false, err
				}
				val, err = br.Read32(11)
				if err != nil {
					return 0, 0, false, err
				}
				level = int32(val)
				if sign {
					level += -2048
				}
				return run, level, false, nil
			}

		} else if dct, ok := val.(DCTCoefficient); ok {

			run = dct.run
			level = dct.level
			sign, err := br.ReadBit()
			if err != nil {
				return 0, 0, false, err
			}
			if sign {
				level *= -1
			}
			return run, level, false, nil

		}

		return 0, 0, false, ErrUnexpectedDecodedValueType
	}
}

var dctCoefficientTables = [2][2]huffman.HuffmanTable{
	{
		{
			"1":            DCTCoefficient{0, 1},
			"011":          DCTCoefficient{1, 1},
			"0100":         DCTCoefficient{0, 2},
			"0101":         DCTCoefficient{2, 1},
			"0010 1":       DCTCoefficient{0, 3},
			"0011 1":       DCTCoefficient{3, 1},
			"0011 0":       DCTCoefficient{4, 1},
			"0001 10":      DCTCoefficient{1, 2},
			"0001 11":      DCTCoefficient{5, 1},
			"0001 01":      DCTCoefficient{6, 1},
			"0001 00":      DCTCoefficient{7, 1},
			"0000 110":     DCTCoefficient{0, 4},
			"0000 100":     DCTCoefficient{2, 2},
			"0000 111":     DCTCoefficient{8, 1},
			"0000 101":     DCTCoefficient{9, 1},
			"0000 01":      DCTSpecial.Escape,
			"0010 0110":    DCTCoefficient{0, 5},
			"0010 0001":    DCTCoefficient{0, 6},
			"0010 0101":    DCTCoefficient{1, 3},
			"0010 0100":    DCTCoefficient{3, 2},
			"0010 0111":    DCTCoefficient{10, 1},
			"0010 0011":    DCTCoefficient{11, 1},
			"0010 0010":    DCTCoefficient{12, 1},
			"0010 0000":    DCTCoefficient{13, 1},
			"0000 0010 10": DCTCoefficient{0, 7},
			"0000 0011 00": DCTCoefficient{1, 4},
			"0000 0010 11": DCTCoefficient{2, 3},
			"0000 0011 11": DCTCoefficient{4, 2},
			"0000 0010 01": DCTCoefficient{5, 2},
			"0000 0011 10": DCTCoefficient{14, 1},
			"0000 0011 01": DCTCoefficient{15, 1},
			"0000 0010 00": DCTCoefficient{16, 1},

			"0000 0001 1101":   DCTCoefficient{0, 8},
			"0000 0001 1000":   DCTCoefficient{0, 9},
			"0000 0001 0011":   DCTCoefficient{0, 10},
			"0000 0001 0000":   DCTCoefficient{0, 11},
			"0000 0001 1011":   DCTCoefficient{1, 5},
			"0000 0001 0100":   DCTCoefficient{2, 4},
			"0000 0001 1100":   DCTCoefficient{3, 3},
			"0000 0001 0010":   DCTCoefficient{4, 3},
			"0000 0001 1110":   DCTCoefficient{6, 2},
			"0000 0001 0101":   DCTCoefficient{7, 2},
			"0000 0001 0001":   DCTCoefficient{8, 2},
			"0000 0001 1111":   DCTCoefficient{17, 1},
			"0000 0001 1010":   DCTCoefficient{18, 1},
			"0000 0001 1001":   DCTCoefficient{19, 1},
			"0000 0001 0111":   DCTCoefficient{20, 1},
			"0000 0001 0110":   DCTCoefficient{21, 1},
			"0000 0000 1101 0": DCTCoefficient{0, 12},
			"0000 0000 1100 1": DCTCoefficient{0, 13},
			"0000 0000 1100 0": DCTCoefficient{0, 14},
			"0000 0000 1011 1": DCTCoefficient{0, 15},
			"0000 0000 1011 0": DCTCoefficient{1, 6},
			"0000 0000 1010 1": DCTCoefficient{1, 7},
			"0000 0000 1010 0": DCTCoefficient{2, 5},
			"0000 0000 1001 1": DCTCoefficient{3, 4},
			"0000 0000 1001 0": DCTCoefficient{5, 3},
			"0000 0000 1000 1": DCTCoefficient{9, 2},
			"0000 0000 1000 0": DCTCoefficient{10, 2},
			"0000 0000 1111 1": DCTCoefficient{22, 1},
			"0000 0000 1111 0": DCTCoefficient{23, 1},
			"0000 0000 1110 1": DCTCoefficient{24, 1},
			"0000 0000 1110 0": DCTCoefficient{25, 1},
			"0000 0000 1101 1": DCTCoefficient{26, 1},

			"0000 0000 0111 11":  DCTCoefficient{0, 16},
			"0000 0000 0111 10":  DCTCoefficient{0, 17},
			"0000 0000 0111 01":  DCTCoefficient{0, 18},
			"0000 0000 0111 00":  DCTCoefficient{0, 19},
			"0000 0000 0110 11":  DCTCoefficient{0, 20},
			"0000 0000 0110 10":  DCTCoefficient{0, 21},
			"0000 0000 0110 01":  DCTCoefficient{0, 22},
			"0000 0000 0110 00":  DCTCoefficient{0, 23},
			"0000 0000 0101 11":  DCTCoefficient{0, 24},
			"0000 0000 0101 10":  DCTCoefficient{0, 25},
			"0000 0000 0101 01":  DCTCoefficient{0, 26},
			"0000 0000 0101 00":  DCTCoefficient{0, 27},
			"0000 0000 0100 11":  DCTCoefficient{0, 28},
			"0000 0000 0100 10":  DCTCoefficient{0, 29},
			"0000 0000 0100 01":  DCTCoefficient{0, 30},
			"0000 0000 0100 00":  DCTCoefficient{0, 31},
			"0000 0000 0011 000": DCTCoefficient{0, 32},
			"0000 0000 0010 111": DCTCoefficient{0, 33},
			"0000 0000 0010 110": DCTCoefficient{0, 34},
			"0000 0000 0010 101": DCTCoefficient{0, 35},
			"0000 0000 0010 100": DCTCoefficient{0, 36},
			"0000 0000 0010 011": DCTCoefficient{0, 37},
			"0000 0000 0010 010": DCTCoefficient{0, 38},
			"0000 0000 0010 001": DCTCoefficient{0, 39},
			"0000 0000 0010 000": DCTCoefficient{0, 40},
			"0000 0000 0011 111": DCTCoefficient{1, 8},
			"0000 0000 0011 110": DCTCoefficient{1, 9},
			"0000 0000 0011 101": DCTCoefficient{1, 10},
			"0000 0000 0011 100": DCTCoefficient{1, 11},
			"0000 0000 0011 011": DCTCoefficient{1, 12},
			"0000 0000 0011 010": DCTCoefficient{1, 13},
			"0000 0000 0011 001": DCTCoefficient{1, 14},

			"0000 0000 0001 0011": DCTCoefficient{1, 15},
			"0000 0000 0001 0010": DCTCoefficient{1, 16},
			"0000 0000 0001 0001": DCTCoefficient{1, 17},
			"0000 0000 0001 0000": DCTCoefficient{1, 18},
			"0000 0000 0001 0100": DCTCoefficient{6, 3},
			"0000 0000 0001 1010": DCTCoefficient{11, 2},
			"0000 0000 0001 1001": DCTCoefficient{12, 2},
			"0000 0000 0001 1000": DCTCoefficient{13, 2},
			"0000 0000 0001 0111": DCTCoefficient{14, 2},
			"0000 0000 0001 0110": DCTCoefficient{15, 2},
			"0000 0000 0001 0101": DCTCoefficient{16, 2},
			"0000 0000 0001 1111": DCTCoefficient{27, 1},
			"0000 0000 0001 1110": DCTCoefficient{28, 1},
			"0000 0000 0001 1101": DCTCoefficient{29, 1},
			"0000 0000 0001 1100": DCTCoefficient{30, 1},
			"0000 0000 0001 1011": DCTCoefficient{31, 1},
		}, {
			"10":           DCTSpecial.EndOfBlock,
			"11":           DCTCoefficient{0, 1},
			"011":          DCTCoefficient{1, 1},
			"0100":         DCTCoefficient{0, 2},
			"0101":         DCTCoefficient{2, 1},
			"0010 1":       DCTCoefficient{0, 3},
			"0011 1":       DCTCoefficient{3, 1},
			"0011 0":       DCTCoefficient{4, 1},
			"0001 10":      DCTCoefficient{1, 2},
			"0001 11":      DCTCoefficient{5, 1},
			"0001 01":      DCTCoefficient{6, 1},
			"0001 00":      DCTCoefficient{7, 1},
			"0000 110":     DCTCoefficient{0, 4},
			"0000 100":     DCTCoefficient{2, 2},
			"0000 111":     DCTCoefficient{8, 1},
			"0000 101":     DCTCoefficient{9, 1},
			"0000 01":      DCTSpecial.Escape,
			"0010 0110":    DCTCoefficient{0, 5},
			"0010 0001":    DCTCoefficient{0, 6},
			"0010 0101":    DCTCoefficient{1, 3},
			"0010 0100":    DCTCoefficient{3, 2},
			"0010 0111":    DCTCoefficient{10, 1},
			"0010 0011":    DCTCoefficient{11, 1},
			"0010 0010":    DCTCoefficient{12, 1},
			"0010 0000":    DCTCoefficient{13, 1},
			"0000 0010 10": DCTCoefficient{0, 7},
			"0000 0011 00": DCTCoefficient{1, 4},
			"0000 0010 11": DCTCoefficient{2, 3},
			"0000 0011 11": DCTCoefficient{4, 2},
			"0000 0010 01": DCTCoefficient{5, 2},
			"0000 0011 10": DCTCoefficient{14, 1},
			"0000 0011 01": DCTCoefficient{15, 1},
			"0000 0010 00": DCTCoefficient{16, 1},

			"0000 0001 1101":   DCTCoefficient{0, 8},
			"0000 0001 1000":   DCTCoefficient{0, 9},
			"0000 0001 0011":   DCTCoefficient{0, 10},
			"0000 0001 0000":   DCTCoefficient{0, 11},
			"0000 0001 1011":   DCTCoefficient{1, 5},
			"0000 0001 0100":   DCTCoefficient{2, 4},
			"0000 0001 1100":   DCTCoefficient{3, 3},
			"0000 0001 0010":   DCTCoefficient{4, 3},
			"0000 0001 1110":   DCTCoefficient{6, 2},
			"0000 0001 0101":   DCTCoefficient{7, 2},
			"0000 0001 0001":   DCTCoefficient{8, 2},
			"0000 0001 1111":   DCTCoefficient{17, 1},
			"0000 0001 1010":   DCTCoefficient{18, 1},
			"0000 0001 1001":   DCTCoefficient{19, 1},
			"0000 0001 0111":   DCTCoefficient{20, 1},
			"0000 0001 0110":   DCTCoefficient{21, 1},
			"0000 0000 1101 0": DCTCoefficient{0, 12},
			"0000 0000 1100 1": DCTCoefficient{0, 13},
			"0000 0000 1100 0": DCTCoefficient{0, 14},
			"0000 0000 1011 1": DCTCoefficient{0, 15},
			"0000 0000 1011 0": DCTCoefficient{1, 6},
			"0000 0000 1010 1": DCTCoefficient{1, 7},
			"0000 0000 1010 0": DCTCoefficient{2, 5},
			"0000 0000 1001 1": DCTCoefficient{3, 4},
			"0000 0000 1001 0": DCTCoefficient{5, 3},
			"0000 0000 1000 1": DCTCoefficient{9, 2},
			"0000 0000 1000 0": DCTCoefficient{10, 2},
			"0000 0000 1111 1": DCTCoefficient{22, 1},
			"0000 0000 1111 0": DCTCoefficient{23, 1},
			"0000 0000 1110 1": DCTCoefficient{24, 1},
			"0000 0000 1110 0": DCTCoefficient{25, 1},
			"0000 0000 1101 1": DCTCoefficient{26, 1},

			"0000 0000 0111 11":  DCTCoefficient{0, 16},
			"0000 0000 0111 10":  DCTCoefficient{0, 17},
			"0000 0000 0111 01":  DCTCoefficient{0, 18},
			"0000 0000 0111 00":  DCTCoefficient{0, 19},
			"0000 0000 0110 11":  DCTCoefficient{0, 20},
			"0000 0000 0110 10":  DCTCoefficient{0, 21},
			"0000 0000 0110 01":  DCTCoefficient{0, 22},
			"0000 0000 0110 00":  DCTCoefficient{0, 23},
			"0000 0000 0101 11":  DCTCoefficient{0, 24},
			"0000 0000 0101 10":  DCTCoefficient{0, 25},
			"0000 0000 0101 01":  DCTCoefficient{0, 26},
			"0000 0000 0101 00":  DCTCoefficient{0, 27},
			"0000 0000 0100 11":  DCTCoefficient{0, 28},
			"0000 0000 0100 10":  DCTCoefficient{0, 29},
			"0000 0000 0100 01":  DCTCoefficient{0, 30},
			"0000 0000 0100 00":  DCTCoefficient{0, 31},
			"0000 0000 0011 000": DCTCoefficient{0, 32},
			"0000 0000 0010 111": DCTCoefficient{0, 33},
			"0000 0000 0010 110": DCTCoefficient{0, 34},
			"0000 0000 0010 101": DCTCoefficient{0, 35},
			"0000 0000 0010 100": DCTCoefficient{0, 36},
			"0000 0000 0010 011": DCTCoefficient{0, 37},
			"0000 0000 0010 010": DCTCoefficient{0, 38},
			"0000 0000 0010 001": DCTCoefficient{0, 39},
			"0000 0000 0010 000": DCTCoefficient{0, 40},
			"0000 0000 0011 111": DCTCoefficient{1, 8},
			"0000 0000 0011 110": DCTCoefficient{1, 9},
			"0000 0000 0011 101": DCTCoefficient{1, 10},
			"0000 0000 0011 100": DCTCoefficient{1, 11},
			"0000 0000 0011 011": DCTCoefficient{1, 12},
			"0000 0000 0011 010": DCTCoefficient{1, 13},
			"0000 0000 0011 001": DCTCoefficient{1, 14},

			"0000 0000 0001 0011": DCTCoefficient{1, 15},
			"0000 0000 0001 0010": DCTCoefficient{1, 16},
			"0000 0000 0001 0001": DCTCoefficient{1, 17},
			"0000 0000 0001 0000": DCTCoefficient{1, 18},
			"0000 0000 0001 0100": DCTCoefficient{6, 3},
			"0000 0000 0001 1010": DCTCoefficient{11, 2},
			"0000 0000 0001 1001": DCTCoefficient{12, 2},
			"0000 0000 0001 1000": DCTCoefficient{13, 2},
			"0000 0000 0001 0111": DCTCoefficient{14, 2},
			"0000 0000 0001 0110": DCTCoefficient{15, 2},
			"0000 0000 0001 0101": DCTCoefficient{16, 2},
			"0000 0000 0001 1111": DCTCoefficient{27, 1},
			"0000 0000 0001 1110": DCTCoefficient{28, 1},
			"0000 0000 0001 1101": DCTCoefficient{29, 1},
			"0000 0000 0001 1100": DCTCoefficient{30, 1},
			"0000 0000 0001 1011": DCTCoefficient{31, 1},
		},
	},
	{
		{
			"10":           DCTCoefficient{0, 1},
			"010":          DCTCoefficient{1, 1},
			"110":          DCTCoefficient{0, 2},
			"0010 1":       DCTCoefficient{2, 1},
			"0111":         DCTCoefficient{0, 3},
			"0011 1":       DCTCoefficient{3, 1},
			"0001 10":      DCTCoefficient{4, 1},
			"0011 0":       DCTCoefficient{1, 2},
			"0001 11":      DCTCoefficient{5, 1},
			"0000 110":     DCTCoefficient{6, 1},
			"0000 100":     DCTCoefficient{7, 1},
			"1110 0":       DCTCoefficient{0, 4},
			"0000 111":     DCTCoefficient{2, 2},
			"0000 101":     DCTCoefficient{8, 1},
			"1111 000":     DCTCoefficient{9, 1},
			"0000 01":      DCTSpecial.Escape,
			"1110 1":       DCTCoefficient{0, 5},
			"0001 01":      DCTCoefficient{0, 6},
			"1111 001":     DCTCoefficient{1, 3},
			"0010 0110":    DCTCoefficient{3, 2},
			"1111 010":     DCTCoefficient{10, 1},
			"0010 0001":    DCTCoefficient{11, 1},
			"0010 0101":    DCTCoefficient{12, 1},
			"0010 0100":    DCTCoefficient{13, 1},
			"0001 00":      DCTCoefficient{0, 7},
			"0010 0111":    DCTCoefficient{1, 4},
			"1111 1100":    DCTCoefficient{2, 3},
			"1111 1101":    DCTCoefficient{4, 2},
			"0000 0010 0":  DCTCoefficient{5, 2},
			"0000 0010 1":  DCTCoefficient{14, 1},
			"0000 0011 1":  DCTCoefficient{15, 1},
			"0000 0011 01": DCTCoefficient{16, 1},

			"1111 011":         DCTCoefficient{0, 8},
			"1111 100":         DCTCoefficient{0, 9},
			"0010 0011":        DCTCoefficient{0, 10},
			"0010 0010":        DCTCoefficient{0, 11},
			"0010 0000":        DCTCoefficient{1, 5},
			"0000 0011 00":     DCTCoefficient{2, 4},
			"0000 0001 1100":   DCTCoefficient{3, 3},
			"0000 0001 0010":   DCTCoefficient{4, 3},
			"0000 0001 1110":   DCTCoefficient{6, 2},
			"0000 0001 0101":   DCTCoefficient{7, 2},
			"0000 0001 0001":   DCTCoefficient{8, 2},
			"0000 0001 1111":   DCTCoefficient{17, 1},
			"0000 0001 1010":   DCTCoefficient{18, 1},
			"0000 0001 1001":   DCTCoefficient{19, 1},
			"0000 0001 0111":   DCTCoefficient{20, 1},
			"0000 0001 0110":   DCTCoefficient{21, 1},
			"1111 1010":        DCTCoefficient{0, 12},
			"1111 1011":        DCTCoefficient{0, 13},
			"1111 1110":        DCTCoefficient{0, 14},
			"1111 1111":        DCTCoefficient{0, 15},
			"0000 0000 1011 0": DCTCoefficient{1, 6},
			"0000 0000 1010 1": DCTCoefficient{1, 7},
			"0000 0000 1010 0": DCTCoefficient{2, 5},
			"0000 0000 1001 1": DCTCoefficient{3, 4},
			"0000 0000 1001 0": DCTCoefficient{5, 3},
			"0000 0000 1000 1": DCTCoefficient{9, 2},
			"0000 0000 1000 0": DCTCoefficient{10, 2},
			"0000 0000 1111 1": DCTCoefficient{22, 1},
			"0000 0000 1111 0": DCTCoefficient{23, 1},
			"0000 0000 1110 1": DCTCoefficient{24, 1},
			"0000 0000 1110 0": DCTCoefficient{25, 1},
			"0000 0000 1101 1": DCTCoefficient{26, 1},

			"0000 0000 0111 11":  DCTCoefficient{0, 16},
			"0000 0000 0111 10":  DCTCoefficient{0, 17},
			"0000 0000 0111 01":  DCTCoefficient{0, 18},
			"0000 0000 0111 00":  DCTCoefficient{0, 19},
			"0000 0000 0110 11":  DCTCoefficient{0, 20},
			"0000 0000 0110 10":  DCTCoefficient{0, 21},
			"0000 0000 0110 01":  DCTCoefficient{0, 22},
			"0000 0000 0110 00":  DCTCoefficient{0, 23},
			"0000 0000 0101 11":  DCTCoefficient{0, 24},
			"0000 0000 0101 10":  DCTCoefficient{0, 25},
			"0000 0000 0101 01":  DCTCoefficient{0, 26},
			"0000 0000 0101 00":  DCTCoefficient{0, 27},
			"0000 0000 0100 11":  DCTCoefficient{0, 28},
			"0000 0000 0100 10":  DCTCoefficient{0, 29},
			"0000 0000 0100 01":  DCTCoefficient{0, 30},
			"0000 0000 0100 00":  DCTCoefficient{0, 31},
			"0000 0000 0011 000": DCTCoefficient{0, 32},
			"0000 0000 0010 111": DCTCoefficient{0, 33},
			"0000 0000 0010 110": DCTCoefficient{0, 34},
			"0000 0000 0010 101": DCTCoefficient{0, 35},
			"0000 0000 0010 100": DCTCoefficient{0, 36},
			"0000 0000 0010 011": DCTCoefficient{0, 37},
			"0000 0000 0010 010": DCTCoefficient{0, 38},
			"0000 0000 0010 001": DCTCoefficient{0, 39},
			"0000 0000 0010 000": DCTCoefficient{0, 40},
			"0000 0000 0011 111": DCTCoefficient{1, 8},
			"0000 0000 0011 110": DCTCoefficient{1, 9},
			"0000 0000 0011 101": DCTCoefficient{1, 10},
			"0000 0000 0011 100": DCTCoefficient{1, 11},
			"0000 0000 0011 011": DCTCoefficient{1, 12},
			"0000 0000 0011 010": DCTCoefficient{1, 13},
			"0000 0000 0011 001": DCTCoefficient{1, 14},

			"0000 0000 0001 0011": DCTCoefficient{1, 15},
			"0000 0000 0001 0010": DCTCoefficient{1, 16},
			"0000 0000 0001 0001": DCTCoefficient{1, 17},
			"0000 0000 0001 0000": DCTCoefficient{1, 18},
			"0000 0000 0001 0100": DCTCoefficient{6, 3},
			"0000 0000 0001 1010": DCTCoefficient{11, 2},
			"0000 0000 0001 1001": DCTCoefficient{12, 2},
			"0000 0000 0001 1000": DCTCoefficient{13, 2},
			"0000 0000 0001 0111": DCTCoefficient{14, 2},
			"0000 0000 0001 0110": DCTCoefficient{15, 2},
			"0000 0000 0001 0101": DCTCoefficient{16, 2},
			"0000 0000 0001 1111": DCTCoefficient{27, 1},
			"0000 0000 0001 1110": DCTCoefficient{28, 1},
			"0000 0000 0001 1101": DCTCoefficient{29, 1},
			"0000 0000 0001 1100": DCTCoefficient{30, 1},
			"0000 0000 0001 1011": DCTCoefficient{31, 1},
		}, {
			"0110":         DCTSpecial.EndOfBlock,
			"10":           DCTCoefficient{0, 1},
			"010":          DCTCoefficient{1, 1},
			"110":          DCTCoefficient{0, 2},
			"0010 1":       DCTCoefficient{2, 1},
			"0111":         DCTCoefficient{0, 3},
			"0011 1":       DCTCoefficient{3, 1},
			"0001 10":      DCTCoefficient{4, 1},
			"0011 0":       DCTCoefficient{1, 2},
			"0001 11":      DCTCoefficient{5, 1},
			"0000 110":     DCTCoefficient{6, 1},
			"0000 100":     DCTCoefficient{7, 1},
			"1110 0":       DCTCoefficient{0, 4},
			"0000 111":     DCTCoefficient{2, 2},
			"0000 101":     DCTCoefficient{8, 1},
			"1111 000":     DCTCoefficient{9, 1},
			"0000 01":      DCTSpecial.Escape,
			"1110 1":       DCTCoefficient{0, 5},
			"0001 01":      DCTCoefficient{0, 6},
			"1111 001":     DCTCoefficient{1, 3},
			"0010 0110":    DCTCoefficient{3, 2},
			"1111 010":     DCTCoefficient{10, 1},
			"0010 0001":    DCTCoefficient{11, 1},
			"0010 0101":    DCTCoefficient{12, 1},
			"0010 0100":    DCTCoefficient{13, 1},
			"0001 00":      DCTCoefficient{0, 7},
			"0010 0111":    DCTCoefficient{1, 4},
			"1111 1100":    DCTCoefficient{2, 3},
			"1111 1101":    DCTCoefficient{4, 2},
			"0000 0010 0":  DCTCoefficient{5, 2},
			"0000 0010 1":  DCTCoefficient{14, 1},
			"0000 0011 1":  DCTCoefficient{15, 1},
			"0000 0011 01": DCTCoefficient{16, 1},

			"1111 011":         DCTCoefficient{0, 8},
			"1111 100":         DCTCoefficient{0, 9},
			"0010 0011":        DCTCoefficient{0, 10},
			"0010 0010":        DCTCoefficient{0, 11},
			"0010 0000":        DCTCoefficient{1, 5},
			"0000 0011 00":     DCTCoefficient{2, 4},
			"0000 0001 1100":   DCTCoefficient{3, 3},
			"0000 0001 0010":   DCTCoefficient{4, 3},
			"0000 0001 1110":   DCTCoefficient{6, 2},
			"0000 0001 0101":   DCTCoefficient{7, 2},
			"0000 0001 0001":   DCTCoefficient{8, 2},
			"0000 0001 1111":   DCTCoefficient{17, 1},
			"0000 0001 1010":   DCTCoefficient{18, 1},
			"0000 0001 1001":   DCTCoefficient{19, 1},
			"0000 0001 0111":   DCTCoefficient{20, 1},
			"0000 0001 0110":   DCTCoefficient{21, 1},
			"1111 1010":        DCTCoefficient{0, 12},
			"1111 1011":        DCTCoefficient{0, 13},
			"1111 1110":        DCTCoefficient{0, 14},
			"1111 1111":        DCTCoefficient{0, 15},
			"0000 0000 1011 0": DCTCoefficient{1, 6},
			"0000 0000 1010 1": DCTCoefficient{1, 7},
			"0000 0000 1010 0": DCTCoefficient{2, 5},
			"0000 0000 1001 1": DCTCoefficient{3, 4},
			"0000 0000 1001 0": DCTCoefficient{5, 3},
			"0000 0000 1000 1": DCTCoefficient{9, 2},
			"0000 0000 1000 0": DCTCoefficient{10, 2},
			"0000 0000 1111 1": DCTCoefficient{22, 1},
			"0000 0000 1111 0": DCTCoefficient{23, 1},
			"0000 0000 1110 1": DCTCoefficient{24, 1},
			"0000 0000 1110 0": DCTCoefficient{25, 1},
			"0000 0000 1101 1": DCTCoefficient{26, 1},

			"0000 0000 0111 11":  DCTCoefficient{0, 16},
			"0000 0000 0111 10":  DCTCoefficient{0, 17},
			"0000 0000 0111 01":  DCTCoefficient{0, 18},
			"0000 0000 0111 00":  DCTCoefficient{0, 19},
			"0000 0000 0110 11":  DCTCoefficient{0, 20},
			"0000 0000 0110 10":  DCTCoefficient{0, 21},
			"0000 0000 0110 01":  DCTCoefficient{0, 22},
			"0000 0000 0110 00":  DCTCoefficient{0, 23},
			"0000 0000 0101 11":  DCTCoefficient{0, 24},
			"0000 0000 0101 10":  DCTCoefficient{0, 25},
			"0000 0000 0101 01":  DCTCoefficient{0, 26},
			"0000 0000 0101 00":  DCTCoefficient{0, 27},
			"0000 0000 0100 11":  DCTCoefficient{0, 28},
			"0000 0000 0100 10":  DCTCoefficient{0, 29},
			"0000 0000 0100 01":  DCTCoefficient{0, 30},
			"0000 0000 0100 00":  DCTCoefficient{0, 31},
			"0000 0000 0011 000": DCTCoefficient{0, 32},
			"0000 0000 0010 111": DCTCoefficient{0, 33},
			"0000 0000 0010 110": DCTCoefficient{0, 34},
			"0000 0000 0010 101": DCTCoefficient{0, 35},
			"0000 0000 0010 100": DCTCoefficient{0, 36},
			"0000 0000 0010 011": DCTCoefficient{0, 37},
			"0000 0000 0010 010": DCTCoefficient{0, 38},
			"0000 0000 0010 001": DCTCoefficient{0, 39},
			"0000 0000 0010 000": DCTCoefficient{0, 40},
			"0000 0000 0011 111": DCTCoefficient{1, 8},
			"0000 0000 0011 110": DCTCoefficient{1, 9},
			"0000 0000 0011 101": DCTCoefficient{1, 10},
			"0000 0000 0011 100": DCTCoefficient{1, 11},
			"0000 0000 0011 011": DCTCoefficient{1, 12},
			"0000 0000 0011 010": DCTCoefficient{1, 13},
			"0000 0000 0011 001": DCTCoefficient{1, 14},

			"0000 0000 0001 0011": DCTCoefficient{1, 15},
			"0000 0000 0001 0010": DCTCoefficient{1, 16},
			"0000 0000 0001 0001": DCTCoefficient{1, 17},
			"0000 0000 0001 0000": DCTCoefficient{1, 18},
			"0000 0000 0001 0100": DCTCoefficient{6, 3},
			"0000 0000 0001 1010": DCTCoefficient{11, 2},
			"0000 0000 0001 1001": DCTCoefficient{12, 2},
			"0000 0000 0001 1000": DCTCoefficient{13, 2},
			"0000 0000 0001 0111": DCTCoefficient{14, 2},
			"0000 0000 0001 0110": DCTCoefficient{15, 2},
			"0000 0000 0001 0101": DCTCoefficient{16, 2},
			"0000 0000 0001 1111": DCTCoefficient{27, 1},
			"0000 0000 0001 1110": DCTCoefficient{28, 1},
			"0000 0000 0001 1101": DCTCoefficient{29, 1},
			"0000 0000 0001 1100": DCTCoefficient{30, 1},
			"0000 0000 0001 1011": DCTCoefficient{31, 1},
		},
	},
}

type DCTDCSizeDecoderFn func(br bitreader.BitReader) (uint, error)

func newDCTDCSizeDecoder(table huffman.HuffmanTable) DCTDCSizeDecoderFn {
	decoder := huffman.NewHuffmanDecoder(table)
	return func(br bitreader.BitReader) (uint, error) {
		val, err := decoder.Decode(br)
		if err != nil {
			return 0, err
		} else if i, ok := val.(uint); ok {
			return i, nil
		} else {
			return 0, ErrUnexpectedDecodedValueType
		}
	}
}

// Table B-12
var dctDcSizeLuminanceTable = huffman.HuffmanTable{
	"100":         uint(0),
	"00":          uint(1),
	"01":          uint(2),
	"101":         uint(3),
	"110":         uint(4),
	"1110":        uint(5),
	"1111 0":      uint(6),
	"1111 10":     uint(7),
	"1111 110":    uint(8),
	"1111 1110":   uint(9),
	"1111 1111 0": uint(10),
	"1111 1111 1": uint(11),
}

// Table B-13
var dctDcSizeChrominanceTable = huffman.HuffmanTable{
	"00":           uint(0),
	"01":           uint(1),
	"10":           uint(2),
	"110":          uint(3),
	"1110":         uint(4),
	"1111 0":       uint(5),
	"1111 10":      uint(6),
	"1111 110":     uint(7),
	"1111 1110":    uint(8),
	"1111 1111 0":  uint(9),
	"1111 1111 10": uint(10),
	"1111 1111 11": uint(11),
}

var DCTDCSizeDecoders = struct {
	Luma   DCTDCSizeDecoderFn
	Chroma DCTDCSizeDecoderFn
}{
	newDCTDCSizeDecoder(dctDcSizeLuminanceTable),
	newDCTDCSizeDecoder(dctDcSizeChrominanceTable),
}

var DCTCoefficientDecoders = struct {
	TableZero DCTCoefficientDecoderFn
	TableOne  DCTCoefficientDecoderFn
}{

	newDCTCoefficientDecoder(dctCoefficientTables[0]),
	newDCTCoefficientDecoder(dctCoefficientTables[1]),
}
