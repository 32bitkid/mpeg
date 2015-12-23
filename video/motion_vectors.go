package video

func motion_vectors(i int) {
	switch i {
	case 0:
		panic("forward motion vectors")
	case 1:
		panic("backwards motion vectors")
	default:
		panic("unknown motion vectors")
	}
}
