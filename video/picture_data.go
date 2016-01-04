package video

import "image"

func (self *VideoSequence) picture_data() (frame *image.YCbCr, err error) {

	w := int(self.SequenceHeader.horizontal_size_value)
	h := int(self.SequenceHeader.vertical_size_value)

	r := image.Rect(0, 0, w, h)

	var subsampleRatio image.YCbCrSubsampleRatio
	switch self.SequenceExtension.chroma_format {
	case ChromaFormat_4_2_0:
		subsampleRatio = image.YCbCrSubsampleRatio420
	case ChromaFormat_4_2_2:
		subsampleRatio = image.YCbCrSubsampleRatio422
	case ChromaFormat_4_4_4:
		subsampleRatio = image.YCbCrSubsampleRatio444
	}

	frame = image.NewYCbCr(r, subsampleRatio)

	for {
		if err := self.slice(frame); err != nil {
			return nil, err
		}

		if nextbits, err := self.Peek32(32); err != nil {
			return nil, err
		} else if StartCode(nextbits).isSlice() == false {
			break
		}
	}

	return frame, self.next_start_code()
}
