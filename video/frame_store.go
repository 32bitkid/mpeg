package video

import "image"

type referenceFrame struct {
	frame              *image.YCbCr
	temporal_reference uint32
}

type frameStore struct {
	frames   [2]referenceFrame
	forward  *image.YCbCr
	backward *image.YCbCr
}

func (fs *frameStore) set(current_temporal_reference uint32) {
	for _, ref := range fs.frames {
		if ref.frame != nil {
			if ref.temporal_reference < current_temporal_reference {
				fs.forward = ref.frame
			} else {
				fs.backward = ref.frame
			}
		}
	}
}

func (fs *frameStore) add(f *image.YCbCr, temporal_reference uint32) {
	fs.frames[1] = fs.frames[0]
	fs.frames[0] = referenceFrame{f, temporal_reference}
}
