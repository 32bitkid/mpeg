package ts_test

import "testing"
import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/mpeg/ts"

func TestAdaptationField(t *testing.T) {
	reader := bitreader.NewBitReader(adaptationFieldReader())
	packet, err := ts.NewPacket(reader)
	if err != nil {
		t.Fatal(err)
	}

	if expected, actual := packet.AdaptationFieldControl, ts.FieldThenPayload; actual != expected {
		t.Fatalf("unexpected AdaptationFieldControl. expected %d, got %d", actual, expected)
	}

	if packet.AdaptationField == nil {
		t.Fatal("exptected adaptation field to be set")
	}

	if expected, actual := 80, len(packet.Payload); actual != expected {
		t.Fatalf("payload was not the correct size. expected %d, got %d", expected, actual)
	}

}
