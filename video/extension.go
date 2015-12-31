package video

import "github.com/32bitkid/bitreader"
import "errors"

type ExtensionID uint32

const (
	_                                  ExtensionID = iota // reserved
	SequenceExtensionID                                   //
	SequenceDisplayExtensionID                            //
	QuantMatrixExtensionID                                //
	CopyrightExtensionID                                  //
	SequenceScalableExtensionID                           //
	_                                                     // reserved
	PictureDisplayExtensionID                             //
	PictureCodingExtensionID                              //
	PictureSpatialScalableExtensionID                     //
	PictureTemporalScalableExtensionID                    //
	_                                                     // reserved
	_                                                     // reserved
	_                                                     // reserved
	_                                                     // reserved
)

var ErrUnexpectedSequenceExtensionID = errors.New("unexpected sequence extension id")

func extension_code_check(br bitreader.BitReader, expected ExtensionID) error {
	actual, err := br.Read32(4)
	if err != nil {
		return err
	}
	if ExtensionID(actual) != expected {
		return ErrUnexpectedSequenceExtensionID
	}
	return nil
}
