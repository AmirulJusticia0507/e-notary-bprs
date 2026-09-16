<script setup>
import { ref } from 'vue'

const props = defineProps({
  label: { type: String, default: 'Dokumen PDF' },
  accept: { type: String, default: '.pdf' },
  maxMb: { type: Number, default: 10 },
  error: { type: String, default: '' },
})
const emit = defineEmits(['update:file'])

const fileName = ref('')
const localError = ref('')

function onChange(e) {
  localError.value = ''
  const file = e.target.files?.[0]
  if (!file) return
  if (file.size > props.maxMb * 1024 * 1024) {
    localError.value = `Ukuran maksimal ${props.maxMb}MB`
    return
  }
  fileName.value = `${file.name} (${(file.size / 1024).toFixed(0)} KB)`
  emit('update:file', file)
}
</script>

<template>
  <div class="space-y-1.5">
    <label class="field-label">{{ props.label }}</label>
    <label class="flex min-h-[5.25rem] cursor-pointer flex-col items-center justify-center gap-1.5 rounded-xl border border-dashed border-slate-300 bg-slate-50 px-4 py-3 text-center transition-colors hover:border-teal-400 hover:bg-teal-50/50">
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 16V4m0 0 4 4m-4-4L8 8" stroke="#0f766e" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><path d="M5 14v4a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-4" stroke="#64748b" stroke-width="1.8" stroke-linecap="round"/></svg>
      <span class="text-xs font-semibold text-slate-600">Pilih file atau tarik ke sini</span>
      <span class="text-[11px] text-slate-400">{{ props.accept }} · maksimal {{ props.maxMb }} MB</span>
      <input type="file" :accept="props.accept" class="sr-only" @change="onChange" />
    </label>
    <p v-if="fileName" class="font-mono text-[11px] text-slate-500">{{ fileName }}</p>
    <p v-if="localError || props.error" class="field-error">{{ localError || props.error }}</p>
  </div>
</template>
