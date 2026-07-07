import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { defineComponent } from 'vue';
import { mount } from '@vue/test-utils';
import { useCountdownToNoon } from '@/composables/useCountdownToNoon.js';

// Mounts a throwaway component so the composable's onMounted/onUnmounted
// lifecycle hooks run, and exposes what it returns.
function withCountdown() {
    let api;
    const wrapper = mount(
        defineComponent({
            setup() {
                api = useCountdownToNoon();
                return () => null;
            },
        })
    );
    return { api, wrapper };
}

describe('useCountdownToNoon', () => {
    beforeEach(() => {
        vi.useFakeTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('counts down to 10:00 when the current time is before it', () => {
        // 08:00:00 local -> 2h to target
        vi.setSystemTime(new Date(2026, 0, 1, 8, 0, 0));
        const { api, wrapper } = withCountdown();

        expect(api.isPastNoon.value).toBe(false);
        expect(api.hours.value).toBe('02');
        expect(api.minutes.value).toBe('00');
        expect(api.seconds.value).toBe('00');
        wrapper.unmount();
    });

    it('zero-pads hours, minutes and seconds', () => {
        // 09:57:05 -> 00:02:55 remaining
        vi.setSystemTime(new Date(2026, 0, 1, 9, 57, 5));
        const { api, wrapper } = withCountdown();

        expect(api.hours.value).toBe('00');
        expect(api.minutes.value).toBe('02');
        expect(api.seconds.value).toBe('55');
        wrapper.unmount();
    });

    it('reports isPastNoon and zeroes the clock at/after 10:00', () => {
        vi.setSystemTime(new Date(2026, 0, 1, 10, 0, 0));
        const { api, wrapper } = withCountdown();

        expect(api.isPastNoon.value).toBe(true);
        expect(api.hours.value).toBe('00');
        expect(api.minutes.value).toBe('00');
        expect(api.seconds.value).toBe('00');
        wrapper.unmount();
    });

    it('updates every second via its interval', () => {
        vi.setSystemTime(new Date(2026, 0, 1, 8, 0, 0));
        const { api, wrapper } = withCountdown();
        expect(api.seconds.value).toBe('00');

        // advanceTimersByTime advances the faked clock too, so after 1s the
        // callback sees 08:00:01.
        vi.advanceTimersByTime(1000);
        // 1h 59m 59s remaining
        expect(api.hours.value).toBe('01');
        expect(api.minutes.value).toBe('59');
        expect(api.seconds.value).toBe('59');
        wrapper.unmount();
    });

    it('clears its interval when the component unmounts', () => {
        vi.setSystemTime(new Date(2026, 0, 1, 8, 0, 0));
        const clearSpy = vi.spyOn(globalThis, 'clearInterval');
        const { wrapper } = withCountdown();
        wrapper.unmount();
        expect(clearSpy).toHaveBeenCalled();
    });
});
