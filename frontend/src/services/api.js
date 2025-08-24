import axios from 'axios'

// Create Axios instance for Picsum API
const apiClient = axios.create({
  baseURL: 'https://picsum.photos/v2',
  timeout: 10000,
})

// Fetch list of photos (default 30 photos)
export const fetchPhotos = async (page = 1, limit = 10) => {
  const response = await apiClient.get('/list', {
    params: { page, limit },
  })
  return response.data
}

// Fetch single photo info by id (just metadata)
export const fetchPhotoById = async (id) => {
  const response = await apiClient.get(`/list/${id}`)
  return response.data
}


// --- NEW: Game API Client ---
// Create a separate Axios instance for the game API
const gameApiClient = axios.create({
  baseURL: 'http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1',
  timeout: 10000,
});

// Intercept requests to automatically add the auth token for the game API
gameApiClient.interceptors.request.use(
  (config) => {
    const token = getToken();
    if (token) {
      config.headers.Authorization = token;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Intercept responses to handle the custom format {ok, error, result} for the game API
gameApiClient.interceptors.response.use(
  (response) => {
    if (response.data && typeof response.data.ok !== 'undefined' && !response.data.ok) {
      return Promise.reject(new Error(response.data.error || 'API request failed'));
    }
    return response.data.result;
  },
  (error) => {
    if (error.response && error.response.data && error.response.data.error) {
      return Promise.reject(new Error(error.response.data.error));
    }
    return Promise.reject(error);
  }
);


// --- NEW: Exported Functions for the Game API ---

// Perform user login
export const login = async (username, password) => {
  logout();
  const response = await gameApiClient.post('/login', { username, password });
  return response; // response is already the "result" object { token: "..." }
};

// Get current player state (authenticated)
// THE ONLY CHANGE IS HERE: Changed from .post to .get
export const getPlayer = async () => {
  const response = await gameApiClient.get('/player');
  return response;
};

// CHECL TRAVEL
export const checkTravel = async (fromId, toId) => {
  const response = await gameApiClient.post('/travel_check', {
    fromIsland: fromId,
    toIsland: toId,
  });
  return response;
};

//TRAVEL
export const travelTo = async (fromIsland, toIsland) => {
  const response = await gameApiClient.post('/travel', {
    fromIsland,
    toIsland,
  });
  return response;
};

// REFUEL CHECK
export const refuelCheck = async () => {
  const response = await gameApiClient.post('/refuel_check', {});
  return response;
};

// REFUEL
export const buyFuel = async (amount) => {
  const response = await gameApiClient.post('/refuel', {
    amount
  });
  return response;
};

// Get current user info (authenticated)
export const getMe = async () => {
  const response = await gameApiClient.get('/me');
  return response;
};

// Fetch territory data by id
export const getTerritory = async (territoryID) => {
  const response = await gameApiClient.get(`/territories/${territoryID}`);
  return response;
};

// Fetch island data by id
export const getIsland = async (islandID) => {
  const response = await gameApiClient.get(`/islands/${islandID}`);
  return response;
};

// Submit an answer for a challenge
export const submitAnswer = async (inputID, formData) => {
  const response = await gameApiClient.post(`/answer/${inputID}`, formData);
  return response;
};

// Helper for auth state
export const isLoggedIn = () => {
  return !!(getToken());
};

export const getToken = () => {
  return localStorage.getItem('authToken') || sessionStorage.getItem('authToken');
}

export const logout = () => {
  localStorage.removeItem('authToken');
  sessionStorage.removeItem('authToken');
};