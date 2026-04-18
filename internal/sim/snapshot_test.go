package sim

import (
	"testing"

	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/physics"
	"github.com/devlikebear/wirecraft/internal/world"
)

func TestBuildSnapshotFromWorld(t *testing.T) {
	w := world.NewDefault()
	if err := w.Set(world.Position{X: 2, Y: 0, Z: 0}, world.BlockSolid); err != nil {
		t.Fatalf("Set(solid) error = %v, want nil", err)
	}
	if err := w.Set(world.Position{X: 1, Y: 0, Z: 0}, world.BlockDebugMover); err != nil {
		t.Fatalf("Set(debug mover) error = %v, want nil", err)
	}

	snapshot := BuildSnapshot(SnapshotInput{
		Tick:         TickID(9),
		ServerTimeMS: 1700000000999,
		World:        w,
		Stats: SnapshotStatsInput{
			ClientCount:        3,
			CommandQueueLength: 5,
		},
		Presence: netproto.PresenceSnapshot{
			Clients: []netproto.ClientPresenceSnapshot{
				{ID: "client-1", DisplayName: "Client 1"},
			},
		},
	})

	if snapshot.Tick != 9 {
		t.Fatalf("snapshot.Tick = %d, want 9", snapshot.Tick)
	}
	if snapshot.Mode != netproto.SnapshotModeFull {
		t.Fatalf("snapshot.Mode = %q, want %q", snapshot.Mode, netproto.SnapshotModeFull)
	}
	if len(snapshot.ChangedBlocks) != 0 || len(snapshot.RemovedBlocks) != 0 || len(snapshot.ChangedEntities) != 0 {
		t.Fatalf("changed-set fields = blocks:%+v removed:%+v entities:%+v, want empty full snapshot primitives", snapshot.ChangedBlocks, snapshot.RemovedBlocks, snapshot.ChangedEntities)
	}
	if snapshot.ServerTimeMS != 1700000000999 {
		t.Fatalf("snapshot.ServerTimeMS = %d, want 1700000000999", snapshot.ServerTimeMS)
	}

	wantBlocks := []netproto.BlockSnapshot{
		{Position: world.Position{X: 1, Y: 0, Z: 0}, BlockType: world.BlockDebugMover},
		{Position: world.Position{X: 2, Y: 0, Z: 0}, BlockType: world.BlockSolid},
	}
	if len(snapshot.Blocks) != len(wantBlocks) {
		t.Fatalf("len(snapshot.Blocks) = %d, want %d", len(snapshot.Blocks), len(wantBlocks))
	}
	for i := range wantBlocks {
		if snapshot.Blocks[i] != wantBlocks[i] {
			t.Fatalf("snapshot.Blocks[%d] = %+v, want %+v", i, snapshot.Blocks[i], wantBlocks[i])
		}
	}
	if len(snapshot.Entities) != 1 {
		t.Fatalf("len(snapshot.Entities) = %d, want 1", len(snapshot.Entities))
	}
	assertDebugMoverEntity(t, snapshot.Entities[0], netproto.Vec3{X: 2.125, Y: 1.25, Z: 1.5625})
	if snapshot.Stats.ClientCount != 3 {
		t.Fatalf("snapshot.Stats.ClientCount = %d, want 3", snapshot.Stats.ClientCount)
	}
	if snapshot.Stats.CommandQueueLength != 5 {
		t.Fatalf("snapshot.Stats.CommandQueueLength = %d, want 5", snapshot.Stats.CommandQueueLength)
	}
	if len(snapshot.Presence.Clients) != 1 || snapshot.Presence.Clients[0].ID != "client-1" {
		t.Fatalf("snapshot.Presence.Clients = %+v, want client-1", snapshot.Presence.Clients)
	}
	if snapshot.Stats.SnapshotBytes <= 0 {
		t.Fatalf("snapshot.Stats.SnapshotBytes = %d, want positive", snapshot.Stats.SnapshotBytes)
	}
}

func TestBuildSnapshotIncludesDeterministicDebugEntity(t *testing.T) {
	first := BuildSnapshot(SnapshotInput{Tick: TickID(1)})
	second := BuildSnapshot(SnapshotInput{Tick: TickID(2)})

	if len(first.Entities) != 1 {
		t.Fatalf("len(first.Entities) = %d, want 1", len(first.Entities))
	}
	if len(second.Entities) != 1 {
		t.Fatalf("len(second.Entities) = %d, want 1", len(second.Entities))
	}

	assertDebugMoverEntity(t, first.Entities[0], netproto.Vec3{X: 1.125, Y: 1.25, Z: 1.0625})
	assertDebugMoverEntity(t, second.Entities[0], netproto.Vec3{X: 1.25, Y: 1.25, Z: 1.125})
	if first.Entities[0].Transform.Position == second.Entities[0].Transform.Position {
		t.Fatalf("debug entity position did not change: %+v", first.Entities[0].Transform.Position)
	}
}

