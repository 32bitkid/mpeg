package ts

type PacketTester func(*Packet) bool

func (pt PacketTester) Not() PacketTester {
	return func(p *Packet) bool { return !pt(p) }
}

func (pt PacketTester) And(other PacketTester) PacketTester {
	return func(p *Packet) bool { return pt(p) && other(p) }
}

func IsPID(pid uint32) PacketTester {
	return func(p *Packet) bool { return p.PID == pid }
}

var IsPayloadUnitStart PacketTester = func(p *Packet) bool {
	return p.PayloadUnitStartIndicator
}
