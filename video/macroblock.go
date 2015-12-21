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
		coded_block_pattern()
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
		decoded, err := br.block(i, &mb)
		if err != nil {
			return 0, err
		}
		updateFrameSlice(i, mbAddress, frameSlice, decoded)

	}

	return mbAddress, nil
}

func updateFrameSlice(i int, mbAddress int, frameSlice *image.YCbCr, b *block) {
	switch i {
	case 0:
		xs := (mbAddress - 1) * 16
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				frameSlice.Y[y*frameSlice.YStride+xs+x] = uint8(b[y*8+x])
			}
		}
	case 1:
		xs := (mbAddress - 1) * 16
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				frameSlice.Y[y*frameSlice.YStride+xs+x+8] = uint8(b[y*8+x])
			}
		}
	case 2:
		xs := (mbAddress - 1) * 16
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				frameSlice.Y[(y+8)*frameSlice.YStride+xs+x] = uint8(b[y*8+x])
			}
		}
	case 3:
		xs := (mbAddress - 1) * 16
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				frameSlice.Y[(y+8)*frameSlice.YStride+xs+x+8] = uint8(b[y*8+x])
			}
		}
	case 4:
		xs := (mbAddress - 1) * 8
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				frameSlice.Cb[y*frameSlice.CStride+xs+x] = uint8(b[y*8+x])
			}
		}
	case 5:
		xs := (mbAddress - 1) * 8
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				frameSlice.Cr[y*frameSlice.CStride+xs+x] = uint8(b[y*8+x])
			}
		}
	}

}

func motion_vectors(i int) {
	switch i {
	case 0:
		panic("forward motion vectors")
	case 1:
		panic("backwards motion vectors")
	default:
		panic("unknown motion vectors")
	}

}

func coded_block_pattern() {
	panic("coding block pattern")
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
