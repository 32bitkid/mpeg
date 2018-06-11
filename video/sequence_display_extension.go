package video

import "github.com/32bitkid/bitreader"

type SequenceDisplayExtension struct {
	video_format             uint32 // 3 uimsbf
	colour_description       bool   // 1 uimsbf
	colour_primaries         uint32 // 8 uimsbf
	transfer_characteristics uint32 // 8 uimsbf
	matrix_coefficients      uint32 // 8 uimsbf
	display_horizontal_size  uint32 // 14 uimsbf
	display_vertical_size    uint32 // 14 uimsbf
}

func sequence_display_extension(br bitreader.BitReader) (*SequenceDisplayExtension, error) {

	err := SequenceDisplayExtensionID.Assert(br)
	if err != nil {
		return nil, err
	}

	sde := SequenceDisplayExtension{}
	sde.video_format, err = br.Read32(3)
	if err != nil {
		return nil, err
	}

	sde.colour_description, err = br.Read1()
	if sde.colour_description {
		sde.colour_primaries, err = br.Read32(8)
		if err != nil {
			return nil, err
		}

		sde.transfer_characteristics, err = br.Read32(8)
		if err != nil {
			return nil, err
		}

		sde.matrix_coefficients, err = br.Read32(8)
		if err != nil {
			return nil, err
		}
	}

	sde.display_horizontal_size, err = br.Read32(14)
	if err != nil {
		return nil, err
	}

	err = marker_bit(br)
	if err != nil {
		return nil, err
	}

	sde.display_vertical_size, err = br.Read32(14)
	if err != nil {
		return nil, err
	}

	return &sde, next_start_code(br)

}
