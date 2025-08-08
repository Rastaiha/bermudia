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
