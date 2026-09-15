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

onMounted(() => store.fetchDetail(props.id).catch((e) => (failed.value = e.response?.data?.error ?? 'Gagal memuat detail')))

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
  } catch (e) {
    failed.value = e.response?.data?.error ?? 'Gagal update status'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="space-y-5">
    <button @click="$router.back()" class="text-xs font-medium text-blue-600 hover:text-blue-800">← Kembali</button>
    <LoadingSpinner v-if="store.loading && !store.detail" />
    <template v-else-if="store.detail">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 class="font-mono text-lg font-bold text-slate-900">{{ store.detail.order_number }}</h2>
          <p class="text-xs text-slate-500">{{ store.detail.customer_name }} · {{ store.detail.notary_name }}</p>
        </div>
        <StatusBadge :status="store.detail.status" />
      </div>
      <p v-if="isOverdue(store.detail)" class="rounded-md bg-red-50 px-3 py-2 text-xs font-semibold text-red-700">
        SLA terlewati ({{ new Date(store.detail.sla_deadline).toLocaleDateString('id-ID') }}) — hanya boleh diselesaikan (completed).
      </p>
      <div class="grid gap-4 lg:grid-cols-2">
        <div class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
          <h3 class="mb-3 text-sm font-bold">Info Order</h3>
          <dl class="space-y-2 text-sm">
            <div class="flex justify-between"><dt class="text-slate-500">PIC</dt><dd>{{ store.detail.assigned_name }}</dd></div>
            <div class="flex justify-between"><dt class="text-slate-500">SLA</dt><dd>{{ new Date(store.detail.sla_deadline).toLocaleDateString('id-ID') }}</dd></div>
            <div class="flex justify-between"><dt class="text-slate-500">Selesai</dt><dd>{{ store.detail.completed_at ? new Date(store.detail.completed_at).toLocaleDateString('id-ID') : '—' }}</dd></div>
          </dl>
          <div class="mt-4 flex gap-2">
            <AppButton
              v-for="s in allowedTransitions(store.detail.status)"
              :key="s"
              :variant="s === 'rejected' ? 'danger' : 'primary'"
              @click="openStatusModal(s)"
            >
              → {{ s }}
            </AppButton>
          </div>
        </div>
        <div class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
          <h3 class="mb-3 text-sm font-bold">Audit Trail</h3>
          <ol class="space-y-2 text-xs">
            <li v-if="!store.logs.length" class="text-slate-400">Belum ada perubahan status</li>
            <li v-for="log in store.logs" :key="log.id" class="flex justify-between border-b border-slate-100 pb-2">
              <span class="font-mono">{{ log.prev_status }} → {{ log.new_status }}</span>
              <span class="text-slate-400">{{ new Date(log.created_at).toLocaleString('id-ID') }}</span>
            </li>
          </ol>
        </div>
      </div>
      <div class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
        <h3 class="mb-3 text-sm font-bold">Dokumen ({{ store.documents.length }})</h3>
        <ul class="space-y-2 text-sm">
          <li v-if="!store.documents.length" class="text-xs text-slate-400">Belum ada dokumen</li>
          <li v-for="d in store.documents" :key="d.id" class="flex flex-wrap items-center justify-between gap-2 border-b border-slate-100 pb-2">
            <span class="font-mono text-xs">{{ d.sha256_hash.slice(0, 16) }}…</span>
            <StatusBadge :status="d.e_sign_status" />
          </li>
        </ul>
      </div>
    </template>
    <p v-if="failed" class="rounded-md bg-red-50 px-3 py-2 text-xs text-red-700">{{ failed }}</p>
    <AppModal :show="showStatusModal" title="Ubah status" @close="showStatusModal = false">
      <p class="text-sm">Ubah status menjadi <b class="font-mono">{{ nextStatus }}</b>?</p>
      <div class="mt-4 flex justify-end gap-2">
        <AppButton variant="secondary" @click="showStatusModal = false">Batal</AppButton>
        <AppButton :loading="saving" @click="confirmStatus">Ya, ubah</AppButton>
      </div>
    </AppModal>
  </div>
</template>
