package netproto

import (
	"errors"

	"github.com/devlikebear/wirecraft/internal/sim"
	"github.com/devlikebear/wirecraft/internal/world"
)

var (
	ErrMissingIdentifier  = errors.New("missing command identifier")
	ErrUnknownCommandType = errors.New("unknown command type")
)

type CommandType string

const (
	CommandPlaceBlock  CommandType = "place_block"
	CommandRemoveBlock CommandType = "remove_block"
)

type Command struct {
	Type      CommandType     `json:"type"`
	ClientID  string          `json:"clientId"`
	CommandID string          `json:"commandId"`
	TickHint  sim.TickID      `json:"tickHint"`
	Position  world.Position  `json:"position"`
	BlockType world.BlockType `json:"blockType"`
}

func (c Command) Validate(bounds world.Dimensions) error {
	if c.ClientID == "" || c.CommandID == "" {
		return ErrMissingIdentifier
	}
	if !positionInBounds(c.Position, bounds) {
		return world.ErrOutOfBounds
	}

	switch c.Type {
	case CommandPlaceBlock:
		if !c.BlockType.Valid() {
			return world.ErrInvalidBlockType
		}
		if c.BlockType == world.BlockAir {
			return world.ErrInvalidBlockType
		}
		return nil
	case CommandRemoveBlock:
		return nil
	default:
		return ErrUnknownCommandType
	}
}

func positionInBounds(pos world.Position, bounds world.Dimensions) bool {
	return pos.X >= 0 && pos.X < bounds.X &&
		pos.Y >= 0 && pos.Y < bounds.Y &&
		pos.Z >= 0 && pos.Z < bounds.Z
}
