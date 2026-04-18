package world

import "errors"

var (
	ErrOutOfBounds      = errors.New("position out of bounds")
	ErrInvalidBlockType = errors.New("invalid block type")
)

type Position struct {
	X int
	Y int
	Z int
}

type Dimensions struct {
	X int
	Y int
	Z int
}

type World struct {
	dimensions Dimensions
	blocks     map[Position]BlockType
}

func NewDefault() *World {
	return New(Dimensions{X: 32, Y: 32, Z: 16})
}

func New(dimensions Dimensions) *World {
	return &World{
		dimensions: dimensions,
		blocks:     make(map[Position]BlockType),
	}
}

func (w *World) InBounds(pos Position) bool {
	return pos.X >= 0 && pos.X < w.dimensions.X &&
		pos.Y >= 0 && pos.Y < w.dimensions.Y &&
		pos.Z >= 0 && pos.Z < w.dimensions.Z
}

func (w *World) Get(pos Position) (BlockType, error) {
	if !w.InBounds(pos) {
		return BlockAir, ErrOutOfBounds
	}

	block, ok := w.blocks[pos]
	if !ok {
		return BlockAir, nil
	}
	return block, nil
}

func (w *World) Set(pos Position, block BlockType) error {
	if !w.InBounds(pos) {
		return ErrOutOfBounds
	}
	if !block.Valid() {
		return ErrInvalidBlockType
	}

	if block == BlockAir {
		delete(w.blocks, pos)
		return nil
	}

	w.blocks[pos] = block
	return nil
}

func (w *World) Remove(pos Position) error {
	if !w.InBounds(pos) {
		return ErrOutOfBounds
	}

	delete(w.blocks, pos)
	return nil
}
