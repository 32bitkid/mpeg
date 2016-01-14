package video

// Read quatisation matricies from the bitstream
//
// The meaning of quantisationMatricies[n] is as follows:
//
//                           ┌───────────────┬───────────────┐
//                           │     4:2:0     │ 4:2:2 & 4:4:4 │
//                           ├───────┬───────┼───────┬───────┤
//                           │  lum  │  chr  │  lum  │  chr  │
//                           │(cc==0)│(cc!=0)│(cc==0)│(cc!=0)│
//  ┌────────────────────────┼───────┼───────┼───────┼───────┤
//  │            intra blocks│   0   │   0   │   0   │   2   │
//  │  (macroblock_intra = 1)│       │       │       │       │
//  ├────────────────────────┼───────┼───────┼───────┼───────┤
//  │        non-intra blocks│   1   │   1   │   1   │   3   │
//  │  (macroblock_intra = 0)│       │       │       │       │
//  └────────────────────────┴───────┴───────┴───────┴───────┘
//
func (vs *VideoSequence) quant_matrix_extension() error {

	QuantMatrixExtensionID.Assert(vs)

	if load, err := vs.ReadBit(); err != nil {
		return err
	} else if load {
		var intra_quantiser_matrix quantisationMatrix
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
		var non_intra_quantiser_matrix quantisationMatrix
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
		var chroma_intra_quantiser_matrix quantisationMatrix
		for v := 0; v < 8; v++ {
			for u := 0; u < 8; u++ {
				if val, err := vs.Read32(8); err != nil {
					return err
				} else {
					chroma_intra_quantiser_matrix[v][u] = uint8(val)
				}
			}
		}
		vs.quantisationMatricies[2] = chroma_intra_quantiser_matrix
	}

	if load, err := vs.ReadBit(); err != nil {
		return err
	} else if load {
		var chroma_non_intra_quantiser_matrix quantisationMatrix
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
