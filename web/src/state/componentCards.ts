import { BlockType } from '../net/protocol';

export interface ComponentPin {
  name: string;
  direction: string;
  signal: string;
}

export interface ComponentCard {
  id: string;
  name: string;
  role: string;
  pins: ComponentPin[];
  wiringNotes: string[];
  warnings: string[];
  simplificationNotes: string[];
}

export const STARTER_COMPONENT_CARDS: ComponentCard[] = [
  {
    id: 'power',
    name: 'Power',
    role: 'Provides a constant High digital signal',
    pins: [{ name: 'OUT', direction: 'output', signal: 'digital high' }],
    wiringNotes: ['Connect OUT to a wire, logic input, or indicator input.'],
    warnings: ['Real supplies need voltage and current limits before powering hardware.'],
    simplificationNotes: ['WireCraft treats this as an ideal always-on digital source for starter circuits.'],
  },
  {
    id: 'wire',
    name: 'Wire',
    role: 'Carries a digital signal between adjacent components',
    pins: [
      { name: 'A', direction: 'bidirectional', signal: 'digital' },
      { name: 'B', direction: 'bidirectional', signal: 'digital' },
    ],
    wiringNotes: ['Place wires next to circuit blocks to bridge their signals through the voxel grid.'],
    warnings: ['Real wires have resistance and safe current ratings.'],
    simplificationNotes: ['Starter wires propagate High, Low, or Unknown without analog voltage drop.'],
  },
  {
    id: 'button',
    name: 'Button',
    role: 'User-controlled digital switch',
    pins: [
      { name: 'IN', direction: 'input', signal: 'digital' },
      { name: 'OUT', direction: 'output', signal: 'digital' },
    ],
    wiringNotes: ['Use set_button commands to press or release the button in the authoritative simulation.'],
    warnings: ['Real buttons bounce and often need pull-up or pull-down resistors.'],
    simplificationNotes: ['WireCraft models a released button as Low and a pressed button as High.'],
  },
  {
    id: 'and_gate',
    name: 'AND Gate',
    role: 'Outputs High only when two inputs are High',
    pins: [
      { name: 'A', direction: 'input', signal: 'digital' },
      { name: 'B', direction: 'input', signal: 'digital' },
      { name: 'OUT', direction: 'output', signal: 'digital' },
    ],
    wiringNotes: ['Feed two input paths into the gate, then connect OUT to another wire or output block.'],
    warnings: ['Real logic chips need a matching supply voltage and ground reference.'],
    simplificationNotes: ['The current starter runtime infers gate inputs from adjacent circuit blocks.'],
  },
  {
    id: 'mcu_output',
    name: 'MCU Output',
    role: 'Represents a microcontroller output pin or observed endpoint',
    pins: [{ name: 'IN', direction: 'input', signal: 'digital' }],
    wiringNotes: ['Connect a wire or gate output to observe the resulting server-authoritative signal.'],
    warnings: ['Real MCU pins cannot drive motors or high-current loads directly.'],
    simplificationNotes: ['WireCraft uses this block as a readable digital endpoint before code runtime exists.'],
  },
  {
    id: 'led',
    name: 'LED',
    role: 'Visual indicator for a digital output',
    pins: [
      { name: 'A', direction: 'input', signal: 'digital high' },
      { name: 'K', direction: 'input', signal: 'ground reference' },
    ],
    wiringNotes: ['Drive the anode through a resistor and connect the cathode to ground in real circuits.'],
    warnings: ['Real LEDs need current-limiting resistors to avoid damage.'],
    simplificationNotes: ['LED is documented as a starter card before it becomes a placeable voxel block.'],
  },
  {
    id: 'resistor',
    name: 'Resistor',
    role: 'Limits current or creates a pull-up/pull-down path',
    pins: [
      { name: 'A', direction: 'passive', signal: 'analog' },
      { name: 'B', direction: 'passive', signal: 'analog' },
    ],
    wiringNotes: ['Use resistors with LEDs and buttons in real hardware to control current and default voltage.'],
    warnings: ['Choose resistance and power rating for the actual voltage and load.'],
    simplificationNotes: ['WireCraft documents resistors now but does not simulate analog resistance in Phase 2.'],
  },
  {
    id: 'motor',
    name: 'Motor',
    role: 'Converts a driver-enabled signal into actuator motion',
    pins: [
      { name: 'DRIVE', direction: 'input', signal: 'motor drive' },
      { name: 'GND', direction: 'input', signal: 'ground reference' },
    ],
    wiringNotes: ['Connect motors through a motor driver or transistor switch, not directly to MCU GPIO.'],
    warnings: ['MCU GPIO pins cannot drive a motor directly.'],
    simplificationNotes: [
      'WireCraft models motor enable as a digital driver signal and does not simulate current draw.',
    ],
  },
  {
    id: 'motor_driver',
    name: 'Motor Driver',
    role: 'Lets a low-current logic signal control a motor load',
    pins: [
      { name: 'IN', direction: 'input', signal: 'digital' },
      { name: 'VM', direction: 'input', signal: 'motor supply' },
      { name: 'OUT', direction: 'output', signal: 'motor drive' },
    ],
    wiringNotes: ['Place a motor driver or transistor switch between logic and motor loads.'],
    warnings: ['Real motor drivers need a rated motor supply, shared ground, and current headroom.'],
    simplificationNotes: ['WireCraft treats the driver as a digital permission gate for motor actuation.'],
  },
  {
    id: 'transistor_switch',
    name: 'Transistor Switch',
    role: 'Switches a motor or other load from a digital control signal',
    pins: [
      { name: 'CTRL', direction: 'input', signal: 'digital' },
      { name: 'LOAD', direction: 'output', signal: 'switched load' },
      { name: 'GND', direction: 'input', signal: 'ground reference' },
    ],
    wiringNotes: ['Use a transistor switch when a logic signal must control a higher-current load.'],
    warnings: ['Real motor switches need flyback protection and correct current ratings.'],
    simplificationNotes: [
      'WireCraft models the switch as digital on/off and ignores analog saturation, heat, and diode behavior.',
    ],
  },
];

