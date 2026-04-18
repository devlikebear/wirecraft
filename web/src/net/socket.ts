import type { Command, Snapshot } from './protocol';
import { parseSnapshot } from './protocol';

export type SnapshotSocketStatus = 'idle' | 'connecting' | 'open' | 'closed' | 'error';

export interface LocationLike {
  protocol: string;
  host: string;
}

export interface SnapshotSocketOptions {
  url?: string;
  path?: string;
  location?: LocationLike;
  onSnapshot?: (snapshot: Snapshot) => void;
  onStatusChange?: (status: SnapshotSocketStatus) => void;
  onError?: (error: Error) => void;
  webSocketFactory?: (url: string) => WebSocket;
}

export function websocketURLFromLocation(location: LocationLike, path = '/ws'): string {
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${protocol}//${location.host}${normalizedPath}`;
}

export class SnapshotSocket {
  private socket: WebSocket | null = null;
  private status: SnapshotSocketStatus = 'idle';
  private readonly url: string;
  private readonly onSnapshot: (snapshot: Snapshot) => void;
  private readonly onStatusChange: (status: SnapshotSocketStatus) => void;
  private readonly onError: (error: Error) => void;
  private readonly webSocketFactory: (url: string) => WebSocket;

  constructor(options: SnapshotSocketOptions = {}) {
    this.url =
      options.url ??
      websocketURLFromLocation(options.location ?? defaultLocation(), options.path ?? '/ws');
    this.onSnapshot = options.onSnapshot ?? (() => undefined);
    this.onStatusChange = options.onStatusChange ?? (() => undefined);
    this.onError = options.onError ?? (() => undefined);
    this.webSocketFactory = options.webSocketFactory ?? ((url) => new WebSocket(url));
  }

  connect(): void {
    if (this.socket && this.socket.readyState !== WebSocket.CLOSED) {
      return;
    }

    this.setStatus('connecting');
    const socket = this.webSocketFactory(this.url);
    this.socket = socket;

    socket.onopen = () => this.setStatus('open');
    socket.onclose = () => this.setStatus('closed');
    socket.onerror = () => {
      this.setStatus('error');
      this.onError(new Error(`WebSocket connection failed: ${this.url}`));
    };
    socket.onmessage = (event) => this.handleMessage(event);
  }

  close(): void {
    if (!this.socket || this.socket.readyState === WebSocket.CLOSED) {
      return;
    }

    this.socket.close();
  }

  sendCommand(command: Command): void {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket is not open');
    }

    this.socket.send(JSON.stringify(command));
  }

  private handleMessage(event: MessageEvent): void {
    if (typeof event.data !== 'string') {
      this.onError(new Error('Expected snapshot WebSocket message to be a JSON string'));
      return;
    }

    let parsed: unknown;
    try {
      parsed = JSON.parse(event.data);
    } catch (error) {
      this.onError(error instanceof Error ? error : new Error('Failed to parse snapshot message'));
      return;
    }

    const snapshot = parseSnapshot(parsed);
    if (snapshot === null) {
      this.onError(new Error('Received invalid snapshot payload'));
      return;
    }

    this.onSnapshot(snapshot);
  }

  private setStatus(status: SnapshotSocketStatus): void {
    if (this.status === status) {
      return;
    }

    this.status = status;
    this.onStatusChange(status);
  }
}

function defaultLocation(): LocationLike {
  if (typeof window === 'undefined') {
    throw new Error('SnapshotSocket requires a browser location or explicit URL');
  }

  return window.location;
}
