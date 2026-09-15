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

onMounted(() => store.fetchReferences().catch((e) => (failed.value = e.response?.data?.error ?? 'Gagal memuat referensi')))

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
    await store.createOrder({
      order_number: form.order_number || `ORD-${Date.now()}`,
      financing_id: Number(form.financing_id),
      notary_id: Number(form.notary_id),
      assigned_to: userStore.user?.id,
    })
    await router.push({ name: 'orders' })
  } catch (e) {
    failed.value = e.response?.data?.error ?? 'Gagal membuat order'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-2xl space-y-5">
    <h2 class="text-lg font-bold text-slate-900">Order Baru</h2>
    <div class="flex items-center justify-between mb-2 text-xs text-slate-500">
      <span v-if="step >= 1" class="font-medium text-slate-500">Langkah 1 dari 3</span>
      <span v-else class="font-medium text-slate-400">Langkah 1 dari 3</span>
      <span v-if="step >= 2" class="font-medium text-slate-500 separator">/</span>
      <span v-else class="font-medium text-slate-400 separator">/</span>
      <span v-if="step >= 3" class="font-medium text-slate-500">Langkah 3 dari 3</span>
      <span v-else class="font-medium text-slate-400">Langkah 3 dari 3</span>
    </div>
    <ol class="flex gap-2 text-xs font-semibold">
      <li v-for="n in 3" :key="n" class="flex-1 rounded-md px-3 py-2" :class="step >= n ? 'bg-blue-600 text-white' : 'bg-slate-200 text-slate-500'">
        Langkah {{ n }}
      </li>
    </ol>
    <div class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
      <div v-if="step === 1" class="space-y-4">
        <FormInput v-model="form.order_number" label="Nomor Order (opsional)" mono placeholder="ORD-2026-0001" />
        <div class="space-y-1">
          <label class="block text-xs font-semibold text-slate-700">Pembiayaan (sync CBS/LOS)</label>
          <select v-model="form.financing_id" class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm">
            <option value="">— pilih —</option>
            <option v-for="f in store.financings" :key="f.id" :value="f.id">{{ f.customer_name }} · {{ f.customer_nik }}</option>
          </select>
        </div>
      </div>
      <div v-if="step === 2" class="space-y-4">
        <div class="space-y-1">
          <label class="block text-xs font-semibold text-slate-700">Notaris Rekanan</label>
          <select v-model="form.notary_id" class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm">
            <option value="">— pilih —</option>
            <option v-for="n in store.notaries" :key="n.id" :value="n.id">{{ n.full_name }} · {{ n.wilayah_kerja }}</option>
          </select>
        </div>
      </div>
      <div v-if="step === 3" class="space-y-2 text-sm">
        <p><span class="text-slate-500">Order:</span> <span class="font-mono font-bold">{{ form.order_number || '(otomatis)' }}</span></p>
        <p><span class="text-slate-500">SLA:</span> 14 hari sejak dibuat (status awal <b>pending</b>)</p>
        <p><span class="text-slate-500">PIC:</span> {{ userStore.user?.email }}</p>
      </div>
      <p v-if="failed" class="mt-4 rounded-md bg-red-50 px-3 py-2 text-xs text-red-700">{{ failed }}</p>
      <div class="mt-5 flex justify-between">
        <AppButton variant="secondary" :disabled="step === 1" @click="step -= 1">Kembali</AppButton>
        <AppButton v-if="step < 3" @click="next">Lanjut</AppButton>
        <AppButton v-else :loading="saving" @click="submit">Buat Order</AppButton>
      </div>
    </div>
  </div>
</template>
