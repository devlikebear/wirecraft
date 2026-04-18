package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/world"
)

func TestWebSocketRejectsPlainHTTP(t *testing.T) {
	handler := NewWithOptions(Options{TickRateHz: 200})
	defer handler.Close()

	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUpgradeRequired {
		t.Fatalf("GET /ws status = %d, want %d", rec.Code, http.StatusUpgradeRequired)
	}
}

func TestWebSocketSendsSnapshots(t *testing.T) {
	handler := NewWithOptions(Options{TickRateHz: 200})
	defer handler.Close()
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, websocketURL(server.URL, "/ws"), nil)
	if err != nil {
		t.Fatalf("Dial(/ws) error = %v, want nil", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	var snapshot netproto.Snapshot
	if err := wsjson.Read(ctx, conn, &snapshot); err != nil {
		t.Fatalf("Read(snapshot) error = %v, want nil", err)
	}

	if snapshot.Tick == 0 {
		t.Fatalf("snapshot.Tick = %d, want positive", snapshot.Tick)
	}
	if snapshot.Stats.SnapshotBytes <= 0 {
		t.Fatalf("snapshot.Stats.SnapshotBytes = %d, want positive", snapshot.Stats.SnapshotBytes)
	}
}

func TestWebSocketAppliesCommandAndReturnsAuthoritativeSnapshot(t *testing.T) {
	handler := NewWithOptions(Options{TickRateHz: 200})
	defer handler.Close()
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, websocketURL(server.URL, "/ws"), nil)
	if err != nil {
		t.Fatalf("Dial(/ws) error = %v, want nil", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	pos := world.Position{X: 4, Y: 5, Z: 6}
	command := netproto.Command{
		Type:      netproto.CommandPlaceBlock,
		ClientID:  "client-1",
		CommandID: "cmd-1",
		Position:  pos,
		BlockType: world.BlockDebugMover,
	}
	if err := wsjson.Write(ctx, conn, command); err != nil {
		t.Fatalf("Write(command) error = %v, want nil", err)
	}

	for {
		var snapshot netproto.Snapshot
		if err := wsjson.Read(ctx, conn, &snapshot); err != nil {
			t.Fatalf("Read(snapshot) error = %v, want nil", err)
		}
		if snapshotContainsBlock(snapshot, pos, world.BlockDebugMover) {
			return
		}
	}
}

func websocketURL(httpURL string, path string) string {
	return "ws" + strings.TrimPrefix(httpURL, "http") + path
}

func snapshotContainsBlock(snapshot netproto.Snapshot, pos world.Position, blockType world.BlockType) bool {
	for _, block := range snapshot.Blocks {
		if block.Position == pos && block.BlockType == blockType {
			return true
		}
	}
	return false
}
