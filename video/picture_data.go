package video

import "image"

func createFrameBuffer(w, h uint32, cf ChromaFormat) *image.YCbCr {

	horizontalMacroblocks := w >> 4
	verticalMacroblocks := h >> 4

	if (w & 0xf) != 0 {
		horizontalMacroblocks++
	}

	if (h & 0xf) != 0 {
		verticalMacroblocks++
	}

	r := image.Rect(0, 0, int(horizontalMacroblocks<<4), int(verticalMacroblocks<<4))

	var subsampleRatio image.YCbCrSubsampleRatio
	switch cf {
	case ChromaFormat420:
		subsampleRatio = image.YCbCrSubsampleRatio420
	case ChromaFormat422:
		subsampleRatio = image.YCbCrSubsampleRatio422
	case ChromaFormat444:
		subsampleRatio = image.YCbCrSubsampleRatio444
	}

	return image.NewYCbCr(r, subsampleRatio)
}

func (self *VideoSequence) picture_data() (frame *image.YCbCr, err error) {

	frame = createFrameBuffer(self.SequenceHeader.horizontal_size_value, self.SequenceHeader.vertical_size_value, self.SequenceExtension.chroma_format)

	for {
		if err := self.slice(frame); err != nil {
			return nil, err
		}

		if nextbits, err := self.Peek32(32); err != nil {
			return nil, err
		} else if StartCode(nextbits).IsSlice() == false {
			break
		}
	}

	return frame, self.next_start_code()
}
