<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import DataTable from '../../components/tables/DataTable.vue'
import StatusBadge from '../../components/tables/StatusBadge.vue'
import TablePagination from '../../components/tables/TablePagination.vue'
import AppButton from '../../components/common/AppButton.vue'
import { useOrderStore } from '../../stores/orderStore'

const store = useOrderStore()
const router = useRouter()
const query = ref('')
const statusFilter = ref('')
const page = ref(1)
const perPage = 10

onMounted(() => store.fetchOrders().catch(() => {}))

const filtered = computed(() => {
  const searchText = query.value.trim().toLowerCase()
  return store.orders.filter((order) => {
    if (statusFilter.value && order.status !== statusFilter.value) return false
    if (!searchText) return true
    return [order.order_number, order.customer_name, order.notary_name].some((value) => String(value ?? '').toLowerCase().includes(searchText))
  })
})

const paged = computed(() => filtered.value.slice((page.value - 1) * perPage, page.value * perPage))
</script>

<template>
  <div class="space-y-6">
    <div class="page-header">
      <div>
        <p class="text-[11px] font-bold uppercase tracking-[0.16em] text-teal-600">Workflow</p>
        <h1 class="page-title">Legal Orders</h1>
        <p class="page-subtitle">Kelola dan lacak seluruh order legalitas akad dalam satu tampilan.</p>
      </div>
      <AppButton @click="router.push({ name: 'order-create' })" class="!min-h-10">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
        Order Baru
      </AppButton>
    </div>

    <section class="panel p-3 sm:p-4">
      <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div class="relative flex-1">
          <svg class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" width="16" height="16" viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle cx="11" cy="11" r="6.5" stroke="currentColor" stroke-width="1.8"/><path d="m16 16 4 4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
          <input v-model="query" type="search" placeholder="Cari order, nasabah, atau notaris…" class="field-input !pl-9" />
        </div>
        <select v-model="statusFilter" class="field-input !w-full !min-w-[11rem] lg:!w-52">
          <option value="">Semua status</option>
          <option value="pending">Pending</option>
          <option value="in_progress">In Progress</option>
          <option value="completed">Completed</option>
          <option value="rejected">Rejected</option>
        </select>
      </div>
    </section>

    <LoadingSpinner v-if="store.loading && !store.orders.length" label="Memuat order…" />
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
        @row-click="(row) => router.push({ name: 'order-detail', params: { id: row.id } })"
      >
        <template #cell-status="{ value }"><StatusBadge :status="value" /></template>
        <template #cell-sla_deadline="{ value }">{{ value ? new Date(value).toLocaleDateString('id-ID') : '—' }}</template>
      </DataTable>
      <TablePagination v-model:page="page" :per-page="perPage" :total="filtered.length" />
    </template>
  </div>
</template>
