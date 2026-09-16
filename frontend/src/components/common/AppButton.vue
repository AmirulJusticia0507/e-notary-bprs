<script setup>
const props = defineProps({
  variant: { type: String, default: 'primary' },
  type: { type: String, default: 'button' },
  disabled: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
})

const classes = {
  primary: 'primary-button',
  secondary: 'secondary-button',
  danger: 'danger-button',
  ghost: 'ghost-button',
}
</script>

<template>
  <button
    :type="props.type"
    :disabled="props.disabled || props.loading"
    class="inline-flex"
    :class="classes[props.variant] ?? classes.primary"
  >
    <span v-if="props.loading" class="loading-spinner" aria-hidden="true" />
    <span>{{ props.loading ? 'Memproses…' : '' }}</span>
    <slot v-if="!props.loading" />
  </button>
</template>
