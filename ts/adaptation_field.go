package ts

import (
	"io"

	"github.com/32bitkid/bitreader"
)

// AdaptationFieldControl is the two bit code that appears in a transport
// stream packet header that determines whether an Adapation Field appears
// in the bit stream.
type AdaptationFieldControl uint32

const (
	_                AdaptationFieldControl = iota
	PayloadOnly                             // 0b01
	FieldOnly                               // 0b10
	FieldThenPayload                        // 0b11
)

// AdaptationField is an optional field in a transport stream packet header.
type AdaptationField struct {
	DiscontinuityIndicator            bool
	RandomAccessIndicator             bool
	ElementaryStreamPriorityIndicator bool
	PCRFlag                           bool
	OPCRFlag                          bool
	SplicingPointFlag                 bool
	TransportPrivateDataFlag          bool
	AdaptationFieldExtensionFlag      bool

	PCR             uint64
	OPCR            uint64
	SpliceCountdown uint8
	// TODO: implement these
	// TransportPrivateData *TransportPrivateData
	// AdapatationExtension *AdaptationExt

	length uint32
	junk   []byte
}

func newAdaptationField(br bitreader.BitReader) (*AdaptationField, uint32, error) {
	af := AdaptationField{}
	remainLength, err := br.Read32(8)
	af.length = remainLength
	if err != nil {
		return nil, 0, err
	}

	// parse the adaptation field flags
	if remainLength >= 1 {
		if af.DiscontinuityIndicator, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		if af.RandomAccessIndicator, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		if af.ElementaryStreamPriorityIndicator, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		if af.PCRFlag, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		if af.OPCRFlag, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		if af.SplicingPointFlag, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		if af.TransportPrivateDataFlag, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		if af.AdaptationFieldExtensionFlag, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		remainLength--
	}

	readPcr := func(br bitreader.BitReader) (uint64, error) {
		base, err1 := br.Read64(33)
		if err1 != nil {
			return 0, err1
		}
		if _, err2 := br.Read64(6); err2 != nil {
			return 0, err2
		}
		extension, err3 := br.Read64(9)
		if err3 != nil {
			return 0, err3
		}
		return base*300 + extension, nil
	}

	if remainLength >= 6 && af.PCRFlag {
		if af.PCR, err = readPcr(br); err != nil {
			return nil, 0, err
		}
		remainLength -= 6
	}

	if remainLength >= 6 && af.OPCRFlag {
		if af.OPCR, err = readPcr(br); err != nil {
			return nil, 0, err
		}
		remainLength -= 6
	}

	if remainLength >= 1 && af.SplicingPointFlag {
		if af.SpliceCountdown, err = br.Read8(8); err != nil {
			return nil, 0, err
		}
		remainLength--
	}

	// skip the remaining bytes
	if remainLength > 0 {
		af.junk = make([]byte, remainLength)
		_, err := io.ReadFull(br, af.junk)
		if err == io.EOF {
			return nil, 0, io.ErrUnexpectedEOF
		} else if err != nil {
			return nil, 0, err
		}
	}

	return &af, af.length, nil
}
