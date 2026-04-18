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
      'motor',
      'motor_driver',
      'transistor_switch',
    ]);
    expect(validateComponentCards(STARTER_COMPONENT_CARDS)).toEqual([]);
  });

  it('maps placeable circuit block types to their starter cards', () => {
    expect(cardForBlockType(BlockType.Power)?.id).toBe('power');
    expect(cardForBlockType(BlockType.Wire)?.id).toBe('wire');
    expect(cardForBlockType(BlockType.Button)?.id).toBe('button');
    expect(cardForBlockType(BlockType.AndGate)?.id).toBe('and_gate');
    expect(cardForBlockType(BlockType.MCUOutput)?.id).toBe('mcu_output');
    expect(cardForBlockType(BlockType.Motor)?.id).toBe('motor');
    expect(cardForBlockType(BlockType.MotorDriver)?.id).toBe('motor_driver');
    expect(cardForBlockType(BlockType.TransistorSwitch)?.id).toBe('transistor_switch');
    expect(cardForBlockType(BlockType.Solid)).toBeNull();
  });

  it('explains motor driver constraints in beginner-facing cards', () => {
    expect(findComponentCard('motor')?.warnings).toContain('MCU GPIO pins cannot drive a motor directly.');
    expect(findComponentCard('motor')?.simplificationNotes).toContain(
      'WireCraft models motor enable as a digital driver signal and does not simulate current draw.',
    );
    expect(findComponentCard('motor_driver')?.wiringNotes).toContain(
      'Place a motor driver or transistor switch between logic and motor loads.',
    );
    expect(findComponentCard('transistor_switch')?.warnings).toContain(
      'Real motor switches need flyback protection and correct current ratings.',
    );
  });

  it('finds cards by id', () => {
    expect(findComponentCard('led')?.name).toBe('LED');
    expect(findComponentCard('missing')).toBeNull();
  });
});
