<script setup>
const props = defineProps({ tone: { type: String, default: 'neutral' } })

const tones = {
  success: 'badge badge-success',
  warning: 'badge badge-warning',
  danger: 'badge badge-danger',
  info: 'badge badge-info',
  neutral: 'badge badge-neutral',
}
</script>

<template>
  <span class="badge" :class="tones[props.tone] ?? tones.neutral">
    <span class="badge-dot" aria-hidden="true" />
    <slot />
  </span>
</template>
