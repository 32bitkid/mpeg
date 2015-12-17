package video

type StartCode uint32

const (
	StuffingByte = 0x00
)

const (
	StartCodePrefix = 0x000001

	PictureStartCode    = 0x00
	ReservedStartCode_1 = 0xB0
	ReservedStartCode_2 = 0xB1
	UserDataStartCode   = 0xB2
	SequenceHeaderCode  = 0xB3
	SequenceErrorCode   = 0xB4
	ExtensionCode       = 0xB5
	ReservedStartCode_3 = 0xB6
	SequenceEndCode     = 0xB7
	GroupCode           = 0xB8

	SequenceHeaderStartCode StartCode = (StartCodePrefix << 8) + SequenceHeaderCode
	ExtensionStartCode                = (StartCodePrefix << 8) + ExtensionCode
	SequenceEndStartCode              = (StartCodePrefix << 8) + SequenceEndCode
	GroupStartCode                    = (StartCodePrefix << 8) + GroupCode
)

// slice_start_code 01 through AF
// system start codes (see note) B9 through FF
