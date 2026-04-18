import type { Quat, Snapshot, TransformSnapshot, Vec3 } from '../net/protocol';

export const DEFAULT_INTERPOLATION_DELAY_MS = 120;

export type SnapshotPairStatus = 'empty' | 'before-first' | 'after-last' | 'exact' | 'between';

export interface SnapshotPair {
  before: Snapshot | null;
  after: Snapshot | null;
  alpha: number;
  status: SnapshotPairStatus;
}

export function findSnapshotPair(buffer: Snapshot[], renderServerTimeMs: number): SnapshotPair {
  if (buffer.length === 0) {
    return { before: null, after: null, alpha: 0, status: 'empty' };
  }

  const snapshots = [...buffer].sort((a, b) => a.serverTimeMs - b.serverTimeMs);
  const first = snapshots[0];
  if (renderServerTimeMs < first.serverTimeMs) {
    return { before: null, after: first, alpha: 0, status: 'before-first' };
  }

  const last = snapshots[snapshots.length - 1];
  if (renderServerTimeMs > last.serverTimeMs) {
    return { before: last, after: null, alpha: 0, status: 'after-last' };
  }

  for (let index = 0; index < snapshots.length; index += 1) {
    const current = snapshots[index];
    if (current.serverTimeMs === renderServerTimeMs) {
      return { before: current, after: current, alpha: 0, status: 'exact' };
    }

    const next = snapshots[index + 1];
    if (!next) {
      break;
    }
    if (current.serverTimeMs < renderServerTimeMs && renderServerTimeMs < next.serverTimeMs) {
      return {
        before: current,
        after: next,
        alpha: calculateAlpha(current.serverTimeMs, next.serverTimeMs, renderServerTimeMs),
        status: 'between',
      };
    }
  }

  return { before: last, after: null, alpha: 0, status: 'after-last' };
}

export function calculateAlpha(startTimeMs: number, endTimeMs: number, renderTimeMs: number): number {
  if (endTimeMs <= startTimeMs) {
    return 0;
  }

  return clamp01((renderTimeMs - startTimeMs) / (endTimeMs - startTimeMs));
}

export function interpolateTransform(
  before: TransformSnapshot,
  after: TransformSnapshot,
  alpha: number,
): TransformSnapshot {
  return {
    position: interpolateVec3(before.position, after.position, alpha),
    rotation: slerpQuat(before.rotation, after.rotation, alpha),
    scale: interpolateVec3(before.scale, after.scale, alpha),
  };
}

export function interpolateVec3(before: Vec3, after: Vec3, alpha: number): Vec3 {
  const t = clamp01(alpha);

  return {
    x: lerp(before.x, after.x, t),
    y: lerp(before.y, after.y, t),
    z: lerp(before.z, after.z, t),
  };
}

export function slerpQuat(before: Quat, after: Quat, alpha: number): Quat {
  const t = clamp01(alpha);
  let target = normalizeQuat(after);
  const source = normalizeQuat(before);
  let dot =
    source.x * target.x + source.y * target.y + source.z * target.z + source.w * target.w;

  if (dot < 0) {
    target = {
      x: -target.x,
      y: -target.y,
      z: -target.z,
      w: -target.w,
    };
    dot = -dot;
  }

  if (dot > 0.9995) {
    return normalizeQuat({
      x: lerp(source.x, target.x, t),
      y: lerp(source.y, target.y, t),
      z: lerp(source.z, target.z, t),
      w: lerp(source.w, target.w, t),
    });
  }

  const theta0 = Math.acos(dot);
  const theta = theta0 * t;
  const sinTheta = Math.sin(theta);
  const sinTheta0 = Math.sin(theta0);
  const scaleSource = Math.cos(theta) - (dot * sinTheta) / sinTheta0;
  const scaleTarget = sinTheta / sinTheta0;

  return normalizeQuat({
    x: scaleSource * source.x + scaleTarget * target.x,
    y: scaleSource * source.y + scaleTarget * target.y,
    z: scaleSource * source.z + scaleTarget * target.z,
    w: scaleSource * source.w + scaleTarget * target.w,
  });
}

function normalizeQuat(quat: Quat): Quat {
  const length = Math.hypot(quat.x, quat.y, quat.z, quat.w);
  if (length === 0) {
    return { x: 0, y: 0, z: 0, w: 1 };
  }

  return {
    x: quat.x / length,
    y: quat.y / length,
    z: quat.z / length,
    w: quat.w / length,
  };
}

function lerp(before: number, after: number, alpha: number): number {
  return before + (after - before) * alpha;
}

function clamp01(value: number): number {
  if (value < 0) {
    return 0;
  }
  if (value > 1) {
    return 1;
  }
  return value;
}
