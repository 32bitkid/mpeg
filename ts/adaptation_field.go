package ts

import br "github.com/32bitkid/bitreader"
import "io"

const (
	_ = iota
	PayloadOnly
	FieldOnly
	FieldThenPayload
)

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

func ReadAdaptationField(tsr br.Reader32) (*AdaptationField, error) {
	var err error

	adaptationField := AdaptationField{}
	adaptationField.Length, err = tsr.Read32(8)
	if err != nil {
		return nil, err
	}

	adaptationField.Junk = make([]byte, adaptationField.Length)
	_, err = io.ReadAtLeast(tsr, adaptationField.Junk, int(adaptationField.Length))

	if isFatalErr(err) {
		return nil, err
	}

	return &adaptationField, nil
}
