import { describe, expect, it } from 'vitest';
import { BlockType, type BlockSnapshot } from '../net/protocol';
import { createVoxelRenderItems, groupVoxelRenderItems } from './VoxelRenderer';

describe('createVoxelRenderItems', () => {
  it('maps snapshot blocks to render positions', () => {
    const blocks: BlockSnapshot[] = [
      { position: { x: 4, y: 2, z: 6 }, blockType: BlockType.Solid },
      { position: { x: 1, y: 0, z: 3 }, blockType: BlockType.DebugMover },
      { position: { x: 2, y: 0, z: 3 }, blockType: BlockType.Power },
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
        key: '2:0:3:3',
        blockType: BlockType.Power,
        blockPosition: { x: 2, y: 0, z: 3 },
        position: { x: 2, y: 0.5, z: 3 },
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
