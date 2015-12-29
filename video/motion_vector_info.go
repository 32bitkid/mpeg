package video

type MotionVectorFormat int

const (
	_ MotionVectorFormat = iota
	MotionVectorFormat_Field
	MotionVectorFormat_Frame
)

func mv_info(fp *VideoSequence, mb *Macroblock) (int, MotionVectorFormat, int) {
	if fp.PictureCodingExtension.frame_pred_frame_dct == 1 {
		return 1, MotionVectorFormat_Frame, 0
	}

	switch mb.frame_motion_type {
	case 01:
		panic("unsupported: field-based motion type")
	case 02:
		return 1, MotionVectorFormat_Frame, 0
	case 03:
		return 1, MotionVectorFormat_Field, 1
	}

	switch mb.field_motion_type {
	case 01:
		return 1, MotionVectorFormat_Field, 0
	case 02:
		return 2, MotionVectorFormat_Field, 0
	case 03:
		return 1, MotionVectorFormat_Field, 1
	}

	panic("invalid motion vector state")
}
