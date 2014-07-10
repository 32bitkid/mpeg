package ts_test

import "testing"
import "github.com/32bitkid/mpeg-go/ts"
import "bytes"
import "io"

var nullPacketPID uint32 = 0x1fff
var nullPacket = []byte{0x47, 0x1f, 0xff, 0x10, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

var dataPacketPID uint32 = 0x21
var dataPacket = []byte{0x47, 0x00, 0x21, 0x1d, 0xba, 0x41, 0xac, 0x59, 0x34, 0x12, 0x2b, 0x6a, 0x61, 0xef, 0x7f, 0xe8, 0xca, 0xaa, 0x7c, 0xbb, 0xa4, 0xe4, 0x45, 0x35, 0x23, 0x25, 0x55, 0x22, 0x2a, 0xe9, 0x3e, 0x13, 0x49, 0xe2, 0xd6, 0x78, 0x7f, 0x77, 0x4f, 0x1a, 0xee, 0xe9, 0xf5, 0xd3, 0x2a, 0x98, 0x21, 0x12, 0x24, 0x41, 0x97, 0xa9, 0xe8, 0x44, 0xf1, 0x5e, 0x45, 0xea, 0x8a, 0x1c, 0xf7, 0xa6, 0xaa, 0x9e, 0xde, 0xea, 0xa7, 0xbb, 0x39, 0xe2, 0x5d, 0xcd, 0x74, 0xf6, 0x72, 0x6e, 0xb7, 0x82, 0x73, 0x6a, 0x80, 0x00, 0x00, 0x01, 0x1b, 0x52, 0x8b, 0x66, 0x92, 0x4b, 0x07, 0x49, 0xfb, 0xd4, 0x85, 0x67, 0x0d, 0x06, 0x18, 0x69, 0xd6, 0xb3, 0x1d, 0x83, 0x39, 0x04, 0x4b, 0xa0, 0x43, 0x73, 0xa1, 0x29, 0x06, 0x80, 0x93, 0x5d, 0x4c, 0x3d, 0xfb, 0x57, 0xce, 0xf5, 0x5e, 0xea, 0x19, 0xe1, 0x94, 0xb9, 0xee, 0x2e, 0x4d, 0xd3, 0xe0, 0x2e, 0x36, 0x46, 0x82, 0x43, 0x7d, 0x0d, 0x6c, 0x04, 0xa0, 0xa4, 0xb2, 0xc0, 0xc1, 0x31, 0x3c, 0xa8, 0x01, 0x98, 0x01, 0xf7, 0xc0, 0x0f, 0xd6, 0x4c, 0x04, 0x9f, 0xf9, 0x1e, 0x6c, 0x03, 0x3b, 0xea, 0xa1, 0xb5, 0x92, 0xcd, 0x01, 0xa0, 0x26, 0x2c, 0x96, 0x19, 0x31, 0xff, 0x40, 0x27, 0x00, 0xc5, 0x88, 0x58, 0x94, 0x05, 0x4b, 0x13}

func nullPacketReader() io.Reader {
	return bytes.NewReader(nullPacket)
}

func dataPacketReader() io.Reader {
	return bytes.NewReader(dataPacket)
}

func TestPacketParsing(t *testing.T) {
	reader := ts.NewReader(nullPacketReader())

	p, err := reader.Next()
	if err != nil {
		t.Fatal(err)
	}
	if p.PID != nullPacketPID {
		t.Fatalf("Unexpected PID. Expected %x, got %x", nullPacketPID, p.PID)
	}
}

func TestIncompletePacket(t *testing.T) {
	reader := ts.NewReader(io.LimitReader(nullPacketReader(), 100))

	_, err := reader.Next()
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("Expected an error %v", err)
	}
}
