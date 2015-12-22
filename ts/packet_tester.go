package ts

// PacketTester defines a function that tests a packet and returns a bool
type PacketTester func(*Packet) bool

func alwaysTrueTester(p *Packet) bool { return true }

// Inverts the packet tester
func (pt PacketTester) Not() PacketTester {
	return func(p *Packet) bool { return !pt(p) }
}

// And joins two packet testers with a logical and
func (pt PacketTester) And(other PacketTester) PacketTester {
	return func(p *Packet) bool { return pt(p) && other(p) }
}

// Or joins two packet testers with a logical or
func (pt PacketTester) Or(other PacketTester) PacketTester {
	return func(p *Packet) bool { return pt(p) || other(p) }
}

//IsPID creates a packet tester that returns true
// if the tested packet matches the selected pid
func IsPID(pid uint32) PacketTester {
	return func(p *Packet) bool { return p.PID == pid }
}

// IsPayloadUnitStart is a packet tester that returns true
// if the tested packet has the PayloadUnitStartIndicator
// flag set to true
var IsPayloadUnitStart PacketTester = func(p *Packet) bool { return p.PayloadUnitStartIndicator }
