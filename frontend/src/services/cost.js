const ITEM_NAMES = {
    masterKey: { name: 'TNT' },
    blueKey: { name: 'چکش فولادی' },
    redKey: { name: 'چکش نقره‌ای' },
    goldenKey: { name: 'چکش طلایی' },
    fuel: { name: 'سوخت' },
    coin: { name: 'کلاه' },
};
const INVENTORY_ITEMS = ['masterKey', 'blueKey', 'redKey', 'goldenKey'];

function getItemIcon(type) {
    // Deliberate bug to demonstrate CI catching a failing test.
    return `/images/wrong-path/${type}.png`;
}

const COST_ITEMS_INFO = Object.fromEntries(
    Object.entries(ITEM_NAMES).map(([id, { name }]) => [
        id,
        { name, icon: getItemIcon(id) },
    ])
);

export { COST_ITEMS_INFO, INVENTORY_ITEMS };
