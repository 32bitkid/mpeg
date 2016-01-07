package video

import "image"

var color_channel = [12]int{0, 0, 0, 0, 1, 2, 1, 2, 1, 2, 1, 2}

type Macroblock struct {
	macroblock_address_increment int
	macroblock_type              *MacroblockType
	spatial_temporal_weight_code uint32
	frame_motion_type            uint32
	field_motion_type            uint32
	dct_type                     bool
	quantiser_scale_code         uint32

	cpb int
}

func (br *VideoSequence) macroblock(mb_row, mb_address int, frameSlice *image.YCbCr) (int, error) {

	mb := Macroblock{}

	for {
		if nextbits, err := br.Peek32(11); err != nil {
			return 0, err
		} else if nextbits == 0x08 { // 0000 0001 000
			br.Trash(11)
			mb.macroblock_address_increment += 33
		} else {
			break
		}
	}

	if incr, err := macroblockAddressIncrementDecoder.Decode(br); err != nil {
		return 0, err
	} else {
		mb.macroblock_address_increment += incr
	}

	if br.PictureHeader.picture_coding_type == PFrame && mb.macroblock_address_increment > 1 {
		copy_macroblocks(mb_row, mb_address+1, mb.macroblock_address_increment-1, frameSlice, br.frameStore.forward)
	}

	if mb.macroblock_address_increment > 1 {
		br.resetDCPredictors()
	}

	// Reset motion vector predictors: P-picture with a skipped macroblock (7.6.4.3)
	if br.PictureHeader.picture_coding_type == PFrame &&
		mb.macroblock_address_increment > 1 {
		br.pMV.reset()
	}

	if err := br.macroblock_mode(&mb); err != nil {
		return 0, err
	}

	if mb.macroblock_type.macroblock_intra == false {
		br.resetDCPredictors()
	}

	// Reset motion vector predictors: intra macroblock without concealment motion vectors (7.6.4.3)
	if mb.macroblock_type.macroblock_intra == true &&
		br.PictureCodingExtension.concealment_motion_vectors == false {
		br.pMV.reset()
	}

	// Reset motion vector predictors: non-intra P-picture with no forward motion vectors (7.6.4.3)
	if br.PictureHeader.picture_coding_type == PFrame &&
		mb.macroblock_type.macroblock_intra == false &&
		mb.macroblock_type.macroblock_motion_forward == false {
		br.pMV.reset()
	}

	if mb.macroblock_type.macroblock_quant {
		if qsc, err := br.Read32(5); err != nil {
			return 0, err
		} else {
			mb.quantiser_scale_code = qsc
			br.currentQSC = qsc
		}
	}

	var mvd motionVectorData

	if mb.macroblock_type.macroblock_motion_forward ||
		(mb.macroblock_type.macroblock_intra && br.PictureCodingExtension.concealment_motion_vectors) {
		if err := br.motion_vectors(0, &mb, &mvd); err != nil {
			return 0, err
		}
	}

	if mb.macroblock_type.macroblock_motion_backward {
		if err := br.motion_vectors(1, &mb, &mvd); err != nil {
			return 0, err
		}
	}

	if mb.macroblock_type.macroblock_intra && br.PictureCodingExtension.concealment_motion_vectors {
		if err := marker_bit(br); err != nil {
			return 0, err
		}
	}

	if mb.macroblock_type.macroblock_pattern {
		if cpb, err := coded_block_pattern(br, br.SequenceExtension.chroma_format); err != nil {
			return 0, nil
		} else {
			mb.cpb = cpb
		}
	}

	var block_count int
	switch br.SequenceExtension.chroma_format {
	case ChromaFormat_420:
		block_count = 6
	case ChromaFormat_422:
		block_count = 8
	case ChromaFormat_444:
		block_count = 12
	}

	mb_address += mb.macroblock_address_increment
	pattern_code := mb.decodePatternCode(br.SequenceExtension.chroma_format)

	var b block

	for i := 0; i < block_count; i++ {
		cc := color_channel[i]

		if pattern_code[i] {
			if err := b.read(br, &br.dcDctPredictors, br.PictureCodingExtension.intra_vlc_format, cc, mb.macroblock_type.macroblock_intra); err != nil {
				return 0, err
			}
			b.decode_block(br, cc, mb.macroblock_type.macroblock_intra)
			b.idct()
		} else {
			b.empty()
		}

		br.motion_compensation(&mvd, i, mb_row, mb_address, &mb, &b)
		updateFrameSlice(i, mb_address, mb.dct_type, frameSlice, &b)
	}

	return mb_address, nil
}

