package video

import "github.com/32bitkid/bitreader"
import "errors"

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

var ErrUnexpectedSequenceExtensionID = errors.New("unexpected sequence extension id")

func (expected ExtensionID) check(br bitreader.BitReader) (bool, error) {
	if nextbits, err := br.Peek32(4); err != nil {
		return false, err
	} else {
		return ExtensionID(nextbits) == expected, nil
	}
}

func (expected ExtensionID) assert(br bitreader.BitReader) error {
	if test, err := expected.check(br); err != nil {
		return err
	} else if test != true {
		return ErrUnexpectedSequenceExtensionID
	}
	if err := br.Trash(4); err != nil {
		return err
	}
	return nil
}
