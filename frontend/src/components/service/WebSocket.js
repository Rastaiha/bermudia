import { ref, onUnmounted } from 'vue';
import { getToken } from '@/services/api';

export const usePlayerWebSocket = (playerRef, nodes) => {
  const ws = ref(null);
  const reconnectInterval = 10000; // 10 seconds
  let reconnectTimeout = null;

  const connect = () => {
    const token = getToken();
    const wsUrl = `ws://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1/events?token=${token}`;
    ws.value = new WebSocket(wsUrl);

    ws.value.onopen = () => console.log('WebSocket connected!');
    

    ws.value.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.playerUpdate) {
          const update = data.playerUpdate;
          const island = nodes.value.find(node => node.id === update.player.atIsland);
          playerRef.value.atTerritory = update.player.atTerritory;
          playerRef.value.atIsland = island;
          playerRef.value.fuel = update.player.fuel;
          playerRef.value.fuelCap = update.player.fuelCap;
          console.log('Player updated via WebSocket:', playerRef.value);
        }
      } catch (err) {
        console.error('Failed to parse WebSocket message:', err);
      }
    };

    ws.value.onerror = (err) => {
      console.error('WebSocket error:', err);
    };

    ws.value.onclose = () => {
      console.warn('WebSocket closed. Reconnecting in 10s...');
      reconnectTimeout = setTimeout(connect, reconnectInterval);
    };
  };

  connect();

  onUnmounted(() => {
    if (ws.value) ws.value.close();
    if (reconnectTimeout) clearTimeout(reconnectTimeout);
  });
};
