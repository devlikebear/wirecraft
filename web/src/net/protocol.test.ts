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
      circuit: { nodes: [] },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 128,
      },
    });

    expect(snapshot?.blocks[0]?.position).toEqual({ x: 1, y: 2, z: 3 });
  });

  it('accepts circuit block types in authoritative snapshots', () => {
    const snapshot = parseSnapshot({
      tick: 8,
      serverTimeMs: 1700000000057,
      blocks: [
        {
          position: { X: 2, Y: 0, Z: 3 },
          blockType: BlockType.Power,
        },
        {
          position: { X: 3, Y: 0, Z: 3 },
          blockType: BlockType.MCUOutput,
        },
      ],
      entities: [],
      circuit: { nodes: [] },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 160,
      },
    });

    expect(snapshot?.blocks.map((block) => block.blockType)).toEqual([
      BlockType.Power,
      BlockType.MCUOutput,
    ]);
  });

  it('accepts actuator block types in authoritative snapshots', () => {
    const snapshot = parseSnapshot({
      tick: 10,
      serverTimeMs: 1700000000157,
      blocks: [
        {
          position: { X: 6, Y: 0, Z: 3 },
          blockType: BlockType.Piston,
        },
        {
          position: { X: 7, Y: 0, Z: 3 },
          blockType: BlockType.Motor,
        },
        {
          position: { X: 8, Y: 0, Z: 3 },
          blockType: BlockType.MotorDriver,
        },
        {
          position: { X: 9, Y: 0, Z: 3 },
          blockType: BlockType.TransistorSwitch,
        },
      ],
      entities: [],
      circuit: { nodes: [] },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 224,
      },
    });

    expect(snapshot?.blocks.map((block) => block.blockType)).toEqual([
      BlockType.Piston,
      BlockType.Motor,
      BlockType.MotorDriver,
      BlockType.TransistorSwitch,
    ]);
  });

  it('parses circuit signal state from authoritative snapshots', () => {
    const snapshot = parseSnapshot({
      tick: 9,
      serverTimeMs: 1700000000107,
      blocks: [],
      entities: [],
      circuit: {
        nodes: [
          {
            position: { X: 4, Y: 1, Z: 2 },
            nodeId: '4:1:2',
            nodeType: 'wire',
            signalState: 'high',
          },
          {
            position: { X: 5, Y: 1, Z: 2 },
            nodeId: '5:1:2',
            nodeType: 'mcu_output',
            signalState: 'low',
          },
        ],
      },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 192,
      },
    });

    expect(snapshot?.circuit.nodes).toEqual([
      {
        position: { x: 4, y: 1, z: 2 },
        nodeId: '4:1:2',
        nodeType: 'wire',
        signalState: 'high',
      },
      {
        position: { x: 5, y: 1, z: 2 },
        nodeId: '5:1:2',
        nodeType: 'mcu_output',
        signalState: 'low',
      },
    ]);
  });

  it('rejects invalid circuit signal state', () => {
    const snapshot = parseSnapshot({
      tick: 9,
      serverTimeMs: 1700000000107,
      blocks: [],
      entities: [],
      circuit: {
        nodes: [
          {
            position: { X: 4, Y: 1, Z: 2 },
            nodeId: '4:1:2',
            nodeType: 'wire',
            signalState: 'energized',
          },
        ],
      },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 192,
      },
    });

    expect(snapshot).toBeNull();
  });

  it('rejects invalid payloads', () => {
    const snapshot = parseSnapshot({
      tick: 'bad',
      serverTimeMs: 1700000000007,
      blocks: [],
      entities: [],
      circuit: { nodes: [] },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 128,
      },
    });

    expect(snapshot).toBeNull();
  });
});
