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
		br: util.NewBitReader(source),
	}
}

type frameProvider struct {
	br util.BitReader32
}

func (fp *frameProvider) Next() (interface{}, error) {
	br := fp.br

	// Align to next start code
	err := next_start_code(br)
	if err != nil {
		panic(err)
	}

	// Read sequence_header
	sqh, err := sequence_header(br)
	if err != nil {
		panic(err)
	}

	// peek for sequence_extension
	val, err := br.Peek32(32)
	if err != nil {
		panic(err)
	}

	if val == ExtensionStartCode {

		se, err := sequence_extension(br)

		log.Printf("%#v\n", se)

		for {
			err = extension_and_user_data(0, br)
			if err != nil {
				panic("extension_and_user_data")
			}

			for {
				nextbits, err := br.Peek32(32)
				if err != nil {
					panic("Peek32")
				}

				if StartCode(nextbits) == GroupStartCode {
					_, err = group_of_pictures_header(br)
					if err != nil {
						panic("group_of_pictures_header")
					}
					err = extension_and_user_data(1, br)
					if err != nil {
						panic("extension_and_user_data")
					}
				}
				_, err = picture_header(br)
				if err != nil {
					panic("picture_header")
				}
				_, err = picture_coding_extension(br)
				if err != nil {
					panic("picture_coding_extension")
				}
				err = extension_and_user_data(2, br)
				if err != nil {
					panic("extension_and_user_data")
				}

				picture_data(br)

				panic("not implemented")
			}

			val, err := br.Peek32(32)
			log.Printf("%x\n", val)
			if err != nil {
				panic("Peek32")
			}

			if val == SequenceEndStartCode {
				break
			}
		}

		err = br.Trash(32)

		return sqh, err
	} else {
		// Stream is MPEG-1 Video
		return nil, ErrUnsupportedVideoStream_ISO_IEC_11172_2
	}

}
