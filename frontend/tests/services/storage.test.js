import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { getJSON, setJSON } from '@/services/storage.js';

describe('storage', () => {
    beforeEach(() => {
        localStorage.clear();
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('getJSON', () => {
        it('returns the parsed value for a stored key', () => {
            localStorage.setItem('k', JSON.stringify({ a: 1, b: [2, 3] }));
            expect(getJSON('k')).toEqual({ a: 1, b: [2, 3] });
        });

        it('returns the default fallback ([]) when the key is absent', () => {
            expect(getJSON('missing')).toEqual([]);
        });

        it('returns a custom fallback when the key is absent', () => {
            expect(getJSON('missing', { x: 1 })).toEqual({ x: 1 });
        });

        it('returns the fallback when stored value is not valid JSON', () => {
            localStorage.setItem('bad', '{not json');
            const spy = vi.spyOn(console, 'error').mockImplementation(() => {});
            expect(getJSON('bad', 'fallback')).toBe('fallback');
            expect(spy).toHaveBeenCalled();
        });

        it('parses stored falsy JSON values correctly', () => {
            localStorage.setItem('zero', JSON.stringify(0));
            expect(getJSON('zero', 99)).toBe(0);
            localStorage.setItem('false', JSON.stringify(false));
            expect(getJSON('false', true)).toBe(false);
        });

        it('returns the fallback for an empty-string value (falsy)', () => {
            localStorage.setItem('empty', '');
            expect(getJSON('empty', 'fb')).toBe('fb');
        });
    });

    describe('setJSON', () => {
        it('serializes and stores the value', () => {
            setJSON('k', { hello: 'world' });
            expect(localStorage.getItem('k')).toBe(
                JSON.stringify({ hello: 'world' })
            );
        });

        it('round-trips with getJSON', () => {
            const value = [{ id: 1 }, { id: 2 }];
            setJSON('list', value);
            expect(getJSON('list')).toEqual(value);
        });

        it('swallows and logs errors when serialization fails', () => {
            const spy = vi.spyOn(console, 'error').mockImplementation(() => {});
            const circular = {};
            circular.self = circular;
            expect(() => setJSON('circ', circular)).not.toThrow();
            expect(spy).toHaveBeenCalled();
        });
    });
});
