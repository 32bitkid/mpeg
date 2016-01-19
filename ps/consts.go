package ps

const StartCodePrefix uint32 = 0x000001

type StartCode uint32

const (
	PackStartCode         StartCode = (StartCodePrefix << 8) | 0xBA
	ProgramEndCode        StartCode = (StartCodePrefix << 8) | 0xB9
	SystemHeaderStartCode StartCode = (StartCodePrefix << 8) | 0xBB
)
