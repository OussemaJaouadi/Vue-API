<script setup lang="ts">
interface EnvironmentItem {
  name: string
  visibility: string
}

defineProps<{
  loading?: boolean
  environments: EnvironmentItem[]
  activeEnvironmentName: string
}>()

defineEmits<{
  select: [name: string]
}>()
</script>

<template>
  <aside class="flex flex-col w-64 border-r bg-card/30 shrink-0 select-none overflow-hidden">
    <div class="flex h-10 items-center justify-between border-b bg-muted/30 px-3 shrink-0">
      <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">Cluster</span>
      <span v-if="!loading" class="font-mono text-[9px] font-black uppercase tracking-widest text-primary/40">{{ environments.length }} nodes</span>
      <div v-else class="h-2 w-12 bg-muted-foreground/15 animate-pulse" />
    </div>

    <div class="flex-1 p-1 space-y-0.5 overflow-y-auto custom-scrollbar">
      <template v-if="loading">
        <div v-for="i in 3" :key="i" class="flex h-14 w-full items-center gap-3 px-3">
          <div class="flex-1 space-y-1.5">
            <div class="h-3 w-24 bg-muted-foreground/15 animate-pulse" />
            <div class="h-2 w-14 bg-muted-foreground/10 animate-pulse" />
          </div>
        </div>
      </template>

      <button
        v-for="environment in environments"
        :key="environment.name"
        class="group relative flex h-14 w-full items-center justify-between px-3 transition-all duration-200 outline-none shrink-0"
        :class="activeEnvironmentName === environment.name ? 'bg-primary/10 text-foreground' : 'text-muted-foreground/70 hover:bg-primary/5 hover:text-foreground'"
        type="button"
        @click="$emit('select', environment.name)"
      >
        <div v-if="activeEnvironmentName === environment.name" class="wb-active-indicator" />
        
        <div class="min-w-0 flex-1 text-left">
          <span class="block truncate font-mono text-[11px] font-black uppercase tracking-tight transition-colors" :class="activeEnvironmentName === environment.name ? 'text-primary' : 'group-hover:text-foreground'">
            {{ environment.name }}
          </span>
          <span 
            class="mt-1 inline-block px-1 py-0.5 border rounded-[2px] font-mono text-[8px] font-black uppercase tracking-tighter transition-all"
            :class="environment.visibility === 'restricted' ? 'border-amber-500/30 bg-amber-500/5 text-amber-600 dark:text-amber-400' : 'border-border/60 text-muted-foreground/80'"
          >
            {{ environment.visibility }}
          </span>
        </div>

        <div v-if="activeEnvironmentName === environment.name" class="size-1.5 bg-primary/40 animate-pulse" />
      </button>

      <div v-if="!loading && environments.length === 0" class="flex flex-col items-center gap-3 px-4 pt-10 text-center">
        <span class="font-mono text-[10px] text-muted-foreground/40 italic">No environments yet</span>
      </div>
    </div>
  </aside>
</template>
