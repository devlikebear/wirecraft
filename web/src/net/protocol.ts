export const BlockType = {
  Air: 0,
  Solid: 1,
  DebugMover: 2,
  Power: 3,
  Wire: 4,
  Button: 5,
  AndGate: 6,
  MCUOutput: 7,
} as const;

export type BlockType = (typeof BlockType)[keyof typeof BlockType];

export type CommandType = 'place_block' | 'remove_block' | 'set_button';

export interface Position {
  x: number;
  y: number;
  z: number;
}

export interface Command {
  type: CommandType;
  clientId: string;
  commandId: string;
  tickHint: number;
  position: Position;
  blockType: BlockType;
  buttonPressed?: boolean;
}

export interface Snapshot {
  tick: number;
  serverTimeMs: number;
  blocks: BlockSnapshot[];
  entities: EntitySnapshot[];
  stats: SnapshotStats;
}

export interface BlockSnapshot {
  position: Position;
  blockType: BlockType;
}

export interface EntitySnapshot {
  id: string;
  type: string;
  transform: TransformSnapshot;
}

export interface TransformSnapshot {
  position: Vec3;
  rotation: Quat;
  scale: Vec3;
}

export interface Vec3 {
  x: number;
  y: number;
  z: number;
}

export interface Quat {
  x: number;
  y: number;
  z: number;
  w: number;
}

export interface SnapshotStats {
  clientCount: number;
  commandQueueLength: number;
  snapshotBytes: number;
}

export function isSnapshot(value: unknown): value is Snapshot {
  return parseSnapshot(value) !== null;
}

export function parseSnapshot(value: unknown): Snapshot | null {
  if (!isRecord(value)) {
    return null;
  }

  const blocks = parseArray(value.blocks, parseBlockSnapshot);
  const entities = parseArray(value.entities, parseEntitySnapshot);
  const stats = parseSnapshotStats(value.stats);

  if (
    typeof value.tick !== 'number' ||
    typeof value.serverTimeMs !== 'number' ||
    blocks === null ||
    entities === null ||
    stats === null
  ) {
    return null;
  }

  return {
    tick: value.tick,
    serverTimeMs: value.serverTimeMs,
    blocks,
    entities,
    stats,
  };
}

function parseBlockSnapshot(value: unknown): BlockSnapshot | null {
  if (!isRecord(value)) {
    return null;
  }

  const position = parsePosition(value.position);
  if (position === null || !isBlockType(value.blockType)) {
    return null;
  }

  return {
    position,
    blockType: value.blockType,
  };
}

function parseEntitySnapshot(value: unknown): EntitySnapshot | null {
  if (!isRecord(value)) {
    return null;
  }

  const transform = parseTransformSnapshot(value.transform);
  if (typeof value.id !== 'string' || typeof value.type !== 'string' || transform === null) {
    return null;
  }

  return {
    id: value.id,
    type: value.type,
    transform,
  };
}

function parseTransformSnapshot(value: unknown): TransformSnapshot | null {
  if (!isRecord(value)) {
    return null;
  }

  const position = parseVec3(value.position);
  const rotation = parseQuat(value.rotation);
  const scale = parseVec3(value.scale);
  if (position === null || rotation === null || scale === null) {
    return null;
  }

  return { position, rotation, scale };
}

function parsePosition(value: unknown): Position | null {
  if (!isRecord(value)) {
    return null;
  }

  if (hasNumberFields(value, 'x', 'y', 'z')) {
    return { x: value.x, y: value.y, z: value.z };
  }
  if (hasNumberFields(value, 'X', 'Y', 'Z')) {
    return { x: value.X, y: value.Y, z: value.Z };
  }

  return null;
}

function parseVec3(value: unknown): Vec3 | null {
  if (!isRecord(value)) {
    return null;
  }

  if (!hasNumberFields(value, 'x', 'y', 'z')) {
    return null;
  }

  return { x: value.x, y: value.y, z: value.z };
}

function parseQuat(value: unknown): Quat | null {
  if (!isRecord(value)) {
    return null;
  }

  if (!hasNumberFields(value, 'x', 'y', 'z', 'w')) {
    return null;
  }

  return { x: value.x, y: value.y, z: value.z, w: value.w };
}

function parseSnapshotStats(value: unknown): SnapshotStats | null {
  if (!isRecord(value)) {
    return null;
  }

  if (
    typeof value.clientCount === 'number' &&
    typeof value.commandQueueLength === 'number' &&
    typeof value.snapshotBytes === 'number'
  ) {
    return {
      clientCount: value.clientCount,
      commandQueueLength: value.commandQueueLength,
      snapshotBytes: value.snapshotBytes,
    };
  }

  return null;
}

function isBlockType(value: unknown): value is BlockType {
  return Object.values(BlockType).includes(value as BlockType);
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null;
}

function hasNumberFields<T extends string>(
  value: Record<string, unknown>,
  ...fields: T[]
): value is Record<T, number> {
  return fields.every((field) => typeof value[field] === 'number');
}

function parseArray<T>(value: unknown, parser: (entry: unknown) => T | null): T[] | null {
  if (!Array.isArray(value)) {
    return null;
  }

  const parsed: T[] = [];
  for (const entry of value) {
    const item = parser(entry);
    if (item === null) {
      return null;
    }
    parsed.push(item);
  }

  return parsed;
}
