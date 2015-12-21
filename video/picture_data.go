package video

import "image"
import "image/png"
import "os"

func (self *VideoSequence) picture_data() (err error) {

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

	frame := image.NewYCbCr(r, subsampleRatio)

	for {
		err := self.slice(frame)
		if err != nil {
			return err
		}

		nextbits, err := self.Peek32(32)
		if err != nil {
			return err
		}

		if !is_slice_start_code(StartCode(nextbits)) {
			break
		}
	}

	f, _ := os.Create("test.png")
	png.Encode(f, frame)

	return self.next_start_code()
}
