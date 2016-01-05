package video

func (vs *VideoSequence) quant_matrix_extension() error {

	QuantMatrixExtensionID.assert(vs)

	if load, err := vs.ReadBit(); err != nil {
		return err
	} else if load {
		log.Println("iq")
		var intra_quantiser_matrix QuantisationMatrix
		for v := 0; v < 8; v++ {
			for u := 0; u < 8; u++ {
				if val, err := vs.Read32(8); err != nil {
					return err
				} else {
					intra_quantiser_matrix[v][u] = uint8(val)
				}
			}
		}
		vs.quantisationMatricies[0] = intra_quantiser_matrix
		vs.quantisationMatricies[2] = intra_quantiser_matrix
	}

	if load, err := vs.ReadBit(); err != nil {
		return err
	} else if load {
		log.Println("niq")
		var non_intra_quantiser_matrix QuantisationMatrix
		for v := 0; v < 8; v++ {
			for u := 0; u < 8; u++ {
				if val, err := vs.Read32(8); err != nil {
					return err
				} else {
					non_intra_quantiser_matrix[v][u] = uint8(val)
				}
			}
		}
		vs.quantisationMatricies[1] = non_intra_quantiser_matrix
		vs.quantisationMatricies[3] = non_intra_quantiser_matrix
	}

	if load, err := vs.ReadBit(); err != nil {
		return err
	} else if load {
		log.Println("ciq")
		var chroma_intra_quantiser_matrix QuantisationMatrix
		for v := 0; v < 8; v++ {
			for u := 0; u < 8; u++ {
				if val, err := vs.Read32(8); err != nil {
					return err
				} else {
					chroma_intra_quantiser_matrix[v][u] = uint8(val)
				}
			}
		}
		vs.quantisationMatricies[1] = chroma_intra_quantiser_matrix
	}

	if load, err := vs.ReadBit(); err != nil {
		return err
	} else if load {
		log.Println("cniq")
		var chroma_non_intra_quantiser_matrix QuantisationMatrix
		for v := 0; v < 8; v++ {
			for u := 0; u < 8; u++ {
				if val, err := vs.Read32(8); err != nil {
					return err
				} else {
					chroma_non_intra_quantiser_matrix[v][u] = uint8(val)
				}
			}
		}
		vs.quantisationMatricies[3] = chroma_non_intra_quantiser_matrix
	}

	return next_start_code(vs)
}
