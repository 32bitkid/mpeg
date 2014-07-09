package ts

import "io"
import "log"
import "errors"
import "github.com/32bitkid/bitreader"

func NewReader(reader io.Reader) TransportStreamReader {
	br := bitreader.NewReader32(reader)
	return &tsReader{br}
}

type TransportStreamReader interface {
	Next() (*TsPacket, error)
}

type tsReader struct {
	bitreader.Reader32
}

const SyncByte = 0x47

type TsPacket struct {
	TransportErrorIndicator    bool
	PayloadUnitStartIndicator  bool
	TransportPriority          bool
	PID                        uint32
	TransportScramblingControl uint32
	AdaptationFieldControl     uint32
	ContinuityCounter          uint32
	Payload                    []byte
}

func (tsr *tsReader) Next() (*TsPacket, error) {

	if !tsr.isAligned() && !tsr.realign() {
		return nil, errors.New("No sync_byte found")
	}

	tsr.Trash(8)

	packet := TsPacket{
		TransportErrorIndicator:   tsr.ReadBit(),
		PayloadUnitStartIndicator: tsr.ReadBit(),
		TransportPriority:         tsr.ReadBit(),
		PID:                       tsr.Read32(13),
		TransportScramblingControl: tsr.Read32(2),
		AdaptationFieldControl:     tsr.Read32(2),
		ContinuityCounter:          tsr.Read32(4),
		Payload:                    make([]byte, 184),
	}

	for i := 0; i < 184; i++ {
		packet.Payload[i] = byte(tsr.Read32(8))
	}

	return &packet, nil
}

func (tsr *tsReader) isAligned() bool {
	return tsr.Peek32(8) == SyncByte
}

func (tsr *tsReader) realign() bool {
	log.Printf("Attempting to realign")
	for i := 0; i < 188; i++ {
		if tsr.isAligned() {
			log.Printf("Realigned after %d bytes.\n", i)
			return true
		}
		tsr.Trash(8)
	}
	return false
}
