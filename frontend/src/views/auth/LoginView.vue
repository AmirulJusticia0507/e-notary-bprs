<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import FormInput from '../../components/forms/FormInput.vue'
import AppButton from '../../components/common/AppButton.vue'
import CaptchaBox from '../../components/forms/CaptchaBox.vue'
import { useFormValidation } from '../../composables/useFormValidation'
import { useUserStore } from '../../stores/userStore'
import { authService } from '../../services/authService'

const router = useRouter()
const route = useRoute()
const store = useUserStore()
const justRegistered = route.query.registered === '1'
const justReset = route.query.reset === '1'
const ssoError = route.query.sso_error ?? ''
const ssoUrl = authService.ssoLoginUrl()
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
  <form @submit.prevent="submit" class="space-y-5">
    <div class="rounded-2xl border border-teal-100 bg-teal-50/70 p-4">
      <div class="flex items-start gap-3">
        <div class="mt-0.5 grid h-8 w-8 shrink-0 place-items-center rounded-full bg-teal-100 text-teal-700">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 3 5 6v5c0 4.2 2.8 8 7 10 4.2-2 7-5.8 7-10V6z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/><path d="m9 12 2 2 4-4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </div>
        <div>
          <p class="text-xs font-bold text-teal-900">Akses aman untuk tim BPRS</p>
          <p class="mt-1 text-[11px] leading-5 text-teal-800/80">Gunakan akun yang telah terdaftar untuk mengelola order dan legalitas akad.</p>
        </div>
      </div>
    </div>

    <div class="grid gap-4 sm:grid-cols-2">
      <FormInput v-model="form.email" label="Email" type="email" placeholder="nama@bprs.local" :error="errors.email" />
      <FormInput v-model="form.password" label="Password" type="password" placeholder="••••••••" :error="errors.password" />
    </div>

    <div v-if="justRegistered || justReset" class="notice notice-success">
      {{ justRegistered ? 'Pendaftaran berhasil, silakan masuk.' : 'Password berhasil direset, silakan masuk.' }}
    </div>
    <div v-if="ssoError" class="notice notice-error">{{ ssoError }}</div>
    <CaptchaBox ref="captchaRef" v-model="captchaOk" />
    <div v-if="failed" class="notice notice-error">{{ failed }}</div>
    <AppButton type="submit" :loading="store.loading" :disabled="!captchaOk" class="w-full !min-h-11">Masuk ke Dashboard</AppButton>

    <div v-if="isDev" class="rounded-xl border border-slate-200 bg-slate-50 px-3 py-2.5 text-center">
      <p class="text-[10px] font-bold uppercase tracking-[0.14em] text-slate-400">Development access</p>
      <p class="mt-1 font-mono text-[11px] text-slate-600">admin@bprs.local / password123</p>
    </div>

    <div class="flex items-center gap-3 text-[11px] font-bold uppercase tracking-[0.14em] text-slate-400">
      <span class="h-px flex-1 bg-slate-200"></span> atau <span class="h-px flex-1 bg-slate-200"></span>
    </div>
    <a
      :href="ssoUrl"
      class="block rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-center text-sm font-bold text-slate-700 hover:bg-slate-50"
    >
      Login dengan SSO Keycloak
    </a>
    <div class="flex items-center justify-center gap-2 text-xs text-slate-500">
      <router-link :to="{ name: 'forgot-password' }" class="font-bold text-teal-700 hover:text-teal-800">Lupa password?</router-link>
      <span class="text-slate-300">·</span>
      <router-link :to="{ name: 'register' }" class="font-bold text-teal-700 hover:text-teal-800">Daftar sebagai nasabah</router-link>
    </div>
  </form>
</template>
