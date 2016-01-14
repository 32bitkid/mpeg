package video

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
