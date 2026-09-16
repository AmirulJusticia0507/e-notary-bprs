<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import FormInput from '../../components/forms/FormInput.vue'
import PhotoInput from '../../components/forms/PhotoInput.vue'
import AppButton from '../../components/common/AppButton.vue'
import { useFormValidation } from '../../composables/useFormValidation'
import { authService } from '../../services/authService'

const router = useRouter()
const { errors, required, isValidEmail, clear } = useFormValidation()
const form = reactive({ full_name: '', email: '', password: '', confirm: '', photo_url: '' })
const failed = ref('')
const loading = ref(false)

async function submit() {
  failed.value = ''
  clear()
  const okName = required(form.full_name, 'full_name', 'Nama lengkap')
  const okEmail = required(form.email, 'email', 'Email') && isValidEmail(form.email)
  const okPass = required(form.password, 'password', 'Password')
  if (form.password && form.password.length < 6) {
    errors.password = 'Password minimal 6 karakter'
  }
  if (form.confirm !== form.password) {
    errors.confirm = 'Konfirmasi password tidak sama'
  }
  if (!okName || !okEmail || !okPass || errors.password || errors.confirm) return
  loading.value = true
  try {
    await authService.register({ full_name: form.full_name, email: form.email, password: form.password, role: 'nasabah', photo_url: form.photo_url })
    await router.push({ name: 'login', query: { registered: '1' } })
  } catch (e) {
    failed.value = e.response?.data?.error ?? 'Pendaftaran gagal, coba lagi'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form @submit.prevent="submit" class="space-y-4">
    <div class="text-center">
      <h2 class="text-lg font-bold text-slate-900">Daftar Akun Nasabah</h2>
      <p class="text-xs text-slate-500">Buat akun untuk mengajukan & memantau pembiayaan</p>
    </div>
    <FormInput v-model="form.full_name" label="Nama lengkap" placeholder="Nama sesuai KTP" :error="errors.full_name" />
    <FormInput v-model="form.email" label="Email" type="email" placeholder="nama@email.com" :error="errors.email" />
    <FormInput v-model="form.password" label="Password" type="password" placeholder="Minimal 6 karakter" :error="errors.password" />
    <FormInput v-model="form.confirm" label="Konfirmasi password" type="password" placeholder="Ulangi password" :error="errors.confirm" />
    <PhotoInput v-model="form.photo_url" label="Foto profil (opsional)" />
    <p v-if="failed" class="rounded-md bg-red-50 px-3 py-2 text-xs text-red-700">{{ failed }}</p>
    <AppButton type="submit" :loading="loading" class="w-full">Daftar</AppButton>
    <p class="text-center text-xs text-slate-500">
      Sudah punya akun?
      <router-link :to="{ name: 'login' }" class="font-semibold text-blue-600 hover:underline">Masuk</router-link>
    </p>
  </form>
</template>
