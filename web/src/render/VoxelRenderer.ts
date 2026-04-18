import {
  BoxGeometry,
  Group,
  InstancedMesh,
  Matrix4,
  MeshStandardMaterial,
  Vector3,
  type Intersection,
  type Object3D,
} from 'three';
import { BlockType, type BlockSnapshot, type Position, type Snapshot } from '../net/protocol';

export interface VoxelRenderItem {
  key: string;
  blockType: Exclude<BlockType, typeof BlockType.Air>;
  blockPosition: Position;
  position: {
    x: number;
    y: number;
    z: number;
  };
}

interface ManagedMesh {
  mesh: InstancedMesh;
  capacity: number;
}

export interface VoxelRaycastHit {
  position: Position;
  faceNormal: Position;
}

const renderableBlockTypes = [BlockType.Solid, BlockType.DebugMover] as const;

export function createVoxelRenderItems(blocks: BlockSnapshot[]): VoxelRenderItem[] {
  const items: VoxelRenderItem[] = [];

  for (const block of blocks) {
    if (block.blockType === BlockType.Air) {
      continue;
    }

    items.push({
      key: `${block.position.x}:${block.position.y}:${block.position.z}:${block.blockType}`,
      blockType: block.blockType,
      blockPosition: block.position,
      position: {
        x: block.position.x,
        y: block.position.y + 0.5,
        z: block.position.z,
      },
    });
  }

  return items;
}

export function groupVoxelRenderItems(items: VoxelRenderItem[]): Map<VoxelRenderItem['blockType'], VoxelRenderItem[]> {
  const grouped = new Map<VoxelRenderItem['blockType'], VoxelRenderItem[]>();

  for (const item of items) {
    const bucket = grouped.get(item.blockType) ?? [];
    bucket.push(item);
    grouped.set(item.blockType, bucket);
  }

  return grouped;
}

export class VoxelRenderer {
  readonly object = new Group();

  private readonly geometry = new BoxGeometry(1, 1, 1);
  private readonly transform = new Matrix4();
  private readonly meshes = new Map<VoxelRenderItem['blockType'], ManagedMesh>();
  private readonly meshItems = new Map<InstancedMesh, VoxelRenderItem[]>();
  private readonly materials = new Map<VoxelRenderItem['blockType'], MeshStandardMaterial>([
    [BlockType.Solid, new MeshStandardMaterial({ color: 0x6d8f7d, roughness: 0.7 })],
    [
      BlockType.DebugMover,
      new MeshStandardMaterial({ color: 0xd6c064, roughness: 0.45, metalness: 0.08 }),
    ],
  ]);

  constructor(parent: Object3D) {
    this.object.name = 'wirecraft-voxels';
    parent.add(this.object);

    for (const blockType of renderableBlockTypes) {
      this.ensureMesh(blockType, 1).mesh.count = 0;
    }
  }

  update(snapshot: Snapshot): void {
    const grouped = groupVoxelRenderItems(createVoxelRenderItems(snapshot.blocks));

    for (const blockType of renderableBlockTypes) {
      const items = grouped.get(blockType) ?? [];
      const { mesh } = this.ensureMesh(blockType, items.length);

      mesh.count = items.length;
      for (let index = 0; index < items.length; index += 1) {
        const { position } = items[index];
        this.transform.makeTranslation(position.x, position.y, position.z);
        mesh.setMatrixAt(index, this.transform);
      }
      mesh.instanceMatrix.needsUpdate = true;
      this.meshItems.set(mesh, items);
    }
  }

  raycastTargets(): Object3D[] {
    return [...this.meshes.values()].map(({ mesh }) => mesh);
  }

  hitFromIntersection(intersection: Intersection<Object3D>): VoxelRaycastHit | null {
    if (!(intersection.object instanceof InstancedMesh)) {
      return null;
    }
    if (intersection.instanceId === undefined) {
      return null;
    }

    const item = this.meshItems.get(intersection.object)?.[intersection.instanceId];
    if (!item) {
      return null;
    }

    const normal = intersection.face?.normal ?? new Vector3(0, 1, 0);
    return {
      position: item.blockPosition,
      faceNormal: {
        x: Math.round(normal.x),
        y: Math.round(normal.y),
        z: Math.round(normal.z),
      },
    };
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

  private ensureMesh(blockType: VoxelRenderItem['blockType'], count: number): ManagedMesh {
    const requiredCapacity = Math.max(1, count);
    const existing = this.meshes.get(blockType);
    if (existing && existing.capacity >= requiredCapacity) {
      return existing;
    }

    if (existing) {
      this.object.remove(existing.mesh);
      this.meshItems.delete(existing.mesh);
      existing.mesh.dispose();
    }

    const material = this.materials.get(blockType);
    if (!material) {
      throw new Error(`Missing material for block type ${blockType}`);
    }

    const mesh = new InstancedMesh(this.geometry, material, requiredCapacity);
    mesh.name = `wirecraft-voxels-${blockType}`;
    mesh.count = 0;
    this.object.add(mesh);

    const managed = { mesh, capacity: requiredCapacity };
    this.meshes.set(blockType, managed);
    return managed;
  }
}
