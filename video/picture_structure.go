package video

type PictureStructure uint32

const (
	_ PictureStructure = iota
	PictureStructure_TopField
	PictureStructure_BottomField
	PictureStructure_FramePicture
)
