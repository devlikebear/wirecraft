import { describe, expect, it } from 'vitest';
import type { CircuitNodeSnapshot } from '../net/protocol';
import { createCircuitOverlayItems, groupCircuitOverlayItems } from './CircuitOverlay';

describe('createCircuitOverlayItems', () => {
  it('maps circuit nodes to markers above their block positions', () => {
    const nodes: CircuitNodeSnapshot[] = [
      {
        position: { x: 1, y: 0, z: 2 },
        nodeId: '1:0:2',
        nodeType: 'wire',
        signalState: 'high',
      },
      {
        position: { x: 2, y: 0, z: 2 },
        nodeId: '2:0:2',
        nodeType: 'mcu_output',
        signalState: 'low',
      },
    ];

    expect(createCircuitOverlayItems(nodes)).toEqual([
      {
        key: '1:0:2:high',
        nodeId: '1:0:2',
        nodeType: 'wire',
        signalState: 'high',
        position: { x: 1, y: 1.08, z: 2 },
      },
      {
        key: '2:0:2:low',
        nodeId: '2:0:2',
        nodeType: 'mcu_output',
        signalState: 'low',
        position: { x: 2, y: 1.08, z: 2 },
      },
    ]);
  });
});

describe('groupCircuitOverlayItems', () => {
  it('groups markers by signal state', () => {
    const items = createCircuitOverlayItems([
      {
        position: { x: 0, y: 0, z: 0 },
        nodeId: '0:0:0',
        nodeType: 'power_source',
        signalState: 'high',
      },
      {
        position: { x: 1, y: 0, z: 0 },
        nodeId: '1:0:0',
        nodeType: 'button',
        signalState: 'low',
      },
      {
        position: { x: 2, y: 0, z: 0 },
        nodeId: '2:0:0',
        nodeType: 'wire',
        signalState: 'unknown',
      },
    ]);

    const grouped = groupCircuitOverlayItems(items);

    expect(grouped.get('high')?.map((item) => item.nodeId)).toEqual(['0:0:0']);
    expect(grouped.get('low')?.map((item) => item.nodeId)).toEqual(['1:0:0']);
    expect(grouped.get('unknown')?.map((item) => item.nodeId)).toEqual(['2:0:0']);
  });
});
