package pes

import "github.com/32bitkid/mpeg/util"

const (
	MinimumHeaderSize uint32 = 3
)

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

func ReadHeader(reader util.BitReader32) (*Header, uint32, error) {

	val, err := reader.Read32(2)
	if val != 2 || err != nil {
		return nil, 0, ErrMarkerNotFound
	}

	header := Header{}

	header.ScramblingControl, err = reader.Read32(2)
	if err != nil {
		return nil, 0, err
	}

	header.Priority, err = reader.ReadBit()
	if err != nil {
		return nil, 0, err
	}

	header.DataAlignmentIndicator, err = reader.ReadBit()
	if err != nil {
		return nil, 0, err
	}

	header.Copyright, err = reader.ReadBit()
	if err != nil {
		return nil, 0, err
	}

	header.OrignalOrCopy, err = reader.ReadBit()
	if err != nil {
		return nil, 0, err
	}

	header.PtsDtsFlags, err = reader.Read32(2)
	if err != nil {
		return nil, 0, err
	}

	header.EscrFlag, err = reader.ReadBit()
	if err != nil {
		return nil, 0, err
	}

	header.EsRateFlag, err = reader.ReadBit()
	if err != nil {
		return nil, 0, err
	}

	header.DsmTrickModeFlag, err = reader.ReadBit()
	if err != nil {
		return nil, 0, err
	}

	header.AdditionalCopyInfoFlag, err = reader.ReadBit()
	if err != nil {
		return nil, 0, err
	}

	header.CrcFlag, err = reader.ReadBit()
	if err != nil {
		return nil, 0, err
	}

	header.ExtensionFlag, err = reader.ReadBit()
	if err != nil {
		return nil, 0, err
	}

	header.HeaderDataLength, err = reader.Read32(8)
	if err != nil {
		return nil, 0, err
	}

	var bytesRead uint32 = MinimumHeaderSize

	if header.PtsDtsFlags == 2 {
		var len uint32
		header.PresentationTimeStamp, len, err = readTimeStamp(2, reader)

		if err != nil {
			return nil, 0, err
		}

		bytesRead += len

	}

	if header.PtsDtsFlags == 3 {
		var len uint32
		header.PresentationTimeStamp, len, err = readTimeStamp(3, reader)

		if err != nil {
			return nil, 0, err
		}

		bytesRead += len

		header.DecodingTimeStamp, len, err = readTimeStamp(1, reader)

		if err != nil {
			return nil, 0, err
		}

		bytesRead += len
	}

	if header.EscrFlag {
		panic("escr")
	}

	if header.EsRateFlag {
		panic("es rate")
	}

	if header.DsmTrickModeFlag {
		panic("dsm trick mode")
	}

	if header.AdditionalCopyInfoFlag {
		panic("additional copy info")
	}

	if header.CrcFlag {
		panic("crc flag")
	}

	if header.ExtensionFlag {
		var len uint32

		header.Extension, len, err = readExtension(reader)
		if err != nil {
			return nil, 0, err
		}

		bytesRead += len
	}

	stuffingLength, err := readStuffingBytes(reader, header.HeaderDataLength-(bytesRead-MinimumHeaderSize))
	if err != nil {
		return nil, 0, err
	}

	bytesRead += stuffingLength

	return &header, bytesRead, nil
}

func readStuffingBytes(reader util.BitReader32, length uint32) (uint32, error) {
	for i := uint32(0); i < length; i++ {
		val, err := reader.Read32(8)
		if err != nil {
			return 0, err
		}
		if val != 255 {
			return 0, ErrInvalidStuffingByte
		}
	}
	return length, nil
}

func readTimeStamp(marker uint32, reader util.BitReader32) (uint32, uint32, error) {

	var (
		timeStamp uint32
		err       error
		val       uint32
	)

	val, err = reader.Read32(4)
	if val != marker || err != nil {
		return 0, 0, ErrMarkerNotFound
	}

	val, err = reader.Read32(3)
	if err != nil {
		return 0, 0, err
	}

	timeStamp = timeStamp | (val << 30)

	val, err = reader.Read32(1)
	if val != 1 || err != nil {
		return 0, 0, ErrMarkerNotFound
	}

	val, err = reader.Read32(15)
	if err != nil {
		return 0, 0, err
	}

	timeStamp = timeStamp | (val << 15)

	val, err = reader.Read32(1)
	if val != 1 || err != nil {
		return 0, 0, ErrMarkerNotFound
	}

	val, err = reader.Read32(15)
	if err != nil {
		return 0, 0, err
	}

	timeStamp = timeStamp | val

	val, err = reader.Read32(1)
	if val != 1 || err != nil {
		return 0, 0, ErrMarkerNotFound
	}

	return timeStamp, 5, nil
}
