package ps

import "errors"
import "github.com/32bitkid/bitreader"

// ErrUnexpectedStartCode indicates that a start code was read from the bitstream that was unexpected.
var ErrUnexpectedStartCode = errors.New("unexpected start code")

type StartCode uint32

const (
	StartCodePrefix = 0x000001

	PackStartCode         StartCode = (StartCodePrefix << 8) | 0xBA
	ProgramEndCode        StartCode = (StartCodePrefix << 8) | 0xB9
	SystemHeaderStartCode StartCode = (StartCodePrefix << 8) | 0xBB
)

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
