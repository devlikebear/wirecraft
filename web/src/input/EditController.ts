import {
  Plane,
  Raycaster,
  Vector2,
  Vector3,
  type Camera,
  type Object3D,
  type WebGLRenderer,
} from 'three';
import { BlockType, type Command, type Position, type Vec3 } from '../net/protocol';
import type { VoxelRaycastHit, VoxelRenderer } from '../render/VoxelRenderer';

export type EditMode = 'place' | 'remove';

export interface EditPointerLike {
  button: number;
  shiftKey: boolean;
}

export interface BuildEditCommandInput {
  mode: EditMode;
  clientId: string;
  commandId: string;
  tickHint: number;
  position: Position;
  blockType: BlockType;
}

export interface EditControllerOptions {
  camera: Camera;
  renderer: WebGLRenderer;
  worldRoot: Object3D;
  voxelRenderer: VoxelRenderer;
  clientId: string;
  blockType?: BlockType;
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

export function buildEditCommand(input: BuildEditCommandInput): Command {
  return {
    type: input.mode === 'place' ? 'place_block' : 'remove_block',
    clientId: input.clientId,
    commandId: input.commandId,
    tickHint: input.tickHint,
    position: input.position,
    blockType: input.mode === 'place' ? input.blockType : BlockType.Air,
  };
}

export class EditController {
  private readonly camera: Camera;
  private readonly renderer: WebGLRenderer;
  private readonly worldRoot: Object3D;
  private readonly voxelRenderer: VoxelRenderer;
  private readonly clientId: string;
  private readonly blockType: BlockType;
  private readonly getTickHint: () => number;
  private readonly sendCommand: (command: Command) => void;
  private readonly raycaster = new Raycaster();
  private readonly pointer = new Vector2();
  private readonly groundPlane = new Plane(new Vector3(0, 1, 0), 0);
  private readonly groundPoint = new Vector3();
  private commandSequence = 0;

  constructor(options: EditControllerOptions) {
    this.camera = options.camera;
    this.renderer = options.renderer;
    this.worldRoot = options.worldRoot;
    this.voxelRenderer = options.voxelRenderer;
    this.clientId = options.clientId;
    this.blockType = options.blockType ?? BlockType.DebugMover;
    this.getTickHint = options.getTickHint;
    this.sendCommand = options.sendCommand;
  }

  connect(): void {
    this.renderer.domElement.addEventListener('pointerdown', this.handlePointerDown);
    this.renderer.domElement.addEventListener('contextmenu', this.preventContextMenu);
  }

  disconnect(): void {
    this.renderer.domElement.removeEventListener('pointerdown', this.handlePointerDown);
    this.renderer.domElement.removeEventListener('contextmenu', this.preventContextMenu);
  }

  private readonly preventContextMenu = (event: Event) => {
    event.preventDefault();
  };

  private readonly handlePointerDown = (event: PointerEvent) => {
    const mode = editModeFromPointer(event);
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
        blockType: this.blockType,
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
