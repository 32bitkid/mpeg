package video

import "github.com/32bitkid/bitreader"
import "errors"

const stuffingByte = 0x00

var (
	ErrUnexpectedNonZeroByte = errors.New("unexpected non-zero byte")
)

func (self *VideoSequence) next_start_code() (err error) {
	return next_start_code(self)
}

func next_start_code(br bitreader.BitReader) error {
	if !br.IsByteAligned() {
		if _, err := br.ByteAlign(); err != nil {
			return err
		}
	}

	for {
		if val, err := br.Peek32(24); err != nil {
			return err
		} else if val == StartCodePrefix {
			return nil
		}

		if val, err := br.Read32(8); err != nil {
			return err
		} else if val != stuffingByte {
			return ErrUnexpectedNonZeroByte
		}
	}
}

func marker_bit(br bitreader.BitReader) error {
	marker, err := br.ReadBit()
	if err != nil {
		return err
	}
	if marker == false {
		return ErrMissingMarkerBit
	}
	return nil
}

func (br *VideoSequence) TrashUntil(startCode StartCode) error {
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
