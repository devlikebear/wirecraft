import {
  type BufferGeometry,
  BoxGeometry,
  CylinderGeometry,
  Group,
  InstancedMesh,
  Matrix4,
  MeshStandardMaterial,
  Vector3,
  type Intersection,
  type Object3D,
} from 'three';
import { BlockType, type BlockFacing, type BlockSnapshot, type Position, type Snapshot } from '../net/protocol';

export interface VoxelRenderItem {
  key: string;
  blockType: Exclude<BlockType, typeof BlockType.Air>;
  facing?: BlockFacing;
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

export type VoxelGeometryProfile =
  | {
      geometry: 'box';
      width: number;
      height: number;
      depth: number;
    }
  | {
      geometry: 'cylinder';
      radiusTop: number;
      radiusBottom: number;
      height: number;
      radialSegments: number;
    };

export interface VoxelRaycastHit {
  position: Position;
  faceNormal: Position;
}

const renderableBlockTypes = [
  BlockType.Solid,
  BlockType.DebugMover,
  BlockType.Power,
  BlockType.Wire,
  BlockType.Button,
  BlockType.AndGate,
  BlockType.MCUOutput,
] as const;

const defaultVoxelProfile: VoxelGeometryProfile = {
  geometry: 'box',
  width: 1,
  height: 1,
  depth: 1,
};

const visualProfiles = new Map<BlockType, VoxelGeometryProfile>([
  [BlockType.Solid, defaultVoxelProfile],
  [BlockType.DebugMover, defaultVoxelProfile],
  [
    BlockType.Power,
    {
      geometry: 'cylinder',
      radiusTop: 0.26,
      radiusBottom: 0.42,
      height: 0.72,
      radialSegments: 18,
    },
  ],
  [
    BlockType.Wire,
    {
      geometry: 'box',
      width: 0.9,
      height: 0.14,
      depth: 0.32,
    },
  ],
  [
    BlockType.Button,
    {
      geometry: 'cylinder',
      radiusTop: 0.34,
      radiusBottom: 0.42,
      height: 0.22,
      radialSegments: 18,
    },
  ],
  [
    BlockType.AndGate,
    {
      geometry: 'box',
      width: 0.9,
      height: 0.46,
      depth: 0.62,
    },
  ],
  [
    BlockType.MCUOutput,
    {
      geometry: 'cylinder',
      radiusTop: 0.36,
      radiusBottom: 0.36,
      height: 0.56,
      radialSegments: 18,
    },
  ],
]);

export function voxelVisualProfileForBlockType(blockType: BlockType): VoxelGeometryProfile {
  return visualProfiles.get(blockType) ?? defaultVoxelProfile;
}

export function blockFacingYaw(facing: BlockFacing | undefined): number {
  switch (facing) {
    case 'north':
      return -Math.PI / 2;
    case 'south':
      return Math.PI / 2;
    case 'west':
      return Math.PI;
    case 'east':
    default:
      return 0;
  }
}

export function createVoxelRenderItems(blocks: BlockSnapshot[]): VoxelRenderItem[] {
  const items: VoxelRenderItem[] = [];

  for (const block of blocks) {
    if (block.blockType === BlockType.Air) {
      continue;
    }

    const blockType = block.blockType as VoxelRenderItem['blockType'];
    const profile = voxelVisualProfileForBlockType(blockType);
    const facingKey = block.facing ? `:${block.facing}` : '';
    items.push({
      key: `${block.position.x}:${block.position.y}:${block.position.z}:${block.blockType}${facingKey}`,
      blockType,
      ...(block.facing ? { facing: block.facing } : {}),
      blockPosition: block.position,
      position: {
        x: block.position.x,
        y: block.position.y + profile.height / 2,
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

  private readonly geometries = new Map<VoxelRenderItem['blockType'], BufferGeometry>(
    renderableBlockTypes.map((blockType) => [
      blockType,
      createGeometry(voxelVisualProfileForBlockType(blockType)),
    ]),
  );
  private readonly transform = new Matrix4();
  private readonly meshes = new Map<VoxelRenderItem['blockType'], ManagedMesh>();
  private readonly meshItems = new Map<InstancedMesh, VoxelRenderItem[]>();
  private readonly materials = new Map<VoxelRenderItem['blockType'], MeshStandardMaterial>([
    [BlockType.Solid, new MeshStandardMaterial({ color: 0x6d8f7d, roughness: 0.7 })],
    [
      BlockType.DebugMover,
      new MeshStandardMaterial({ color: 0xd6c064, roughness: 0.45, metalness: 0.08 }),
    ],
    [
      BlockType.Power,
      new MeshStandardMaterial({
        color: 0xf05a4f,
        emissive: 0x7f1f18,
        emissiveIntensity: 0.35,
        roughness: 0.42,
      }),
    ],
    [
      BlockType.Wire,
      new MeshStandardMaterial({ color: 0xc9824a, roughness: 0.58, metalness: 0.12 }),
    ],
    [BlockType.Button, new MeshStandardMaterial({ color: 0x4f8df0, roughness: 0.5 })],
    [
      BlockType.AndGate,
      new MeshStandardMaterial({ color: 0x7d6cf2, roughness: 0.45, metalness: 0.04 }),
    ],
    [
      BlockType.MCUOutput,
      new MeshStandardMaterial({
        color: 0x35b58a,
        emissive: 0x0c5f43,
        emissiveIntensity: 0.26,
        roughness: 0.48,
      }),
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
        const { facing, position } = items[index];
        this.transform.makeRotationY(blockFacingYaw(facing));
        this.transform.setPosition(position.x, position.y, position.z);
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
    for (const geometry of this.geometries.values()) {
      geometry.dispose();
    }
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
    const geometry = this.geometries.get(blockType);
    if (!geometry) {
      throw new Error(`Missing geometry for block type ${blockType}`);
    }

    const mesh = new InstancedMesh(geometry, material, requiredCapacity);
    mesh.name = `wirecraft-voxels-${blockType}`;
    mesh.count = 0;
    this.object.add(mesh);

    const managed = { mesh, capacity: requiredCapacity };
    this.meshes.set(blockType, managed);
    return managed;
  }
}

function createGeometry(profile: VoxelGeometryProfile): BufferGeometry {
  switch (profile.geometry) {
    case 'box':
      return new BoxGeometry(profile.width, profile.height, profile.depth);
    case 'cylinder':
      return new CylinderGeometry(
        profile.radiusTop,
        profile.radiusBottom,
        profile.height,
        profile.radialSegments,
      );
  }
}
