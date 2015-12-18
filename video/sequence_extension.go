package video

import "github.com/32bitkid/mpeg/util"
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

func extension_code_check(br util.BitReader32, expected ExtensionID) error {
	actual, err := br.Read32(4)
	if err != nil {
		return err
	}
	if ExtensionID(actual) != expected {
		return ErrUnexpectedSequenceExtensionID
	}
	return nil
}

const (
	_ uint32 = iota // reserved
	ChromaFormat_4_2_0
	ChromaFormat_4_2_2
	ChromaFormat_4_4_4
)

type SequenceExtension struct {
	extension_start_code            uint32 // 32 bslbf
	extension_start_code_identifier uint32 // 4 uimsbf
	profile_and_level_indication    uint32 // 8 uimsbf
	progressive_sequence            uint32 // 1 uimsbf
	chroma_format                   uint32 // 2 uimsbf
	horizontal_size_extension       uint32 // 2 uimsbf
	vertical_size_extension         uint32 // 2 uimsbf
	bit_rate_extension              uint32 // 12 uimsbf
	marker_bit                      uint32 // 1 bslbf
	vbv_buffer_size_extension       uint32 // 8 uimsbf
	low_delay                       uint32 // 1 uimsbf
	frame_rate_extension_n          uint32 // 2 uimsbf
	frame_rate_extension_d          uint32 // 5 uimsbf
}

func sequence_extension(br util.BitReader32) (*SequenceExtension, error) {

	err := start_code_check(br, ExtensionStartCode)
	if err != nil {
		return nil, err
	}

	err = extension_code_check(br, SequenceExtensionID)
	if err != nil {
		return nil, err
	}

	se := SequenceExtension{}

	se.profile_and_level_indication, err = br.Read32(8)
	if err != nil {
		return nil, err
	}

	se.progressive_sequence, err = br.Read32(1)
	if err != nil {
		return nil, err
	}

	se.chroma_format, err = br.Read32(2)
	if err != nil {
		return nil, err
	}

	se.horizontal_size_extension, err = br.Read32(2)
	if err != nil {
		return nil, err
	}

	se.vertical_size_extension, err = br.Read32(2)
	if err != nil {
		return nil, err
	}

	se.bit_rate_extension, err = br.Read32(12)
	if err != nil {
		return nil, err
	}

	err = marker_bit(br)
	if err != nil {
		return nil, err
	}

	se.vbv_buffer_size_extension, err = br.Read32(8)
	if err != nil {
		return nil, err
	}

	se.low_delay, err = br.Read32(1)
	if err != nil {
		return nil, err
	}

	se.frame_rate_extension_n, err = br.Read32(2)
	if err != nil {
		return nil, err
	}

	se.frame_rate_extension_d, err = br.Read32(5)
	if err != nil {
		return nil, err
	}

	return &se, next_start_code(br)
}
