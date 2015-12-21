package video

import "github.com/32bitkid/bitreader"

type GroupOfPicturesHeader struct {
	time_code   uint32 // 25 bslbf
	closed_gop  bool   // 1 uimsbf
	broken_link bool   // 1 uimsbf

}

func group_of_pictures_header(br bitreader.BitReader) (*GroupOfPicturesHeader, error) {
	err := start_code_check(br, GroupStartCode)
	if err != nil {
		return nil, err
	}

	goph := GroupOfPicturesHeader{}
	goph.time_code, err = br.Read32(25)
	if err != nil {
		return nil, err
	}

	goph.closed_gop, err = br.ReadBit()
	if err != nil {
		return nil, err
	}

	goph.broken_link, err = br.ReadBit()
	if err != nil {
		return nil, err
	}

	return &goph, next_start_code(br)
}
