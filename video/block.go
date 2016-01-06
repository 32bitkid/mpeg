package video

import "github.com/32bitkid/bitreader"

const (
	blockSide = 8
	blockSize = blockSide * blockSide
)

type block [blockSize]int32
type intermediaryblock [blockSide][blockSide]int32
type clampedblock [blockSize]uint8

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

func (src *block) clamp(dest *clampedblock) {
	for i := 0; i < blockSize; i++ {
		if src[i] > 255 {
			dest[i] = 255
		} else if src[i] < 0 {
			dest[i] = 0
		} else {
			dest[i] = uint8(src[i])
		}
	}
}

func (QFS *block) read(br bitreader.BitReader, dcDctPredictors *dcDctPredictors, intra_vlc_format uint32, cc int, macroblock_intra bool) error {

	eob_not_read, n := true, 0

	// 7.2.1
	if macroblock_intra {
		var dcSizeDecoder dctDCSizeDecoderFn
		if cc == 0 {
			dcSizeDecoder = dctDCSizeDecoders.Luma
		} else {
			dcSizeDecoder = dctDCSizeDecoders.Chroma
		}

		dc_dct_size, err := dcSizeDecoder(br)
		if err != nil {
			return err
		}

		dc_dct_differential, err := br.Read32(dc_dct_size)
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

		QFS[0] = dcDctPredictors[cc] + dct_diff
		dcDctPredictors[cc] = QFS[0]
		n = 1
	}

	for eob_not_read {
		var dctDecoder dctCoefficientDecoderFn

		if macroblock_intra && intra_vlc_format == 1 {
			dctDecoder = dctCoefficientDecoders.TableOne
		} else {
			dctDecoder = dctCoefficientDecoders.TableZero
		}

		run, level, end, err := dctDecoder(br, n)
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
