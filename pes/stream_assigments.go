package pes

const (
	program_stream_map                  = 0xbc // program_stream_map
	private_stream_1                    = 0xbd // private_stream_1
	padding_stream                      = 0xbe // padding_stream
	private_stream_2                    = 0xbf // private_stream_2
	ecm_stream                          = 0xf0 // ECM_stream © ISO/IEC ISO/IEC 13818-1: 1994(E) ITU-T Rec. H.222.0 (1995 E) 37
	emm_stream                          = 0xf1 // EMM_stream
	itu_t_rec_h_222_0                   = 0xf2 // ITU-T Rec. H.222.0 | ISO/IEC 13818-1 Annex A
	dsmcc_stream                        = 0xf2 // ISO/IEC 13818-6_DSMCC_stream
	iso_iec_13522_stream                = 0xf3 // ISO/IEC_13522_stream
	itu_t_rec_h_222_1_type_a            = 0xf4 // ITU-T Rec. H.222.1 type A
	itu_t_rec_h_222_1_type_b            = 0xf5 // ITU-T Rec. H.222.1 type B
	itu_t_rec_h_222_1_type_c            = 0xf6 // ITU-T Rec. H.222.1 type C
	itu_t_rec_h_222_1_type_d            = 0xf7 // ITU-T Rec. H.222.1 type D
	itu_t_rec_h_222_1_type_e            = 0xf8 // ITU-T Rec. H.222.1 type E
	ancillary_stream                    = 0xf9 // ancillary_stream
	iso_iec14496_1_sl_packetized_stream = 0xfa
	iso_iec14496_1_flexmux_stream       = 0xfb
	program_stream_directory            = 0xff // program_stream_directory
)

func hasPESHeader(streamID uint32) bool {
	return streamID != program_stream_map &&
		streamID != padding_stream &&
		streamID != private_stream_2 &&
		streamID != ecm_stream &&
		streamID != emm_stream &&
		streamID != program_stream_directory &&
		streamID != dsmcc_stream &&
		streamID != itu_t_rec_h_222_1_type_e
}
