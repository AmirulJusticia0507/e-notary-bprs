<script setup>
const props = defineProps({
  columns: { type: Array, required: true }, // [{ key, label, mono?, align? }]
  rows: { type: Array, default: () => [] },
  rowKey: { type: String, default: 'id' },
})
const emit = defineEmits(['row-click'])
</script>

<template>
  <div class="overflow-x-auto rounded-lg border border-slate-200 bg-white shadow-sm">
    <table class="w-full text-left text-sm text-slate-700">
      <thead class="bg-slate-900 text-xs uppercase text-slate-200">
        <tr>
          <th
            v-for="col in props.columns"
            :key="col.key"
            class="whitespace-nowrap px-4 py-3"
            :class="col.align === 'right' ? 'text-right' : ''"
          >
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-200">
        <tr v-if="!props.rows.length">
          <td :colspan="props.columns.length" class="px-4 py-6 text-center text-xs text-slate-400">Belum ada data</td>
        </tr>
        <tr
          v-for="row in props.rows"
          :key="row[props.rowKey]"
          class="hover:bg-slate-50"
          :class="{ 'cursor-pointer': $attrs['onRow-click'] }"
          @click="emit('row-click', row)"
        >
          <td
            v-for="col in props.columns"
            :key="col.key"
            class="max-w-64 truncate px-4 py-3"
            :class="[col.mono ? 'font-mono' : '', col.align === 'right' ? 'text-right' : '']"
            :title="String(row[col.key] ?? '')"
          >
            <slot :name="`cell-${col.key}`" :row="row" :value="row[col.key]">
              {{ row[col.key] ?? '—' }}
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
