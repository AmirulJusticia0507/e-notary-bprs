<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import FormInput from '../../components/forms/FormInput.vue'
import AppButton from '../../components/common/AppButton.vue'
import { useFormValidation } from '../../composables/useFormValidation'
import { authService } from '../../services/authService'

const route = useRoute()
const router = useRouter()
const { errors, required, clear } = useFormValidation()
const form = reactive({ token: route.query.token ?? '', new_password: '', confirm: '' })
const loading = ref(false)
const failed = ref('')

async function submit() {
  failed.value = ''
  clear()
  const okToken = required(form.token, 'token', 'Kode reset')
  const okPass = required(form.new_password, 'new_password', 'Password baru')
  if (form.new_password && form.new_password.length < 6) {
    errors.new_password = 'Password minimal 6 karakter'
  }
  if (form.confirm !== form.new_password) {
    errors.confirm = 'Konfirmasi password tidak sama'
  }
  if (!okToken || !okPass || errors.new_password || errors.confirm) return
  loading.value = true
  try {
    await authService.resetPassword(form.token.trim(), form.new_password)
    await router.push({ name: 'login', query: { reset: '1' } })
  } catch (e) {
    failed.value = e.response?.data?.error ?? 'Reset gagal, kode mungkin kedaluwarsa'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form @submit.prevent="submit" class="space-y-4">
    <div class="text-center">
      <h2 class="text-lg font-bold text-slate-900">Reset Password</h2>
      <p class="text-xs text-slate-500">Masukkan kode dari admin/CS + password baru</p>
    </div>
    <FormInput v-model="form.token" label="Kode reset" placeholder="Tempel kode dari admin" :error="errors.token" />
    <FormInput v-model="form.new_password" label="Password baru" type="password" placeholder="Minimal 6 karakter" :error="errors.new_password" />
    <FormInput v-model="form.confirm" label="Konfirmasi password baru" type="password" placeholder="Ulangi password" :error="errors.confirm" />
    <p v-if="failed" class="rounded-md bg-red-50 px-3 py-2 text-xs text-red-700">{{ failed }}</p>
    <AppButton type="submit" :loading="loading" class="w-full">Simpan Password Baru</AppButton>
    <p class="text-center text-xs text-slate-500">
      <router-link :to="{ name: 'login' }" class="font-semibold text-blue-600 hover:underline">Kembali masuk</router-link>
    </p>
  </form>
</template>
