<script setup>
// Layout halaman auth: kartu terpusat di atas canvas netral.
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-slate-900 p-6">
    <div class="w-full max-w-md rounded-lg bg-white p-6 shadow-sm">
      <div class="mb-5 text-center">
        <p class="text-lg font-bold text-slate-900">e-Notary BPRS</p>
        <p class="text-xs text-slate-500">Masuk untuk mengelola legalitas akad</p>
      </div>
      <router-view />
    </div>
  </div>
</template>
