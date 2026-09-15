import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api/v1',
  headers: { 'Content-Type': 'application/json' },
  timeout: 15000,
})

// Injeksi JWT ke setiap request terproteksi.
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('enotary_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Backend membungkus payload: { status, data?, error? }.
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status
    if (status === 401) {
      localStorage.removeItem('enotary_token')
      localStorage.removeItem('enotary_user')
      if (window.location.pathname !== '/login') {
        window.location.assign('/login')
      }
    }
    return Promise.reject(error)
  },
)

// Ambil `data` dari envelope backend { status, data, error }.
export function unwrap(response) {
  return response.data?.data ?? response.data
}

export function apiErrorMessage(error, fallback = 'Terjadi kesalahan') {
  return error.response?.data?.error ?? error.message ?? fallback
}

export default api
