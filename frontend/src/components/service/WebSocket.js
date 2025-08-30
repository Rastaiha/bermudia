// frontend/src/components/service/WebSocket.js

import { onMounted, onUnmounted } from 'vue';
import { getToken } from '@/services/api.js';
import { API_ENDPOINTS } from '@/services/apiConfig.js';

function connectWebSocket(player, reconnectCallback) {
  const token = getToken();
  if (!token) {
    console.error("No auth token found, WebSocket connection aborted.");
    return null;
  }

  const socket = new WebSocket(`${API_ENDPOINTS.events}?token=${token}`);

  socket.onopen = () => {
    console.log("WebSocket connection established.");
    // Reset reconnect attempts on successful connection
    if (reconnectCallback) {
      reconnectCallback.resetAttempts();
    }
  };

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      console.log("WebSocket message received:", data);

      if (data.playerUpdate) {
        player.value = data.playerUpdate.player;
      }
    } catch (error) {
      console.error("Error parsing WebSocket message:", error);
    }
  };

  socket.onclose = (event) => {
    console.log("WebSocket connection closed:", event.reason, "Code:", event.code);
    // Attempt to reconnect unless it was a clean close (1000) or manual close
    if (event.code !== 1000 && reconnectCallback) {
      reconnectCallback.scheduleReconnect();
    }
  };

  socket.onerror = (error) => {
    console.error("WebSocket error:", error);
  };

  return socket;
}

export function usePlayerWebSocket(player) {
  let socket = null;
  let reconnectTimeoutId = null;
  let reconnectAttempts = 0;
  const maxReconnectAttempts = 10;
  const baseReconnectDelay = 1000; // 1 second

  const cleanup = () => {
    if (reconnectTimeoutId) {
      clearTimeout(reconnectTimeoutId);
      reconnectTimeoutId = null;
    }
    if (socket) {
      // Close with code 1000 (normal closure) to prevent reconnection
      socket.close(1000, "Manual close");
      socket = null;
    }
  };

  const scheduleReconnect = () => {
    if (reconnectAttempts >= maxReconnectAttempts) {
      console.error("Max WebSocket reconnection attempts reached. Giving up.");
      return;
    }

    // Exponential backoff: 1s, 2s, 4s, 8s, 16s, 32s, then cap at 32s
    const delay = Math.min(baseReconnectDelay * Math.pow(2, reconnectAttempts), 32000);
    reconnectAttempts++;

    console.log(`Scheduling WebSocket reconnection attempt ${reconnectAttempts} in ${delay}ms`);
    
    reconnectTimeoutId = setTimeout(() => {
      setupWebSocket();
    }, delay);
  };

  const resetAttempts = () => {
    reconnectAttempts = 0;
  };

  const setupWebSocket = () => {
    if (socket) {
      cleanup();
    }
    
    socket = connectWebSocket(player, {
      scheduleReconnect,
      resetAttempts
    });
  };

  onMounted(setupWebSocket);

  onUnmounted(() => {
    cleanup();
  });
}
