import type { Snapshot } from '../net/protocol';

export interface SnapshotStoreOptions {
  maxSnapshots?: number;
}

interface SnapshotEntry {
  snapshot: Snapshot;
  receivedAtMs: number;
}

export class SnapshotStore {
  private readonly maxSnapshots: number;
  private entries: SnapshotEntry[] = [];

  constructor(options: SnapshotStoreOptions = {}) {
    this.maxSnapshots = Math.max(1, options.maxSnapshots ?? 64);
  }

  get length(): number {
    return this.entries.length;
  }

  append(snapshot: Snapshot, receivedAtMs = Date.now()): void {
    this.entries = [...this.entries, { snapshot, receivedAtMs }].slice(-this.maxSnapshots);
  }

  latest(): Snapshot | null {
    return this.entries.at(-1)?.snapshot ?? null;
  }

  all(): Snapshot[] {
    return this.entries.map((entry) => entry.snapshot);
  }

  renderServerTimeMs(clientTimeMs: number, interpolationDelayMs: number): number | null {
    const latest = this.entries.at(-1);
    if (!latest) {
      return null;
    }

    const elapsedSinceLatestMs = Math.max(0, clientTimeMs - latest.receivedAtMs);
    return latest.snapshot.serverTimeMs + elapsedSinceLatestMs - Math.max(0, interpolationDelayMs);
  }

  clear(): void {
    this.entries = [];
  }
}
