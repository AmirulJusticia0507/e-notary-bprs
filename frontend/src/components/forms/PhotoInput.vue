<script setup>
import { ref } from 'vue'

const modelValue = defineModel({ default: '' })
const props = defineProps({
  label: { type: String, default: 'Foto profil' },
  error: { type: String, default: '' },
})

const preview = ref('')
const localError = ref('')

function onChange(e) {
  localError.value = ''
  const file = e.target.files?.[0]
  if (!file) return
  if (!/^image\/(jpeg|png|webp)$/.test(file.type)) {
    localError.value = 'Format harus JPG, PNG, atau WebP'
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    localError.value = 'Ukuran file maksimal 5MB'
    return
  }
  const img = new Image()
  const url = URL.createObjectURL(file)
  img.onload = () => {
    URL.revokeObjectURL(url)
    const MAX = 512
    const scale = Math.min(1, MAX / Math.max(img.width, img.height))
    const w = Math.round(img.width * scale)
    const h = Math.round(img.height * scale)
    const canvas = document.createElement('canvas')
    canvas.width = w
    canvas.height = h
    canvas.getContext('2d').drawImage(img, 0, 0, w, h)
    const dataUrl = canvas.toDataURL('image/jpeg', 0.8)
    preview.value = dataUrl
    modelValue.value = dataUrl
  }
  img.onerror = () => {
    URL.revokeObjectURL(url)
    localError.value = 'File gambar tidak bisa dibaca'
  }
  img.src = url
}

function clear() {
  preview.value = ''
  modelValue.value = ''
}
</script>

<template>
  <div class="space-y-1.5">
    <label class="field-label">{{ props.label }}</label>
    <div class="flex items-center gap-3 rounded-xl border border-slate-200 bg-slate-50/70 p-3">
      <div class="relative flex h-16 w-16 shrink-0 items-center justify-center overflow-hidden rounded-full border-2 border-white bg-gradient-to-br from-teal-100 to-sky-100 shadow">
        <img v-if="preview" :src="preview" alt="Pratinjau foto" class="h-full w-full object-cover" />
        <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle cx="12" cy="8" r="3.2" stroke="#64748b" stroke-width="1.7"/><path d="M5.5 20c.7-3.2 3.1-5 6.5-5s5.8 1.8 6.5 5" stroke="#64748b" stroke-width="1.7" stroke-linecap="round"/></svg>
      </div>
      <div class="flex min-w-0 flex-1 items-center justify-between gap-3">
        <div class="min-w-0">
          <p class="truncate text-xs font-semibold text-slate-700">{{ preview ? 'Foto sudah dipilih' : 'Belum ada foto' }}</p>
          <p class="text-[11px] text-slate-400">JPG, PNG, WebP · maks. 5 MB</p>
        </div>
        <label v-if="!preview" for="photo-upload" class="primary-button !min-h-0 !px-3 !py-2">Pilih</label>
        <button v-else type="button" class="text-[11px] font-semibold text-rose-600 hover:text-rose-700" @click="clear">Hapus</button>
      </div>
      <input v-if="!preview" id="photo-upload" type="file" accept="image/jpeg,image/png,image/webp" class="sr-only" @change="onChange" />
    </div>
    <p v-if="localError || props.error" class="field-error">{{ localError || props.error }}</p>
  </div>
</template>
