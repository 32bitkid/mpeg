package video

import "testing"
import "bytes"
import "github.com/32bitkid/mpeg/util"

func TestSequenceHeader(t *testing.T) {
	testData := []byte{
		0x00, 0x00, 0x01, 0xB3, 0x50, 0x02, 0xD0, 0x37, 0x6D,
		0xDD, 0x2F, 0x41, 0x10, 0x11, 0x11, 0x12, 0x12, 0x12,
		0x13, 0x13, 0x13, 0x13, 0x14, 0x14, 0x14, 0x14, 0x14,
		0x15, 0x15, 0x15, 0x15, 0x15, 0x15, 0x16, 0x16, 0x16,
		0x16, 0x16, 0x16, 0x16, 0x17, 0x17, 0x17, 0x17, 0x17,
		0x17, 0x17, 0x17, 0x18, 0x18, 0x18, 0x18, 0x18, 0x18,
		0x18, 0x19, 0x19, 0x19, 0x19, 0x19, 0x19, 0x1A, 0x1A,
		0x1A, 0x1A, 0x1A, 0x1B, 0x1B, 0x1B, 0x1B, 0x1C, 0x1C,
		0x1C, 0x1D, 0x1D, 0x1E}

	br := util.NewSimpleReader32(bytes.NewReader(testData))

	actual, err := sequence_header(br)
	if err != nil {
		t.Fatal(err)
	}

	expected := sequence_header{
		horizontal_size_value:       1280,
		vertical_size_value:         720,
		aspect_ratio_information:    3,
		frame_rate_code:             7,
		bit_rate_value:              112500,
		vbv_buffer_size_value:       488,
		constrained_parameters_flag: false,
		intra_quantiser_matrix:      [...]byte{8, 16, 19, 22, 26, 27, 29, 34, 16, 16, 22, 24, 27, 29, 34, 37, 19, 22, 26, 27, 29, 34, 34, 38, 22, 22, 26, 27, 29, 34, 37, 40, 22, 26, 27, 29, 32, 35, 40, 48, 26, 27, 29, 32, 35, 40, 48, 58, 26, 27, 29, 34, 38, 46, 56, 69, 27, 29, 35, 38, 46, 56, 69, 83},
		non_intra_quantizer_matrix:  [...]byte{16, 17, 17, 18, 18, 18, 19, 19, 19, 19, 20, 20, 20, 20, 20, 21, 21, 21, 21, 21, 21, 22, 22, 22, 22, 22, 22, 22, 23, 23, 23, 23, 23, 23, 23, 23, 24, 24, 24, 24, 24, 24, 24, 25, 25, 25, 25, 25, 25, 26, 26, 26, 26, 26, 27, 27, 27, 27, 28, 28, 28, 29, 29, 30},
	}

	if actual.horizontal_size_value != expected.horizontal_size_value {
		t.Fail()
	}

	if actual.vertical_size_value != expected.vertical_size_value {
		t.Fail()
	}

	if actual.aspect_ratio_information != expected.aspect_ratio_information {
		t.Fail()
	}

	if actual.frame_rate_code != expected.frame_rate_code {
		t.Fail()
	}

	if actual.bit_rate_value != expected.bit_rate_value {
		t.Fail()
	}

	if actual.vbv_buffer_size_value != expected.vbv_buffer_size_value {
		t.Fail()
	}

	if actual.constrained_parameters_flag != expected.constrained_parameters_flag {
		t.Fail()
	}

	if actual.intra_quantiser_matrix != expected.intra_quantiser_matrix {
		t.Fail()
	}

	if actual.non_intra_quantizer_matrix != expected.non_intra_quantizer_matrix {
		t.Fail()
	}

}
