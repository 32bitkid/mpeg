package ts_test

import "testing"
import "io"
import "github.com/32bitkid/mpeg/ts"

func TestReadingASinglePacket(t *testing.T) {
	source := io.MultiReader(nullPacketReader())
	payload := ts.NewPayloadReader(source, ts.IsPID(nullPacketPID))

	data := make([]byte, ts.MaxPayloadSize)
	_, err := io.ReadAtLeast(payload, data, ts.MaxPayloadSize)

	if err != nil {
		t.Fatal(err)
	}

	for _, val := range data {
		if val != 255 {
			t.Fatal("Unexpected data")
		}
	}
}

func TestReadingTwoPackets(t *testing.T) {
	source := io.MultiReader(nullPacketReader(), nullPacketReader())
	payload := ts.NewPayloadReader(source, ts.IsPID(nullPacketPID))

	data := make([]byte, ts.MaxPayloadSize*2)
	_, err := io.ReadAtLeast(payload, data, ts.MaxPayloadSize*2)

	if err != nil {
		t.Fatal(err)
	}

	for _, val := range data {
		if val != 255 {
			t.Fatal("Unexpected data %d", val)
		}
	}
}

func TestReadingTooMuch(t *testing.T) {
	source := io.MultiReader(nullPacketReader(), nullPacketReader())
	payload := ts.NewPayloadReader(source, ts.IsPID(nullPacketPID))

	data := make([]byte, ts.MaxPayloadSize*3)
	_, err := io.ReadAtLeast(payload, data, ts.MaxPayloadSize*3)

	if err != io.ErrUnexpectedEOF {
		t.Fatal(err)
	}
}
