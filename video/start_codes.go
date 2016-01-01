package video

import "github.com/32bitkid/bitreader"

type StartCode uint32

const (
	StuffingByte = 0x00
)

const (
	StartCodePrefix = 0x000001

	pictureCode        = 0x00
	reservedCode1      = 0xB0
	reservedCode2      = 0xB1
	userDataCode       = 0xB2
	sequenceHeaderCode = 0xB3
	sequenceErrorCode  = 0xB4
	extensionCode      = 0xB5
	reservedCode3      = 0xB6
	sequenceEndCode    = 0xB7
	groupCode          = 0xB8

	SequenceHeaderStartCode StartCode = (StartCodePrefix << 8) + sequenceHeaderCode
	ExtensionStartCode      StartCode = (StartCodePrefix << 8) + extensionCode
	SequenceEndStartCode    StartCode = (StartCodePrefix << 8) + sequenceEndCode
	GroupStartCode          StartCode = (StartCodePrefix << 8) + groupCode
	PictureStartCode        StartCode = (StartCodePrefix << 8) + pictureCode
	UserDataStartCode       StartCode = (StartCodePrefix << 8) + userDataCode

	MinSliceStartCode StartCode = (StartCodePrefix << 8) + 0x01
	MaxSliceStartCode StartCode = (StartCodePrefix << 8) + 0xAF
)

func (code StartCode) isSlice() bool {
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
