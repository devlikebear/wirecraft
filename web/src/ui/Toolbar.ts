import { BlockType } from '../net/protocol';

export interface BlockTool {
  blockType: BlockType;
  name: string;
  color: string;
}

export interface ToolbarOptions {
  selectedBlockType: BlockType;
  onSelectBlockType: (blockType: BlockType) => void;
}

export interface Toolbar {
  element: HTMLElement;
  setSelectedBlockType(blockType: BlockType): void;
}

export const BLOCK_TOOLS: BlockTool[] = [
  { blockType: BlockType.Solid, name: 'Solid', color: '#6d8f7d' },
  { blockType: BlockType.Power, name: 'Power', color: '#f05a4f' },
  { blockType: BlockType.Wire, name: 'Wire', color: '#c9824a' },
  { blockType: BlockType.Button, name: 'Button', color: '#4f8df0' },
  { blockType: BlockType.AndGate, name: 'AND Gate', color: '#7d6cf2' },
  { blockType: BlockType.MCUOutput, name: 'MCU Output', color: '#35b58a' },
  { blockType: BlockType.Piston, name: 'Piston', color: '#b8c4d6' },
  { blockType: BlockType.Motor, name: 'Motor', color: '#45b7a8' },
  { blockType: BlockType.MotorDriver, name: 'Motor Driver', color: '#f2c14e' },
  { blockType: BlockType.TransistorSwitch, name: 'Transistor Switch', color: '#d973a8' },
];

export function createToolbar(options: ToolbarOptions): Toolbar {
  let selectedBlockType = options.selectedBlockType;
  const element = document.createElement('nav');
  element.className = 'block-toolbar';
  element.setAttribute('aria-label', 'Block tools');

  const buttons = new Map<BlockType, HTMLButtonElement>();
  for (const tool of BLOCK_TOOLS) {
    const button = document.createElement('button');
    button.type = 'button';
    button.className = 'block-toolbar__button';
    button.title = tool.name;
    button.setAttribute('aria-label', tool.name);
    button.dataset.blockType = String(tool.blockType);

    const swatch = document.createElement('span');
    swatch.className = 'block-toolbar__swatch';
    swatch.style.backgroundColor = tool.color;
    button.appendChild(swatch);

    button.addEventListener('click', () => {
      selectedBlockType = tool.blockType;
      updateSelectedState();
      options.onSelectBlockType(tool.blockType);
    });

    buttons.set(tool.blockType, button);
    element.appendChild(button);
  }

  function updateSelectedState(): void {
    for (const [blockType, button] of buttons) {
      const selected = blockType === selectedBlockType;
      button.classList.toggle('block-toolbar__button--selected', selected);
      button.setAttribute('aria-pressed', String(selected));
    }
  }

  updateSelectedState();

  return {
    element,
    setSelectedBlockType(blockType: BlockType) {
      selectedBlockType = blockType;
      updateSelectedState();
    },
  };
}
