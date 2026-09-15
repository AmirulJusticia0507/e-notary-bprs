import api, { unwrap } from './api'

export const orderService = {
  async list() {
    return unwrap(await api.get('/orders'))
  },
  async getById(id) {
    return unwrap(await api.get(`/orders/${id}`))
  },
  async listByNotary(notaryId) {
    return unwrap(await api.get(`/orders/notary/${notaryId}`))
  },
  async listByAssignee(userId) {
    return unwrap(await api.get(`/orders/assigned/${userId}`))
  },
  async create(payload) {
    return unwrap(await api.post('/orders', payload))
  },
  async updateStatus(id, status, changedBy) {
    return unwrap(await api.patch(`/orders/${id}/status`, { status, changed_by: changedBy }))
  },
  async logs(id) {
    return unwrap(await api.get(`/orders/${id}/logs`))
  },
}

export const financingService = {
  async list() {
    return unwrap(await api.get('/financings'))
  },
  async getById(id) {
    return unwrap(await api.get(`/financings/${id}`))
  },
  async create(payload) {
    return unwrap(await api.post('/financings', payload))
  },
  async sync(apps) {
    return unwrap(await api.post('/financings/sync', apps))
  },
}

export const notaryService = {
  async list() {
    return unwrap(await api.get('/notaries'))
  },
  async getById(id) {
    return unwrap(await api.get(`/notaries/${id}`))
  },
  async create(payload) {
    return unwrap(await api.post('/notaries', payload))
  },
  async update(id, payload) {
    return unwrap(await api.put(`/notaries/${id}`, payload))
  },
  async remove(id) {
    return unwrap(await api.delete(`/notaries/${id}`))
  },
}

export const documentService = {
  async listByOrder(orderId) {
    return unwrap(await api.get(`/documents/order/${orderId}`))
  },
  async getById(id) {
    return unwrap(await api.get(`/documents/${id}`))
  },
  async create(payload) {
    return unwrap(await api.post('/documents', payload))
  },
  async updateESignStatus(id, status) {
    return unwrap(await api.patch(`/documents/${id}/sign`, { status }))
  },
}
