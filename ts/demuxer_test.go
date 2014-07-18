package ts_test

import "testing"
import "io"
import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/mpeg/ts"

func TestDemuxingASinglePacket(t *testing.T) {

	source := io.MultiReader(nullPacketReader(), nullPacketReader())
	demux := ts.NewDemuxer(bitreader.NewReader32(source))

	nullStream := demux.Where(ts.IsPID(nullPacketPID))

	stop := demux.Go()

	var done = false
	for done == false {
		select {
		case p := <-nullStream:
			if p.PID != nullPacketPID {
				t.Fatalf("Unexpected PID. Expected %x, got %x", nullPacketPID, p.PID)
			}
			done = true
		case <-stop:
			t.Fatalf("unexpected stop")
			done = true
		}
	}

	if demux.Err() != nil {
		t.Fatalf("Unxpected error: %s", demux.Err())
	}
}

func TestDemuxingASingleStream(t *testing.T) {
	source := io.MultiReader(nullPacketReader(), nullPacketReader(), nullPacketReader(), nullPacketReader())
	demux := ts.NewDemuxer(bitreader.NewReader32(source))

	nullStream := demux.Where(ts.IsPID(nullPacketPID))
	stop := demux.Go()

	var done = false
	count := 0
	for done == false {
		select {
		case _, ok := <-nullStream:
			if ok {
				count++
			}
		case <-stop:
			done = true
		}
	}

	if demux.Err() != bitreader.ErrNotAvailable {
		t.Fatalf("Unxpected error: %s", demux.Err())
	}

	if count != 4 {
		t.Fatalf("Not enough packets read. Expected %d, got %d", 4, count)
	}
}

func TestDemuxingUsingWheres(t *testing.T) {
	source := io.MultiReader(nullPacketReader(), dataPacketReader(), nullPacketReader(), dataPacketReader(), nullPacketReader())
	demux := ts.NewDemuxer(bitreader.NewReader32(source))
	dataStream := demux.Where(ts.IsPID(dataPacketPID))
	junkStream := demux.Where(ts.IsPID(dataPacketPID).Not())

	stop := demux.Go()

	var done = false
	dataCount := 0
	junkCount := 0
	for done == false {
		select {
		case _, ok := <-dataStream:
			if ok {
				dataCount++
			}
		case _, ok := <-junkStream:
			if ok {
				junkCount++
			}
		case <-stop:
			done = true
		}
	}

	if demux.Err() != bitreader.ErrNotAvailable {
		t.Fatalf("Unxpected error: %s", demux.Err())
	}

	if dataCount != 2 {
		t.Fatalf("Not enough packets read. Expected %d, got %d", 2, dataCount)
	}

	if junkCount != 3 {
		t.Fatalf("Not enough packets read. Expected %d, got %d", 3, junkCount)
	}
}

func TestDemuxingRange(t *testing.T) {
	source := io.MultiReader(fivePacketReader())
	demux := ts.NewDemuxer(bitreader.NewReader32(source))
	allStream := demux.Where(func(p *ts.Packet) bool { return true })

	count := 0

	demux.SkipUntil(ts.IsPID(0x31))
	demux.TakeWhile(ts.IsPID(0x41).Not())
	stop := demux.Go()

	var done = false
	for done == false {
		select {
		case _, ok := <-allStream:
			if ok {
				count++
			}
		case <-stop:
			done = true
		}
	}

	if demux.Err() != nil {
		t.Fatalf("Unxpected error: %s", demux.Err())
	}

	if count != 2 {
		t.Fatalf("Not enough packets read. Expected %d, got %d", 2, count)
	}
}
