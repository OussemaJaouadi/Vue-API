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
  <div class="flex flex-1 flex-col overflow-hidden">
    <!-- Tab Bar + Layout Toggle -->
    <div class="flex h-9 items-center justify-between border-b bg-muted/20">
      <WorkbenchTabBar class="h-full border-b-0" />
      
      <UiTooltip>
        <UiTooltipTrigger as-child>
          <button 
            @click="toggleLayout"
            class="mr-2 grid size-6 place-items-center rounded-none text-muted-foreground transition-all hover:bg-primary/10 hover:text-primary"
            :class="workbench.responsePosition.value === 'right' && 'rotate-180'"
          >
            <component :is="workbench.responsePosition.value === 'bottom' ? PhColumns : PhRows" class="size-3.5" />
          </button>
        </UiTooltipTrigger>
        <UiTooltipContent side="bottom">Toggle Response Layout</UiTooltipContent>
      </UiTooltip>
    </div>
    
    <WorkbenchCommandBar />

    <!-- Dynamic Layout Engine -->
    <div 
      class="flex flex-1 min-h-0 overflow-hidden"
      :class="workbench.responsePosition.value === 'bottom' ? 'flex-col' : 'flex-row'"
    >
      <!-- Request Panel -->
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

      <!-- Response Inspector -->
      <section
        class="min-h-35 min-w-75 overflow-hidden bg-card"
        :class="workbench.responsePosition.value === 'bottom' ? 'flex-1' : 'shrink-0'"
        :style="workbench.responsePosition.value === 'right' ? { width: `${workbench.responseWidth.value}px` } : {}"
      >
        <slot name="response" />
      </section>
    </div>
  </div>
</template>
