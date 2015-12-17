package video

import "github.com/32bitkid/mpeg/util"

type SequenceDisplayExtension struct {
	video_format             uint32 // 3 uimsbf
	colour_description       bool   // 1 uimsbf
	colour_primaries         uint32 // 8 uimsbf
	transfer_characteristics uint32 // 8 uimsbf
	matrix_coefficients      uint32 // 8 uimsbf
	display_horizontal_size  uint32 // 14 uimsbf
	display_vertical_size    uint32 // 14 uimsbf
}

func sequence_display_extension(br util.BitReader32) (*SequenceDisplayExtension, error) {
	val, err := br.Read32(4)
	if err != nil {
		return nil, err
	} else if ExtensionID(val) != SequenceDisplayExtensionID {
		return nil, ErrUnexpectedSequenceExtensionID
	}

	sde := SequenceDisplayExtension{}
	sde.video_format, err = br.Read32(3)
	if err != nil {
		return nil, err
	}

	sde.colour_description, err = br.ReadBit()
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
