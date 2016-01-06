package video

const blockSize = 64 // A DCT block is 8x8.
type block [blockSize]int32

func (dest *block) sum(other *block) {
	for i := 0; i < blockSize; i++ {
		dest[i] += other[i]
	}
}

func (b *block) empty() {
	for i := 0; i < blockSize; i++ {
		b[i] = 0
	}
}

func (self *VideoSequence) block(cc int, mb *Macroblock, QFS *block) error {

	eob_not_read, n := true, 0

	// 7.2.1
	if mb.macroblock_type.macroblock_intra {
		var dcSizeDecoder dctDCSizeDecoderFn
		if cc == 0 {
			dcSizeDecoder = dctDCSizeDecoders.Luma
		} else {
			dcSizeDecoder = dctDCSizeDecoders.Chroma
		}

		dc_dct_size, err := dcSizeDecoder(self)
		if err != nil {
			return err
		}

		dc_dct_differential, err := self.Read32(dc_dct_size)
		if err != nil {
			return err
		}

		var dct_diff int32
		if dc_dct_size == 0 {
			dct_diff = 0
		} else {
			half_range := uint32(1) << (dc_dct_size - 1)
			if dc_dct_differential >= half_range {
				dct_diff = int32(dc_dct_differential)
			} else {
				dct_diff = int32(dc_dct_differential+1) - int32(2*half_range)
			}
		}

		QFS[0] = self.dcDctPredictors[cc] + dct_diff
		self.dcDctPredictors[cc] = QFS[0]
		n = 1

		if QFS[0] < 0 || QFS[0] > (1<<(8+self.PictureCodingExtension.intra_dc_precision)-1) {
			panic("DC is out of range")
		}
	}

	for eob_not_read {
		var dctDecoder dctCoefficientDecoderFn

		if mb.macroblock_type.macroblock_intra &&
			self.PictureCodingExtension.intra_vlc_format == 1 {
			dctDecoder = dctCoefficientDecoders.TableOne
		} else {
			dctDecoder = dctCoefficientDecoders.TableZero
		}

		run, level, end, err := dctDecoder(self, n)
		if err != nil {
			return err
		} else if end {
			eob_not_read = false
			for n < blockSize {
				QFS[n] = 0
				n = n + 1
			}
		} else {
			for m := 0; m < run; m++ {
				QFS[n] = 0
				n = n + 1
			}
			QFS[n] = level
			n = n + 1
		}
	}

	return nil
}
