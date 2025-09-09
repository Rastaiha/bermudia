import { onUnmounted, watch } from 'vue';
import { getToken } from '@/services/api.js';
import { API_ENDPOINTS } from '@/services/apiConfig.js';
import emitter from '@/services/eventBus.js';

export function usePlayerWebSocket(player, territoryId, router) {
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

        socket.onmessage = event => {
            try {
                const data = JSON.parse(event.data);
                console.log('WebSocket message received:', data);

                if (data.playerUpdate) {
                    const oldPlayerState = JSON.parse(
                        JSON.stringify(player.value)
                    );
                    const newPlayerState = data.playerUpdate.player;

                    player.value = newPlayerState;

                    if (data.playerUpdate.reason === 'unlockTreasure') {
                        emitter.emit('treasure-unlocked', {
                            oldPlayerState,
                            newPlayerState,
                        });
                    }

                    if (
                        territoryId &&
                        territoryId.value &&
                        newPlayerState.atTerritory != territoryId.value
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
