import { onUnmounted, watch } from 'vue';
import { getToken, getInboxMessages } from '@/services/api/index.js';
import { API_ENDPOINTS } from '@/services/api/config.js';
import emitter from '@/services/eventBus.js';
import { notificationService } from '@/services/notificationService.js';
import { uiState } from '@/services/uiState.js';
import { messages as inboxMessages } from '@/services/inboxWebsocket.js';
import { createReconnectingSocket } from '@/services/ws/reconnectingSocket.js';

export function usePlayerWebSocket(player, territoryId, route, router) {
    const handleMessage = async data => {
        console.log('WebSocket message received:', data);

        if (!data.playerUpdate) return;

        const reason = data.playerUpdate.reason;
        if (reason === 'correction' || reason === 'ownOfferAccepted') {
            try {
                const result = await getInboxMessages(null, 20);
                const allMessages = result?.messages || result || [];

                inboxMessages.value = allMessages;

                notificationService.setReceivedMessages(allMessages);

                if (uiState.isInboxOpen) {
                    notificationService.markAllAsSeen();
                }
            } catch (apiError) {
                console.error(
                    'Failed to fetch inbox messages after playerUpdate event:',
                    apiError
                );
            }
        }

        const oldPlayerState = JSON.parse(JSON.stringify(player.value));
        const newPlayerState = data.playerUpdate.player;

        player.value = newPlayerState;

        if (reason === 'unlockTreasure') {
            emitter.emit('treasure-unlocked', {
                oldPlayerState,
                newPlayerState,
            });
        }

        if (
            router &&
            route &&
            territoryId &&
            territoryId.value &&
            newPlayerState.atTerritory != territoryId.value &&
            !route.params.islandId
        ) {
            router.push({
                name: 'Territory',
                params: { id: newPlayerState.atTerritory },
            });
        }
    };

    const { connect, disconnect } = createReconnectingSocket({
        buildUrl: () => {
            const token = getToken();
            return token ? `${API_ENDPOINTS.events}?token=${token}` : null;
        },
        onMessage: handleMessage,
        label: 'WebSocket',
    });

    watch(
        player,
        newPlayer => {
            if (newPlayer) {
                connect();
            } else {
                disconnect();
            }
        },
        { immediate: true }
    );

    onUnmounted(() => {
        disconnect();
    });
}
