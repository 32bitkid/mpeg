package video

import "io"
import "image"
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
	frameCounter uint32
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

func (vs *VideoSequence) Size() (width int, height int) {
	if vs.SequenceHeader == nil {
		return -1, -1
	}
	width = int(vs.SequenceHeader.horizontal_size_value)
	height = int(vs.SequenceHeader.vertical_size_value)
	if vs.SequenceExtension != nil {
		width |= int(vs.SequenceExtension.horizontal_size_extension << 12)
		height |= int(vs.SequenceExtension.vertical_size_extension << 12)
	}
	return
}

// Next() will return the next frame of video decoded from the video stream.
func (vs *VideoSequence) Next() (*image.YCbCr, error) {
	// Try to get a previously decoded frame out of the frameStore.
	if img := vs.frameStore.tryGet(vs.frameCounter); img != nil {
		vs.frameCounter++
		return img, nil
	}

	// Step until a temporal match is found.
	for {
		if img, err := vs.step(); err != nil {
			return nil, err
		} else if vs.PictureHeader.temporal_reference == vs.frameCounter {
			vs.frameCounter++
			return img, nil
		}
	}
}
