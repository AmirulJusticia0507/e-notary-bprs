<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import DataTable from '../../components/tables/DataTable.vue'
import StatusBadge from '../../components/tables/StatusBadge.vue'
import { useOrderStore } from '../../stores/orderStore'
import { useUserStore } from '../../stores/userStore'

const store = useOrderStore()
const userStore = useUserStore()
const router = useRouter()

const cards = [
  { label: 'Total Order', key: 'total', style: 'border-slate-200' },
  { label: 'Pending', key: 'pending', style: 'border-amber-200' },
  { label: 'In Progress', key: 'inProgress', style: 'border-blue-200' },
  { label: 'Completed', key: 'completed', style: 'border-emerald-200' },
]

onMounted(() => store.fetchOrders().catch(() => {}))
</script>

<template>
  <div class="space-y-5">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-bold text-slate-900">Dashboard</h2>
      <button
        v-if="userStore.role !== 'notary'"
        @click="router.push({ name: 'order-create' })"
        class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
      >
        + Order Baru
      </button>
    </div>
    <LoadingSpinner v-if="store.loading && !store.orders.length" />
    <div v-else class="grid grid-cols-2 gap-4 lg:grid-cols-4">
      <div v-for="card in cards" :key="card.key" class="rounded-lg border bg-white p-5 shadow-sm" :class="card.style">
        <p class="text-xs font-semibold text-slate-500">{{ card.label }}</p>
        <p class="mt-1 text-2xl font-bold text-slate-900">{{ store.counts[card.key] }}</p>
      </div>
    </div>
    <div class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
      <h3 class="mb-3 text-sm font-bold text-slate-900">Order terbaru</h3>
      <DataTable
        :columns="[
          { key: 'order_number', label: 'Order', mono: true },
          { key: 'customer_name', label: 'Nasabah' },
          { key: 'notary_name', label: 'Notaris' },
          { key: 'status', label: 'Status' },
        ]"
        :rows="store.orders.slice(0, 5)"
        @row-click="(r) => router.push({ name: 'order-detail', params: { id: r.id } })"
      >
        <template #cell-status="{ value }"><StatusBadge :status="value" /></template>
      </DataTable>
    </div>
  </div>
</template>
