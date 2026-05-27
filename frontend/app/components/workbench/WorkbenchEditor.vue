<script setup lang="ts">
import {
  PhColumns,
  PhRows,
} from '@phosphor-icons/vue'

const workbench = useWorkbench()

const toggleLayout = () => {
  workbench.responsePosition.value = workbench.responsePosition.value === 'bottom' ? 'right' : 'bottom'
}
</script>

<template>
  <div class="flex flex-1 flex-col overflow-hidden bg-background select-none">
    <!-- Tab Bar + Layout Toggle -->
    <div class="flex h-10 items-center justify-between border-b bg-muted/30">
      <WorkbenchTabBar class="h-full border-b-0" />
      
      <div class="flex items-center gap-1 px-3">
        <UiTooltip>
          <UiTooltipTrigger as-child>
            <button 
              class="group grid size-7 place-items-center border-2 border-primary/10 bg-background transition-all hover:border-primary/40 hover:bg-primary/5 text-muted-foreground hover:text-primary active:scale-95 shadow-sm"
              @click="toggleLayout"
            >
              <component 
                :is="workbench.responsePosition.value === 'bottom' ? PhColumns : PhRows" 
                class="size-3.5 transition-transform" 
                :class="workbench.responsePosition.value === 'right' && 'rotate-180'"
              />
            </button>
          </UiTooltipTrigger>
          <UiTooltipContent side="bottom" class="font-mono text-[9px] font-black uppercase tracking-widest">Toggle Terminal Orientation</UiTooltipContent>
        </UiTooltip>
      </div>
    </div>
    
    <WorkbenchCommandBar />

    <!-- Dynamic Layout Engine -->
    <div 
      class="flex flex-1 min-h-0 overflow-hidden bg-card/5"
      :class="workbench.responsePosition.value === 'bottom' ? 'flex-col' : 'flex-row'"
    >
      <!-- Request Panel (The Workspace) -->
      <WorkbenchRequestPanel
        :style="workbench.responsePosition.value === 'bottom'
          ? { height: `${workbench.editorPaneHeight.value}%` }
          : { flex: '1 1 0%' }"
        class="min-h-45 min-w-75"
      />
      
      <WorkbenchResizer
        v-if="workbench.responsePosition.value === 'bottom'"
        v-model="workbench.editorPaneHeight.value"
        orientation="vertical"
        :min="20"
        :max="80"
      />

      <WorkbenchResizer
        v-else
        v-model="workbench.responseWidth.value"
        orientation="horizontal"
        :reverse="true"
        :min="300"
        :max="800"
      />

      <!-- Response Inspector (The Terminal) -->
      <section
        class="min-h-35 min-w-75 overflow-hidden bg-background"
        :class="workbench.responsePosition.value === 'bottom' ? 'flex-1' : 'shrink-0 shadow-[-4px_0_12px_rgba(0,0,0,0.05)]'"
        :style="workbench.responsePosition.value === 'right' ? { width: `${workbench.responseWidth.value}px` } : {}"
      >
        <slot name="response" />
      </section>
    </div>
  </div>
</template>
