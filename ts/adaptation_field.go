package ts

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

func ReadAdaptationField(tsr *tsReader) (*AdaptationField, error) {
	var err error

	adaptationField := AdaptationField{}
	adaptationField.Length, err = tsr.Read32(8)
	if err != nil {
		return nil, err
	}

	adaptationField.Junk = make([]byte, adaptationField.Length)

	var val uint32
	var i uint32
	for i = 0; i < adaptationField.Length; i++ {
		val, err = tsr.Read32(8)
		if isFatalErr(err) {
			return nil, err
		}
		adaptationField.Junk[i] = byte(val)
	}

	return &adaptationField, nil
}
