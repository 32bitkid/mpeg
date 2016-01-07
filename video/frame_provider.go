package video

import "io"
import "github.com/32bitkid/bitreader"
import "errors"
import "image"

var ErrUnsupportedVideoStream_ISO_IEC_11172_2 = errors.New("unsupported video stream ISO/IEC 11172-2")

func NewFrameProvider(source io.Reader) *VideoSequence {
	return &VideoSequence{BitReader: bitreader.NewBitReader(source)}
}

func (self *VideoSequence) Next() (image.Image, error) {
	// align to next start code
	if err := next_start_code(self); err != nil {
		panic(err)
	}

	// read sequence_header
	if err := self.sequence_header(); err != nil {
		panic(err)
	}

	// peek for sequence_extension
	if val, err := self.Peek32(32); err != nil {
		panic(err)
	} else if StartCode(val) == ExtensionStartCode {
		if err := self.sequence_extension(); err != nil {
			panic(err)
		}

		for {
			if err := self.extension_and_user_data(0); err != nil {
				panic("extension_and_user_data: " + err.Error())
			}

			for {
				if nextbits, err := self.Peek32(32); err != nil {
					panic("Peek32")
				} else if StartCode(nextbits) == GroupStartCode {
					if err := self.group_of_pictures_header(); err != nil {
						panic("group_of_pictures_header: " + err.Error())
					}
					if err := self.extension_and_user_data(1); err != nil {
						panic("extension_and_user_data:" + err.Error())
					}
				}

				if err := self.picture_header(); err != nil {
					panic("picture_header: " + err.Error())
				}

				if err := self.picture_coding_extension(); err != nil {
					panic("picture_coding_extension: " + err.Error())
				}

				if err := self.extension_and_user_data(2); err != nil {
					panic("extension_and_user_data: " + err.Error())
				}

				{
					frame, err := self.picture_data()
					if err != nil {
						panic(err)
					}
					if self.PictureHeader.picture_coding_type == IFrame {
						self.frameStore.forward = frame
					}
					if true {
						return frame, nil
					}
				}

				if nextbits, err := self.Peek32(32); err != nil {
					panic("peeking: " + err.Error())
				} else if StartCode(nextbits) == PictureStartCode {
					continue
				} else if StartCode(nextbits) == GroupStartCode {
					continue
				} else {
					break
				}
			}

			if nextbits, err := self.Peek32(32); err != nil {
				panic("Peek32")
			} else if StartCode(nextbits) == SequenceEndStartCode {
				break
			}

			if err := self.sequence_header(); err != nil {
				panic(err)
			}

			if err := self.sequence_extension(); err != nil {
				panic(err)
			}
		}

		// SequenceEndStartCode
		return nil, self.Trash(32)
	} else {
		// Stream is MPEG-1 Video
		return nil, ErrUnsupportedVideoStream_ISO_IEC_11172_2
	}

}
