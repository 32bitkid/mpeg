package video

import "io"
import "github.com/32bitkid/mpeg/util"
import "errors"

var ErrUnsupportedVideoStream_ISO_IEC_11172_2 = errors.New("unsupported video stream ISO/IEC 11172-2")

type FrameProvider interface {
	Next() (interface{}, error)
}

func NewFrameProvider(source io.Reader) FrameProvider {
	return &frameProvider{
		NewVideoSequence(util.NewBitReader(source)),
	}
}

type frameProvider struct {
	VideoSequence
}

func (self *frameProvider) Next() (interface{}, error) {

	// Align to next start code
	err := next_start_code(self)
	if err != nil {
		panic(err)
	}

	// Read sequence_header
	err = self.sequence_header()
	if err != nil {
		panic(err)
	}

	// peek for sequence_extension
	val, err := self.Peek32(32)
	if err != nil {
		panic(err)
	}

	if val == ExtensionStartCode {

		err := self.sequence_extension()

		for {
			err = extension_and_user_data(0, self)
			if err != nil {
				panic("extension_and_user_data")
			}

			for {
				nextbits, err := self.Peek32(32)
				if err != nil {
					panic("Peek32")
				}

				if StartCode(nextbits) == GroupStartCode {
					err = self.group_of_pictures_header()
					if err != nil {
						panic("group_of_pictures_header")
					}
					err = extension_and_user_data(1, self)
					if err != nil {
						panic("extension_and_user_data")
					}
				}
				err = self.picture_header()
				if err != nil {
					panic("picture_header")
				}
				err = self.picture_coding_extension()
				if err != nil {
					panic("picture_coding_extension")
				}
				err = extension_and_user_data(2, self)
				if err != nil {
					panic("extension_and_user_data")
				}

				self.picture_data()

				panic("not implemented")
			}

			val, err := self.Peek32(32)
			log.Printf("%x\n", val)
			if err != nil {
				panic("Peek32")
			}

			if val == SequenceEndStartCode {
				break
			}
		}

		err = self.Trash(32)

		return nil, err
	} else {
		// Stream is MPEG-1 Video
		return nil, ErrUnsupportedVideoStream_ISO_IEC_11172_2
	}

}
