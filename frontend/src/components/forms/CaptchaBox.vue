<script setup>
import { onMounted, ref, watch } from 'vue'

const solved = defineModel({ default: false })
const input = ref('')
const code = ref('')
const canvasRef = ref(null)
const CHARS = 'ABCDEFGHJKMNPQRSTUVWXYZ23456789'

function randomCode() {
  let s = ''
  for (let i = 0; i < 5; i++) s += CHARS[Math.floor(Math.random() * CHARS.length)]
  return s
}

function draw() {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  const W = canvas.width
  const H = canvas.height
  ctx.clearRect(0, 0, W, H)
  const gradient = ctx.createLinearGradient(0, 0, W, H)
  gradient.addColorStop(0, '#f0fdfa')
  gradient.addColorStop(1, '#eff6ff')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, W, H)
  for (let i = 0; i < 4; i++) {
    ctx.strokeStyle = `rgba(13,148,136,${0.18 + Math.random() * 0.25})`
    ctx.lineWidth = 1 + Math.random()
    ctx.beginPath()
    ctx.moveTo(Math.random() * W, Math.random() * H)
    ctx.bezierCurveTo(Math.random() * W, Math.random() * H, Math.random() * W, Math.random() * H, Math.random() * W, Math.random() * H)
    ctx.stroke()
  }
  const step = W / (code.value.length + 1)
  ;[...code.value].forEach((ch, i) => {
    const x = step * (i + 1) + (Math.random() * 6 - 3)
    const y = H / 2 + (Math.random() * 8 - 4)
    const angle = (Math.random() * 50 - 25) * (Math.PI / 180)
    ctx.save()
    ctx.translate(x, y)
    ctx.rotate(angle)
    ctx.font = `bold ${22 + Math.floor(Math.random() * 6)}px monospace`
    ctx.fillStyle = i % 2 ? '#0f766e' : '#1e40af'
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillText(ch, 0, 0)
    ctx.restore()
  })
  for (let i = 0; i < 35; i++) {
    ctx.fillStyle = `rgba(15,118,110,${0.12 + Math.random() * 0.25})`
    ctx.fillRect(Math.random() * W, Math.random() * H, 1.5, 1.5)
  }
}

function refresh() {
  code.value = randomCode()
  input.value = ''
  solved.value = false
  draw()
}

watch(input, (v) => {
  solved.value = v.trim().toUpperCase() === code.value
})

defineExpose({ refresh })
onMounted(refresh)
</script>

<template>
  <div class="space-y-2">
    <p class="field-label">Captcha <span class="font-normal text-slate-400">ketik kode di gambar</span></p>
    <div class="flex items-center gap-2">
      <canvas ref="canvasRef" width="150" height="44" class="rounded-lg border border-slate-200 shadow-sm" />
      <button type="button" title="Muat ulang kode" class="grid h-11 w-11 place-items-center rounded-lg border border-slate-200 bg-white text-teal-700 shadow-sm transition-colors hover:border-teal-300 hover:bg-teal-50" @click="refresh">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M20 11a8 8 0 1 0-2.3 5.7M20 5v6h-6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/></svg>
      </button>
    </div>
    <input v-model="input" type="text" maxlength="5" autocomplete="off" placeholder="Kode captcha" class="field-input font-mono uppercase tracking-[0.2em]" :class="{ 'border-emerald-400': solved }" />
    <p v-if="solved" class="text-[11px] font-medium text-emerald-600">Captcha benar, silakan klik Masuk.</p>
  </div>
</template>
