package netproto

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/devlikebear/wirecraft/internal/sim"
	"github.com/devlikebear/wirecraft/internal/world"
)

func TestCommandValidateAcceptsPlaceAndRemove(t *testing.T) {
	bounds := world.Dimensions{X: 32, Y: 32, Z: 16}

	place := Command{
		Type:      CommandPlaceBlock,
		ClientID:  "client-1",
		CommandID: "cmd-1",
		TickHint:  sim.TickID(7),
		Position:  world.Position{X: 1, Y: 2, Z: 3},
		BlockType: world.BlockSolid,
	}
	if err := place.Validate(bounds); err != nil {
		t.Fatalf("Validate(place) error = %v, want nil", err)
	}

	remove := Command{
		Type:      CommandRemoveBlock,
		ClientID:  "client-1",
		CommandID: "cmd-2",
		TickHint:  sim.TickID(8),
		Position:  world.Position{X: 1, Y: 2, Z: 3},
		BlockType: world.BlockAir,
	}
	if err := remove.Validate(bounds); err != nil {
		t.Fatalf("Validate(remove) error = %v, want nil", err)
	}
}

func TestCommandValidateRejectsMissingIdentifiers(t *testing.T) {
	bounds := world.Dimensions{X: 32, Y: 32, Z: 16}

	command := Command{
		Type:      CommandPlaceBlock,
		Position:  world.Position{X: 1, Y: 2, Z: 3},
		BlockType: world.BlockSolid,
	}

	if err := command.Validate(bounds); !errors.Is(err, ErrMissingIdentifier) {
		t.Fatalf("Validate(missing identifiers) error = %v, want %v", err, ErrMissingIdentifier)
	}
}

func TestCommandValidateRejectsOutOfBoundsPosition(t *testing.T) {
	bounds := world.Dimensions{X: 32, Y: 32, Z: 16}

	command := Command{
		Type:      CommandPlaceBlock,
		ClientID:  "client-1",
		CommandID: "cmd-1",
		Position:  world.Position{X: 32, Y: 0, Z: 0},
		BlockType: world.BlockSolid,
	}

	if err := command.Validate(bounds); !errors.Is(err, world.ErrOutOfBounds) {
		t.Fatalf("Validate(out of bounds) error = %v, want %v", err, world.ErrOutOfBounds)
	}
}

func TestCommandValidateRejectsInvalidBlockType(t *testing.T) {
	bounds := world.Dimensions{X: 32, Y: 32, Z: 16}

	command := Command{
		Type:      CommandPlaceBlock,
		ClientID:  "client-1",
		CommandID: "cmd-1",
		Position:  world.Position{X: 1, Y: 2, Z: 3},
		BlockType: world.BlockType(255),
	}

	if err := command.Validate(bounds); !errors.Is(err, world.ErrInvalidBlockType) {
		t.Fatalf("Validate(invalid block) error = %v, want %v", err, world.ErrInvalidBlockType)
	}
}

func TestCommandValidateRejectsUnknownCommandType(t *testing.T) {
	bounds := world.Dimensions{X: 32, Y: 32, Z: 16}

	command := Command{
		Type:      CommandType("teleport_block"),
		ClientID:  "client-1",
		CommandID: "cmd-1",
		Position:  world.Position{X: 1, Y: 2, Z: 3},
		BlockType: world.BlockSolid,
	}

	if err := command.Validate(bounds); !errors.Is(err, ErrUnknownCommandType) {
		t.Fatalf("Validate(unknown command) error = %v, want %v", err, ErrUnknownCommandType)
	}
}

func TestCommandJSONRoundTrip(t *testing.T) {
	command := Command{
		Type:      CommandPlaceBlock,
		ClientID:  "client-1",
		CommandID: "cmd-1",
		TickHint:  sim.TickID(11),
		Position:  world.Position{X: 4, Y: 5, Z: 6},
		BlockType: world.BlockDebugMover,
	}

	encoded, err := json.Marshal(command)
	if err != nil {
		t.Fatalf("Marshal(command) error = %v, want nil", err)
	}

	var decoded Command
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("Unmarshal(command) error = %v, want nil", err)
	}

	if decoded != command {
		t.Fatalf("decoded command = %+v, want %+v", decoded, command)
	}
}
