import { describe, expect, it } from 'vitest';
import { BlockType, EntityType, type EntitySnapshot, type Snapshot } from '../net/protocol';
import { SnapshotStore } from './snapshotStore';

describe('SnapshotStore', () => {
  it('keeps the latest snapshot', () => {
    const store = new SnapshotStore({ maxSnapshots: 3 });
    const first = snapshot(1);
    const second = snapshot(2);

    store.append(first);
    store.append(second);

    expect(store.latest()).toEqual(second);
    expect(store.length).toBe(2);
  });

  it('caps the snapshot buffer', () => {
    const store = new SnapshotStore({ maxSnapshots: 2 });

    store.append(snapshot(1));
    store.append(snapshot(2));
    store.append(snapshot(3));

    expect(store.all().map((entry) => entry.tick)).toEqual([2, 3]);
  });

  it('returns copies of buffered snapshots', () => {
    const store = new SnapshotStore({ maxSnapshots: 2 });
    store.append(snapshot(1));

    const buffered = store.all();
    buffered.length = 0;

    expect(store.length).toBe(1);
  });

  it('estimates render server time from latest snapshot receipt time and interpolation delay', () => {
    const store = new SnapshotStore({ maxSnapshots: 2 });

    store.append(snapshot(1, 1000), 5000);

    expect(store.renderServerTimeMs(5170, 120)).toBe(1050);
  });

  it('returns null render time when empty', () => {
    const store = new SnapshotStore({ maxSnapshots: 2 });

    expect(store.renderServerTimeMs(5170, 120)).toBeNull();
  });

  it('applies changed-set block updates onto the matching latest base snapshot', () => {
    const store = new SnapshotStore({ maxSnapshots: 3 });

    store.append(
      snapshot(1, 1000, {
        blocks: [
          { position: { x: 0, y: 0, z: 0 }, blockType: BlockType.Power },
          { position: { x: 1, y: 0, z: 0 }, blockType: BlockType.Button },
        ],
      }),
    );
    store.append(
      snapshot(2, 1050, {
        mode: 'changed_set',
        baseTick: 1,
        blocks: [],
        changedBlocks: [{ position: { x: 2, y: 0, z: 0 }, blockType: BlockType.Wire }],
        removedBlocks: [{ x: 1, y: 0, z: 0 }],
      }),
    );

    expect(store.latest()?.blocks).toEqual([
      { position: { x: 0, y: 0, z: 0 }, blockType: BlockType.Power },
      { position: { x: 2, y: 0, z: 0 }, blockType: BlockType.Wire },
    ]);
    expect(store.latest()?.mode).toBe('changed_set');
  });

  it('applies changed-set entities onto buffered snapshots for interpolation', () => {
    const store = new SnapshotStore({ maxSnapshots: 3 });

    store.append(
      snapshot(1, 1000, {
        entities: [entity('motor:1', 1)],
      }),
    );
    store.append(
      snapshot(2, 1050, {
        mode: 'changed_set',
        baseTick: 1,
        entities: [],
        changedEntities: [entity('motor:1', 2), entity('piston:1', 3, EntityType.Piston)],
      }),
    );

    expect(store.all().map((entry) => entry.entities.map((snapshotEntity) => snapshotEntity.id))).toEqual([
      ['motor:1'],
      ['motor:1', 'piston:1'],
    ]);
    expect(store.latest()?.entities.map((snapshotEntity) => snapshotEntity.transform.position.x)).toEqual([
      2, 3,
    ]);
  });

  it('ignores changed-set snapshots when the matching base is unavailable', () => {
    const store = new SnapshotStore({ maxSnapshots: 3 });
    store.append(snapshot(1, 1000));

    store.append(
      snapshot(2, 1050, {
        mode: 'changed_set',
        baseTick: 99,
        changedBlocks: [{ position: { x: 2, y: 0, z: 0 }, blockType: BlockType.Wire }],
      }),
    );

    expect(store.length).toBe(1);
    expect(store.latest()?.tick).toBe(1);
  });
});

function snapshot(
  tick: number,
  serverTimeMs = 1700000000000 + tick,
  overrides: Partial<Snapshot> = {},
): Snapshot {
  return {
    mode: 'full',
    tick,
    serverTimeMs,
    baseTick: undefined,
    blocks: [],
    changedBlocks: [],
    removedBlocks: [],
    entities: [],
    changedEntities: [],
    circuit: { nodes: [] },
    presence: { clients: [] },
    commandAcks: [],
    stats: {
      clientCount: 1,
      commandQueueLength: 0,
      snapshotBytes: 128,
    },
    ...overrides,
  };
}

function entity(id: string, x: number, type = EntityType.Motor): EntitySnapshot {
  return {
    id,
    type,
    transform: {
      position: { x, y: 0.5, z: 0 },
      rotation: { x: 0, y: 0, z: 0, w: 1 },
      scale: { x: 1, y: 1, z: 1 },
    },
  };
}
