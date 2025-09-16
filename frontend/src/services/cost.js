const ITEM_NAMES = {
    masterKey: { name: 'دینامیت' },
    blueKey: { name: 'چکش آهنی' },
    redKey: { name: 'چکش نقره‌ای' },
    goldenKey: { name: 'چکش طلایی' },
    fuel: { name: 'سوخت' },
    coin: { name: 'کلاه' },
};
const INVENTORY_ITEMS = ['masterKey', 'blueKey', 'redKey', 'goldenKey'];

function getItemIcon(type) {
    return `/images/icons/${type}.png`;
}

const COST_ITEMS_INFO = Object.fromEntries(
    Object.entries(ITEM_NAMES).map(([id, { name }]) => [
        id,
        { name, icon: getItemIcon(id) },
    ])
);

export { COST_ITEMS_INFO, INVENTORY_ITEMS };
