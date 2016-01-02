package video

func (br *VideoSequence) extension_and_user_data(i int) error {
	for {
		if nextbits, err := br.Peek32(32); err != nil {
			return err
		} else if StartCode(nextbits) != ExtensionStartCode && StartCode(nextbits) != UserDataStartCode {
			break
		} else if (i != 1) && (StartCode(nextbits) == ExtensionStartCode) {
			if err := br.extension_data(i); err != nil {
				return err
			}
		} else if StartCode(nextbits) == UserDataStartCode {
			if _, err := br.user_data(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (br *VideoSequence) extension_data(i int) error {
	for {
		if nextbits, err := br.Peek32(32); err != nil {
			return err
		} else if StartCode(nextbits) != ExtensionStartCode {
			break
		}

		br.Trash(32)

		switch i {
		case 0: /* follows sequence_extension() */
			nextbits, err := br.Peek32(4)
			if err != nil {
				return err
			}

			switch ExtensionID(nextbits) {
			case SequenceDisplayExtensionID:
				if _, err := sequence_display_extension(br); err != nil {
					return err
				}
			default:
				if _, err := sequence_scalable_extension(br); err != nil {
					return err
				}
			}

		case 1: /* NOTE - i never takes the value 1 because extension_data()
			never follows a group_of_pictures_header() */
			break
		case 2: /* follows picture_coding_extension() */
			nextbits, err := br.Peek32(4)

			if err != nil {
				return err
			}

			switch ExtensionID(nextbits) {
			case QuantMatrixExtensionID:
				quant_matrix_extension()
			case CopyrightExtensionID:
				copyright_extension()
			case PictureDisplayExtensionID:
				picture_display_extension()
			case PictureSpatialScalableExtensionID:
				picture_spatial_scalable_extension()
			default:
				if _, err := picture_temporal_scalable_extension(br); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func quant_matrix_extension() {
	panic("unsupported: quant_matrix_extension")
}

func copyright_extension() {
	panic("unsupported: copyright_extension")
}

func picture_display_extension() {
	panic("unsupported: picture_display_extension")
}

func picture_spatial_scalable_extension() {
	panic("unsupported: picture_spatial_scalable_extension")
}
