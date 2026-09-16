<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import DataTable from '../../components/tables/DataTable.vue'
import StatusBadge from '../../components/tables/StatusBadge.vue'
import AppButton from '../../components/common/AppButton.vue'
import { useOrderStore } from '../../stores/orderStore'
import { useUserStore } from '../../stores/userStore'

const store = useOrderStore()
const userStore = useUserStore()
const router = useRouter()

const cards = [
  { label: 'Total Order', key: 'total', color: '#0f766e', bg: 'rgba(20,184,166,.11)', glow: 'rgba(20,184,166,.1)', icon: 'folder' },
  { label: 'Pending', key: 'pending', color: '#b45309', bg: 'rgba(245,158,11,.12)', glow: 'rgba(245,158,11,.1)', icon: 'clock' },
  { label: 'In Progress', key: 'inProgress', color: '#0369a1', bg: 'rgba(14,165,233,.11)', glow: 'rgba(14,165,233,.1)', icon: 'activity' },
  { label: 'Completed', key: 'completed', color: '#047857', bg: 'rgba(16,185,129,.11)', glow: 'rgba(16,185,129,.1)', icon: 'check' },
]

onMounted(() => store.fetchOrders().catch(() => {}))
</script>

<template>
  <div class="space-y-6">
    <div class="page-header">
      <div>
        <p class="text-[11px] font-bold uppercase tracking-[0.16em] text-teal-600">Overview</p>
        <h1 class="page-title">Dashboard</h1>
        <p class="page-subtitle">Pantau order, SLA, dan progres legalitas akad hari ini.</p>
      </div>
      <AppButton v-if="userStore.role !== 'notary'" @click="router.push({ name: 'order-create' })" class="!min-h-10">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
        Order Baru
      </AppButton>
    </div>

    <div class="rounded-2xl border border-slate-200 bg-gradient-to-r from-teal-600 to-sky-600 p-5 text-white shadow-lg shadow-teal-900/10 sm:p-6">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div>
          <p class="text-xs font-semibold text-teal-100">Halo, {{ userStore.user?.full_name || 'Tim Legal' }}</p>
          <h2 class="mt-1 text-lg font-extrabold tracking-tight sm:text-xl">Semangat menyelesaikan akad hari ini.</h2>
          <p class="mt-1 max-w-xl text-xs leading-5 text-teal-50 sm:text-sm">Gunakan dashboard ini untuk melihat order yang perlu ditindaklanjuti dan menjaga SLA tetap berada pada jalurnya.</p>
        </div>
        <div class="grid grid-cols-3 gap-2 rounded-2xl border border-white/15 bg-white/10 p-3 backdrop-blur sm:min-w-[280px]">
          <div class="text-center"><p class="text-xl font-extrabold">{{ store.counts.pending }}</p><p class="text-[10px] text-teal-100">Pending</p></div>
          <div class="divider hidden sm:block" />
          <div class="text-center"><p class="text-xl font-extrabold">{{ store.counts.inProgress }}</p><p class="text-[10px] text-teal-100">Proses</p></div>
          <div class="divider hidden sm:block" />
          <div class="text-center"><p class="text-xl font-extrabold">{{ store.counts.completed }}</p><p class="text-[10px] text-teal-100">Selesai</p></div>
        </div>
      </div>
    </div>

    <LoadingSpinner v-if="store.loading && !store.orders.length" label="Memuat dashboard…" />
    <div v-else class="grid grid-cols-2 gap-4 lg:grid-cols-4">
      <div v-for="card in cards" :key="card.key" class="stat-card p-4 sm:p-5" :style="`--stat-color:${card.color};--stat-bg:${card.bg};--stat-glow:${card.glow}`">
        <div class="relative z-10 flex items-start justify-between">
          <div>
            <p class="text-[11px] font-bold text-slate-500">{{ card.label }}</p>
            <p class="mt-2 text-2xl font-extrabold tracking-tight text-slate-900">{{ store.counts[card.key] }}</p>
          </div>
          <span class="stat-icon">
            <svg v-if="card.icon === 'folder'" width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M4 6.5h6l2 2h8v10H4z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/></svg>
            <svg v-else-if="card.icon === 'clock'" width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle cx="12" cy="12" r="8.5" stroke="currentColor" stroke-width="1.8"/><path d="M12 7.5V12l3 2" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
            <svg v-else-if="card.icon === 'activity'" width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M4 12h4l2-6 4 12 2-6h4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/></svg>
            <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="m5 12 4 4L19 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </span>
        </div>
      </div>
    </div>

    <section class="panel">
      <div class="panel-header">
        <div>
          <h3 class="panel-title">Order terbaru</h3>
          <p class="mt-1 text-[11px] text-slate-400">Aktivitas order yang paling baru diperbarui</p>
        </div>
        <button type="button" class="ghost-button !min-h-0 !px-2 !py-1.5 text-xs" @click="router.push({ name: 'orders' })">Lihat semua</button>
      </div>
      <div class="p-2 sm:p-3">
        <DataTable
          :columns="[
            { key: 'order_number', label: 'Order', mono: true },
            { key: 'customer_name', label: 'Nasabah' },
            { key: 'notary_name', label: 'Notaris' },
            { key: 'status', label: 'Status' },
          ]"
          :rows="store.orders.slice(0, 5)"
          @row-click="(row) => router.push({ name: 'order-detail', params: { id: row.id } })"
        >
          <template #cell-status="{ value }"><StatusBadge :status="value" /></template>
        </DataTable>
      </div>
    </section>
  </div>
</template>
