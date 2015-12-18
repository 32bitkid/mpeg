package video

import "github.com/32bitkid/mpeg/util"

type VideoSequence struct {
	util.BitReader32

	// Sequence Headers
	*SequenceHeader
	*SequenceExtension

	// Extensions
	*SequenceScalableExtension

	// Picture Data
	*GroupOfPicturesHeader
	*PictureHeader
	*PictureCodingExtension
	*PictureData
}

func NewVideoSequence(br util.BitReader32) VideoSequence {
	return VideoSequence{
		BitReader32: br,
	}
}

func (vs *VideoSequence) sequence_header() (err error) {
	vs.SequenceHeader, err = sequence_header(vs)
	return
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
