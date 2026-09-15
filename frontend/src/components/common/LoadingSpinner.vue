<script setup>
const props = defineProps({ label: { type: String, default: 'Memuat…' } })
</script>

<template>
  <div class="flex items-center gap-2 py-8 text-sm text-slate-500" role="status">
    <span class="h-4 w-4 animate-spin rounded-full border-2 border-slate-300 border-t-blue-600" />
    {{ props.label }}
  </div>
</template>
