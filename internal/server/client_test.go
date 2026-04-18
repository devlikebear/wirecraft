package server

import (
	"testing"

	"github.com/devlikebear/wirecraft/internal/netproto"
)

func TestNewClientPresenceUsesStableDebugIdentity(t *testing.T) {
	presence := NewClientPresence(3)

	if presence.ID != ClientID("client-3") {
		t.Fatalf("presence.ID = %q, want client-3", presence.ID)
	}
	if presence.DisplayName != "Client 3" {
		t.Fatalf("presence.DisplayName = %q, want Client 3", presence.DisplayName)
	}
	if got := presence.Snapshot(); got != (netproto.ClientPresenceSnapshot{ID: "client-3", DisplayName: "Client 3"}) {
		t.Fatalf("presence.Snapshot() = %+v, want client-3 snapshot", got)
	}
}
