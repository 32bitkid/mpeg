package pes_test

import "testing"
import "bytes"
import "github.com/32bitkid/mpeg/pes"
import "github.com/32bitkid/mpeg/util"

func TestPacketWithExtensionFlag(t *testing.T) {
	br := util.NewSimpleBitReader(bytes.NewReader(packetWithExtensionFlag))

	p, err := pes.ReadPacketFrom(br)
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
