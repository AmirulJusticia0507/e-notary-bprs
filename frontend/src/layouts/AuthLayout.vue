<script setup>
// Layout halaman auth: kartu terpusat di atas canvas netral.
</script>

<template>
  <div class="min-h-screen bg-[#f5f8fb]">
    <div class="grid min-h-screen lg:grid-cols-[minmax(0,1.05fr)_minmax(380px,.95fr)]">
      <section class="relative hidden overflow-hidden bg-[#07111f] p-10 lg:flex lg:flex-col lg:justify-between">
        <div class="absolute inset-0 bg-[radial-gradient(circle_at_20%_15%,rgba(45,212,191,.22),transparent_26rem),radial-gradient(circle_at_85%_80%,rgba(59,130,246,.16),transparent_28rem)]" aria-hidden="true" />
        <div class="absolute inset-0 opacity-[.16] [background-image:linear-gradient(rgba(255,255,255,.22)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,.22)_1px,transparent_1px)] [background-size:36px_36px]" aria-hidden="true" />
        <div class="relative flex items-center gap-3">
          <div class="grid h-11 w-11 place-items-center rounded-2xl bg-gradient-to-br from-teal-300 to-sky-500 font-mono text-lg font-bold text-[#06202b] shadow-lg shadow-teal-950/30">N</div>
          <div>
            <p class="text-base font-bold tracking-tight text-white">e-Notary BPRS</p>
            <p class="text-xs text-slate-300">Legal Akad Syariah</p>
          </div>
        </div>

        <div class="relative max-w-xl">
          <div class="mb-6 inline-flex items-center gap-2 rounded-full border border-teal-200/20 bg-white/5 px-3 py-1.5 text-xs font-semibold text-teal-200 backdrop-blur">
            <span class="h-1.5 w-1.5 rounded-full bg-teal-300" />
            Platform legalitas digital terpercaya
          </div>
          <h1 class="text-4xl font-extrabold leading-[1.12] tracking-[-0.045em] text-white sm:text-5xl">
            Akad lebih jelas.<br />
            Proses lebih tenang.
          </h1>
          <p class="mt-5 max-w-lg text-sm leading-6 text-slate-300 sm:text-base">
            Kelola pembiayaan, order notaris, dokumen, dan status akad dalam satu ruang kerja yang rapi, cepat, dan mudah dipantau.
          </p>
          <div class="mt-8 grid gap-3 sm:grid-cols-3">
            <div class="rounded-2xl border border-white/10 bg-white/5 p-4 backdrop-blur">
              <p class="text-2xl font-bold text-teal-200">3–5</p>
              <p class="mt-1 text-[11px] leading-4 text-slate-300">hari target closing</p>
            </div>
            <div class="rounded-2xl border border-white/10 bg-white/5 p-4 backdrop-blur">
              <p class="text-2xl font-bold text-sky-200">14</p>
              <p class="mt-1 text-[11px] leading-4 text-slate-300">hari SLA order</p>
            </div>
            <div class="rounded-2xl border border-white/10 bg-white/5 p-4 backdrop-blur">
              <p class="text-2xl font-bold text-indigo-200">24/7</p>
              <p class="mt-1 text-[11px] leading-4 text-slate-300">akses dashboard</p>
            </div>
          </div>
        </div>

        <div class="relative flex items-center justify-between border-t border-white/10 pt-5 text-xs text-slate-400">
          <span>© 2026 e-Notary BPRS</span>
          <span class="inline-flex items-center gap-2"><span class="h-1.5 w-1.5 rounded-full bg-emerald-300" /> Sistem siap digunakan</span>
        </div>
      </section>

      <section class="relative flex items-center justify-center overflow-hidden p-5 sm:p-10">
        <div class="absolute -left-24 -top-24 h-72 w-72 rounded-full bg-teal-100/50 blur-3xl" aria-hidden="true" />
        <div class="absolute -right-24 -bottom-24 h-72 w-72 rounded-full bg-sky-100/60 blur-3xl" aria-hidden="true" />
        <div class="relative w-full max-w-[440px]">
          <div class="mb-6 flex items-center gap-3 lg:hidden">
            <div class="grid h-10 w-10 place-items-center rounded-xl bg-gradient-to-br from-teal-500 to-sky-600 font-mono text-base font-bold text-white shadow">N</div>
            <div>
              <p class="text-sm font-bold text-slate-900">e-Notary BPRS</p>
              <p class="text-xs text-slate-500">Legal Akad Syariah</p>
            </div>
          </div>
          <div class="auth-card p-6 sm:p-7">
            <div class="mb-6">
              <p class="text-[11px] font-bold uppercase tracking-[0.18em] text-teal-600">Secure workspace</p>
              <h2 class="mt-2 auth-heading">Selamat datang kembali</h2>
              <p class="auth-subheading">Masuk untuk melanjutkan pengelolaan legalitas akad.</p>
            </div>
            <router-view />
          </div>
          <p class="mt-5 text-center text-[11px] text-slate-400">Dilindungi autentikasi JWT dan pembatasan akses berbasis role.</p>
        </div>
      </section>
    </div>
  </div>
</template>
