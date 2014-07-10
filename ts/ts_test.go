package ts_test

import "testing"
import "io"
import "github.com/32bitkid/mpeg-go/ts"

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
