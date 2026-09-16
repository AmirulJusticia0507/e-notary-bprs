<script setup>
const props = defineProps({ label: { type: String, default: 'Memuat…' } })
</script>

<template>
  <div class="loading flex items-center justify-center gap-2.5 py-10 text-sm" role="status">
    <span class="loading-spinner" aria-hidden="true" />
    <span>{{ props.label }}</span>
  </div>
</template>
