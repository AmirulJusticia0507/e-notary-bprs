<script setup>
import AppBadge from '../common/AppBadge.vue'

const props = defineProps({ status: { type: String, default: '' } })

const toneByStatus = {
  pending: 'warning',
  in_progress: 'info',
  completed: 'success',
  approved: 'success',
  signed: 'success',
  rejected: 'danger',
  failed: 'danger',
}
</script>

<template>
  <AppBadge :tone="toneByStatus[props.status] ?? 'neutral'">{{ props.status?.replaceAll('_', ' ') || '—' }}</AppBadge>
</template>
