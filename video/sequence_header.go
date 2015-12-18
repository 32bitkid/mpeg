package video

import "github.com/32bitkid/mpeg/util"
import "errors"
import "io"

var ErrUnexpectedStartCode = errors.New("unexpected start code")
var ErrMissingMarkerBit = errors.New("missing marker bit")

var DefaultIntraQuantiserMatrix = [...]byte{
	8, 16, 19, 22, 26, 27, 29, 34,
	16, 16, 22, 24, 27, 29, 34, 37,
	19, 22, 26, 27, 29, 34, 34, 38,
	22, 22, 26, 27, 29, 34, 37, 40,
	22, 26, 27, 29, 32, 35, 40, 48,
	26, 27, 29, 32, 35, 40, 48, 58,
	26, 27, 29, 34, 38, 46, 56, 69,
	27, 29, 35, 38, 46, 56, 69, 83,
}

var DefaultNonIntraQuantiserMatrix = [...]byte{
	16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16,
}

type SequenceHeader struct {
	horizontal_size_value       uint32
	vertical_size_value         uint32
	aspect_ratio_information    uint32
	frame_rate_code             uint32
	bit_rate_value              uint32
	vbv_buffer_size_value       uint32
	constrained_parameters_flag bool
	intra_quantiser_matrix      [64]byte
	non_intra_quantizer_matrix  [64]byte
}

func sequence_header(br util.BitReader32) (*SequenceHeader, error) {
	var err error

	err = start_code_check(br, SequenceHeaderStartCode)
	if err != nil {
		return nil, err
	}

	sh := SequenceHeader{}

	if sh.horizontal_size_value, err = br.Read32(12); err != nil {
		return nil, err
	}

	if sh.vertical_size_value, err = br.Read32(12); err != nil {
		return nil, err
	}

	if sh.aspect_ratio_information, err = br.Read32(4); err != nil {
		return nil, err
	}

	if sh.frame_rate_code, err = br.Read32(4); err != nil {
		return nil, err
	}

	if sh.bit_rate_value, err = br.Read32(18); err != nil {
		return nil, err
	}

	err = marker_bit(br)
	if err != nil {
		return nil, err
	}

	if sh.vbv_buffer_size_value, err = br.Read32(10); err != nil {
		return nil, err
	}

	if sh.constrained_parameters_flag, err = br.ReadBit(); err != nil {
		return nil, err
	}

	load_intra_quantiser_matrix, err := br.ReadBit()
	if err != nil {
		return nil, err
	}
	if load_intra_quantiser_matrix {
		io.ReadAtLeast(br, sh.intra_quantiser_matrix[:], 64)
	} else {
		copy(sh.intra_quantiser_matrix[:], DefaultIntraQuantiserMatrix[:])
	}

	load_non_intra_quantiser_matrix, err := br.ReadBit()
	if err != nil {
		return nil, err
	}
	if load_non_intra_quantiser_matrix {
		io.ReadAtLeast(br, sh.non_intra_quantizer_matrix[:], 64)
	} else {
		copy(sh.non_intra_quantizer_matrix[:], DefaultNonIntraQuantiserMatrix[:])
	}

	next_start_code(br)

	return &sh, nil
}
