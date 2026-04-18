import { describe, expect, it } from 'vitest';
import { BlockType } from '../net/protocol';
import {
  adjacentPosition,
  buildSetButtonCommand,
  buildEditCommand,
  editModeFromPointer,
  positionFromGroundPoint,
} from './EditController';

describe('editModeFromPointer', () => {
  it('uses left click for place', () => {
    expect(editModeFromPointer({ button: 0, shiftKey: false })).toBe('place');
  });

  it('uses shift left click and right click for remove', () => {
    expect(editModeFromPointer({ button: 0, shiftKey: true })).toBe('remove');
    expect(editModeFromPointer({ button: 2, shiftKey: false })).toBe('remove');
  });

  it('ignores non-edit pointer buttons', () => {
    expect(editModeFromPointer({ button: 1, shiftKey: false })).toBeNull();
  });
});

describe('position helpers', () => {
  it('places adjacent to the clicked face', () => {
    expect(
      adjacentPosition({
        position: { x: 3, y: 4, z: 5 },
        faceNormal: { x: 0, y: 1, z: 0 },
      }),
    ).toEqual({ x: 3, y: 5, z: 5 });
  });

  it('rounds ground hits to block coordinates', () => {
    expect(positionFromGroundPoint({ x: 2.4, y: 0, z: 3.6 })).toEqual({ x: 2, y: 0, z: 4 });
  });
});

describe('buildEditCommand', () => {
  it('builds place commands without mutating local world state', () => {
    expect(
      buildEditCommand({
        mode: 'place',
        clientId: 'client-1',
        commandId: 'cmd-1',
        tickHint: 42,
        position: { x: 1, y: 2, z: 3 },
        blockType: BlockType.DebugMover,
      }),
    ).toEqual({
      type: 'place_block',
      clientId: 'client-1',
      commandId: 'cmd-1',
      tickHint: 42,
      position: { x: 1, y: 2, z: 3 },
      blockType: BlockType.DebugMover,
    });
  });

  it('builds place commands with the selected circuit block type', () => {
    expect(
      buildEditCommand({
        mode: 'place',
        clientId: 'client-1',
        commandId: 'cmd-3',
        tickHint: 44,
        position: { x: 5, y: 0, z: 2 },
        blockType: BlockType.AndGate,
      }),
    ).toEqual({
      type: 'place_block',
      clientId: 'client-1',
      commandId: 'cmd-3',
      tickHint: 44,
      position: { x: 5, y: 0, z: 2 },
      blockType: BlockType.AndGate,
    });
  });

  it('builds place commands with selected actuator block types', () => {
    const actuatorBlockTypes = [
      BlockType.Piston,
      BlockType.Motor,
      BlockType.MotorDriver,
      BlockType.TransistorSwitch,
    ];

    for (const blockType of actuatorBlockTypes) {
      expect(
        buildEditCommand({
          mode: 'place',
          clientId: 'client-1',
          commandId: `cmd-${blockType}`,
          tickHint: 45,
          position: { x: blockType, y: 0, z: 2 },
          blockType,
        }),
      ).toEqual({
        type: 'place_block',
        clientId: 'client-1',
        commandId: `cmd-${blockType}`,
        tickHint: 45,
        position: { x: blockType, y: 0, z: 2 },
        blockType,
      });
    }
  });

  it('builds remove commands with air block type', () => {
    expect(
      buildEditCommand({
        mode: 'remove',
        clientId: 'client-1',
        commandId: 'cmd-2',
        tickHint: 43,
        position: { x: 1, y: 2, z: 3 },
        blockType: BlockType.DebugMover,
      }),
    ).toEqual({
      type: 'remove_block',
      clientId: 'client-1',
      commandId: 'cmd-2',
      tickHint: 43,
      position: { x: 1, y: 2, z: 3 },
      blockType: BlockType.Air,
    });
  });
});

describe('buildSetButtonCommand', () => {
  it('builds server-authoritative button input commands', () => {
    expect(
      buildSetButtonCommand({
        clientId: 'client-1',
        commandId: 'cmd-4',
        tickHint: 45,
        position: { x: 2, y: 0, z: 1 },
        buttonPressed: true,
      }),
    ).toEqual({
      type: 'set_button',
      clientId: 'client-1',
      commandId: 'cmd-4',
      tickHint: 45,
      position: { x: 2, y: 0, z: 1 },
      blockType: BlockType.Air,
      buttonPressed: true,
    });
  });
});
