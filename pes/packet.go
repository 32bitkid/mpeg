package pes

import "errors"
import "io"
import "io/ioutil"
import "github.com/32bitkid/mpeg/util"
import "github.com/32bitkid/mpeg/ts"

const (
	StartCodePrefix = 0x000001
)

var (
	ErrStartCodePrefixNotFound = errors.New("start code prefix not found")
	ErrMarkerNotFound          = errors.New("marker not found")
	ErrInvalidStuffingByte     = errors.New("invalid stuffing byte")
)

type Packet struct {
	StreamID     uint32
	PacketLength uint32
	*Header
	Payload []byte
}

func NewPacket(br util.BitReader32) (packet *Packet, err error) {
	packet = new(Packet)
	err = packet.readFrom(br)
	return packet, err
}

func (packet *Packet) readFrom(reader util.BitReader32) error {

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
		packet.Header, headerLen, err = ReadHeader(reader)
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
