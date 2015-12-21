package video

import "github.com/32bitkid/bitreader"

type DC_DCT_Predictors [3]int

type VideoSequence struct {
	bitreader.BitReader

	// Sequence Headers
	*SequenceHeader
	*SequenceExtension

	// Extensions
	*SequenceScalableExtension

	// Picture Data
	*GroupOfPicturesHeader
	*PictureHeader
	*PictureCodingExtension

	dcDctPredictors        [3]int32
	quantisationMatricies  [4]QuantisationMatrix
	lastQuantiserScaleCode uint32
}

func NewVideoSequence(br bitreader.BitReader) VideoSequence {
	return VideoSequence{
		BitReader: br,
	}
}

func (pred *VideoSequence) resetPredictors() {
	resetValue := int32(1) << (7 + pred.PictureCodingExtension.intra_dc_precision)
	pred.dcDctPredictors[0] = resetValue
	pred.dcDctPredictors[1] = resetValue
	pred.dcDctPredictors[2] = resetValue
}

func (vs *VideoSequence) sequence_extension() (err error) {
	vs.SequenceExtension, err = sequence_extension(vs)
	return
}

func (vs *VideoSequence) group_of_pictures_header() (err error) {
	vs.GroupOfPicturesHeader, err = group_of_pictures_header(vs)
	return
}

func (vs *VideoSequence) picture_header() (err error) {
	vs.PictureHeader, err = picture_header(vs)
	return
}

func (vs *VideoSequence) picture_coding_extension() (err error) {
	vs.PictureCodingExtension, err = picture_coding_extension(vs)
	return
}

func (vs *VideoSequence) VerticalSize() uint32 {
	return vs.SequenceHeader.vertical_size_value
}
