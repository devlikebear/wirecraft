import { describe, expect, it } from 'vitest';
import {
  DEBUG_OVERLAY_TEST_IDS,
  calculateFps,
  formatDebugOverlayRows,
  normalizeDebugOverlayState,
} from './DebugOverlay';

describe('formatDebugOverlayRows', () => {
  it('formats runtime state into stable debug rows', () => {
    const rows = formatDebugOverlayRows({
      wsStatus: 'open',
      serverTick: 42,
      snapshotBuffer: 64,
      renderedEntities: 1,
      fps: 59.6,
    });

    expect(rows).toEqual([
      { key: 'wsStatus', label: 'WS', testId: DEBUG_OVERLAY_TEST_IDS.wsStatus, value: 'open' },
      { key: 'serverTick', label: 'Tick', testId: DEBUG_OVERLAY_TEST_IDS.serverTick, value: '42' },
      {
        key: 'snapshotBuffer',
        label: 'Buffer',
        testId: DEBUG_OVERLAY_TEST_IDS.snapshotBuffer,
        value: '64',
      },
      {
        key: 'renderedEntities',
        label: 'Entities',
        testId: DEBUG_OVERLAY_TEST_IDS.renderedEntities,
        value: '1',
      },
      { key: 'fps', label: 'FPS', testId: DEBUG_OVERLAY_TEST_IDS.fps, value: '60' },
    ]);
  });

  it('uses readable fallback values before runtime data arrives', () => {
    const rows = formatDebugOverlayRows(normalizeDebugOverlayState({}));

    expect(rows.map((row) => row.value)).toEqual(['idle', 'n/a', '0', '0', 'n/a']);
  });
});

describe('calculateFps', () => {
  it('calculates rounded FPS from animation frame timestamps', () => {
    expect(calculateFps(1000, 1016.67)).toBe(60);
  });

  it('ignores invalid or first-frame timestamps', () => {
    expect(calculateFps(null, 1016.67)).toBeNull();
    expect(calculateFps(1000, 1000)).toBeNull();
  });
});
