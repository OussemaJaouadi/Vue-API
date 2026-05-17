<script setup lang="ts">
import {
  PhPlus,
  PhX,
} from '@phosphor-icons/vue'
import { METHOD_COLORS, METHOD_LABELS, METHOD_STRIP_COLORS } from '~/composables/useWorkbench'

const workbench = useWorkbench()
</script>

<template>
  <div class="flex h-11 items-center border-b-2 border-primary/20 bg-muted/20">
    <div
      v-for="tab in workbench.openTabs.value"
      :key="tab.id"
      class="group relative flex h-full shrink-0 items-center gap-2 border-r border-primary/10 px-3.5 text-[11px] font-mono font-bold uppercase tracking-widest transition-colors"
      :class="[
        tab.id === workbench.activeRequest.value.id
          ? 'bg-background text-foreground'
          : 'text-muted-foreground/70 hover:bg-primary/3 hover:text-foreground',
      ]"
    >
      <div
        v-if="tab.id === workbench.activeRequest.value.id"
        class="absolute inset-x-0 top-0 h-0.5 bg-primary"
      />
      <div class="h-4 w-0.75 rounded-full" :class="METHOD_STRIP_COLORS[tab.method]" />
      <span :class="METHOD_COLORS[tab.method]" class="text-[9px] font-black">{{ METHOD_LABELS[tab.method] }}</span>
      <button
        class="max-w-48 truncate text-left text-foreground/90 outline-none"
        type="button"
        @click="workbench.setActiveRequest(tab)"
      >{{ tab.name }}</button>
      <button
        v-if="workbench.openTabs.value.length > 1"
        class="flex size-4 items-center justify-center opacity-0 transition-opacity hover:bg-destructive/10 hover:text-destructive focus:opacity-100 group-hover:opacity-100"
        type="button"
        @click.stop="workbench.closeTab(tab)"
      >
        <PhX class="size-2.5" />
      </button>
    </div>

    <button class="ml-0.5 flex size-9 items-center justify-center border-l border-primary/20 text-muted-foreground transition-colors hover:bg-primary/5 hover:text-primary" type="button">
      <PhPlus class="size-3.5" />
    </button>
  </div>
</template>
