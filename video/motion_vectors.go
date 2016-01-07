package video

import "github.com/32bitkid/huffman"
import "github.com/32bitkid/bitreader"

type motionVectors [2][2][2]int

func absInt(in int) int {
	if in < 0 {
		return -in
	}
	return in
}

type motionVectorData struct {
	motion_code                  motionVectors
	motion_residual              motionVectors
	motion_vertical_field_select [2][2]uint32
}

func (mvD motionVectorData) calc(r, s, t int, f_code FCode, pMV *motionVectorPredictions) int {

	motion_code := mvD.motion_code
	motion_residual := mvD.motion_residual

	r_size := f_code[s][t] - 1
	f := 1 << r_size
	high := (16 * f) - 1
	low := -16 * f
	_range := 32 * f

	var delta int
	if f == 1 || motion_code[r][s][t] == 0 {
		delta = motion_code[r][s][t]
	} else {
		delta = ((absInt(motion_code[r][s][t]) - 1) * f) + motion_residual[r][s][t] + 1
		if motion_code[r][s][t] < 0 {
			delta = -delta
		}
	}

	prediction := pMV[r][s][t]

	vector := prediction + delta
	if vector < low {
		vector += _range
	}
	if vector > high {
		vector -= _range
	}

	pMV[r][s][t] = vector

	return vector
}

type motionVectorPredictions [2][2][2]int

func (pMV *motionVectorPredictions) reset() {
	pMV[0][0][0] = 0
	pMV[0][0][1] = 0
	pMV[0][1][0] = 0
	pMV[0][1][1] = 0
	pMV[1][0][0] = 0
	pMV[1][0][1] = 0
	pMV[1][1][0] = 0
	pMV[1][1][1] = 0
}

func (fp *VideoSequence) motion_vectors(s int, mb *Macroblock, mvd *motionVectorData) error {

	f_code := fp.PictureCodingExtension.f_code

	mv_count, mv_format, dmv := mv_info(fp, mb)

	motion_vector_part := func(r, s, t int) error {
		if code, err := decodeMotionCode(fp); err != nil {
			return err
		} else {
			mvd.motion_code[r][s][t] = code
		}
		if f_code[s][t] != 1 && mvd.motion_code[r][s][t] != 0 {
			r_size := uint(f_code[s][t] - 1)

			if code, err := fp.Read32(r_size); err != nil {
				return err
			} else {
				mvd.motion_residual[r][s][t] = int(code)
			}
		}
		if dmv == 1 {
			panic("unsupported: dmv[]")
		}
		return nil
	}

	motion_vector := func(r, s int) error {
		err := motion_vector_part(r, s, 0)
		if err != nil {
			return err
		}
		return motion_vector_part(r, s, 1)
	}

	if mv_count == 1 {
		if mv_format == MotionVectorFormat_Field && dmv != 1 {
			if val, err := fp.Read32(1); err != nil {
				return err
			} else {
				mvd.motion_vertical_field_select[0][s] = val
			}
		}
		return motion_vector(0, s)
	} else {
		if val, err := fp.Read32(1); err != nil {
			return err
		} else {
			mvd.motion_vertical_field_select[0][s] = val
		}
		if err := motion_vector(0, s); err != nil {
			return err
		}
		if val, err := fp.Read32(1); err != nil {
			return err
		} else {
			mvd.motion_vertical_field_select[1][s] = val
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
