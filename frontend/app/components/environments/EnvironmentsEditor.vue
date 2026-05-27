<script setup lang="ts">
import {
  PhLockKey,
  PhPlus,
  PhStack,
} from '@phosphor-icons/vue'

interface Variable {
  key: string
  value: string
  secret: boolean
}

defineProps<{
  environmentName: string
  variables: Variable[]
}>()

defineEmits<{
  addVariable: []
}>()
</script>

<template>
  <section class="flex-1 flex flex-col border-r bg-card/10 select-none overflow-hidden">
    <div class="flex h-12 items-center justify-between border-b bg-muted/30 px-4 shrink-0">
      <div class="min-w-0">
        <h2 class="truncate font-mono text-[11px] font-black uppercase tracking-tight text-primary">
          {{ environmentName }} Variables
        </h2>
        <p class="font-mono text-[9px] font-bold uppercase tracking-widest text-muted-foreground/90">Storage / Resolution Map</p>
      </div>
      <button
        class="group flex h-8 items-center gap-2 border-2 border-primary/20 bg-primary/5 px-2.5 font-mono text-[9px] font-black uppercase tracking-widest text-primary transition-all hover:border-primary/50 hover:bg-primary/10 hover:shadow-[3px_3px_0_0_rgba(16,185,129,0.1)] active:translate-x-0.5 active:translate-y-0.5"
        type="button"
        @click="$emit('addVariable')"
      >
        <PhPlus class="size-3 group-hover:scale-110 transition-transform" />
        New Key
      </button>
    </div>

    <div class="flex-1 overflow-y-auto custom-scrollbar">
      <table class="w-full border-collapse">
        <thead class="sticky top-0 z-10 bg-muted/50 backdrop-blur-md border-b">
          <tr class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">
            <th class="h-10 px-4 text-left font-black border-r border-border/10 w-[35%]">Key</th>
            <th class="h-10 px-4 text-left font-black border-r border-border/10">Value</th>
            <th class="h-10 px-4 text-right font-black w-[100px]">Authority</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border/20">
          <tr
            v-for="variable in variables"
            :key="variable.key"
            class="group hover:bg-primary/[0.03] transition-colors"
          >
            <td class="h-12 px-4 border-r border-border/5">
              <span class="font-mono text-[11px] font-black uppercase tracking-tight text-foreground group-hover:text-primary transition-colors">
                {{ variable.key }}
              </span>
            </td>
            <td class="h-12 px-4 border-r border-border/5">
              <div class="flex items-center gap-3">
                <PhLockKey v-if="variable.secret" class="size-3.5 text-amber-500/50" />
                <code class="font-mono text-[10px] text-muted-foreground truncate max-w-[400px]">
                  {{ variable.value }}
                </code>
              </div>
            </td>
            <td class="h-12 px-4 text-right">
              <span
                class="font-mono text-[8px] font-black uppercase tracking-widest px-1.5 py-0.5 border-2 transition-all"
                :class="variable.secret 
                  ? 'border-amber-500/20 bg-amber-500/5 text-amber-600 dark:text-amber-400 shadow-[0_0_10px_rgba(245,158,11,0.08)]' 
                  : 'border-border/60 bg-muted/20 text-muted-foreground/80'"
              >
                {{ variable.secret ? 'Secret' : 'Plain' }}
              </span>
            </td>
          </tr>
          <tr v-if="variables.length === 0">
            <td colspan="3" class="h-48 text-center font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground/15 italic">
              No parameters configured for this node
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="h-14 border-t bg-muted/10 p-2 shrink-0">
      <button
        class="flex h-full w-full items-center justify-center gap-3 border-2 border-dashed border-primary/10 bg-background/50 hover:border-primary/40 hover:bg-primary/[0.03] transition-all group"
        type="button"
        @click="$emit('addVariable')"
      >
        <PhPlus class="size-4 text-muted-foreground/80 group-hover:text-primary group-hover:scale-110 transition-all" />
        <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground group-hover:text-primary">Append Parameter</span>
      </button>
    </div>
  </section>
</template>
