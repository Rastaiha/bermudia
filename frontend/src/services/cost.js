const ITEM_NAMES = {
    masterKey: { name: 'شاه‌کلید' },
    blueKey: { name: 'کلید آبی' },
    redKey: { name: 'کلید قرمز' },
    goldenKey: { name: 'کلید طلایی' },
    fuel: { name: 'سوخت' },
    coin: { name: 'سکه' },
    book: { name: 'کتاب' },
};

function getItemIcon(type) {
    return `/images/icons/${type}.png`;
}

const COST_ITEMS_INFO = Object.fromEntries(
    Object.entries(ITEM_NAMES).map(([id, { name }]) => [
        id,
        { name, icon: getItemIcon(id) },
    ])
);
export default COST_ITEMS_INFO;
