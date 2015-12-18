package video

import "github.com/32bitkid/mpeg/util"
import "errors"

var ErrInvalidReservedBits = errors.New("invalid reserved bits")

type Slice struct {
	slice_start_code                  StartCode // 32 bslbf
	slice_vertical_position_extension uint32    // 3 uimsbf
	priority_breakpoint               uint32    // 7 uimsbf
	quantiser_scale_code              uint32    // 5 uimsbf
	intra_slice_flag                  bool      // 1 bslbf
	intra_slice                       bool      // 1 uimsbf
	extra_information                 []byte
	macroblocks                       []interface{}
}

func slice(br util.BitReader32) (*Slice, error) {

	s := Slice{}

	code, err := br.Read32(32)
	if err != nil {
		return nil, err
	}

	if is_slice_start_code(StartCode(code)) == false {
		return nil, ErrUnexpectedStartCode
	}

	s.slice_start_code = StartCode(code)

	if false /* (vertical_size > 2800) */ {
		s.slice_vertical_position_extension, err = br.Read32(3)
		if err != nil {
			return nil, err
		}
	}

	if false /* (<sequence_scalable_extension() is present in the bitstream>) */ {
		if false /* (scalable_mode == “data partitioning” ) */ {
			s.priority_breakpoint, err = br.Read32(7)
			if err != nil {
				return nil, err
			}
		}
	}

	s.quantiser_scale_code, err = br.Read32(5)
	if err != nil {
		return nil, err
	}

	nextbits, err := br.Peek32(1)
	if err != nil {
		return nil, err
	}

	if nextbits == 1 {
		s.intra_slice_flag, err = br.ReadBit()
		if err != nil {
			return nil, err
		}

		s.intra_slice, err = br.ReadBit()
		if err != nil {
			return nil, err
		}

		reserved_bits, err := br.Read32(7)
		if err != nil {
			return nil, err
		} else if reserved_bits != 0 {
			return nil, ErrInvalidReservedBits
		}

		for {
			nextbits, err = br.Peek32(1)
			if err != nil {
				return nil, err
			} else if nextbits != 1 {
				break
			}

			err = br.Trash(1)
			if err != nil {
				return nil, err
			}

			data, err := br.Read32(8)
			if err != nil {
				return nil, err
			}
			s.extra_information = append(s.extra_information, byte(data))
		}
	}

	err = br.Trash(1)
	if err != nil {
		return nil, err
	}

	for {
		mb, err := macroblock(br)
		if err != nil {
			return nil, err
		}
		s.macroblocks = append(s.macroblocks, mb)

		nextbits, err = br.Peek32(23)
		if err != nil {
			return nil, err
		}

		if nextbits == 0 {
			break
		}

	}

	return &s, next_start_code(br)
}
