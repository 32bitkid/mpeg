package video

import "testing"
import "io"
import "bytes"
import "github.com/32bitkid/bitreader"
import "encoding/base64"

func TestSequenceHeader(t *testing.T) {
	testData, _ := base64.StdEncoding.DecodeString(
		`AAABs1AC0Ddt3S9BEBEREhISExMTExQUFBQUFRUVFRUVFhYWFhYWF` +
			`hcXFxcXFxcXGBgYGBgYGBkZGRkZGRoaGhoaGxsbGxwcHB0dHg==`)

	br := bitreader.NewBitReader(bytes.NewReader(testData))

	actual, err := sequence_header(br)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}

	expected := SequenceHeader{
		horizontal_size_value:       1280,
		vertical_size_value:         720,
		aspect_ratio_information:    3,
		frame_rate_code:             7,
		bit_rate_value:              112500,
		vbv_buffer_size_value:       488,
		constrained_parameters_flag: false,
		load_intra_quantiser_matrix: false,
		non_intra_quantizer_matrix: QuantisationMatrix{
			{16, 17, 17, 18, 18, 18, 19, 19},
			{19, 19, 20, 20, 20, 20, 20, 21},
			{21, 21, 21, 21, 21, 22, 22, 22},
			{22, 22, 22, 22, 23, 23, 23, 23},
			{23, 23, 23, 23, 24, 24, 24, 24},
			{24, 24, 24, 25, 25, 25, 25, 25},
			{25, 26, 26, 26, 26, 26, 27, 27},
			{27, 27, 28, 28, 28, 29, 29, 30},
		},
	}

	if actual.horizontal_size_value != expected.horizontal_size_value {
		t.Fatal("horizontal_size_value")
	}

	if actual.vertical_size_value != expected.vertical_size_value {
		t.Fatal("vertical_size_value")
	}

	if actual.aspect_ratio_information != expected.aspect_ratio_information {
		t.Fatal("aspect_ratio_information")
	}

	if actual.frame_rate_code != expected.frame_rate_code {
		t.Fatal("frame_rate_code")
	}

	if actual.bit_rate_value != expected.bit_rate_value {
		t.Fatal("bit_rate_value")
	}

	if actual.vbv_buffer_size_value != expected.vbv_buffer_size_value {
		t.Fatal("vbv_buffer_size_value")
	}

	if actual.constrained_parameters_flag != expected.constrained_parameters_flag {
		t.Fatal("constrained_parameters_flag")
	}

	if actual.load_intra_quantiser_matrix != expected.load_intra_quantiser_matrix {
		t.Fatal("intra_quantiser_matrix")
	}

	if actual.non_intra_quantizer_matrix != expected.non_intra_quantizer_matrix {
		t.Fatal("non_intra_quantizer_matrix")
	}

}
