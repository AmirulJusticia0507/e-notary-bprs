<script setup>
import { ref } from 'vue'

// PhotoInput: pilih foto -> downscale via canvas (maks 512px, JPEG 0.8)
// -> v-model berupa data-URL base64 yang siap dikirim sebagai photo_url.
// Hasil典型 30–100KB, aman untuk kolom TEXT dan serverless.
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
  <div class="space-y-2">
    <label class="block text-xs font-semibold text-slate-700">{{ props.label }}</label>
    <div class="flex items-center gap-3">
      <div class="flex h-16 w-16 shrink-0 items-center justify-center overflow-hidden rounded-full border border-slate-300 bg-slate-100">
        <img v-if="preview" :src="preview" alt="Pratinjau foto" class="h-full w-full object-cover" />
        <span v-else class="text-xl text-slate-400">&#128100;</span>
      </div>
      <div class="flex-1">
        <input
          type="file"
          accept="image/jpeg,image/png,image/webp"
          @change="onChange"
          class="block w-full rounded-md border border-dashed border-slate-300 bg-slate-50 px-3 py-2 text-sm file:mr-3 file:rounded-md file:border-0 file:bg-blue-600 file:px-3 file:py-1.5 file:text-xs file:font-medium file:text-white hover:file:bg-blue-700"
        />
        <button v-if="preview" type="button" @click="clear" class="mt-1 text-[11px] text-slate-500 hover:underline">
          Hapus foto
        </button>
      </div>
    </div>
    <p v-if="localError || props.error" class="text-[11px] text-red-600">{{ localError || props.error }}</p>
  </div>
</template>
