package video

type PictureCodingType uint32

func (pct PictureCodingType) String() string {
	switch pct {
	case IntraCoded:
		return "I"
	case PredictiveCoded:
		return "P"
	case BidirectionallyPredictiveCoded:
		return "B"
	case DCIntraCoded:
		return "D"
	}
	return string(uint32(pct))
}

const (
	_                              PictureCodingType = iota // 000 forbidden
	IntraCoded                                              // 001
	PredictiveCoded                                         // 010
	BidirectionallyPredictiveCoded                          // 011
	DCIntraCoded                                            // 100 (Not Used in ISO/IEC11172-2)
	_                                                       // 101 reserved
	_                                                       // 110 reserved
	_                                                       // 111 reserved

	IFrame = IntraCoded
	PFrame = PredictiveCoded
	BFrame = BidirectionallyPredictiveCoded
)
