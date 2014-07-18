package pes

import "errors"
import "io"
import . "github.com/32bitkid/mpeg_go"

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

type Header struct {
	ScramblingControl      uint32
	Priority               bool
	DataAlignmentIndicator bool
	Copyright              bool
	OrignalOrCopy          bool
	PtsDtsFlags            uint32
	EscrFlag               bool
	EsRateFlag             bool
	DsmTrickModeFlag       bool
	AdditionalCopyInfoFlag bool
	CrcFlag                bool
	ExtensionFlag          bool
	HeaderDataLength       uint32

	PresentationTimeStamp uint32
	DecodingTimeStamp     uint32

	Extension *Extension
}

func ReadPacket(reader BitReader, total int) (*Packet, error) {

	var (
		val uint32
		err error
	)

	val, err = reader.Peek32(24)

	if val != 0x0001 || err != nil {
		return nil, ErrStartCodePrefixNotFound
	}

	reader.Trash(24)

	packet := Packet{}

	packet.StreamID, err = reader.Read32(8)
	if err != nil {
		return nil, err
	}

	packet.PacketLength, err = reader.Read32(16)
	if err != nil {
		return nil, err
	}

	if hasPESHeader(packet.StreamID) {

		var len uint32
		packet.Header, len, err = ReadHeader(reader)
		if err != nil {
			return nil, err
		}

		var payloadLen int

		if total > 0 {
			payloadLen = total - int(packet.Header.HeaderDataLength) - 3 - 6
		} else {
			payloadLen = int(packet.PacketLength - len)
		}

		packet.Payload = make([]byte, payloadLen)

		_, err = io.ReadAtLeast(reader, packet.Payload, payloadLen)
		if err != nil {
			return nil, err
		}

	} else if packet.StreamID == padding_stream {
		payloadLen := int(packet.PacketLength)
		junk := make([]byte, payloadLen)
		_, err = io.ReadAtLeast(reader, junk, payloadLen)
		if err != nil {
			return nil, err
		}
	}

	return &packet, nil
}
