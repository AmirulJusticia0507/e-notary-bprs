<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import DataTable from '../../components/tables/DataTable.vue'
import StatusBadge from '../../components/tables/StatusBadge.vue'
import FormInput from '../../components/forms/FormInput.vue'
import CurrencyInput from '../../components/forms/CurrencyInput.vue'
import AppButton from '../../components/common/AppButton.vue'
import { useUserStore } from '../../stores/userStore'
import { useFormValidation } from '../../composables/useFormValidation'
import { financingService } from '../../services/financingService'

const userStore = useUserStore()
const { errors, required, clear } = useFormValidation()

const apps = ref([])
const loading = ref(false)
const loadError = ref('')
const showForm = ref(false)
const saving = ref(false)
const saveError = ref('')
const saveOk = ref('')

const form = reactive({ customer_name: '', customer_nik: '', financing_amount: 0, collateral_type: '', collateral_details: '' })

const pendingCount = computed(() => apps.value.filter((a) => a.status === 'pending').length)

function rupiah(n) {
  return `Rp ${Number(n ?? 0).toLocaleString('id-ID')}`
}

async function fetchMine() {
  loading.value = true
  loadError.value = ''
  try {
    const data = await financingService.listMine()
    apps.value = Array.isArray(data) ? data : []
  } catch (e) {
    loadError.value = e.response?.data?.error ?? 'Gagal memuat pengajuan'
  } finally {
    loading.value = false
  }
}

async function submitApply() {
  saveError.value = ''
  saveOk.value = ''
  clear()
  const ok =
    required(form.customer_name, 'customer_name', 'Nama') &
    required(form.customer_nik, 'customer_nik', 'NIK') &
    required(form.collateral_type, 'collateral_type', 'Jenis agunan')
  if (!ok || Number(form.financing_amount) <= 0) {
    if (Number(form.financing_amount) <= 0) errors.financing_amount = 'Plafon harus lebih dari 0'
    return
  }
  saving.value = true
  try {
    await financingService.apply({ ...form, financing_amount: Number(form.financing_amount) })
    saveOk.value = 'Pengajuan terkirim, menunggu verifikasi BPRS'
    Object.assign(form, { customer_name: '', customer_nik: '', financing_amount: 0, collateral_type: '', collateral_details: '' })
    showForm.value = false
    await fetchMine()
  } catch (e) {
    saveError.value = e.response?.data?.error ?? 'Pengajuan gagal, coba lagi'
  } finally {
    saving.value = false
  }
}

onMounted(fetchMine)
</script>

<template>
  <div class="space-y-5">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-lg font-bold text-slate-900">Halo, {{ userStore.user?.full_name ?? 'Nasabah' }}</h2>
        <p class="text-xs text-slate-500">Pantau status pengajuan pembiayaan Anda di sini</p>
      </div>
      <button
        @click="showForm = !showForm"
        class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
      >
        {{ showForm ? 'Tutup' : '+ Pengajuan Baru' }}
      </button>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
        <p class="text-xs font-semibold text-slate-500">Total Pengajuan</p>
        <p class="mt-1 text-2xl font-bold text-slate-900">{{ apps.length }}</p>
      </div>
      <div class="rounded-lg border border-amber-200 bg-white p-5 shadow-sm">
        <p class="text-xs font-semibold text-slate-500">Menunggu Verifikasi</p>
        <p class="mt-1 text-2xl font-bold text-slate-900">{{ pendingCount }}</p>
      </div>
    </div>

    <div v-if="showForm" class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
      <h3 class="mb-3 text-sm font-bold text-slate-900">Form Pengajuan Pembiayaan</h3>
      <form @submit.prevent="submitApply" class="grid gap-4 md:grid-cols-2">
        <FormInput v-model="form.customer_name" label="Nama lengkap (sesuai KTP)" :error="errors.customer_name" />
        <FormInput v-model="form.customer_nik" label="NIK" placeholder="16 digit" :error="errors.customer_nik" />
        <CurrencyInput v-model="form.financing_amount" label="Plafon pengajuan (Rp)" :error="errors.financing_amount" />
        <FormInput v-model="form.collateral_type" label="Jenis agunan" placeholder="Tanah / Rumah / Kendaraan" :error="errors.collateral_type" />
        <FormInput v-model="form.collateral_details" label="Detail agunan" placeholder="Alamat / tipe / tahun" class="md:col-span-2" />
        <p v-if="saveError" class="rounded-md bg-red-50 px-3 py-2 text-xs text-red-700 md:col-span-2">{{ saveError }}</p>
        <div class="md:col-span-2">
          <AppButton type="submit" :loading="saving">Kirim Pengajuan</AppButton>
        </div>
      </form>
    </div>

    <p v-if="saveOk" class="rounded-md bg-emerald-50 px-3 py-2 text-xs text-emerald-700">{{ saveOk }}</p>

    <div class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
      <h3 class="mb-3 text-sm font-bold text-slate-900">Pengajuan Saya</h3>
      <LoadingSpinner v-if="loading && !apps.length" />
      <p v-else-if="loadError" class="text-xs text-red-600">{{ loadError }}</p>
      <p v-else-if="!apps.length" class="text-xs text-slate-500">Belum ada pengajuan. Klik "+ Pengajuan Baru" untuk mulai.</p>
      <DataTable
        v-else
        :columns="[
          { key: 'customer_nik', label: 'NIK', mono: true },
          { key: 'financing_amount', label: 'Plafon' },
          { key: 'collateral_type', label: 'Agunan' },
          { key: 'status', label: 'Status' },
        ]"
        :rows="apps"
      >
        <template #cell-financing_amount="{ value }">{{ rupiah(value) }}</template>
        <template #cell-status="{ value }"><StatusBadge :status="value" /></template>
      </DataTable>
    </div>
  </div>
</template>
