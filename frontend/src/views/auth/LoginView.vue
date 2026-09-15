<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import FormInput from '../../components/forms/FormInput.vue'
import AppButton from '../../components/common/AppButton.vue'
import { useFormValidation } from '../../composables/useFormValidation'
import { useUserStore } from '../../stores/userStore'

const router = useRouter()
const store = useUserStore()
const { errors, required, isValidEmail } = useFormValidation()
const form = reactive({ email: '', password: '' })
const failed = ref('')

async function submit() {
  failed.value = ''
  const okEmail = required(form.email, 'email', 'Email') && isValidEmail(form.email)
  const okPass = required(form.password, 'password', 'Password')
  if (!okEmail || !okPass) return
  try {
    await store.login(form.email, form.password)
    await router.push({ name: 'dashboard' })
  } catch {
    failed.value = store.error || 'Email atau password salah'
  }
}
</script>

<template>
  <form @submit.prevent="submit" class="space-y-4">
    <FormInput v-model="form.email" label="Email" type="email" placeholder="nama@bprs.local" :error="errors.email" />
    <FormInput v-model="form.password" label="Password" type="password" placeholder="••••••••" :error="errors.password" />
    <p v-if="failed" class="rounded-md bg-red-50 px-3 py-2 text-xs text-red-700">{{ failed }}</p>
    <AppButton type="submit" :loading="store.loading" class="w-full">Masuk</AppButton>
    <p class="text-center font-mono text-[11px] text-slate-400">dev: admin@bprs.local / admin123</p>
  </form>
</template>
