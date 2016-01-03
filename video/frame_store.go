package video

import "image"

type referenceFrame struct {
	frame              *image.YCbCr
	temporal_reference uint32
}

type frameStore struct {
	referenceFrames [2]referenceFrame
}
