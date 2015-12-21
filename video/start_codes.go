package video

import "github.com/32bitkid/bitreader"

type StartCode uint32

const (
	StuffingByte = 0x00
)

const (
	StartCodePrefix = 0x000001

	PictureCode        = 0x00
	ReservedCode_1     = 0xB0
	ReservedCode_2     = 0xB1
	UserDataCode       = 0xB2
	SequenceHeaderCode = 0xB3
	SequenceErrorCode  = 0xB4
	ExtensionCode      = 0xB5
	ReservedCode_3     = 0xB6
	SequenceEndCode    = 0xB7
	GroupCode          = 0xB8

	SequenceHeaderStartCode StartCode = (StartCodePrefix << 8) + SequenceHeaderCode
	ExtensionStartCode                = (StartCodePrefix << 8) + ExtensionCode
	SequenceEndStartCode              = (StartCodePrefix << 8) + SequenceEndCode
	GroupStartCode                    = (StartCodePrefix << 8) + GroupCode
	PictureStartCode                  = (StartCodePrefix << 8) + PictureCode
	UserDataStartCode                 = (StartCodePrefix << 8) + UserDataCode

	MinSliceStartCode StartCode = (StartCodePrefix << 8) + 0x01
	MaxSliceStartCode StartCode = (StartCodePrefix << 8) + 0xAF
)

func is_slice_start_code(code StartCode) bool {
	return code >= MinSliceStartCode && code <= MaxSliceStartCode
}

// slice_start_code 01 through AF
// system start codes (see note) B9 through FF

func start_code_check(br bitreader.BitReader, expected StartCode) error {
	actual, err := br.Read32(32)
	if err != nil {
		return err
	}
	if StartCode(actual) != expected {
		return ErrUnexpectedStartCode
	}
	return nil
}
