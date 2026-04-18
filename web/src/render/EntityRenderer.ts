import {
  BoxGeometry,
  Group,
  Mesh,
  MeshStandardMaterial,
  type Object3D,
} from 'three';
import type { EntitySnapshot, Snapshot, TransformSnapshot } from '../net/protocol';
import { findSnapshotPair, interpolateTransform } from '../sim/interpolation';

export const DEBUG_MOVER_ENTITY_TYPE = 'debug_mover';

export interface EntityRenderItem {
  id: string;
  type: string;
  transform: TransformSnapshot;
}

export function selectInterpolatedEntityRenderItems(
  snapshots: Snapshot[],
  renderServerTimeMs: number,
): EntityRenderItem[] {
  const pair = findSnapshotPair(snapshots, renderServerTimeMs);
  const beforeEntities = renderableEntities(pair.before?.entities ?? []);
  const afterEntities = renderableEntities(pair.after?.entities ?? []);

  if (pair.before && pair.after && pair.before !== pair.after) {
    const afterByID = new Map(afterEntities.map((entity) => [entity.id, entity]));
    return beforeEntities.map((before) => {
      const after = afterByID.get(before.id);
      return {
        id: before.id,
        type: before.type,
        transform:
          after && after.type === before.type
            ? interpolateTransform(before.transform, after.transform, pair.alpha)
            : before.transform,
      };
    });
  }

  return beforeEntities.length > 0
    ? beforeEntities.map(toRenderItem)
    : afterEntities.map(toRenderItem);
}

export class EntityRenderer {
  readonly object = new Group();

  private readonly geometry = new BoxGeometry(1, 1, 1);
  private readonly material = new MeshStandardMaterial({
    color: 0x58a6ff,
    emissive: 0x0c2d57,
    emissiveIntensity: 0.45,
    roughness: 0.35,
    metalness: 0.12,
  });
  private readonly meshes = new Map<string, Mesh>();

  constructor(private readonly parent: Object3D) {
    this.object.name = 'wirecraft-entities';
    parent.add(this.object);
  }

  get count(): number {
    return this.meshes.size;
  }

  updateFromSnapshots(snapshots: Snapshot[], renderServerTimeMs: number): void {
    const items = selectInterpolatedEntityRenderItems(snapshots, renderServerTimeMs);
    const activeIDs = new Set<string>();

    for (const item of items) {
      activeIDs.add(item.id);
      this.applyTransform(this.ensureMesh(item), item.transform);
    }

    for (const [id, mesh] of this.meshes) {
      if (!activeIDs.has(id)) {
        this.object.remove(mesh);
        mesh.geometry.dispose();
        this.meshes.delete(id);
      }
    }
  }

  dispose(): void {
    for (const mesh of this.meshes.values()) {
      this.object.remove(mesh);
      mesh.geometry.dispose();
    }
    this.meshes.clear();
    this.parent.remove(this.object);
    this.geometry.dispose();
    this.material.dispose();
  }

  private ensureMesh(item: EntityRenderItem): Mesh {
    const existing = this.meshes.get(item.id);
    if (existing) {
      return existing;
    }

    const mesh = new Mesh(this.geometry.clone(), this.material);
    mesh.name = `wirecraft-entity-${item.id}`;
    this.object.add(mesh);
    this.meshes.set(item.id, mesh);
    return mesh;
  }

  private applyTransform(mesh: Mesh, transform: TransformSnapshot): void {
    mesh.position.set(transform.position.x, transform.position.y, transform.position.z);
    mesh.quaternion.set(
      transform.rotation.x,
      transform.rotation.y,
      transform.rotation.z,
      transform.rotation.w,
    );
    mesh.scale.set(transform.scale.x, transform.scale.y, transform.scale.z);
  }
}

function renderableEntities(entities: EntitySnapshot[]): EntitySnapshot[] {
  return entities.filter((entity) => entity.type === DEBUG_MOVER_ENTITY_TYPE);
}

function toRenderItem(entity: EntitySnapshot): EntityRenderItem {
  return {
    id: entity.id,
    type: entity.type,
    transform: entity.transform,
  };
}
