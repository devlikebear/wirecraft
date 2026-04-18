package netproto

import (
	"encoding/json"
	"testing"

	"github.com/devlikebear/wirecraft/internal/world"
)

func TestSnapshotJSONRoundTrip(t *testing.T) {
	snapshot := Snapshot{
		Tick:         12,
		ServerTimeMS: 1700000000123,
		Blocks: []BlockSnapshot{
			{
				Position:  world.Position{X: 1, Y: 2, Z: 3},
				BlockType: world.BlockSolid,
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
	if decoded.ServerTimeMS != snapshot.ServerTimeMS {
		t.Fatalf("decoded ServerTimeMS = %d, want %d", decoded.ServerTimeMS, snapshot.ServerTimeMS)
	}
	if len(decoded.Blocks) != 1 || decoded.Blocks[0] != snapshot.Blocks[0] {
		t.Fatalf("decoded Blocks = %+v, want %+v", decoded.Blocks, snapshot.Blocks)
	}
	if len(decoded.Entities) != 1 || decoded.Entities[0] != snapshot.Entities[0] {
		t.Fatalf("decoded Entities = %+v, want %+v", decoded.Entities, snapshot.Entities)
	}
	if decoded.Stats != snapshot.Stats {
		t.Fatalf("decoded Stats = %+v, want %+v", decoded.Stats, snapshot.Stats)
	}
}
