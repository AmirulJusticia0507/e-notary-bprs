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
  <div class="space-y-1">
    <label class="block text-xs font-semibold text-slate-700">{{ props.label }}</label>
    <input
      type="file"
      :accept="props.accept"
      @change="onChange"
      class="block w-full rounded-md border border-dashed border-slate-300 bg-slate-50 px-3 py-2 text-sm file:mr-3 file:rounded-md file:border-0 file:bg-blue-600 file:px-3 file:py-1.5 file:text-xs file:font-medium file:text-white hover:file:bg-blue-700"
    />
    <p v-if="fileName" class="font-mono text-[11px] text-slate-600">{{ fileName }}</p>
    <p v-if="localError || props.error" class="text-[11px] text-red-600">{{ localError || props.error }}</p>
  </div>
</template>
