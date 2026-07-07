import { describe, it, expect, beforeEach } from 'vitest';
import { notificationService } from '@/services/notificationService.js';

const msg = (createdAt, reason = 'reward') => ({
    createdAt,
    content: { [reason]: {} },
});

describe('notificationService', () => {
    beforeEach(() => {
        // Reset persisted state between tests, then clear any received
        // messages carried over via the module-level reactive state.
        localStorage.clear();
        notificationService.setReceivedMessages([]);
        notificationService.markAllAsSeen();
    });

    describe('setReceivedMessages', () => {
        it('keeps createdAt and derives reason from the first content key', () => {
            notificationService.setReceivedMessages([
                msg('2026-01-01', 'coin'),
            ]);
            const stored = JSON.parse(
                localStorage.getItem('received_messages')
            );
            expect(stored).toEqual([
                { createdAt: '2026-01-01', reason: 'coin' },
            ]);
        });

        it('filters out messages missing createdAt or content', () => {
            notificationService.setReceivedMessages([
                msg('2026-01-01'),
                { content: { x: {} } }, // no createdAt
                { createdAt: '2026-01-02' }, // no content
                null,
                undefined,
            ]);
            const stored = JSON.parse(
                localStorage.getItem('received_messages')
            );
            expect(stored).toEqual([
                { createdAt: '2026-01-01', reason: 'reward' },
            ]);
        });

        it('persists the received messages to localStorage', () => {
            notificationService.setReceivedMessages([msg('a'), msg('b')]);
            expect(
                JSON.parse(localStorage.getItem('received_messages'))
            ).toHaveLength(2);
        });
    });

    describe('hasUnreadMessages', () => {
        it('is false when there are no received messages', () => {
            notificationService.setReceivedMessages([]);
            expect(notificationService.hasUnreadMessages.value).toBe(false);
        });

        it('is true when a received message has not been seen', () => {
            notificationService.setReceivedMessages([msg('2026-01-01')]);
            expect(notificationService.hasUnreadMessages.value).toBe(true);
        });

        it('becomes false after markAllAsSeen', () => {
            notificationService.setReceivedMessages([msg('a'), msg('b')]);
            expect(notificationService.hasUnreadMessages.value).toBe(true);
            notificationService.markAllAsSeen();
            expect(notificationService.hasUnreadMessages.value).toBe(false);
        });

        it('is true again when a new unseen message arrives after seeing', () => {
            notificationService.setReceivedMessages([msg('a')]);
            notificationService.markAllAsSeen();
            expect(notificationService.hasUnreadMessages.value).toBe(false);

            notificationService.setReceivedMessages([msg('a'), msg('b')]);
            expect(notificationService.hasUnreadMessages.value).toBe(true);
        });

        it('persists seen messages so they survive re-init', () => {
            notificationService.setReceivedMessages([msg('a')]);
            notificationService.markAllAsSeen();
            expect(JSON.parse(localStorage.getItem('seen_messages'))).toEqual([
                { createdAt: 'a', reason: 'reward' },
            ]);
        });
    });
});
