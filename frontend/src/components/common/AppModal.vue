<script setup>
const props = defineProps({ show: { type: Boolean, default: false }, title: { type: String, default: '' } })
const emit = defineEmits(['close'])
</script>

<template>
  <div v-if="props.show" class="fixed inset-0 z-50 flex items-center justify-center bg-[#0f172a]/55 p-4 backdrop-blur-sm" @click.self="emit('close')">
    <div class="panel w-full max-w-lg overflow-hidden">
      <div class="panel-header">
        <h2 class="panel-title">{{ props.title }}</h2>
        <button type="button" aria-label="Tutup" class="grid h-8 w-8 place-items-center rounded-full text-slate-400 transition-colors hover:bg-slate-100 hover:text-slate-700" @click="emit('close')">
          <svg width="15" height="15" viewBox="0 0 16 16" fill="none" aria-hidden="true"><path d="m3 3 10 10M13 3 3 13" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
        </button>
      </div>
      <div class="panel-body"><slot /></div>
    </div>
  </div>
</template>
