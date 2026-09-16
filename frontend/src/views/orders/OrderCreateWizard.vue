<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppButton from '../../components/common/AppButton.vue'
import FormInput from '../../components/forms/FormInput.vue'
import { useUserStore } from '../../stores/userStore'
import { useOrderStore } from '../../stores/orderStore'

const router = useRouter()
const userStore = useUserStore()
const store = useOrderStore()
const step = ref(1)
const saving = ref(false)
const failed = ref('')
const form = reactive({ order_number: '', financing_id: '', notary_id: '' })

onMounted(() => store.fetchReferences().catch((error) => (failed.value = error.response?.data?.error ?? 'Gagal memuat referensi')))

function next() {
  failed.value = ''
  if (step.value === 1 && !form.financing_id) failed.value = 'Pilih pembiayaan dulu'
  else if (step.value === 2 && !form.notary_id) failed.value = 'Pilih notaris dulu'
  else step.value += 1
}

async function submit() {
  saving.value = true
  failed.value = ''
  try {
    await store.createOrder({ order_number: form.order_number || `ORD-${Date.now()}`, financing_id: Number(form.financing_id), notary_id: Number(form.notary_id), assigned_to: userStore.user?.id })
    await router.push({ name: 'orders' })
  } catch (error) {
    failed.value = error.response?.data?.error ?? 'Gagal membuat order'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-3xl space-y-6">
    <div class="page-header !mb-0">
      <div>
        <p class="text-[11px] font-bold uppercase tracking-[0.16em] text-teal-600">Legal workflow</p>
        <h1 class="page-title">Order Baru</h1>
        <p class="page-subtitle">Buat order dari pembiayaan yang telah terverifikasi.</p>
      </div>
    </div>

    <div class="flex items-center gap-2 sm:gap-3" aria-label="Tahapan pembuatan order">
      <div v-for="item in [{ n: 1, label: 'Pembiayaan' }, { n: 2, label: 'Notaris' }, { n: 3, label: 'Review' }]" :key="item.n" class="flex flex-1 items-center gap-2">
        <div class="step flex-1" :class="{ 'step-active': step >= item.n }"><span class="mr-2 font-mono">{{ item.n }}</span>{{ item.label }}</div>
        <span v-if="item.n < 3" class="h-px flex-1 bg-slate-200"></span>
      </div>
    </div>

    <section class="panel p-5 sm:p-6">
      <div v-if="step === 1" class="space-y-5">
        <div>
          <h2 class="text-base font-extrabold text-slate-900">Pilih pembiayaan</h2>
          <p class="mt-1 text-xs text-slate-500">Hubungkan order dengan aplikasi pembiayaan dari CBS/LOS.</p>
        </div>
        <FormInput v-model="form.order_number" label="Nomor Order (opsional)" mono placeholder="ORD-2026-0001" />
        <div class="space-y-1.5">
          <label class="field-label">Pembiayaan (sync CBS/LOS)</label>
          <select v-model="form.financing_id" class="field-input">
            <option value="">— pilih pembiayaan —</option>
            <option v-for="financing in store.financings" :key="financing.id" :value="financing.id">{{ financing.customer_name }} · {{ financing.customer_nik }}</option>
          </select>
        </div>
      </div>

      <div v-else-if="step === 2" class="space-y-5">
        <div>
          <h2 class="text-base font-extrabold text-slate-900">Pilih notaris rekanan</h2>
          <p class="mt-1 text-xs text-slate-500">Tentukan notaris yang akan memproses dokumen akad.</p>
        </div>
        <div class="space-y-1.5">
          <label class="field-label">Notaris Rekanan</label>
          <select v-model="form.notary_id" class="field-input">
            <option value="">— pilih notaris —</option>
            <option v-for="notary in store.notaries" :key="notary.id" :value="notary.id">{{ notary.full_name }} · {{ notary.wilayah_kerja }}</option>
          </select>
        </div>
      </div>

      <div v-else class="space-y-5">
        <div>
          <h2 class="text-base font-extrabold text-slate-900">Review order</h2>
          <p class="mt-1 text-xs text-slate-500">Periksa kembali informasi sebelum order dibuat.</p>
        </div>
        <div class="rounded-xl border border-slate-200 bg-slate-50 p-4 text-sm">
          <dl class="space-y-3">
            <div class="flex justify-between gap-3"><dt class="text-slate-500">Order</dt><dd class="font-mono font-bold text-slate-900">{{ form.order_number || '(otomatis)' }}</dd></div>
            <div class="flex justify-between gap-3"><dt class="text-slate-500">SLA</dt><dd class="font-medium text-slate-900">14 hari sejak dibuat</dd></div>
            <div class="flex justify-between gap-3"><dt class="text-slate-500">PIC</dt><dd class="font-medium text-slate-900">{{ userStore.user?.email }}</dd></div>
          </dl>
        </div>
      </div>

      <div v-if="failed" class="notice notice-error mt-5">{{ failed }}</div>
      <div class="mt-6 flex justify-between border-t border-slate-100 pt-5">
        <AppButton variant="secondary" :disabled="step === 1" @click="step -= 1">Kembali</AppButton>
        <AppButton v-if="step < 3" @click="next">Lanjut</AppButton>
        <AppButton v-else :loading="saving" @click="submit">Buat Order</AppButton>
      </div>
    </section>
  </div>
</template>
