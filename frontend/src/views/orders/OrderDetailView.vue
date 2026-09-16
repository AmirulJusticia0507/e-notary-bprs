<script setup>
import { onMounted, ref } from 'vue'
import AppButton from '../../components/common/AppButton.vue'
import AppModal from '../../components/common/AppModal.vue'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import StatusBadge from '../../components/tables/StatusBadge.vue'
import { useLegalOrder } from '../../composables/useLegalOrder'
import { useUserStore } from '../../stores/userStore'
import { useOrderStore } from '../../stores/orderStore'

const props = defineProps({ id: { type: [String, Number], required: true } })
const store = useOrderStore()
const userStore = useUserStore()
const { allowedTransitions, isOverdue } = useLegalOrder()
const showStatusModal = ref(false)
const nextStatus = ref('')
const failed = ref('')
const saving = ref(false)

onMounted(() => store.fetchDetail(props.id).catch((error) => (failed.value = error.response?.data?.error ?? 'Gagal memuat detail')))

function openStatusModal(status) {
  nextStatus.value = status
  showStatusModal.value = true
}

async function confirmStatus() {
  saving.value = true
  failed.value = ''
  try {
    await store.updateStatus(props.id, nextStatus.value, userStore.user?.id)
    showStatusModal.value = false
  } catch (error) {
    failed.value = error.response?.data?.error ?? 'Gagal update status'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="space-y-6">
    <button type="button" class="secondary-button !min-h-0 !px-0 !py-0 !text-xs !text-teal-700 !shadow-none" @click="$router.back()">← Kembali ke daftar order</button>
    <LoadingSpinner v-if="store.loading && !store.detail" label="Memuat detail order…" />
    <template v-else-if="store.detail">
      <div class="panel p-5 sm:p-6">
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div>
            <p class="font-mono text-xl font-extrabold tracking-tight text-slate-900">{{ store.detail.order_number }}</p>
            <p class="mt-1 text-xs text-slate-500">{{ store.detail.customer_name }} · {{ store.detail.notary_name }}</p>
          </div>
          <StatusBadge :status="store.detail.status" />
        </div>
        <div v-if="isOverdue(store.detail)" class="notice notice-error mt-5">SLA terlewati ({{ new Date(store.detail.sla_deadline).toLocaleDateString('id-ID') }}) — hanya boleh diselesaikan (completed).</div>
      </div>

      <div class="grid gap-5 lg:grid-cols-5">
        <section class="panel p-5 lg:col-span-3">
          <div class="flex items-center justify-between">
            <div><h2 class="text-base font-extrabold text-slate-900">Informasi Order</h2><p class="mt-1 text-xs text-slate-500">Detail penanggung jawab dan batas waktu</p></div>
          </div>
          <dl class="mt-5 grid gap-4 sm:grid-cols-3">
            <div class="rounded-xl bg-slate-50 p-4"><dt class="text-[11px] font-bold text-slate-400">PIC</dt><dd class="mt-1 text-sm font-bold text-slate-800">{{ store.detail.assigned_name }}</dd></div>
            <div class="rounded-xl bg-slate-50 p-4"><dt class="text-[11px] font-bold text-slate-400">SLA</dt><dd class="mt-1 text-sm font-bold text-slate-800">{{ new Date(store.detail.sla_deadline).toLocaleDateString('id-ID') }}</dd></div>
            <div class="rounded-xl bg-slate-50 p-4"><dt class="text-[11px] font-bold text-slate-400">Selesai</dt><dd class="mt-1 text-sm font-bold text-slate-800">{{ store.detail.completed_at ? new Date(store.detail.completed_at).toLocaleDateString('id-ID') : '—' }}</dd></div>
          </dl>
          <div class="mt-6 flex flex-wrap gap-2 border-t border-slate-100 pt-5">
            <AppButton v-for="status in allowedTransitions(store.detail.status)" :key="status" :variant="status === 'rejected' ? 'danger' : 'primary'" @click="openStatusModal(status)">→ {{ status.replaceAll('_', ' ') }}</AppButton>
          </div>
        </section>

        <section class="panel p-5 lg:col-span-2">
          <h2 class="text-base font-extrabold text-slate-900">Audit Trail</h2>
          <p class="mt-1 text-xs text-slate-500">Riwayat perubahan status</p>
          <ol class="mt-4 space-y-3 text-xs">
            <li v-if="!store.logs.length" class="rounded-xl bg-slate-50 p-4 text-center text-slate-400">Belum ada perubahan status</li>
            <li v-for="log in store.logs" :key="log.id" class="flex items-center justify-between gap-3 rounded-xl border border-slate-100 p-3">
              <span class="font-mono font-bold text-slate-700">{{ log.prev_status.replaceAll('_', ' ') }} <span class="text-slate-300">→</span> {{ log.new_status.replaceAll('_', ' ') }}</span>
              <span class="text-slate-400">{{ new Date(log.created_at).toLocaleDateString('id-ID') }}</span>
            </li>
          </ol>
        </section>
      </div>

      <section class="panel">
        <div class="panel-header"><div><h3 class="panel-title">Dokumen</h3><p class="mt-1 text-[11px] text-slate-400">Dokumen terkait order ini</p></div><span class="badge badge-info">{{ store.documents.length }} file</span></div>
        <div class="p-2 sm:p-3">
          <div v-if="!store.documents.length" class="empty-state">Belum ada dokumen</div>
          <ul v-else class="divide-y divide-slate-100">
            <li v-for="document in store.documents" :key="document.id" class="flex flex-wrap items-center justify-between gap-3 px-3 py-3">
              <span class="max-w-full truncate font-mono text-xs text-slate-600">{{ document.sha256_hash.slice(0, 18) }}…</span>
              <StatusBadge :status="document.e_sign_status" />
            </li>
          </ul>
        </div>
      </section>
    </template>
    <div v-if="failed" class="notice notice-error">{{ failed }}</div>
    <AppModal :show="showStatusModal" title="Ubah status order" @close="showStatusModal = false">
      <p class="text-sm leading-6 text-slate-600">Ubah status menjadi <b class="font-mono text-teal-700">{{ nextStatus.replaceAll('_', ' ') }}</b>?</p>
      <div v-if="failed" class="notice notice-error mt-4">{{ failed }}</div>
      <div class="mt-5 flex justify-end gap-2"><AppButton variant="secondary" @click="showStatusModal = false">Batal</AppButton><AppButton :loading="saving" @click="confirmStatus">Ya, ubah</AppButton></div>
    </AppModal>
  </div>
</template>
