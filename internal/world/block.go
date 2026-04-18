package world

type BlockType uint8

const (
	BlockAir BlockType = iota
	BlockSolid
	BlockDebugMover
)

func (b BlockType) Valid() bool {
	switch b {
	case BlockAir, BlockSolid, BlockDebugMover:
		return true
	default:
		return false
	}
}

func (b BlockType) String() string {
	switch b {
	case BlockAir:
		return "air"
	case BlockSolid:
		return "solid"
	case BlockDebugMover:
		return "debug_mover"
	default:
		return "invalid"
	}
}
