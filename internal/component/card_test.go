package component

import "testing"

func TestStarterCardsIncludeRequiredComponents(t *testing.T) {
	cards := StarterCards()

	wantIDs := []string{"power", "wire", "button", "and_gate", "mcu_output", "led", "resistor"}
	gotIDs := make([]string, 0, len(cards))
	for _, card := range cards {
		gotIDs = append(gotIDs, card.ID)
	}

	if len(gotIDs) != len(wantIDs) {
		t.Fatalf("len(StarterCards()) = %d, want %d: %+v", len(gotIDs), len(wantIDs), gotIDs)
	}
	for i := range wantIDs {
		if gotIDs[i] != wantIDs[i] {
			t.Fatalf("StarterCards()[%d].ID = %q, want %q", i, gotIDs[i], wantIDs[i])
		}
	}
}

func TestValidateCardsAcceptsStarterCards(t *testing.T) {
	if err := ValidateCards(StarterCards()); err != nil {
		t.Fatalf("ValidateCards(StarterCards()) error = %v, want nil", err)
	}
}

func TestValidateCardsRejectsMissingRequiredFields(t *testing.T) {
	cards := []Card{
		{
			ID:                  "incomplete",
			Name:                "Incomplete",
			Role:                "test fixture",
			Pins:                []Pin{{Name: "A", Direction: "input", Signal: "digital"}},
			WiringNotes:         []string{"Wire it."},
			Warnings:            []string{"Do not do the thing."},
			SimplificationNotes: nil,
		},
	}

	if err := ValidateCards(cards); err == nil {
		t.Fatal("ValidateCards(incomplete) error = nil, want non-nil")
	}
}

func TestFindCard(t *testing.T) {
	card, ok := FindCard("and_gate")
	if !ok {
		t.Fatal("FindCard(and_gate) ok = false, want true")
	}
	if card.Name != "AND Gate" {
		t.Fatalf("FindCard(and_gate).Name = %q, want AND Gate", card.Name)
	}

	if _, ok := FindCard("missing"); ok {
		t.Fatal("FindCard(missing) ok = true, want false")
	}
}
