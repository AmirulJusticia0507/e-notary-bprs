
# 📊 Table Specifications (Vue 3 + Tailwind)

## Pattern Standard

* High-density layout dengan sticky header.
* Truncate text untuk kolom panjang + tooltip.
* State Badge menggunakan warna dinamis dari Tailwind.

## Sample Vue 3 Component Structure (`LegalOrdersTable.vue`)

```vue
<template>
  <div class="overflow-x-auto rounded-lg border border-slate-200 shadow-sm">
    <table class="w-full text-left text-sm text-slate-700">
      <thead class="bg-slate-900 text-xs uppercase text-slate-200">
        <tr>
          <th class="px-4 py-3">Order ID</th>
          <th class="px-4 py-3">Nasabah</th>
          <th class="px-4 py-3">Notaris</th>
          <th class="px-4 py-3">Status</th>
          <th class="px-4 py-3 text-right">Aksi</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-200 bg-white">
        <tr v-for="order in orders" :key="order.id" class="hover:bg-slate-50">
          <td class="px-4 py-3 font-mono font-bold text-slate-900">{{ order.order_number }}</td>
          <td class="px-4 py-3">{{ order.customer_name }}</td>
          <td class="px-4 py-3">{{ order.notary_name }}</td>
          <td class="px-4 py-3">
            <span :class="statusBadgeClass(order.status)" class="px-2.5 py-1 rounded-full text-xs font-semibold">
              {{ order.status }}
            </span>
          </td>
          <td class="px-4 py-3 text-right">
            <button @click="$emit('view', order.id)" class="text-blue-600 hover:text-blue-800 font-medium">Detail</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
```
