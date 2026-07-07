import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { defineComponent } from 'vue';
import { mount } from '@vue/test-utils';
import { useNow } from '@/composables/useNow.js';

function withNow(interval) {
    let api;
    const wrapper = mount(
        defineComponent({
            setup() {
                api = useNow(interval);
                return () => null;
            },
        })
    );
    return { api, wrapper };
}

describe('useNow', () => {
    beforeEach(() => {
        vi.useFakeTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('initializes now to the current timestamp', () => {
        vi.setSystemTime(new Date(2026, 0, 1, 0, 0, 0));
        const { api, wrapper } = withNow();
        expect(api.now.value).toBe(Date.now());
        wrapper.unmount();
    });

    it('advances now on each default (1000ms) tick', () => {
        vi.setSystemTime(1000);
        const { api, wrapper } = withNow();
        expect(api.now.value).toBe(1000);

        // advanceTimersByTime also advances the (faked) system clock, so the
        // tick reads Date.now() === 2000.
        vi.advanceTimersByTime(1000);
        expect(api.now.value).toBe(2000);
        wrapper.unmount();
    });

    it('respects a custom interval', () => {
        vi.setSystemTime(0);
        const { api, wrapper } = withNow(5000);

        vi.advanceTimersByTime(3000);
        // interval not yet elapsed
        expect(api.now.value).toBe(0);

        vi.advanceTimersByTime(2000);
        expect(api.now.value).toBe(5000);
        wrapper.unmount();
    });

    it('stops updating after unmount', () => {
        vi.setSystemTime(0);
        const { api, wrapper } = withNow();
        wrapper.unmount();

        vi.setSystemTime(10000);
        vi.advanceTimersByTime(10000);
        expect(api.now.value).toBe(0);
    });
});
