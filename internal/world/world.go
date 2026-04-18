package world

import (
	"errors"
	"sort"
)

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

type Block struct {
	Position  Position
	BlockType BlockType
}

type World struct {
	dimensions Dimensions
	blocks     map[Position]BlockType
}

func NewDefault() *World {
	return New(DefaultDimensions())
}

func DefaultDimensions() Dimensions {
	return Dimensions{X: 32, Y: 32, Z: 16}
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

func (w *World) OccupiedBlocks() []Block {
	blocks := make([]Block, 0, len(w.blocks))
	for pos, blockType := range w.blocks {
		if blockType == BlockAir {
			continue
		}
		blocks = append(blocks, Block{
			Position:  pos,
			BlockType: blockType,
		})
	}

	sort.Slice(blocks, func(i, j int) bool {
		a := blocks[i].Position
		b := blocks[j].Position
		if a.X != b.X {
			return a.X < b.X
		}
		if a.Y != b.Y {
			return a.Y < b.Y
		}
		return a.Z < b.Z
	})

	return blocks
}
