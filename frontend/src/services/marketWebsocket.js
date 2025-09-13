import { onMounted, onUnmounted } from 'vue';
import { getToken } from '@/services/api/index.js';
import { API_ENDPOINTS } from '@/services/api/config.js';

export function useMarketWebSocket(
    mySyncTrade,
    otherSyncTrade,
    myOffers,
    otherOffers
) {
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
            console.log('Market WebSocket connection cleanly disconnected.');
        }
    };

    const connect = () => {
        if (socket) return;

        const token = getToken();
        if (!token) {
            console.error(
                'Market WebSocket: No auth token found, connection aborted.'
            );
            return;
        }

        console.log('Attempting to connect Market WebSocket...');
        socket = new WebSocket(`${API_ENDPOINTS.marketEvents}?token=${token}`);

        socket.onopen = () => {
            console.log('Market WebSocket connection established.');
            reconnectAttempts = 0;
        };

        socket.onmessage = event => {
            try {
                const data = JSON.parse(event.data);
                console.log('Market WebSocket message received:', data);

                if (data.sync) {
                    mySyncTrade.value = data.sync.offset;
                    otherSyncTrade.value = data.sync.offset;
                }
                if (data.new_offer) {
                    if (data.new_offer.offer.byMe) {
                        myOffers.value.unshift(data.new_offer.offer);
                        mySyncTrade.value = data.new_offer.offer.created_at;
                    } else {
                        otherOffers.value.unshift(data.new_offer.offer);
                        otherSyncTrade.value = data.new_offer.offer.created_at;
                    }
                }
                if (data.deleted_offer) {
                    if (data.deleted_offer.byMe) {
                        for (let offer of myOffers.value) {
                            if (offer.id == data.deleted_offer.offerID) {
                                myOffers.value.splice(
                                    myOffers.value.indexOf(offer),
                                    1
                                );
                                break;
                            }
                        }
                    } else {
                        for (let offer of otherOffers.value) {
                            if (offer.id == data.deleted_offer.offerID) {
                                otherOffers.value.splice(
                                    otherOffers.value.indexOf(offer),
                                    1
                                );
                                break;
                            }
                        }
                    }
                }
            } catch (error) {
                console.error('Error parsing Market WebSocket message:', error);
            }
        };

        socket.onclose = event => {
            console.log(
                'Market WebSocket connection closed. Code:',
                event.code
            );
            socket = null;

            if (event.code !== 1000) {
                scheduleReconnect();
            }
        };

        socket.onerror = error => {
            console.error('Market WebSocket error:', error);
        };
    };

    const scheduleReconnect = () => {
        if (reconnectAttempts >= maxReconnectAttempts) {
            console.error(
                'Max Market WebSocket reconnection attempts reached. Giving up.'
            );
            return;
        }
        const delay = Math.min(
            baseReconnectDelay * Math.pow(2, reconnectAttempts),
            30000
        );
        reconnectAttempts++;

        console.log(
            `Scheduling Market WebSocket reconnection attempt ${reconnectAttempts} in ${delay}ms`
        );
        reconnectTimeoutId = setTimeout(connect, delay);
    };

    onMounted(() => {
        if (!socket) {
            connect();
        }
    });

    onUnmounted(() => {
        disconnect();
    });
}
