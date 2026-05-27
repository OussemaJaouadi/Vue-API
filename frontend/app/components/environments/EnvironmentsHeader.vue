<script setup lang="ts">
import {
  PhDatabase,
  PhPlus,
} from '@phosphor-icons/vue'

defineProps<{
  activeEnvironmentName: string
}>()

defineEmits<{
  add: []
}>()
</script>

<template>
  <header class="flex h-16 items-center justify-between border-b bg-muted/30 px-6 shrink-0 select-none">
    <div class="flex items-center gap-4">
      <div class="grid size-10 place-items-center border-2 border-primary shadow-[3px_3px_0_0_rgba(16,185,129,0.15)] bg-primary/10 text-primary transition-transform hover:scale-105">
        <PhDatabase class="size-5" />
      </div>
      <div>
        <h1 class="font-heading text-lg font-black uppercase tracking-tight text-foreground">Environments</h1>
        <p class="font-mono text-[10px] font-black uppercase tracking-widest text-primary/60">
          Selected Context: {{ activeEnvironmentName }}
        </p>
      </div>
    </div>

    <button
      class="btn-tactile-primary flex h-9 items-center gap-2 px-4 font-mono text-[10px] font-black uppercase tracking-widest outline-none"
      type="button"
      @click="$emit('add')"
    >
      <PhPlus class="size-3.5" />
      <span class="hidden sm:inline">New Environment</span>
      <span class="sm:hidden">New</span>
    </button>
  </header>
</template>
