package video

import "github.com/32bitkid/mpeg/util"

type PictureCodingExtension struct {
	fcode                      [4][4]byte
	intra_dc_precision         uint32 // 2 uimsbf
	picture_structure          uint32 // 2 uimsbf
	top_field_first            bool   // 1 uimsbf
	frame_pred_frame_dct       bool   // 1 uimsbf
	concealment_motion_vectors bool   // 1 uimsbf
	q_scale_type               bool   // 1 uimsbf
	intra_vlc_format           bool   // 1 uimsbf
	alternate_scan             bool   // 1 uimsbf
	repeat_first_field         bool   // 1 uimsbf
	chroma_420_type            bool   // 1 uimsbf
	progressive_frame          bool   // 1 uimsbf
	composite_display_flag     bool   // 1 uimsbf

	v_axis            bool   // 1 uimsbf
	field_sequence    uint32 // 3 uimsbf
	sub_carrier       bool   // 1 uimsbf
	burst_amplitude   uint32 // 7 uimsbf
	sub_carrier_phase uint32 // 8 uimsbf
}

func picture_coding_extension(br util.BitReader32) (*PictureCodingExtension, error) {

	err := start_code_check(br, ExtensionStartCode)
	if err != nil {
		return nil, err
	}

	err = extension_code_check(br, PictureCodingExtensionID)
	if err != nil {
		return nil, err
	}

	// extension_start_code 32 bslbf
	// extension_start_code_identifier 4 uimsbf
	// f_code[0][0] /* forward horizontal */ 4 uimsbf
	// f_code[0][1] /* forward vertical */ 4 uimsbf
	// f_code[1][0] /* backward horizontal */ 4 uimsbf
	// f_code[1][1] /* backward vertical */ 4 uimsbf
	// intra_dc_precision 2 uimsbf
	// picture_structure 2 uimsbf
	// top_field_first 1 uimsbf
	// frame_pred_frame_dct 1 uimsbf
	// concealment_motion_vectors 1 uimsbf
	// q_scale_type 1 uimsbf
	// intra_vlc_format 1 uimsbf
	// alternate_scan 1 uimsbf
	// repeat_first_field 1 uimsbf
	// chroma_420_type 1 uimsbf
	// progressive_frame 1 uimsbf
	// composite_display_flag 1 uimsbf
	// if ( composite_display_flag ) {
	// v_axis 1 uimsbf
	// field_sequence 3 uimsbf
	// sub_carrier 1 uimsbf
	// burst_amplitude 7 uimsbf
	// sub_carrier_phase 8 uimsbf
	// }
	next_start_code(br)

	panic("not supported: picture_coding_extension")
}
