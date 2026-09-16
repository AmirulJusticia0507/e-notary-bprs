<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import Sidebar from './Sidebar.vue'

const route = useRoute()
const { store, logout } = useAuth()

const mobileLinks = [
  { to: { name: 'dashboard' }, label: 'Dashboard', names: ['dashboard'], roles: ['admin', 'legal_officer', 'notary'] },
  { to: { name: 'orders' }, label: 'Order', names: ['orders', 'order-detail', 'order-create'], roles: ['admin', 'legal_officer', 'notary'] },
  { to: { name: 'order-create' }, label: 'Baru', names: ['order-create'], roles: ['admin', 'legal_officer'] },
  { to: { name: 'nasabah-dashboard' }, label: 'Pengajuan', names: ['nasabah-dashboard'], roles: ['nasabah'] },
  { to: { name: 'change-password' }, label: 'Password', names: ['change-password'], roles: ['admin', 'legal_officer', 'notary', 'nasabah'] },
]
const visibleMobileLinks = computed(() => mobileLinks.filter((link) => link.roles.includes(store.role)))
</script>

<template>
  <div class="app-shell bg-[#f5f8fb]">
    <Sidebar />
    <div class="lg:pl-[17rem]">
      <header class="topbar">
        <div class="page-container flex items-center justify-between gap-4 px-5 py-3.5 sm:px-7">
          <div class="min-w-0">
            <p class="truncate text-sm font-bold text-slate-800">Sistem Manajemen Legalitas Akad</p>
            <p class="mt-0.5 hidden text-[11px] text-slate-400 sm:block">Ruang kerja digital untuk order, dokumen, dan pembiayaan</p>
          </div>
          <div class="flex shrink-0 items-center gap-2.5 sm:gap-3.5">
            <div class="hidden items-center gap-2 rounded-xl border border-slate-200 bg-white px-3 py-2 shadow-sm sm:flex">
              <span class="relative flex h-2 w-2">
                <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-60"></span>
                <span class="relative inline-flex h-2 w-2 rounded-full bg-emerald-500"></span>
              </span>
              <span class="max-w-[11rem] truncate text-xs text-slate-500">{{ store.user?.email }}</span>
            </div>
            <span class="hidden rounded-lg bg-teal-50 px-2.5 py-1.5 text-[11px] font-bold capitalize text-teal-700 sm:inline-block">{{ store.user?.role?.replaceAll('_', ' ') }}</span>
            <button type="button" class="secondary-button !min-h-0 !px-3 !py-2 text-xs" @click="logout">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M10 5H5a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h5m4-3h7m0 0-3-3m3 3-3 3" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/></svg>
              Keluar
            </button>
          </div>
        </div>
      </header>
      <main class="page-container px-5 pb-24 pt-6 sm:px-7 sm:pb-10 sm:py-8">
        <router-view />
      </main>
      <nav class="fixed inset-x-0 bottom-0 z-20 flex border-t border-slate-200 bg-white/95 px-2 pb-[env(safe-area-inset-bottom)] pt-1 shadow-[0_-8px_24px_rgba(15,23,42,.06)] backdrop-blur lg:hidden" aria-label="Navigasi mobile">
        <router-link v-for="link in visibleMobileLinks" :key="link.label" :to="link.to" class="flex min-w-0 flex-1 flex-col items-center gap-1 px-1 py-2 text-[10px] font-bold" :class="link.names.includes(route.name) ? 'text-teal-700' : 'text-slate-400'">
          <span class="grid h-7 w-7 place-items-center rounded-xl" :class="link.names.includes(route.name) ? 'bg-teal-50 text-teal-600' : 'bg-slate-100 text-slate-500'">
            <svg v-if="link.label === 'Dashboard'" width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden="true"><rect x="3" y="3" width="7.5" height="7.5" rx="2" stroke="currentColor" stroke-width="1.8"/><rect x="13.5" y="3" width="7.5" height="7.5" rx="2" stroke="currentColor" stroke-width="1.8"/><rect x="3" y="13.5" width="7.5" height="7.5" rx="2" stroke="currentColor" stroke-width="1.8"/><rect x="13.5" y="13.5" width="7.5" height="7.5" rx="2" stroke="currentColor" stroke-width="1.8"/></svg>
            <svg v-else-if="link.label === 'Order'" width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M7 3h7l4 4v14H7z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/><path d="M14 3v5h4M10 12h5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
            <svg v-else-if="link.label === 'Baru'" width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
            <svg v-else-if="link.label === 'Pengajuan'" width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M4 5.5h16v13H4z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/><path d="M8 9h8M8 13h5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
            <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 3 5 6v5c0 4.2 2.8 8 7 10 4.2-2 7-5.8 7-10V6z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/><path d="m9 12 2 2 4-4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </span>
          {{ link.label }}
        </router-link>
      </nav>
    </div>
  </div>
</template>
