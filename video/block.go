package video

func (self *VideoSequence) block(cc int, mb *Macroblock) (*block, error) {

	var QFS block

	eob_not_read, n := true, 0

	// 7.2.1
	if mb.macroblock_type.macroblock_intra {
		var dcSizeDecoder DCTDCSizeDecoderFn
		if cc == 0 {
			dcSizeDecoder = DCTDCSizeDecoders.Luma
		} else {
			dcSizeDecoder = DCTDCSizeDecoders.Chroma
		}

		dc_dct_size, err := dcSizeDecoder(self)
		if err != nil {
			return nil, err
		}

		dc_dct_differential, err := self.Read32(dc_dct_size)
		if err != nil {
			return nil, err
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

		var dctDecoder DCTCoefficientDecoderFn

		if mb.macroblock_type.macroblock_intra &&
			self.PictureCodingExtension.intra_vlc_format == 1 {
			dctDecoder = DCTCoefficientDecoders.TableOne
		} else {
			dctDecoder = DCTCoefficientDecoders.TableZero
		}

		run, level, end, err := dctDecoder(self, n)
		if err != nil {
			return nil, err
		} else if end {
			eob_not_read = false
			for n < 64 {
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

	return &QFS, nil
}
