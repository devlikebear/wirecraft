import { describe, expect, it } from 'vitest';
import { cameraPanVectorForKey, shouldIgnoreKeyboardNavigation } from './CameraController';

describe('cameraPanVectorForKey', () => {
  it('maps keyboard navigation keys to planar pan vectors', () => {
    expect(cameraPanVectorForKey('w')).toEqual({ x: 0, z: -1 });
    expect(cameraPanVectorForKey('ArrowUp')).toEqual({ x: 0, z: -1 });
    expect(cameraPanVectorForKey('d')).toEqual({ x: 1, z: 0 });
    expect(cameraPanVectorForKey('ArrowLeft')).toEqual({ x: -1, z: 0 });
  });

  it('ignores non-navigation keys', () => {
    expect(cameraPanVectorForKey('r')).toBeNull();
  });
});

describe('shouldIgnoreKeyboardNavigation', () => {
  it('ignores keyboard navigation while text fields are focused', () => {
    expect(shouldIgnoreKeyboardNavigation({ tagName: 'INPUT' })).toBe(true);
    expect(shouldIgnoreKeyboardNavigation({ tagName: 'TEXTAREA' })).toBe(true);
    expect(shouldIgnoreKeyboardNavigation({ tagName: 'BUTTON' })).toBe(false);
    expect(shouldIgnoreKeyboardNavigation(null)).toBe(false);
  });
});
