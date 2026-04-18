import { describe, expect, it } from 'vitest';
import { findComponentCard } from '../state/componentCards';
import { inspectPanelText } from './InspectPanel';

describe('inspectPanelText', () => {
  it('formats a starter component card for compact display', () => {
    const card = findComponentCard('button');
    if (!card) {
      throw new Error('missing button card fixture');
    }

    expect(inspectPanelText(card)).toContain('Button');
    expect(inspectPanelText(card)).toContain('Role: User-controlled digital switch');
    expect(inspectPanelText(card)).toContain('Pins: IN input, OUT output');
    expect(inspectPanelText(card)).toContain('Warning: Real buttons bounce');
  });

  it('uses a neutral prompt when no card is selected', () => {
    expect(inspectPanelText(null)).toBe('Select a circuit component.');
  });
});
