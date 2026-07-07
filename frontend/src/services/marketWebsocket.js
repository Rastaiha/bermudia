import { onMounted, onUnmounted } from 'vue';
import { getToken } from '@/services/api/index.js';
import { API_ENDPOINTS } from '@/services/api/config.js';
import { createReconnectingSocket } from '@/services/ws/reconnectingSocket.js';
import { logger } from '@/services/logger.js';

export function useMarketWebSocket(
    mySyncTrade,
    otherSyncTrade,
    myOffers,
    otherOffers
) {
    const handleMessage = data => {
        logger.log('Market WebSocket message received:', data);

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
                        myOffers.value.splice(myOffers.value.indexOf(offer), 1);
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
    };

    const { connect, disconnect } = createReconnectingSocket({
        buildUrl: () => {
            const token = getToken();
            return token
                ? `${API_ENDPOINTS.marketEvents}?token=${token}`
                : null;
        },
        onMessage: handleMessage,
        label: 'Market WebSocket',
    });

    onMounted(() => {
        connect();
    });

    onUnmounted(() => {
        disconnect();
    });
}
