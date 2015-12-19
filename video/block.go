package video

func (self *VideoSequence) block(i int, mb *Macroblock) (interface{}, error) {

	var QFS [64]int

	var cc int
	if i < 4 {
		cc = 0
	} else if i&1 == 0 {
		cc = 1
	} else {
		cc = 2
	}

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

		var dct_diff int
		if dc_dct_size == 0 {
			dct_diff = 0
		} else {
			half_range := uint32(1) << (dc_dct_size - 1)
			if dc_dct_differential >= half_range {
				dct_diff = int(dc_dct_differential)
			} else {
				dct_diff = int(dc_dct_differential+1) - int(2*half_range)
			}
		}

		QFS[0] = /* dc_dct_pred[cc] + */ dct_diff
		//dc_dct_pred[cc] = QFS[0]
		n = 1
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
		}

		if end {
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

	log.Println(QFS)

	return nil, nil
	/*
		eob_not_read = 1;
		while ( eob_not_read )
		{
		<decode VLC, decode Escape coded coefficient if required>
		if ( <decoded VLC indicates End of block> ) {
		eob_not_read = 0;

		} else {
		for ( m = 0; m < run; m++ ) {
		QFS[n] = 0;
		n = n + 1;
		}
		QFS[n] = signed_level
		n = n + 1;
		}
		}
	*/
}
