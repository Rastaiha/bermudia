// Enumerates image assets that already exist in the app's `public/` directory.
// The game references these by absolute public path (e.g.
// `/images/islands/educational/10.png`), so the picker must emit that exact
// form. We use Vite's import.meta.glob against the source `public/` tree at
// build time to discover what's available, then map each key back to its
// public URL.

const iconModules = import.meta.glob(
    '/public/images/islands/**/*.{png,jpg,jpeg,webp,svg}',
    { eager: true, query: '?url', import: 'default' }
);

const backgroundModules = import.meta.glob(
    '/public/images/backgrounds/territory/*.{png,jpg,jpeg,webp}',
    { eager: true, query: '?url', import: 'default' }
);

// Convert a glob key like `/public/images/islands/educational/10.png` into the
// public path the game stores: `/images/islands/educational/10.png`.
const toPublicPath = key => key.replace(/^\/public/, '');

// A natural sort so "2.png" comes before "10.png".
const naturalCompare = (a, b) =>
    a.localeCompare(b, undefined, { numeric: true, sensitivity: 'base' });

const buildIslandIcons = () => {
    // Group by island type (the sub-directory under islands/).
    const groups = {};
    for (const key of Object.keys(iconModules)) {
        const path = toPublicPath(key);
        const m = path.match(/\/images\/islands\/([^/]+)\//);
        const type = m ? m[1] : 'other';
        (groups[type] ||= []).push(path);
    }
    for (const type of Object.keys(groups)) {
        groups[type].sort(naturalCompare);
    }
    return groups;
};

const buildBackgrounds = () =>
    Object.keys(backgroundModules).map(toPublicPath).sort(naturalCompare);

// Categorized island icons: { challenge: [...], educational: [...], ... }
export const islandIconGroups = buildIslandIcons();

// Flat list of every island icon (for pickers that don't care about type).
export const allIslandIcons = Object.values(islandIconGroups)
    .flat()
    .sort(naturalCompare);

// Territory background images.
export const territoryBackgrounds = buildBackgrounds();

// The island type names we discovered, in a stable order.
export const islandTypes = Object.keys(islandIconGroups).sort(naturalCompare);
