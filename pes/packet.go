package pes

import "errors"
import "io"
import "io/ioutil"
import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/mpeg/ts"

// StartCodePrefix is the prefix that signals the start a PES packet.
const StartCodePrefix = 0x000001

var (
	ErrStartCodePrefixNotFound = errors.New("start code prefix not found")
	ErrMarkerNotFound          = errors.New("marker not found")
	ErrInvalidStuffingByte     = errors.New("invalid stuffing byte")
)

// Packet is a parsed PES packet from a bitstream. A PES packet consists,
// at minimum, of a start_code_prefix, stream_id, packet_length, followed
// by a variable number of bytes of payload. It can optionally, for certain
// stream types, contain a Header.
//
//  ┌──────────────────────┬──────┬──────────────┬────────────────────────-
//  │start_code_prefix     │stream│packet_length │payload                 -
//  │                  (24)│id (8)│          (16)│                        -
//  └──────────────────────┴──────┴──────────────┴────────────────────────-
//                                               Λ
//                                              ╱ ╲
//       ╱─────────────(optional)──────────────╱   ╲
//      ╱                                           ╲
//      ┌───────────────────────────────────────────┐
//      │PES Header                                 │
//      │                                 (variable)│
//      └───────────────────────────────────────────┘
type Packet struct {
	StreamID     uint32
	PacketLength uint32
	*Header
	Payload []byte
}

// Creates a new packet and reads it from the bitstream.
func NewPacket(br bitreader.BitReader) (packet *Packet, err error) {
	packet = new(Packet)
	err = packet.Next(br)
	return
}

// Next reads the next packet from the bitstream
func (packet *Packet) Next(reader bitreader.BitReader) error {

	var (
		val uint32
		err error
	)

	val, err = reader.Peek32(24)

	if err != nil {
		return err
	}

	if val != StartCodePrefix {
		return ErrStartCodePrefixNotFound
	}

	reader.Trash(24)

	packet.StreamID, err = reader.Read32(8)
	if err != nil {
		return err
	}

	packet.PacketLength, err = reader.Read32(16)
	if err != nil {
		return err
	}

	switch {
	case hasPESHeader(packet.StreamID):
		var headerLen uint32
		packet.Header, headerLen, err = readHeader(reader)
		if err != nil {
			return err
		}

		if packet.PacketLength > 0 {
			var payloadLen = int(packet.PacketLength - headerLen)
			packet.Payload = make([]byte, payloadLen)
			_, err = io.ReadAtLeast(reader, packet.Payload, payloadLen)
			if err != nil {
				return err
			}
		} else {
			// Read until end of buffer
			// TODO dont create a new slice but read into a shared sized buffer
			packet.Payload, err = ioutil.ReadAll(reader)
			if err != nil && err != ts.EOP {
				return err
			}
		}
	case packet.StreamID == padding_stream:
		payloadLen := int(packet.PacketLength)
		junk := make([]byte, payloadLen)
		_, err = io.ReadAtLeast(reader, junk, payloadLen)
		if err != nil {
			return err
		}
	}

	return nil
}
