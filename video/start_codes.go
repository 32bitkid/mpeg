package video

import "github.com/32bitkid/bitreader"
import "fmt"

type StartCode uint32

const (
	StartCodePrefix = 0x000001

	PictureStartCode StartCode = (StartCodePrefix << 8) | 0x00

	// slice_start_code 01 through AF
	MinSliceStartCode StartCode = (StartCodePrefix << 8) | 0x01
	MaxSliceStartCode StartCode = (StartCodePrefix << 8) | 0xAF

	UserDataStartCode       StartCode = (StartCodePrefix << 8) | 0xB2
	SequenceHeaderStartCode StartCode = (StartCodePrefix << 8) | 0xB3
	ExtensionStartCode      StartCode = (StartCodePrefix << 8) | 0xB5
	SequenceEndStartCode    StartCode = (StartCodePrefix << 8) | 0xB7
	GroupStartCode          StartCode = (StartCodePrefix << 8) | 0xB8
)

func (code StartCode) IsSlice() bool {
	return code >= MinSliceStartCode && code <= MaxSliceStartCode
}

func (code StartCode) String() string {
	return fmt.Sprintf("[%08x]", uint32(code))
}

// Check if the next bits in the bitstream matches the expected code
func (expected StartCode) check(br bitreader.BitReader) (bool, error) {
	if nextbits, err := br.Peek32(32); err != nil {
		return false, err
	} else {
		return StartCode(nextbits) == expected, nil
	}
}

// Assert the next bits in the bit stream match the expected startcode
func (expected StartCode) assert(br bitreader.BitReader) error {
	if test, err := expected.check(br); err != nil {
		return err
	} else if test != true {
		return ErrUnexpectedStartCode
	}
	if err := br.Trash(32); err != nil {
		return err
	}
	return nil
}
