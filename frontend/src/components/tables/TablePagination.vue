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
  <div class="flex items-center justify-between py-3 text-xs text-slate-500">
    <span>Total {{ props.total }} data · hal. {{ props.page }}/{{ totalPages() }}</span>
    <div class="flex gap-2">
      <button
        :disabled="props.page <= 1"
        @click="emit('update:page', props.page - 1)"
        class="rounded-md border border-slate-300 px-3 py-1.5 font-medium disabled:opacity-40"
      >
        ←
      </button>
      <button
        :disabled="props.page >= totalPages()"
        @click="emit('update:page', props.page + 1)"
        class="rounded-md border border-slate-300 px-3 py-1.5 font-medium disabled:opacity-40"
      >
        →
      </button>
    </div>
  </div>
</template>
