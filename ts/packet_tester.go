package ts

type PacketTester func(*TsPacket) bool

func (pt PacketTester) Not() PacketTester {
	return func(p *TsPacket) bool { return !pt(p) }
}

func IsPID(pid uint32) PacketTester {
	return func(p *TsPacket) bool { return p.PID == pid }
}

func IsPayloadUnitStart() PacketTester {
	return func(p *TsPacket) bool { return p.PayloadUnitStartIndicator }
}
