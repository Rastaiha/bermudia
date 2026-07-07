import { ref } from 'vue';
import {
    getPlayer,
    getMe,
    getToken,
    getTerritory,
} from '@/services/api/index.js';
import { logger } from '@/services/logger.js';

const calculateViewBox = (islands, padding = 0.1) => {
    if (!islands || islands.length === 0) return '0 0 1 1';
    const bounds = islands.reduce(
        (acc, island) => ({
            minX: Math.min(acc.minX, island.x - island.width / 2),
            maxX: Math.max(acc.maxX, island.x + island.width / 2),
            minY: Math.min(acc.minY, island.y - island.height / 2),
            maxY: Math.max(acc.maxY, island.y + island.height / 2),
        }),
        { minX: Infinity, maxX: -Infinity, minY: Infinity, maxY: -Infinity }
    );
    const { minX, minY, maxX, maxY } = bounds;
    return `${minX - padding} ${minY - padding} ${maxX - minX + padding * 2} ${
        maxY - minY + padding * 2
    }`;
};

export function useTerritoryData(router, { onBeforeLoad, onLoaded } = {}) {
    const territoryId = ref(null);
    const islands = ref([]);
    const refuelIslands = ref([]);
    const terminalIslands = ref([]);
    const edges = ref([]);
    const territoryName = ref('');
    const player = ref(null);
    const username = ref('...');
    const backgroundImage = ref('');
    const dynamicViewBox = ref('0 0 1 1');
    const isLoading = ref(true);
    const loadingProgress = ref(0);

    const fetchTerritoryData = async id => getTerritory(id);

    const fetchPlayerAndUserData = async () => {
        if (!getToken()) {
            router.push({ name: 'Login' });
            throw new Error('User not authenticated');
        }
        const [playerData, meData] = await Promise.all([getPlayer(), getMe()]);
        return { playerData, meData };
    };

    const setupTerritoryData = territoryData => {
        backgroundImage.value = territoryData.backgroundAsset;
        territoryName.value = territoryData.name;
        islands.value = territoryData.islands;
        edges.value = territoryData.edges;
        refuelIslands.value = territoryData.refuelIslands;
        terminalIslands.value = territoryData.terminalIslands;
        dynamicViewBox.value = calculateViewBox(territoryData.islands);
    };

    const setupPlayerAndUserData = (playerAndUserData, currentTerritoryId) => {
        if (!playerAndUserData) return;
        const { playerData, meData } = playerAndUserData;

        if (
            playerData.atTerritory.toString() !== currentTerritoryId.toString()
        ) {
            router.push({
                name: 'Territory',
                params: { id: playerData.atTerritory },
            });
            throw new Error('Redirecting to correct territory');
        }

        username.value = meData.name;
        player.value = playerData;
    };

    const loadPageData = async id => {
        if (!id) return;
        isLoading.value = true;
        loadingProgress.value = 0;
        territoryId.value = id;
        onBeforeLoad?.();

        let progressInterval = null;

        try {
            progressInterval = setInterval(() => {
                if (loadingProgress.value < 90) {
                    loadingProgress.value += 5;
                }
            }, 100);

            const [territoryData, playerAndUserData] = await Promise.all([
                fetchTerritoryData(id),
                fetchPlayerAndUserData(),
            ]);

            loadingProgress.value = 100;

            setupTerritoryData(territoryData);
            setupPlayerAndUserData(playerAndUserData, id);
        } catch (error) {
            logger.error('Failed to load page data:', error.message);
            if (
                !error.message.includes('authenticated') &&
                !error.message.includes('Redirecting')
            ) {
                router.push({ name: 'Login' });
            }
        } finally {
            clearInterval(progressInterval);
            setTimeout(() => {
                isLoading.value = false;
                onLoaded?.();
            }, 500);
        }
    };

    return {
        territoryId,
        islands,
        refuelIslands,
        terminalIslands,
        edges,
        territoryName,
        player,
        username,
        backgroundImage,
        dynamicViewBox,
        isLoading,
        loadingProgress,
        loadPageData,
    };
}
