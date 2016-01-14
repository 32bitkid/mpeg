package video

import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/huffman"

type dctCoefficient struct {
	run   int
	level int32
}

type dctEndOfBlock struct{}
type dctEscape struct{}

type dctCoefficientDecoderFn func(br bitreader.BitReader, n int) (int, int32, bool, error)

func newDCTCoefficientDecoder(tables [2]huffman.HuffmanTable) dctCoefficientDecoderFn {
	inital := huffman.NewHuffmanDecoder(tables[0])
	rest := huffman.NewHuffmanDecoder(tables[1])
	return func(br bitreader.BitReader, n int) (run int, level int32, end bool, err error) {

		var decoder huffman.HuffmanDecoder
		if n == 0 {
			decoder = inital
		} else {
			decoder = rest
		}

		if val, err := decoder.Decode(br); err != nil {
			return 0, 0, false, err
		} else if _, ok := val.(dctEndOfBlock); ok {
			return 0, 0, true, nil
		} else if _, ok := val.(dctEscape); ok {
			if val, err := br.Read32(6); err != nil {
				return 0, 0, false, err
			} else {
				run = int(val)
			}

			if sign, err := br.ReadBit(); err != nil {
				return 0, 0, false, err
			} else if sign {
				level = -2048
			}

			if val, err := br.Read32(11); err != nil {
				return 0, 0, false, err
			} else {
				level += int32(val)
			}

			return run, level, false, nil
		} else if dct, ok := val.(dctCoefficient); ok {
			run = dct.run
			level = dct.level
			if sign, err := br.ReadBit(); err != nil {
				return 0, 0, false, err
			} else if sign {
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
			"1":            dctCoefficient{0, 1},
			"011":          dctCoefficient{1, 1},
			"0100":         dctCoefficient{0, 2},
			"0101":         dctCoefficient{2, 1},
			"0010 1":       dctCoefficient{0, 3},
			"0011 1":       dctCoefficient{3, 1},
			"0011 0":       dctCoefficient{4, 1},
			"0001 10":      dctCoefficient{1, 2},
			"0001 11":      dctCoefficient{5, 1},
			"0001 01":      dctCoefficient{6, 1},
			"0001 00":      dctCoefficient{7, 1},
			"0000 110":     dctCoefficient{0, 4},
			"0000 100":     dctCoefficient{2, 2},
			"0000 111":     dctCoefficient{8, 1},
			"0000 101":     dctCoefficient{9, 1},
			"0000 01":      dctEscape{},
			"0010 0110":    dctCoefficient{0, 5},
			"0010 0001":    dctCoefficient{0, 6},
			"0010 0101":    dctCoefficient{1, 3},
			"0010 0100":    dctCoefficient{3, 2},
			"0010 0111":    dctCoefficient{10, 1},
			"0010 0011":    dctCoefficient{11, 1},
			"0010 0010":    dctCoefficient{12, 1},
			"0010 0000":    dctCoefficient{13, 1},
			"0000 0010 10": dctCoefficient{0, 7},
			"0000 0011 00": dctCoefficient{1, 4},
			"0000 0010 11": dctCoefficient{2, 3},
			"0000 0011 11": dctCoefficient{4, 2},
			"0000 0010 01": dctCoefficient{5, 2},
			"0000 0011 10": dctCoefficient{14, 1},
			"0000 0011 01": dctCoefficient{15, 1},
			"0000 0010 00": dctCoefficient{16, 1},

			"0000 0001 1101":   dctCoefficient{0, 8},
			"0000 0001 1000":   dctCoefficient{0, 9},
			"0000 0001 0011":   dctCoefficient{0, 10},
			"0000 0001 0000":   dctCoefficient{0, 11},
			"0000 0001 1011":   dctCoefficient{1, 5},
			"0000 0001 0100":   dctCoefficient{2, 4},
			"0000 0001 1100":   dctCoefficient{3, 3},
			"0000 0001 0010":   dctCoefficient{4, 3},
			"0000 0001 1110":   dctCoefficient{6, 2},
			"0000 0001 0101":   dctCoefficient{7, 2},
			"0000 0001 0001":   dctCoefficient{8, 2},
			"0000 0001 1111":   dctCoefficient{17, 1},
			"0000 0001 1010":   dctCoefficient{18, 1},
			"0000 0001 1001":   dctCoefficient{19, 1},
			"0000 0001 0111":   dctCoefficient{20, 1},
			"0000 0001 0110":   dctCoefficient{21, 1},
			"0000 0000 1101 0": dctCoefficient{0, 12},
			"0000 0000 1100 1": dctCoefficient{0, 13},
			"0000 0000 1100 0": dctCoefficient{0, 14},
			"0000 0000 1011 1": dctCoefficient{0, 15},
			"0000 0000 1011 0": dctCoefficient{1, 6},
			"0000 0000 1010 1": dctCoefficient{1, 7},
			"0000 0000 1010 0": dctCoefficient{2, 5},
			"0000 0000 1001 1": dctCoefficient{3, 4},
			"0000 0000 1001 0": dctCoefficient{5, 3},
			"0000 0000 1000 1": dctCoefficient{9, 2},
			"0000 0000 1000 0": dctCoefficient{10, 2},
			"0000 0000 1111 1": dctCoefficient{22, 1},
			"0000 0000 1111 0": dctCoefficient{23, 1},
			"0000 0000 1110 1": dctCoefficient{24, 1},
			"0000 0000 1110 0": dctCoefficient{25, 1},
			"0000 0000 1101 1": dctCoefficient{26, 1},

			"0000 0000 0111 11":  dctCoefficient{0, 16},
			"0000 0000 0111 10":  dctCoefficient{0, 17},
			"0000 0000 0111 01":  dctCoefficient{0, 18},
			"0000 0000 0111 00":  dctCoefficient{0, 19},
			"0000 0000 0110 11":  dctCoefficient{0, 20},
			"0000 0000 0110 10":  dctCoefficient{0, 21},
			"0000 0000 0110 01":  dctCoefficient{0, 22},
			"0000 0000 0110 00":  dctCoefficient{0, 23},
			"0000 0000 0101 11":  dctCoefficient{0, 24},
			"0000 0000 0101 10":  dctCoefficient{0, 25},
			"0000 0000 0101 01":  dctCoefficient{0, 26},
			"0000 0000 0101 00":  dctCoefficient{0, 27},
			"0000 0000 0100 11":  dctCoefficient{0, 28},
			"0000 0000 0100 10":  dctCoefficient{0, 29},
			"0000 0000 0100 01":  dctCoefficient{0, 30},
			"0000 0000 0100 00":  dctCoefficient{0, 31},
			"0000 0000 0011 000": dctCoefficient{0, 32},
			"0000 0000 0010 111": dctCoefficient{0, 33},
			"0000 0000 0010 110": dctCoefficient{0, 34},
			"0000 0000 0010 101": dctCoefficient{0, 35},
			"0000 0000 0010 100": dctCoefficient{0, 36},
			"0000 0000 0010 011": dctCoefficient{0, 37},
			"0000 0000 0010 010": dctCoefficient{0, 38},
			"0000 0000 0010 001": dctCoefficient{0, 39},
			"0000 0000 0010 000": dctCoefficient{0, 40},
			"0000 0000 0011 111": dctCoefficient{1, 8},
			"0000 0000 0011 110": dctCoefficient{1, 9},
			"0000 0000 0011 101": dctCoefficient{1, 10},
			"0000 0000 0011 100": dctCoefficient{1, 11},
			"0000 0000 0011 011": dctCoefficient{1, 12},
			"0000 0000 0011 010": dctCoefficient{1, 13},
			"0000 0000 0011 001": dctCoefficient{1, 14},

			"0000 0000 0001 0011": dctCoefficient{1, 15},
			"0000 0000 0001 0010": dctCoefficient{1, 16},
			"0000 0000 0001 0001": dctCoefficient{1, 17},
			"0000 0000 0001 0000": dctCoefficient{1, 18},
			"0000 0000 0001 0100": dctCoefficient{6, 3},
			"0000 0000 0001 1010": dctCoefficient{11, 2},
			"0000 0000 0001 1001": dctCoefficient{12, 2},
			"0000 0000 0001 1000": dctCoefficient{13, 2},
			"0000 0000 0001 0111": dctCoefficient{14, 2},
			"0000 0000 0001 0110": dctCoefficient{15, 2},
			"0000 0000 0001 0101": dctCoefficient{16, 2},
			"0000 0000 0001 1111": dctCoefficient{27, 1},
			"0000 0000 0001 1110": dctCoefficient{28, 1},
			"0000 0000 0001 1101": dctCoefficient{29, 1},
			"0000 0000 0001 1100": dctCoefficient{30, 1},
			"0000 0000 0001 1011": dctCoefficient{31, 1},
		}, {
			"10":           dctEndOfBlock{},
			"11":           dctCoefficient{0, 1},
			"011":          dctCoefficient{1, 1},
			"0100":         dctCoefficient{0, 2},
			"0101":         dctCoefficient{2, 1},
			"0010 1":       dctCoefficient{0, 3},
			"0011 1":       dctCoefficient{3, 1},
			"0011 0":       dctCoefficient{4, 1},
			"0001 10":      dctCoefficient{1, 2},
			"0001 11":      dctCoefficient{5, 1},
			"0001 01":      dctCoefficient{6, 1},
			"0001 00":      dctCoefficient{7, 1},
			"0000 110":     dctCoefficient{0, 4},
			"0000 100":     dctCoefficient{2, 2},
			"0000 111":     dctCoefficient{8, 1},
			"0000 101":     dctCoefficient{9, 1},
			"0000 01":      dctEscape{},
			"0010 0110":    dctCoefficient{0, 5},
			"0010 0001":    dctCoefficient{0, 6},
			"0010 0101":    dctCoefficient{1, 3},
			"0010 0100":    dctCoefficient{3, 2},
			"0010 0111":    dctCoefficient{10, 1},
			"0010 0011":    dctCoefficient{11, 1},
			"0010 0010":    dctCoefficient{12, 1},
			"0010 0000":    dctCoefficient{13, 1},
			"0000 0010 10": dctCoefficient{0, 7},
			"0000 0011 00": dctCoefficient{1, 4},
			"0000 0010 11": dctCoefficient{2, 3},
			"0000 0011 11": dctCoefficient{4, 2},
			"0000 0010 01": dctCoefficient{5, 2},
			"0000 0011 10": dctCoefficient{14, 1},
			"0000 0011 01": dctCoefficient{15, 1},
			"0000 0010 00": dctCoefficient{16, 1},

			"0000 0001 1101":   dctCoefficient{0, 8},
			"0000 0001 1000":   dctCoefficient{0, 9},
			"0000 0001 0011":   dctCoefficient{0, 10},
			"0000 0001 0000":   dctCoefficient{0, 11},
			"0000 0001 1011":   dctCoefficient{1, 5},
			"0000 0001 0100":   dctCoefficient{2, 4},
			"0000 0001 1100":   dctCoefficient{3, 3},
			"0000 0001 0010":   dctCoefficient{4, 3},
			"0000 0001 1110":   dctCoefficient{6, 2},
			"0000 0001 0101":   dctCoefficient{7, 2},
			"0000 0001 0001":   dctCoefficient{8, 2},
			"0000 0001 1111":   dctCoefficient{17, 1},
			"0000 0001 1010":   dctCoefficient{18, 1},
			"0000 0001 1001":   dctCoefficient{19, 1},
			"0000 0001 0111":   dctCoefficient{20, 1},
			"0000 0001 0110":   dctCoefficient{21, 1},
			"0000 0000 1101 0": dctCoefficient{0, 12},
			"0000 0000 1100 1": dctCoefficient{0, 13},
			"0000 0000 1100 0": dctCoefficient{0, 14},
			"0000 0000 1011 1": dctCoefficient{0, 15},
			"0000 0000 1011 0": dctCoefficient{1, 6},
			"0000 0000 1010 1": dctCoefficient{1, 7},
			"0000 0000 1010 0": dctCoefficient{2, 5},
			"0000 0000 1001 1": dctCoefficient{3, 4},
			"0000 0000 1001 0": dctCoefficient{5, 3},
			"0000 0000 1000 1": dctCoefficient{9, 2},
			"0000 0000 1000 0": dctCoefficient{10, 2},
			"0000 0000 1111 1": dctCoefficient{22, 1},
			"0000 0000 1111 0": dctCoefficient{23, 1},
			"0000 0000 1110 1": dctCoefficient{24, 1},
			"0000 0000 1110 0": dctCoefficient{25, 1},
			"0000 0000 1101 1": dctCoefficient{26, 1},

			"0000 0000 0111 11":  dctCoefficient{0, 16},
			"0000 0000 0111 10":  dctCoefficient{0, 17},
			"0000 0000 0111 01":  dctCoefficient{0, 18},
			"0000 0000 0111 00":  dctCoefficient{0, 19},
			"0000 0000 0110 11":  dctCoefficient{0, 20},
			"0000 0000 0110 10":  dctCoefficient{0, 21},
			"0000 0000 0110 01":  dctCoefficient{0, 22},
			"0000 0000 0110 00":  dctCoefficient{0, 23},
			"0000 0000 0101 11":  dctCoefficient{0, 24},
			"0000 0000 0101 10":  dctCoefficient{0, 25},
			"0000 0000 0101 01":  dctCoefficient{0, 26},
			"0000 0000 0101 00":  dctCoefficient{0, 27},
			"0000 0000 0100 11":  dctCoefficient{0, 28},
			"0000 0000 0100 10":  dctCoefficient{0, 29},
			"0000 0000 0100 01":  dctCoefficient{0, 30},
			"0000 0000 0100 00":  dctCoefficient{0, 31},
			"0000 0000 0011 000": dctCoefficient{0, 32},
			"0000 0000 0010 111": dctCoefficient{0, 33},
			"0000 0000 0010 110": dctCoefficient{0, 34},
			"0000 0000 0010 101": dctCoefficient{0, 35},
			"0000 0000 0010 100": dctCoefficient{0, 36},
			"0000 0000 0010 011": dctCoefficient{0, 37},
			"0000 0000 0010 010": dctCoefficient{0, 38},
			"0000 0000 0010 001": dctCoefficient{0, 39},
			"0000 0000 0010 000": dctCoefficient{0, 40},
			"0000 0000 0011 111": dctCoefficient{1, 8},
			"0000 0000 0011 110": dctCoefficient{1, 9},
			"0000 0000 0011 101": dctCoefficient{1, 10},
			"0000 0000 0011 100": dctCoefficient{1, 11},
			"0000 0000 0011 011": dctCoefficient{1, 12},
			"0000 0000 0011 010": dctCoefficient{1, 13},
			"0000 0000 0011 001": dctCoefficient{1, 14},

			"0000 0000 0001 0011": dctCoefficient{1, 15},
			"0000 0000 0001 0010": dctCoefficient{1, 16},
			"0000 0000 0001 0001": dctCoefficient{1, 17},
			"0000 0000 0001 0000": dctCoefficient{1, 18},
			"0000 0000 0001 0100": dctCoefficient{6, 3},
			"0000 0000 0001 1010": dctCoefficient{11, 2},
			"0000 0000 0001 1001": dctCoefficient{12, 2},
			"0000 0000 0001 1000": dctCoefficient{13, 2},
			"0000 0000 0001 0111": dctCoefficient{14, 2},
			"0000 0000 0001 0110": dctCoefficient{15, 2},
			"0000 0000 0001 0101": dctCoefficient{16, 2},
			"0000 0000 0001 1111": dctCoefficient{27, 1},
			"0000 0000 0001 1110": dctCoefficient{28, 1},
			"0000 0000 0001 1101": dctCoefficient{29, 1},
			"0000 0000 0001 1100": dctCoefficient{30, 1},
			"0000 0000 0001 1011": dctCoefficient{31, 1},
		},
	},
	{
		{
			"10":           dctCoefficient{0, 1},
			"010":          dctCoefficient{1, 1},
			"110":          dctCoefficient{0, 2},
			"0010 1":       dctCoefficient{2, 1},
			"0111":         dctCoefficient{0, 3},
			"0011 1":       dctCoefficient{3, 1},
			"0001 10":      dctCoefficient{4, 1},
			"0011 0":       dctCoefficient{1, 2},
			"0001 11":      dctCoefficient{5, 1},
			"0000 110":     dctCoefficient{6, 1},
			"0000 100":     dctCoefficient{7, 1},
			"1110 0":       dctCoefficient{0, 4},
			"0000 111":     dctCoefficient{2, 2},
			"0000 101":     dctCoefficient{8, 1},
			"1111 000":     dctCoefficient{9, 1},
			"0000 01":      dctEscape{},
			"1110 1":       dctCoefficient{0, 5},
			"0001 01":      dctCoefficient{0, 6},
			"1111 001":     dctCoefficient{1, 3},
			"0010 0110":    dctCoefficient{3, 2},
			"1111 010":     dctCoefficient{10, 1},
			"0010 0001":    dctCoefficient{11, 1},
			"0010 0101":    dctCoefficient{12, 1},
			"0010 0100":    dctCoefficient{13, 1},
			"0001 00":      dctCoefficient{0, 7},
			"0010 0111":    dctCoefficient{1, 4},
			"1111 1100":    dctCoefficient{2, 3},
			"1111 1101":    dctCoefficient{4, 2},
			"0000 0010 0":  dctCoefficient{5, 2},
			"0000 0010 1":  dctCoefficient{14, 1},
			"0000 0011 1":  dctCoefficient{15, 1},
			"0000 0011 01": dctCoefficient{16, 1},

			"1111 011":         dctCoefficient{0, 8},
			"1111 100":         dctCoefficient{0, 9},
			"0010 0011":        dctCoefficient{0, 10},
			"0010 0010":        dctCoefficient{0, 11},
			"0010 0000":        dctCoefficient{1, 5},
			"0000 0011 00":     dctCoefficient{2, 4},
			"0000 0001 1100":   dctCoefficient{3, 3},
			"0000 0001 0010":   dctCoefficient{4, 3},
			"0000 0001 1110":   dctCoefficient{6, 2},
			"0000 0001 0101":   dctCoefficient{7, 2},
			"0000 0001 0001":   dctCoefficient{8, 2},
			"0000 0001 1111":   dctCoefficient{17, 1},
			"0000 0001 1010":   dctCoefficient{18, 1},
			"0000 0001 1001":   dctCoefficient{19, 1},
			"0000 0001 0111":   dctCoefficient{20, 1},
			"0000 0001 0110":   dctCoefficient{21, 1},
			"1111 1010":        dctCoefficient{0, 12},
			"1111 1011":        dctCoefficient{0, 13},
			"1111 1110":        dctCoefficient{0, 14},
			"1111 1111":        dctCoefficient{0, 15},
			"0000 0000 1011 0": dctCoefficient{1, 6},
			"0000 0000 1010 1": dctCoefficient{1, 7},
			"0000 0000 1010 0": dctCoefficient{2, 5},
			"0000 0000 1001 1": dctCoefficient{3, 4},
			"0000 0000 1001 0": dctCoefficient{5, 3},
			"0000 0000 1000 1": dctCoefficient{9, 2},
			"0000 0000 1000 0": dctCoefficient{10, 2},
			"0000 0000 1111 1": dctCoefficient{22, 1},
			"0000 0000 1111 0": dctCoefficient{23, 1},
			"0000 0000 1110 1": dctCoefficient{24, 1},
			"0000 0000 1110 0": dctCoefficient{25, 1},
			"0000 0000 1101 1": dctCoefficient{26, 1},

			"0000 0000 0111 11":  dctCoefficient{0, 16},
			"0000 0000 0111 10":  dctCoefficient{0, 17},
			"0000 0000 0111 01":  dctCoefficient{0, 18},
			"0000 0000 0111 00":  dctCoefficient{0, 19},
			"0000 0000 0110 11":  dctCoefficient{0, 20},
			"0000 0000 0110 10":  dctCoefficient{0, 21},
			"0000 0000 0110 01":  dctCoefficient{0, 22},
			"0000 0000 0110 00":  dctCoefficient{0, 23},
			"0000 0000 0101 11":  dctCoefficient{0, 24},
			"0000 0000 0101 10":  dctCoefficient{0, 25},
			"0000 0000 0101 01":  dctCoefficient{0, 26},
			"0000 0000 0101 00":  dctCoefficient{0, 27},
			"0000 0000 0100 11":  dctCoefficient{0, 28},
			"0000 0000 0100 10":  dctCoefficient{0, 29},
			"0000 0000 0100 01":  dctCoefficient{0, 30},
			"0000 0000 0100 00":  dctCoefficient{0, 31},
			"0000 0000 0011 000": dctCoefficient{0, 32},
			"0000 0000 0010 111": dctCoefficient{0, 33},
			"0000 0000 0010 110": dctCoefficient{0, 34},
			"0000 0000 0010 101": dctCoefficient{0, 35},
			"0000 0000 0010 100": dctCoefficient{0, 36},
			"0000 0000 0010 011": dctCoefficient{0, 37},
			"0000 0000 0010 010": dctCoefficient{0, 38},
			"0000 0000 0010 001": dctCoefficient{0, 39},
			"0000 0000 0010 000": dctCoefficient{0, 40},
			"0000 0000 0011 111": dctCoefficient{1, 8},
			"0000 0000 0011 110": dctCoefficient{1, 9},
			"0000 0000 0011 101": dctCoefficient{1, 10},
			"0000 0000 0011 100": dctCoefficient{1, 11},
			"0000 0000 0011 011": dctCoefficient{1, 12},
			"0000 0000 0011 010": dctCoefficient{1, 13},
			"0000 0000 0011 001": dctCoefficient{1, 14},

			"0000 0000 0001 0011": dctCoefficient{1, 15},
			"0000 0000 0001 0010": dctCoefficient{1, 16},
			"0000 0000 0001 0001": dctCoefficient{1, 17},
			"0000 0000 0001 0000": dctCoefficient{1, 18},
			"0000 0000 0001 0100": dctCoefficient{6, 3},
			"0000 0000 0001 1010": dctCoefficient{11, 2},
			"0000 0000 0001 1001": dctCoefficient{12, 2},
			"0000 0000 0001 1000": dctCoefficient{13, 2},
			"0000 0000 0001 0111": dctCoefficient{14, 2},
			"0000 0000 0001 0110": dctCoefficient{15, 2},
			"0000 0000 0001 0101": dctCoefficient{16, 2},
			"0000 0000 0001 1111": dctCoefficient{27, 1},
			"0000 0000 0001 1110": dctCoefficient{28, 1},
			"0000 0000 0001 1101": dctCoefficient{29, 1},
			"0000 0000 0001 1100": dctCoefficient{30, 1},
			"0000 0000 0001 1011": dctCoefficient{31, 1},
		}, {
			"0110":         dctEndOfBlock{},
			"10":           dctCoefficient{0, 1},
			"010":          dctCoefficient{1, 1},
			"110":          dctCoefficient{0, 2},
			"0010 1":       dctCoefficient{2, 1},
			"0111":         dctCoefficient{0, 3},
			"0011 1":       dctCoefficient{3, 1},
			"0001 10":      dctCoefficient{4, 1},
			"0011 0":       dctCoefficient{1, 2},
			"0001 11":      dctCoefficient{5, 1},
			"0000 110":     dctCoefficient{6, 1},
			"0000 100":     dctCoefficient{7, 1},
			"1110 0":       dctCoefficient{0, 4},
			"0000 111":     dctCoefficient{2, 2},
			"0000 101":     dctCoefficient{8, 1},
			"1111 000":     dctCoefficient{9, 1},
			"0000 01":      dctEscape{},
			"1110 1":       dctCoefficient{0, 5},
			"0001 01":      dctCoefficient{0, 6},
			"1111 001":     dctCoefficient{1, 3},
			"0010 0110":    dctCoefficient{3, 2},
			"1111 010":     dctCoefficient{10, 1},
			"0010 0001":    dctCoefficient{11, 1},
			"0010 0101":    dctCoefficient{12, 1},
			"0010 0100":    dctCoefficient{13, 1},
			"0001 00":      dctCoefficient{0, 7},
			"0010 0111":    dctCoefficient{1, 4},
			"1111 1100":    dctCoefficient{2, 3},
			"1111 1101":    dctCoefficient{4, 2},
			"0000 0010 0":  dctCoefficient{5, 2},
			"0000 0010 1":  dctCoefficient{14, 1},
			"0000 0011 1":  dctCoefficient{15, 1},
			"0000 0011 01": dctCoefficient{16, 1},

			"1111 011":         dctCoefficient{0, 8},
			"1111 100":         dctCoefficient{0, 9},
			"0010 0011":        dctCoefficient{0, 10},
			"0010 0010":        dctCoefficient{0, 11},
			"0010 0000":        dctCoefficient{1, 5},
			"0000 0011 00":     dctCoefficient{2, 4},
			"0000 0001 1100":   dctCoefficient{3, 3},
			"0000 0001 0010":   dctCoefficient{4, 3},
			"0000 0001 1110":   dctCoefficient{6, 2},
			"0000 0001 0101":   dctCoefficient{7, 2},
			"0000 0001 0001":   dctCoefficient{8, 2},
			"0000 0001 1111":   dctCoefficient{17, 1},
			"0000 0001 1010":   dctCoefficient{18, 1},
			"0000 0001 1001":   dctCoefficient{19, 1},
			"0000 0001 0111":   dctCoefficient{20, 1},
			"0000 0001 0110":   dctCoefficient{21, 1},
			"1111 1010":        dctCoefficient{0, 12},
			"1111 1011":        dctCoefficient{0, 13},
			"1111 1110":        dctCoefficient{0, 14},
			"1111 1111":        dctCoefficient{0, 15},
			"0000 0000 1011 0": dctCoefficient{1, 6},
			"0000 0000 1010 1": dctCoefficient{1, 7},
			"0000 0000 1010 0": dctCoefficient{2, 5},
			"0000 0000 1001 1": dctCoefficient{3, 4},
			"0000 0000 1001 0": dctCoefficient{5, 3},
			"0000 0000 1000 1": dctCoefficient{9, 2},
			"0000 0000 1000 0": dctCoefficient{10, 2},
			"0000 0000 1111 1": dctCoefficient{22, 1},
			"0000 0000 1111 0": dctCoefficient{23, 1},
			"0000 0000 1110 1": dctCoefficient{24, 1},
			"0000 0000 1110 0": dctCoefficient{25, 1},
			"0000 0000 1101 1": dctCoefficient{26, 1},

			"0000 0000 0111 11":  dctCoefficient{0, 16},
			"0000 0000 0111 10":  dctCoefficient{0, 17},
			"0000 0000 0111 01":  dctCoefficient{0, 18},
			"0000 0000 0111 00":  dctCoefficient{0, 19},
			"0000 0000 0110 11":  dctCoefficient{0, 20},
			"0000 0000 0110 10":  dctCoefficient{0, 21},
			"0000 0000 0110 01":  dctCoefficient{0, 22},
			"0000 0000 0110 00":  dctCoefficient{0, 23},
			"0000 0000 0101 11":  dctCoefficient{0, 24},
			"0000 0000 0101 10":  dctCoefficient{0, 25},
			"0000 0000 0101 01":  dctCoefficient{0, 26},
			"0000 0000 0101 00":  dctCoefficient{0, 27},
			"0000 0000 0100 11":  dctCoefficient{0, 28},
			"0000 0000 0100 10":  dctCoefficient{0, 29},
			"0000 0000 0100 01":  dctCoefficient{0, 30},
			"0000 0000 0100 00":  dctCoefficient{0, 31},
			"0000 0000 0011 000": dctCoefficient{0, 32},
			"0000 0000 0010 111": dctCoefficient{0, 33},
			"0000 0000 0010 110": dctCoefficient{0, 34},
			"0000 0000 0010 101": dctCoefficient{0, 35},
			"0000 0000 0010 100": dctCoefficient{0, 36},
			"0000 0000 0010 011": dctCoefficient{0, 37},
			"0000 0000 0010 010": dctCoefficient{0, 38},
			"0000 0000 0010 001": dctCoefficient{0, 39},
			"0000 0000 0010 000": dctCoefficient{0, 40},
			"0000 0000 0011 111": dctCoefficient{1, 8},
			"0000 0000 0011 110": dctCoefficient{1, 9},
			"0000 0000 0011 101": dctCoefficient{1, 10},
			"0000 0000 0011 100": dctCoefficient{1, 11},
			"0000 0000 0011 011": dctCoefficient{1, 12},
			"0000 0000 0011 010": dctCoefficient{1, 13},
			"0000 0000 0011 001": dctCoefficient{1, 14},

			"0000 0000 0001 0011": dctCoefficient{1, 15},
			"0000 0000 0001 0010": dctCoefficient{1, 16},
			"0000 0000 0001 0001": dctCoefficient{1, 17},
			"0000 0000 0001 0000": dctCoefficient{1, 18},
			"0000 0000 0001 0100": dctCoefficient{6, 3},
			"0000 0000 0001 1010": dctCoefficient{11, 2},
			"0000 0000 0001 1001": dctCoefficient{12, 2},
			"0000 0000 0001 1000": dctCoefficient{13, 2},
			"0000 0000 0001 0111": dctCoefficient{14, 2},
			"0000 0000 0001 0110": dctCoefficient{15, 2},
			"0000 0000 0001 0101": dctCoefficient{16, 2},
			"0000 0000 0001 1111": dctCoefficient{27, 1},
			"0000 0000 0001 1110": dctCoefficient{28, 1},
			"0000 0000 0001 1101": dctCoefficient{29, 1},
			"0000 0000 0001 1100": dctCoefficient{30, 1},
			"0000 0000 0001 1011": dctCoefficient{31, 1},
		},
	},
}

type dctDCSizeDecoderFn func(br bitreader.BitReader) (uint, error)

func newDCTDCSizeDecoder(table huffman.HuffmanTable) dctDCSizeDecoderFn {
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

var dctDCSizeDecoders = struct {
	Luma   dctDCSizeDecoderFn
	Chroma dctDCSizeDecoderFn
}{
	newDCTDCSizeDecoder(dctDcSizeLuminanceTable),
	newDCTDCSizeDecoder(dctDcSizeChrominanceTable),
}

var dctCoefficientDecoders = struct {
	TableZero dctCoefficientDecoderFn
	TableOne  dctCoefficientDecoderFn
}{

	newDCTCoefficientDecoder(dctCoefficientTables[0]),
	newDCTCoefficientDecoder(dctCoefficientTables[1]),
}
