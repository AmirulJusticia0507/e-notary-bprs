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
const pendingCount = computed(() => apps.value.filter((app) => app.status === 'pending').length)

function rupiah(value) { return `Rp ${Number(value ?? 0).toLocaleString('id-ID')}` }
async function fetchMine() { loading.value = true; loadError.value = ''; try { const data = await financingService.listMine(); apps.value = Array.isArray(data) ? data : [] } catch (error) { loadError.value = error.response?.data?.error ?? 'Gagal memuat pengajuan' } finally { loading.value = false } }
async function submitApply() { saveError.value = ''; saveOk.value = ''; clear(); const ok = required(form.customer_name, 'customer_name', 'Nama') & required(form.customer_nik, 'customer_nik', 'NIK') & required(form.collateral_type, 'collateral_type', 'Jenis agunan'); if (!ok || Number(form.financing_amount) <= 0) { if (Number(form.financing_amount) <= 0) errors.financing_amount = 'Plafon harus lebih dari 0'; return } saving.value = true; try { await financingService.apply({ ...form, financing_amount: Number(form.financing_amount) }); saveOk.value = 'Pengajuan terkirim, menunggu verifikasi BPRS'; Object.assign(form, { customer_name: '', customer_nik: '', financing_amount: 0, collateral_type: '', collateral_details: '' }); showForm.value = false; await fetchMine() } catch (error) { saveError.value = error.response?.data?.error ?? 'Pengajuan gagal, coba lagi' } finally { saving.value = false } }
onMounted(fetchMine)
</script>

<template>
  <div class="space-y-6">
    <div class="rounded-2xl border border-slate-200 bg-gradient-to-r from-sky-500 to-teal-600 p-5 text-white shadow-lg shadow-teal-900/10 sm:p-6">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <div class="grid h-12 w-12 place-items-center rounded-2xl border-2 border-white/30 bg-white/15 font-bold text-white">{{ (userStore.user?.full_name || 'N')[0].toUpperCase() }}</div>
          <div><p class="text-xs font-semibold text-sky-50">Halo, {{ userStore.user?.full_name || 'Nasabah' }}</p><h1 class="mt-1 text-lg font-extrabold tracking-tight">Pengajuan Pembiayaan</h1><p class="mt-1 text-xs text-sky-100">Pantau status pengajuan Anda dengan lebih mudah.</p></div>
        </div>
        <AppButton @click="showForm = !showForm" class="!min-h-10 border-white/20 bg-white text-teal-700 hover:bg-sky-50">{{ showForm ? 'Tutup Form' : '+ Pengajuan Baru' }}</AppButton>
      </div>
    </div>
    <div class="grid grid-cols-2 gap-4 sm:grid-cols-2">
      <div class="stat-card p-4" style="--stat-color:#0f766e;--stat-bg:rgba(20,184,166,.11);--stat-glow:rgba(20,184,166,.1)"><p class="text-[11px] font-bold text-slate-500">Total Pengajuan</p><p class="mt-2 text-2xl font-extrabold text-slate-900">{{ apps.length }}</p></div>
      <div class="stat-card p-4" style="--stat-color:#b45309;--stat-bg:rgba(245,158,11,.12);--stat-glow:rgba(245,158,11,.1)"><p class="text-[11px] font-bold text-slate-500">Menunggu Verifikasi</p><p class="mt-2 text-2xl font-extrabold text-slate-900">{{ pendingCount }}</p></div>
    </div>
    <div v-if="showForm" class="panel p-5 sm:p-6"><h2 class="text-base font-extrabold text-slate-900">Form Pengajuan Pembiayaan</h2><p class="mt-1 text-xs text-slate-500">Lengkapi data di bawah ini untuk memulai proses verifikasi.</p><form class="mt-5 grid gap-4 md:grid-cols-2" @submit.prevent="submitApply"><FormInput v-model="form.customer_name" label="Nama lengkap (sesuai KTP)" :error="errors.customer_name" /><FormInput v-model="form.customer_nik" label="NIK" placeholder="16 digit" :error="errors.customer_nik" /><CurrencyInput v-model="form.financing_amount" label="Plafon pengajuan (Rp)" :error="errors.financing_amount" /><FormInput v-model="form.collateral_type" label="Jenis agunan" placeholder="Tanah / Rumah / Kendaraan" :error="errors.collateral_type" /><FormInput v-model="form.collateral_details" label="Detail agunan" placeholder="Alamat / tipe / tahun" class="md:col-span-2" :error="errors.collateral_details" /><div v-if="saveError" class="notice notice-error md:col-span-2">{{ saveError }}</div><div class="md:col-span-2"><AppButton type="submit" :loading="saving" class="w-full">Kirim Pengajuan</AppButton></div></form></div>
    <div v-if="saveOk" class="notice notice-success">{{ saveOk }}</div>
    <section class="panel"><div class="panel-header"><div><h3 class="panel-title">Pengajuan Saya</h3><p class="mt-1 text-[11px] text-slate-400">Riwayat pengajuan pembiayaan Anda</p></div></div><div class="p-2 sm:p-3"><LoadingSpinner v-if="loading && !apps.length" label="Memuat pengajuan…" /><p v-else-if="loadError" class="empty-state">{{ loadError }}</p><p v-else-if="!apps.length" class="empty-state">Belum ada pengajuan. Klik <b>+ Pengajuan Baru</b> untuk mulai.</p><DataTable v-else :columns="[{ key: 'customer_nik', label: 'NIK', mono: true }, { key: 'financing_amount', label: 'Plafon' }, { key: 'collateral_type', label: 'Agunan' }, { key: 'status', label: 'Status' }]" :rows="apps"><template #cell-financing_amount="{ value }">{{ rupiah(value) }}</template><template #cell-status="{ value }"><StatusBadge :status="value" /></template></DataTable></div></section>
  </div>
</template>
