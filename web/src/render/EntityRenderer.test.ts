import { describe, expect, it } from 'vitest';
import { Group } from 'three';
import { EntityType, type EntitySnapshot, type Snapshot } from '../net/protocol';
import {
  DEBUG_MOVER_ENTITY_TYPE,
  EntityRenderer,
  selectInterpolatedEntityRenderItems,
} from './EntityRenderer';

describe('selectInterpolatedEntityRenderItems', () => {
  it('hides debug mover entities from the default runtime view', () => {
    const items = selectInterpolatedEntityRenderItems(
      [snapshot(1, 1000, entity('debug-mover-1', 0)), snapshot(2, 1100, entity('debug-mover-1', 10))],
      1050,
    );

    expect(items).toEqual([]);
  });

  it('falls back to the nearest available snapshot when outside the buffered range', () => {
    const items = selectInterpolatedEntityRenderItems(
      [snapshot(1, 1000, entity('piston:2:0:0', 3, EntityType.Piston))],
      1200,
    );

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

  it('interpolates actuator entities between buffered snapshots', () => {
    const items = selectInterpolatedEntityRenderItems(
      [
        snapshot(1, 1000, entity('piston:2:0:0', 2, EntityType.Piston)),
        snapshot(2, 1100, entity('piston:2:0:0', 4, EntityType.Piston)),
      ],
      1050,
    );

    expect(items).toHaveLength(1);
    expect(items[0].id).toBe('piston:2:0:0');
    expect(items[0].type).toBe(EntityType.Piston);
    expect(items[0].transform.position).toEqual({ x: 3, y: 1.25, z: 1 });
  });
});

describe('EntityRenderer', () => {
  it('does not create debug mover meshes by default', () => {
    const parent = new Group();
    const renderer = new EntityRenderer(parent);

    renderer.updateFromSnapshots(
      [snapshot(1, 1000, entity('debug-mover-1', 0)), snapshot(2, 1100, entity('debug-mover-1', 10))],
      1050,
    );

    expect(parent.children).toContain(renderer.object);
    expect(renderer.count).toBe(0);
    expect(renderer.object.getObjectByName('wirecraft-entity-debug-mover-1')).toBeUndefined();

    renderer.dispose();
    expect(parent.children).not.toContain(renderer.object);
  });

  it('creates actuator meshes by entity type and ID', () => {
    const parent = new Group();
    const renderer = new EntityRenderer(parent);

    renderer.updateFromSnapshots(
      [
        snapshot(
          1,
          1000,
          entity('debug-mover-1', 0),
          entity('piston:2:0:0', 2, EntityType.Piston),
          entity('motor:4:0:0', 4, EntityType.Motor),
          entity('servo-1', 6, 'servo'),
        ),
      ],
      1000,
    );

    expect(renderer.count).toBe(2);
    expect(renderer.object.getObjectByName('wirecraft-entity-debug-mover-1')).toBeUndefined();
    expect(renderer.object.getObjectByName('wirecraft-entity-piston:2:0:0')).toBeTruthy();
    expect(renderer.object.getObjectByName('wirecraft-entity-motor:4:0:0')).toBeTruthy();
    expect(renderer.object.getObjectByName('wirecraft-entity-servo-1')).toBeUndefined();

    renderer.dispose();
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
