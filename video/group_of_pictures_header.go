package video

import "github.com/32bitkid/bitreader"

type GroupOfPicturesHeader struct {
	time_code   uint32 // 25 bslbf
	closed_gop  bool   // 1 uimsbf
	broken_link bool   // 1 uimsbf

}

func group_of_pictures_header(br bitreader.BitReader) (*GroupOfPicturesHeader, error) {

	if err := GroupStartCode.assert(br); err != nil {
		return nil, err
	}

	goph := GroupOfPicturesHeader{}

	if time_code, err := br.Read32(25); err != nil {
		return nil, err
	} else {
		goph.time_code = time_code
	}

	if closed_gop, err := br.ReadBit(); err != nil {
		return nil, err
	} else {
		goph.closed_gop = closed_gop
	}

	if broken_link, err := br.ReadBit(); err != nil {
		return nil, err
	} else {
		goph.broken_link = broken_link
	}

	return &goph, next_start_code(br)
}
