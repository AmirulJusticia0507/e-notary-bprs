<script setup>
const props = defineProps({
  label: { type: String, required: true },
  modelValue: { type: [String, Number], default: '' },
  type: { type: String, default: 'text' },
  placeholder: { type: String, default: '' },
  error: { type: String, default: '' },
  mono: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <div class="space-y-1">
    <label class="block text-xs font-semibold text-slate-700">{{ props.label }}</label>
    <input
      :type="props.type"
      :value="props.modelValue"
      :placeholder="props.placeholder"
      @input="emit('update:modelValue', $event.target.value)"
      class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
      :class="{ 'font-mono': props.mono, 'border-red-400': props.error }"
    />
    <p v-if="props.error" class="text-[11px] text-red-600">{{ props.error }}</p>
  </div>
</template>
