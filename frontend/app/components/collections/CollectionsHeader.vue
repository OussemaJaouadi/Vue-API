<script setup lang="ts">
import {
  PhCaretDown,
  PhCheck,
  PhDownloadSimple,
  PhFolderOpen,
  PhPlus,
  PhUploadSimple,
} from '@phosphor-icons/vue'

defineProps<{
  groupCount: number
  requestCount: number
  collectionNames: string[]
  activeCollectionName: string
  exporting?: boolean
}>()

defineEmits<{
  import: []
  export: []
  addCollection: []
  selectCollection: [name: string]
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
      <UiDropdownMenu>
        <UiDropdownMenuTrigger as-child>
          <button class="flex h-9 min-w-40 items-center justify-between gap-3 border-2 border-primary/15 bg-background px-3 font-mono text-[10px] font-black uppercase tracking-widest text-primary transition-all hover:border-primary/50 hover:bg-primary/5">
            <span class="truncate">{{ activeCollectionName === 'all' ? 'All collections' : activeCollectionName }}</span>
            <PhCaretDown class="size-3 opacity-70" />
          </button>
        </UiDropdownMenuTrigger>
        <UiDropdownMenuContent align="end" class="w-64 rounded-none border-2 p-1">
          <UiDropdownMenuLabel class="p-2 text-[9px] font-black uppercase tracking-[0.2em] text-primary/80">
            Active Collection
          </UiDropdownMenuLabel>
          <UiDropdownMenuItem
            class="grid grid-cols-[minmax(0,1fr)_14px] items-center gap-2 rounded-none px-2 py-2.5 font-mono text-[10px] font-black uppercase tracking-widest"
            :class="activeCollectionName === 'all' ? 'bg-primary/8 text-primary' : 'text-muted-foreground'"
            @select="$emit('selectCollection', 'all')"
          >
            <span class="truncate">All collections</span>
            <PhCheck v-if="activeCollectionName === 'all'" class="size-3 text-primary" />
          </UiDropdownMenuItem>
          <UiDropdownMenuItem
            v-for="name in collectionNames"
            :key="name"
            class="grid grid-cols-[minmax(0,1fr)_14px] items-center gap-2 rounded-none px-2 py-2.5 font-mono text-[10px] font-black uppercase tracking-widest"
            :class="name === activeCollectionName ? 'bg-primary/8 text-primary' : 'text-muted-foreground'"
            @select="$emit('selectCollection', name)"
          >
            <span class="truncate">{{ name }}</span>
            <PhCheck v-if="name === activeCollectionName" class="size-3 text-primary" />
          </UiDropdownMenuItem>
        </UiDropdownMenuContent>
      </UiDropdownMenu>

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
          :disabled="exporting"
          :class="exporting && 'cursor-wait opacity-60'"
          type="button"
          @click="$emit('export')"
        >
          <PhDownloadSimple class="size-3.5" />
          <span class="hidden sm:inline">{{ exporting ? 'Exporting' : 'Export' }}</span>
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
