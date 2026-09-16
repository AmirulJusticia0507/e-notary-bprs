<script setup>
import { reactive, ref } from 'vue'
import FormInput from '../../components/forms/FormInput.vue'
import AppButton from '../../components/common/AppButton.vue'
import { useFormValidation } from '../../composables/useFormValidation'
import { authService } from '../../services/authService'

const { errors, required, clear } = useFormValidation()
const form = reactive({ current_password: '', new_password: '', confirm: '' })
const loading = ref(false)
const failed = ref('')
const done = ref('')
async function submit() { failed.value = ''; done.value = ''; clear(); const okCur = required(form.current_password, 'current_password', 'Password lama'); const okNew = required(form.new_password, 'new_password', 'Password baru'); if (form.new_password && form.new_password.length < 6) errors.new_password = 'Password minimal 6 karakter'; if (form.confirm !== form.new_password) errors.confirm = 'Konfirmasi password tidak sama'; if (!okCur || !okNew || errors.new_password || errors.confirm) return; loading.value = true; try { await authService.changePassword(form.current_password, form.new_password); done.value = 'Password berhasil diganti'; Object.assign(form, { current_password: '', new_password: '', confirm: '' }) } catch (error) { failed.value = error.response?.data?.error ?? 'Ganti password gagal' } finally { loading.value = false } }
</script>

<template>
  <div class="mx-auto max-w-xl space-y-6">
    <div class="page-header !mb-0"><div><p class="text-[11px] font-bold uppercase tracking-[0.16em] text-teal-600">Account security</p><h1 class="page-title">Ganti Password</h1><p class="page-subtitle">Jaga keamanan akun dengan password yang kuat dan unik.</p></div></div>
    <form class="panel p-5 sm:p-6" @submit.prevent="submit"><div class="space-y-4"><FormInput v-model="form.current_password" label="Password lama" type="password" :error="errors.current_password" /><FormInput v-model="form.new_password" label="Password baru" type="password" placeholder="Minimal 6 karakter" :error="errors.new_password" /><FormInput v-model="form.confirm" label="Konfirmasi password baru" type="password" :error="errors.confirm" /></div><div v-if="done" class="notice notice-success mt-5">{{ done }}</div><div v-if="failed" class="notice notice-error mt-5">{{ failed }}</div><AppButton type="submit" :loading="loading" class="mt-5">Simpan Password</AppButton></form>
  </div>
</template>
