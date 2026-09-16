<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../../stores/userStore'

function homeFor(role) {
  return role === 'nasabah' ? 'nasabah-dashboard' : 'dashboard'
}

const route = useRoute()
const router = useRouter()
const store = useUserStore()
const failed = ref('')

onMounted(async () => {
  const token = route.query.token
  if (!token) {
    failed.value = route.query.sso_error || 'Login SSO gagal, silakan coba lagi'
    return
  }
  store.setSession(token, {
    id: Number(route.query.id ?? 0),
    full_name: route.query.full_name ?? '',
    email: route.query.email ?? '',
    role: route.query.role ?? 'nasabah',
    photo_url: route.query.photo_url ?? '',
  })
  await router.push({ name: homeFor(store.role) })
})
</script>

<template>
  <div class="space-y-4 text-center">
    <p v-if="!failed" class="text-sm text-slate-600">Menyelesaikan login SSO…</p>
    <p v-else class="rounded-md bg-red-50 px-3 py-2 text-xs text-red-700">{{ failed }}</p>
    <router-link v-if="failed" :to="{ name: 'login' }" class="text-xs font-semibold text-blue-600 hover:underline">
      Kembali masuk
    </router-link>
  </div>
</template>
