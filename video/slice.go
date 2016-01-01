package video

import "errors"
import "image"

var ErrInvalidReservedBits = errors.New("invalid reserved bits")

type Slice struct {
	slice_start_code                  StartCode // 32 bslbf
	slice_vertical_position_extension uint32    // 3 uimsbf
	priority_breakpoint               uint32    // 7 uimsbf
	quantiser_scale_code              uint32    // 5 uimsbf
	intra_slice_flag                  bool      // 1 bslbf
	intra_slice                       bool      // 1 uimsbf
	extra_information                 []byte
}

func (br *VideoSequence) slice(frame *image.YCbCr) error {

	br.resetPredictors()

	s := Slice{}

	code, err := br.Read32(32)
	if err != nil {
		return err
	}

	if StartCode(code).isSlice() == false {
		return ErrUnexpectedStartCode
	}

	s.slice_start_code = StartCode(code)

	yPos := int((code & 0xFF) - 1)
	frameSlice := &image.YCbCr{
		Y:       frame.Y[16*yPos*frame.YStride:],
		Cb:      frame.Cb[8*yPos*frame.CStride:],
		Cr:      frame.Cr[8*yPos*frame.CStride:],
		YStride: frame.YStride,
		CStride: frame.CStride,
		Rect:    image.Rect(0, 0, 16, frame.YStride),
	}

	if br.VerticalSize() > 2800 {
		s.slice_vertical_position_extension, err = br.Read32(3)
		if err != nil {
			return err
		}
	}

	// TODO
	if false /* (<sequence_scalable_extension() is present in the bitstream>) */ {
		if false /* (scalable_mode == “data partitioning” ) */ {
			s.priority_breakpoint, err = br.Read32(7)
			if err != nil {
				return err
			}
		}
	}

	s.quantiser_scale_code, err = br.Read32(5)
	if err != nil {
		return err
	}

	br.lastQuantiserScaleCode = s.quantiser_scale_code

	nextbits, err := br.Peek32(1)
	if err != nil {
		return err
	}

	if nextbits == 1 {
		s.intra_slice_flag, err = br.ReadBit()
		if err != nil {
			return err
		}

		s.intra_slice, err = br.ReadBit()
		if err != nil {
			return err
		}

		reserved_bits, err := br.Read32(7)
		if err != nil {
			return err
		} else if reserved_bits != 0 {
			return ErrInvalidReservedBits
		}

		for {
			nextbits, err = br.Peek32(1)
			if err != nil {
				return err
			} else if nextbits != 1 {
				break
			}

			err = br.Trash(1)
			if err != nil {
				return err
			}

			data, err := br.Read32(8)
			if err != nil {
				return err
			}
			s.extra_information = append(s.extra_information, byte(data))
		}
	}

	err = br.Trash(1)
	if err != nil {
		return err
	}

	var mbAddress int = 0
	for {
		mbAddress, err = br.macroblock(mbAddress, frameSlice)
		if err != nil {
			return err
		}

		nextbits, err = br.Peek32(23)
		if err != nil {
			return err
		}

		if nextbits == 0 {
			break
		}

	}

	return next_start_code(br)
}
