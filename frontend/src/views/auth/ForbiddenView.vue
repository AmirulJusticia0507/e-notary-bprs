<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div class="mx-auto flex min-h-[65vh] max-w-md items-center justify-center text-center">
    <div class="panel w-full p-7">
      <div class="mx-auto grid h-14 w-14 place-items-center rounded-2xl bg-rose-50 text-2xl font-extrabold text-rose-500">!</div>
      <p class="mt-5 font-mono text-4xl font-extrabold tracking-tight text-slate-900">403</p>
      <h2 class="mt-2 text-lg font-extrabold text-slate-900">Akses ditolak</h2>
      <p class="mx-auto mt-2 max-w-sm text-xs leading-5 text-slate-500">Role akun Anda tidak memiliki izin untuk membuka halaman ini. Hubungi admin bila Anda merasa ini keliru.</p>
      <AppButton class="mt-6 w-full" @click="router.push({ name: 'dashboard' })">Kembali ke Dashboard</AppButton>
    </div>
  </div>
</template>
