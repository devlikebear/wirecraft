import { describe, expect, it } from 'vitest';
import { BlockType, EntityType, parseSnapshot } from './protocol';

describe('parseSnapshot', () => {
  it('defaults missing snapshot mode to full and empty changed sets', () => {
    const snapshot = parseSnapshot({
      tick: 6,
      serverTimeMs: 1700000000006,
      blocks: [],
      entities: [],
      circuit: { nodes: [] },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 128,
      },
    });

    expect(snapshot?.mode).toBe('full');
    expect(snapshot?.baseTick).toBeUndefined();
    expect(snapshot?.changedBlocks).toEqual([]);
    expect(snapshot?.removedBlocks).toEqual([]);
    expect(snapshot?.changedEntities).toEqual([]);
  });

  it('normalizes Go block positions into client positions', () => {
    const snapshot = parseSnapshot({
      tick: 7,
      serverTimeMs: 1700000000007,
      blocks: [
        {
          position: { X: 1, Y: 2, Z: 3 },
          blockType: BlockType.DebugMover,
          facing: 'east',
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
    expect(snapshot?.blocks[0]?.facing).toBe('east');
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

  it('parses actuator entity transforms from authoritative snapshots', () => {
    const snapshot = parseSnapshot({
      tick: 11,
      serverTimeMs: 1700000000207,
      blocks: [],
      entities: [
        {
          id: 'piston:2:0:0',
          type: EntityType.Piston,
          transform: {
            position: { x: 3, y: 0.5, z: 0 },
            rotation: { x: 0, y: 0, z: 0, w: 1 },
            scale: { x: 1, y: 1, z: 1 },
          },
        },
        {
          id: 'motor:4:0:0',
          type: EntityType.Motor,
          transform: {
            position: { x: 4, y: 0.5, z: 0 },
            rotation: { x: 0, y: 0, z: 0, w: 1 },
            scale: { x: 1, y: 1, z: 1 },
          },
        },
      ],
      circuit: { nodes: [] },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 256,
      },
    });

    expect(snapshot?.entities.map((entity) => entity.type)).toEqual([
      EntityType.Piston,
      EntityType.Motor,
    ]);
    expect(snapshot?.entities[0]?.transform.position).toEqual({ x: 3, y: 0.5, z: 0 });
  });

  it('parses changed-set snapshot payloads from the server', () => {
    const snapshot = parseSnapshot({
      mode: 'changed_set',
      tick: 14,
      baseTick: 13,
      serverTimeMs: 1700000000357,
      blocks: [],
      changedBlocks: [
        {
          position: { X: 4, Y: 0, Z: 2 },
          blockType: BlockType.Wire,
        },
      ],
      removedBlocks: [{ X: 3, Y: 0, Z: 2 }],
      entities: [],
      changedEntities: [
        {
          id: 'piston:4:0:2',
          type: EntityType.Piston,
          transform: {
            position: { x: 4, y: 0.5, z: 2 },
            rotation: { x: 0, y: 0, z: 0, w: 1 },
            scale: { x: 1, y: 1, z: 1 },
          },
        },
      ],
      circuit: { nodes: [] },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 196,
      },
    });

    expect(snapshot?.mode).toBe('changed_set');
    expect(snapshot?.baseTick).toBe(13);
    expect(snapshot?.changedBlocks).toEqual([
      {
        position: { x: 4, y: 0, z: 2 },
        blockType: BlockType.Wire,
      },
    ]);
    expect(snapshot?.removedBlocks).toEqual([{ x: 3, y: 0, z: 2 }]);
    expect(snapshot?.changedEntities[0]?.id).toBe('piston:4:0:2');
  });

  it('parses optional presence metadata from authoritative snapshots', () => {
    const snapshot = parseSnapshot({
      tick: 12,
      serverTimeMs: 1700000000257,
      blocks: [],
      entities: [],
      circuit: { nodes: [] },
      presence: {
        clients: [
          { id: 'client-1', displayName: 'Client 1' },
          { id: 'client-2', displayName: 'Client 2' },
        ],
      },
      stats: {
        clientCount: 2,
        commandQueueLength: 0,
        snapshotBytes: 144,
      },
    });

    expect(snapshot?.presence.clients).toEqual([
      { id: 'client-1', displayName: 'Client 1' },
      { id: 'client-2', displayName: 'Client 2' },
    ]);
  });

  it('parses optional command acknowledgements from authoritative snapshots', () => {
    const snapshot = parseSnapshot({
      tick: 13,
      serverTimeMs: 1700000000307,
      blocks: [],
      entities: [],
      circuit: { nodes: [] },
      commandAcks: [
        { clientId: 'client-1', commandId: 'cmd-1', status: 'accepted' },
        {
          clientId: 'client-1',
          commandId: 'cmd-1',
          status: 'rejected',
          reason: 'duplicate_command',
        },
      ],
      stats: {
        clientCount: 1,
        commandQueueLength: 2,
        snapshotBytes: 188,
      },
    });

    expect(snapshot?.commandAcks).toEqual([
      { clientId: 'client-1', commandId: 'cmd-1', status: 'accepted' },
      {
        clientId: 'client-1',
        commandId: 'cmd-1',
        status: 'rejected',
        reason: 'duplicate_command',
      },
    ]);
  });

  it('rejects invalid command acknowledgement statuses', () => {
    const snapshot = parseSnapshot({
      tick: 13,
      serverTimeMs: 1700000000307,
      blocks: [],
      entities: [],
      circuit: { nodes: [] },
      commandAcks: [{ clientId: 'client-1', commandId: 'cmd-1', status: 'pending' }],
      stats: {
        clientCount: 1,
        commandQueueLength: 1,
        snapshotBytes: 188,
      },
    });

    expect(snapshot).toBeNull();
  });

  it('rejects invalid snapshot modes', () => {
    const snapshot = parseSnapshot({
      mode: 'delta',
      tick: 14,
      serverTimeMs: 1700000000357,
      blocks: [],
      entities: [],
      circuit: { nodes: [] },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 196,
      },
    });

    expect(snapshot).toBeNull();
  });

  it('rejects invalid block facing values', () => {
    const snapshot = parseSnapshot({
      tick: 14,
      serverTimeMs: 1700000000357,
      blocks: [
        {
          position: { X: 1, Y: 0, Z: 0 },
          blockType: BlockType.Wire,
          facing: 'diagonal',
        },
      ],
      entities: [],
      circuit: { nodes: [] },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 196,
      },
    });

    expect(snapshot).toBeNull();
  });

  it('rejects changed-set snapshots without a numeric base tick', () => {
    const snapshot = parseSnapshot({
      mode: 'changed_set',
      tick: 14,
      serverTimeMs: 1700000000357,
      blocks: [],
      changedBlocks: [],
      removedBlocks: [],
      entities: [],
      changedEntities: [],
      circuit: { nodes: [] },
      stats: {
        clientCount: 1,
        commandQueueLength: 0,
        snapshotBytes: 196,
      },
    });

    expect(snapshot).toBeNull();
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
