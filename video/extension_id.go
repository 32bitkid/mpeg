package video

import "github.com/32bitkid/bitreader"
import "errors"

// ExtensionID is a 4 bit code, that immediately follows an ExtensionStartCode,
// used identify the following data.
type ExtensionID uint32

const (
	_ ExtensionID = iota // reserved
	SequenceExtensionID
	SequenceDisplayExtensionID
	QuantMatrixExtensionID
	CopyrightExtensionID
	SequenceScalableExtensionID
	_ // reserved
	PictureDisplayExtensionID
	PictureCodingExtensionID
	PictureSpatialScalableExtensionID
	PictureTemporalScalableExtensionID
	_ // reserved
	_ // reserved
	_ // reserved
	_ // reserved
)

// ErrUnexpectedExtensionID indicates that a Extension ID was read from the bitstream that was unexpected.
var ErrUnexpectedExtensionID = errors.New("unexpected sequence extension id")

// IsReserved() returns true if the extension id is described as "reserved".
func (id ExtensionID) IsReserved() bool {
	return id == ExtensionID(0) ||
		id == ExtensionID(6) ||
		id > ExtensionID(11)
}

// Check() returns true if the following bits in the bitstream match the expected extension id.
// Check() does not consume any bits from the bitstream and will only return
// an error if there is a underlying error attempting to peek into the bitstream.
func (expected ExtensionID) Check(br bitreader.BitReader) (bool, error) {
	if nextbits, err := br.Peek32(4); err != nil {
		return false, err
	} else {
		return ExtensionID(nextbits) == expected, nil
	}
}

// Assert() returns an ErrUnexpectedExtensionID if the following bits in the bitstream do not match
// the expected extension id. If the expected code is present, the the bits
// are consumed from the bitstream.
func (expected ExtensionID) Assert(br bitreader.BitReader) error {
	if test, err := expected.Check(br); err != nil {
		return err
	} else if test != true {
		return ErrUnexpectedExtensionID
	}
	if err := br.Trash(4); err != nil {
		return err
	}
	return nil
}
