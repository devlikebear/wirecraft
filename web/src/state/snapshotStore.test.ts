import { describe, expect, it } from 'vitest';
import type { Snapshot } from '../net/protocol';
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
});

function snapshot(tick: number): Snapshot {
  return {
    tick,
    serverTimeMs: 1700000000000 + tick,
    blocks: [],
    entities: [],
    stats: {
      clientCount: 1,
      commandQueueLength: 0,
      snapshotBytes: 128,
    },
  };
}
