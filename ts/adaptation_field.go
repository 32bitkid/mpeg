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

type TransportPrivateData struct {
	// length of the transport private data section excludes the length byte
	length uint32
	data   []byte
}

func newTransportPrivateData(br bitreader.BitReader) (*TransportPrivateData, uint32, error) {
	tpd := TransportPrivateData{}
	var err error
	if tpd.length, err = br.Read32(8); err != nil {
		return nil, 0, err
	}
	if tpd.length > 0 {
		tpd.data = make([]byte, tpd.length)
		_, err = io.ReadFull(br, tpd.data)
		if err == io.EOF {
			return nil, 0, io.ErrUnexpectedEOF
		} else if err != nil {
			return nil, 0, err
		}
	}
	return &tpd, tpd.length, nil
}

type AdaptationExt struct {
	LegalTimeWindowFlag bool
	PiecewiseRateFlag   bool
	SeamlessSpliceFlag  bool

	// optional fields
	LTWValidFlag  bool
	LTWOffset     uint16
	PiecewiseRate uint32
	SpliceType    uint8
	DTSNextAU     uint64

	// length of the adaptation extension section excludes the length byte
	length uint32
}

func newAdaptationExt(br bitreader.BitReader) (*AdaptationExt, uint32, error) {
	ae := AdaptationExt{}
	var err error
	if ae.length, err = br.Read32(8); err != nil {
		return nil, 0, err
	}
	remainLength := ae.length

	if remainLength >= 1 {
		if ae.LegalTimeWindowFlag, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		if ae.PiecewiseRateFlag, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		if ae.SeamlessSpliceFlag, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		br.Skip(5)
		remainLength--
	}

	if remainLength >= 2 && ae.LegalTimeWindowFlag {
		if ae.LTWValidFlag, err = br.Read1(); err != nil {
			return nil, 0, err
		}
		if ae.LTWOffset, err = br.Read16(15); err != nil {
			return nil, 0, err
		}
		remainLength -= 2
	}

	if remainLength >= 3 && ae.PiecewiseRateFlag {
		br.Skip(2)
		if ae.PiecewiseRate, err = br.Read32(22); err != nil {
			return nil, 0, err
		}
		remainLength -= 3
	}

	if remainLength >= 5 && ae.SeamlessSpliceFlag {
		if ae.SpliceType, err = br.Read8(4); err != nil {
			return nil, 0, err
		}
		var v uint64
		if v, err = br.Read64(3); err != nil {
			return nil, 0, err
		}
		ae.DTSNextAU = v << 30
		br.Skip(1)
		if v, err = br.Read64(15); err != nil {
			return nil, 0, err
		}
		ae.DTSNextAU += v << 15
		br.Skip(1)
		if v, err = br.Read64(15); err != nil {
			return nil, 0, err
		}
		ae.DTSNextAU += v
	}

	return &ae, ae.length, nil
}

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

	// optional fields
	PCR                  uint64
	OPCR                 uint64
	SpliceCountdown      uint8
	TransportPrivateData *TransportPrivateData
	AdapatationExtension *AdaptationExt

	// length of the adaptation field section excludes the length byte
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

	if remainLength >= 1 && af.TransportPrivateDataFlag {
		var length uint32
		if af.TransportPrivateData, length, err = newTransportPrivateData(br); err != nil {
			return nil, 0, err
		}
		remainLength -= length + 1
	}

	if remainLength >= 2 && af.AdaptationFieldExtensionFlag {
		var length uint32
		if af.AdapatationExtension, length, err = newAdaptationExt(br); err != nil {
			return nil, 0, err
		}
		remainLength -= length + 1
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
