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
  if (form.new_password && form.new_password.length < 6) errors.new_password = 'Password minimal 6 karakter'
  if (form.confirm !== form.new_password) errors.confirm = 'Konfirmasi password tidak sama'
  if (!okToken || !okPass || errors.new_password || errors.confirm) return
  loading.value = true
  try {
    await authService.resetPassword(form.token.trim(), form.new_password)
    await router.push({ name: 'login', query: { reset: '1' } })
  } catch (error) {
    failed.value = error.response?.data?.error ?? 'Reset gagal, kode mungkin kedaluwarsa'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form @submit.prevent="submit" class="space-y-5">
    <div><p class="text-[11px] font-bold uppercase tracking-[0.16em] text-teal-600">Account recovery</p><h2 class="mt-2 auth-heading">Reset Password</h2><p class="auth-subheading">Masukkan kode dari admin/CS dan password baru.</p></div>
    <div class="space-y-4"><FormInput v-model="form.token" label="Kode reset" placeholder="Tempel kode dari admin" :error="errors.token" /><FormInput v-model="form.new_password" label="Password baru" type="password" placeholder="Minimal 6 karakter" :error="errors.new_password" /><FormInput v-model="form.confirm" label="Konfirmasi password baru" type="password" placeholder="Ulangi password" :error="errors.confirm" /></div>
    <div v-if="failed" class="notice notice-error">{{ failed }}</div>
    <AppButton type="submit" :loading="loading" class="w-full">Simpan Password Baru</AppButton>
    <p class="text-center text-xs text-slate-500"><router-link :to="{ name: 'login' }" class="font-bold text-teal-700 hover:text-teal-800">Kembali masuk</router-link></p>
  </form>
</template>
