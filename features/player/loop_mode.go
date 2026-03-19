package player

type LoopMode int

const (
	LoopOff LoopMode = iota
	LoopTrack
	LoopQueue
)

func (l LoopMode) String() string {
	switch l {
	case LoopTrack:
		return "リピート: トラック"
	case LoopQueue:
		return "リピート: キュー"
	default:
		return "リピートオフ"
	}
}

func (l LoopMode) Next() LoopMode {
	return (l + 1) % 3
}
