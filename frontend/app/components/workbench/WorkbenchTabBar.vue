<script setup lang="ts">
import {
  PhPlus,
  PhTerminalWindow,
  PhX,
} from '@phosphor-icons/vue'
import { METHOD_BADGE_COLORS } from '~/composables/useWorkbench'

const workbench = useWorkbench()
</script>

<template>
  <nav class="flex h-10 w-full items-center gap-px overflow-hidden border-b bg-muted/20 select-none">
    <div class="flex h-full min-w-0 flex-1 items-center overflow-x-auto overflow-y-hidden custom-scrollbar">
      <div
        v-for="tab in workbench.openTabs.value"
        :key="tab.id"
        class="group relative flex h-full min-w-[120px] max-w-[200px] shrink-0 cursor-pointer items-center justify-between border-r px-3 transition-all duration-200 outline-none"
        :class="workbench.activeRequestId.value === tab.id ? 'bg-background text-primary' : 'bg-muted/10 text-muted-foreground hover:bg-muted/35 hover:text-foreground'"
        @click="workbench.activeRequestId.value = tab.id"
      >
        <div v-if="workbench.activeRequestId.value === tab.id" class="absolute bottom-0 left-0 h-0.75 w-full bg-primary" />
        
        <div class="flex min-w-0 items-center gap-2">
          <div 
            class="flex size-4.5 items-center justify-center border text-[8px] font-black tracking-tighter transition-all"
            :class="METHOD_BADGE_COLORS[tab.method]"
          >
            {{ tab.method.charAt(0) }}
          </div>
          <span class="truncate font-mono text-[10px] font-black uppercase tracking-tight">{{ tab.name }}</span>
        </div>

        <button
          class="ml-2 flex size-5 items-center justify-center opacity-0 transition-all hover:bg-destructive/10 hover:text-destructive group-hover:opacity-100"
          type="button"
          @click.stop="workbench.closeTab(tab.id)"
        >
          <PhX class="size-3" />
        </button>
      </div>

      <button
        class="flex h-full w-10 shrink-0 items-center justify-center text-muted-foreground/70 transition-colors hover:bg-muted/40 hover:text-primary"
        type="button"
        @click="workbench.addRequest()"
      >
        <PhPlus class="size-4" />
      </button>
    </div>

    <!-- Right Side Global Context -->
    <div class="flex h-full shrink-0 items-center gap-2 border-l bg-muted/5 px-3">
      <div class="flex items-center gap-2 px-2 py-1 border-2 border-primary/10 bg-primary/5 transition-all hover:border-primary/30">
        <PhTerminalWindow class="size-3.5 text-primary" />
        <span class="font-mono text-[9px] font-black uppercase tracking-widest text-primary/80">Localhost:8080</span>
      </div>
    </div>
  </nav>
</template>
