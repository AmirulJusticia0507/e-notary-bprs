<script setup>
import { onMounted, ref, watch } from 'vue'

// CaptchaBox: kode acak 5 karakter di canvas + input verifikasi.
// v-model = true bila jawaban benar. Tombol login wajib menunggu ini.
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
  ctx.fillStyle = '#f1f5f9'
  ctx.fillRect(0, 0, W, H)
  // Garis noise.
  for (let i = 0; i < 4; i++) {
    ctx.strokeStyle = `rgba(100,116,139,${0.3 + Math.random() * 0.4})`
    ctx.lineWidth = 1 + Math.random()
    ctx.beginPath()
    ctx.moveTo(Math.random() * W, Math.random() * H)
    ctx.bezierCurveTo(Math.random() * W, Math.random() * H, Math.random() * W, Math.random() * H, Math.random() * W, Math.random() * H)
    ctx.stroke()
  }
  // Karakter acak dengan rotasi.
  const step = W / (code.value.length + 1)
  ;[...code.value].forEach((ch, i) => {
    const x = step * (i + 1) + (Math.random() * 6 - 3)
    const y = H / 2 + (Math.random() * 8 - 4)
    const angle = (Math.random() * 50 - 25) * (Math.PI / 180)
    ctx.save()
    ctx.translate(x, y)
    ctx.rotate(angle)
    ctx.font = `bold ${22 + Math.floor(Math.random() * 6)}px monospace`
    ctx.fillStyle = `rgb(${30 + Math.floor(Math.random() * 60)},${40 + Math.floor(Math.random() * 60)},${80 + Math.floor(Math.random() * 80)})`
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillText(ch, 0, 0)
    ctx.restore()
  })
  // Titik noise.
  for (let i = 0; i < 40; i++) {
    ctx.fillStyle = `rgba(71,85,105,${0.2 + Math.random() * 0.5})`
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
    <p class="text-xs font-semibold text-slate-600">Captcha <span class="font-normal text-slate-400">(ketik kode di gambar)</span></p>
    <div class="flex items-center gap-2">
      <canvas ref="canvasRef" width="150" height="44" class="rounded-md border border-slate-300" />
      <button
        type="button"
        @click="refresh"
        title="Muat ulang kode"
        class="rounded-md border border-slate-300 px-2.5 py-2 text-sm text-slate-600 hover:bg-slate-100"
      >
        &#10227;
      </button>
    </div>
    <input
      v-model="input"
      type="text"
      maxlength="5"
      autocomplete="off"
      placeholder="Kode captcha"
      class="w-full rounded-md border border-slate-300 px-3 py-2 font-mono text-sm uppercase tracking-widest outline-none focus:border-blue-500"
      :class="solved ? 'border-emerald-500' : ''"
    />
    <p v-if="solved" class="text-[11px] font-medium text-emerald-600">Captcha benar, silakan klik Masuk.</p>
  </div>
</template>
