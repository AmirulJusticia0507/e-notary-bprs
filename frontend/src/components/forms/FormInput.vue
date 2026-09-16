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
    <label class="field-label">{{ props.label }}</label>
    <input
      :type="props.type"
      :value="props.modelValue"
      :placeholder="props.placeholder"
      :aria-invalid="Boolean(props.error)"
      @input="emit('update:modelValue', $event.target.value)"
      class="field-input"
      :class="{ 'font-mono': props.mono, 'has-error': props.error }"
    />
    <p v-if="props.error" class="field-error">{{ props.error }}</p>
  </div>
</template>
