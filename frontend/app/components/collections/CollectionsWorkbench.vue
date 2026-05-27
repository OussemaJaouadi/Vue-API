<script setup lang="ts">
import {
  PhCaretDown,
  PhCaretRight,
  PhPlus,
  PhTrash,
} from '@phosphor-icons/vue'
import { METHOD_BADGE_COLORS, METHOD_LABELS } from '~/composables/useWorkbench'

interface RequestItem {
  id: string
  method: string
  name: string
  path: string
}

interface CollectionGroup {
  name: string
  requests: RequestItem[]
}

defineProps<{
  activeCollection: CollectionGroup | null
  activeRequestCount: number
  displayedCollections: CollectionGroup[]
  expandedCollections: Record<string, boolean>
}>()

defineEmits<{
  addRequest: [collectionName?: string]
  toggleExpand: [name: string]
  deleteRequest: [id: string]
  deleteCollection: [name: string]
}>()
</script>

<template>
  <section class="flex-1 flex flex-col min-w-0 border-r bg-card/10 select-none overflow-hidden">
    <div class="flex h-12 items-center justify-between border-b bg-muted/30 px-4 shrink-0">
      <div class="min-w-0">
        <h2 class="truncate font-mono text-[11px] font-black uppercase tracking-tight text-primary">
          {{ activeCollection ? activeCollection.name : 'All collections' }}
        </h2>
        <p class="font-mono text-[9px] font-bold uppercase tracking-widest text-muted-foreground/90">Content Browser / Visibility Grid</p>
      </div>
      <button
        class="group flex h-8 items-center gap-2 border-2 border-primary/20 bg-primary/5 px-2.5 font-mono text-[9px] font-black uppercase tracking-widest text-primary transition-all hover:border-primary/50 hover:bg-primary/10 hover:shadow-[3px_3px_0_0_rgba(16,185,129,0.1)] active:translate-x-0.5 active:translate-y-0.5"
        type="button"
        @click="$emit('addRequest', activeCollection?.name)"
      >
        <PhPlus class="size-3 group-hover:scale-110 transition-transform" />
        New Request
      </button>
    </div>

    <div class="flex-1 overflow-y-auto custom-scrollbar p-3 space-y-4">
      <article
        v-for="group in displayedCollections"
        :key="group.name"
        class="border bg-background transition-shadow hover:shadow-sm"
      >
        <div class="flex h-11 items-center justify-between border-b bg-muted/10 px-3">
          <button 
            class="flex items-center gap-2.5 font-mono text-[10px] font-black uppercase tracking-widest text-foreground hover:text-primary transition-colors outline-none"
            @click="$emit('toggleExpand', group.name)"
          >
            <component :is="expandedCollections[group.name] ? PhCaretDown : PhCaretRight" class="size-3.5 transition-transform" :class="expandedCollections[group.name] ? 'text-primary' : 'text-muted-foreground/80 group-hover:text-primary/60'" />
            {{ group.name }}
          </button>
          
          <div class="flex items-center gap-2">
            <span class="font-mono text-[9px] font-black text-muted-foreground/70">{{ group.requests.length }} items</span>
            <button
              class="flex size-6 items-center justify-center text-muted-foreground/50 transition-all hover:bg-destructive/10 hover:text-destructive active:translate-x-0.5 active:translate-y-0.5"
              type="button"
              @click="$emit('deleteCollection', group.name)"
            >
              <PhTrash class="size-3" />
            </button>
          </div>
        </div>

        <div v-if="expandedCollections[group.name]" class="divide-y divide-border/20">
          <div
            v-for="request in group.requests"
            :key="request.id"
            class="group relative flex h-12 items-center gap-4 px-4 transition-all duration-200 hover:bg-primary/[0.03] cursor-pointer"
            @click="navigateTo('/')"
          >
            <div
              class="flex h-5 w-12 items-center justify-center border-2 text-[8px] font-black tracking-tight uppercase transition-all group-hover:shadow-sm"
              :class="METHOD_BADGE_COLORS[request.method as keyof typeof METHOD_BADGE_COLORS] || 'border-muted bg-muted/10 text-muted-foreground'"
            >
              {{ METHOD_LABELS[request.method as keyof typeof METHOD_LABELS] || request.method }}
            </div>
            
            <div class="min-w-0 flex-1">
              <span class="block truncate font-mono text-[11px] font-black uppercase tracking-tight text-foreground transition-colors group-hover:text-foreground">
                {{ request.name }}
              </span>
              <code class="block truncate font-mono text-[9px] text-muted-foreground/80 group-hover:text-muted-foreground/70">{{ request.path }}</code>
            </div>

            <button
              class="opacity-0 group-hover:opacity-100 flex size-8 items-center justify-center text-muted-foreground/80 transition-all hover:bg-destructive/10 hover:text-destructive outline-none"
              type="button"
              @click.stop="$emit('deleteRequest', request.id)"
            >
              <PhTrash class="size-3.5" />
            </button>

            <div class="absolute left-0 h-[40%] w-0.75 bg-primary/0 transition-all group-hover:bg-primary/40 group-hover:h-[60%]" />
          </div>

          <div v-if="group.requests.length === 0" class="flex h-16 items-center justify-center px-3 font-mono text-[9px] uppercase tracking-widest text-muted-foreground/50 italic">
            Collection contains no records
          </div>
        </div>
      </article>
    </div>
  </section>
</template>
