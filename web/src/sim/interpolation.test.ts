import { describe, expect, it } from 'vitest';
import type { Quat, Snapshot, TransformSnapshot, Vec3 } from '../net/protocol';
import {
  calculateAlpha,
  findSnapshotPair,
  interpolateTransform,
  interpolateVec3,
  slerpQuat,
} from './interpolation';

describe('findSnapshotPair', () => {
  it('finds snapshots around the render server time', () => {
    const pair = findSnapshotPair([snapshot(1000), snapshot(1050), snapshot(1100)], 1075);

    expect(pair.before?.serverTimeMs).toBe(1050);
    expect(pair.after?.serverTimeMs).toBe(1100);
    expect(pair.alpha).toBe(0.5);
    expect(pair.status).toBe('between');
  });

  it('returns the first snapshot when render time is too early', () => {
    const pair = findSnapshotPair([snapshot(1000), snapshot(1050)], 900);

    expect(pair.before).toBeNull();
    expect(pair.after?.serverTimeMs).toBe(1000);
    expect(pair.alpha).toBe(0);
    expect(pair.status).toBe('before-first');
  });

  it('returns the last snapshot when render time is too late', () => {
    const pair = findSnapshotPair([snapshot(1000), snapshot(1050)], 1200);

    expect(pair.before?.serverTimeMs).toBe(1050);
    expect(pair.after).toBeNull();
    expect(pair.alpha).toBe(0);
    expect(pair.status).toBe('after-last');
  });

  it('handles empty buffers', () => {
    const pair = findSnapshotPair([], 1000);

    expect(pair.before).toBeNull();
    expect(pair.after).toBeNull();
    expect(pair.alpha).toBe(0);
    expect(pair.status).toBe('empty');
  });
});

describe('calculateAlpha', () => {
  it('clamps alpha to the snapshot range', () => {
    expect(calculateAlpha(1000, 1100, 950)).toBe(0);
    expect(calculateAlpha(1000, 1100, 1050)).toBe(0.5);
    expect(calculateAlpha(1000, 1100, 1200)).toBe(1);
  });

  it('returns zero when timestamps collapse', () => {
    expect(calculateAlpha(1000, 1000, 1000)).toBe(0);
  });
});

describe('transform interpolation', () => {
  it('linearly interpolates vec3 values', () => {
    expect(interpolateVec3({ x: 0, y: 2, z: 4 }, { x: 10, y: 4, z: 8 }, 0.25)).toEqual({
      x: 2.5,
      y: 2.5,
      z: 5,
    });
  });

  it('slerps quaternion-compatible values and normalizes the result', () => {
    const got = slerpQuat(
      { x: 0, y: 0, z: 0, w: 1 },
      { x: 0, y: 1, z: 0, w: 0 },
      0.5,
    );

    expect(got.x).toBeCloseTo(0);
    expect(got.y).toBeCloseTo(Math.SQRT1_2);
    expect(got.z).toBeCloseTo(0);
    expect(got.w).toBeCloseTo(Math.SQRT1_2);
  });

  it('interpolates complete transforms', () => {
    const got = interpolateTransform(transform({ x: 0, y: 0, z: 0 }), transform({ x: 10, y: 4, z: 2 }), 0.5);

    expect(got.position).toEqual({ x: 5, y: 2, z: 1 });
    expect(got.scale).toEqual({ x: 1, y: 1, z: 1 });
    expect(got.rotation.w).toBeCloseTo(1);
  });
});

function snapshot(serverTimeMs: number): Snapshot {
  return {
    tick: serverTimeMs / 50,
    serverTimeMs,
    blocks: [],
    entities: [],
    circuit: { nodes: [] },
    stats: {
      clientCount: 1,
      commandQueueLength: 0,
      snapshotBytes: 128,
    },
  };
}

function transform(position: Vec3, rotation: Quat = { x: 0, y: 0, z: 0, w: 1 }): TransformSnapshot {
  return {
    position,
    rotation,
    scale: { x: 1, y: 1, z: 1 },
  };
}
