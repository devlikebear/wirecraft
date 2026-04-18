package sim

import (
	"sort"

	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/world"
)

func BuildChangedSetSnapshot(base netproto.Snapshot, next netproto.Snapshot) netproto.Snapshot {
	snapshot := next
	snapshot.Mode = netproto.SnapshotModeChangedSet
	snapshot.BaseTick = base.Tick
	snapshot.Blocks = []netproto.BlockSnapshot{}
	snapshot.Entities = []netproto.EntitySnapshot{}
	snapshot.ChangedBlocks = changedBlocks(base.Blocks, next.Blocks)
	snapshot.RemovedBlocks = removedBlockPositions(base.Blocks, next.Blocks)
	snapshot.ChangedEntities = changedEntities(base.Entities, next.Entities)

	if snapshot.ChangedBlocks == nil {
		snapshot.ChangedBlocks = []netproto.BlockSnapshot{}
	}
	if snapshot.RemovedBlocks == nil {
		snapshot.RemovedBlocks = []world.Position{}
	}
	if snapshot.ChangedEntities == nil {
		snapshot.ChangedEntities = []netproto.EntitySnapshot{}
	}
	if snapshot.CommandAcks == nil {
		snapshot.CommandAcks = []netproto.CommandAckSnapshot{}
	}

	return finalizeSnapshot(snapshot)
}

func changedBlocks(base []netproto.BlockSnapshot, next []netproto.BlockSnapshot) []netproto.BlockSnapshot {
	baseByPosition := make(map[world.Position]world.BlockType, len(base))
	for _, block := range base {
		baseByPosition[block.Position] = block.BlockType
	}

	changed := make([]netproto.BlockSnapshot, 0)
	for _, block := range next {
		if baseType, ok := baseByPosition[block.Position]; !ok || baseType != block.BlockType {
			changed = append(changed, block)
		}
	}
	sort.Slice(changed, func(i, j int) bool {
		return positionLess(changed[i].Position, changed[j].Position)
	})
	return changed
}

func removedBlockPositions(base []netproto.BlockSnapshot, next []netproto.BlockSnapshot) []world.Position {
	nextPositions := make(map[world.Position]struct{}, len(next))
	for _, block := range next {
		nextPositions[block.Position] = struct{}{}
	}

	removed := make([]world.Position, 0)
	for _, block := range base {
		if _, ok := nextPositions[block.Position]; !ok {
			removed = append(removed, block.Position)
		}
	}
	sort.Slice(removed, func(i, j int) bool {
		return positionLess(removed[i], removed[j])
	})
	return removed
}

func changedEntities(base []netproto.EntitySnapshot, next []netproto.EntitySnapshot) []netproto.EntitySnapshot {
	baseByID := make(map[string]netproto.EntitySnapshot, len(base))
	for _, entity := range base {
		baseByID[entity.ID] = entity
	}

	changed := make([]netproto.EntitySnapshot, 0)
	for _, entity := range next {
		if baseEntity, ok := baseByID[entity.ID]; !ok || baseEntity != entity {
			changed = append(changed, entity)
		}
	}
	sort.Slice(changed, func(i, j int) bool {
		return changed[i].ID < changed[j].ID
	})
	return changed
}

func positionLess(left world.Position, right world.Position) bool {
	if left.X != right.X {
		return left.X < right.X
	}
	if left.Y != right.Y {
		return left.Y < right.Y
	}
	return left.Z < right.Z
}