const cardsByID = new Map(STARTER_COMPONENT_CARDS.map((card) => [card.id, card]));
const cardIDsByBlockType = new Map<BlockType, string>([
  [BlockType.Power, 'power'],
  [BlockType.Wire, 'wire'],
  [BlockType.Button, 'button'],
  [BlockType.AndGate, 'and_gate'],
  [BlockType.MCUOutput, 'mcu_output'],
  [BlockType.Motor, 'motor'],
  [BlockType.MotorDriver, 'motor_driver'],
  [BlockType.TransistorSwitch, 'transistor_switch'],
]);

export function findComponentCard(id: string): ComponentCard | null {
  return cardsByID.get(id) ?? null;
}

export function cardForBlockType(blockType: BlockType): ComponentCard | null {
  const id = cardIDsByBlockType.get(blockType);
  return id ? findComponentCard(id) : null;
}

export function validateComponentCards(cards: ComponentCard[]): string[] {
  const errors: string[] = [];
  const seenIDs = new Set<string>();

  for (const card of cards) {
    if (!card.id || !card.name || !card.role) {
      errors.push(`card ${card.id || '<missing>'} is missing id, name, or role`);
    }
    if (seenIDs.has(card.id)) {
      errors.push(`duplicate card id ${card.id}`);
    }
    seenIDs.add(card.id);

    if (
      card.pins.length === 0 ||
      card.wiringNotes.length === 0 ||
      card.warnings.length === 0 ||
      card.simplificationNotes.length === 0
    ) {
      errors.push(`card ${card.id} is missing required content`);
    }

    for (const pin of card.pins) {
      if (!pin.name || !pin.direction || !pin.signal) {
        errors.push(`card ${card.id} has an incomplete pin`);
      }
    }
  }

  return errors;
}
