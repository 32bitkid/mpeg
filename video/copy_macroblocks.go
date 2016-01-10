package video

import "image"

func pframe_copy_macroblocks(row, col, n int, dest, src *image.YCbCr) {
	// Copy Y
	{
		y := row * 16
		x := col * 16
		w := n * 16

		for v := 0; v < 16; v++ {
			si := ((y + v) * src.YStride) + x
			di := v*dest.YStride + x
			copy(dest.Y[di:di+w], src.Y[si:si+w])
		}
	}

	// Copy Cb/Cr
	{
		y := row * 8
		x := col * 8
		w := int(n) * 8

		for v := 0; v < 8; v++ {
			si := ((y + v) * src.CStride) + x
			di := v*dest.CStride + x
			copy(dest.Cb[di:di+w], src.Cb[si:si+w])
			copy(dest.Cr[di:di+w], src.Cr[si:si+w])
		}
	}
}

func bframe_copy_macroblocks(mb_row, mb_col, mb_count int, mvd motionVectorData, fs frameStore, dest *image.YCbCr) {
	var b block
	var cb clampedblock

	for mb_addr := mb_col; mb_addr < mb_col+mb_count; mb_addr++ {
		for i := 0; i < 6; i++ {
			b.zero()
			b.motion_compensation(mvd, i, mb_row, mb_addr, fs)
			b.clamp(&cb)
			updateFrameSlice(i, mb_addr, false, dest, &cb)
		}
	}
}
