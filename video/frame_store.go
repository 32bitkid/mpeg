package video

import "image"

type referenceFrame struct {
	frame              *image.YCbCr
	temporal_reference uint32
	previous_gop       bool
}

type frameStore struct {
	frames [2]referenceFrame
	past   *image.YCbCr
	future *image.YCbCr
}

func (fs *frameStore) gop() {
	fs.frames[0].previous_gop = true
	fs.frames[1].previous_gop = true
}

func (fs *frameStore) set(current_temporal_reference uint32) {
	for _, ref := range fs.frames {
		if ref.frame != nil {
			if ref.previous_gop == true ||
				ref.temporal_reference < current_temporal_reference {
				fs.past = ref.frame
			} else {
				fs.future = ref.frame
			}
		}
	}
}

func (fs *frameStore) add(f *image.YCbCr, temporal_reference uint32) {
	fs.frames[0] = fs.frames[1]
	fs.frames[1] = referenceFrame{f, temporal_reference, false}
}

func (fs *frameStore) tryGet(temporal_reference uint32) *image.YCbCr {
	for _, ref := range fs.frames {
		if ref.frame != nil && ref.previous_gop == false && ref.temporal_reference == temporal_reference {
			return ref.frame
		}
	}
	return nil
}
