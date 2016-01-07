package video

import "image"

func copy_macroblocks(row, col, n int, dest, src *image.YCbCr) {
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
