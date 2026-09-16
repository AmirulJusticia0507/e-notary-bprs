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
async function submit() { failed.value = ''; done.value = false; clear(); if (!required(form.email, 'email', 'Email') || !isValidEmail(form.email)) return; loading.value = true; try { await authService.forgotPassword(form.email); done.value = true } catch (error) { failed.value = error.response?.data?.error ?? 'Gagal memproses permintaan' } finally { loading.value = false } }
</script>

<template>
  <form @submit.prevent="submit" class="space-y-5">
    <div><p class="text-[11px] font-bold uppercase tracking-[0.16em] text-teal-600">Account recovery</p><h2 class="mt-2 auth-heading">Lupa Password</h2><p class="auth-subheading">Masukkan email akun Anda untuk membuat kode reset.</p></div>
    <FormInput v-model="form.email" label="Email" type="email" placeholder="nama@email.com" :error="errors.email" />
    <div v-if="done" class="notice notice-info">Jika email terdaftar, kode reset telah dibuat (berlaku 1 jam). Hubungi admin/CS BPRS untuk mendapatkan kode tersebut, lalu buka halaman reset password.</div>
    <div v-if="failed" class="notice notice-error">{{ failed }}</div>
    <AppButton type="submit" :loading="loading" class="w-full">Minta Kode Reset</AppButton>
    <p class="text-center text-xs text-slate-500"><router-link :to="{ name: 'login' }" class="font-bold text-teal-700 hover:text-teal-800">Kembali masuk</router-link><span class="mx-2">·</span><router-link :to="{ name: 'reset-password' }" class="font-bold text-teal-700 hover:text-teal-800">Saya sudah punya kode</router-link></p>
  </form>
</template>
