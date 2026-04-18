import { describe, expect, it } from 'vitest';
import { BlockType, parseSnapshot } from './protocol';

describe('parseSnapshot', () => {
  it('normalizes Go block positions into client positions', () => {
    const snapshot = parseSnapshot({
      tick: 7,
      serverTimeMs: 1700000000007,
      blocks: [
        {
          position: { X: 1, Y: 2, Z: 3 },
          blockType: BlockType.DebugMover,
        },
      ],
      entities: [],
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 128,
      },
    });

    expect(snapshot?.blocks[0]?.position).toEqual({ x: 1, y: 2, z: 3 });
  });

  it('rejects invalid payloads', () => {
    const snapshot = parseSnapshot({
      tick: 'bad',
      serverTimeMs: 1700000000007,
      blocks: [],
      entities: [],
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 128,
      },
    });

    expect(snapshot).toBeNull();
  });
});
