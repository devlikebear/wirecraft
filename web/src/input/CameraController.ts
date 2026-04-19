import {
  Vector3,
  type Camera,
  type WebGLRenderer,
} from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js';

export interface PlanarPanVector {
  x: number;
  z: number;
}

export interface KeyboardTargetLike {
  tagName?: string;
}

export interface CameraControllerOptions {
  camera: Camera;
  renderer: WebGLRenderer;
  target?: Vector3;
  panStep?: number;
}

export function cameraPanVectorForKey(key: string): PlanarPanVector | null {
  switch (key) {
    case 'w':
    case 'W':
    case 'ArrowUp':
      return { x: 0, z: -1 };
    case 's':
    case 'S':
    case 'ArrowDown':
      return { x: 0, z: 1 };
    case 'a':
    case 'A':
    case 'ArrowLeft':
      return { x: -1, z: 0 };
    case 'd':
    case 'D':
    case 'ArrowRight':
      return { x: 1, z: 0 };
    default:
      return null;
  }
}

export function shouldIgnoreKeyboardNavigation(target: KeyboardTargetLike | null): boolean {
  const tagName = target?.tagName?.toUpperCase();
  return tagName === 'INPUT' || tagName === 'TEXTAREA' || tagName === 'SELECT';
}

export class CameraController {
  readonly controls: OrbitControls;

  private readonly panStep: number;

  constructor(options: CameraControllerOptions) {
    this.controls = new OrbitControls(options.camera, options.renderer.domElement);
    this.controls.enableDamping = true;
    this.controls.dampingFactor = 0.08;
    this.controls.screenSpacePanning = false;
    this.controls.minDistance = 4;
    this.controls.maxDistance = 36;
    this.controls.maxPolarAngle = Math.PI * 0.48;
    this.controls.target.copy(options.target ?? new Vector3(0, 0, 0));
    this.panStep = options.panStep ?? 0.75;
  }

  connect(): void {
    window.addEventListener('keydown', this.handleKeyDown);
  }

  update(): void {
    this.controls.update();
  }

  dispose(): void {
    window.removeEventListener('keydown', this.handleKeyDown);
    this.controls.dispose();
  }

  private readonly handleKeyDown = (event: KeyboardEvent) => {
    if (shouldIgnoreKeyboardNavigation(event.target as KeyboardTargetLike | null)) {
      return;
    }

    const vector = cameraPanVectorForKey(event.key);
    if (vector === null) {
      return;
    }

    event.preventDefault();
    const offset = new Vector3(vector.x * this.panStep, 0, vector.z * this.panStep);
    this.controls.object.position.add(offset);
    this.controls.target.add(offset);
    this.controls.update();
  };
}
