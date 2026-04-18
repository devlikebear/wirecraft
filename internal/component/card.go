package component

import (
	"errors"
	"fmt"
)

var ErrInvalidCard = errors.New("invalid component card")

type Pin struct {
	Name      string
	Direction string
	Signal    string
}

type Card struct {
	ID                  string
	Name                string
	Role                string
	Pins                []Pin
	WiringNotes         []string
	Warnings            []string
	SimplificationNotes []string
}

var starterCards = []Card{
	{
		ID:   "power",
		Name: "Power",
		Role: "Provides a constant High digital signal",
		Pins: []Pin{
			{Name: "OUT", Direction: "output", Signal: "digital high"},
		},
		WiringNotes: []string{
			"Connect OUT to a wire, logic input, or indicator input.",
		},
		Warnings: []string{
			"Real supplies need voltage and current limits before powering hardware.",
		},
		SimplificationNotes: []string{
			"WireCraft treats this as an ideal always-on digital source for starter circuits.",
		},
	},
	{
		ID:   "wire",
		Name: "Wire",
		Role: "Carries a digital signal between adjacent components",
		Pins: []Pin{
			{Name: "A", Direction: "bidirectional", Signal: "digital"},
			{Name: "B", Direction: "bidirectional", Signal: "digital"},
		},
		WiringNotes: []string{
			"Place wires next to circuit blocks to bridge their signals through the voxel grid.",
		},
		Warnings: []string{
			"Real wires have resistance and safe current ratings.",
		},
		SimplificationNotes: []string{
			"Starter wires propagate High, Low, or Unknown without analog voltage drop.",
		},
	},
	{
		ID:   "button",
		Name: "Button",
		Role: "User-controlled digital switch",
		Pins: []Pin{
			{Name: "IN", Direction: "input", Signal: "digital"},
			{Name: "OUT", Direction: "output", Signal: "digital"},
		},
		WiringNotes: []string{
			"Use set_button commands to press or release the button in the authoritative simulation.",
		},
		Warnings: []string{
			"Real buttons bounce and often need pull-up or pull-down resistors.",
		},
		SimplificationNotes: []string{
			"WireCraft models a released button as Low and a pressed button as High.",
		},
	},
	{
		ID:   "and_gate",
		Name: "AND Gate",
		Role: "Outputs High only when two inputs are High",
		Pins: []Pin{
			{Name: "A", Direction: "input", Signal: "digital"},
			{Name: "B", Direction: "input", Signal: "digital"},
			{Name: "OUT", Direction: "output", Signal: "digital"},
		},
		WiringNotes: []string{
			"Feed two input paths into the gate, then connect OUT to another wire or output block.",
		},
		Warnings: []string{
			"Real logic chips need a matching supply voltage and ground reference.",
		},
		SimplificationNotes: []string{
			"The current starter runtime infers gate inputs from adjacent circuit blocks.",
		},
	},
	{
		ID:   "mcu_output",
		Name: "MCU Output",
		Role: "Represents a microcontroller output pin or observed endpoint",
		Pins: []Pin{
			{Name: "IN", Direction: "input", Signal: "digital"},
		},
		WiringNotes: []string{
			"Connect a wire or gate output to observe the resulting server-authoritative signal.",
		},
		Warnings: []string{
			"Real MCU pins cannot drive motors or high-current loads directly.",
		},
		SimplificationNotes: []string{
			"WireCraft uses this block as a readable digital endpoint before code runtime exists.",
		},
	},
	{
		ID:   "led",
		Name: "LED",
		Role: "Visual indicator for a digital output",
		Pins: []Pin{
			{Name: "A", Direction: "input", Signal: "digital high"},
			{Name: "K", Direction: "input", Signal: "ground reference"},
		},
		WiringNotes: []string{
			"Drive the anode through a resistor and connect the cathode to ground in real circuits.",
		},
		Warnings: []string{
			"Real LEDs need current-limiting resistors to avoid damage.",
		},
		SimplificationNotes: []string{
			"LED is documented as a starter card before it becomes a placeable voxel block.",
		},
	},
	{
		ID:   "resistor",
		Name: "Resistor",
		Role: "Limits current or creates a pull-up/pull-down path",
		Pins: []Pin{
			{Name: "A", Direction: "passive", Signal: "analog"},
			{Name: "B", Direction: "passive", Signal: "analog"},
		},
		WiringNotes: []string{
			"Use resistors with LEDs and buttons in real hardware to control current and default voltage.",
		},
		Warnings: []string{
			"Choose resistance and power rating for the actual voltage and load.",
		},
		SimplificationNotes: []string{
			"WireCraft documents resistors now but does not simulate analog resistance in Phase 2.",
		},
	},
}

func StarterCards() []Card {
	return append([]Card(nil), starterCards...)
}

func FindCard(id string) (Card, bool) {
	for _, card := range starterCards {
		if card.ID == id {
			return card, true
		}
	}
	return Card{}, false
}

func ValidateCards(cards []Card) error {
	seenIDs := make(map[string]struct{}, len(cards))
	for _, card := range cards {
		if card.ID == "" || card.Name == "" || card.Role == "" {
			return fmt.Errorf("%w: card %q is missing id, name, or role", ErrInvalidCard, card.ID)
		}
		if _, ok := seenIDs[card.ID]; ok {
			return fmt.Errorf("%w: duplicate card id %q", ErrInvalidCard, card.ID)
		}
		seenIDs[card.ID] = struct{}{}

		if len(card.Pins) == 0 || len(card.WiringNotes) == 0 ||
			len(card.Warnings) == 0 || len(card.SimplificationNotes) == 0 {
			return fmt.Errorf("%w: card %q is missing required content", ErrInvalidCard, card.ID)
		}

		for _, pin := range card.Pins {
			if pin.Name == "" || pin.Direction == "" || pin.Signal == "" {
				return fmt.Errorf("%w: card %q has an incomplete pin", ErrInvalidCard, card.ID)
			}
		}
	}
	return nil
}
