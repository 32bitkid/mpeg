package video

import "io"
import "github.com/32bitkid/bitreader"

type sequenceHeaders struct {
	*SequenceHeader
	*SequenceExtension
}

type pictureHeaders struct {
	*GroupOfPicturesHeader
	*PictureHeader
	*PictureCodingExtension
}

type VideoSequence struct {
	bitreader.BitReader

	sequenceHeaders
	pictureHeaders

	quantisationMatricies [4]quantisationMatrix
	frameStore
}

func NewVideoSequence(r io.Reader) *VideoSequence {
	return &VideoSequence{
		BitReader: bitreader.NewBitReader(r),
	}
}

// AlignTo will trash all bits until the stream is aligned with the desired start code or error is produced.
func (br *VideoSequence) AlignTo(startCode StartCode) error {
	if !br.IsByteAligned() {
		if _, err := br.ByteAlign(); err != nil {
			return err
		}
	}

	for {
		if val, err := br.Peek32(32); err != nil {
			return err
		} else if StartCode(val) == startCode {
			return nil
		} else if err := br.Trash(8); err != nil {
			return err
		}
	}
}

func (vs *VideoSequence) sequence_extension() (err error) {
	vs.SequenceExtension, err = sequence_extension(vs)
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
