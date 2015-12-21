package video

import "errors"
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
}

func (br *frameProvider) sequence_header() error {
	var err error

	err = start_code_check(br, SequenceHeaderStartCode)
	if err != nil {
		return err
	}

	sh := SequenceHeader{}

	if sh.horizontal_size_value, err = br.Read32(12); err != nil {
		return err
	}

	if sh.vertical_size_value, err = br.Read32(12); err != nil {
		return err
	}

	if sh.aspect_ratio_information, err = br.Read32(4); err != nil {
		return err
	}

	if sh.frame_rate_code, err = br.Read32(4); err != nil {
		return err
	}

	if sh.bit_rate_value, err = br.Read32(18); err != nil {
		return err
	}

	err = marker_bit(br)
	if err != nil {
		return err
	}

	if sh.vbv_buffer_size_value, err = br.Read32(10); err != nil {
		return err
	}

	if sh.constrained_parameters_flag, err = br.ReadBit(); err != nil {
		return err
	}

	load_intra_quantiser_matrix, err := br.ReadBit()
	if err != nil {
		return err
	}
	if load_intra_quantiser_matrix {
		var intraQuantiserMatrix QuantisationMatrix
		for i := 0; i < 8; i++ {
			_, err := io.ReadFull(br, intraQuantiserMatrix[i][:])
			if err != nil {
				return err
			}
		}
		br.quantisationMatricies[0] = intraQuantiserMatrix
		br.quantisationMatricies[2] = intraQuantiserMatrix
	} else {
		br.quantisationMatricies[0] = DefaultQuantisationMatrices.Intra
		br.quantisationMatricies[2] = DefaultQuantisationMatrices.Intra
	}

	load_non_intra_quantiser_matrix, err := br.ReadBit()
	if err != nil {
		return err
	}
	if load_non_intra_quantiser_matrix {
		var nonIntraQuantiserMatrix QuantisationMatrix
		for i := 0; i < 8; i++ {
			_, err := io.ReadFull(br, nonIntraQuantiserMatrix[i][:])
			if err != nil {
				return err
			}
		}
		br.quantisationMatricies[1] = nonIntraQuantiserMatrix
		br.quantisationMatricies[3] = nonIntraQuantiserMatrix
	} else {
		br.quantisationMatricies[1] = DefaultQuantisationMatrices.NonIntra
		br.quantisationMatricies[3] = DefaultQuantisationMatrices.NonIntra
	}

	next_start_code(br)

	br.SequenceHeader = &sh

	return nil
}
