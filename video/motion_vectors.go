package video

import "github.com/32bitkid/huffman"
import "github.com/32bitkid/bitreader"

type motionVectors [2][2][2]int
type motionVectorPredictions motionVectors

func absInt(in int) int {
	if in < 0 {
		return -in
	}
	return in
}

type motionVectorsFormed uint

const (
	motionVectorsFormed_None         = motionVectorsFormed(0)
	motionVectorsFormed_FrameForward = motionVectorsFormed(1 << (iota - 1))
	motionVectorsFormed_FrameBackward
)

func (mvf *motionVectorsFormed) set(mb_type *MacroblockType, pct PictureCodingType) {
	switch {
	case mb_type.macroblock_intra:
		*mvf = motionVectorsFormed_None
	case pct == PFrame &&
		mb_type.macroblock_intra == false &&
		mb_type.macroblock_motion_forward == false &&
		mb_type.macroblock_motion_backward == false:
		*mvf = motionVectorsFormed_FrameForward
	case mb_type.macroblock_motion_forward && mb_type.macroblock_motion_backward:
		*mvf = motionVectorsFormed_FrameForward | motionVectorsFormed_FrameBackward
	case mb_type.macroblock_motion_forward:
		*mvf = motionVectorsFormed_FrameForward
	case mb_type.macroblock_motion_backward:
		*mvf = motionVectorsFormed_FrameBackward
	}
}

type motionVectorData struct {
	info                  motionVectorInfo
	code                  motionVectors
	residual              motionVectors
	vertical_field_select [2][2]uint32

	predictions motionVectorPredictions
	actual      motionVectors
	previous    motionVectorsFormed
}

func (motionVector *motionVectorData) update_actual(r, s, t int, f_code fCode) {

	code := motionVector.code
	residual := motionVector.residual

	r_size := f_code[s][t] - 1
	f := 1 << r_size
	high := (16 * f) - 1
	low := -16 * f
	_range := 32 * f

	var delta int
	if f == 1 || code[r][s][t] == 0 {
		delta = code[r][s][t]
	} else {
		delta = ((absInt(code[r][s][t]) - 1) * f) + residual[r][s][t] + 1
		if code[r][s][t] < 0 {
			delta = -delta
		}
	}

	prediction := motionVector.predictions[r][s][t]

	vector := prediction + delta
	if vector < low {
		vector += _range
	}
	if vector > high {
		vector -= _range
	}

	motionVector.predictions[r][s][t] = vector
	motionVector.actual[r][s][t] = vector
}

func (mvd *motionVectorData) reset() {
	mvd.predictions[0][0][0] = 0
	mvd.predictions[0][0][1] = 0
	mvd.predictions[0][1][0] = 0
	mvd.predictions[0][1][1] = 0
	mvd.predictions[1][0][0] = 0
	mvd.predictions[1][0][1] = 0
	mvd.predictions[1][1][0] = 0
	mvd.predictions[1][1][1] = 0
	mvd.clear_actual(0)
	mvd.clear_actual(1)
}

func (mvd *motionVectorData) clear_actual(s int) {
	// First Motion Vector
	mvd.actual[0][s][0] = 0
	mvd.actual[0][s][1] = 0

	// Second motion vector
	mvd.actual[1][s][0] = 0
	mvd.actual[1][s][1] = 0
}

func (fp *VideoSequence) motion_vectors(s int, mb *Macroblock, mvd *motionVectorData) error {

	f_code := fp.PictureCodingExtension.f_code

	mvd.info = mv_info(fp, mb)

	motion_vector_part := func(r, s, t int) error {
		if code, err := decodeMotionCode(fp); err != nil {
			return err
		} else {
			mvd.code[r][s][t] = code
		}
		if f_code[s][t] != 1 && mvd.code[r][s][t] != 0 {
			r_size := uint(f_code[s][t] - 1)

			if code, err := fp.Read32(r_size); err != nil {
				return err
			} else {
				mvd.residual[r][s][t] = int(code)
			}
		}
		if mvd.info.dmv == 1 {
			panic("unsupported: dmv[]")
		}

		mvd.update_actual(r, s, t, f_code)

		return nil
	}

	motion_vector := func(r, s int) error {
		if err := motion_vector_part(r, s, 0); err != nil {
			return err
		}
		if err := motion_vector_part(r, s, 1); err != nil {
			return err
		}
		return nil
	}

	if mvd.info.motion_vector_count == 1 {
		if mvd.info.motion_vector_format == motionVectorFormat_Field && mvd.info.dmv != 1 {
			if val, err := fp.Read32(1); err != nil {
				return err
			} else {
				mvd.vertical_field_select[0][s] = val
			}
		}
		return motion_vector(0, s)
	} else {
		if val, err := fp.Read32(1); err != nil {
			return err
		} else {
			mvd.vertical_field_select[0][s] = val
		}
		if err := motion_vector(0, s); err != nil {
			return err
		}
		if val, err := fp.Read32(1); err != nil {
			return err
		} else {
			mvd.vertical_field_select[1][s] = val
		}
		if err := motion_vector(1, s); err != nil {
			return err
		}
	}

	return nil
}

func decodeMotionCode(br bitreader.BitReader) (int, error) {
	val, err := motionCodeDecoder.Decode(br)
	if err != nil {
		return 0, err
	} else if code, ok := val.(int); ok {
		return code, nil
	} else {
		return 0, huffman.ErrMissingHuffmanValue
	}
}

var motionCodeDecoder = huffman.NewHuffmanDecoder(huffman.HuffmanTable{
	"0000 0011 001 ": -16,
	"0000 0011 011 ": -15,
	"0000 0011 101 ": -14,
	"0000 0011 111 ": -13,
	"0000 0100 001 ": -12,
	"0000 0100 011 ": -11,
	"0000 0100 11 ":  -10,
	"0000 0101 01 ":  -9,
	"0000 0101 11 ":  -8,
	"0000 0111 ":     -7,
	"0000 1001 ":     -6,
	"0000 1011 ":     -5,
	"0000 111 ":      -4,
	"0001 1 ":        -3,
	"0011 ":          -2,
	"011 ":           -1,
	"1":              0,
	"010":            1,
	"0010":           2,
	"0001 0":         3,
	"0000 110":       4,
	"0000 1010":      5,
	"0000 1000":      6,
	"0000 0110":      7,
	"0000 0101 10":   8,
	"0000 0101 00":   9,
	"0000 0100 10":   10,
	"0000 0100 010":  11,
	"0000 0100 000":  12,
	"0000 0011 110":  13,
	"0000 0011 100":  14,
	"0000 0011 010":  15,
	"0000 0011 000":  16,
})
