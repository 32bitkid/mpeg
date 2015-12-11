package ts_test

import "testing"
import "io"
import "bytes"
import "github.com/32bitkid/mpeg/util"
import "github.com/32bitkid/mpeg/ts"

func TestPacketParsing(t *testing.T) {
	reader := util.NewSimpleBitReader(nullPacketReader())
	packet := &ts.Packet{}

	err := packet.ReadFrom(reader)
	if err != nil {
		t.Fatal(err)
	}
	if packet.PID != nullPacketPID {
		t.Fatalf("unexpected PID. expected %x, got %x", nullPacketPID, packet.PID)
	}
}

func TestIncompletePacket(t *testing.T) {
	reader := util.NewSimpleBitReader(io.LimitReader(nullPacketReader(), 4))
	_, err := ts.NewPacket(reader)
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("unexpected error. expected %v, got %v", io.ErrUnexpectedEOF, err)
	}
}

func TestAdaptationField(t *testing.T) {
	reader := util.NewSimpleBitReader(bytes.NewReader(adaptationFieldData))
	packet := &ts.Packet{}

	err := packet.ReadFrom(reader)

	if err != nil {
		t.Fatal(err)
	}

	if packet.AdaptationFieldControl != ts.FieldThenPayload {
		t.Fatalf("unexpected AdaptationFieldControl. expected %d, got %d", ts.FieldThenPayload, packet.AdaptationFieldControl)
	}

	if packet.AdaptationField == nil {
		t.Fatal("exptected adaptation field to be set")
	}

	if packet.AdaptationField.Length != uint32(adaptationFieldData[4]) {
		t.Fatal("unexpected Length. expected %d, got %d", adaptationFieldData[4], packet.AdaptationField.Length)
	}

	if len(packet.Payload) != 184-int(adaptationFieldData[4])-1 {
		t.Fatal("payload was not the correct size")
	}
}
