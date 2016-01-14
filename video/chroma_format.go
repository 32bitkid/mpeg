package video

// ChromaFormat is chroma subsample ratio used in the video bitstream. It
// defines the number and location of the blocks encoded in a macroblock.
//
// A macroblock in a 4:2:0 image contains 6 blocks: 4 Y, and 1 Cb, and 1 Cr,
// arranged in the following order:
//
//  ┌───┬───┐
//  │ 0 │ 1 │  ┌───┐  ┌───┐
//  ├───┼───┤  │ 4 │  │ 5 │
//  │ 2 │ 3 │  └───┘  └───┘
//  └───┴───┘
//      Y        Cb     Cr
//
// A macroblock in a 4:2:2 image contains 8 blocks: 4 Y, and 2 Cb, and 2 Cr,
// arranged in the following order:
//
//  ┌───┬───┐  ┌───┐  ┌───┐
//  │ 0 │ 1 │  │ 4 │  │ 5 │
//  ├───┼───┤  ├───┤  ├───┤
//  │ 2 │ 3 │  │ 6 │  │ 7 │
//  └───┴───┘  └───┘  └───┘
//      Y        Cb     Cr
//
// A macroblock in a 4:4:4 image contains 12 blocks: 4 Y, and 4 Cb, and 4 Cr,
// arranged in the following order:
//
//  ┌───┬───┐ ┌───┬───┐ ┌───┬───┐
//  │ 0 │ 1 │ │ 4 │ 8 │ │ 5 │ 9 │
//  ├───┼───┤ ├───┼───┤ ├───┼───┤
//  │ 2 │ 3 │ │ 6 │10 │ │ 7 │11 │
//  └───┴───┘ └───┴───┘ └───┴───┘
//      Y        Cb        Cr
//
// Note: At present, this library only supports decoding video with subsample
// ratio of 4:2:0.
type ChromaFormat uint32

const (
	_ ChromaFormat = iota // reserved
	ChromaFormat420
	ChromaFormat422
	ChromaFormat444
)

func (cf ChromaFormat) String() string {
	switch cf {
	case ChromaFormat420:
		return "4:2:0"
	case ChromaFormat422:
		return "4:2:2"
	case ChromaFormat444:
		return "4:4:4"
	}
	return "[Invalid]"
}
