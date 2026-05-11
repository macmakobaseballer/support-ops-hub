<script setup lang="ts">
const config = useRuntimeConfig()
</script>

<template>
  <main class="min-h-screen flex flex-col items-center justify-center gap-4 bg-slate-50 text-slate-900">
    <h1 class="text-3xl font-bold">Support Ops Hub</h1>
    <p class="text-sm text-slate-600">M0 — project skeleton</p>
    <p class="text-xs text-slate-400">API base: {{ config.public.apiBase }}</p>
  </main>
</template>
