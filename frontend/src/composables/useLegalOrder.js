import { computed } from 'vue'
import { useOrderStore } from '../stores/orderStore'

const NEXT_STATUS = {
  pending: ['in_progress'],
  in_progress: ['completed', 'rejected'],
  completed: [],
  rejected: [],
}

// useLegalOrder: transisi status yang diizinkan + flag SLA.
export function useLegalOrder() {
  const store = useOrderStore()

  const allowedTransitions = computed(() => (status) => NEXT_STATUS[status] ?? [])

  function isOverdue(order) {
    if (!order?.sla_deadline) return false
    return new Date(order.sla_deadline).getTime() < Date.now() && order.status !== 'completed'
  }

  return { store, allowedTransitions, isOverdue, NEXT_STATUS }
}
