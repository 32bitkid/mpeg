package ts

import "github.com/32bitkid/bitreader"
import "io"

// AdaptationFieldControl is the two bit code that appears in a transport
// stream packet header that determines whether an Adapation Field appears
// in the bit stream.
type AdaptationFieldControl uint32

const (
	_ AdaptationFieldControl = iota
	PayloadOnly
	FieldOnly
	FieldThenPayload
)

// AdaptationField is an optional field in a transport stream packet header.
// TODO(jh): Needs implementation
type AdaptationField struct {
	Length uint32
	//DiscontinuityIndicator bool
	//RandomAccessIndicator bool
	//ElementaryStreamPriorityIndicator bool
	//PCRFlag bool
	//OPCRFlag bool
	//SplicingPointFlag bool
	//TransportPrivateDataFlag bool
	//AdaptationFieldExtensionFlag   bool
	Junk []byte
}

// ReadAdaptationField reads an AdaptationField from a bit stream.
func ReadAdaptationField(br bitreader.BitReader) (*AdaptationField, error) {
	var err error

	adaptationField := AdaptationField{}
	adaptationField.Length, err = br.Read32(8)
	if err != nil {
		return nil, err
	}

	adaptationField.Junk = make([]byte, adaptationField.Length)
	_, err = io.ReadFull(br, adaptationField.Junk)
	if err == io.EOF {
		return nil, io.ErrUnexpectedEOF
	} else if err != nil {
		return nil, err
	}

	return &adaptationField, nil
}
