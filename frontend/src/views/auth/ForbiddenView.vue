<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div class="flex min-h-[60vh] flex-col items-center justify-center space-y-4 text-center">
    <p class="text-5xl font-bold text-slate-300">403</p>
    <h2 class="text-lg font-bold text-slate-900">Akses ditolak</h2>
    <p class="max-w-sm text-sm text-slate-500">
      Role akun Anda tidak memiliki izin untuk membuka halaman ini. Hubungi admin bila Anda merasa ini keliru.
    </p>
    <button
      @click="router.push({ name: 'dashboard' })"
      class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
    >
      Kembali ke Dashboard
    </button>
  </div>
</template>
