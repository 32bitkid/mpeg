package video

import "image"

type Macroblock struct {
	macroblock_address_increment uint32
	macroblock_type              *MacroblockType
	spatial_temporal_weight_code uint32
	frame_motion_type            uint32
	field_motion_type            uint32
	dct_type                     bool
	quantiser_scale_code         uint32
}

func (br *VideoSequence) macroblock(mbAddress int, frameSlice *image.YCbCr) (int, error) {

	mb := Macroblock{}

	nextbits, err := br.Peek32(11)
	if err != nil {
		return 0, err
	}
	if nextbits == 0x08 { // 0000 0001 000
		br.Trash(11)
		mb.macroblock_address_increment += 33
	}

	incr, err := MacroblockAddressIncrementDecoder.Decode(br)
	if err != nil {
		return 0, err
	}
	mb.macroblock_address_increment += incr

	mbAddress += int(mb.macroblock_address_increment)

	if incr > 1 {
		br.resetPredictors()
	}

	err = br.macroblock_mode(&mb)
	if err != nil {
		return 0, err
	}

	if !mb.macroblock_type.macroblock_intra {
		br.resetPredictors()
	}

	if mb.macroblock_type.macroblock_quant {
		mb.quantiser_scale_code, err = br.Read32(5)
		if err != nil {
			return 0, err
		}
		br.lastQuantiserScaleCode = mb.quantiser_scale_code
	}

	if mb.macroblock_type.macroblock_motion_forward ||
		(mb.macroblock_type.macroblock_intra && br.PictureCodingExtension.concealment_motion_vectors) {
		motion_vectors(0)
	}

	if mb.macroblock_type.macroblock_motion_backward {
		motion_vectors(1)
	}

	if mb.macroblock_type.macroblock_intra && br.PictureCodingExtension.concealment_motion_vectors {
		err := marker_bit(br)
		if err != nil {
			return 0, err
		}
	}

	if mb.macroblock_type.macroblock_pattern {
		coded_block_pattern(br, br.SequenceExtension.chroma_format)
	}

	var block_count int
	switch br.SequenceExtension.chroma_format {
	case ChromaFormat_4_2_0:
		block_count = 6
	case ChromaFormat_4_2_2:
		block_count = 8
	case ChromaFormat_4_4_4:
		block_count = 12
	}

	for i := 0; i < block_count; i++ {

		cc := calcCC(i)

		block, err := br.block(cc, &mb)
		if err != nil {
			return mbAddress, err
		}
		decoded, err := br.decode_block(cc, block, mb.macroblock_type.macroblock_intra)
		if err != nil {
			return mbAddress, err
		}
		idct(decoded)
		updateFrameSlice(i, mbAddress, frameSlice, decoded)
	}

	return mbAddress, nil
}

func calcCC(i int) int {
	switch {
	case i < 4:
		return 0
	case i&1 == 0:
		return 1
	default:
		return 2
	}
}

type clampedBlock [blockSize]uint8

func clamp(dest *clampedBlock, src *block) {
	for i := 0; i < 64; i++ {
		if src[i] > 255 {
			dest[i] = 255
		} else if src[i] < 0 {
			dest[i] = 0
		} else {
			dest[i] = uint8(src[i])
		}
	}
}

func updateFrameSlice(i int, mbAddress int, frameSlice *image.YCbCr, b *block) {

	var cb clampedBlock
	clamp(&cb, b)

	switch i {
	case 0:
		xs := (mbAddress - 1) * 16
		for y := 0; y < 8; y++ {
			si := y * 8
			di := y*frameSlice.YStride + xs
			copy(frameSlice.Y[di:di+8], cb[si:si+8])
		}
	case 1:
		xs := (mbAddress - 1) * 16
		for y := 0; y < 8; y++ {
			si := y * 8
			di := y*frameSlice.YStride + xs + 8
			copy(frameSlice.Y[di:di+8], cb[si:si+8])
		}
	case 2:
		xs := (mbAddress - 1) * 16
		for y := 0; y < 8; y++ {
			si := y * 8
			di := (y+8)*frameSlice.YStride + xs
			copy(frameSlice.Y[di:di+8], cb[si:si+8])
		}
	case 3:
		xs := (mbAddress - 1) * 16
		for y := 0; y < 8; y++ {
			si := y * 8
			di := (y+8)*frameSlice.YStride + xs + 8
			copy(frameSlice.Y[di:di+8], cb[si:si+8])
		}
	case 4:
		xs := (mbAddress - 1) * 8
		for y := 0; y < 8; y++ {
			si := y * 8
			di := y*frameSlice.CStride + xs
			copy(frameSlice.Cb[di:di+8], cb[si:si+8])
		}
	case 5:
		xs := (mbAddress - 1) * 8
		for y := 0; y < 8; y++ {
			si := y * 8
			di := y*frameSlice.CStride + xs
			copy(frameSlice.Cr[di:di+8], cb[si:si+8])
		}
	}

}

func (br *VideoSequence) macroblock_mode(mb *Macroblock) (err error) {

	var typeDecoder macroblockTypeDecoder
	switch br.PictureHeader.picture_coding_type {
	case IntraCoded:
		typeDecoder = MacroblockTypeDecoder.IFrame
	case PredictiveCoded:
		typeDecoder = MacroblockTypeDecoder.PFrame
	case BidirectionallyPredictiveCoded:
		typeDecoder = MacroblockTypeDecoder.PFrame
	default:
		panic("not implemented: macroblock type decoder")
	}

	mb.macroblock_type, err = typeDecoder(br)

	if err != nil {
		return err
	}

	if mb.macroblock_type.spatial_temporal_weight_code_flag &&
		false /* ( spatial_temporal_weight_code_table_index != ‘00’) */ {
		mb.spatial_temporal_weight_code, err = br.Read32(2)
		if err != nil {
			return err
		}
	}

	if mb.macroblock_type.macroblock_motion_forward ||
		mb.macroblock_type.macroblock_motion_backward {
		if br.PictureCodingExtension.picture_structure == PictureStructure_FramePicture {
			if br.PictureCodingExtension.frame_pred_frame_dct == 0 {
				mb.frame_motion_type, err = br.Read32(2)
				if err != nil {
					return err
				}
			}
		} else {
			mb.field_motion_type, err = br.Read32(2)
			if err != nil {
				return err
			}
		}
	}

	if br.PictureCodingExtension.picture_structure == PictureStructure_FramePicture &&
		br.PictureCodingExtension.frame_pred_frame_dct == 0 &&
		(mb.macroblock_type.macroblock_intra || mb.macroblock_type.macroblock_pattern) {
		mb.dct_type, err = br.PeekBit() //dct_type 1 uimsbf
		if err != nil {
			return err
		}
	}

	return nil
}
