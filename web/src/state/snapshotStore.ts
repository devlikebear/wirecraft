import type { BlockSnapshot, EntitySnapshot, Position, Snapshot } from '../net/protocol';

export interface SnapshotStoreOptions {
  maxSnapshots?: number;
}

interface SnapshotEntry {
  snapshot: Snapshot;
  receivedAtMs: number;
}

export class SnapshotStore {
  private readonly maxSnapshots: number;
  private entries: SnapshotEntry[] = [];

  constructor(options: SnapshotStoreOptions = {}) {
    this.maxSnapshots = Math.max(1, options.maxSnapshots ?? 64);
  }

  get length(): number {
    return this.entries.length;
  }

  append(snapshot: Snapshot, receivedAtMs = Date.now()): void {
    const materialized = this.materialize(snapshot);
    if (materialized === null) {
      return;
    }

    this.entries = [...this.entries, { snapshot: materialized, receivedAtMs }].slice(-this.maxSnapshots);
  }

  latest(): Snapshot | null {
    return this.entries.at(-1)?.snapshot ?? null;
  }

  all(): Snapshot[] {
    return this.entries.map((entry) => entry.snapshot);
  }

  renderServerTimeMs(clientTimeMs: number, interpolationDelayMs: number): number | null {
    const latest = this.entries.at(-1);
    if (!latest) {
      return null;
    }

    const elapsedSinceLatestMs = Math.max(0, clientTimeMs - latest.receivedAtMs);
    return latest.snapshot.serverTimeMs + elapsedSinceLatestMs - Math.max(0, interpolationDelayMs);
  }

  clear(): void {
    this.entries = [];
  }

  private materialize(snapshot: Snapshot): Snapshot | null {
    if (snapshot.mode !== 'changed_set') {
      return snapshot;
    }

    const base = this.entries.at(-1)?.snapshot;
    if (!base || base.tick !== snapshot.baseTick) {
      return null;
    }

    return {
      ...snapshot,
      blocks: applyBlockChanges(base.blocks, snapshot.changedBlocks ?? [], snapshot.removedBlocks ?? []),
      entities: applyEntityChanges(base.entities, snapshot.changedEntities ?? []),
    };
  }
}

function applyBlockChanges(
  baseBlocks: BlockSnapshot[],
  changedBlocks: BlockSnapshot[],
  removedBlocks: Position[],
): BlockSnapshot[] {
  const byPosition = new Map<string, BlockSnapshot>();
  for (const block of baseBlocks) {
    byPosition.set(positionKey(block.position), block);
  }
  for (const position of removedBlocks) {
    byPosition.delete(positionKey(position));
  }
  for (const block of changedBlocks) {
    byPosition.set(positionKey(block.position), block);
  }

  return [...byPosition.values()].sort((left, right) => comparePosition(left.position, right.position));
}

function applyEntityChanges(
  baseEntities: EntitySnapshot[],
  changedEntities: EntitySnapshot[],
): EntitySnapshot[] {
  const byID = new Map<string, EntitySnapshot>();
  for (const entity of baseEntities) {
    byID.set(entity.id, entity);
  }
  for (const entity of changedEntities) {
    byID.set(entity.id, entity);
  }

  return [...byID.values()].sort((left, right) => left.id.localeCompare(right.id));
}

function positionKey(position: Position): string {
  return `${position.x}:${position.y}:${position.z}`;
}

function comparePosition(left: Position, right: Position): number {
  if (left.x !== right.x) {
    return left.x - right.x;
  }
  if (left.y !== right.y) {
    return left.y - right.y;
  }
  return left.z - right.z;
}
