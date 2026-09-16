<script setup>
import { reactive, ref } from 'vue'
import FormInput from '../../components/forms/FormInput.vue'
import AppButton from '../../components/common/AppButton.vue'
import { useFormValidation } from '../../composables/useFormValidation'
import { authService } from '../../services/authService'

const { errors, required, isValidEmail, clear } = useFormValidation()
const form = reactive({ email: '' })
const loading = ref(false)
const done = ref(false)
const failed = ref('')

async function submit() {
  failed.value = ''
  done.value = false
  clear()
  if (!required(form.email, 'email', 'Email') || !isValidEmail(form.email)) return
  loading.value = true
  try {
    await authService.forgotPassword(form.email)
    done.value = true
  } catch (e) {
    failed.value = e.response?.data?.error ?? 'Gagal memproses permintaan'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form @submit.prevent="submit" class="space-y-4">
    <div class="text-center">
      <h2 class="text-lg font-bold text-slate-900">Lupa Password</h2>
      <p class="text-xs text-slate-500">Masukkan email akun Anda</p>
    </div>
    <FormInput v-model="form.email" label="Email" type="email" placeholder="nama@email.com" :error="errors.email" />
    <p v-if="done" class="rounded-md bg-emerald-50 px-3 py-2 text-xs text-emerald-700">
      Jika email terdaftar, kode reset telah dibuat (berlaku 1 jam). Hubungi admin/CS BPRS untuk mendapatkan kode tersebut, lalu buka halaman reset password.
    </p>
    <p v-if="failed" class="rounded-md bg-red-50 px-3 py-2 text-xs text-red-700">{{ failed }}</p>
    <AppButton type="submit" :loading="loading" class="w-full">Minta Kode Reset</AppButton>
    <p class="text-center text-xs text-slate-500">
      <router-link :to="{ name: 'login' }" class="font-semibold text-blue-600 hover:underline">Kembali masuk</router-link>
      <span class="mx-1">·</span>
      <router-link :to="{ name: 'reset-password' }" class="font-semibold text-blue-600 hover:underline">Saya sudah punya kode</router-link>
    </p>
  </form>
</template>
