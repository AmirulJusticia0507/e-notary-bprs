<script setup>
const props = defineProps({
  page: { type: Number, default: 1 },
  perPage: { type: Number, default: 10 },
  total: { type: Number, default: 0 },
})
const emit = defineEmits(['update:page'])

const totalPages = () => Math.max(1, Math.ceil(props.total / props.perPage))
</script>

<template>
  <div class="flex flex-wrap items-center justify-between gap-3 border-t border-slate-100 px-1 py-3 text-xs text-slate-500">
    <span>Total <b class="font-semibold text-slate-700">{{ props.total }}</b> data · hal. <b class="font-semibold text-slate-700">{{ props.page }}</b>/<b class="font-semibold text-slate-700">{{ totalPages() }}</b></span>
    <div class="flex gap-2">
      <button :disabled="props.page <= 1" aria-label="Halaman sebelumnya" class="grid h-8 w-8 place-items-center rounded-lg border border-slate-200 bg-white text-slate-600 transition-colors hover:border-teal-300 hover:text-teal-700 disabled:cursor-not-allowed disabled:opacity-40" @click="emit('update:page', props.page - 1)">←</button>
      <button :disabled="props.page >= totalPages()" aria-label="Halaman berikutnya" class="grid h-8 w-8 place-items-center rounded-lg border border-slate-200 bg-white text-slate-600 transition-colors hover:border-teal-300 hover:text-teal-700 disabled:cursor-not-allowed disabled:opacity-40" @click="emit('update:page', props.page + 1)">→</button>
    </div>
  </div>
</template>
