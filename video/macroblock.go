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

	cpb int
}

func (br *VideoSequence) macroblock(
	// location
	mb_address, mb_row int,
	// dct predictors
	dcp *dcDctPredictors, resetDCPredictors dcDctPredictorResetter,
	mvd *motionVectorData,
	qsc *uint32,
	frameSlice *image.YCbCr) (int, error) {

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

	// Copy skipped macroblocks for PFrames and BFrames
	if mb.macroblock_address_increment > 1 {
		switch br.PictureHeader.picture_coding_type {
		case PFrame:
			pframe_copy_macroblocks(
				mb_row, mb_address+1,
				mb.macroblock_address_increment-1,
				frameSlice, br.frameStore.past)
		case BFrame:
			bframe_copy_macroblocks(
				mb_row, mb_address+1,
				mb.macroblock_address_increment-1,
				*mvd,
				br.frameStore,
				frameSlice)
		}
	}

	// Reset dcDctPredictors: whenever a macroblock is skipped. (7.2.1)
	if mb.macroblock_address_increment > 1 {
		resetDCPredictors()
	}

	// Reset motion vector predictors: P-picture with a skipped macroblock (7.6.4.3)
	if br.PictureHeader.picture_coding_type == PFrame &&
		mb.macroblock_address_increment > 1 {
		mvd.reset()
	}

	if err := br.macroblock_mode(&mb); err != nil {
		return 0, err
	}

	// Reset dcDctPredictors: whenever a non-intra macroblock is decoded. (7.2.1)
	if mb.macroblock_type.macroblock_intra == false {
		resetDCPredictors()
	}

	// Reset motion vector predictors: intra macroblock without concealment motion vectors (7.6.4.3)
	if mb.macroblock_type.macroblock_intra == true &&
		br.PictureCodingExtension.concealment_motion_vectors == false {
		mvd.reset()
	}

	// Reset motion vector predictors: non-intra P-picture with no forward motion vectors (7.6.4.3)
	if br.PictureHeader.picture_coding_type == PFrame &&
		mb.macroblock_type.macroblock_intra == false &&
		mb.macroblock_type.macroblock_motion_forward == false {
		mvd.reset()
	}

	if mb.macroblock_type.macroblock_quant {
		if mb_qsc, err := br.Read32(5); err != nil {
			return 0, err
		} else {
			*qsc = mb_qsc
		}
	}

	if mb.macroblock_type.macroblock_motion_forward ||
		(mb.macroblock_type.macroblock_intra && br.PictureCodingExtension.concealment_motion_vectors) {
		if err := br.motion_vectors(0, &mb, mvd); err != nil {
			return 0, err
		}
	}

	if mb.macroblock_type.macroblock_motion_backward {
		if err := br.motion_vectors(1, &mb, mvd); err != nil {
			return 0, err
		}
	}

	mvd.previous.set(mb.macroblock_type, br.PictureHeader.picture_coding_type)

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
	case ChromaFormat420:
		block_count = 6
	case ChromaFormat422:
		block_count = 8
	case ChromaFormat444:
		block_count = 12
	}

	mb_address += mb.macroblock_address_increment
	pattern_code := mb.decodePatternCode(br.SequenceExtension.chroma_format)

	var b block
	var cb clampedblock

	for i := 0; i < block_count; i++ {
		cc := color_channel[i]

		if pattern_code[i] {
			if err := b.read(br, dcp, br.PictureCodingExtension.intra_vlc_format, cc, mb.macroblock_type.macroblock_intra); err != nil {
				return 0, err
			}
			b.decode_block(br, cc, *qsc, mb.macroblock_type.macroblock_intra)
			b.idct()
		} else {
			b.zero()
		}

		b.motion_compensation(*mvd, i, mb_row, mb_address, br.frameStore)
		b.clamp(&cb)
		updateFrameSlice(i, mb_address, mb.dct_type, frameSlice, &cb)
	}

	return mb_address, nil
}

func updateFrameSlice(i, mb_address int, interlaced bool, frameSlice *image.YCbCr, cb *clampedblock) {

	var (
		base_i  int
		channel []uint8
		stride  int
	)

	// channel switch
	switch i {
	case 0, 1, 2, 3:
		channel = frameSlice.Y
	case 4:
		channel = frameSlice.Cb
	case 5:
		channel = frameSlice.Cr
	}

	// base address and stride switch
	switch i {
	case 0, 1, 2, 3:
		stride = frameSlice.YStride
		base_i = mb_address * 16
	case 4, 5:
		stride = frameSlice.CStride
		base_i = mb_address * 8
	}

	// position switch
	if interlaced {
		// Field DCT coding alternates lines from each block:
		//
		//  <-8px-> <-8px->
		//  ───0───│───1───
		//  ───2───│───3───
		//  ───0───│───1───
		//  ───2───│───3───
		//  ───0───│───1───
		//  ───2───│───3───
		//  ───0───│───1───
		//  ───2───│───3───
		switch i {
		case 0, 1, 2, 3:
			base_i += (i & 1) << 3
			base_i += ((i & 2) >> 1) * stride
			stride *= 2
		}
	} else {
		// Frame DCT coding are mapped in the follow order:
		//
		//  <-8px-> <-8px->
		//  ───0───│───1───
		//  ───0───│───1───
		//  ───0───│───1───
		//  ───0───│───1───
		//  ───2───│───3───
		//  ───2───│───3───
		//  ───2───│───3───
		//  ───2───│───3───
		switch i {
		case 0, 1, 2, 3:
			base_i += (i & 1) << 3            // horiztonal positioning
			base_i += ((i & 2) << 2) * stride // vertical positioning
		}
	}

	// perform copy
	for y := 0; y < 8; y++ {
		si := y * 8
		di := base_i + (y * stride)
		copy(channel[di:di+8], cb[si:si+8])
	}

}

type PatternCode [12]bool

func (mb *Macroblock) decodePatternCode(chroma_format ChromaFormat) (pattern_code PatternCode) {
	for i := 0; i < 12; i++ {
		if mb.macroblock_type.macroblock_intra {
			pattern_code[i] = true
		} else {
			pattern_code[i] = false
		}
	}

	if mb.macroblock_type.macroblock_pattern {
		for i := 0; i < 6; i++ {
			if mask := 1 << uint(5-i); mb.cpb&mask == mask {
				pattern_code[i] = true
			}
		}

		if chroma_format == ChromaFormat422 || chroma_format == ChromaFormat444 {
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
