import {
  Group,
  InstancedMesh,
  Matrix4,
  MeshStandardMaterial,
  SphereGeometry,
  type Object3D,
} from 'three';
import type { CircuitNodeSnapshot, SignalState, Snapshot, Vec3 } from '../net/protocol';

export interface CircuitOverlayItem {
  key: string;
  nodeId: string;
  nodeType: string;
  signalState: SignalState;
  position: Vec3;
}

interface ManagedMesh {
  mesh: InstancedMesh;
  capacity: number;
}

const markerOffsetY = 1.08;
const signalStates = ['high', 'low', 'unknown'] as const satisfies readonly SignalState[];

export function createCircuitOverlayItems(nodes: CircuitNodeSnapshot[]): CircuitOverlayItem[] {
  return nodes.map((node) => ({
    key: `${node.nodeId}:${node.signalState}`,
    nodeId: node.nodeId,
    nodeType: node.nodeType,
    signalState: node.signalState,
    position: {
      x: node.position.x,
      y: node.position.y + markerOffsetY,
      z: node.position.z,
    },
  }));
}

export function groupCircuitOverlayItems(
  items: CircuitOverlayItem[],
): Map<CircuitOverlayItem['signalState'], CircuitOverlayItem[]> {
  const grouped = new Map<CircuitOverlayItem['signalState'], CircuitOverlayItem[]>();

  for (const item of items) {
    const bucket = grouped.get(item.signalState) ?? [];
    bucket.push(item);
    grouped.set(item.signalState, bucket);
  }

  return grouped;
}

export class CircuitOverlay {
  readonly object = new Group();

  private readonly geometry = new SphereGeometry(0.16, 16, 8);
  private readonly transform = new Matrix4();
  private readonly meshes = new Map<CircuitOverlayItem['signalState'], ManagedMesh>();
  private readonly materials = new Map<CircuitOverlayItem['signalState'], MeshStandardMaterial>([
    [
      'high',
      new MeshStandardMaterial({
        color: 0xfff06a,
        emissive: 0xf0b629,
        emissiveIntensity: 0.85,
        roughness: 0.35,
      }),
    ],
    [
      'low',
      new MeshStandardMaterial({
        color: 0x446375,
        emissive: 0x1d3846,
        emissiveIntensity: 0.22,
        roughness: 0.65,
      }),
    ],
    [
      'unknown',
      new MeshStandardMaterial({
        color: 0x9aa3a0,
        emissive: 0x303735,
        emissiveIntensity: 0.16,
        roughness: 0.72,
      }),
    ],
  ]);

  constructor(parent: Object3D) {
    this.object.name = 'wirecraft-circuit-overlay';
    parent.add(this.object);

    for (const signalState of signalStates) {
      this.ensureMesh(signalState, 1).mesh.count = 0;
    }
  }

  update(snapshot: Snapshot): void {
    const grouped = groupCircuitOverlayItems(createCircuitOverlayItems(snapshot.circuit.nodes));

    for (const signalState of signalStates) {
      const items = grouped.get(signalState) ?? [];
      const { mesh } = this.ensureMesh(signalState, items.length);

      mesh.count = items.length;
      for (let index = 0; index < items.length; index += 1) {
        const { position } = items[index];
        this.transform.makeTranslation(position.x, position.y, position.z);
        mesh.setMatrixAt(index, this.transform);
      }
      mesh.instanceMatrix.needsUpdate = true;
    }
  }

  dispose(): void {
    for (const { mesh } of this.meshes.values()) {
      this.object.remove(mesh);
      mesh.dispose();
    }
    this.meshes.clear();
    this.geometry.dispose();
    for (const material of this.materials.values()) {
      material.dispose();
    }
  }

  private ensureMesh(signalState: CircuitOverlayItem['signalState'], count: number): ManagedMesh {
    const requiredCapacity = Math.max(1, count);
    const existing = this.meshes.get(signalState);
    if (existing && existing.capacity >= requiredCapacity) {
      return existing;
    }

    if (existing) {
      this.object.remove(existing.mesh);
      existing.mesh.dispose();
    }

    const material = this.materials.get(signalState);
    if (!material) {
      throw new Error(`Missing material for circuit state ${signalState}`);
    }

    const mesh = new InstancedMesh(this.geometry, material, requiredCapacity);
    mesh.name = `wirecraft-circuit-overlay-${signalState}`;
    mesh.count = 0;
    this.object.add(mesh);

    const managed = { mesh, capacity: requiredCapacity };
    this.meshes.set(signalState, managed);
    return managed;
  }
}
