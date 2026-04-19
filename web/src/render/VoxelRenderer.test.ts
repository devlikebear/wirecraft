import { describe, expect, it } from 'vitest';
import { BoxGeometry, CylinderGeometry, InstancedMesh, Scene } from 'three';
import { BlockType, type BlockSnapshot } from '../net/protocol';
import {
  VoxelRenderer,
  blockFacingYaw,
  createVoxelRenderItems,
  groupVoxelRenderItems,
  voxelVisualProfileForBlockType,
} from './VoxelRenderer';

describe('createVoxelRenderItems', () => {
  it('maps snapshot blocks to render positions', () => {
    const blocks: BlockSnapshot[] = [
      { position: { x: 4, y: 2, z: 6 }, blockType: BlockType.Solid },
      { position: { x: 1, y: 0, z: 3 }, blockType: BlockType.DebugMover },
      { position: { x: 2, y: 0, z: 3 }, blockType: BlockType.Power, facing: 'south' },
    ];

    expect(createVoxelRenderItems(blocks)).toEqual([
      {
        key: '4:2:6:1',
        blockType: BlockType.Solid,
        blockPosition: { x: 4, y: 2, z: 6 },
        position: { x: 4, y: 2.5, z: 6 },
      },
      {
        key: '1:0:3:2',
        blockType: BlockType.DebugMover,
        blockPosition: { x: 1, y: 0, z: 3 },
        position: { x: 1, y: 0.5, z: 3 },
      },
      {
        key: '2:0:3:3:south',
        blockType: BlockType.Power,
        facing: 'south',
        blockPosition: { x: 2, y: 0, z: 3 },
        position: { x: 2, y: 0.36, z: 3 },
      },
    ]);
  });

  it('omits air blocks', () => {
    const blocks: BlockSnapshot[] = [
      { position: { x: 0, y: 0, z: 0 }, blockType: BlockType.Air },
      { position: { x: 1, y: 0, z: 0 }, blockType: BlockType.Solid },
    ];

    expect(createVoxelRenderItems(blocks).map((item) => item.blockType)).toEqual([BlockType.Solid]);
  });
});

describe('blockFacingYaw', () => {
  it('maps cardinal block facings to horizontal yaw rotations', () => {
    expect(blockFacingYaw(undefined)).toBe(0);
    expect(blockFacingYaw('east')).toBe(0);
    expect(blockFacingYaw('south')).toBeCloseTo(Math.PI / 2);
    expect(blockFacingYaw('west')).toBeCloseTo(Math.PI);
    expect(blockFacingYaw('north')).toBeCloseTo(-Math.PI / 2);
  });
});

describe('voxelVisualProfileForBlockType', () => {
  it('uses distinct circuit block profiles instead of identical cubes', () => {
    expect(voxelVisualProfileForBlockType(BlockType.Wire)).toMatchObject({
      geometry: 'box',
      width: 0.9,
      height: 0.14,
      depth: 0.32,
    });
    expect(voxelVisualProfileForBlockType(BlockType.Power)).toMatchObject({
      geometry: 'cylinder',
      radiusTop: 0.26,
      radiusBottom: 0.42,
      height: 0.72,
    });
    expect(voxelVisualProfileForBlockType(BlockType.Button)).toMatchObject({
      geometry: 'cylinder',
      radiusTop: 0.34,
      radiusBottom: 0.42,
      height: 0.22,
    });
    expect(voxelVisualProfileForBlockType(BlockType.AndGate)).toMatchObject({
      geometry: 'box',
      width: 0.9,
      height: 0.46,
      depth: 0.62,
    });
    expect(voxelVisualProfileForBlockType(BlockType.MCUOutput)).toMatchObject({
      geometry: 'cylinder',
      radiusTop: 0.36,
      radiusBottom: 0.36,
      height: 0.56,
    });
  });
});

describe('VoxelRenderer', () => {
  it('builds profile-specific mesh geometry for circuit block types', () => {
    const scene = new Scene();
    const renderer = new VoxelRenderer(scene);

    renderer.update({
      tick: 1,
      serverTimeMs: 1000,
      blocks: [
        { position: { x: 0, y: 0, z: 0 }, blockType: BlockType.Wire },
        { position: { x: 1, y: 0, z: 0 }, blockType: BlockType.Power },
        { position: { x: 2, y: 0, z: 0 }, blockType: BlockType.Button },
        { position: { x: 3, y: 0, z: 0 }, blockType: BlockType.AndGate },
        { position: { x: 4, y: 0, z: 0 }, blockType: BlockType.MCUOutput },
      ],
      entities: [],
      circuit: { nodes: [] },
      presence: { clients: [] },
      commandAcks: [],
      stats: { clientCount: 0, commandQueueLength: 0, snapshotBytes: 0 },
    });

    const wireGeometry = meshGeometry(renderer, BlockType.Wire);
    const powerGeometry = meshGeometry(renderer, BlockType.Power);
    const buttonGeometry = meshGeometry(renderer, BlockType.Button);
    const andGateGeometry = meshGeometry(renderer, BlockType.AndGate);
    const outputGeometry = meshGeometry(renderer, BlockType.MCUOutput);

    expect(wireGeometry).toBeInstanceOf(BoxGeometry);
    expect((wireGeometry as BoxGeometry).parameters.height).toBeCloseTo(0.14);
    expect(powerGeometry).toBeInstanceOf(CylinderGeometry);
    expect(buttonGeometry).toBeInstanceOf(CylinderGeometry);
    expect(andGateGeometry).toBeInstanceOf(BoxGeometry);
    expect((andGateGeometry as BoxGeometry).parameters.depth).toBeCloseTo(0.62);
    expect(outputGeometry).toBeInstanceOf(CylinderGeometry);

    renderer.dispose();
  });
});

function meshGeometry(renderer: VoxelRenderer, blockType: BlockType) {
  const mesh = renderer.object.getObjectByName(`wirecraft-voxels-${blockType}`);
  expect(mesh).toBeInstanceOf(InstancedMesh);
  return (mesh as InstancedMesh).geometry;
}

describe('groupVoxelRenderItems', () => {
  it('groups render items by block type', () => {
    const items = createVoxelRenderItems([
      { position: { x: 0, y: 0, z: 0 }, blockType: BlockType.Solid },
      { position: { x: 1, y: 0, z: 0 }, blockType: BlockType.DebugMover },
      { position: { x: 2, y: 0, z: 0 }, blockType: BlockType.Solid },
      { position: { x: 3, y: 0, z: 0 }, blockType: BlockType.Wire },
    ]);

    const grouped = groupVoxelRenderItems(items);

    expect(grouped.get(BlockType.Solid)?.map((item) => item.position.x)).toEqual([0, 2]);
    expect(grouped.get(BlockType.DebugMover)?.map((item) => item.position.x)).toEqual([1]);
    expect(grouped.get(BlockType.Wire)?.map((item) => item.position.x)).toEqual([3]);
  });
});
