export const DEBUG_OVERLAY_TEST_IDS = {
  root: 'wirecraft-debug-overlay',
  wsStatus: 'debug-ws-status',
  serverTick: 'debug-server-tick',
  clientCount: 'debug-client-count',
  snapshotBuffer: 'debug-snapshot-buffer',
  renderedEntities: 'debug-rendered-entities',
  fps: 'debug-fps',
} as const;

export interface DebugOverlayState {
  wsStatus: string;
  serverTick: number | null;
  clientCount: number;
  snapshotBuffer: number;
  renderedEntities: number;
  fps: number | null;
}

export type DebugOverlayStateInput = Partial<DebugOverlayState>;

export interface DebugOverlayRow {
  key: keyof DebugOverlayState;
  label: string;
  testId: string;
  value: string;
}

const rowMeta: Array<Omit<DebugOverlayRow, 'value'>> = [
  { key: 'wsStatus', label: 'WS', testId: DEBUG_OVERLAY_TEST_IDS.wsStatus },
  { key: 'serverTick', label: 'Tick', testId: DEBUG_OVERLAY_TEST_IDS.serverTick },
  { key: 'clientCount', label: 'Clients', testId: DEBUG_OVERLAY_TEST_IDS.clientCount },
  { key: 'snapshotBuffer', label: 'Buffer', testId: DEBUG_OVERLAY_TEST_IDS.snapshotBuffer },
  { key: 'renderedEntities', label: 'Entities', testId: DEBUG_OVERLAY_TEST_IDS.renderedEntities },
  { key: 'fps', label: 'FPS', testId: DEBUG_OVERLAY_TEST_IDS.fps },
];

export function normalizeDebugOverlayState(input: DebugOverlayStateInput): DebugOverlayState {
  return {
    wsStatus: input.wsStatus ?? 'idle',
    serverTick: input.serverTick ?? null,
    clientCount: Math.max(0, input.clientCount ?? 0),
    snapshotBuffer: Math.max(0, input.snapshotBuffer ?? 0),
    renderedEntities: Math.max(0, input.renderedEntities ?? 0),
    fps: input.fps ?? null,
  };
}

export function formatDebugOverlayRows(input: DebugOverlayStateInput): DebugOverlayRow[] {
  const state = normalizeDebugOverlayState(input);

  return rowMeta.map((meta) => ({
    ...meta,
    value: formatDebugValue(meta.key, state),
  }));
}

export function calculateFps(previousTimestampMs: number | null, currentTimestampMs: number): number | null {
  if (previousTimestampMs === null) {
    return null;
  }

  const deltaMs = currentTimestampMs - previousTimestampMs;
  if (deltaMs <= 0) {
    return null;
  }

  return Math.round(1000 / deltaMs);
}

export class DebugOverlay {
  readonly element: HTMLDivElement;

  private readonly valueElements = new Map<keyof DebugOverlayState, HTMLSpanElement>();

  constructor(parent: HTMLElement) {
    this.element = document.createElement('div');
    this.element.className = 'debug-overlay';
    this.element.dataset.testid = DEBUG_OVERLAY_TEST_IDS.root;
    this.element.setAttribute('aria-label', 'WireCraft debug overlay');

    for (const row of formatDebugOverlayRows({})) {
      const rowElement = document.createElement('div');
      rowElement.className = 'debug-overlay__row';

      const label = document.createElement('span');
      label.className = 'debug-overlay__label';
      label.textContent = row.label;

      const value = document.createElement('span');
      value.className = 'debug-overlay__value';
      value.dataset.testid = row.testId;
      value.textContent = row.value;

      rowElement.append(label, value);
      this.element.append(rowElement);
      this.valueElements.set(row.key, value);
    }

    parent.append(this.element);
  }

  update(input: DebugOverlayStateInput): void {
    for (const row of formatDebugOverlayRows(input)) {
      const element = this.valueElements.get(row.key);
      if (element && element.textContent !== row.value) {
        element.textContent = row.value;
      }
    }
  }

  dispose(): void {
    this.element.remove();
    this.valueElements.clear();
  }
}

function formatDebugValue(key: keyof DebugOverlayState, state: DebugOverlayState): string {
  switch (key) {
    case 'wsStatus':
      return state.wsStatus;
    case 'serverTick':
      return state.serverTick === null ? 'n/a' : String(Math.round(state.serverTick));
    case 'clientCount':
      return String(Math.round(state.clientCount));
    case 'snapshotBuffer':
      return String(Math.round(state.snapshotBuffer));
    case 'renderedEntities':
      return String(Math.round(state.renderedEntities));
    case 'fps':
      return state.fps === null ? 'n/a' : String(Math.round(state.fps));
  }
}
