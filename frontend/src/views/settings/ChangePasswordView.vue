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

async function submit() {
  failed.value = ''
  done.value = ''
  clear()
  const okCur = required(form.current_password, 'current_password', 'Password lama')
  const okNew = required(form.new_password, 'new_password', 'Password baru')
  if (form.new_password && form.new_password.length < 6) {
    errors.new_password = 'Password minimal 6 karakter'
  }
  if (form.confirm !== form.new_password) {
    errors.confirm = 'Konfirmasi password tidak sama'
  }
  if (!okCur || !okNew || errors.new_password || errors.confirm) return
  loading.value = true
  try {
    await authService.changePassword(form.current_password, form.new_password)
    done.value = 'Password berhasil diganti'
    Object.assign(form, { current_password: '', new_password: '', confirm: '' })
  } catch (e) {
    failed.value = e.response?.data?.error ?? 'Ganti password gagal'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-md space-y-5">
    <div>
      <h2 class="text-lg font-bold text-slate-900">Ganti Password</h2>
      <p class="text-xs text-slate-500">Berlaku untuk semua role: nasabah, legal, notaris, admin</p>
    </div>
    <form @submit.prevent="submit" class="space-y-4 rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
      <FormInput v-model="form.current_password" label="Password lama" type="password" :error="errors.current_password" />
      <FormInput v-model="form.new_password" label="Password baru" type="password" placeholder="Minimal 6 karakter" :error="errors.new_password" />
      <FormInput v-model="form.confirm" label="Konfirmasi password baru" type="password" :error="errors.confirm" />
      <p v-if="done" class="rounded-md bg-emerald-50 px-3 py-2 text-xs text-emerald-700">{{ done }}</p>
      <p v-if="failed" class="rounded-md bg-red-50 px-3 py-2 text-xs text-red-700">{{ failed }}</p>
      <AppButton type="submit" :loading="loading">Simpan</AppButton>
    </form>
  </div>
</template>
