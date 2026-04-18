import { describe, expect, it } from 'vitest';
import { websocketURLFromLocation } from './socket';

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
