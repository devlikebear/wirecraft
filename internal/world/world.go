package world

import (
	"errors"
	"sort"
)

var (
	ErrOutOfBounds      = errors.New("position out of bounds")
	ErrInvalidBlockType = errors.New("invalid block type")
	ErrInvalidFacing    = errors.New("invalid block facing")
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
	Facing    Facing
}

type World struct {
	dimensions Dimensions
	blocks     map[Position]Block
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
		blocks:     make(map[Position]Block),
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
	return block.BlockType, nil
}

func (w *World) GetBlock(pos Position) (Block, error) {
	if !w.InBounds(pos) {
		return Block{}, ErrOutOfBounds
	}

	block, ok := w.blocks[pos]
	if !ok {
		return Block{Position: pos, BlockType: BlockAir}, nil
	}
	return block, nil
}

func (w *World) Set(pos Position, block BlockType) error {
	return w.SetBlock(Block{Position: pos, BlockType: block})
}

func (w *World) SetBlock(block Block) error {
	pos := block.Position
	if !w.InBounds(pos) {
		return ErrOutOfBounds
	}
	if !block.BlockType.Valid() {
		return ErrInvalidBlockType
	}
	if !block.Facing.Valid() {
		return ErrInvalidFacing
	}

	if block.BlockType == BlockAir {
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
	for _, block := range w.blocks {
		if block.BlockType == BlockAir {
			continue
		}
		blocks = append(blocks, block)
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
