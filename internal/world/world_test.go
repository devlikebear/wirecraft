package world

import (
	"errors"
	"testing"
)

func TestDefaultWorldBounds(t *testing.T) {
	w := NewDefault()

	inBounds := []Position{
		{X: 0, Y: 0, Z: 0},
		{X: 31, Y: 31, Z: 15},
	}
	for _, pos := range inBounds {
		if !w.InBounds(pos) {
			t.Fatalf("InBounds(%+v) = false, want true", pos)
		}
	}

	outOfBounds := []Position{
		{X: -1, Y: 0, Z: 0},
		{X: 32, Y: 0, Z: 0},
		{X: 0, Y: 32, Z: 0},
		{X: 0, Y: 0, Z: 16},
	}
	for _, pos := range outOfBounds {
		if w.InBounds(pos) {
			t.Fatalf("InBounds(%+v) = true, want false", pos)
		}
	}
}

func TestWorldGetSetRemove(t *testing.T) {
	w := NewDefault()
	pos := Position{X: 2, Y: 3, Z: 4}

	block, err := w.Get(pos)
	if err != nil {
		t.Fatalf("Get(empty) error = %v, want nil", err)
	}
	if block != BlockAir {
		t.Fatalf("Get(empty) = %s, want %s", block, BlockAir)
	}

	if err := w.Set(pos, BlockSolid); err != nil {
		t.Fatalf("Set(solid) error = %v, want nil", err)
	}
	block, err = w.Get(pos)
	if err != nil {
		t.Fatalf("Get(solid) error = %v, want nil", err)
	}
	if block != BlockSolid {
		t.Fatalf("Get(solid) = %s, want %s", block, BlockSolid)
	}

	if err := w.Remove(pos); err != nil {
		t.Fatalf("Remove() error = %v, want nil", err)
	}
	block, err = w.Get(pos)
	if err != nil {
		t.Fatalf("Get(after remove) error = %v, want nil", err)
	}
	if block != BlockAir {
		t.Fatalf("Get(after remove) = %s, want %s", block, BlockAir)
	}
}

func TestWorldRejectsOutOfBoundsOperations(t *testing.T) {
	w := NewDefault()
	pos := Position{X: 32, Y: 0, Z: 0}

	if _, err := w.Get(pos); !errors.Is(err, ErrOutOfBounds) {
		t.Fatalf("Get(out of bounds) error = %v, want %v", err, ErrOutOfBounds)
	}
	if err := w.Set(pos, BlockSolid); !errors.Is(err, ErrOutOfBounds) {
		t.Fatalf("Set(out of bounds) error = %v, want %v", err, ErrOutOfBounds)
	}
	if err := w.Remove(pos); !errors.Is(err, ErrOutOfBounds) {
		t.Fatalf("Remove(out of bounds) error = %v, want %v", err, ErrOutOfBounds)
	}
}

func TestWorldRejectsInvalidBlockType(t *testing.T) {
	w := NewDefault()

	err := w.Set(Position{X: 0, Y: 0, Z: 0}, BlockType(255))
	if !errors.Is(err, ErrInvalidBlockType) {
		t.Fatalf("Set(invalid block) error = %v, want %v", err, ErrInvalidBlockType)
	}
}

func TestWorldSetRemoveIsDeterministic(t *testing.T) {
	w := NewDefault()
	pos := Position{X: 7, Y: 8, Z: 9}

	for range 3 {
		if err := w.Set(pos, BlockDebugMover); err != nil {
			t.Fatalf("Set(debug mover) error = %v, want nil", err)
		}
	}

	block, err := w.Get(pos)
	if err != nil {
		t.Fatalf("Get(debug mover) error = %v, want nil", err)
	}
	if block != BlockDebugMover {
		t.Fatalf("Get(debug mover) = %s, want %s", block, BlockDebugMover)
	}

	for range 3 {
		if err := w.Remove(pos); err != nil {
			t.Fatalf("Remove(debug mover) error = %v, want nil", err)
		}
	}

	block, err = w.Get(pos)
	if err != nil {
		t.Fatalf("Get(after repeated remove) error = %v, want nil", err)
	}
	if block != BlockAir {
		t.Fatalf("Get(after repeated remove) = %s, want %s", block, BlockAir)
	}
}
