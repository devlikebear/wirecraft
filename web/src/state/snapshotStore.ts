import type { Snapshot } from '../net/protocol';

export interface SnapshotStoreOptions {
  maxSnapshots?: number;
}

export class SnapshotStore {
  private readonly maxSnapshots: number;
  private snapshots: Snapshot[] = [];

  constructor(options: SnapshotStoreOptions = {}) {
    this.maxSnapshots = Math.max(1, options.maxSnapshots ?? 64);
  }

  get length(): number {
    return this.snapshots.length;
  }

  append(snapshot: Snapshot): void {
    this.snapshots = [...this.snapshots, snapshot].slice(-this.maxSnapshots);
  }

  latest(): Snapshot | null {
    return this.snapshots.at(-1) ?? null;
  }

  all(): Snapshot[] {
    return [...this.snapshots];
  }

  clear(): void {
    this.snapshots = [];
  }
}
