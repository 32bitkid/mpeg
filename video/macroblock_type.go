package video

import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/huffman"

type macroblockTypeDecoder func(bitreader.BitReader) (*MacroblockType, error)

func newMacroblockTypeDecoder(table huffman.HuffmanTable) macroblockTypeDecoder {
	decoder := huffman.NewHuffmanDecoder(table)
	return func(br bitreader.BitReader) (*MacroblockType, error) {
		val, err := decoder.Decode(br)
		if err != nil {
			return nil, err
		} else if mb_type, ok := val.(*MacroblockType); ok {
			return mb_type, nil
		} else {
			return nil, ErrUnexpectedDecodedValueType
		}
	}
}

type SpatialTemporalWeightClass int

const (
	SpatialTemporalWeightClass_0 = 1 << iota
	SpatialTemporalWeightClass_1
	SpatialTemporalWeightClass_2
	SpatialTemporalWeightClass_3
	SpatialTemporalWeightClass_4
)

type MacroblockType struct {
	macroblock_quant                  bool
	macroblock_motion_forward         bool
	macroblock_motion_backward        bool
	macroblock_pattern                bool
	macroblock_intra                  bool
	spatial_temporal_weight_code_flag bool
	spatial_temporal_weight_classes   SpatialTemporalWeightClass
}

var iFrameMacroblockTypesTable = huffman.HuffmanTable{
	"1":  &MacroblockType{false, false, false, false, true, false, SpatialTemporalWeightClass_0},
	"01": &MacroblockType{true, false, false, false, true, false, SpatialTemporalWeightClass_0},
}

var pFrameMacroblockTypesTable = huffman.HuffmanTable{
	"1":       &MacroblockType{false, true, false, true, false, false, SpatialTemporalWeightClass_0},
	"01":      &MacroblockType{false, false, false, true, false, false, SpatialTemporalWeightClass_0},
	"001":     &MacroblockType{false, true, false, false, false, false, SpatialTemporalWeightClass_0},
	"0001 1":  &MacroblockType{false, false, false, false, true, false, SpatialTemporalWeightClass_0},
	"0001 0":  &MacroblockType{true, true, false, true, false, false, SpatialTemporalWeightClass_0},
	"0000 1":  &MacroblockType{true, false, false, true, false, false, SpatialTemporalWeightClass_0},
	"0000 01": &MacroblockType{true, false, false, false, true, false, SpatialTemporalWeightClass_0},
}

var bFrameMacroblockTypesTable = huffman.HuffmanTable{
	"10":      &MacroblockType{false, true, true, false, false, false, SpatialTemporalWeightClass_0},
	"11":      &MacroblockType{false, true, true, true, false, false, SpatialTemporalWeightClass_0},
	"010":     &MacroblockType{false, false, true, false, false, false, SpatialTemporalWeightClass_0},
	"011":     &MacroblockType{false, false, true, true, false, false, SpatialTemporalWeightClass_0},
	"0010":    &MacroblockType{false, true, false, false, false, false, SpatialTemporalWeightClass_0},
	"0011":    &MacroblockType{false, true, false, true, false, false, SpatialTemporalWeightClass_0},
	"0001 1":  &MacroblockType{false, false, false, false, true, false, SpatialTemporalWeightClass_0},
	"0001 0":  &MacroblockType{true, true, true, true, false, false, SpatialTemporalWeightClass_0},
	"0000 11": &MacroblockType{true, true, false, true, false, false, SpatialTemporalWeightClass_0},
	"0000 10": &MacroblockType{true, false, true, true, false, false, SpatialTemporalWeightClass_0},
	"0000 01": &MacroblockType{true, false, false, false, true, false, SpatialTemporalWeightClass_0},
}

var MacroblockTypeDecoder = struct {
	IFrame macroblockTypeDecoder
	PFrame macroblockTypeDecoder
	BFrame macroblockTypeDecoder
}{
	newMacroblockTypeDecoder(iFrameMacroblockTypesTable),
	newMacroblockTypeDecoder(pFrameMacroblockTypesTable),
	newMacroblockTypeDecoder(bFrameMacroblockTypesTable),
}