func TestBuildSnapshotIncludesActuatorEntitiesAfterDebugMover(t *testing.T) {
	snapshot := BuildSnapshot(SnapshotInput{
		Tick: TickID(4),
		ActuatorEntities: []physics.DynamicEntity{
			{
				ID:   physics.EntityID("piston:2:0:0"),
				Type: physics.EntityTypePiston,
				Transform: physics.Transform{
					Position: physics.Vec3{X: 3, Y: 0.5, Z: 0},
					Rotation: physics.IdentityQuat(),
					Scale:    physics.UnitVec3(),
				},
			},
			{
				ID:   physics.EntityID("motor:4:0:0"),
				Type: physics.EntityTypeMotor,
				Transform: physics.Transform{
					Position: physics.Vec3{X: 4, Y: 0.5, Z: 0},
					Rotation: physics.IdentityQuat(),
					Scale:    physics.UnitVec3(),
				},
			},
		},
	})

	if len(snapshot.Entities) != 3 {
		t.Fatalf("len(snapshot.Entities) = %d, want 3: %+v", len(snapshot.Entities), snapshot.Entities)
	}
	if snapshot.Entities[0].ID != netproto.EntityIDDebugMover {
		t.Fatalf("snapshot.Entities[0].ID = %q, want debug mover first", snapshot.Entities[0].ID)
	}
	if snapshot.Entities[1].ID != "motor:4:0:0" || snapshot.Entities[1].Type != netproto.EntityTypeMotor {
		t.Fatalf("snapshot.Entities[1] = %+v, want motor entity", snapshot.Entities[1])
	}
	if snapshot.Entities[2].ID != "piston:2:0:0" || snapshot.Entities[2].Type != netproto.EntityTypePiston {
		t.Fatalf("snapshot.Entities[2] = %+v, want piston entity", snapshot.Entities[2])
	}
	if snapshot.Entities[2].Transform.Position != (netproto.Vec3{X: 3, Y: 0.5, Z: 0}) {
		t.Fatalf("piston position = %+v, want exported actuator transform", snapshot.Entities[2].Transform.Position)
	}
}

func TestBuildSnapshotOmitsAirBlocks(t *testing.T) {
	w := world.NewDefault()
	pos := world.Position{X: 1, Y: 2, Z: 3}

	if err := w.Set(pos, world.BlockSolid); err != nil {
		t.Fatalf("Set(solid) error = %v, want nil", err)
	}
	if err := w.Remove(pos); err != nil {
		t.Fatalf("Remove() error = %v, want nil", err)
	}

	snapshot := BuildSnapshot(SnapshotInput{
		Tick:  TickID(1),
		World: w,
	})

	if len(snapshot.Blocks) != 0 {
		t.Fatalf("snapshot.Blocks = %+v, want empty", snapshot.Blocks)
	}
}

func assertDebugMoverEntity(t *testing.T, entity netproto.EntitySnapshot, wantPosition netproto.Vec3) {
	t.Helper()

	if entity.ID != netproto.EntityIDDebugMover {
		t.Fatalf("entity.ID = %q, want %q", entity.ID, netproto.EntityIDDebugMover)
	}
	if entity.Type != netproto.EntityTypeDebugMover {
		t.Fatalf("entity.Type = %q, want %q", entity.Type, netproto.EntityTypeDebugMover)
	}
	if entity.Transform.Position != wantPosition {
		t.Fatalf("entity.Transform.Position = %+v, want %+v", entity.Transform.Position, wantPosition)
	}
	if entity.Transform.Rotation != (netproto.Quat{X: 0, Y: 0, Z: 0, W: 1}) {
		t.Fatalf("entity.Transform.Rotation = %+v, want identity", entity.Transform.Rotation)
	}
	if entity.Transform.Scale != (netproto.Vec3{X: 0.5, Y: 0.5, Z: 0.5}) {
		t.Fatalf("entity.Transform.Scale = %+v, want half scale", entity.Transform.Scale)
	}
}
