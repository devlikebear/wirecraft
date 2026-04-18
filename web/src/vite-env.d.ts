/// <reference types="vite/client" />

import type { BlockType, Position } from './net/protocol';
import type { SnapshotStore } from './state/snapshotStore';

declare global {
  interface Window {
    wirecraft: {
      placeBlock(position: Position, blockType?: BlockType): void;
      removeBlock(position: Position): void;
      setButton(position: Position, buttonPressed: boolean): void;
      snapshots: SnapshotStore;
    };
  }
}
