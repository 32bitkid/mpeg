package ts_test

import "testing"
import "io"
import "io/ioutil"
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
		if expected, actual := byte(255), val; expected != actual {
			t.Fatal("Unexpected value read. Expected %d, got %d", expected, actual)
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

func TestStreamWithMultipleParts(t *testing.T) {
	source := fivePacketReader()
	payload := ts.NewPayloadReader(source, ts.IsPID(dataPacketPID))

	dest, err := ioutil.ReadAll(payload)
	if err != nil {
		t.Fatal(err)
	}

	if expected, actual := ts.MaxPayloadSize*3, len(dest); expected != actual {
		t.Fatalf("dest length is incorrect. expected %d, got %d", expected, actual)
	}
}
