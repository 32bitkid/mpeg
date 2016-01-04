package video

var scan = [2][8][8]int{
	{
		{0, 1, 5, 6, 14, 15, 27, 28},
		{2, 4, 7, 13, 16, 26, 29, 42},
		{3, 8, 12, 17, 25, 30, 41, 43},
		{9, 11, 18, 24, 31, 40, 44, 53},
		{10, 19, 23, 32, 39, 45, 52, 54},
		{20, 22, 33, 38, 46, 51, 55, 60},
		{21, 34, 37, 47, 50, 56, 59, 61},
		{35, 36, 48, 49, 57, 58, 62, 63},
	},
	{
		{0, 4, 6, 20, 22, 36, 38, 52},
		{1, 5, 7, 21, 23, 37, 39, 53},
		{2, 8, 19, 24, 34, 40, 50, 54},
		{3, 9, 18, 25, 35, 41, 51, 55},
		{10, 17, 26, 30, 42, 46, 56, 60},
		{11, 16, 27, 31, 43, 47, 57, 61},
		{12, 15, 28, 32, 44, 48, 58, 62},
		{13, 14, 29, 33, 45, 49, 59, 63},
	},
}

func sign(i int32) int32 {
	if i > 0 {
		return 1
	} else if i < 0 {
		return -1
	}
	return 0
}

type intermediaryblock [8][8]int32

func (self *VideoSequence) decode_block(cc int, QFS *block, F *block, macroblock_intra bool) error {
	var QF intermediaryblock
	var Fpp intermediaryblock
	var Fp intermediaryblock

	// inverse scan
	{
		alternate_scan := self.PictureCodingExtension.alternate_scan
		for v := 0; v < 8; v++ {
			for u := 0; u < 8; u++ {
				QF[v][u] = QFS[scan[alternate_scan][v][u]]
			}
		}
	}

	// Inverse quantisation
	{
		q_scale_type := self.PictureCodingExtension.q_scale_type
		quantiser_scale_code := self.currentQSC
		quantiser_scale := quantiser_scale_tables[q_scale_type][quantiser_scale_code]

		var w int
		if cc == 0 {
			if macroblock_intra {
				w = 0
			} else {
				w = 1
			}
		} else {
			if self.SequenceExtension.chroma_format == ChromaFormat_4_2_0 {
				if macroblock_intra {
					w = 0
				} else {
					w = 1
				}
			} else {
				if macroblock_intra {
					w = 2
				} else {
					w = 3
				}
			}
		}

		W := self.quantisationMatricies

		for v := 0; v < 8; v++ {
			for u := 0; u < 8; u++ {
				if (u == 0) && (v == 0) && (macroblock_intra) {
					// Table 7-4
					intra_dc_mult := int32(1) << (3 - self.PictureCodingExtension.intra_dc_precision)
					Fpp[v][u] = intra_dc_mult * QF[v][u]
				} else {
					if macroblock_intra {
						Fpp[v][u] = (QF[v][u] * int32(W[w][v][u]) * quantiser_scale * 2) / 32
					} else {
						Fpp[v][u] = (((QF[v][u] * 2) + sign(QF[v][u])) * int32(W[w][v][u]) *
							quantiser_scale) / 32
					}
				}
			}
		}
	}

	{
		// Saturation
		var sum int32 = 0
		for v := 0; v < 8; v++ {
			for u := 0; u < 8; u++ {
				if Fpp[v][u] > 2047 {
					Fp[v][u] = 2047
				} else if Fpp[v][u] < -2048 {
					Fp[v][u] = -2048
				} else {
					Fp[v][u] = Fpp[v][u]
				}
				sum = sum + Fp[v][u]
				F[v*8+u] = Fp[v][u]
			}
		}

		// Mismatch control
		if (sum & 1) == 0 {
			if (F[7*8+7] & 1) != 0 {
				F[7*8+7] = Fp[7][7] - 1
			} else {
				F[7*8+7] = Fp[7][7] + 1
			}
		}
	}

	return nil
}
