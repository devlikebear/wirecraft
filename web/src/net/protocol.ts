export const BlockType = {
  Air: 0,
  Solid: 1,
  DebugMover: 2,
  Power: 3,
  Wire: 4,
  Button: 5,
  AndGate: 6,
  MCUOutput: 7,
  Piston: 8,
  Motor: 9,
  MotorDriver: 10,
  TransistorSwitch: 11,
} as const;

export type BlockType = (typeof BlockType)[keyof typeof BlockType];

export const EntityType = {
  DebugMover: 'debug_mover',
  Piston: 'piston',
  Motor: 'motor',
} as const;

export type EntityType = (typeof EntityType)[keyof typeof EntityType];

export type CommandType = 'place_block' | 'remove_block' | 'set_button';
export type SnapshotMode = 'full' | 'changed_set';

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
  mode: SnapshotMode;
  tick: number;
  baseTick?: number;
  serverTimeMs: number;
  blocks: BlockSnapshot[];
  changedBlocks: BlockSnapshot[];
  removedBlocks: Position[];
  entities: EntitySnapshot[];
  changedEntities: EntitySnapshot[];
  circuit: CircuitSnapshot;
  presence: PresenceSnapshot;
  commandAcks: CommandAckSnapshot[];
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

export type SignalState = 'unknown' | 'low' | 'high';

export interface CircuitSnapshot {
  nodes: CircuitNodeSnapshot[];
}

export interface CircuitNodeSnapshot {
  position: Position;
  nodeId: string;
  nodeType: string;
  signalState: SignalState;
}

export interface PresenceSnapshot {
  clients: ClientPresenceSnapshot[];
}

export interface ClientPresenceSnapshot {
  id: string;
  displayName: string;
}

export type CommandAckStatus = 'accepted' | 'rejected';

export interface CommandAckSnapshot {
  clientId: string;
  commandId: string;
  status: CommandAckStatus;
  reason?: string;
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

  const mode = parseSnapshotMode(value.mode);
  if (mode === null) {
    return null;
  }
  let baseTick: number | undefined;
  if (typeof value.baseTick !== 'undefined') {
    if (typeof value.baseTick !== 'number') {
      return null;
    }
    baseTick = value.baseTick;
  }
  if (mode === 'changed_set' && typeof baseTick !== 'number') {
    return null;
  }

  const blocks = parseArray(value.blocks, parseBlockSnapshot);
  const changedBlocks = parseOptionalArray(value.changedBlocks, parseBlockSnapshot);
  const removedBlocks = parseOptionalArray(value.removedBlocks, parsePosition);
  const entities = parseArray(value.entities, parseEntitySnapshot);
  const changedEntities = parseOptionalArray(value.changedEntities, parseEntitySnapshot);
  const circuit = parseCircuitSnapshot(value.circuit);
  const presence = parsePresenceSnapshot(value.presence);
  const commandAcks = parseCommandAcks(value.commandAcks);
  const stats = parseSnapshotStats(value.stats);

  if (
    typeof value.tick !== 'number' ||
    typeof value.serverTimeMs !== 'number' ||
    blocks === null ||
    changedBlocks === null ||
    removedBlocks === null ||
    entities === null ||
    changedEntities === null ||
    circuit === null ||
    presence === null ||
    commandAcks === null ||
    stats === null
  ) {
    return null;
  }

  return {
    mode,
    tick: value.tick,
    ...(typeof baseTick === 'number' ? { baseTick } : {}),
    serverTimeMs: value.serverTimeMs,
    blocks,
    changedBlocks,
    removedBlocks,
    entities,
    changedEntities,
    circuit,
    presence,
    commandAcks,
    stats,
  };
}

function parseSnapshotMode(value: unknown): SnapshotMode | null {
  if (typeof value === 'undefined') {
    return 'full';
  }
  if (value === 'full' || value === 'changed_set') {
    return value;
  }
  return null;
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

function parseCircuitSnapshot(value: unknown): CircuitSnapshot | null {
  if (!isRecord(value)) {
    return null;
  }

  const nodes = parseArray(value.nodes, parseCircuitNodeSnapshot);
  if (nodes === null) {
    return null;
  }

  return { nodes };
}

function parseCircuitNodeSnapshot(value: unknown): CircuitNodeSnapshot | null {
  if (!isRecord(value)) {
    return null;
  }

  const position = parsePosition(value.position);
  if (
    position === null ||
    typeof value.nodeId !== 'string' ||
    typeof value.nodeType !== 'string' ||
    !isSignalState(value.signalState)
  ) {
    return null;
  }

  return {
    position,
    nodeId: value.nodeId,
    nodeType: value.nodeType,
    signalState: value.signalState,
  };
}

function parsePresenceSnapshot(value: unknown): PresenceSnapshot | null {
  if (typeof value === 'undefined') {
    return { clients: [] };
  }
  if (!isRecord(value)) {
    return null;
  }

  const clients = parseArray(value.clients, parseClientPresenceSnapshot);
  if (clients === null) {
    return null;
  }

  return { clients };
}

function parseClientPresenceSnapshot(value: unknown): ClientPresenceSnapshot | null {
  if (!isRecord(value)) {
    return null;
  }
  if (typeof value.id !== 'string' || typeof value.displayName !== 'string') {
    return null;
  }
  return { id: value.id, displayName: value.displayName };
}

function parseCommandAcks(value: unknown): CommandAckSnapshot[] | null {
  if (typeof value === 'undefined') {
    return [];
  }
  return parseArray(value, parseCommandAckSnapshot);
}

function parseCommandAckSnapshot(value: unknown): CommandAckSnapshot | null {
  if (!isRecord(value)) {
    return null;
  }
  if (
    typeof value.clientId !== 'string' ||
    typeof value.commandId !== 'string' ||
    !isCommandAckStatus(value.status)
  ) {
    return null;
  }
  if (typeof value.reason !== 'undefined' && typeof value.reason !== 'string') {
    return null;
  }

  return typeof value.reason === 'string'
    ? {
        clientId: value.clientId,
        commandId: value.commandId,
        status: value.status,
        reason: value.reason,
      }
    : {
        clientId: value.clientId,
        commandId: value.commandId,
        status: value.status,
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

function isSignalState(value: unknown): value is SignalState {
  return value === 'unknown' || value === 'low' || value === 'high';
}

function isCommandAckStatus(value: unknown): value is CommandAckStatus {
  return value === 'accepted' || value === 'rejected';
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

function parseOptionalArray<T>(value: unknown, parser: (entry: unknown) => T | null): T[] | null {
  if (typeof value === 'undefined') {
    return [];
  }
  return parseArray(value, parser);
}
