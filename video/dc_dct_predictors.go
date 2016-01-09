package video

type dcDctPredictors [3]int32
type dcDctPredictorResetter func()

func (pred *dcDctPredictors) createResetter(intra_dc_precision uint32) dcDctPredictorResetter {
	return func() {
		resetValue := int32(1) << (7 + intra_dc_precision)
		pred[0] = resetValue
		pred[1] = resetValue
		pred[2] = resetValue
	}
}
