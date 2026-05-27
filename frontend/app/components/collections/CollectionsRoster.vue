<script setup lang="ts">
import { WORKBENCH_ICONS } from '~/composables/useWorkbench'

interface CollectionGroup {
  name: string
  icon: keyof typeof WORKBENCH_ICONS
  requests: any[]
}

defineProps<{
  groups: CollectionGroup[]
  activeCollectionName: string
  totalRequestCount: number
  policyResolver: (name: string) => { defaultEnvironment: string }
}>()

defineEmits<{
  select: [name: string]
}>()
</script>

<template>
  <aside class="flex flex-col w-64 border-r bg-card/30 shrink-0 select-none overflow-hidden">
    <div class="flex h-10 items-center justify-between border-b bg-muted/30 px-3 shrink-0">
      <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">Registry</span>
      <span class="font-mono text-[9px] font-black uppercase tracking-widest text-primary/40">{{ groups.length }} groups</span>
    </div>

    <div class="flex-1 p-1 space-y-0.5 overflow-y-auto custom-scrollbar">
      <button
        class="group relative flex h-12 w-full items-center justify-between px-3 transition-all duration-200 outline-none shrink-0"
        :class="activeCollectionName === 'all' ? 'bg-primary/10 text-foreground' : 'text-muted-foreground/70 hover:bg-primary/5 hover:text-foreground'"
        type="button"
        @click="$emit('select', 'all')"
      >
        <div v-if="activeCollectionName === 'all'" class="wb-active-indicator" />
        <span class="font-mono text-[11px] font-black uppercase tracking-tight" :class="activeCollectionName === 'all' ? 'text-primary' : 'group-hover:text-foreground'">
          All Collections
        </span>
        <span class="font-mono text-[9px] font-black opacity-70 group-hover:opacity-60">{{ totalRequestCount }}</span>
      </button>

      <div class="h-px bg-border/20 my-1 mx-2" />

      <button
        v-for="group in groups"
        :key="group.name"
        class="group relative flex h-12 w-full items-center gap-3 px-3 transition-all duration-200 outline-none shrink-0"
        :class="activeCollectionName === group.name ? 'bg-primary/10 text-foreground' : 'text-muted-foreground/70 hover:bg-primary/5 hover:text-foreground'"
        type="button"
        @click="$emit('select', group.name)"
      >
        <div v-if="activeCollectionName === group.name" class="wb-active-indicator" />
        
        <component 
          :is="WORKBENCH_ICONS[group.icon]" 
          class="size-4 shrink-0 transition-colors" 
          :class="activeCollectionName === group.name ? 'text-primary' : 'text-muted-foreground/80 group-hover:text-primary/60'" 
        />
        
        <span class="min-w-0 flex-1 text-left">
          <span class="block truncate font-mono text-[11px] font-black uppercase tracking-tight transition-colors" :class="activeCollectionName === group.name ? 'text-primary' : 'group-hover:text-foreground'">
            {{ group.name }}
          </span>
          <span class="mt-0.5 flex items-center justify-between gap-2 opacity-80 group-hover:opacity-100 transition-opacity">
            <span class="font-mono text-[8px] font-black uppercase tracking-widest truncate">{{ policyResolver(group.name).defaultEnvironment }}</span>
            <span class="font-mono text-[9px] font-black">{{ group.requests.length }}</span>
          </span>
        </span>
      </button>
    </div>
  </aside>
</template>
