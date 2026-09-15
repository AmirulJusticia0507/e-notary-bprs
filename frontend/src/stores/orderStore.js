import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { documentService, financingService, notaryService, orderService } from '../services/legalOrderService'

export const useOrderStore = defineStore('order', () => {
  const orders = ref([])
  const detail = ref(null)
  const logs = ref([])
  const documents = ref([])
  const financings = ref([])
  const notaries = ref([])
  const loading = ref(false)
  const error = ref('')

  const counts = computed(() => ({
    total: orders.value.length,
    pending: orders.value.filter((o) => o.status === 'pending').length,
    inProgress: orders.value.filter((o) => o.status === 'in_progress').length,
    completed: orders.value.filter((o) => o.status === 'completed').length,
    rejected: orders.value.filter((o) => o.status === 'rejected').length,
  }))

  async function fetchOrders() {
    loading.value = true
    error.value = ''
    try {
      orders.value = (await orderService.list()) ?? []
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Gagal memuat order'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchDetail(id) {
    loading.value = true
    try {
      detail.value = await orderService.getById(id)
      const [l, d] = await Promise.all([orderService.logs(id), documentService.listByOrder(id)])
      logs.value = l ?? []
      documents.value = d ?? []
    } finally {
      loading.value = false
    }
  }

  async function createOrder(payload) {
    return orderService.create(payload)
  }

  async function updateStatus(id, status, changedBy) {
    await orderService.updateStatus(id, status, changedBy)
    await fetchDetail(id)
  }

  async function fetchReferences() {
    const [f, n] = await Promise.all([financingService.list(), notaryService.list()])
    financings.value = f ?? []
    notaries.value = (n ?? []).filter((x) => x.is_available !== false)
  }

  return {
    orders,
    detail,
    logs,
    documents,
    financings,
    notaries,
    loading,
    error,
    counts,
    fetchOrders,
    fetchDetail,
    createOrder,
    updateStatus,
    fetchReferences,
  }
})
