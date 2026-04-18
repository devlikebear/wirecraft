package sensor

import (
	"testing"

	"github.com/devlikebear/wirecraft/internal/world"
)

func TestSensorInputStatePrimitives(t *testing.T) {
	pos := world.Position{X: 1, Y: 2, Z: 3}

	button := ButtonInput(pos, true)
	if button.Type != SensorTypeButton {
		t.Fatalf("button.Type = %q, want %q", button.Type, SensorTypeButton)
	}
	if button.Position != pos {
		t.Fatalf("button.Position = %+v, want %+v", button.Position, pos)
	}
	if !button.Active {
		t.Fatalf("button.Active = false, want true")
	}

	proximity := ProximityStubInput(pos)
	if proximity.Type != SensorTypeProximityStub {
		t.Fatalf("proximity.Type = %q, want %q", proximity.Type, SensorTypeProximityStub)
	}
	if proximity.Active {
		t.Fatalf("proximity.Active = true, want inert stub")
	}
}

func TestStoreExportsOnlyActiveButtonStates(t *testing.T) {
	store := NewStore()
	buttonPos := world.Position{X: 1, Y: 0, Z: 0}
	proximityPos := world.Position{X: 2, Y: 0, Z: 0}

	store.Set(ButtonInput(buttonPos, true))
	store.Set(InputState{
		Type:     SensorTypeProximityStub,
		Position: proximityPos,
		Active:   true,
	})

	buttonStates := store.ButtonStates()
	if !buttonStates[buttonPos] {
		t.Fatalf("buttonStates[%+v] = false, want true", buttonPos)
	}
	if buttonStates[proximityPos] {
		t.Fatalf("buttonStates[%+v] = true, want proximity stub to stay inert", proximityPos)
	}
}

func TestStoreClearsInactiveAndRemovedInputs(t *testing.T) {
	store := NewStore()
	pos := world.Position{X: 1, Y: 0, Z: 0}

	store.Set(ButtonInput(pos, true))
	store.Set(ButtonInput(pos, false))
	if len(store.ButtonStates()) != 0 {
		t.Fatalf("ButtonStates after inactive button = %+v, want empty", store.ButtonStates())
	}

	store.Set(ButtonInput(pos, true))
	store.Clear(pos)
	if len(store.ButtonStates()) != 0 {
		t.Fatalf("ButtonStates after Clear = %+v, want empty", store.ButtonStates())
	}
}
