package video

import "github.com/32bitkid/bitreader"
import "errors"

// StartCode is a 32 bit code that acts as a marker in a coded bitstream.
// They usually signal the structure of following bits and/or how the bits
// should be interpreted.
//
// Start codes always begin with the 24 bit integer 0x000001 followed by a
// specific identifier.
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

// ErrUnexpectedStartCode indicates that a start code was read from the bitstream that was unexpected.
var ErrUnexpectedStartCode = errors.New("unexpected start code")

// IsSlice() returns true if the StartCode falls within the
// acceptable range of codes designated as slice start codes.
func (code StartCode) IsSlice() bool {
	return code >= MinSliceStartCode && code <= MaxSliceStartCode
}

// Check() will return true if the next bits in the bitstream match the expected code.
// Check() does not consume any bits from the bitstream and will only return
// an error if there is a underlying error attempting to peek into the bitstream.
func (expected StartCode) Check(br bitreader.BitReader) (bool, error) {
	if nextbits, err := br.Peek32(32); err != nil {
		return false, err
	} else {
		return StartCode(nextbits) == expected, nil
	}
}

// Assert() returns an ErrUnexpectedStartCode if the next bits in the bitstream do not match the expected code.
// If the expected code is present, the the bits are consumed from the bitstream.
func (expected StartCode) Assert(br bitreader.BitReader) error {
	if test, err := expected.Check(br); err != nil {
		return err
	} else if test != true {
		return ErrUnexpectedStartCode
	}
	if err := br.Trash(32); err != nil {
		return err
	}
	return nil
}
