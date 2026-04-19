package netproto

import (
	"encoding/json"
	"testing"

	"github.com/devlikebear/wirecraft/internal/world"
)

func TestSnapshotJSONRoundTrip(t *testing.T) {
	snapshot := Snapshot{
		Mode:         SnapshotModeChangedSet,
		Tick:         12,
		BaseTick:     11,
		ServerTimeMS: 1700000000123,
		Blocks: []BlockSnapshot{
			{
				Position:  world.Position{X: 1, Y: 2, Z: 3},
				BlockType: world.BlockSolid,
				Facing:    world.FacingSouth,
			},
		},
		Entities: []EntitySnapshot{
			{
				ID:   "debug-mover-1",
				Type: "debug_mover",
				Transform: TransformSnapshot{
					Position: Vec3{X: 1.25, Y: 2.5, Z: 3.75},
					Rotation: Quat{X: 0, Y: 0, Z: 0, W: 1},
					Scale:    Vec3{X: 1, Y: 1, Z: 1},
				},
			},
			{
				ID:   "piston:2:0:0",
				Type: EntityTypePiston,
				Transform: TransformSnapshot{
					Position: Vec3{X: 3, Y: 0.5, Z: 0},
					Rotation: Quat{X: 0, Y: 0, Z: 0, W: 1},
					Scale:    Vec3{X: 1, Y: 1, Z: 1},
				},
			},
		},
		Presence: PresenceSnapshot{
			Clients: []ClientPresenceSnapshot{
				{ID: "client-1", DisplayName: "Client 1"},
				{ID: "client-2", DisplayName: "Client 2"},
			},
		},
		CommandAcks: []CommandAckSnapshot{
			{ClientID: "client-1", CommandID: "cmd-1", Status: CommandAckAccepted},
			{ClientID: "client-1", CommandID: "cmd-1", Status: CommandAckRejected, Reason: "duplicate_command"},
		},
		ChangedBlocks: []BlockSnapshot{
			{
				Position:  world.Position{X: 4, Y: 0, Z: 0},
				BlockType: world.BlockWire,
				Facing:    world.FacingEast,
			},
		},
		RemovedBlocks: []world.Position{
			{X: 5, Y: 0, Z: 0},
		},
		ChangedEntities: []EntitySnapshot{
			{
				ID:   "motor:4:0:0",
				Type: EntityTypeMotor,
				Transform: TransformSnapshot{
					Position: Vec3{X: 4, Y: 0.5, Z: 0},
					Rotation: Quat{X: 0, Y: 0, Z: 0, W: 1},
					Scale:    Vec3{X: 1, Y: 1, Z: 1},
				},
			},
		},
		Stats: SnapshotStats{
			ClientCount:        2,
			CommandQueueLength: 4,
			SnapshotBytes:      512,
		},
	}

	encoded, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatalf("Marshal(snapshot) error = %v, want nil", err)
	}

	var decoded Snapshot
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("Unmarshal(snapshot) error = %v, want nil", err)
	}

	if decoded.Tick != snapshot.Tick {
		t.Fatalf("decoded Tick = %d, want %d", decoded.Tick, snapshot.Tick)
	}
	if decoded.Mode != snapshot.Mode {
		t.Fatalf("decoded Mode = %q, want %q", decoded.Mode, snapshot.Mode)
	}
	if decoded.BaseTick != snapshot.BaseTick {
		t.Fatalf("decoded BaseTick = %d, want %d", decoded.BaseTick, snapshot.BaseTick)
	}
	if decoded.ServerTimeMS != snapshot.ServerTimeMS {
		t.Fatalf("decoded ServerTimeMS = %d, want %d", decoded.ServerTimeMS, snapshot.ServerTimeMS)
	}
	if len(decoded.Blocks) != 1 || decoded.Blocks[0] != snapshot.Blocks[0] {
		t.Fatalf("decoded Blocks = %+v, want %+v", decoded.Blocks, snapshot.Blocks)
	}
	if len(decoded.Entities) != len(snapshot.Entities) {
		t.Fatalf("len(decoded.Entities) = %d, want %d", len(decoded.Entities), len(snapshot.Entities))
	}
	for i := range snapshot.Entities {
		if decoded.Entities[i] != snapshot.Entities[i] {
			t.Fatalf("decoded Entities[%d] = %+v, want %+v", i, decoded.Entities[i], snapshot.Entities[i])
		}
	}
	if len(decoded.Presence.Clients) != len(snapshot.Presence.Clients) {
		t.Fatalf("len(decoded.Presence.Clients) = %d, want %d", len(decoded.Presence.Clients), len(snapshot.Presence.Clients))
	}
	for i := range snapshot.Presence.Clients {
		if decoded.Presence.Clients[i] != snapshot.Presence.Clients[i] {
			t.Fatalf("decoded Presence.Clients[%d] = %+v, want %+v", i, decoded.Presence.Clients[i], snapshot.Presence.Clients[i])
		}
	}
	if len(decoded.CommandAcks) != len(snapshot.CommandAcks) {
		t.Fatalf("len(decoded.CommandAcks) = %d, want %d", len(decoded.CommandAcks), len(snapshot.CommandAcks))
	}
	for i := range snapshot.CommandAcks {
		if decoded.CommandAcks[i] != snapshot.CommandAcks[i] {
			t.Fatalf("decoded CommandAcks[%d] = %+v, want %+v", i, decoded.CommandAcks[i], snapshot.CommandAcks[i])
		}
	}
	if len(decoded.ChangedBlocks) != len(snapshot.ChangedBlocks) {
		t.Fatalf("len(decoded.ChangedBlocks) = %d, want %d", len(decoded.ChangedBlocks), len(snapshot.ChangedBlocks))
	}
	for i := range snapshot.ChangedBlocks {
		if decoded.ChangedBlocks[i] != snapshot.ChangedBlocks[i] {
			t.Fatalf("decoded ChangedBlocks[%d] = %+v, want %+v", i, decoded.ChangedBlocks[i], snapshot.ChangedBlocks[i])
		}
	}
	if len(decoded.RemovedBlocks) != len(snapshot.RemovedBlocks) {
		t.Fatalf("len(decoded.RemovedBlocks) = %d, want %d", len(decoded.RemovedBlocks), len(snapshot.RemovedBlocks))
	}
	for i := range snapshot.RemovedBlocks {
		if decoded.RemovedBlocks[i] != snapshot.RemovedBlocks[i] {
			t.Fatalf("decoded RemovedBlocks[%d] = %+v, want %+v", i, decoded.RemovedBlocks[i], snapshot.RemovedBlocks[i])
		}
	}
	if len(decoded.ChangedEntities) != len(snapshot.ChangedEntities) {
		t.Fatalf("len(decoded.ChangedEntities) = %d, want %d", len(decoded.ChangedEntities), len(snapshot.ChangedEntities))
	}
	for i := range snapshot.ChangedEntities {
		if decoded.ChangedEntities[i] != snapshot.ChangedEntities[i] {
			t.Fatalf("decoded ChangedEntities[%d] = %+v, want %+v", i, decoded.ChangedEntities[i], snapshot.ChangedEntities[i])
		}
	}
	if decoded.Stats != snapshot.Stats {
		t.Fatalf("decoded Stats = %+v, want %+v", decoded.Stats, snapshot.Stats)
	}
}
