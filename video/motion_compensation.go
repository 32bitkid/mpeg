package video

func (vs *VideoSequence) motion_compensation(motionVectors *motionVectorData, i, mb_row, mb_addr int, mb *Macroblock, b *block) {

	if vs.PictureHeader.picture_coding_type == PFrame {

		if !mb.macroblock_type.macroblock_intra {

			horizontal, vertical := motionVectors.actual[0][0][0], motionVectors.actual[0][0][1]

			// Scale Cb/Cr vectors
			if i >= 4 {
				horizontal >>= 1
				vertical >>= 1
			}

			// is half?
			h_half, v_half := (horizontal&1) == 1, (vertical&1) == 1

			var (
				srcX    int
				srcY    int
				channel []uint8
				stride  int
			)

			image := vs.frameStore.past

			switch i {
			case 0, 1, 2, 3:
				channel = image.Y
				stride = image.YStride
				srcX = (mb_addr * 16) + (i&1)<<3
				srcY = (mb_row * 16) + (i&2)<<2
			case 4:
				channel = image.Cb
				stride = image.CStride
				srcX = mb_addr * 8
				srcY = mb_row * 8
			case 5:
				channel = image.Cr
				stride = image.CStride
				srcX = mb_addr * 8
				srcY = mb_row * 8
			}

			srcX += horizontal >> 1
			srcY += vertical >> 1

			for v := 0; v < 8; v++ {
				for u := 0; u < 8; u++ {
					i := ((srcY + v) * stride) + (srcX + u)
					switch {
					case !h_half && !v_half:
						b[v*8+u] += int32(channel[i])
					case h_half && !v_half:
						b[v*8+u] += (int32(channel[i]) + int32(channel[i+1])) / 2
					case !h_half && v_half:
						b[v*8+u] += (int32(channel[i]) + int32(channel[i+stride])) / 2
					case h_half && v_half:
						b[v*8+u] += (int32(channel[i]) + int32(channel[i+1]) + int32(channel[i+stride]) + int32(channel[i+stride+1])) / 4
					}
				}
			}
		}
	}
}
