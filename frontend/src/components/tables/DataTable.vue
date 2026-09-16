<script setup>
const props = defineProps({
  columns: { type: Array, required: true },
  rows: { type: Array, default: () => [] },
  rowKey: { type: String, default: 'id' },
})
const emit = defineEmits(['row-click'])
</script>

<template>
  <div class="table-shell">
    <table class="w-full border-collapse text-left text-sm text-slate-600">
      <thead>
        <tr>
          <th v-for="col in props.columns" :key="col.key" class="table-header-cell" :class="col.align === 'right' ? 'text-right' : ''">
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="!props.rows.length">
          <td :colspan="props.columns.length" class="empty-state !min-h-[8rem]">
            <svg width="30" height="30" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M4 6.5h16M4 12h10M4 17.5h7" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"/></svg>
            <span>Belum ada data</span>
          </td>
        </tr>
        <tr v-for="row in props.rows" :key="row[props.rowKey]" class="table-row" :class="{ 'cursor-pointer': $attrs['onRow-click'] }" @click="emit('row-click', row)">
          <td v-for="col in props.columns" :key="col.key" class="max-w-64 truncate px-4 py-3.5" :class="[col.mono ? 'font-mono' : '', col.align === 'right' ? 'text-right' : '']" :title="String(row[col.key] ?? '')">
            <slot :name="`cell-${col.key}`" :row="row" :value="row[col.key]">
              {{ row[col.key] ?? '—' }}
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
