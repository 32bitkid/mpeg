package video

import "image"

type referenceFrame struct {
	frame              *image.YCbCr
	temporal_reference uint32
}

type frameStore struct {
	frames [2]referenceFrame
	past   *image.YCbCr
	future *image.YCbCr
}

func (fs *frameStore) set(current_temporal_reference uint32) {
	for _, ref := range fs.frames {
		if ref.frame != nil {
			if ref.temporal_reference < current_temporal_reference {
				fs.past = ref.frame
			} else {
				fs.future = ref.frame
			}
		}
	}
}

func (fs *frameStore) add(f *image.YCbCr, temporal_reference uint32) {
	fs.frames[0] = fs.frames[1]
	fs.frames[1] = referenceFrame{f, temporal_reference}
}
