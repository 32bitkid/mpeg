package video

type PictureData struct {
	slices []*Slice
}

func (self *VideoSequence) picture_data() (err error) {

	pd := PictureData{}

	for {
		s, err := self.slice()
		if err != nil {
			return err
		}

		pd.slices = append(pd.slices, s)

		nextbits, err := self.Peek32(32)
		if err != nil {
			return err
		}

		if !is_slice_start_code(StartCode(nextbits)) {
			break
		}
	}

	self.PictureData = &pd

	return self.next_start_code()
}
