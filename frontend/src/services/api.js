// frontend/src/services/api.js

import { API_ENDPOINTS } from './apiConfig.js';

// --- Helper Functions ---
const getAuthHeaders = () => {
  const token = getToken();
  return token ? { 'Authorization': `Bearer ${token}` } : {};
};

const handleResponse = async (response) => {
  const data = await response.json();
  if (!response.ok || (data.ok === false)) {
    throw new Error(data.error || 'An unknown error occurred');
  }
  return data.result;
};

// --- Token Management ---
export const getToken = () => localStorage.getItem('authToken') || sessionStorage.getItem('authToken');

export const logout = () => {
  localStorage.removeItem('authToken');
  sessionStorage.removeItem('authToken');
};

// --- API Functions ---
export const login = async (username, password) => {
  const response = await fetch(API_ENDPOINTS.login, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  return handleResponse(response);
};

export const getMe = async () => {
  const response = await fetch(API_ENDPOINTS.getMe, {
    method: 'GET',
    headers: getAuthHeaders(),
  });
  return handleResponse(response);
};

export const getPlayer = async () => {
  const response = await fetch(API_ENDPOINTS.getPlayer, {
    method: 'GET',
    headers: getAuthHeaders(),
  });
  return handleResponse(response);
};

export const getTerritory = async (id) => {
  const response = await fetch(API_ENDPOINTS.getTerritory(id), {
    headers: getAuthHeaders(),
  });
  return handleResponse(response);
};

export const travelCheck = async (from, dest) => {
  const response = await fetch(API_ENDPOINTS.travelCheck, {
    method: 'POST',
    headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify({ fromIsland: from, toIsland: dest }),
  });
  return handleResponse(response);
};

export const anchorCheck = async (island) => {
  const response = await fetch(API_ENDPOINTS.anchorCheck, {
    method: 'POST',
    headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify({ fromIsland: island}),
  });
  return handleResponse(response);
};

export const travelTo = async (from, dest) => {
  const response = await fetch(API_ENDPOINTS.travelTo, {
    method: 'POST',
    headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify({ fromIsland: from, toIsland: dest }),
  });
  return handleResponse(response);
};

export const dropAnchorAtIsland = async (currentIsland) => {
  const response = await fetch(API_ENDPOINTS.dropAnchor, {
    method: 'POST',
    headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify({ island: currentIsland}),
  });
  return handleResponse(response);
};

export const getIsland = async (id) => {
  const response = await fetch(API_ENDPOINTS.getIsland(id), {
    headers: getAuthHeaders(),
  });
  return handleResponse(response);
};

export const submitAnswer = async (id, formData) => {
  const response = await fetch(API_ENDPOINTS.submitAnswer(id), {
    method: 'POST',
    headers: getAuthHeaders(),
    body: formData,
  });
  return handleResponse(response);
};

export const refuelCheck = async () => {
  const response = await fetch(API_ENDPOINTS.refuelCheck, {
    method: 'POST',
    headers: getAuthHeaders(),
  });
  return handleResponse(response);
};

export const buyFuel = async (amount) => {
  const response = await fetch(API_ENDPOINTS.buyFuel, {
    method: 'POST',
    headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify({ amount }),
  });
  return handleResponse(response);
};
