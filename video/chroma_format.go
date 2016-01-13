package video

type chromaFormat uint32

const (
	_ chromaFormat = iota // reserved
	ChromaFormat_420
	ChromaFormat_422
	ChromaFormat_444
)

func (cf chromaFormat) String() string {
	switch cf {
	case ChromaFormat_420:
		return "4:2:0"
	case ChromaFormat_422:
		return "4:2:2"
	case ChromaFormat_444:
		return "4:4:4"
	}
	return "[Invalid]"
}
