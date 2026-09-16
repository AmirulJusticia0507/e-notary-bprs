<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import FormInput from '../../components/forms/FormInput.vue'
import AppButton from '../../components/common/AppButton.vue'
import CaptchaBox from '../../components/forms/CaptchaBox.vue'
import { useFormValidation } from '../../composables/useFormValidation'
import { useUserStore } from '../../stores/userStore'

const router = useRouter()
const route = useRoute()
const store = useUserStore()
const justRegistered = route.query.registered === '1'
const { errors, required, isValidEmail } = useFormValidation()
const form = reactive({ email: '', password: '' })
const failed = ref('')
const isDev = import.meta.env.DEV
const captchaOk = ref(false)
const captchaRef = ref(null)

async function submit() {
  failed.value = ''
  const okEmail = required(form.email, 'email', 'Email') && isValidEmail(form.email)
  const okPass = required(form.password, 'password', 'Password')
  if (!okEmail || !okPass) return
  if (!captchaOk.value) {
    failed.value = 'Isi captcha dulu sebelum masuk'
    return
  }
  try {
    const user = await store.login(form.email, form.password)
    await router.push({ name: user?.role === 'nasabah' ? 'nasabah-dashboard' : 'dashboard' })
  } catch {
    failed.value = store.error || 'Email atau password salah'
    captchaRef.value?.refresh()
  }
}
</script>

<template>
  <form @submit.prevent="submit" class="space-y-4">
    <FormInput v-model="form.email" label="Email" type="email" placeholder="nama@bprs.local" :error="errors.email" />
    <FormInput v-model="form.password" label="Password" type="password" placeholder="••••••••" :error="errors.password" />
    <p v-if="justRegistered" class="rounded-md bg-emerald-50 px-3 py-2 text-xs text-emerald-700">Pendaftaran berhasil, silakan masuk.</p>
    <CaptchaBox ref="captchaRef" v-model="captchaOk" />
    <p v-if="failed" class="rounded-md bg-red-50 px-3 py-2 text-xs text-red-700">{{ failed }}</p>
    <AppButton type="submit" :loading="store.loading" :disabled="!captchaOk" class="w-full">Masuk</AppButton>
    <p v-if="isDev" class="text-center font-mono text-[11px] text-slate-400">dev: admin@bprs.local / password123</p>
    <p class="text-center text-xs text-slate-500">
      Belum punya akun?
      <router-link :to="{ name: 'register' }" class="font-semibold text-blue-600 hover:underline">Daftar sebagai nasabah</router-link>
    </p>
  </form>
</template>
