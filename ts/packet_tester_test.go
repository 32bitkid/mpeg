package ts_test

import "testing"
import "github.com/32bitkid/mpeg/ts"

func TestBasicPacketTester(t *testing.T) {
	test := ts.IsPID(0x01)
	packetA := &ts.Packet{
		PID: 0x01,
	}
	packetB := &ts.Packet{
		PID: 0x02,
	}
	if test(packetA) != true {
		t.Fatal("test failed")
	}
	if test(packetB) != false {
		t.Fatal("test failed")
	}
}

func TestAndingPacketTesters(t *testing.T) {
	pidTest := ts.IsPID(0x01)
	test := pidTest.And(ts.IsPayloadUnitStart)

	packetA := &ts.Packet{
		PID: 0x01,
		PayloadUnitStartIndicator: false,
	}

	packetB := &ts.Packet{
		PID: 0x01,
		PayloadUnitStartIndicator: true,
	}

	packetC := &ts.Packet{
		PID: 0x02,
		PayloadUnitStartIndicator: true,
	}

	if test(packetA) != false {
		t.Fatal("test failed")
	}

	if test(packetB) != true {
		t.Fatal("test failed")
	}

	if test(packetC) != false {
		t.Fatal("test failed")
	}
}
