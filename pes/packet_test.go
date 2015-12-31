package pes_test

import "testing"
import "github.com/32bitkid/mpeg/pes"
import "github.com/32bitkid/bitreader"

func TestPacketWithExtensionFlag(t *testing.T) {
	br := bitreader.NewBitReader(packetWithExtensionFlag())

	p, err := pes.NewPacket(br)
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
