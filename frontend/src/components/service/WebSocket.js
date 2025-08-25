// frontend/src/components/service/WebSocket.js

import { ref, onMounted, onUnmounted, watch } from 'vue';
import { getToken } from '@/services/api.js'; // Also added .js here for consistency
import { API_ENDPOINTS } from '@/services/apiConfig.js'; // Fix: Added .js extension

function connectWebSocket(player, nodes) {
  const token = getToken();
  if (!token) {
    console.error("No auth token found, WebSocket connection aborted.");
    return;
  }

  const socket = new WebSocket(`${API_ENDPOINTS.events}?token=${token}`);

  socket.onopen = () => {
    console.log("WebSocket connection established.");
  };

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      console.log("WebSocket message received:", data);

      if (data.type === "player-update" && player.value) {
        const payload = data.payload;
        if (payload.fuel !== undefined) player.value.fuel = payload.fuel;

        if (payload.atIsland) {
          const newIsland = nodes.value.find(node => node.id === payload.atIsland);
          if (newIsland) {
            player.value.atIsland = newIsland;
          }
        }
      }
    } catch (error) {
      console.error("Error parsing WebSocket message:", error);
    }
  };

  socket.onclose = (event) => {
    console.log("WebSocket connection closed:", event.reason);
  };

  socket.onerror = (error) => {
    console.error("WebSocket error:", error);
  };

  return socket;
}

export function usePlayerWebSocket(player, nodes) {
  let socket = null;

  const setupWebSocket = () => {
    if (player.value && nodes.value.length > 0) {
      if (socket) {
        socket.close();
      }
      socket = connectWebSocket(player, nodes);
    }
  };

  onMounted(setupWebSocket);

  watch([player, nodes], setupWebSocket, { deep: true });

  onUnmounted(() => {
    if (socket) {
      socket.close();
    }
  });
}