package video

import "github.com/32bitkid/mpeg/util"

type PictureData struct {
	slices []*Slice
}

func picture_data(br util.BitReader32) (*PictureData, error) {

	pd := PictureData{}

	for {
		s, err := slice(br)
		if err != nil {
			return nil, err
		}

		pd.slices = append(pd.slices, s)

		nextbits, err := br.Peek32(32)
		if err != nil {
			return nil, err
		}
		if is_slice_start_code(StartCode(nextbits)) {
			break
		}
	}
	return &pd, next_start_code(br)
}
