package video

import "github.com/32bitkid/bitreader"

func extension_and_user_data(i int, br bitreader.BitReader) error {

	for {
		nextbits, err := br.Peek32(32)
		if err != nil {
			return err
		}

		if nextbits != ExtensionStartCode && nextbits != UserDataStartCode {
			break
		}

		if (i != 1) && (nextbits == ExtensionStartCode) {
			extension_data(i, br)
		}

		if nextbits == UserDataStartCode {
			user_data(br)
		}

	}

	return nil
}

func extension_data(i int, br bitreader.BitReader) error {
	for {
		nextbits, err := br.Peek32(32)
		if err != nil {
			return err
		} else if nextbits != ExtensionStartCode {
			break
		}

		br.Trash(32)

		switch i {
		case 0:
			/* follows sequence_extension() */

			nextbits, err = br.Peek32(4)
			if err != nil {
				return err
			}

			switch ExtensionID(nextbits) {
			case SequenceDisplayExtensionID:
				sequence_display_extension(br)
			default:
				sequence_scalable_extension(br)
			}

		case 1:
			/* NOTE - i never takes the value 1 because extension_data()
			never follows a group_of_pictures_header() */
		case 2:
			/* follows picture_coding_extension() */

			nextbits, err = br.Peek32(4)

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
				picture_temporal_scalable_extension(br)
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
