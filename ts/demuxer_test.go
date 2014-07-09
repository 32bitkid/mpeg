package ts_test

import "testing"
import "io"
import "github.com/32bitkid/mpeg-go/ts"

func TestDemuxingASinglePacket(t *testing.T) {

	source := io.MultiReader(nullPacketReader(), nullPacketReader())
	demux := ts.Demux(source)

	nullStream := demux.PID(nullPacketPID)

	eos := demux.Begin()

	var done = false
	for done == false {
		select {
		case p := <-nullStream:
			if p.PID != nullPacketPID {
				t.Fatalf("Unexpected PID. Expected %x, got %x", nullPacketPID, p.PID)
			}
			done = true
		case <-eos:
			done = true
		}
	}

	if demux.Err() != nil {
		t.Fatalf("Unxpected error: %s", demux.Err())
	}
}

func TestDemuxingASingleStream(t *testing.T) {
	source := io.MultiReader(nullPacketReader(), nullPacketReader(), nullPacketReader(), nullPacketReader())
	demux := ts.Demux(source)

	nullStream := demux.PID(nullPacketPID)
	eos := demux.Begin()

	var done = false
	count := 0
	for done == false {
		select {
		case <-nullStream:
			count++
		case <-eos:
			done = true
		}
	}

	if demux.Err() != io.ErrUnexpectedEOF {
		t.Fatalf("Unxpected error: %s", demux.Err())
	}

	if count != 4 {
		t.Fatalf("Not enough packets read. Expected %d, got %d", 4, count)
	}
}
