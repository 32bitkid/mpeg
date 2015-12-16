package ts_test

import "testing"
import "github.com/32bitkid/mpeg/util"
import "github.com/32bitkid/mpeg/ts"

func TestAdaptationField(t *testing.T) {
	reader := util.NewBitReader(adaptationFieldReader())
	packet := new(ts.Packet)

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

	if packet.AdaptationField.Length != uint32(adapationFieldPacket[4]) {
		t.Fatal("unexpected Length. expected %d, got %d", adapationFieldPacket[4], packet.AdaptationField.Length)
	}

	if len(packet.Payload) != 184-int(adapationFieldPacket[4])-1 {
		t.Fatal("payload was not the correct size")
	}
}
