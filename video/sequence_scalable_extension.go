package video

import "github.com/32bitkid/bitreader"

type ScalableMode uint32

const (
	DataPartitioning ScalableMode = iota
	SpatialScalability
	SNRScalability
	TemporalScalability
)

type SequenceScalableExtension struct {
	scalable_mode ScalableMode // 2 uimsbf
	layer_id      uint32       // 4 uimsbf

	lower_layer_prediction_horizontal_size uint32 // 14 uimsbf
	marker_bit                             uint32 // 1 bslbf
	lower_layer_prediction_vertical_size   uint32 // 14 uimsbf
	horizontal_subsampling_factor_m        uint32 // 5 uimsbf
	horizontal_subsampling_factor_n        uint32 // 5 uimsbf
	vertical_subsampling_factor_m          uint32 // 5 uimsbf
	vertical_subsampling_factor_n          uint32 // 5 uimsbf

	picture_mux_enable          uint32 // 1 uimsbf
	mux_to_progressive_sequence uint32 // 1 uimsbf
	picture_mux_order           uint32 // 3 uimsbf
	picture_mux_factor          uint32 // 3 uimsbf
}

func sequence_scalable_extension(br bitreader.BitReader) (*SequenceScalableExtension, error) {

	err := SequenceScalableExtensionID.Assert(br)
	if err != nil {
		return nil, err
	}

	sse := SequenceScalableExtension{}

	val, err := br.Read32(2)
	if err != nil {
		return nil, err
	}

	sse.scalable_mode = ScalableMode(val)
	sse.layer_id, err = br.Read32(4)
	if err != nil {
		return nil, err
	}

	panic(sse.scalable_mode)

	/*
		scalable_mode 2 uimsbf
		layer_id 4 uimsbf
		if (scalable_mode == “spatial scalability”) {
		lower_layer_prediction_horizontal_size 14 uimsbf
		marker_bit 1 bslbf
		lower_layer_prediction_vertical_size 14 uimsbf
		horizontal_subsampling_factor_m 5 uimsbf
		horizontal_subsampling_factor_n 5 uimsbf
		vertical_subsampling_factor_m 5 uimsbf
		vertical_subsampling_factor_n 5 uimsbf
		}
		if ( scalable_mode == “temporal scalability” ) {
		picture_mux_enable 1 uimsbf
		if ( picture_mux_enable )
		mux_to_progressive_sequence 1 uimsbf
		picture_mux_order 3 uimsbf
		picture_mux_factor 3 uimsbf
		}
	*/
	next_start_code(br)
	panic("unsupported: sequence_scalable_extension")
}
