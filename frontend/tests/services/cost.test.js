import { describe, it, expect } from 'vitest';
import { COST_ITEMS_INFO, INVENTORY_ITEMS } from '@/services/cost.js';

describe('cost', () => {
    describe('COST_ITEMS_INFO', () => {
        it('exposes an entry for every known item type', () => {
            expect(Object.keys(COST_ITEMS_INFO).sort()).toEqual(
                [
                    'blueKey',
                    'coin',
                    'fuel',
                    'goldenKey',
                    'masterKey',
                    'redKey',
                ].sort()
            );
        });

        it('maps each item to a name and an icon path', () => {
            expect(COST_ITEMS_INFO.masterKey).toEqual({
                name: 'TNT',
                icon: '/images/icons/masterKey.png',
            });
        });

        it('derives the icon path from the item id', () => {
            for (const [id, info] of Object.entries(COST_ITEMS_INFO)) {
                expect(info.icon).toBe(`/images/icons/${id}.png`);
                expect(typeof info.name).toBe('string');
                expect(info.name.length).toBeGreaterThan(0);
            }
        });
    });

    describe('INVENTORY_ITEMS', () => {
        it('lists exactly the four key types', () => {
            expect(INVENTORY_ITEMS).toEqual([
                'masterKey',
                'blueKey',
                'redKey',
                'goldenKey',
            ]);
        });

        it('only contains ids that exist in COST_ITEMS_INFO', () => {
            for (const id of INVENTORY_ITEMS) {
                expect(COST_ITEMS_INFO).toHaveProperty(id);
            }
        });
    });
});
