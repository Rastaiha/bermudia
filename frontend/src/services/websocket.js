import { onUnmounted, watch } from 'vue';
import { getToken, getInboxMessages } from '@/services/api/index.js';
import { API_ENDPOINTS } from '@/services/api/config.js';
import emitter from '@/services/eventBus.js';
import { notificationService } from '@/services/notificationService.js';
import { uiState } from '@/services/uiState.js';
import { messages as inboxMessages } from '@/services/inboxWebsocket.js';

export function usePlayerWebSocket(player, territoryId, route, router) {
    let socket = null;
    let reconnectTimeoutId = null;
    let reconnectAttempts = 0;
    const maxReconnectAttempts = 10;
    const baseReconnectDelay = 1000;

    const disconnect = () => {
        if (reconnectTimeoutId) clearTimeout(reconnectTimeoutId);

        if (socket) {
            socket.onclose = null;
            socket.onmessage = null;
            socket.onerror = null;
            socket.onopen = null;
            socket.close(1000, 'Connection closed intentionally by client');
            socket = null;
            console.log('WebSocket connection cleanly disconnected.');
        }
    };

    const connect = () => {
        if (socket) return;

        const token = getToken();
        if (!token) {
            console.error(
                'WebSocket: No auth token found, connection aborted.'
            );
            return;
        }

        console.log('Attempting to connect WebSocket...');
        socket = new WebSocket(`${API_ENDPOINTS.events}?token=${token}`);

        socket.onopen = () => {
            console.log('WebSocket connection established.');
            reconnectAttempts = 0;
        };

        socket.onmessage = async event => {
            try {
                const data = JSON.parse(event.data);
                console.log('WebSocket message received:', data);

                if (data.playerUpdate) {
                    const reason = data.playerUpdate.reason;
                    if (
                        reason === 'correction' ||
                        reason === 'ownOfferAccepted'
                    ) {
                        try {
                            const result = await getInboxMessages(null, 20);
                            const allMessages =
                                result?.messages || result || [];

                            inboxMessages.value = allMessages;

                            notificationService.setReceivedMessages(
                                allMessages
                            );

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

                    const oldPlayerState = JSON.parse(
                        JSON.stringify(player.value)
                    );
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
                }
            } catch (error) {
                console.error('Error parsing WebSocket message:', error);
            }
        };

        socket.onclose = event => {
            console.log('WebSocket connection closed. Code:', event.code);
            socket = null;

            if (event.code !== 1000) {
                scheduleReconnect();
            }
        };

        socket.onerror = error => {
            console.error('WebSocket error:', error);
        };
    };

    const scheduleReconnect = () => {
        if (reconnectAttempts >= maxReconnectAttempts) {
            console.error(
                'Max WebSocket reconnection attempts reached. Giving up.'
            );
            return;
        }
        const delay = Math.min(
            baseReconnectDelay * Math.pow(2, reconnectAttempts),
            30000
        );
        reconnectAttempts++;

        console.log(
            `Scheduling WebSocket reconnection attempt ${reconnectAttempts} in ${delay}ms`
        );
        reconnectTimeoutId = setTimeout(connect, delay);
    };

    watch(
        player,
        newPlayer => {
            if (newPlayer && !socket) {
                connect();
            } else if (!newPlayer && socket) {
                disconnect();
            }
        },
        { immediate: true }
    );

    onUnmounted(() => {
        disconnect();
    });
}
