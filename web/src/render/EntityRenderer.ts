import {
  Group,
  Mesh,
  type Object3D,
} from 'three';
import { EntityType, type EntitySnapshot, type Snapshot, type TransformSnapshot } from '../net/protocol';
import { findSnapshotPair, interpolateTransform } from '../sim/interpolation';
import { createEntityMesh, disposeEntityMesh, isRenderableEntityType } from './ActuatorMeshes';

export const DEBUG_MOVER_ENTITY_TYPE = EntityType.DebugMover;

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
        disposeEntityMesh(mesh);
        this.meshes.delete(id);
      }
    }
  }

  dispose(): void {
    for (const mesh of this.meshes.values()) {
      this.object.remove(mesh);
      disposeEntityMesh(mesh);
    }
    this.meshes.clear();
    this.parent.remove(this.object);
  }

  private ensureMesh(item: EntityRenderItem): Mesh {
    const existing = this.meshes.get(item.id);
    if (existing && existing.userData.entityType === item.type) {
      return existing;
    }
    if (existing) {
      this.object.remove(existing);
      disposeEntityMesh(existing);
      this.meshes.delete(item.id);
    }

    const mesh = createEntityMesh(item.type);
    mesh.name = `wirecraft-entity-${item.id}`;
    mesh.userData.entityType = item.type;
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
  return entities.filter((entity) => isRenderableEntityType(entity.type));
}

function toRenderItem(entity: EntitySnapshot): EntityRenderItem {
  return {
    id: entity.id,
    type: entity.type,
    transform: entity.transform,
  };
}
