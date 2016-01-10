package video

type motionVectorFormat int

const (
	_ motionVectorFormat = iota
	motionVectorFormat_Field
	motionVectorFormat_Frame
)

type motionVectorInfo struct {
	motion_vector_count  int
	motion_vector_format motionVectorFormat
	dmv                  int
}

func mv_info(fp *VideoSequence, mb *Macroblock) motionVectorInfo {
	if fp.PictureCodingExtension.frame_pred_frame_dct == 1 {
		return motionVectorInfo{1, motionVectorFormat_Frame, 0}
	}

	switch mb.frame_motion_type {
	case 1:
		return motionVectorInfo{2, motionVectorFormat_Field, 0}
	case 2:
		return motionVectorInfo{1, motionVectorFormat_Frame, 0}
	case 3:
		return motionVectorInfo{1, motionVectorFormat_Field, 1}
	}

	switch mb.field_motion_type {
	case 1:
		return motionVectorInfo{1, motionVectorFormat_Field, 0}
	case 2:
		return motionVectorInfo{2, motionVectorFormat_Field, 0}
	case 3:
		return motionVectorInfo{1, motionVectorFormat_Field, 1}
	}

	panic("invalid motion vector state")
}
