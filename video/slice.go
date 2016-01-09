package video

import "errors"
import "image"

var ErrInvalidReservedBits = errors.New("invalid reserved bits")

type Slice struct {
	slice_start_code                  StartCode // 32 bslbf
	slice_vertical_position_extension uint32    // 3 uimsbf
	priority_breakpoint               uint32    // 7 uimsbf
	intra_slice_flag                  bool      // 1 bslbf
	intra_slice                       bool      // 1 uimsbf
	extra_information                 []uint8
}

func (br *VideoSequence) slice(frame *image.YCbCr) error {

	code, err := br.Read32(32)
	if err != nil {
		return err
	}

	if StartCode(code).isSlice() == false {
		return ErrUnexpectedStartCode
	}

	var dcp dcDctPredictors
	resetDcPredictors := dcp.createResetter(br.PictureCodingExtension.intra_dc_precision)
	// Reset dcDctPredictors: at start of slice (7.2.1)
	resetDcPredictors()

	// Reset motion vector predictors: Start if each slice (7.6.4.3)
	var mvd motionVectorData
	mvd.reset()

	s := Slice{}
	s.slice_start_code = StartCode(code)

	mb_row := int((code & 0xFF) - 1)

	frameSlice := &image.YCbCr{
		Y:       frame.Y[16*mb_row*frame.YStride:],
		Cb:      frame.Cb[8*mb_row*frame.CStride:],
		Cr:      frame.Cr[8*mb_row*frame.CStride:],
		YStride: frame.YStride,
		CStride: frame.CStride,
		Rect:    image.Rect(frame.Rect.Min.X, mb_row*16, frame.Rect.Max.X, mb_row*16+16),
	}

	if br.SequenceHeader.vertical_size_value > 2800 {
		s.slice_vertical_position_extension, err = br.Read32(3)
		if err != nil {
			return err
		}
	}

	// TODO(jh): sequence_scalable support
	if false /* (<sequence_scalable_extension() is present in the bitstream>) */ {
		if false /* (scalable_mode == “data partitioning” ) */ {
			s.priority_breakpoint, err = br.Read32(7)
			if err != nil {
				return err
			}
		}
	}

	var quantiser_scale_code uint32
	if qsc, err := br.Read32(5); err != nil {
		return err
	} else {
		quantiser_scale_code = qsc
	}

	if nextbits, err := br.Peek32(1); err != nil {
		return err
	} else if nextbits == 1 {
		s.intra_slice_flag, err = br.ReadBit()
		if err != nil {
			return err
		}

		s.intra_slice, err = br.ReadBit()
		if err != nil {
			return err
		}

		if reserved_bits, err := br.Read32(7); err != nil {
			return err
		} else if reserved_bits != 0 {
			return ErrInvalidReservedBits
		}

		for {
			if nextbits, err := br.Peek32(1); err != nil {
				return err
			} else if nextbits != 1 {
				break
			} else if err := br.Trash(1); err != nil {
				return err
			}

			if data, err := br.Read32(8); err != nil {
				return err
			} else {
				s.extra_information = append(s.extra_information, uint8(data))
			}
		}
	}

	if err := br.Trash(1); err != nil {
		return err
	}

	var mb_address int = -1
	for {
		mb_address, err = br.macroblock(
			mb_address, mb_row,
			&dcp, resetDcPredictors,
			&mvd,
			&quantiser_scale_code,
			frameSlice)

		if err != nil {
			return err
		}

		if nextbits, err := br.Peek32(23); err != nil {
			return err
		} else if nextbits == 0 {
			break
		}
	}

	return next_start_code(br)
}
