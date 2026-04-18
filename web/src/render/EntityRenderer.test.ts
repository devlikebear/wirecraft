import { describe, expect, it } from 'vitest';
import { Group } from 'three';
import type { EntitySnapshot, Snapshot } from '../net/protocol';
import {
  DEBUG_MOVER_ENTITY_TYPE,
  EntityRenderer,
  selectInterpolatedEntityRenderItems,
} from './EntityRenderer';

describe('selectInterpolatedEntityRenderItems', () => {
  it('interpolates matching debug mover entities between buffered snapshots', () => {
    const items = selectInterpolatedEntityRenderItems(
      [snapshot(1, 1000, entity('debug-mover-1', 0)), snapshot(2, 1100, entity('debug-mover-1', 10))],
      1050,
    );

    expect(items).toHaveLength(1);
    expect(items[0].id).toBe('debug-mover-1');
    expect(items[0].transform.position).toEqual({ x: 5, y: 1.25, z: 1 });
    expect(items[0].transform.scale).toEqual({ x: 0.5, y: 0.5, z: 0.5 });
  });

  it('falls back to the nearest available snapshot when outside the buffered range', () => {
    const items = selectInterpolatedEntityRenderItems([snapshot(1, 1000, entity('debug-mover-1', 3))], 1200);

    expect(items).toHaveLength(1);
    expect(items[0].transform.position).toEqual({ x: 3, y: 1.25, z: 1 });
  });

  it('ignores non-debug mover entities', () => {
    const items = selectInterpolatedEntityRenderItems(
      [snapshot(1, 1000, entity('servo-1', 0, 'servo'))],
      1000,
    );

    expect(items).toEqual([]);
  });
});

describe('EntityRenderer', () => {
  it('creates and updates debug mover meshes by entity ID', () => {
    const parent = new Group();
    const renderer = new EntityRenderer(parent);

    renderer.updateFromSnapshots(
      [snapshot(1, 1000, entity('debug-mover-1', 0)), snapshot(2, 1100, entity('debug-mover-1', 10))],
      1050,
    );

    expect(parent.children).toContain(renderer.object);
    expect(renderer.count).toBe(1);

    const mesh = renderer.object.getObjectByName('wirecraft-entity-debug-mover-1');
    expect(mesh?.position.x).toBe(5);
    expect(mesh?.position.y).toBe(1.25);
    expect(mesh?.scale.x).toBe(0.5);

    renderer.dispose();
    expect(parent.children).not.toContain(renderer.object);
  });
});

function snapshot(tick: number, serverTimeMs: number, ...entities: EntitySnapshot[]): Snapshot {
  return {
    tick,
    serverTimeMs,
    blocks: [],
    entities,
    circuit: { nodes: [] },
    stats: {
      clientCount: 1,
      commandQueueLength: 0,
      snapshotBytes: 128,
    },
  };
}

function entity(id: string, x: number, type = DEBUG_MOVER_ENTITY_TYPE): EntitySnapshot {
  return {
    id,
    type,
    transform: {
      position: { x, y: 1.25, z: 1 },
      rotation: { x: 0, y: 0, z: 0, w: 1 },
      scale: { x: 0.5, y: 0.5, z: 0.5 },
    },
  };
}
