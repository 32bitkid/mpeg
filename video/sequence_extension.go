package video

import "github.com/32bitkid/bitreader"

type SequenceExtension struct {
	profile_and_level_indication uint32       // 8 uimsbf
	progressive_sequence         uint32       // 1 uimsbf
	chroma_format                ChromaFormat // 2 uimsbf
	horizontal_size_extension    uint32       // 2 uimsbf
	vertical_size_extension      uint32       // 2 uimsbf
	bit_rate_extension           uint32       // 12 uimsbf
	vbv_buffer_size_extension    uint32       // 8 uimsbf
	low_delay                    uint32       // 1 uimsbf
	frame_rate_extension_n       uint32       // 2 uimsbf
	frame_rate_extension_d       uint32       // 5 uimsbf
}

func sequence_extension(br bitreader.BitReader) (*SequenceExtension, error) {

	err := ExtensionStartCode.assert(br)
	if err != nil {
		return nil, err
	}

	err = SequenceExtensionID.assert(br)
	if err != nil {
		return nil, err
	}

	se := SequenceExtension{}

	se.profile_and_level_indication, err = br.Read32(8)
	if err != nil {
		return nil, err
	}

	se.progressive_sequence, err = br.Read32(1)
	if err != nil {
		return nil, err
	}

	if chroma_format, err := br.Read32(2); err != nil {
		return nil, err
	} else {
		se.chroma_format = ChromaFormat(chroma_format)
	}

	se.horizontal_size_extension, err = br.Read32(2)
	if err != nil {
		return nil, err
	}

	se.vertical_size_extension, err = br.Read32(2)
	if err != nil {
		return nil, err
	}

	se.bit_rate_extension, err = br.Read32(12)
	if err != nil {
		return nil, err
	}

	err = marker_bit(br)
	if err != nil {
		return nil, err
	}

	se.vbv_buffer_size_extension, err = br.Read32(8)
	if err != nil {
		return nil, err
	}

	se.low_delay, err = br.Read32(1)
	if err != nil {
		return nil, err
	}

	se.frame_rate_extension_n, err = br.Read32(2)
	if err != nil {
		return nil, err
	}

	se.frame_rate_extension_d, err = br.Read32(5)
	if err != nil {
		return nil, err
	}

	return &se, next_start_code(br)
}
