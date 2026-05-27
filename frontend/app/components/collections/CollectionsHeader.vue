<script setup lang="ts">
import {
  PhDownloadSimple,
  PhFolderOpen,
  PhPlus,
  PhUploadSimple,
} from '@phosphor-icons/vue'

defineProps<{
  groupCount: number
  requestCount: number
}>()

defineEmits<{
  import: []
  export: []
  addCollection: []
}>()
</script>

<template>
  <header class="flex h-16 items-center justify-between border-b bg-muted/30 px-6 shrink-0 select-none">
    <div class="flex items-center gap-4">
      <div class="grid size-10 place-items-center border-2 border-primary shadow-[3px_3px_0_0_rgba(16,185,129,0.15)] bg-primary/10 text-primary transition-transform hover:scale-105">
        <PhFolderOpen class="size-5" />
      </div>
      <div>
        <h1 class="font-heading text-lg font-black uppercase tracking-tight text-foreground">Collections</h1>
        <p class="font-mono text-[10px] font-black uppercase tracking-widest text-primary/60">
          {{ groupCount }} groups / {{ requestCount }} requests
        </p>
      </div>
    </div>

    <div class="flex items-center gap-3">
      <div class="flex items-center gap-2">
        <button
          class="btn-tactile-muted flex h-9 items-center gap-2 px-3 font-mono text-[10px] font-black uppercase tracking-widest outline-none"
          type="button"
          @click="$emit('import')"
        >
          <PhUploadSimple class="size-3.5" />
          <span class="hidden sm:inline">Import</span>
        </button>
        <button
          class="btn-tactile-muted flex h-9 items-center gap-2 px-3 font-mono text-[10px] font-black uppercase tracking-widest outline-none"
          type="button"
          @click="$emit('export')"
        >
          <PhDownloadSimple class="size-3.5" />
          <span class="hidden sm:inline">Export</span>
        </button>
      </div>

      <div class="h-6 w-px bg-border/40 mx-1" />

      <button
        class="btn-tactile-primary flex h-9 items-center gap-2 px-4 font-mono text-[10px] font-black uppercase tracking-widest outline-none"
        type="button"
        @click="$emit('addCollection')"
      >
        <PhPlus class="size-3.5" />
        <span class="hidden sm:inline">New Collection</span>
        <span class="sm:hidden">New</span>
      </button>
    </div>
  </header>
</template>
