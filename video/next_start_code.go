package video

import "github.com/32bitkid/bitreader"
import "errors"

const stuffingByte = 0x00

var ErrUnexpectedNonZeroByte = errors.New("unexpected non-zero byte")

func (self *VideoSequence) next_start_code() (err error) {
	return next_start_code(self)
}

func next_start_code(br bitreader.BitReader) error {
	if !br.IsAligned() {
		if _, err := br.Align(); err != nil {
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
	if marker, err := br.Read1(); err != nil {
		return err
	} else if marker == false {
		return ErrMissingMarkerBit
	}
	return nil
}
