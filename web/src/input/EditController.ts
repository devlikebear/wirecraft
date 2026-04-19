import {
  Plane,
  Raycaster,
  Vector2,
  Vector3,
  type Camera,
  type Object3D,
  type WebGLRenderer,
} from 'three';
import { BlockType, type BlockFacing, type Command, type Position, type Vec3 } from '../net/protocol';
import type { VoxelRaycastHit, VoxelRenderer } from '../render/VoxelRenderer';

export type EditMode = 'place' | 'remove';
export const BLOCK_FACINGS: BlockFacing[] = ['north', 'east', 'south', 'west'];
const EDIT_CLICK_DRAG_THRESHOLD_PX = 5;

export interface EditPointerLike {
  button: number;
  shiftKey: boolean;
}

export interface PointerPointLike {
  x: number;
  y: number;
}

export interface BuildEditCommandInput {
  mode: EditMode;
  clientId: string;
  commandId: string;
  tickHint: number;
  position: Position;
  blockType: BlockType;
  facing?: BlockFacing;
}

export interface BuildSetButtonCommandInput {
  clientId: string;
  commandId: string;
  tickHint: number;
  position: Position;
  buttonPressed: boolean;
}

export interface EditControllerOptions {
  camera: Camera;
  renderer: WebGLRenderer;
  worldRoot: Object3D;
  voxelRenderer: VoxelRenderer;
  clientId: string;
  blockType?: BlockType;
  getBlockType?: () => BlockType;
  getFacing?: () => BlockFacing;
  getTickHint: () => number;
  sendCommand: (command: Command) => void;
}

export function editModeFromPointer(event: EditPointerLike): EditMode | null {
  if (event.button === 2 || (event.button === 0 && event.shiftKey)) {
    return 'remove';
  }
  if (event.button === 0) {
    return 'place';
  }
  return null;
}

export function adjacentPosition(hit: VoxelRaycastHit): Position {
  return {
    x: hit.position.x + hit.faceNormal.x,
    y: hit.position.y + hit.faceNormal.y,
    z: hit.position.z + hit.faceNormal.z,
  };
}

export function positionFromGroundPoint(point: Vec3): Position {
  return {
    x: Math.round(point.x),
    y: 0,
    z: Math.round(point.z),
  };
}

export function nextBlockFacing(facing: BlockFacing): BlockFacing {
  const index = BLOCK_FACINGS.indexOf(facing);
  return BLOCK_FACINGS[(index + 1) % BLOCK_FACINGS.length];
}

export function isEditClickWithinDragThreshold(
  start: PointerPointLike,
  end: PointerPointLike,
  thresholdPx = EDIT_CLICK_DRAG_THRESHOLD_PX,
): boolean {
  return Math.hypot(end.x - start.x, end.y - start.y) <= thresholdPx;
}

export function buildEditCommand(input: BuildEditCommandInput): Command {
  return {
    type: input.mode === 'place' ? 'place_block' : 'remove_block',
    clientId: input.clientId,
    commandId: input.commandId,
    tickHint: input.tickHint,
    position: input.position,
    blockType: input.mode === 'place' ? input.blockType : BlockType.Air,
    ...(input.mode === 'place' && input.facing ? { facing: input.facing } : {}),
  };
}

export function buildSetButtonCommand(input: BuildSetButtonCommandInput): Command {
  return {
    type: 'set_button',
    clientId: input.clientId,
    commandId: input.commandId,
    tickHint: input.tickHint,
    position: input.position,
    blockType: BlockType.Air,
    buttonPressed: input.buttonPressed,
  };
}

export class EditController {
  private readonly camera: Camera;
  private readonly renderer: WebGLRenderer;
  private readonly worldRoot: Object3D;
  private readonly voxelRenderer: VoxelRenderer;
  private readonly clientId: string;
  private readonly getBlockType: () => BlockType;
  private readonly getFacing: () => BlockFacing;
  private readonly getTickHint: () => number;
  private readonly sendCommand: (command: Command) => void;
  private readonly raycaster = new Raycaster();
  private readonly pointer = new Vector2();
  private readonly groundPlane = new Plane(new Vector3(0, 1, 0), 0);
  private readonly groundPoint = new Vector3();
  private commandSequence = 0;
  private pendingEditPointer: (EditPointerLike & PointerPointLike) | null = null;

  constructor(options: EditControllerOptions) {
    this.camera = options.camera;
    this.renderer = options.renderer;
    this.worldRoot = options.worldRoot;
    this.voxelRenderer = options.voxelRenderer;
    this.clientId = options.clientId;
    this.getBlockType = options.getBlockType ?? (() => options.blockType ?? BlockType.DebugMover);
    this.getFacing = options.getFacing ?? (() => 'north');
    this.getTickHint = options.getTickHint;
    this.sendCommand = options.sendCommand;
  }

  connect(): void {
    this.renderer.domElement.addEventListener('pointerdown', this.handlePointerDown);
    this.renderer.domElement.addEventListener('pointerup', this.handlePointerUp);
    this.renderer.domElement.addEventListener('contextmenu', this.preventContextMenu);
  }

  disconnect(): void {
    this.renderer.domElement.removeEventListener('pointerdown', this.handlePointerDown);
    this.renderer.domElement.removeEventListener('pointerup', this.handlePointerUp);
    this.renderer.domElement.removeEventListener('contextmenu', this.preventContextMenu);
  }

  private readonly preventContextMenu = (event: Event) => {
    event.preventDefault();
  };

  private readonly handlePointerDown = (event: PointerEvent) => {
    const mode = editModeFromPointer(event);
    if (mode === null) {
      this.pendingEditPointer = null;
      return;
    }

    this.pendingEditPointer = {
      button: event.button,
      shiftKey: event.shiftKey,
      x: event.clientX,
      y: event.clientY,
    };
  };

  private readonly handlePointerUp = (event: PointerEvent) => {
    const pending = this.pendingEditPointer;
    this.pendingEditPointer = null;
    if (!pending || pending.button !== event.button) {
      return;
    }
    if (!isEditClickWithinDragThreshold(pending, { x: event.clientX, y: event.clientY })) {
      return;
    }

    const mode = editModeFromPointer(pending);
    if (mode === null) {
      return;
    }
    const position = this.pickPosition(event, mode);
    if (position === null) {
      return;
    }

    this.sendCommand(
      buildEditCommand({
        mode,
        clientId: this.clientId,
        commandId: `${this.clientId}-${++this.commandSequence}`,
        tickHint: this.getTickHint(),
        position,
        blockType: this.getBlockType(),
        facing: this.getFacing(),
      }),
    );
  };

  private pickPosition(event: PointerEvent, mode: EditMode): Position | null {
    this.updatePointer(event);
    this.raycaster.setFromCamera(this.pointer, this.camera);

    const [intersection] = this.raycaster.intersectObjects(this.voxelRenderer.raycastTargets(), false);
    if (intersection) {
      const hit = this.voxelRenderer.hitFromIntersection(intersection);
      if (hit) {
        return mode === 'place' ? adjacentPosition(hit) : hit.position;
      }
    }

    if (mode === 'remove') {
      return null;
    }

    if (!this.raycaster.ray.intersectPlane(this.groundPlane, this.groundPoint)) {
      return null;
    }

    const localPoint = this.worldRoot.worldToLocal(this.groundPoint.clone());
    return positionFromGroundPoint(localPoint);
  }

  private updatePointer(event: PointerEvent): void {
    const rect = this.renderer.domElement.getBoundingClientRect();
    this.pointer.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
    this.pointer.y = -(((event.clientY - rect.top) / rect.height) * 2 - 1);
  }
}
