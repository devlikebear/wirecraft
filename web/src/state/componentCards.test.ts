import { describe, expect, it } from 'vitest';
import { BlockType } from '../net/protocol';
import {
  STARTER_COMPONENT_CARDS,
  cardForBlockType,
  findComponentCard,
  validateComponentCards,
} from './componentCards';

describe('starter component cards', () => {
  it('defines the expected beginner cards with required fields', () => {
    expect(STARTER_COMPONENT_CARDS.map((card) => card.id)).toEqual([
      'power',
      'wire',
      'button',
      'and_gate',
      'mcu_output',
      'led',
      'resistor',
    ]);
    expect(validateComponentCards(STARTER_COMPONENT_CARDS)).toEqual([]);
  });

  it('maps placeable circuit block types to their starter cards', () => {
    expect(cardForBlockType(BlockType.Power)?.id).toBe('power');
    expect(cardForBlockType(BlockType.Wire)?.id).toBe('wire');
    expect(cardForBlockType(BlockType.Button)?.id).toBe('button');
    expect(cardForBlockType(BlockType.AndGate)?.id).toBe('and_gate');
    expect(cardForBlockType(BlockType.MCUOutput)?.id).toBe('mcu_output');
    expect(cardForBlockType(BlockType.Solid)).toBeNull();
  });

  it('finds cards by id', () => {
    expect(findComponentCard('led')?.name).toBe('LED');
    expect(findComponentCard('missing')).toBeNull();
  });
});
