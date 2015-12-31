package ts

import "github.com/32bitkid/bitreader"
import "io"

// AdaptationFieldControl is the two bit code that appears in a transport
// stream packet header that determines whether an Adapation Field appears
// in the bit stream.
type AdaptationFieldControl uint32

const (
	_                AdaptationFieldControl = iota
	PayloadOnly                             // 0b01
	FieldOnly                               // 0b10
	FieldThenPayload                        //0b11
)

// AdaptationField is an optional field in a transport stream packet header.
// TODO(jh): Needs implementation
type AdaptationField struct {
	//DiscontinuityIndicator bool
	//RandomAccessIndicator bool
	//ElementaryStreamPriorityIndicator bool
	//PCRFlag bool
	//OPCRFlag bool
	//SplicingPointFlag bool
	//TransportPrivateDataFlag bool
	//AdaptationFieldExtensionFlag   bool

	length uint32
	junk   []byte
}

func newAdaptationField(br bitreader.BitReader) (*AdaptationField, uint32, error) {
	adaptationField := AdaptationField{}
	length, err := br.Read32(8)
	if err != nil {
		return nil, 0, err
	}

	adaptationField.junk = make([]byte, length)
	_, err = io.ReadFull(br, adaptationField.junk)
	if err == io.EOF {
		return nil, 0, io.ErrUnexpectedEOF
	} else if err != nil {
		return nil, 0, err
	}

	return &adaptationField, length, nil
}
