<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: { type: String, required: true },
  modelValue: { type: [String, Number], default: '' },
  error: { type: String, default: '' },
})
const emit = defineEmits(['update:modelValue'])

function formatRupiah(raw) {
  const digits = String(raw ?? '').replace(/\D/g, '').slice(0, 15)
  if (!digits) return ''
  return new Intl.NumberFormat('id-ID').format(Number(digits))
}

const display = computed(() => formatRupiah(props.modelValue))

function onInput(e) {
  emit('update:modelValue', e.target.value.replace(/\D/g, ''))
}
</script>

<template>
  <div class="space-y-1">
    <label class="block text-xs font-semibold text-slate-700">{{ props.label }}</label>
    <div class="relative">
      <span class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-sm text-slate-400">Rp</span>
      <input
        :value="display"
        inputmode="numeric"
        @input="onInput"
        placeholder="0"
        class="w-full rounded-md border border-slate-300 py-2 pl-9 pr-3 text-right font-mono text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
        :class="{ 'border-red-400': props.error }"
      />
    </div>
    <p v-if="props.error" class="text-[11px] text-red-600">{{ props.error }}</p>
  </div>
</template>
