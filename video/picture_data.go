package video

import "image"

func createFrameBuffer(w, h int, cf ChromaFormat) *image.YCbCr {

	horizontalMacroblocks := (w + 15) >> 4
	verticalMacroblocks := (h + 15) >> 4

	r := image.Rect(0, 0, horizontalMacroblocks<<4, verticalMacroblocks<<4)

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

	w, h := self.Size()
	frame = createFrameBuffer(w, h, self.SequenceExtension.chroma_format)

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
