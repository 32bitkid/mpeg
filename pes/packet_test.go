package pes_test

import "testing"
import "os"
import "bytes"
import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/mpeg_go/pes"

func TestBasicPacketParsing(t *testing.T) {

	datafile := "testdata/football.pes"
	fi, _ := os.Stat(datafile)
	data, _ := os.Open(datafile)

	br := bitreader.NewReader32(data)
	packet, err := pes.ReadPacket(br, int(fi.Size()))

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var expectedStreamID uint32 = 0xe0
	if packet.StreamID != expectedStreamID {
		t.Fatalf("incorrect stream id. expected %#x got %#x", expectedStreamID, packet.StreamID)
	}

	if packet.PacketLength != 0 {
		t.Fatalf("incorrect packet length. expected %#x got %#x", 0, packet.PacketLength)
	}

	if packet.Header == nil {
		t.Fatalf("expected header")
	}
}

func TestPacketWithExtensionFlag(t *testing.T) {
	br := bitreader.NewReader32(bytes.NewReader(packetWithExtensionFlag))

	p, err := pes.ReadPacket(br, -1)
	if err != nil {
		t.Fatal(err)
	}

	if p.Header.Extension == nil {
		t.Fatal("expected a header extension")
	}

	val, err := br.Peek32(32)
	if err != nil {
		t.Fatal(err)
	}
	if val != 0xffffffff {
		t.Fatal("maker not found")
	}
}
