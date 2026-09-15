<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const links = [
  { to: { name: 'dashboard' }, label: 'Dashboard', match: ['dashboard'] },
  { to: { name: 'orders' }, label: 'Legal Orders', match: ['orders', 'order-detail', 'order-create'] },
]

const isActive = computed(() => (match) => match.includes(route.name))
</script>

<template>
  <aside class="fixed inset-y-0 left-0 flex w-60 flex-col bg-slate-900">
    <div class="flex items-center gap-2 border-b border-slate-800 px-5 py-4">
      <div class="flex h-8 w-8 items-center justify-center rounded-md bg-blue-600 font-mono text-sm font-bold text-white">N</div>
      <div>
        <p class="text-sm font-bold text-white">e-Notary BPRS</p>
        <p class="text-[11px] text-slate-400">Legal Akad Syariah</p>
      </div>
    </div>
    <nav class="flex-1 space-y-1 p-3">
      <router-link
        v-for="link in links"
        :key="link.label"
        :to="link.to"
        class="block rounded-md px-3 py-2 text-sm font-medium transition-colors"
        :class="isActive(link.match) ? 'bg-blue-600 text-white' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
      >
        {{ link.label }}
      </router-link>
    </nav>
    <div class="border-t border-slate-800 p-4 text-[11px] text-slate-500">SLA closing target 3–5 hari</div>
  </aside>
</template>
