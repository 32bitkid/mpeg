package video

import "errors"
import "github.com/32bitkid/bitreader"

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
	load_non_intra_quantiser_matrix bool

	intra_quantiser_matrix     quantisationMatrix
	non_intra_quantiser_matrix quantisationMatrix
}

// ReadSequenceHeader reads a sequence header from the bit stream.
func ReadSequenceHeader(br bitreader.BitReader) (*SequenceHeader, error) {

	var err error

	err = SequenceHeaderStartCode.Assert(br)
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

	if err := marker_bit(br); err != nil {
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
		for v := 0; v < 8; v++ {
			for u := 0; u < 8; u++ {
				if val, err := br.Read32(8); err != nil {
					return nil, err
				} else {
					sh.intra_quantiser_matrix[v][u] = uint8(val)
				}
			}
		}
	}

	sh.load_non_intra_quantiser_matrix, err = br.ReadBit()
	if err != nil {
		return nil, err
	}
	if sh.load_non_intra_quantiser_matrix {
		for v := 0; v < 8; v++ {
			for u := 0; u < 8; u++ {
				if val, err := br.Read32(8); err != nil {
					return nil, err
				} else {
					sh.non_intra_quantiser_matrix[v][u] = uint8(val)
				}
			}
		}
	}

	return &sh, next_start_code(br)
}

func (vs *VideoSequence) sequence_header() (err error) {

	sh, err := ReadSequenceHeader(vs)
	if err != nil {
		return err
	}

	if sh.load_intra_quantiser_matrix {
		vs.quantisationMatricies[0] = sh.intra_quantiser_matrix
		vs.quantisationMatricies[2] = sh.intra_quantiser_matrix
	} else {
		vs.quantisationMatricies[0] = defaultQuantisationMatrices.Intra
		vs.quantisationMatricies[2] = defaultQuantisationMatrices.Intra
	}

	if sh.load_non_intra_quantiser_matrix {
		vs.quantisationMatricies[1] = sh.non_intra_quantiser_matrix
		vs.quantisationMatricies[3] = sh.non_intra_quantiser_matrix
	} else {
		vs.quantisationMatricies[1] = defaultQuantisationMatrices.NonIntra
		vs.quantisationMatricies[3] = defaultQuantisationMatrices.NonIntra
	}

	vs.SequenceHeader = sh

	return nil
}
