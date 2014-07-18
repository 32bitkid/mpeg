package pes

import . "github.com/32bitkid/mpeg_go"
import "io"

type Extension struct {
	PrivateDataFlag                  bool
	PackHeaderFieldFlag              bool
	ProgramPacketSequenceCounterFlag bool
	P_STD_BufferFlag                 bool
	ExtensionFlag2                   bool

	PrivateData []byte
	*P_STD_Buffer
}

type P_STD_Buffer struct {
	Scale uint32
	Size  uint32
}

func readExtension(r BitReader) (*Extension, uint32, error) {
	var err error
	var bytesRead uint32 = 0
	extension := Extension{}

	if extension.PrivateDataFlag, err = r.ReadBit(); err != nil {
		return nil, 0, err
	}
	if extension.PackHeaderFieldFlag, err = r.ReadBit(); err != nil {
		return nil, 0, err
	}
	if extension.ProgramPacketSequenceCounterFlag, err = r.ReadBit(); err != nil {
		return nil, 0, err
	}
	if extension.P_STD_BufferFlag, err = r.ReadBit(); err != nil {
		return nil, 0, err
	}

	if err = r.Trash(3); err != nil {
		return nil, 0, err
	}

	if extension.ExtensionFlag2, err = r.ReadBit(); err != nil {
		return nil, 0, err
	}

	bytesRead += 1

	if extension.PrivateDataFlag {
		extension.PrivateData = make([]byte, 16)
		_, err := io.ReadAtLeast(r, extension.PrivateData, 16)
		bytesRead += 16
		if err != nil {
			return nil, 0, err
		}
	}

	if extension.PackHeaderFieldFlag {
		panic("pack header field")
	}

	if extension.ProgramPacketSequenceCounterFlag {
		panic("program packet sequence counter")
	}

	if extension.P_STD_BufferFlag {
		if v, err := r.Read32(2); v != 1 || err != nil {
			return nil, 0, ErrMarkerNotFound
		}

		scale, err := r.Read32(1)
		if err != nil {
			return nil, 0, err
		}
		size, err := r.Read32(13)
		if err != nil {
			return nil, 0, err
		}
		extension.P_STD_Buffer = &P_STD_Buffer{scale, size}
		bytesRead += 2
	}

	if extension.ExtensionFlag2 {
		panic("extension 2")
	}

	return &extension, bytesRead, nil
}
