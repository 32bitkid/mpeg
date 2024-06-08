package video

import "fmt"

// PictureCodingType defines the encoding used for a frame within a video
// bitstream.
//
// Intra coded pictures are encoded independently of any other pictures. They
// are similar to other self-contained image formats like JPEG.
//
// Predictive coded pictures are pictures coded using a motion compensated
// sample from a previous decoded reference frame to be reconstructed. If
// a macroblock is encoded with motion vectors, then to decode the block
// a motion prediction is formed by projecting a portion of a "past" reference
// frame "forward" into the current frame. Using the motion vector to sample image
// data from a previously decoded reference frame, the samples are then added
// to the encoded samples in the bitstream to reconstruct the final block.
//
//	x──────────────────────┐   x───┬──────────────────┐
//	│╲               ┌ ─ ─ ┼ ─ ▶   │                  │
//	│ MV                   │   ├───┘                  │
//	│  ╲───┐         │     │   │                      │
//	│  │   │─ ─ ─ ─ ─      │   │                      │
//	│  └───┘               │   │                      │
//	│                      │   │                      │
//	│                      │   │                      │
//	│                      │   │                      │
//	└──────────────────────┘   └──────────────────────┘
//	    Reference Frame                P-Frame
//
// Bi-directionally coded picture is a picture coded using motion compensated
// sample from either a past frame and/or a future frame. B-Frames can contain
// "forward" and/or "backward" projections.
//
//	x──────────────────────┐   x───┬──────────────────┐   x──────────────────────┐
//	│╲               ┌ ─ ─ ┼ ─ ▶   ◀ ─ ─ ─ ─ ─ ─ ─    │   │╲─MV──┌───┐           │
//	│ MV                   │   ├───┘              └ ─ ┼ ─ ┼ ─ ─ ─│   │           │
//	│  ╲───┐         │     │   │                      │   │      └───┘           │
//	│  │   │─ ─ ─ ─ ─      │   │                      │   │                      │
//	│  └───┘               │   │                      │   │                      │
//	│                      │   │                      │   │                      │
//	│                      │   │                      │   │                      │
//	│                      │   │                      │   │                      │
//	└──────────────────────┘   └──────────────────────┘   └──────────────────────┘
//	  Past Reference Frame             B-Frame             Future Reference Frame
type PictureCodingType uint32

const (
	_ PictureCodingType = iota // 0b000 forbidden

	PictureCodingType_IntraCoded
	PictureCodingType_PredictiveCoded
	PictureCodingType_BidirectionallyPredictiveCoded

	PictureCodingType_DCIntraCoded // Shall not be used (Used in ISO/IEC11172-2)
	_                              // 0b101 reserved
	_                              // 0b110 reserved
	_                              // 0b111 reserved

	IFrame = PictureCodingType_IntraCoded
	PFrame = PictureCodingType_PredictiveCoded
	BFrame = PictureCodingType_BidirectionallyPredictiveCoded
)

func (pct PictureCodingType) String() string {
	switch pct {
	case PictureCodingType_IntraCoded:
		return "I"
	case PictureCodingType_PredictiveCoded:
		return "P"
	case PictureCodingType_BidirectionallyPredictiveCoded:
		return "B"
	case PictureCodingType_DCIntraCoded:
		return "D"
	}
	return fmt.Sprintf("%v?", uint32(pct))
}
