package video

import "github.com/32bitkid/bitreader"
import "errors"

var (
	ErrUnexpectedNonZeroByte = errors.New("unexpected non-zero byte")
)

func (self *VideoSequence) next_start_code() (err error) {
	return next_start_code(self)
}

func next_start_code(br bitreader.BitReader) (err error) {
	for !br.IsByteAligned() {
		_, err = br.ByteAlign()
		if err != nil {
			return err
		}
	}

	var val uint32

	for {
		val, err = br.Peek32(24)

		if err != nil {
			return err
		}

		if val == StartCodePrefix {
			return nil
		}

		val, err = br.Read32(8)

		if err != nil {
			return err
		}

		if val != StuffingByte {
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
