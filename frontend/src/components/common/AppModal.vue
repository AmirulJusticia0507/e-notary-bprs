<script setup>
const props = defineProps({ show: { type: Boolean, default: false }, title: { type: String, default: '' } })
const emit = defineEmits(['close'])
</script>

<template>
  <div v-if="props.show" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4" @click.self="emit('close')">
    <div class="w-full max-w-lg rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-sm font-bold text-slate-900">{{ props.title }}</h2>
        <button @click="emit('close')" class="text-slate-400 hover:text-slate-600">✕</button>
      </div>
      <slot />
    </div>
  </div>
</template>
