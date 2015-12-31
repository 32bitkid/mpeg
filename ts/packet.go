package ts

import "github.com/32bitkid/bitreader"
import "io"
import "errors"
import "fmt"

// SyncByte the the fixed 8-bit value that marks the start of a TS packet.
const SyncByte = 0x47

// MaxPayloadSize is the maximum payload, in bytes, that a
// Transport Stream packet can contain.
const MaxPayloadSize = 184

// ErrNoSyncByte is the error returned if a sync byte cannot be located in the bitstream.
var ErrNoSyncByte = errors.New("no sync byte")

// Packet is a parsed Transport Stream packet from the bit stream.
//
//                   ┌────────────────┬──────────────────────────────────────────────────────────────────────┐
//                   │ header         │ payload                                                              │
//                   │                │                                                           (184 bytes)│
//                   └────────────────┴──────────────────────────────────────────────────────────────────────┘
//                   │                ╲
//                   │                 ╲────────────────────────────────────────────────────────────────────╲
//                   │                                                                                       ╲
//                   ┌──────────────────────┬─┬─┬─┬─────────────────────────────────────┬────┬────┬──────────┐
//                   │ sync_byte (8)        │ │ │ │ PID (13)                            │    │    │continuity│
//                   │                      │ │ │ │                                     │    │    │counter(4)│
//                   └──────────────────────┴─┴─┴─┴─────────────────────────────────────┴────┴────┴──────────┘
//                                          ╱     ╲                                     ╱         ╲
//                                         ╱       ╲                                   ╱           ╲
//   ╱────────────────────────────────────╱         ╲───────╲    ╱────────────────────╱             ╲────────────────────────╲
//  ╱                                                        ╲  ╱                                                             ╲
//  ┌──────────────────┬──────────────────┬──────────────────┐  ┌──────────────────────────────┬──────────────────────────────┐
//  │ transport_error  │payload_unit_start│transport_priority│  │ transport_scrambling_control │  adaptation_field_control    │
//  │               (1)│               (1)│               (1)│  │                           (2)│                           (2)│
//  └──────────────────┴──────────────────┴──────────────────┘  └──────────────────────────────┴──────────────────────────────┘
type Packet struct {
	TransportErrorIndicator    bool
	PayloadUnitStartIndicator  bool
	TransportPriority          bool
	PID                        uint32
	TransportScramblingControl uint32
	AdaptationFieldControl     AdaptationFieldControl
	ContinuityCounter          uint32
	AdaptationField            *AdaptationField
	Payload                    []byte

	payloadBuffer [MaxPayloadSize]byte
}

// Create a new Packet and read it from the bit stream.
func NewPacket(br bitreader.BitReader) (packet *Packet, err error) {
	packet = new(Packet)
	err = packet.Next(br)
	return
}

//  Read the next packet from the bit stream into an existing Packet.
func (packet *Packet) Next(br bitreader.BitReader) (err error) {

	aligned, err := isAligned(br)
	if err != nil {
		return
	}

	if !aligned {
		if err = realign(br); err != nil {
			return
		}
	}

	if err = br.Trash(8); err != nil {
		return
	}

	packet.TransportErrorIndicator, err = br.ReadBit()
	if err != nil {
		return
	}

	packet.PayloadUnitStartIndicator, err = br.ReadBit()
	if err != nil {
		return
	}

	packet.TransportPriority, err = br.ReadBit()
	if err != nil {
		return
	}

	packet.PID, err = br.Read32(13)
	if err != nil {
		return
	}

	packet.TransportScramblingControl, err = br.Read32(2)
	if err != nil {
		return
	}

	afc, err := br.Read32(2)
	if err != nil {
		return
	}

	packet.AdaptationFieldControl = AdaptationFieldControl(afc)

	packet.ContinuityCounter, err = br.Read32(4)
	if err != nil {
		return
	}

	var payloadSize uint32 = MaxPayloadSize

	if packet.AdaptationFieldControl == FieldOnly || packet.AdaptationFieldControl == FieldThenPayload {
		var length uint32
		packet.AdaptationField, length, err = newAdaptationField(br)

		if err != nil {
			return
		}
		payloadSize -= length + 1
	}

	if packet.AdaptationFieldControl == PayloadOnly || packet.AdaptationFieldControl == FieldThenPayload {
		packet.Payload = packet.payloadBuffer[0:payloadSize]

		_, err = io.ReadFull(br, packet.Payload)
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		} else if err != nil {
			return
		}
	}

	return nil
}

func isAligned(br bitreader.BitReader) (bool, error) {
	if br.IsByteAligned() == false {
		return false, nil
	}

	val, err := br.Peek32(8)
	if err != nil {
		return false, err
	}

	return val == SyncByte, nil
}

func realign(br bitreader.BitReader) error {
	if br.IsByteAligned() == false {
		br.ByteAlign()
	}

	for i := 0; i < 188; i++ {
		val, err := br.Peek32(8)
		if err != nil {
			return err
		}
		if val == SyncByte {
			return nil
		}
		if err := br.Trash(8); err != nil {
			return err
		}
	}
	return ErrNoSyncByte
}

func (p *Packet) String() string {
	return fmt.Sprintf("{ PID: 0x%02x, Counter: %1x }", p.PID, p.ContinuityCounter)
}
