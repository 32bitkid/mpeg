package video

import "github.com/32bitkid/bitreader"

type PictureCodingType uint32

func (pct PictureCodingType) String() string {
	switch pct {
	case IntraCoded:
		return "I"
	case PredictiveCoded:
		return "P"
	case BidirectionallyPredictiveCoded:
		return "B"
	case DCIntraCoded:
		return "D"
	}
	return string(uint32(pct))
}

const (
	_                              PictureCodingType = iota // 000 forbidden
	IntraCoded                                              // 001
	PredictiveCoded                                         // 010
	BidirectionallyPredictiveCoded                          // 011
	DCIntraCoded                                            // 100 (Not Used in ISO/IEC11172-2)
	_                                                       // 101 reserved
	_                                                       // 110 reserved
	_                                                       // 111 reserved

	IFrame = IntraCoded
	PFrame = PredictiveCoded
	BFrame = BidirectionallyPredictiveCoded
)

type PictureHeader struct {
	temporal_reference       uint32            // 10 uimsbf
	picture_coding_type      PictureCodingType // 3 uimsbf
	vbv_delay                uint32            // 16 uimsbf
	full_pel_forward_vector  uint32            // 1 bslbf
	forward_f_code           uint32            // 3 bslbf
	full_pel_backward_vector uint32            // 1 bslbf
	backward_f_code          uint32            // 3 bslbf

	extra_information []byte
}

func picture_header(br bitreader.BitReader) (*PictureHeader, error) {

	err := start_code_check(br, PictureStartCode)
	if err != nil {
		return nil, err
	}

	ph := PictureHeader{}

	ph.temporal_reference, err = br.Read32(10)
	if err != nil {
		return nil, err
	}

	picture_coding_type, err := br.Read32(3)
	if err != nil {
		return nil, err
	}

	ph.picture_coding_type = PictureCodingType(picture_coding_type)

	ph.vbv_delay, err = br.Read32(16)
	if err != nil {
		return nil, err
	}

	if ph.picture_coding_type == PredictiveCoded || ph.picture_coding_type == BidirectionallyPredictiveCoded {
		ph.full_pel_forward_vector, err = br.Read32(1)
		if err != nil {
			return nil, err
		}

		ph.forward_f_code, err = br.Read32(3)
		if err != nil {
			return nil, err
		}
	}
	if ph.picture_coding_type == BidirectionallyPredictiveCoded {
		ph.full_pel_backward_vector, err = br.Read32(1)
		if err != nil {
			return nil, err
		}

		ph.backward_f_code, err = br.Read32(3)
		if err != nil {
			return nil, err
		}
	}

	for {
		if extraBit, err := br.PeekBit(); err != nil {
			return nil, err
		} else if extraBit == false {
			break
		}

		if err := br.Trash(1); err != nil {
			return nil, err
		}

		// extra_information_picture
		if data, err := br.Read32(8); err != nil {
			return nil, err
		} else {
			ph.extra_information = append(ph.extra_information, byte(data))
		}
	}

	if err := br.Trash(1); err != nil {
		return nil, err
	}

	return &ph, next_start_code(br)
}
