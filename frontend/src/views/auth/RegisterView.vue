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
async function submit() { failed.value = ''; clear(); const okName = required(form.full_name, 'full_name', 'Nama lengkap'); const okEmail = required(form.email, 'email', 'Email') && isValidEmail(form.email); const okPass = required(form.password, 'password', 'Password'); if (form.password && form.password.length < 6) errors.password = 'Password minimal 6 karakter'; if (form.confirm !== form.password) errors.confirm = 'Konfirmasi password tidak sama'; if (!okName || !okEmail || !okPass || errors.password || errors.confirm) return; loading.value = true; try { await authService.register({ full_name: form.full_name, email: form.email, password: form.password, role: 'nasabah', photo_url: form.photo_url }); await router.push({ name: 'login', query: { registered: '1' } }) } catch (error) { failed.value = error.response?.data?.error ?? 'Pendaftaran gagal, coba lagi' } finally { loading.value = false } }
</script>

<template>
  <form @submit.prevent="submit" class="space-y-5">
    <div><p class="text-[11px] font-bold uppercase tracking-[0.16em] text-teal-600">Nasabah onboarding</p><h2 class="mt-2 auth-heading">Daftar Akun Nasabah</h2><p class="auth-subheading">Buat akun untuk mengajukan dan memantau pembiayaan.</p></div>
    <div class="grid gap-4 sm:grid-cols-2"><FormInput v-model="form.full_name" label="Nama lengkap" placeholder="Nama sesuai KTP" :error="errors.full_name" /><FormInput v-model="form.email" label="Email" type="email" placeholder="nama@email.com" :error="errors.email" /><FormInput v-model="form.password" label="Password" type="password" placeholder="Minimal 6 karakter" :error="errors.password" /><FormInput v-model="form.confirm" label="Konfirmasi password" type="password" placeholder="Ulangi password" :error="errors.confirm" /></div>
    <PhotoInput v-model="form.photo_url" label="Foto profil (opsional)" />
    <div v-if="failed" class="notice notice-error">{{ failed }}</div>
    <AppButton type="submit" :loading="loading" class="w-full">Daftar Akun</AppButton>
    <p class="text-center text-xs text-slate-500">Sudah punya akun? <router-link :to="{ name: 'login' }" class="font-bold text-teal-700 hover:text-teal-800">Masuk</router-link></p>
  </form>
</template>
