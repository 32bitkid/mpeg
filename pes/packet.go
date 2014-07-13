package pes

import "errors"
import "log"
import "github.com/32bitkid/bitreader"

var (
	ErrStartCodePrefixNotFound = errors.New("start code prefix not found")
	ErrMissingMarker           = errors.New("missing marker")
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
}

func ReadPacket(reader bitreader.Reader32) (*Packet, error) {

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
		packet.Header, err = ReadHeader(reader)
		if err != nil {
			return nil, err
		}
	}

	return &packet, nil
}

func ReadHeader(reader bitreader.Reader32) (*Header, error) {
	val, err := reader.Read32(2)
	if val != 2 || err != nil {
		return nil, ErrMissingMarker
	}

	header := Header{}

	header.ScramblingControl, err = reader.Read32(2)
	if err != nil {
		return nil, err
	}

	header.Priority, err = reader.ReadBit()
	if err != nil {
		return nil, err
	}

	header.DataAlignmentIndicator, err = reader.ReadBit()
	if err != nil {
		return nil, err
	}

	header.Copyright, err = reader.ReadBit()
	if err != nil {
		return nil, err
	}

	header.OrignalOrCopy, err = reader.ReadBit()
	if err != nil {
		return nil, err
	}

	header.PtsDtsFlags, err = reader.Read32(2)
	if err != nil {
		return nil, err
	}

	header.EscrFlag, err = reader.ReadBit()
	if err != nil {
		return nil, err
	}

	header.EsRateFlag, err = reader.ReadBit()
	if err != nil {
		return nil, err
	}

	header.DsmTrickModeFlag, err = reader.ReadBit()
	if err != nil {
		return nil, err
	}

	header.AdditionalCopyInfoFlag, err = reader.ReadBit()
	if err != nil {
		return nil, err
	}

	header.CrcFlag, err = reader.ReadBit()
	if err != nil {
		return nil, err
	}

	header.ExtensionFlag, err = reader.ReadBit()
	if err != nil {
		return nil, err
	}

	header.HeaderDataLength, err = reader.Read32(8)
	if err != nil {
		return nil, err
	}

	var headerBytesLeft uint32 = header.HeaderDataLength

	if header.PtsDtsFlags == 2 {
		header.PresentationTimeStamp, err = ReadTimeStamp(2, reader)

		if err != nil {
			return nil, err
		}

		headerBytesLeft -= 5
	}

	if header.PtsDtsFlags == 3 {
		header.PresentationTimeStamp, err = ReadTimeStamp(3, reader)

		if err != nil {
			return nil, err
		}

		headerBytesLeft -= 5

		header.DecodingTimeStamp, err = ReadTimeStamp(1, reader)

		if err != nil {
			return nil, err
		}

		headerBytesLeft -= 5
	}

	err = ReadStuffingBytes(reader, headerBytesLeft)
	if err != nil {
		return nil, err
	}

	val, _ = reader.Peek32(32)
	log.Printf(">>> %0#8x", val)

	return &header, nil
}

func ReadStuffingBytes(reader bitreader.Reader32, headerBytesLeft uint32) error {
	for headerBytesLeft > 0 {

		val, err := reader.Read32(8)
		headerBytesLeft--
		if err != nil {
			return err
		}
		if val != 255 {
			return ErrInvalidStuffingByte
		}
	}
	return nil
}

func ReadTimeStamp(marker uint32, reader bitreader.Reader32) (uint32, error) {

	var (
		timeStamp uint32
		err       error
		val       uint32
	)

	val, err = reader.Read32(4)
	if val != marker || err != nil {
		return 0, ErrMissingMarker
	}

	val, err = reader.Read32(3)
	if err != nil {
		return 0, err
	}

	timeStamp = timeStamp | (val << 30)

	val, err = reader.Read32(1)
	if val != 1 || err != nil {
		return 0, ErrMissingMarker
	}

	val, err = reader.Read32(15)
	if err != nil {
		return 0, err
	}

	timeStamp = timeStamp | (val << 15)

	val, err = reader.Read32(1)
	if val != 1 || err != nil {
		return 0, ErrMissingMarker
	}

	val, err = reader.Read32(15)
	if err != nil {
		return 0, err
	}

	timeStamp = timeStamp | val

	val, err = reader.Read32(1)
	if val != 1 || err != nil {
		return 0, ErrMissingMarker
	}

	return timeStamp, nil
}