func updateFrameSlice(i, mb_address int, interlaced bool, frameSlice *image.YCbCr, b *block) {

	var cb clampedblock
	b.clamp(&cb)

	if interlaced {
		switch i {
		case 0:
			xs := mb_address * 16
			for y := 0; y < 8; y++ {
				si := y * 8
				di := (y*2)*frameSlice.YStride + xs
				copy(frameSlice.Y[di:di+8], cb[si:si+8])
			}
		case 1:
			xs := mb_address * 16
			for y := 0; y < 8; y++ {
				si := y * 8
				di := (y*2)*frameSlice.YStride + (xs + 8)
				copy(frameSlice.Y[di:di+8], cb[si:si+8])
			}
		case 2:
			xs := mb_address * 16
			for y := 0; y < 8; y++ {
				si := y * 8
				di := ((y*2)+1)*frameSlice.YStride + xs
				copy(frameSlice.Y[di:di+8], cb[si:si+8])
			}
		case 3:
			xs := mb_address * 16
			for y := 0; y < 8; y++ {
				si := y * 8
				di := ((y*2)+1)*frameSlice.YStride + (xs + 8)
				copy(frameSlice.Y[di:di+8], cb[si:si+8])
			}
		case 4:
			xs := mb_address * 8
			for y := 0; y < 8; y++ {
				si := y * 8
				di := y*frameSlice.CStride + xs
				copy(frameSlice.Cb[di:di+8], cb[si:si+8])
			}
		case 5:
			xs := mb_address * 8
			for y := 0; y < 8; y++ {
				si := y * 8
				di := y*frameSlice.CStride + xs
				copy(frameSlice.Cr[di:di+8], cb[si:si+8])
			}
		}
	} else {
		switch i {
		case 0:
			xs := mb_address * 16
			for y := 0; y < 8; y++ {
				si := y * 8
				di := y*frameSlice.YStride + xs
				copy(frameSlice.Y[di:di+8], cb[si:si+8])
			}
		case 1:
			xs := mb_address * 16
			for y := 0; y < 8; y++ {
				si := y * 8
				di := y*frameSlice.YStride + (xs + 8)
				copy(frameSlice.Y[di:di+8], cb[si:si+8])
			}
		case 2:
			xs := mb_address * 16
			for y := 0; y < 8; y++ {
				si := y * 8
				di := (y+8)*frameSlice.YStride + xs
				copy(frameSlice.Y[di:di+8], cb[si:si+8])
			}
		case 3:
			xs := mb_address * 16
			for y := 0; y < 8; y++ {
				si := y * 8
				di := (y+8)*frameSlice.YStride + (xs + 8)
				copy(frameSlice.Y[di:di+8], cb[si:si+8])
			}
		case 4:
			xs := mb_address * 8
			for y := 0; y < 8; y++ {
				si := y * 8
				di := y*frameSlice.CStride + xs
				copy(frameSlice.Cb[di:di+8], cb[si:si+8])
			}
		case 5:
			xs := mb_address * 8
			for y := 0; y < 8; y++ {
				si := y * 8
				di := y*frameSlice.CStride + xs
				copy(frameSlice.Cr[di:di+8], cb[si:si+8])
			}
		}
	}

}

type PatternCode [12]bool

func (mb *Macroblock) decodePatternCode(chroma_format chromaFormat) (pattern_code PatternCode) {
	for i := 0; i < 12; i++ {
		if mb.macroblock_type.macroblock_intra {
			pattern_code[i] = true
		} else {
			pattern_code[i] = false
		}
	}

	if mb.macroblock_type.macroblock_pattern {
		for i := 0; i < 6; i++ {
			mask := 1 << uint(5-i)
			if mb.cpb&mask == mask {
				pattern_code[i] = true
			}
		}

		if chroma_format == ChromaFormat_422 || chroma_format == ChromaFormat_444 {
			panic("unsupported: coded block pattern chroma format")
		}
	}

	return
}

func (br *VideoSequence) macroblock_mode(mb *Macroblock) (err error) {

	var typeDecoder macroblockTypeDecoderFn
	switch br.PictureHeader.picture_coding_type {
	case IntraCoded:
		typeDecoder = macroblockTypeDecoder.IFrame
	case PredictiveCoded:
		typeDecoder = macroblockTypeDecoder.PFrame
	case BidirectionallyPredictiveCoded:
		typeDecoder = macroblockTypeDecoder.BFrame
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
		mb.dct_type, err = br.ReadBit() //dct_type 1 uimsbf
		if err != nil {
			return err
		}
	}

	return nil
}
