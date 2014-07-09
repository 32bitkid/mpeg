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

var ErrNoSyncByte = errors.New("no sync byte")

type TsPacket struct {
	TransportErrorIndicator    bool
	PayloadUnitStartIndicator  bool
	TransportPrriority         bool
	PID                        uint32
	TransportScramblingControl uint32
	AdaptationFieldControl     uint32
	ContinuityCounter          uint32
	Payload                    []byte
}

func isFatalErr(err error) bool {
	return err != nil && err != io.EOF
}

func (tsr *tsReader) Next() (*TsPacket, error) {

	var err error

	aligned, err := tsr.isAligned()
	if isFatalErr(err) {
		return nil, err
	}

	if !aligned {
		err = tsr.realign()
		if err != nil {
			return nil, ErrNoSyncByte
		}
	}

	err = tsr.Trash(8)
	if isFatalErr(err) {
		return nil, err
	}

	packet := TsPacket{}

	packet.TransportErrorIndicator, err = tsr.ReadBit()
	if isFatalErr(err) {
		return nil, err
	}

	packet.PayloadUnitStartIndicator, err = tsr.ReadBit()
	if isFatalErr(err) {
		return nil, err
	}

	packet.TransportPrriority, err = tsr.ReadBit()
	if isFatalErr(err) {
		return nil, err
	}

	packet.PID, err = tsr.Read32(13)
	if isFatalErr(err) {
		return nil, err
	}

	packet.TransportScramblingControl, err = tsr.Read32(2)
	if isFatalErr(err) {
		return nil, err
	}

	packet.AdaptationFieldControl, err = tsr.Read32(2)
	if isFatalErr(err) {
		return nil, err
	}

	packet.ContinuityCounter, err = tsr.Read32(4)
	if isFatalErr(err) {
		return nil, err
	}

	// TODO handle adaptation field

	packet.Payload = make([]byte, 184)

	var val uint32
	for i := 0; i < 184; i++ {
		val, err = tsr.Read32(8)
		if isFatalErr(err) {
			return nil, err
		}
		packet.Payload[i] = byte(val)
	}

	return &packet, nil
}

func (tsr *tsReader) isAligned() (bool, error) {
	val, err := tsr.Peek32(8)
	return val == SyncByte, err
}

func (tsr *tsReader) realign() error {
	log.Printf("Attempting to realign")
	for i := 0; i < 188; i++ {
		if err := tsr.Trash(8); isFatalErr(err) {
			return err
		}
		isAligned, err := tsr.isAligned()
		if isFatalErr(err) {
			return err
		}
		if isAligned {
			log.Printf("Realigned after %d bytes.\n", i)
			return nil
		}
	}
	return ErrNoSyncByte
}
