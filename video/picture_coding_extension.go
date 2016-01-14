package video

import "github.com/32bitkid/bitreader"

type FCode [2][2]uint32

type PictureCodingExtension struct {
	f_code                     FCode
	intra_dc_precision         uint32           // 2 uimsbf
	picture_structure          PictureStructure // 2 uimsbf
	top_field_first            bool             // 1 uimsbf
	frame_pred_frame_dct       uint32           // 1 uimsbf
	concealment_motion_vectors bool             // 1 uimsbf
	q_scale_type               uint32           // 1 uimsbf
	intra_vlc_format           uint32           // 1 uimsbf
	alternate_scan             uint32           // 1 uimsbf
	repeat_first_field         bool             // 1 uimsbf
	chroma_420_type            bool             // 1 uimsbf
	progressive_frame          bool             // 1 uimsbf
	composite_display_flag     bool             // 1 uimsbf

	v_axis            bool   // 1 uimsbf
	field_sequence    uint32 // 3 uimsbf
	sub_carrier       bool   // 1 uimsbf
	burst_amplitude   uint32 // 7 uimsbf
	sub_carrier_phase uint32 // 8 uimsbf
}

func picture_coding_extension(br bitreader.BitReader) (*PictureCodingExtension, error) {

	err := ExtensionStartCode.Assert(br)
	if err != nil {
		return nil, err
	}

	err = PictureCodingExtensionID.Assert(br)
	if err != nil {
		return nil, err
	}

	pce := PictureCodingExtension{}

	pce.f_code[0][0], err = br.Read32(4)
	if err != nil {
		return nil, err
	}
	pce.f_code[0][1], err = br.Read32(4)
	if err != nil {
		return nil, err
	}
	pce.f_code[1][0], err = br.Read32(4)
	if err != nil {
		return nil, err
	}
	pce.f_code[1][1], err = br.Read32(4)
	if err != nil {
		return nil, err
	}

	pce.intra_dc_precision, err = br.Read32(2)
	if err != nil {
		return nil, err
	}

	if picture_structure, err := br.Read32(2); err != nil {
		return nil, err
	} else {
		pce.picture_structure = PictureStructure(picture_structure)
	}

	pce.top_field_first, err = br.ReadBit()
	if err != nil {
		return nil, err
	}

	pce.frame_pred_frame_dct, err = br.Read32(1)
	if err != nil {
		return nil, err
	}

	pce.concealment_motion_vectors, err = br.ReadBit()
	if err != nil {
		return nil, err
	}

	pce.q_scale_type, err = br.Read32(1)
	if err != nil {
		return nil, err
	}

	pce.intra_vlc_format, err = br.Read32(1)
	if err != nil {
		return nil, err
	}

	pce.alternate_scan, err = br.Read32(1)
	if err != nil {
		return nil, err
	}

	pce.repeat_first_field, err = br.ReadBit()
	if err != nil {
		return nil, err
	}

	pce.chroma_420_type, err = br.ReadBit()
	if err != nil {
		return nil, err
	}

	pce.progressive_frame, err = br.ReadBit()
	if err != nil {
		return nil, err
	}

	pce.composite_display_flag, err = br.ReadBit()
	if err != nil {
		return nil, err
	}

	if pce.composite_display_flag {
		pce.v_axis, err = br.ReadBit()
		if err != nil {
			return nil, err
		}
		pce.field_sequence, err = br.Read32(3)
		if err != nil {
			return nil, err
		}
		pce.sub_carrier, err = br.ReadBit()
		if err != nil {
			return nil, err
		}
		pce.burst_amplitude, err = br.Read32(7)
		if err != nil {
			return nil, err
		}
		pce.sub_carrier_phase, err = br.Read32(8)
		if err != nil {
			return nil, err
		}
	}

	return &pce, next_start_code(br)
}
