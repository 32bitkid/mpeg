package ts

import "github.com/32bitkid/bitreader"
import "io"
import "errors"

const (
	SyncByte       = 0x47
	MaxPayloadSize = 184
)

var (
	ErrNoSyncByte = errors.New("no sync byte")
)

type Packet struct {
	TransportErrorIndicator    bool
	PayloadUnitStartIndicator  bool
	TransportPriority          bool
	PID                        uint32
	TransportScramblingControl uint32
	AdaptationFieldControl     uint32
	ContinuityCounter          uint32
	AdaptationField            *AdaptationField
	Payload                    []byte

	payloadBuffer [MaxPayloadSize]byte
}

func NewPacket(br bitreader.BitReader) (packet *Packet, err error) {
	packet = new(Packet)
	err = packet.ReadFrom(br)
	return
}

func (packet *Packet) ReadFrom(tsr bitreader.BitReader) (err error) {

	aligned, err := isAligned(tsr)
	if err != nil {
		return
	}

	if !aligned {
		if err = realign(tsr); err != nil {
			return
		}
	}

	if err = tsr.Trash(8); err != nil {
		return
	}

	packet.TransportErrorIndicator, err = tsr.ReadBit()
	if err != nil {
		return
	}

	packet.PayloadUnitStartIndicator, err = tsr.ReadBit()
	if err != nil {
		return
	}

	packet.TransportPriority, err = tsr.ReadBit()
	if err != nil {
		return
	}

	packet.PID, err = tsr.Read32(13)
	if err != nil {
		return
	}

	packet.TransportScramblingControl, err = tsr.Read32(2)
	if err != nil {
		return
	}

	packet.AdaptationFieldControl, err = tsr.Read32(2)
	if err != nil {
		return
	}

	packet.ContinuityCounter, err = tsr.Read32(4)
	if err != nil {
		return
	}

	var payloadSize uint32 = MaxPayloadSize

	if packet.AdaptationFieldControl == FieldOnly || packet.AdaptationFieldControl == FieldThenPayload {
		packet.AdaptationField, err = ReadAdaptationField(tsr)

		if err != nil {
			return
		}
		payloadSize -= packet.AdaptationField.Length + 1
	}

	if packet.AdaptationFieldControl == PayloadOnly || packet.AdaptationFieldControl == FieldThenPayload {
		packet.Payload = packet.payloadBuffer[0:payloadSize]

		_, err = io.ReadFull(tsr, packet.Payload)
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		} else if err != nil {
			return
		}
	}

	return nil
}

func isAligned(tsr bitreader.BitReader) (bool, error) {
	if tsr.IsByteAligned() == false {
		return false, nil
	}

	val, err := tsr.Peek32(8)
	if err != nil {
		return false, err
	}

	return val == SyncByte, nil
}

func realign(tsr bitreader.BitReader) error {
	if tsr.IsByteAligned() == false {
		tsr.ByteAlign()
	}

	for i := 0; i < 188; i++ {
		val, err := tsr.Peek32(8)
		if err != nil {
			return err
		}
		if val == SyncByte {
			return nil
		}
		if err := tsr.Trash(8); err != nil {
			return err
		}
	}
	return ErrNoSyncByte
}
