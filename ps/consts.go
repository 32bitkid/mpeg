package ps

const StartCodePrefix uint32 = 0x000001

type StartCode uint32

const (
	packCode         uint32 = 0xba
	programEnd       uint32 = 0xb9
	systemHeaderCode uint32 = 0xbb

	PackStartCode         StartCode = StartCode((StartCodePrefix << 8) | packCode)
	ProgramEndCode        StartCode = StartCode((StartCodePrefix << 8) | programEnd)
	SystemHeaderStartCode StartCode = StartCode((StartCodePrefix << 8) | systemHeaderCode)
)
