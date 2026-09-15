import api, { unwrap } from './api'

export const authService = {
  async login(email, password) {
    const res = await api.post('/auth/login', { email, password })
    return unwrap(res) // { token, user }
  },
  async register(payload) {
    const res = await api.post('/auth/register', payload)
    return unwrap(res)
  },
}
