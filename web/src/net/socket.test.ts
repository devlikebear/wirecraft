import { describe, expect, it } from 'vitest';
import { SnapshotSocket, websocketURLFromLocation } from './socket';

describe('websocketURLFromLocation', () => {
  it('derives a ws URL from an http page origin', () => {
    const url = websocketURLFromLocation({ protocol: 'http:', host: '127.0.0.1:5173' });

    expect(url).toBe('ws://127.0.0.1:5173/ws');
  });

  it('derives a wss URL from an https page origin', () => {
    const url = websocketURLFromLocation({ protocol: 'https:', host: 'wirecraft.example' });

    expect(url).toBe('wss://wirecraft.example/ws');
  });

  it('normalizes custom paths', () => {
    const url = websocketURLFromLocation({ protocol: 'http:', host: 'localhost:5173' }, 'sim');

    expect(url).toBe('ws://localhost:5173/sim');
  });
});

describe('SnapshotSocket', () => {
  it('forwards parsed command acknowledgements with snapshots', () => {
    let fakeSocket: FakeWebSocket | null = null;
    const snapshots: unknown[] = [];
    const socket = new SnapshotSocket({
      url: 'ws://wirecraft.test/ws',
      webSocketFactory: (url) => {
        fakeSocket = new FakeWebSocket(url);
        return fakeSocket as unknown as WebSocket;
      },
      onSnapshot: (snapshot) => snapshots.push(snapshot),
      onError: (error) => {
        throw error;
      },
    });

    socket.connect();
    fakeSocket?.emitMessage(
      JSON.stringify({
        tick: 1,
        serverTimeMs: 1700000000001,
        blocks: [],
        entities: [],
        circuit: { nodes: [] },
        commandAcks: [
          { clientId: 'client-1', commandId: 'cmd-1', status: 'accepted' },
          {
            clientId: 'client-1',
            commandId: 'cmd-1',
            status: 'rejected',
            reason: 'duplicate_command',
          },
        ],
        stats: {
          clientCount: 1,
          commandQueueLength: 2,
          snapshotBytes: 128,
        },
      }),
    );

    expect(snapshots).toHaveLength(1);
    expect(snapshots[0]).toMatchObject({
      commandAcks: [
        { clientId: 'client-1', commandId: 'cmd-1', status: 'accepted' },
        {
          clientId: 'client-1',
          commandId: 'cmd-1',
          status: 'rejected',
          reason: 'duplicate_command',
        },
      ],
    });
  });
});

class FakeWebSocket {
  readonly url: string;
  readyState = WebSocket.OPEN;
  onopen: ((event: Event) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;

  constructor(url: string) {
    this.url = url;
  }

  close(): void {
    this.readyState = WebSocket.CLOSED;
  }

  send(): void {
    return undefined;
  }

  emitMessage(data: string): void {
    this.onmessage?.({ data } as MessageEvent);
  }
}
