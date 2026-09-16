import api, { unwrap } from './api'

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api/v1'

export const authService = {
  async login(email, password) {
    const res = await api.post('/auth/login', { email, password })
    return unwrap(res) // { token, user }
  },
  async register(payload) {
    const res = await api.post('/auth/register', payload)
    return unwrap(res)
  },
  async forgotPassword(email) {
    const res = await api.post('/auth/forgot-password', { email })
    return unwrap(res)
  },
  async resetPassword(token, new_password) {
    const res = await api.post('/auth/reset-password', { token, new_password })
    return unwrap(res)
  },
  async changePassword(current_password, new_password) {
    const res = await api.patch('/auth/change-password', { current_password, new_password })
    return unwrap(res)
  },
  ssoLoginUrl() {
    return `${API_BASE}/auth/sso/login`
  },
}
