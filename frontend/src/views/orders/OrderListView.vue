<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import DataTable from '../../components/tables/DataTable.vue'
import StatusBadge from '../../components/tables/StatusBadge.vue'
import TablePagination from '../../components/tables/TablePagination.vue'
import { useOrderStore } from '../../stores/orderStore'

const store = useOrderStore()
const router = useRouter()
const query = ref('')
const statusFilter = ref('')
const page = ref(1)
const perPage = 10

onMounted(() => store.fetchOrders().catch(() => {}))

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  return store.orders.filter((o) => {
    if (statusFilter.value && o.status !== statusFilter.value) return false
    if (!q) return true
    return [o.order_number, o.customer_name, o.notary_name].some((v) => String(v ?? '').toLowerCase().includes(q))
  })
})

const paged = computed(() => filtered.value.slice((page.value - 1) * perPage, page.value * perPage))
</script>

<template>
  <div class="space-y-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h2 class="text-lg font-bold text-slate-900">Legal Orders</h2>
      <button @click="router.push({ name: 'order-create' })" class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700">
        + Order Baru
      </button>
    </div>
    <div class="flex flex-wrap gap-3">
      <input
        v-model="query"
        placeholder="Cari order / nasabah / notaris…"
        class="w-72 rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
      />
      <select v-model="statusFilter" class="rounded-md border border-slate-300 px-3 py-2 text-sm">
        <option value="">Semua status</option>
        <option value="pending">pending</option>
        <option value="in_progress">in_progress</option>
        <option value="completed">completed</option>
        <option value="rejected">rejected</option>
      </select>
    </div>
    <LoadingSpinner v-if="store.loading && !store.orders.length" />
    <template v-else>
      <DataTable
        :columns="[
          { key: 'order_number', label: 'Order', mono: true },
          { key: 'customer_name', label: 'Nasabah' },
          { key: 'notary_name', label: 'Notaris' },
          { key: 'assigned_name', label: 'PIC' },
          { key: 'status', label: 'Status' },
          { key: 'sla_deadline', label: 'SLA' },
        ]"
        :rows="paged"
        @row-click="(r) => router.push({ name: 'order-detail', params: { id: r.id } })"
      >
        <template #cell-status="{ value }"><StatusBadge :status="value" /></template>
        <template #cell-sla_deadline="{ value }">{{ value ? new Date(value).toLocaleDateString('id-ID') : '—' }}</template>
      </DataTable>
      <TablePagination v-model:page="page" :per-page="perPage" :total="filtered.length" />
    </template>
  </div>
</template>
