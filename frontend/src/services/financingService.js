import api, { unwrap } from './api'

export const financingService = {
  async listMine() {
    const res = await api.get('/financings/mine')
    return unwrap(res)
  },
  async apply(payload) {
    const res = await api.post('/financings/apply', payload)
    return unwrap(res)
  },
}
