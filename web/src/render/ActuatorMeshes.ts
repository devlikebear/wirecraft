import {
  BoxGeometry,
  Mesh,
  MeshStandardMaterial,
} from 'three';
import { EntityType } from '../net/protocol';

export function isRenderableEntityType(type: string): boolean {
  return type === EntityType.Piston || type === EntityType.Motor;
}

export function createEntityMesh(type: string): Mesh {
  switch (type) {
    case EntityType.Piston:
      return createMesh(
        new BoxGeometry(0.72, 0.52, 0.72),
        new MeshStandardMaterial({
          color: 0xb8c4d6,
          emissive: 0x1f344c,
          emissiveIntensity: 0.25,
          roughness: 0.32,
          metalness: 0.35,
        }),
      );
    case EntityType.Motor:
      return createMesh(
        new BoxGeometry(0.82, 0.82, 0.82),
        new MeshStandardMaterial({
          color: 0x45b7a8,
          emissive: 0x0f3d38,
          emissiveIntensity: 0.28,
          roughness: 0.36,
          metalness: 0.24,
        }),
      );
    default:
      return createMesh(
        new BoxGeometry(1, 1, 1),
        new MeshStandardMaterial({
          color: 0x58a6ff,
          emissive: 0x0c2d57,
          emissiveIntensity: 0.45,
          roughness: 0.35,
          metalness: 0.12,
        }),
      );
  }
}

export function disposeEntityMesh(mesh: Mesh): void {
  mesh.geometry.dispose();
  if (Array.isArray(mesh.material)) {
    for (const material of mesh.material) {
      material.dispose();
    }
    return;
  }
  mesh.material.dispose();
}

function createMesh(geometry: BoxGeometry, material: MeshStandardMaterial): Mesh {
  return new Mesh(geometry, material);
}
