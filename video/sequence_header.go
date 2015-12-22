package video

import "errors"
import "github.com/32bitkid/bitreader"
import "io"

var ErrUnexpectedStartCode = errors.New("unexpected start code")
var ErrMissingMarkerBit = errors.New("missing marker bit")

type SequenceHeader struct {
	horizontal_size_value       uint32
	vertical_size_value         uint32
	aspect_ratio_information    uint32
	frame_rate_code             uint32
	bit_rate_value              uint32
	vbv_buffer_size_value       uint32
	constrained_parameters_flag bool

	load_intra_quantiser_matrix     bool
	load_non_intra_quantizer_matrix bool

	intra_quantiser_matrix     QuantisationMatrix
	non_intra_quantizer_matrix QuantisationMatrix
}

func sequence_header(br bitreader.BitReader) (*SequenceHeader, error) {

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

	sh.load_intra_quantiser_matrix, err = br.ReadBit()
	if err != nil {
		return nil, err
	}
	if sh.load_intra_quantiser_matrix {
		for i := 0; i < 8; i++ {
			_, err := io.ReadFull(br, sh.intra_quantiser_matrix[i][:])
			if err != nil {
				return nil, err
			}
		}
	}

	sh.load_non_intra_quantizer_matrix, err = br.ReadBit()
	if err != nil {
		return nil, err
	}
	if sh.load_non_intra_quantizer_matrix {
		for i := 0; i < 8; i++ {
			_, err := io.ReadFull(br, sh.non_intra_quantizer_matrix[i][:])
			if err != nil {
				return nil, err
			}
		}
	}

	return &sh, next_start_code(br)
}

func (fp *frameProvider) sequence_header() (err error) {

	sh, err := sequence_header(fp)
	if err != nil {
		return err
	}

	if sh.load_intra_quantiser_matrix {
		fp.quantisationMatricies[0] = sh.intra_quantiser_matrix
		fp.quantisationMatricies[2] = sh.intra_quantiser_matrix
	} else {
		fp.quantisationMatricies[0] = DefaultQuantisationMatrices.Intra
		fp.quantisationMatricies[2] = DefaultQuantisationMatrices.Intra
	}

	if sh.load_non_intra_quantizer_matrix {
		fp.quantisationMatricies[1] = sh.non_intra_quantizer_matrix
		fp.quantisationMatricies[3] = sh.non_intra_quantizer_matrix
	} else {
		fp.quantisationMatricies[1] = DefaultQuantisationMatrices.NonIntra
		fp.quantisationMatricies[3] = DefaultQuantisationMatrices.NonIntra
	}

	fp.SequenceHeader = sh

	return nil
}
