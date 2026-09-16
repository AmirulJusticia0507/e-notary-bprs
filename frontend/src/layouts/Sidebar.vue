<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '../stores/userStore'

const route = useRoute()
const userStore = useUserStore()

const links = [
  { to: { name: 'dashboard' }, label: 'Dashboard', match: ['dashboard'], roles: ['admin', 'legal_officer', 'notary'], icon: 'dashboard' },
  { to: { name: 'orders' }, label: 'Legal Orders', match: ['orders', 'order-detail', 'order-create'], roles: ['admin', 'legal_officer', 'notary'], icon: 'orders' },
  { to: { name: 'order-create' }, label: 'Order Baru', match: ['order-create'], roles: ['admin', 'legal_officer'], icon: 'plus' },
  { to: { name: 'nasabah-dashboard' }, label: 'Pengajuan Saya', match: ['nasabah-dashboard'], roles: ['nasabah'], icon: 'applications' },
  { to: { name: 'change-password' }, label: 'Ganti Password', match: ['change-password'], roles: ['admin', 'legal_officer', 'notary', 'nasabah'], icon: 'security' },
]

const visibleLinks = computed(() => links.filter((link) => link.roles.includes(userStore.role)))
const isActive = computed(() => (match) => match.includes(route.name))

const roleLabels = {
  admin: 'Administrator',
  legal_officer: 'Legal Officer',
  notary: 'Notaris',
  nasabah: 'Nasabah',
}
</script>

<template>
  <aside class="sidebar hidden lg:flex">
    <div class="sidebar-content">
      <div class="flex items-center gap-3 border-b border-white/10 px-5 py-5">
        <div class="grid h-10 w-10 shrink-0 place-items-center rounded-2xl bg-gradient-to-br from-teal-300 to-sky-500 font-mono text-lg font-bold text-[#06202b] shadow-lg shadow-teal-950/30">N</div>
        <div class="min-w-0">
          <p class="truncate text-sm font-bold text-white">e-Notary BPRS</p>
          <p class="text-[11px] text-slate-400">Legal Akad Syariah</p>
        </div>
      </div>

      <nav class="flex-1 space-y-1.5 overflow-y-auto px-3.5 py-5" aria-label="Navigasi utama">
        <p class="mb-2 px-2.5 text-[10px] font-bold uppercase tracking-[0.16em] text-slate-500">Workspace</p>
        <router-link v-for="link in visibleLinks" :key="link.label" :to="link.to" class="sidebar-link" :class="{ 'active-sidebar-link': isActive(link.match) }">
          <svg v-if="link.icon === 'dashboard'" width="17" height="17" viewBox="0 0 24 24" fill="none" aria-hidden="true"><rect x="3" y="3" width="7.5" height="7.5" rx="2" stroke="currentColor" stroke-width="1.8"/><rect x="13.5" y="3" width="7.5" height="7.5" rx="2" stroke="currentColor" stroke-width="1.8"/><rect x="3" y="13.5" width="7.5" height="7.5" rx="2" stroke="currentColor" stroke-width="1.8"/><rect x="13.5" y="13.5" width="7.5" height="7.5" rx="2" stroke="currentColor" stroke-width="1.8"/></svg>
          <svg v-else-if="link.icon === 'orders'" width="17" height="17" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M7 3h7l4 4v14H7z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/><path d="M14 3v5h4M10 12h5m-5 4h5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
          <svg v-else-if="link.icon === 'plus'" width="17" height="17" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
          <svg v-else-if="link.icon === 'applications'" width="17" height="17" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M4 5.5h16v13H4z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/><path d="M8 9h8M8 13h5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
          <svg v-else width="17" height="17" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 3 5 6v5c0 4.2 2.8 8 7 10 4.2-2 7-5.8 7-10V6z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/><path d="m9 12 2 2 4-4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/></svg>
          <span>{{ link.label }}</span>
        </router-link>
      </nav>

      <div class="border-t border-white/10 p-4">
        <div class="rounded-2xl border border-teal-200/20 bg-gradient-to-br from-teal-400/10 to-sky-500/10 p-3.5">
          <div class="flex items-center gap-2.5">
            <div class="grid h-8 w-8 shrink-0 place-items-center rounded-full bg-teal-300/20 text-xs font-bold text-teal-200">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle cx="12" cy="12" r="8.5" stroke="currentColor" stroke-width="1.8"/><path d="M12 7.5V12l3 2" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
            </div>
            <div class="min-w-0">
              <p class="truncate text-[11px] font-bold text-slate-200">SLA closing</p>
              <p class="text-[10px] text-slate-400">Target 3–5 hari</p>
            </div>
          </div>
        </div>
        <div class="mt-3.5 flex items-center gap-2.5 px-1">
          <div v-if="userStore.user?.photo_url" class="h-9 w-9 overflow-hidden rounded-full border-2 border-teal-300/50">
            <img :src="userStore.user.photo_url" alt="Foto profil" class="h-full w-full object-cover" />
          </div>
          <div v-else class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-gradient-to-br from-teal-200 to-sky-300 text-xs font-bold text-teal-950">{{ (userStore.user?.full_name || 'U').charAt(0).toUpperCase() }}</div>
          <div class="min-w-0 flex-1">
            <p class="truncate text-[11px] font-bold text-slate-200">{{ userStore.user?.full_name || 'Pengguna' }}</p>
            <p class="text-[10px] text-slate-400">{{ roleLabels[userStore.role] || userStore.role || 'Role belum diatur' }}</p>
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>
