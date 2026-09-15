<script setup>
const props = defineProps({
  variant: { type: String, default: 'primary' }, // primary | secondary | danger
  type: { type: String, default: 'button' },
  disabled: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
})

const classes = {
  primary: 'bg-blue-600 text-white hover:bg-blue-700 disabled:bg-blue-300',
  secondary: 'border border-slate-300 bg-white text-slate-700 hover:bg-slate-50 disabled:opacity-50',
  danger: 'bg-red-600 text-white hover:bg-red-700 disabled:bg-red-300',
}
</script>

<template>
  <button
    :type="props.type"
    :disabled="props.disabled || props.loading"
    class="rounded-md px-4 py-2 text-sm font-medium transition-colors duration-150"
    :class="classes[props.variant] ?? classes.primary"
  >
    <span v-if="props.loading">Memproses…</span>
    <slot v-else />
  </button>
</template>
