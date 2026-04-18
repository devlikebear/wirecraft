import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { BlockType } from '../net/protocol';
import { BLOCK_TOOLS, createToolbar } from './Toolbar';

class FakeClassList {
  private readonly classes = new Set<string>();

  toggle(className: string, force: boolean): void {
    if (force) {
      this.classes.add(className);
    } else {
      this.classes.delete(className);
    }
  }
}

class FakeElement {
  readonly classList = new FakeClassList();
  readonly children: FakeElement[] = [];
  readonly dataset: Record<string, string> = {};
  readonly style: Record<string, string> = {};
  className = '';
  title = '';
  type = '';

  private readonly attributes = new Map<string, string>();
  private readonly listeners = new Map<string, Array<() => void>>();

  constructor(readonly tagName: string) {}

  appendChild(child: FakeElement): FakeElement {
    this.children.push(child);
    return child;
  }

  setAttribute(name: string, value: string): void {
    this.attributes.set(name, value);
  }

  getAttribute(name: string): string | null {
    return this.attributes.get(name) ?? null;
  }

  addEventListener(name: string, listener: () => void): void {
    const listeners = this.listeners.get(name) ?? [];
    listeners.push(listener);
    this.listeners.set(name, listeners);
  }

  click(): void {
    for (const listener of this.listeners.get('click') ?? []) {
      listener();
    }
  }

  querySelectorAll<T extends FakeElement>(selector: string): T[] {
    if (selector !== 'button[data-block-type]') {
      return [];
    }
    return this.children.filter(
      (child): child is T => child.tagName === 'button' && child.dataset.blockType !== undefined,
    );
  }
}

const documentStub = {
  createElement(tagName: string): FakeElement {
    return new FakeElement(tagName);
  },
};

describe('createToolbar', () => {
  let originalDocument: Document;

  beforeEach(() => {
    originalDocument = globalThis.document;
    globalThis.document = documentStub as unknown as Document;
  });

  afterEach(() => {
    globalThis.document = originalDocument;
  });

  it('defines starter circuit and actuator placement tools', () => {
    expect(BLOCK_TOOLS.map((tool) => tool.blockType)).toEqual([
      BlockType.Solid,
      BlockType.Power,
      BlockType.Wire,
      BlockType.Button,
      BlockType.AndGate,
      BlockType.MCUOutput,
      BlockType.Piston,
      BlockType.Motor,
      BlockType.MotorDriver,
      BlockType.TransistorSwitch,
    ]);
    expect(BLOCK_TOOLS.map((tool) => tool.name)).toEqual([
      'Solid',
      'Power',
      'Wire',
      'Button',
      'AND Gate',
      'MCU Output',
      'Piston',
      'Motor',
      'Motor Driver',
      'Transistor Switch',
    ]);
  });

  it('renders block tool buttons and reports selected block type', () => {
    const selected: BlockType[] = [];
    const toolbar = createToolbar({
      selectedBlockType: BlockType.Solid,
      onSelectBlockType: (blockType) => selected.push(blockType),
    });

    const buttons = [...toolbar.element.querySelectorAll<HTMLButtonElement>('button[data-block-type]')];
    expect(buttons.map((button) => Number(button.dataset.blockType))).toEqual(
      BLOCK_TOOLS.map((tool) => tool.blockType),
    );
    expect(buttons[0].getAttribute('aria-pressed')).toBe('true');

    buttons[2].click();

    expect(selected).toEqual([BlockType.Wire]);
    expect(buttons[0].getAttribute('aria-pressed')).toBe('false');
    expect(buttons[2].getAttribute('aria-pressed')).toBe('true');
  });

  it('reports actuator block selections', () => {
    const selected: BlockType[] = [];
    const toolbar = createToolbar({
      selectedBlockType: BlockType.Power,
      onSelectBlockType: (blockType) => selected.push(blockType),
    });

    const buttons = [...toolbar.element.querySelectorAll<HTMLButtonElement>('button[data-block-type]')];
    const pistonButton = buttons.find((button) => Number(button.dataset.blockType) === BlockType.Piston);
    const motorDriverButton = buttons.find((button) => Number(button.dataset.blockType) === BlockType.MotorDriver);

    expect(pistonButton?.getAttribute('aria-label')).toBe('Piston');
    expect(motorDriverButton?.getAttribute('aria-label')).toBe('Motor Driver');

    pistonButton?.click();
    motorDriverButton?.click();

    expect(selected).toEqual([BlockType.Piston, BlockType.MotorDriver]);
  });
});
