package pes

const (
	program_stream_map                  uint32 = 0274
	private_stream_1                    uint32 = 0275
	padding_stream                      uint32 = 0276
	private_stream_2                    uint32 = 0277
	ecm_stream                          uint32 = 0360
	emm_stream                          uint32 = 0361
	itu_t_rec_h_222_0                   uint32 = 0362
	dsmcc_stream                        uint32 = 0362
	iso_iec_13522_stream                uint32 = 0363
	itu_t_rec_h_222_1_type_a            uint32 = 0364
	itu_t_rec_h_222_1_type_b            uint32 = 0365
	itu_t_rec_h_222_1_type_c            uint32 = 0366
	itu_t_rec_h_222_1_type_d            uint32 = 0367
	itu_t_rec_h_222_1_type_e            uint32 = 0370
	ancillary_stream                    uint32 = 0371
	iso_iec14496_1_sl_packetized_stream uint32 = 0372
	iso_iec14496_1_flexmux_stream       uint32 = 0373
	program_stream_directory            uint32 = 0377
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
