<script setup lang="ts">
import {
  PhCheck,
  PhDotsSixVertical,
  PhX,
} from '@phosphor-icons/vue'
import type { QueryParamItem } from '~/composables/useWorkbench'

const workbench = useWorkbench()
const draggedParamId = ref<string | null>(null)
const dropTarget = ref<{ id: string, position: 'before' | 'after' } | null>(null)

const toggleParam = (param: QueryParamItem) => {
  param.enabled = !param.enabled
}

const isGhostParam = (param: QueryParamItem, index: number) => {
  const isLast = index === workbench.queryParams.value.length - 1
  return isLast && !param.key && !param.value
}

const startDrag = (param: QueryParamItem, index: number, event: DragEvent) => {
  if (isGhostParam(param, index)) return

  draggedParamId.value = param.id
  event.dataTransfer?.setData('text/plain', param.id)
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
  }
}

const updateDropTarget = (param: QueryParamItem, event: DragEvent) => {
  if (!draggedParamId.value || draggedParamId.value === param.id) return

  const row = event.currentTarget as HTMLElement
  const rect = row.getBoundingClientRect()
  const position = event.clientY < rect.top + rect.height / 2 ? 'before' : 'after'
  dropTarget.value = { id: param.id, position }
}

const dropParam = (param: QueryParamItem, event: DragEvent) => {
  const sourceId = draggedParamId.value || event.dataTransfer?.getData('text/plain')
  if (!sourceId || sourceId === param.id) return

  const targetIndex = workbench.queryParams.value.findIndex(item => item.id === param.id)
  if (targetIndex === -1) return

  const insertIndex = dropTarget.value?.position === 'after' ? targetIndex + 1 : targetIndex
  workbench.moveQueryParam(sourceId, insertIndex)
  draggedParamId.value = null
  dropTarget.value = null
}

const endDrag = () => {
  draggedParamId.value = null
  dropTarget.value = null
}

// Watch for changes on rows to auto-add ghost row
watch(workbench.queryParams.value, (newParams) => {
  const lastParam = newParams[newParams.length - 1]
  if (lastParam && (lastParam.key !== '' || lastParam.value !== '')) {
    workbench.addQueryParam()
  }
}, { deep: true })

onMounted(() => {
  // Ensure we start with at least one empty row if none exist
  if (workbench.queryParams.value.length === 0) {
    workbench.addQueryParam()
  }
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col overflow-hidden border-0 bg-card font-mono text-xs">
    <div class="grid shrink-0 grid-cols-[36px_42px_minmax(120px,1fr)_minmax(120px,1.3fr)_36px] border-b bg-muted/20 py-1.5 text-[9px] font-black uppercase tracking-[0.18em] text-muted-foreground/60">
      <div class="flex items-center justify-center border-r border-primary/5">#</div>
      <div class="flex items-center justify-center border-r border-primary/5">On</div>
      <span class="flex items-center border-r border-primary/5 px-3">Key</span>
      <span class="flex items-center border-r border-primary/5 px-3">Value</span>
      <span />
    </div>

    <div class="min-h-0 flex-1 overflow-auto">
      <div
        v-for="(param, index) in workbench.queryParams.value"
        :key="param.id"
        class="group relative grid min-h-9 grid-cols-[36px_42px_minmax(120px,1fr)_minmax(120px,1.3fr)_36px] items-stretch border-b text-[11px] transition-colors"
        :class="[
          param.enabled ? 'bg-background hover:bg-primary/3' : 'bg-muted/10 text-muted-foreground/55 hover:bg-muted/25',
          draggedParamId === param.id && 'opacity-35',
        ]"
        @dragover.prevent="updateDropTarget(param, $event)"
        @dragleave="dropTarget?.id === param.id && (dropTarget = null)"
        @drop.prevent="dropParam(param, $event)"
      >
        <div
          v-if="dropTarget?.id === param.id"
          class="pointer-events-none absolute inset-x-0 z-10 h-0.5 bg-primary shadow-[0_0_10px_rgba(16,185,129,0.45)]"
          :class="dropTarget.position === 'before' ? 'top-0' : 'bottom-0'"
        />

        <div
          class="flex cursor-grab items-center justify-center border-r border-primary/5 text-muted-foreground/20 transition-colors active:cursor-grabbing group-hover:text-primary/40"
          draggable="true"
          @dragstart="startDrag(param, index, $event)"
          @dragend="endDrag"
        >
          <span class="font-mono text-[8px] font-black group-hover:hidden">{{ index + 1 }}</span>
          <PhDotsSixVertical class="hidden size-3.5 group-hover:block" />
        </div>

        <div class="flex items-center justify-center border-r border-primary/5">
          <UiTooltip>
            <UiTooltipTrigger as-child>
              <button
                class="grid size-5 place-items-center border transition-all"
                :class="param.enabled ? 'border-primary bg-primary text-primary-foreground shadow-[0_0_10px_rgba(16,185,129,0.25)]' : 'border-muted-foreground/20 bg-background text-transparent hover:text-muted-foreground'"
                type="button"
                @click="toggleParam(param)"
              >
                <PhCheck class="size-3" />
              </button>
            </UiTooltipTrigger>
            <UiTooltipContent side="top">{{ param.enabled ? 'Disable parameter' : 'Enable parameter' }}</UiTooltipContent>
          </UiTooltip>
        </div>

        <div class="flex items-center border-r border-primary/5 focus-within:bg-background">
          <input
            v-model="param.key"
            class="h-full w-full bg-transparent px-3 font-bold text-foreground/85 outline-none transition-colors placeholder:text-muted-foreground/25 focus:text-primary"
            placeholder="parameter_key"
            spellcheck="false"
          >
        </div>

        <div class="flex items-center border-r border-primary/5 focus-within:bg-background">
          <input
            v-model="param.value"
            class="h-full w-full bg-transparent px-3 text-muted-foreground outline-none transition-colors placeholder:text-muted-foreground/25 focus:text-foreground"
            placeholder="value or {{variable}}"
            spellcheck="false"
          >
        </div>

        <div class="flex items-center justify-center">
          <UiTooltip v-if="index < workbench.queryParams.value.length - 1 || param.key || param.value">
            <UiTooltipTrigger as-child>
              <button
                class="flex size-full items-center justify-center text-muted-foreground/20 transition-all hover:bg-destructive/10 hover:text-destructive active:scale-90"
                type="button"
                @click="workbench.removeQueryParam(param.id)"
              >
                <PhX class="size-3" />
              </button>
            </UiTooltipTrigger>
            <UiTooltipContent side="top">Remove parameter</UiTooltipContent>
          </UiTooltip>
        </div>
      </div>
    </div>
  </div>
</template>
