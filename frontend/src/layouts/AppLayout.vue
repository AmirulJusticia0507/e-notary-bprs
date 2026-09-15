<script setup>
import { useAuth } from '../composables/useAuth'
import Sidebar from './Sidebar.vue'

const { store, logout } = useAuth()
</script>

<template>
  <div class="min-h-screen bg-slate-50">
    <Sidebar />
    <div class="pl-60">
      <header class="sticky top-0 z-10 flex items-center justify-between border-b border-slate-200 bg-white px-6 py-3">
        <h1 class="text-sm font-semibold text-slate-700">Sistem Manajemen Legalitas Akad</h1>
        <div class="flex items-center gap-3">
          <span class="text-xs text-slate-500">{{ store.user?.email }} · {{ store.user?.role }}</span>
          <button @click="logout" class="rounded-md border border-slate-300 px-3 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-50">
            Keluar
          </button>
        </div>
      </header>
      <main class="p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>
